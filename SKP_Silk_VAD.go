package silk

// Reviewed by wdvxdr1123 2022-01-03

// SKP_Silk_VAD_Init
/* Initialization of the Silk VAD */
func SKP_Silk_VAD_Init(psSilk_VAD *SKP_Silk_VAD_state) int32 {
	/* reset state memory */
	*psSilk_VAD = SKP_Silk_VAD_state{}

	/* init noise levels */
	/* Initialize array with approx pink noise levels (psd proportional to inverse of frequency) */
	for b := 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NoiseLevelBias[b] = SKP_max_32(int32(VAD_NOISE_LEVELS_BIAS/(int64(b)+1)), 1)
	}

	/* Initialize state */
	for b := 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NL[b] = int32(int64(psSilk_VAD.NoiseLevelBias[b]) * 100)
		psSilk_VAD.Inv_NL[b] = int32(SKP_int32_MAX / int64(psSilk_VAD.NL[b]))
	}
	psSilk_VAD.Counter = 15

	/* init smoothed energy-to-noise ratio*/
	for b := 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NrgRatioSmth_Q8[b] = 100 * 256
	}
	return 0
}

/* Weighting factors for tilt measure */
var tiltWeights = [4]int32{30000, 6000, -12000, -12000}

// SKP_Silk_VAD_GetSA_Q8
/* Get the speech activity level in Q8 */
func SKP_Silk_VAD_GetSA_Q8(
	psSilk_VAD *SKP_Silk_VAD_state, /* I/O  Silk VAD state                  */
	pSA_Q8 *int32, /* O    Speech activity level in Q8     */
	pSNR_dB_Q7 *int32, /* O    SNR for current frame in Q7     */
	pQuality_Q15 [4]int32, /* O    Smoothed SNR for each band      */
	pTilt_Q15 *int32, /* O    current frame's frequency tilt  */
	pIn []int16, /* I    PCM input       [frameLength]   */
	frameLength int32, /* I    Input frame length              */
) int32 {
	var (
		input_tilt, SNR_Q7, sumSquared, smooth_coef_Q16, speech_nrg int32
		X                                                           [4][240]int16
		Xnrg                                                        [4]int32
		NrgToNoiseRatio_Q8                                          [4]int32
	)

	/* Safety checks */
	SKP_assert(VAD_N_BANDS == 4)
	SKP_assert(MAX_FRAME_LENGTH >= frameLength)
	SKP_assert(frameLength <= 512)

	/***********************/
	/* Filter and Decimate */
	/***********************/
	/* 0-8 kHz to 0-4 kHz and 4-8 kHz */
	ana_filt_bank_1(pIn, psSilk_VAD.AnaState[:], X[0][:frameLength], X[3][:frameLength])

	/* 0-4 kHz to 0-2 kHz and 2-4 kHz */
	ana_filt_bank_1(X[0][:frameLength>>1], psSilk_VAD.AnaState1[:], X[0][:frameLength>>1], X[2][:frameLength>>1])

	/* 0-2 kHz to 0-1 kHz and 1-2 kHz */
	ana_filt_bank_1(X[0][:frameLength>>2], psSilk_VAD.AnaState2[:], X[0][:frameLength>>2], X[1][:frameLength>>2])

	/*********************************************/
	/* HP filter on lowest band (differentiator) */
	/*********************************************/
	decimated_framelength := frameLength >> 3
	X[0][decimated_framelength-1] = (X[0][decimated_framelength-1]) >> 1
	HPstateTmp := X[0][decimated_framelength-1]
	for i := decimated_framelength - 1; i > 0; i-- {
		X[0][i-1] = (X[0][i-1]) >> 1
		X[0][i] -= X[0][i-1]
	}
	X[0][0] -= psSilk_VAD.HPstate
	psSilk_VAD.HPstate = HPstateTmp

	/*************************************/
	/* Calculate the energy in each band */
	/*************************************/
	for b := 0; b < VAD_N_BANDS; b++ {
		/* Find the decimated framelength in the non-uniformly divided bands */
		decimated_framelength = frameLength >> SKP_min_32(int32(VAD_N_BANDS-b), VAD_N_BANDS-1)

		/* Split length into subframe lengths */
		dec_subframe_length := decimated_framelength >> VAD_INTERNAL_SUBFRAMES_LOG2
		var dec_subframe_offset int32

		/* Compute energy per sub-frame */
		/* initialize with summed energy of last subframe */
		Xnrg[b] = psSilk_VAD.XnrgSubfr[b]
		for s := 0; s < (1 << VAD_INTERNAL_SUBFRAMES_LOG2); s++ {
			sumSquared = 0
			for i := int32(0); i < dec_subframe_length; i++ {
				/* The energy will be less than dec_subframe_length * ( SKP_int16_MIN / 8 ) ^ 2.            */
				/* Therefore we can accumulate with no risk of overflow (unless dec_subframe_length > 128)  */
				x_tmp := int32(X[b][i+dec_subframe_offset] >> 3)
				sumSquared = SKP_SMLABB(sumSquared, x_tmp, x_tmp)

				/* Safety check */
				SKP_assert(sumSquared >= 0)
			}

			/* Add/saturate summed energy of current subframe */
			if s < (1<<VAD_INTERNAL_SUBFRAMES_LOG2)-1 {
				Xnrg[b] = SKP_ADD_POS_SAT32(Xnrg[b], sumSquared)
			} else {
				/* Look-ahead subframe */
				Xnrg[b] = SKP_ADD_POS_SAT32(Xnrg[b], sumSquared>>1)
			}

			dec_subframe_offset += dec_subframe_length
		}
		psSilk_VAD.XnrgSubfr[b] = sumSquared
	}

	/********************/
	/* Noise estimation */
	/********************/
	psSilk_VAD.getNoiseLevels(Xnrg)

	/***********************************************/
	/* Signal-plus-noise to noise ratio estimation */
	/***********************************************/
	sumSquared = 0
	input_tilt = 0
	for b := 0; b < VAD_N_BANDS; b++ {
		speech_nrg = Xnrg[b] - psSilk_VAD.NL[b]
		if speech_nrg > 0 {
			/* Divide, with sufficient resolution */
			if uint32(Xnrg[b])&0xFF800000 == 0 {
				NrgToNoiseRatio_Q8[b] = (Xnrg[b] << 8) / (psSilk_VAD.NL[b] + 1)
			} else {
				NrgToNoiseRatio_Q8[b] = Xnrg[b] / ((psSilk_VAD.NL[b] >> 8) + 1)
			}

			/* Convert to log domain */
			SNR_Q7 = SKP_Silk_lin2log(NrgToNoiseRatio_Q8[b]) - 8*128

			/* Sum-of-squares */
			sumSquared = SKP_SMLABB(sumSquared, SNR_Q7, SNR_Q7) /* Q14 */

			/* Tilt measure */
			if speech_nrg < (1 << 20) {
				/* Scale down SNR value for small subband speech energies */
				SNR_Q7 = SKP_SMULWB(SKP_Silk_SQRT_APPROX(speech_nrg)<<6, SNR_Q7)
			}
			input_tilt = SKP_SMLAWB(input_tilt, tiltWeights[b], SNR_Q7)
		} else {
			NrgToNoiseRatio_Q8[b] = 256
		}
	}

	/* Mean-of-squares */
	sumSquared = sumSquared / VAD_N_BANDS

	/* Root-mean-square approximation, scale to dBs, and write to output pointer */
	*pSNR_dB_Q7 = int32(int16(SKP_Silk_SQRT_APPROX(sumSquared) * 3))

	/*********************************/
	/* Speech Probability Estimation */
	/*********************************/
	SA_Q15 := SKP_Silk_sigm_Q15(SKP_SMULWB(VAD_SNR_FACTOR_Q16, *pSNR_dB_Q7) - VAD_NEGATIVE_OFFSET_Q5)

	/**************************/
	/* Frequency Tilt Measure */
	/**************************/
	*pTilt_Q15 = (SKP_Silk_sigm_Q15(input_tilt) - 0x4000) << 1

	/**************************************************/
	/* Scale the sigmoid output based on power levels */
	/**************************************************/
	speech_nrg = 0
	for b := 0; b < VAD_N_BANDS; b++ {
		/* Accumulate signal-without-noise energies, higher frequency bands have more weight */
		speech_nrg += int32(b+1) * ((Xnrg[b] - psSilk_VAD.NL[b]) >> 4)
	}

	/* Power scaling */
	if speech_nrg <= 0 {
		SA_Q15 = SA_Q15 >> 1
	} else if speech_nrg < 0x8000 {
		/* square-root */
		speech_nrg = SKP_Silk_SQRT_APPROX(speech_nrg << 15)
		SA_Q15 = SKP_SMULWB(speech_nrg+0x8000, SA_Q15)
	}

	/* Copy the resulting speech activity in Q8 to *pSA_Q8 */
	*pSA_Q8 = SKP_min_32(SA_Q15>>7, SKP_uint8_MAX)

	/***********************************/
	/* Energy Level and SNR estimation */
	/***********************************/
	/* Smoothing coefficient */
	smooth_coef_Q16 = SKP_SMULWB(VAD_SNR_SMOOTH_COEF_Q18, SKP_SMULWB(SA_Q15, SA_Q15))
	for b := 0; b < VAD_N_BANDS; b++ {
		/* compute smoothed energy-to-noise ratio per band */
		psSilk_VAD.NrgRatioSmth_Q8[b] = SKP_SMLAWB(psSilk_VAD.NrgRatioSmth_Q8[b],
			NrgToNoiseRatio_Q8[b]-psSilk_VAD.NrgRatioSmth_Q8[b], smooth_coef_Q16)

		/* signal to noise ratio in dB per band */
		SNR_Q7 = 3 * (SKP_Silk_lin2log(psSilk_VAD.NrgRatioSmth_Q8[b]) - 8*128)

		/* quality = sigmoid( 0.25 * ( SNR_dB - 16 ) ); */
		pQuality_Q15[b] = SKP_Silk_sigm_Q15((SNR_Q7 - 16*128) >> 4)
	}
	return 0
}

// getNoiseLevels Noise level estimation
func (vad *SKP_Silk_VAD_state) getNoiseLevels(pX [4]int32) {
	var coef, min_coef int32
	/* Initially faster smoothing */
	if vad.Counter < 1000 {
		min_coef = SKP_int16_MAX / ((vad.Counter >> 4) + 1)
	} else {
		min_coef = 0
	}

	for k := 0; k < VAD_N_BANDS; k++ {
		/* Get old noise level estimate for current band */
		nl := vad.NL[k]
		SKP_assert(nl >= 0)

		/* Add bias */
		nrg := SKP_ADD_POS_SAT32(pX[k], vad.NoiseLevelBias[k])
		SKP_assert(nrg > 0)

		/* Invert energies */
		inv_nrg := SKP_int32_MAX / nrg
		SKP_assert(inv_nrg >= 0)

		/* Less update when subband energy is high */
		if nrg > (nl << 3) {
			coef = VAD_NOISE_LEVEL_SMOOTH_COEF_Q16 >> 3
		} else if nrg < nl {
			coef = VAD_NOISE_LEVEL_SMOOTH_COEF_Q16
		} else {
			coef = SKP_SMULWB(SKP_SMULWW(inv_nrg, nl), VAD_NOISE_LEVEL_SMOOTH_COEF_Q16<<1)
		}

		/* Initially faster smoothing */
		coef = SKP_max_32(coef, min_coef)

		/* Smooth inverse energies */
		vad.Inv_NL[k] = SKP_SMLAWB(vad.Inv_NL[k], inv_nrg-vad.Inv_NL[k], coef)
		SKP_assert(vad.Inv_NL[k] >= 0)

		/* Compute noise level by inverting again */
		nl = SKP_int32_MAX / vad.Inv_NL[k]
		SKP_assert(nl >= 0)

		/* Limit noise levels (guarantee 7 bits of head room) */
		if nl > 0xFFFFFF {
			nl = 0xFFFFFF
		}
		/* Store as part of state */
		vad.NL[k] = nl
	}
	/* Increment frame counter */
	vad.Counter++
}
