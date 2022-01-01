package silk

import "unsafe"

func SKP_Silk_find_pred_coefs_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX, res_pitch []int16) {
	var (
		i                int32
		WLTP             [100]int32
		invGains_Q16     [4]int32
		local_gains      [4]int32
		Wght_Q15         [4]int32
		NLSF_Q15         [16]int32
		x_ptr            *int16
		x_pre_ptr        *int16
		LPC_in_pre       [544]int16
		tmp              int32
		min_gain_Q16     int32
		LTP_corrs_rshift [4]int32
	)
	min_gain_Q16 = SKP_int32_MAX >> 6
	for i = 0; int64(i) < NB_SUBFR; i++ {
		if int64(min_gain_Q16) < int64(psEncCtrl.Gains_Q16[i]) {
			min_gain_Q16 = min_gain_Q16
		} else {
			min_gain_Q16 = psEncCtrl.Gains_Q16[i]
		}
	}
	for i = 0; int64(i) < NB_SUBFR; i++ {
		invGains_Q16[i] = SKP_DIV32_varQ(min_gain_Q16, psEncCtrl.Gains_Q16[i], 16-2)
		if int64(invGains_Q16[i]) > 363 {
			invGains_Q16[i] = invGains_Q16[i]
		} else {
			invGains_Q16[i] = 363
		}
		tmp = SKP_SMULWB(invGains_Q16[i], invGains_Q16[i])
		Wght_Q15[i] = tmp >> 1
		local_gains[i] = int32((1 << 16) / int64(invGains_Q16[i]))
	}
	if int64(psEncCtrl.SCmn.Sigtype) == SIG_TYPE_VOICED {
		SKP_Silk_find_LTP_FIX(psEncCtrl.LTPCoef_Q14, WLTP, &psEncCtrl.LTPredCodGain_Q7, res_pitch, ([]int16)(&res_pitch[int64(psEnc.SCmn.Frame_length)>>1]), psEncCtrl.SCmn.PitchL, Wght_Q15, psEnc.SCmn.Subfr_length, psEnc.SCmn.Frame_length, LTP_corrs_rshift)
		SKP_Silk_quant_LTP_gains_FIX(psEncCtrl.LTPCoef_Q14[:], psEncCtrl.SCmn.LTPIndex[:], &psEncCtrl.SCmn.PERIndex, WLTP[:], psEnc.Mu_LTP_Q8, psEnc.SCmn.LTPQuantLowComplexity)
		SKP_Silk_LTP_scale_ctrl_FIX(psEnc, psEncCtrl)
		SKP_Silk_LTP_analysis_filter_FIX(&LPC_in_pre[0], (*int16)(unsafe.Add(unsafe.Pointer(&psEnc.X_buf[psEnc.SCmn.Frame_length]), -int(unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.PredictLPCOrder)))), psEncCtrl.LTPCoef_Q14, psEncCtrl.SCmn.PitchL, invGains_Q16, psEnc.SCmn.Subfr_length, psEnc.SCmn.PredictLPCOrder)
	} else {
		x_ptr = (*int16)(unsafe.Add(unsafe.Pointer(&psEnc.X_buf[psEnc.SCmn.Frame_length]), -int(unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.PredictLPCOrder))))
		x_pre_ptr = &LPC_in_pre[0]
		for i = 0; int64(i) < NB_SUBFR; i++ {
			SKP_Silk_scale_copy_vector16(x_pre_ptr, x_ptr, invGains_Q16[i], int32(int64(psEnc.SCmn.Subfr_length)+int64(psEnc.SCmn.PredictLPCOrder)))
			x_pre_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_pre_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(psEnc.SCmn.Subfr_length)+int64(psEnc.SCmn.PredictLPCOrder))))
			x_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.Subfr_length)))
		}
		memset(unsafe.Pointer(&psEncCtrl.LTPCoef_Q14[0]), 0, size_t(uintptr(NB_SUBFR*LTP_ORDER)*unsafe.Sizeof(int16(0))))
		psEncCtrl.LTPredCodGain_Q7 = 0
	}
	SKP_Silk_find_LPC_FIX(NLSF_Q15[:], &psEncCtrl.SCmn.NLSFInterpCoef_Q2, psEnc.SPred.Prev_NLSFq_Q15[:], int32(int64(psEnc.SCmn.UseInterpolatedNLSFs)*(1-int64(psEnc.SCmn.First_frame_after_reset))), psEnc.SCmn.PredictLPCOrder, LPC_in_pre[:], int32(int64(psEnc.SCmn.Subfr_length)+int64(psEnc.SCmn.PredictLPCOrder)))
	SKP_Silk_process_NLSFs_FIX(psEnc, psEncCtrl, &NLSF_Q15[0])
	SKP_Silk_residual_energy_FIX(psEncCtrl.ResNrg, psEncCtrl.ResNrgQ, LPC_in_pre[:], psEncCtrl.PredCoef_Q12, local_gains, psEnc.SCmn.Subfr_length, psEnc.SCmn.PredictLPCOrder)
	memcpy(unsafe.Pointer(&psEnc.SPred.Prev_NLSFq_Q15[0]), unsafe.Pointer(&NLSF_Q15[0]), size_t(uintptr(psEnc.SCmn.PredictLPCOrder)*unsafe.Sizeof(int32(0))))
}
