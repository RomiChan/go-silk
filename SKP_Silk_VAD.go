package silk

import "unsafe"

func SKP_Silk_VAD_Init(psSilk_VAD *SKP_Silk_VAD_state) int32 {
	var (
		b   int32
		ret int32 = 0
	)
	memset(unsafe.Pointer(psSilk_VAD), 0, unsafe.Sizeof(SKP_Silk_VAD_state{}))
	for b = 0; int64(b) < VAD_N_BANDS; b++ {
		psSilk_VAD.NoiseLevelBias[b] = SKP_max_32(int32(VAD_NOISE_LEVELS_BIAS/(int64(b)+1)), 1)
	}
	for b = 0; int64(b) < VAD_N_BANDS; b++ {
		psSilk_VAD.NL[b] = int32(int64(psSilk_VAD.NoiseLevelBias[b]) * 100)
		psSilk_VAD.Inv_NL[b] = int32(SKP_int32_MAX / int64(psSilk_VAD.NL[b]))
	}
	psSilk_VAD.Counter = 15
	for b = 0; int64(b) < VAD_N_BANDS; b++ {
		psSilk_VAD.NrgRatioSmth_Q8[b] = 100 * 256
	}
	return ret
}

var tiltWeights = [4]int32{30000, 6000, -12000, -12000}

func SKP_Silk_VAD_GetSA_Q8(psSilk_VAD *SKP_Silk_VAD_state, pSA_Q8 *int32, pSNR_dB_Q7 *int32, pQuality_Q15 [4]int32, pTilt_Q15 *int32, pIn []int16, framelength int32) int32 {
	var (
		input_tilt            int32
		decimated_framelength int32
		dec_subframe_length   int32
		dec_subframe_offset   int32
		SNR_Q7                int32
		i                     int32
		b                     int32
		s                     int32
		sumSquared            int32
		smooth_coef_Q16       int32
		HPstateTmp            int16
		X                     [4][240]int16
		Xnrg                  [4]int32
		NrgToNoiseRatio_Q8    [4]int32
		speech_nrg            int32
		x_tmp                 int32
		ret                   int32 = 0
	)
	ana_filt_bank_1(unsafe.Slice(&pIn[0], framelength), unsafe.Slice(&psSilk_VAD.AnaState[0], 2), X[0][:framelength], X[3][:framelength])
	ana_filt_bank_1(X[0][:framelength>>1], unsafe.Slice(&psSilk_VAD.AnaState1[0], 2), X[0][:framelength>>1], X[2][:framelength>>1])
	ana_filt_bank_1(X[0][:framelength>>2], unsafe.Slice(&psSilk_VAD.AnaState2[0], 2), X[0][:framelength>>2], X[1][:framelength>>2])
	decimated_framelength = framelength >> 3
	X[0][int64(decimated_framelength)-1] = (X[0][int64(decimated_framelength)-1]) >> 1
	HPstateTmp = X[0][int64(decimated_framelength)-1]
	for i = int32(int64(decimated_framelength) - 1); int64(i) > 0; i-- {
		X[0][int64(i)-1] = (X[0][int64(i)-1]) >> 1
		X[0][i] -= X[0][int64(i)-1]
	}
	X[0][0] -= psSilk_VAD.HPstate
	psSilk_VAD.HPstate = HPstateTmp
	for b = 0; int64(b) < VAD_N_BANDS; b++ {
		decimated_framelength = framelength >> SKP_min_int(int32(VAD_N_BANDS-int64(b)), VAD_N_BANDS-1)
		dec_subframe_length = decimated_framelength >> VAD_INTERNAL_SUBFRAMES_LOG2
		dec_subframe_offset = 0
		Xnrg[b] = psSilk_VAD.XnrgSubfr[b]
		for s = 0; int64(s) < (1 << VAD_INTERNAL_SUBFRAMES_LOG2); s++ {
			sumSquared = 0
			for i = 0; int64(i) < int64(dec_subframe_length); i++ {
				x_tmp = int32((X[b][int64(i)+int64(dec_subframe_offset)]) >> 3)
				sumSquared = SKP_SMLABB(sumSquared, x_tmp, x_tmp)
			}
			if int64(s) < (1<<VAD_INTERNAL_SUBFRAMES_LOG2)-1 {
				Xnrg[b] = SKP_ADD_POS_SAT32(Xnrg[b], sumSquared)
			} else {
				Xnrg[b] = SKP_ADD_POS_SAT32(Xnrg[b], sumSquared>>1)
			}
			dec_subframe_offset += dec_subframe_length
		}
		psSilk_VAD.XnrgSubfr[b] = sumSquared
	}
	psSilk_VAD.getNoiseLevels(Xnrg)
	sumSquared = 0
	input_tilt = 0
	for b = 0; int64(b) < VAD_N_BANDS; b++ {
		speech_nrg = int32(int64(Xnrg[b]) - int64(psSilk_VAD.NL[b]))
		if int64(speech_nrg) > 0 {
			if (int64(Xnrg[b]) & 0xFF800000) == 0 {
				NrgToNoiseRatio_Q8[b] = int32((int64(Xnrg[b]) << 8) / (int64(psSilk_VAD.NL[b]) + 1))
			} else {
				NrgToNoiseRatio_Q8[b] = int32(int64(Xnrg[b]) / ((int64(psSilk_VAD.NL[b]) >> 8) + 1))
			}
			SNR_Q7 = int32(int64(SKP_Silk_lin2log(NrgToNoiseRatio_Q8[b])) - 8*128)
			sumSquared = SKP_SMLABB(sumSquared, SNR_Q7, SNR_Q7)
			if int64(speech_nrg) < (1 << 20) {
				SNR_Q7 = SKP_SMULWB(int32(int64(SKP_Silk_SQRT_APPROX(speech_nrg))<<6), SNR_Q7)
			}
			input_tilt = SKP_SMLAWB(input_tilt, tiltWeights[b], SNR_Q7)
		} else {
			NrgToNoiseRatio_Q8[b] = 256
		}
	}
	sumSquared = int32(int64(sumSquared) / VAD_N_BANDS)
	*pSNR_dB_Q7 = int32(int16(int64(SKP_Silk_SQRT_APPROX(sumSquared)) * 3))
	SA_Q15 := SKP_Silk_sigm_Q15(int32(int64(SKP_SMULWB(VAD_SNR_FACTOR_Q16, *pSNR_dB_Q7)) - VAD_NEGATIVE_OFFSET_Q5))
	*pTilt_Q15 = int32((int64(SKP_Silk_sigm_Q15(input_tilt)) - 0x4000) << 1)
	speech_nrg = 0
	for b = 0; int64(b) < VAD_N_BANDS; b++ {
		speech_nrg += int32((int64(b) + 1) * ((int64(Xnrg[b]) - int64(psSilk_VAD.NL[b])) >> 4))
	}
	if int64(speech_nrg) <= 0 {
		SA_Q15 = SA_Q15 >> 1
	} else if int64(speech_nrg) < 0x8000 {
		speech_nrg = SKP_Silk_SQRT_APPROX(speech_nrg << 15)
		SA_Q15 = SKP_SMULWB(int32(int64(speech_nrg)+0x8000), SA_Q15)
	}
	*pSA_Q8 = SKP_min_int(SA_Q15>>7, SKP_uint8_MAX)
	smooth_coef_Q16 = SKP_SMULWB(VAD_SNR_SMOOTH_COEF_Q18, SKP_SMULWB(SA_Q15, SA_Q15))
	for b = 0; int64(b) < VAD_N_BANDS; b++ {
		psSilk_VAD.NrgRatioSmth_Q8[b] = SKP_SMLAWB(psSilk_VAD.NrgRatioSmth_Q8[b], int32(int64(NrgToNoiseRatio_Q8[b])-int64(psSilk_VAD.NrgRatioSmth_Q8[b])), smooth_coef_Q16)
		SNR_Q7 = int32((int64(SKP_Silk_lin2log(psSilk_VAD.NrgRatioSmth_Q8[b])) - 8*128) * 3)
		pQuality_Q15[b] = SKP_Silk_sigm_Q15(int32((int64(SNR_Q7) - 16*128) >> 4))
	}
	return ret
}

func (vad *SKP_Silk_VAD_state) getNoiseLevels(pX [4]int32) {
	var coef, min_coef int32
	if int64(vad.Counter) < 1000 {
		min_coef = int32(SKP_int16_MAX / ((int64(vad.Counter) >> 4) + 1))
	} else {
		min_coef = 0
	}
	for k := 0; k < VAD_N_BANDS; k++ {
		nl := vad.NL[k]
		nrg := SKP_ADD_POS_SAT32(pX[k], vad.NoiseLevelBias[k])
		inv_nrg := int32(SKP_int32_MAX / int64(nrg))
		if int64(nrg) > (int64(nl) << 3) {
			coef = VAD_NOISE_LEVEL_SMOOTH_COEF_Q16 >> 3
		} else if int64(nrg) < int64(nl) {
			coef = VAD_NOISE_LEVEL_SMOOTH_COEF_Q16
		} else {
			coef = SKP_SMULWB(SKP_SMULWW(inv_nrg, nl), VAD_NOISE_LEVEL_SMOOTH_COEF_Q16<<1)
		}
		coef = SKP_max_int(coef, min_coef)
		vad.Inv_NL[k] = SKP_SMLAWB(vad.Inv_NL[k], int32(int64(inv_nrg)-int64(vad.Inv_NL[k])), coef)
		nl = int32(SKP_int32_MAX / int64(vad.Inv_NL[k]))
		if int64(nl) < 0xFFFFFF {
		} else {
			nl = 0xFFFFFF
		}
		vad.NL[k] = nl
	}
	vad.Counter++
}
