package silk

// Reviewed: 2022-01-02

// Processing of gains

func SKP_Silk_process_gains_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX) {
	psShapeSt := &psEnc.SShape

	/* Gain reduction when LTP coding gain is high */
	if psEncCtrl.SCmn.Sigtype == SIG_TYPE_VOICED {
		/*s = -0.5f * SKP_sigmoid( 0.25f * ( psEncCtrl->LTPredCodGain - 12.0f ) ); */
		s_Q16 := -SKP_Silk_sigm_Q15(SKP_RSHIFT_ROUND(psEncCtrl.LTPredCodGain_Q7-SKP_FIX_CONST(12.0, 7), 4))
		for k := 0; k < NB_SUBFR; k++ {
			psEncCtrl.Gains_Q16[k] = SKP_SMLAWB(psEncCtrl.Gains_Q16[k], psEncCtrl.Gains_Q16[k], s_Q16)
		}
	}

	/* Limit the quantized signal */
	InvMaxSqrVal_Q16 := SKP_Silk_log2lin(SKP_SMULWB(SKP_FIX_CONST(70.0, 7)-
		psEncCtrl.Current_SNR_dB_Q7, SKP_FIX_CONST(0.33, 16))) / psEnc.SCmn.Subfr_length

	for k := 0; k < NB_SUBFR; k++ {
		/* Soft limit on ratio residual energy and squared gains */
		ResNrg := psEncCtrl.ResNrg[k]
		ResNrgPart := SKP_SMULWW(ResNrg, InvMaxSqrVal_Q16)
		if psEncCtrl.ResNrgQ[k] > 0 {
			if psEncCtrl.ResNrgQ[k] < 32 {
				ResNrgPart = SKP_RSHIFT_ROUND(ResNrgPart, psEncCtrl.ResNrgQ[k])
			} else {
				ResNrgPart = 0
			}
		} else if psEncCtrl.ResNrgQ[k] != 0 {
			if ResNrgPart > (SKP_int32_MAX >> (-psEncCtrl.ResNrgQ[k])) {
				ResNrgPart = SKP_int32_MAX
			} else {
				ResNrgPart = ResNrgPart << (-psEncCtrl.ResNrgQ[k])
			}
		}
		gain := psEncCtrl.Gains_Q16[k]
		gain_squared := SKP_ADD_SAT32(ResNrgPart, SKP_SMMUL(gain, gain))
		if gain_squared < SKP_int16_MAX {
			/* recalculate with higher precision */
			gain_squared = SKP_SMLAWW(ResNrgPart<<16, gain, gain)
			gain = SKP_Silk_SQRT_APPROX(gain_squared)          // Q8
			psEncCtrl.Gains_Q16[k] = SKP_LSHIFT_SAT32(gain, 8) // Q16
		} else {
			gain = SKP_Silk_SQRT_APPROX(gain_squared)           // Q0
			psEncCtrl.Gains_Q16[k] = SKP_LSHIFT_SAT32(gain, 16) // Q16
		}
	}

	/* Noise shaping quantization */
	SKP_Silk_gains_quant(psEncCtrl.SCmn.GainsIndices[:], psEncCtrl.Gains_Q16[:],
		&psShapeSt.LastGainIndex, psEnc.SCmn.NFramesInPayloadBuf)
	/* Set quantizer offset for voiced signals. Larger offset when LTP coding gain is low or tilt is high (ie low-pass) */
	if psEncCtrl.SCmn.Sigtype == SIG_TYPE_VOICED {
		if psEncCtrl.LTPredCodGain_Q7+(psEncCtrl.Input_tilt_Q15>>8) > SKP_FIX_CONST(1.0, 7) {
			psEncCtrl.SCmn.QuantOffsetType = 0
		} else {
			psEncCtrl.SCmn.QuantOffsetType = 1
		}
	}

	/* Quantizer boundary adjustment */
	quant_offset_Q10 := int32(SKP_Silk_Quantization_Offsets_Q10[psEncCtrl.SCmn.Sigtype][psEncCtrl.SCmn.QuantOffsetType])
	psEncCtrl.Lambda_Q10 = SKP_FIX_CONST(LAMBDA_OFFSET, 10) +
		SKP_SMULBB(SKP_FIX_CONST(-0.05, 10), psEnc.SCmn.NStatesDelayedDecision) +
		SKP_SMULWB(SKP_FIX_CONST(-0.3, 18), psEnc.Speech_activity_Q8) +
		SKP_SMULWB(SKP_FIX_CONST(-0.2, 12), psEncCtrl.Input_quality_Q14) +
		SKP_SMULWB(SKP_FIX_CONST(-0.1, 12), psEncCtrl.Coding_quality_Q14) +
		SKP_SMULWB(SKP_FIX_CONST(LAMBDA_QUANT_OFFSET, 16), quant_offset_Q10)

	SKP_assert(psEncCtrl.Lambda_Q10 > 0)
	SKP_assert(psEncCtrl.Lambda_Q10 < SKP_FIX_CONST(2, 10))
}
