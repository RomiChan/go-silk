package silk

import (
	"math"
	"unsafe"
)

func warped_gain(coefs_Q24 *int32, lambda_Q16 int32, order int32) int32 {
	var (
		i        int32
		gain_Q24 int32
	)
	lambda_Q16 = -lambda_Q16
	gain_Q24 = *(*int32)(unsafe.Add(unsafe.Pointer(coefs_Q24), unsafe.Sizeof(int32(0))*uintptr(order-1)))
	for i = order - 2; i >= 0; i-- {
		gain_Q24 = SKP_SMLAWB(*(*int32)(unsafe.Add(unsafe.Pointer(coefs_Q24), unsafe.Sizeof(int32(0))*uintptr(i))), gain_Q24, lambda_Q16)
	}
	gain_Q24 = SKP_SMLAWB(SKP_FIX_CONST(1.0, 24), gain_Q24, -lambda_Q16)
	return SKP_INVERSE32_varQ(gain_Q24, 40)
}
func limit_warped_coefs(coefs_syn_Q24 *int32, coefs_ana_Q24 *int32, lambda_Q16 int32, limit_Q24 int32, order int32) {
	var (
		i            int32
		iter         int32
		ind          int32 = 0
		tmp          int32
		maxabs_Q24   int32
		chirp_Q16    int32
		gain_syn_Q16 int32
		gain_ana_Q16 int32
		nom_Q16      int32
		den_Q24      int32
	)
	lambda_Q16 = -lambda_Q16
	for i = order - 1; i > 0; i-- {
		*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))) = SKP_SMLAWB(*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))), lambda_Q16)
		*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))) = SKP_SMLAWB(*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))), lambda_Q16)
	}
	lambda_Q16 = -lambda_Q16
	nom_Q16 = SKP_SMLAWB(SKP_FIX_CONST(1.0, 16), -lambda_Q16, lambda_Q16)
	den_Q24 = SKP_SMLAWB(SKP_FIX_CONST(1.0, 24), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*0)), lambda_Q16)
	gain_syn_Q16 = SKP_DIV32_varQ(nom_Q16, den_Q24, 24)
	den_Q24 = SKP_SMLAWB(SKP_FIX_CONST(1.0, 24), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*0)), lambda_Q16)
	gain_ana_Q16 = SKP_DIV32_varQ(nom_Q16, den_Q24, 24)
	for i = 0; i < order; i++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_SMULWW(gain_syn_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))
		*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_SMULWW(gain_ana_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))
	}
	for iter = 0; iter < 10; iter++ {
		maxabs_Q24 = -1
		for i = 0; i < order; i++ {
			if (((*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i)))) ^ (*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))>>31) - ((*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i)))) >> 31)) > (((*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i)))) ^ (*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))>>31) - ((*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i)))) >> 31)) {
				tmp = ((*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i)))) ^ (*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))>>31) - ((*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i)))) >> 31)
			} else {
				tmp = ((*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i)))) ^ (*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))>>31) - ((*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i)))) >> 31)
			}
			if tmp > maxabs_Q24 {
				maxabs_Q24 = tmp
				ind = i
			}
		}
		if maxabs_Q24 <= limit_Q24 {
			return
		}
		for i = 1; i < order; i++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))) = SKP_SMLAWB(*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))), lambda_Q16)
			*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))) = SKP_SMLAWB(*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))), lambda_Q16)
		}
		gain_syn_Q16 = SKP_INVERSE32_varQ(gain_syn_Q16, 32)
		gain_ana_Q16 = SKP_INVERSE32_varQ(gain_ana_Q16, 32)
		for i = 0; i < order; i++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_SMULWW(gain_syn_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))
			*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_SMULWW(gain_ana_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))
		}
		chirp_Q16 = SKP_FIX_CONST(0.99, 16) - SKP_DIV32_varQ(SKP_SMULWB(maxabs_Q24-limit_Q24, SKP_SMLABB(SKP_FIX_CONST(0.8, 10), SKP_FIX_CONST(0.1, 10), iter)), maxabs_Q24*(ind+1), 22)
		SKP_Silk_bwexpander_32(coefs_syn_Q24, order, chirp_Q16)
		SKP_Silk_bwexpander_32(coefs_ana_Q24, order, chirp_Q16)
		lambda_Q16 = -lambda_Q16
		for i = order - 1; i > 0; i-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))) = SKP_SMLAWB(*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))), lambda_Q16)
			*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))) = SKP_SMLAWB(*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i-1))), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))), lambda_Q16)
		}
		lambda_Q16 = -lambda_Q16
		nom_Q16 = SKP_SMLAWB(SKP_FIX_CONST(1.0, 16), -lambda_Q16, lambda_Q16)
		den_Q24 = SKP_SMLAWB(SKP_FIX_CONST(1.0, 24), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*0)), lambda_Q16)
		gain_syn_Q16 = SKP_DIV32_varQ(nom_Q16, den_Q24, 24)
		den_Q24 = SKP_SMLAWB(SKP_FIX_CONST(1.0, 24), *(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*0)), lambda_Q16)
		gain_ana_Q16 = SKP_DIV32_varQ(nom_Q16, den_Q24, 24)
		for i = 0; i < order; i++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_SMULWW(gain_syn_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(coefs_syn_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))
			*(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_SMULWW(gain_ana_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(coefs_ana_Q24), unsafe.Sizeof(int32(0))*uintptr(i))))
		}
	}
}
func SKP_Silk_noise_shape_analysis_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX, pitch_res *int16, x *int16) {
	var (
		psShapeSt           *SKP_Silk_shape_state_FIX = &psEnc.SShape
		k                   int32
		i                   int32
		nSamples            int32
		Qnrg                int32
		b_Q14               int32
		warping_Q16         int32
		scale               int32 = 0
		SNR_adj_dB_Q7       int32
		HarmBoost_Q16       int32
		HarmShapeGain_Q16   int32
		Tilt_Q16            int32
		tmp32               int32
		nrg                 int32
		pre_nrg_Q30         int32
		log_energy_Q7       int32
		log_energy_prev_Q7  int32
		energy_variation_Q7 int32
		delta_Q16           int32
		BWExp1_Q16          int32
		BWExp2_Q16          int32
		gain_mult_Q16       int32
		gain_add_Q16        int32
		strength_Q16        int32
		b_Q8                int32
		auto_corr           [17]int32
		refl_coef_Q16       [16]int32
		AR1_Q24             [16]int32
		AR2_Q24             [16]int32
		x_windowed          [360]int16
		x_ptr               *int16
		pitch_res_ptr       *int16
	)
	x_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x), -int(unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.La_shape))))
	psEncCtrl.Current_SNR_dB_Q7 = psEnc.SNR_dB_Q7 - SKP_SMULWB(psEnc.BufferedInChannel_ms<<7, SKP_FIX_CONST(0.05, 16))
	if psEnc.Speech_activity_Q8 > SKP_FIX_CONST(0.5, 8) {
		psEncCtrl.Current_SNR_dB_Q7 -= psEnc.InBandFEC_SNR_comp_Q8 >> 1
	}
	psEncCtrl.Input_quality_Q14 = (psEncCtrl.Input_quality_bands_Q15[0] + psEncCtrl.Input_quality_bands_Q15[1]) >> 2
	psEncCtrl.Coding_quality_Q14 = SKP_Silk_sigm_Q15(SKP_RSHIFT_ROUND(psEncCtrl.Current_SNR_dB_Q7-SKP_FIX_CONST(18.0, 7), 4)) >> 1
	b_Q8 = SKP_FIX_CONST(1.0, 8) - psEnc.Speech_activity_Q8
	b_Q8 = SKP_SMULWB(b_Q8<<8, b_Q8)
	SNR_adj_dB_Q7 = SKP_SMLAWB(psEncCtrl.Current_SNR_dB_Q7, SKP_SMULBB(SKP_FIX_CONST(-4.0, 7)>>(4+1), b_Q8), SKP_SMULWB(SKP_FIX_CONST(1.0, 14)+psEncCtrl.Input_quality_Q14, psEncCtrl.Coding_quality_Q14))
	if psEncCtrl.SCmn.Sigtype == SIG_TYPE_VOICED {
		SNR_adj_dB_Q7 = SKP_SMLAWB(SNR_adj_dB_Q7, SKP_FIX_CONST(2.0, 8), psEnc.LTPCorr_Q15)
	} else {
		SNR_adj_dB_Q7 = SKP_SMLAWB(SNR_adj_dB_Q7, SKP_SMLAWB(SKP_FIX_CONST(6.0, 9), -SKP_FIX_CONST(0.4, 18), psEncCtrl.Current_SNR_dB_Q7), SKP_FIX_CONST(1.0, 14)-psEncCtrl.Input_quality_Q14)
	}
	if psEncCtrl.SCmn.Sigtype == SIG_TYPE_VOICED {
		psEncCtrl.SCmn.QuantOffsetType = 0
		psEncCtrl.Sparseness_Q8 = 0
	} else {
		nSamples = psEnc.SCmn.Fs_kHz << 1
		energy_variation_Q7 = 0
		log_energy_prev_Q7 = 0
		pitch_res_ptr = pitch_res
		for k = 0; k < FRAME_LENGTH_MS/2; k++ {
			SKP_Silk_sum_sqr_shift(&nrg, &scale, pitch_res_ptr, nSamples)
			nrg += nSamples >> scale
			log_energy_Q7 = SKP_Silk_lin2log(nrg)
			if k > 0 {
				energy_variation_Q7 += int32(SKP_abs(int64(log_energy_Q7 - log_energy_prev_Q7)))
			}
			log_energy_prev_Q7 = log_energy_Q7
			pitch_res_ptr = (*int16)(unsafe.Add(unsafe.Pointer(pitch_res_ptr), unsafe.Sizeof(int16(0))*uintptr(nSamples)))
		}
		psEncCtrl.Sparseness_Q8 = SKP_Silk_sigm_Q15(SKP_SMULWB(energy_variation_Q7-SKP_FIX_CONST(5.0, 7), SKP_FIX_CONST(0.1, 16))) >> 7
		if psEncCtrl.Sparseness_Q8 > SKP_FIX_CONST(0.75, 8) {
			psEncCtrl.SCmn.QuantOffsetType = 0
		} else {
			psEncCtrl.SCmn.QuantOffsetType = 1
		}
		SNR_adj_dB_Q7 = SKP_SMLAWB(SNR_adj_dB_Q7, SKP_FIX_CONST(2.0, 15), psEncCtrl.Sparseness_Q8-SKP_FIX_CONST(0.5, 8))
	}
	strength_Q16 = SKP_SMULWB(psEncCtrl.PredGain_Q16, SKP_FIX_CONST(0.001, 16))
	BWExp1_Q16 = func() int32 {
		BWExp2_Q16 = SKP_DIV32_varQ(SKP_FIX_CONST(0.95, 16), SKP_SMLAWW(SKP_FIX_CONST(1.0, 16), strength_Q16, strength_Q16), 16)
		return BWExp2_Q16
	}()
	delta_Q16 = SKP_SMULWB(SKP_FIX_CONST(1.0, 16)-SKP_SMULBB(3, psEncCtrl.Coding_quality_Q14), SKP_FIX_CONST(0.01, 16))
	BWExp1_Q16 = BWExp1_Q16 - delta_Q16
	BWExp2_Q16 = BWExp2_Q16 + delta_Q16
	BWExp1_Q16 = (BWExp1_Q16 << 14) / (BWExp2_Q16 >> 2)
	if psEnc.SCmn.Warping_Q16 > 0 {
		warping_Q16 = SKP_SMLAWB(psEnc.SCmn.Warping_Q16, psEncCtrl.Coding_quality_Q14, SKP_FIX_CONST(0.01, 18))
	} else {
		warping_Q16 = 0
	}
	for k = 0; k < NB_SUBFR; k++ {
		var (
			shift      int32
			slope_part int32
			flat_part  int32
		)
		flat_part = psEnc.SCmn.Fs_kHz * 5
		slope_part = (psEnc.SCmn.ShapeWinLength - flat_part) >> 1
		SKP_Silk_apply_sine_window(x_windowed[:], ([]int16)(x_ptr), 1, slope_part)
		shift = slope_part
		memcpy(unsafe.Pointer(&x_windowed[shift]), unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(shift)))), size_t(uintptr(flat_part)*unsafe.Sizeof(int16(0))))
		shift += flat_part
		SKP_Silk_apply_sine_window(([]int16)(&x_windowed[shift]), ([]int16)((*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(shift)))), 2, slope_part)
		x_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.Subfr_length)))
		if psEnc.SCmn.Warping_Q16 > 0 {
			SKP_Silk_warped_autocorrelation_FIX(&auto_corr[0], &scale, &x_windowed[0], int16(warping_Q16), psEnc.SCmn.ShapeWinLength, psEnc.SCmn.ShapingLPCOrder)
		} else {
			SKP_Silk_autocorr(auto_corr[:], &scale, x_windowed[:], psEnc.SCmn.ShapeWinLength, psEnc.SCmn.ShapingLPCOrder+1)
		}
		auto_corr[0] = (auto_corr[0]) + SKP_max_32(SKP_SMULWB((auto_corr[0])>>4, SKP_FIX_CONST(1e-05, 20)), 1)
		nrg = SKP_Silk_schur64(refl_coef_Q16[:], auto_corr[:], psEnc.SCmn.ShapingLPCOrder)
		SKP_Silk_k2a_Q16(AR2_Q24[:], refl_coef_Q16[:], psEnc.SCmn.ShapingLPCOrder)
		Qnrg = -scale
		if Qnrg&1 != 0 {
			Qnrg -= 1
			nrg >>= 1
		}
		tmp32 = SKP_Silk_SQRT_APPROX(nrg)
		Qnrg >>= 1
		psEncCtrl.Gains_Q16[k] = SKP_LSHIFT_SAT32(tmp32, 16-Qnrg)
		if psEnc.SCmn.Warping_Q16 > 0 {
			gain_mult_Q16 = warped_gain(&AR2_Q24[0], warping_Q16, psEnc.SCmn.ShapingLPCOrder)
			psEncCtrl.Gains_Q16[k] = SKP_SMULWW(psEncCtrl.Gains_Q16[k], gain_mult_Q16)
			if psEncCtrl.Gains_Q16[k] < 0 {
				psEncCtrl.Gains_Q16[k] = SKP_int32_MAX
			}
		}
		SKP_Silk_bwexpander_32(&AR2_Q24[0], psEnc.SCmn.ShapingLPCOrder, BWExp2_Q16)
		memcpy(unsafe.Pointer(&AR1_Q24[0]), unsafe.Pointer(&AR2_Q24[0]), size_t(uintptr(psEnc.SCmn.ShapingLPCOrder)*unsafe.Sizeof(int32(0))))
		SKP_Silk_bwexpander_32(&AR1_Q24[0], psEnc.SCmn.ShapingLPCOrder, BWExp1_Q16)
		SKP_Silk_LPC_inverse_pred_gain_Q24(&pre_nrg_Q30, AR2_Q24[:], psEnc.SCmn.ShapingLPCOrder)
		SKP_Silk_LPC_inverse_pred_gain_Q24(&nrg, AR1_Q24[:], psEnc.SCmn.ShapingLPCOrder)
		pre_nrg_Q30 = SKP_SMULWB(pre_nrg_Q30, SKP_FIX_CONST(0.7, 15)) << 1
		psEncCtrl.GainsPre_Q14[k] = SKP_FIX_CONST(0.3, 14) + SKP_DIV32_varQ(pre_nrg_Q30, nrg, 14)
		limit_warped_coefs(&AR2_Q24[0], &AR1_Q24[0], warping_Q16, SKP_FIX_CONST(3.999, 24), psEnc.SCmn.ShapingLPCOrder)
		for i = 0; i < psEnc.SCmn.ShapingLPCOrder; i++ {
			psEncCtrl.AR1_Q13[k*MAX_SHAPE_LPC_ORDER+i] = SKP_SAT16(SKP_RSHIFT_ROUND(AR1_Q24[i], 11))
			psEncCtrl.AR2_Q13[k*MAX_SHAPE_LPC_ORDER+i] = SKP_SAT16(SKP_RSHIFT_ROUND(AR2_Q24[i], 11))
		}
	}
	gain_mult_Q16 = SKP_Silk_log2lin(-SKP_SMLAWB(-SKP_FIX_CONST(16.0, 7), SNR_adj_dB_Q7, SKP_FIX_CONST(0.16, 16)))
	gain_add_Q16 = SKP_Silk_log2lin(SKP_SMLAWB(SKP_FIX_CONST(16.0, 7), SKP_FIX_CONST(NOISE_FLOOR_dB, 7), SKP_FIX_CONST(0.16, 16)))
	tmp32 = SKP_Silk_log2lin(SKP_SMLAWB(SKP_FIX_CONST(16.0, 7), SKP_FIX_CONST(-50.0, 7), SKP_FIX_CONST(0.16, 16)))
	tmp32 = SKP_SMULWW(psEnc.AvgGain_Q16, tmp32)
	if ((gain_add_Q16 + tmp32) & math.MinInt32) == 0 {
		if ((gain_add_Q16 & tmp32) & math.MinInt32) != 0 {
			gain_add_Q16 = math.MinInt32
		} else {
			gain_add_Q16 = gain_add_Q16 + tmp32
		}
	} else if ((gain_add_Q16 | tmp32) & math.MinInt32) == 0 {
		gain_add_Q16 = SKP_int32_MAX
	} else {
		gain_add_Q16 = gain_add_Q16 + tmp32
	}
	for k = 0; k < NB_SUBFR; k++ {
		psEncCtrl.Gains_Q16[k] = SKP_SMULWW(psEncCtrl.Gains_Q16[k], gain_mult_Q16)
		if psEncCtrl.Gains_Q16[k] < 0 {
			psEncCtrl.Gains_Q16[k] = SKP_int32_MAX
		}
	}
	for k = 0; k < NB_SUBFR; k++ {
		psEncCtrl.Gains_Q16[k] = SKP_ADD_POS_SAT32(psEncCtrl.Gains_Q16[k], gain_add_Q16)
		if ((psEnc.AvgGain_Q16 + SKP_SMULWB(psEncCtrl.Gains_Q16[k]-psEnc.AvgGain_Q16, SKP_RSHIFT_ROUND(SKP_SMULBB(psEnc.Speech_activity_Q8, SKP_FIX_CONST(GAIN_SMOOTHING_COEF, 10)), 2))) & math.MinInt32) == 0 {
			if ((psEnc.AvgGain_Q16 & SKP_SMULWB(psEncCtrl.Gains_Q16[k]-psEnc.AvgGain_Q16, SKP_RSHIFT_ROUND(SKP_SMULBB(psEnc.Speech_activity_Q8, SKP_FIX_CONST(GAIN_SMOOTHING_COEF, 10)), 2))) & math.MinInt32) != 0 {
				psEnc.AvgGain_Q16 = math.MinInt32
			} else {
				psEnc.AvgGain_Q16 = psEnc.AvgGain_Q16 + SKP_SMULWB(psEncCtrl.Gains_Q16[k]-psEnc.AvgGain_Q16, SKP_RSHIFT_ROUND(SKP_SMULBB(psEnc.Speech_activity_Q8, SKP_FIX_CONST(GAIN_SMOOTHING_COEF, 10)), 2))
			}
		} else if ((psEnc.AvgGain_Q16 | SKP_SMULWB(psEncCtrl.Gains_Q16[k]-psEnc.AvgGain_Q16, SKP_RSHIFT_ROUND(SKP_SMULBB(psEnc.Speech_activity_Q8, SKP_FIX_CONST(GAIN_SMOOTHING_COEF, 10)), 2))) & math.MinInt32) == 0 {
			psEnc.AvgGain_Q16 = SKP_int32_MAX
		} else {
			psEnc.AvgGain_Q16 = psEnc.AvgGain_Q16 + SKP_SMULWB(psEncCtrl.Gains_Q16[k]-psEnc.AvgGain_Q16, SKP_RSHIFT_ROUND(SKP_SMULBB(psEnc.Speech_activity_Q8, SKP_FIX_CONST(GAIN_SMOOTHING_COEF, 10)), 2))
		}
	}
	gain_mult_Q16 = SKP_FIX_CONST(1.0, 16) + SKP_RSHIFT_ROUND(SKP_FIX_CONST(INPUT_TILT, 26)+psEncCtrl.Coding_quality_Q14*SKP_FIX_CONST(HIGH_RATE_INPUT_TILT, 12), 10)
	if psEncCtrl.Input_tilt_Q15 <= 0 && psEncCtrl.SCmn.Sigtype == SIG_TYPE_UNVOICED {
		if psEnc.SCmn.Fs_kHz == 24 {
			var essStrength_Q15 int32 = SKP_SMULWW(-psEncCtrl.Input_tilt_Q15, SKP_SMULBB(psEnc.Speech_activity_Q8, SKP_FIX_CONST(1.0, 8)-psEncCtrl.Sparseness_Q8))
			tmp32 = SKP_Silk_log2lin(SKP_FIX_CONST(16.0, 7) - SKP_SMULWB(essStrength_Q15, SKP_SMULWB(SKP_FIX_CONST(DE_ESSER_COEF_SWB_dB, 7), SKP_FIX_CONST(0.16, 17))))
			gain_mult_Q16 = SKP_SMULWW(gain_mult_Q16, tmp32)
		} else if psEnc.SCmn.Fs_kHz == 16 {
			var essStrength_Q15 int32 = SKP_SMULWW(-psEncCtrl.Input_tilt_Q15, SKP_SMULBB(psEnc.Speech_activity_Q8, SKP_FIX_CONST(1.0, 8)-psEncCtrl.Sparseness_Q8))
			tmp32 = SKP_Silk_log2lin(SKP_FIX_CONST(16.0, 7) - SKP_SMULWB(essStrength_Q15, SKP_SMULWB(SKP_FIX_CONST(DE_ESSER_COEF_WB_dB, 7), SKP_FIX_CONST(0.16, 17))))
			gain_mult_Q16 = SKP_SMULWW(gain_mult_Q16, tmp32)
		} else {
		}
	}
	for k = 0; k < NB_SUBFR; k++ {
		psEncCtrl.GainsPre_Q14[k] = SKP_SMULWB(gain_mult_Q16, psEncCtrl.GainsPre_Q14[k])
	}
	strength_Q16 = SKP_FIX_CONST(LOW_FREQ_SHAPING, 0) * (SKP_FIX_CONST(1.0, 16) + SKP_SMULBB(SKP_FIX_CONST(LOW_QUALITY_LOW_FREQ_SHAPING_DECR, 1), psEncCtrl.Input_quality_bands_Q15[0]-SKP_FIX_CONST(1.0, 15)))
	if psEncCtrl.SCmn.Sigtype == SIG_TYPE_VOICED {
		var fs_kHz_inv int32 = SKP_FIX_CONST(0.2, 14) / psEnc.SCmn.Fs_kHz
		for k = 0; k < NB_SUBFR; k++ {
			b_Q14 = fs_kHz_inv + SKP_FIX_CONST(3.0, 14)/(psEncCtrl.SCmn.PitchL[k])
			psEncCtrl.LF_shp_Q14[k] = (SKP_FIX_CONST(1.0, 14) - b_Q14 - SKP_SMULWB(strength_Q16, b_Q14)) << 16
			psEncCtrl.LF_shp_Q14[k] |= int32(uint16(int16(b_Q14 - SKP_FIX_CONST(1.0, 14))))
		}
		Tilt_Q16 = -SKP_FIX_CONST(HP_NOISE_COEF, 16) - SKP_SMULWB(SKP_FIX_CONST(1.0, 16)-SKP_FIX_CONST(HP_NOISE_COEF, 16), SKP_SMULWB(SKP_FIX_CONST(HARM_HP_NOISE_COEF, 24), psEnc.Speech_activity_Q8))
	} else {
		b_Q14 = 0x5333 / psEnc.SCmn.Fs_kHz
		psEncCtrl.LF_shp_Q14[0] = (SKP_FIX_CONST(1.0, 14) - b_Q14 - SKP_SMULWB(strength_Q16, SKP_SMULWB(SKP_FIX_CONST(0.6, 16), b_Q14))) << 16
		psEncCtrl.LF_shp_Q14[0] |= int32(uint16(int16(b_Q14 - SKP_FIX_CONST(1.0, 14))))
		for k = 1; k < NB_SUBFR; k++ {
			psEncCtrl.LF_shp_Q14[k] = psEncCtrl.LF_shp_Q14[0]
		}
		Tilt_Q16 = -SKP_FIX_CONST(HP_NOISE_COEF, 16)
	}
	HarmBoost_Q16 = SKP_SMULWB(SKP_SMULWB(SKP_FIX_CONST(1.0, 17)-(psEncCtrl.Coding_quality_Q14<<3), psEnc.LTPCorr_Q15), SKP_FIX_CONST(LOW_RATE_HARMONIC_BOOST, 16))
	HarmBoost_Q16 = SKP_SMLAWB(HarmBoost_Q16, SKP_FIX_CONST(1.0, 16)-(psEncCtrl.Input_quality_Q14<<2), SKP_FIX_CONST(LOW_INPUT_QUALITY_HARMONIC_BOOST, 16))
	if USE_HARM_SHAPING != 0 && psEncCtrl.SCmn.Sigtype == SIG_TYPE_VOICED {
		HarmShapeGain_Q16 = SKP_SMLAWB(SKP_FIX_CONST(HARMONIC_SHAPING, 16), SKP_FIX_CONST(1.0, 16)-SKP_SMULWB(SKP_FIX_CONST(1.0, 18)-(psEncCtrl.Coding_quality_Q14<<4), psEncCtrl.Input_quality_Q14), SKP_FIX_CONST(HIGH_RATE_OR_LOW_QUALITY_HARMONIC_SHAPING, 16))
		HarmShapeGain_Q16 = SKP_SMULWB(HarmShapeGain_Q16<<1, SKP_Silk_SQRT_APPROX(psEnc.LTPCorr_Q15<<15))
	} else {
		HarmShapeGain_Q16 = 0
	}
	for k = 0; k < NB_SUBFR; k++ {
		psShapeSt.HarmBoost_smth_Q16 = SKP_SMLAWB(psShapeSt.HarmBoost_smth_Q16, HarmBoost_Q16-psShapeSt.HarmBoost_smth_Q16, SKP_FIX_CONST(SUBFR_SMTH_COEF, 16))
		psShapeSt.HarmShapeGain_smth_Q16 = SKP_SMLAWB(psShapeSt.HarmShapeGain_smth_Q16, HarmShapeGain_Q16-psShapeSt.HarmShapeGain_smth_Q16, SKP_FIX_CONST(SUBFR_SMTH_COEF, 16))
		psShapeSt.Tilt_smth_Q16 = SKP_SMLAWB(psShapeSt.Tilt_smth_Q16, Tilt_Q16-psShapeSt.Tilt_smth_Q16, SKP_FIX_CONST(SUBFR_SMTH_COEF, 16))
		psEncCtrl.HarmBoost_Q14[k] = SKP_RSHIFT_ROUND(psShapeSt.HarmBoost_smth_Q16, 2)
		psEncCtrl.HarmShapeGain_Q14[k] = SKP_RSHIFT_ROUND(psShapeSt.HarmShapeGain_smth_Q16, 2)
		psEncCtrl.Tilt_Q14[k] = SKP_RSHIFT_ROUND(psShapeSt.Tilt_smth_Q16, 2)
	}
}
