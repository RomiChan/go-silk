package silk

import "unsafe"

func SKP_Silk_find_pitch_lags_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX, res []int16, x []int16) {
	var (
		psPredSt   = &psEnc.SPred
		buf_len    int32
		i          int32
		scale      int32
		thrhld_Q15 int32
		res_nrg    int32
		x_buf      *int16
		x_buf_ptr  *int16
		Wsig       [576]int16
		Wsig_ptr   *int16
		auto_corr  [17]int32
		rc_Q15     [16]int16
		A_Q24      [16]int32
		FiltState  [16]int32
		A_Q12      [16]int16
	)
	buf_len = psEnc.SCmn.La_pitch + (psEnc.SCmn.Frame_length << 1)
	x_buf = (*int16)(unsafe.Add(unsafe.Pointer(&x[0]), -int(unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.Frame_length))))
	x_buf_ptr = (*int16)(unsafe.Add(unsafe.Add(unsafe.Pointer(x_buf), unsafe.Sizeof(int16(0))*uintptr(buf_len)), -int(unsafe.Sizeof(int16(0))*uintptr(psPredSt.Pitch_LPC_win_length))))
	Wsig_ptr = &Wsig[0]
	SKP_Silk_apply_sine_window(([]int16)(Wsig_ptr), ([]int16)(x_buf_ptr), 1, psEnc.SCmn.La_pitch)
	Wsig_ptr = (*int16)(unsafe.Add(unsafe.Pointer(Wsig_ptr), unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.La_pitch)))
	x_buf_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_buf_ptr), unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.La_pitch)))
	memcpy(unsafe.Pointer(Wsig_ptr), unsafe.Pointer(x_buf_ptr), uintptr(psPredSt.Pitch_LPC_win_length-(psEnc.SCmn.La_pitch<<1))*unsafe.Sizeof(int16(0)))
	Wsig_ptr = (*int16)(unsafe.Add(unsafe.Pointer(Wsig_ptr), unsafe.Sizeof(int16(0))*uintptr(psPredSt.Pitch_LPC_win_length-(psEnc.SCmn.La_pitch<<1))))
	x_buf_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_buf_ptr), unsafe.Sizeof(int16(0))*uintptr(psPredSt.Pitch_LPC_win_length-(psEnc.SCmn.La_pitch<<1))))
	SKP_Silk_apply_sine_window(([]int16)(Wsig_ptr), ([]int16)(x_buf_ptr), 2, psEnc.SCmn.La_pitch)
	SKP_Silk_autocorr(auto_corr[:], &scale, Wsig[:], psPredSt.Pitch_LPC_win_length, psEnc.SCmn.PitchEstimationLPCOrder+1)
	auto_corr[0] = SKP_SMLAWB(auto_corr[0], auto_corr[0], SKP_FIX_CONST(0.001, 16))
	res_nrg = SKP_Silk_schur(&rc_Q15[0], &auto_corr[0], psEnc.SCmn.PitchEstimationLPCOrder)
	psEncCtrl.PredGain_Q16 = SKP_DIV32_varQ(auto_corr[0], SKP_max_int(res_nrg, 1), 16)
	SKP_Silk_k2a(A_Q24[:], rc_Q15[:], psEnc.SCmn.PitchEstimationLPCOrder)
	for i = 0; i < psEnc.SCmn.PitchEstimationLPCOrder; i++ {
		A_Q12[i] = SKP_SAT16((A_Q24[i]) >> 12)
	}
	SKP_Silk_bwexpander(&A_Q12[0], psEnc.SCmn.PitchEstimationLPCOrder, SKP_FIX_CONST(0.99, 16))
	memset(unsafe.Pointer(&FiltState[0]), 0, uintptr(psEnc.SCmn.PitchEstimationLPCOrder)*unsafe.Sizeof(int32(0)))
	SKP_Silk_MA_Prediction(x_buf, &A_Q12[0], &FiltState[0], &res[0], buf_len, psEnc.SCmn.PitchEstimationLPCOrder)
	memset(unsafe.Pointer(&res[0]), 0, uintptr(psEnc.SCmn.PitchEstimationLPCOrder)*unsafe.Sizeof(int16(0)))
	thrhld_Q15 = SKP_FIX_CONST(0.45, 15)
	thrhld_Q15 = SKP_SMLABB(thrhld_Q15, SKP_FIX_CONST(-0.004, 15), psEnc.SCmn.PitchEstimationLPCOrder)
	thrhld_Q15 = SKP_SMLABB(thrhld_Q15, SKP_FIX_CONST(-0.1, 7), psEnc.Speech_activity_Q8)
	thrhld_Q15 = SKP_SMLABB(thrhld_Q15, SKP_FIX_CONST(0.15, 15), psEnc.SCmn.Prev_sigtype)
	thrhld_Q15 = SKP_SMLAWB(thrhld_Q15, SKP_FIX_CONST(-0.1, 16), psEncCtrl.Input_tilt_Q15)
	thrhld_Q15 = int32(SKP_SAT16(thrhld_Q15))
	psEncCtrl.SCmn.Sigtype = SKP_Silk_pitch_analysis_core(&res[0], &psEncCtrl.SCmn.PitchL[0], &psEncCtrl.SCmn.LagIndex, &psEncCtrl.SCmn.ContourIndex, &psEnc.LTPCorr_Q15, psEnc.SCmn.PrevLag, psEnc.SCmn.PitchEstimationThreshold_Q16, int32(int16(thrhld_Q15)), psEnc.SCmn.Fs_kHz, psEnc.SCmn.PitchEstimationComplexity, SKP_FALSE)
}
