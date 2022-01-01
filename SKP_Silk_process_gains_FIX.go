package silk

import "math"

func SKP_Silk_process_gains_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX) {
	var (
		psShapeSt        *SKP_Silk_shape_state_FIX = &psEnc.SShape
		k                int32
		s_Q16            int32
		InvMaxSqrVal_Q16 int32
		gain             int32
		gain_squared     int32
		ResNrg           int32
		ResNrgPart       int32
		quant_offset_Q10 int32
	)
	if int64(psEncCtrl.SCmn.Sigtype) == SIG_TYPE_VOICED {
		s_Q16 = -SKP_Silk_sigm_Q15(SKP_RSHIFT_ROUND(int32(int64(psEncCtrl.LTPredCodGain_Q7)-int64(SKP_FIX_CONST(12.0, 7))), 4))
		for k = 0; int64(k) < NB_SUBFR; k++ {
			psEncCtrl.Gains_Q16[k] = SKP_SMLAWB(psEncCtrl.Gains_Q16[k], psEncCtrl.Gains_Q16[k], s_Q16)
		}
	}
	InvMaxSqrVal_Q16 = int32(int64(SKP_Silk_log2lin(SKP_SMULWB(int32(int64(SKP_FIX_CONST(70.0, 7))-int64(psEncCtrl.Current_SNR_dB_Q7)), SKP_FIX_CONST(0.33, 16)))) / int64(psEnc.SCmn.Subfr_length))
	for k = 0; int64(k) < NB_SUBFR; k++ {
		ResNrg = psEncCtrl.ResNrg[k]
		ResNrgPart = SKP_SMULWW(ResNrg, InvMaxSqrVal_Q16)
		if int64(psEncCtrl.ResNrgQ[k]) > 0 {
			if int64(psEncCtrl.ResNrgQ[k]) < 32 {
				ResNrgPart = SKP_RSHIFT_ROUND(ResNrgPart, psEncCtrl.ResNrgQ[k])
			} else {
				ResNrgPart = 0
			}
		} else if int64(psEncCtrl.ResNrgQ[k]) != 0 {
			if int64(ResNrgPart) > (SKP_int32_MAX >> int64(-psEncCtrl.ResNrgQ[k])) {
				ResNrgPart = SKP_int32_MAX
			} else {
				ResNrgPart = int32(int64(ResNrgPart) << int64(-psEncCtrl.ResNrgQ[k]))
			}
		}
		gain = psEncCtrl.Gains_Q16[k]
		if ((int64(ResNrgPart) + int64(SKP_SMMUL(gain, gain))) & 0x80000000) == 0 {
			if ((int64(ResNrgPart) & int64(SKP_SMMUL(gain, gain))) & 0x80000000) != 0 {
				gain_squared = math.MinInt32
			} else {
				gain_squared = int32(int64(ResNrgPart) + int64(SKP_SMMUL(gain, gain)))
			}
		} else if ((int64(ResNrgPart) | int64(SKP_SMMUL(gain, gain))) & 0x80000000) == 0 {
			gain_squared = SKP_int32_MAX
		} else {
			gain_squared = int32(int64(ResNrgPart) + int64(SKP_SMMUL(gain, gain)))
		}
		if int64(gain_squared) < SKP_int16_MAX {
			gain_squared = SKP_SMLAWW(int32(int64(ResNrgPart)<<16), gain, gain)
			gain = SKP_Silk_SQRT_APPROX(gain_squared)
			psEncCtrl.Gains_Q16[k] = int32((func() int64 {
				if (int64(math.MinInt32) >> 8) > (SKP_int32_MAX >> 8) {
					if int64(gain) > (int64(math.MinInt32) >> 8) {
						return int64(math.MinInt32) >> 8
					}
					if int64(gain) < (SKP_int32_MAX >> 8) {
						return SKP_int32_MAX >> 8
					}
					return int64(gain)
				}
				if int64(gain) > (SKP_int32_MAX >> 8) {
					return SKP_int32_MAX >> 8
				}
				if int64(gain) < (int64(math.MinInt32) >> 8) {
					return int64(math.MinInt32) >> 8
				}
				return int64(gain)
			}()) << 8)
		} else {
			gain = SKP_Silk_SQRT_APPROX(gain_squared)
			psEncCtrl.Gains_Q16[k] = int32((func() int64 {
				if (int64(math.MinInt32) >> 16) > (SKP_int32_MAX >> 16) {
					if int64(gain) > (int64(math.MinInt32) >> 16) {
						return int64(math.MinInt32) >> 16
					}
					if int64(gain) < (SKP_int32_MAX >> 16) {
						return SKP_int32_MAX >> 16
					}
					return int64(gain)
				}
				if int64(gain) > (SKP_int32_MAX >> 16) {
					return SKP_int32_MAX >> 16
				}
				if int64(gain) < (int64(math.MinInt32) >> 16) {
					return int64(math.MinInt32) >> 16
				}
				return int64(gain)
			}()) << 16)
		}
	}
	SKP_Silk_gains_quant(psEncCtrl.SCmn.GainsIndices[:], psEncCtrl.Gains_Q16[:], &psShapeSt.LastGainIndex, psEnc.SCmn.NFramesInPayloadBuf)
	if int64(psEncCtrl.SCmn.Sigtype) == SIG_TYPE_VOICED {
		if int64(psEncCtrl.LTPredCodGain_Q7)+(int64(psEncCtrl.Input_tilt_Q15)>>8) > int64(SKP_FIX_CONST(1.0, 7)) {
			psEncCtrl.SCmn.QuantOffsetType = 0
		} else {
			psEncCtrl.SCmn.QuantOffsetType = 1
		}
	}
	quant_offset_Q10 = int32(SKP_Silk_Quantization_Offsets_Q10[psEncCtrl.SCmn.Sigtype][psEncCtrl.SCmn.QuantOffsetType])
	psEncCtrl.Lambda_Q10 = int32(int64(SKP_FIX_CONST(LAMBDA_OFFSET, 10)) + int64(SKP_SMULBB(SKP_FIX_CONST(-0.05, 10), psEnc.SCmn.NStatesDelayedDecision)) + int64(SKP_SMULWB(SKP_FIX_CONST(-0.3, 18), psEnc.Speech_activity_Q8)) + int64(SKP_SMULWB(SKP_FIX_CONST(-0.2, 12), psEncCtrl.Input_quality_Q14)) + int64(SKP_SMULWB(SKP_FIX_CONST(-0.1, 12), psEncCtrl.Coding_quality_Q14)) + int64(SKP_SMULWB(SKP_FIX_CONST(LAMBDA_QUANT_OFFSET, 16), quant_offset_Q10)))
}
