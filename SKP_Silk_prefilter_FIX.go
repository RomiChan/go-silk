package silk

import "unsafe"

func SKP_Silk_warped_LPC_analysis_filter_FIX(state []int32, res []int16, coef_Q13 []int16, input []int16, lambda_Q16 int16, length int, order int) {
	var (
		acc_Q11 int32
		tmp1    int32
		tmp2    int32
	)
	for n := 0; n < length; n++ {
		tmp2 = SKP_SMLAWB(state[0], state[1], int32(lambda_Q16))
		state[0] = int32(input[n] << 14)
		tmp1 = SKP_SMLAWB(state[1], int32(int64(state[2])-int64(tmp2)), int32(lambda_Q16))
		state[1] = tmp2
		acc_Q11 = SKP_SMULWB(tmp2, int32(coef_Q13[0]))
		for i := 2; i < order; i += 2 {
			tmp2 = SKP_SMLAWB(state[i], state[i+1]-tmp1, int32(lambda_Q16))
			state[i] = tmp1
			acc_Q11 = SKP_SMLAWB(acc_Q11, tmp1, int32(coef_Q13[i-1]))
			tmp1 = SKP_SMLAWB(state[i+1], state[i+2]-tmp2, int32(lambda_Q16))
			state[i+1] = tmp2
			acc_Q11 = SKP_SMLAWB(acc_Q11, tmp2, int32(coef_Q13[i]))
		}
		state[order] = tmp1
		acc_Q11 = SKP_SMLAWB(acc_Q11, tmp1, int32(coef_Q13[order-1]))
		res[n] = SKP_SAT16(int32(input[n]) - SKP_RSHIFT_ROUND(acc_Q11, 11))
	}
}
func SKP_Silk_prefilter_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX, xw []int16, x []int16) {
	var (
		P                      *SKP_Silk_prefilter_state_FIX = &psEnc.SPrefilt
		j                      int32
		k                      int32
		lag                    int32
		tmp_32                 int32
		AR1_shp_Q13            *int16
		px                     *int16
		pxw                    *int16
		HarmShapeGain_Q12      int32
		Tilt_Q14               int32
		HarmShapeFIRPacked_Q12 int32
		LF_shp_Q14             int32
		x_filt_Q12             [120]int32
		st_res                 [136]int16
		B_Q12                  [2]int16
	)
	px = &x[0]
	pxw = &xw[0]
	lag = P.LagPrev
	for k = 0; k < NB_SUBFR; k++ {
		if psEncCtrl.SCmn.Sigtype == SIG_TYPE_VOICED {
			lag = psEncCtrl.SCmn.PitchL[k]
		}
		HarmShapeGain_Q12 = SKP_SMULWB(psEncCtrl.HarmShapeGain_Q14[k], 0x4000-psEncCtrl.HarmBoost_Q14[k])
		HarmShapeFIRPacked_Q12 = HarmShapeGain_Q12 >> 2
		HarmShapeFIRPacked_Q12 |= (HarmShapeGain_Q12 >> 1) << 16
		Tilt_Q14 = psEncCtrl.Tilt_Q14[k]
		LF_shp_Q14 = psEncCtrl.LF_shp_Q14[k]
		AR1_shp_Q13 = &psEncCtrl.AR1_Q13[k*MAX_SHAPE_LPC_ORDER]
		SKP_Silk_warped_LPC_analysis_filter_FIX(P.SAR_shp[:], st_res[:], ([]int16)(AR1_shp_Q13), ([]int16)(px), int16(psEnc.SCmn.Warping_Q16), psEnc.SCmn.Subfr_length, psEnc.SCmn.ShapingLPCOrder)
		B_Q12[0] = int16(SKP_RSHIFT_ROUND(psEncCtrl.GainsPre_Q14[k], 2))
		tmp_32 = SKP_SMLABB(SKP_FIX_CONST(INPUT_TILT, 26), psEncCtrl.HarmBoost_Q14[k], HarmShapeGain_Q12)
		tmp_32 = SKP_SMLABB(tmp_32, psEncCtrl.Coding_quality_Q14, SKP_FIX_CONST(HIGH_RATE_INPUT_TILT, 12))
		tmp_32 = SKP_SMULWB(tmp_32, -psEncCtrl.GainsPre_Q14[k])
		tmp_32 = SKP_RSHIFT_ROUND(tmp_32, 12)
		B_Q12[1] = SKP_SAT16(int16(tmp_32))
		x_filt_Q12[0] = SKP_SMLABB(SKP_SMULBB(int32(st_res[0]), int32(B_Q12[0])), P.SHarmHP, int32(B_Q12[1]))
		for j = 1; j < psEnc.SCmn.Subfr_length; j++ {
			x_filt_Q12[j] = SKP_SMLABB(SKP_SMULBB(int32(st_res[j]), int32(B_Q12[0])), int32(st_res[j-1]), int32(B_Q12[1]))
		}
		P.SHarmHP = int32(st_res[psEnc.SCmn.Subfr_length-1])
		SKP_Silk_prefilt_FIX(P, x_filt_Q12[:], ([]int16)(pxw), HarmShapeFIRPacked_Q12, Tilt_Q14, LF_shp_Q14, lag, psEnc.SCmn.Subfr_length)
		px = (*int16)(unsafe.Add(unsafe.Pointer(px), unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.Subfr_length)))
		pxw = (*int16)(unsafe.Add(unsafe.Pointer(pxw), unsafe.Sizeof(int16(0))*uintptr(psEnc.SCmn.Subfr_length)))
	}
	P.LagPrev = psEncCtrl.SCmn.PitchL[NB_SUBFR-1]
}
func SKP_Silk_prefilt_FIX(P *SKP_Silk_prefilter_state_FIX, st_res_Q12 []int32, xw []int16, HarmShapeFIRPacked_Q12 int32, Tilt_Q14 int32, LF_shp_Q14 int32, lag int32, length int32) {
	var (
		i               int32
		idx             int32
		LTP_shp_buf_idx int32
		n_LTP_Q12       int32
		n_Tilt_Q10      int32
		n_LF_Q10        int32
		sLF_MA_shp_Q12  int32
		sLF_AR_shp_Q12  int32
		LTP_shp_buf     *int16
	)
	LTP_shp_buf = &P.SLTP_shp[0]
	LTP_shp_buf_idx = P.SLTP_shp_buf_idx
	sLF_AR_shp_Q12 = P.SLF_AR_shp_Q12
	sLF_MA_shp_Q12 = P.SLF_MA_shp_Q12
	for i = 0; i < length; i++ {
		if lag > 0 {
			idx = lag + LTP_shp_buf_idx
			n_LTP_Q12 = SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(LTP_shp_buf), unsafe.Sizeof(int16(0))*uintptr((idx-HARM_SHAPE_FIR_TAPS/2-1)&(LTP_BUF_LENGTH-1))))), HarmShapeFIRPacked_Q12)
			n_LTP_Q12 = SKP_SMLABT(n_LTP_Q12, int32(*(*int16)(unsafe.Add(unsafe.Pointer(LTP_shp_buf), unsafe.Sizeof(int16(0))*uintptr((idx-HARM_SHAPE_FIR_TAPS/2)&(LTP_BUF_LENGTH-1))))), HarmShapeFIRPacked_Q12)
			n_LTP_Q12 = SKP_SMLABB(n_LTP_Q12, int32(*(*int16)(unsafe.Add(unsafe.Pointer(LTP_shp_buf), unsafe.Sizeof(int16(0))*uintptr((idx-HARM_SHAPE_FIR_TAPS/2+1)&(LTP_BUF_LENGTH-1))))), HarmShapeFIRPacked_Q12)
		} else {
			n_LTP_Q12 = 0
		}
		n_Tilt_Q10 = SKP_SMULWB(sLF_AR_shp_Q12, Tilt_Q14)
		n_LF_Q10 = SKP_SMLAWB(SKP_SMULWT(sLF_AR_shp_Q12, LF_shp_Q14), sLF_MA_shp_Q12, LF_shp_Q14)
		sLF_AR_shp_Q12 = (st_res_Q12[i]) - (n_Tilt_Q10 << 2)
		sLF_MA_shp_Q12 = sLF_AR_shp_Q12 - (n_LF_Q10 << 2)
		LTP_shp_buf_idx = (LTP_shp_buf_idx - 1) & (LTP_BUF_LENGTH - 1)
		*(*int16)(unsafe.Add(unsafe.Pointer(LTP_shp_buf), unsafe.Sizeof(int16(0))*uintptr(LTP_shp_buf_idx))) = SKP_SAT16(SKP_RSHIFT_ROUND(sLF_MA_shp_Q12, 12))
		xw[i] = SKP_SAT16(SKP_RSHIFT_ROUND(sLF_MA_shp_Q12-n_LTP_Q12, 12))
	}
	P.SLF_AR_shp_Q12 = sLF_AR_shp_Q12
	P.SLF_MA_shp_Q12 = sLF_MA_shp_Q12
	P.SLTP_shp_buf_idx = LTP_shp_buf_idx
}
