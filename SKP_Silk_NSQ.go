package silk

import "unsafe"

func SKP_Silk_NSQ(psEncC *SKP_Silk_encoder_state, psEncCtrlC *SKP_Silk_encoder_control, NSQ *SKP_Silk_nsq_state, x []int16, q []int8, LSFInterpFactor_Q2 int32, PredCoef_Q12 [32]int16, LTPCoef_Q14 [20]int16, AR2_Q13 [64]int16, HarmShapeGain_Q14 [4]int32, Tilt_Q14 [4]int32, LF_shp_Q14 [4]int32, Gains_Q16 [4]int32, Lambda_Q10 int32, LTP_scale_Q14 int32) {
	var (
		k                      int32
		lag                    int32
		start_idx              int32
		LSF_interpolation_flag int32
		A_Q12                  *int16
		B_Q14                  *int16
		AR_shp_Q13             *int16
		pxq                    *int16
		sLTP_Q16               [960]int32
		sLTP                   [960]int16
		HarmShapeFIRPacked_Q14 int32
		offset_Q10             int32
		FiltState              [16]int32
		x_sc_Q10               [120]int32
	)
	NSQ.Rand_seed = psEncCtrlC.Seed
	lag = NSQ.LagPrev
	SKP_assert(NSQ.Prev_inv_gain_Q16 != 0)
	offset_Q10 = int32(SKP_Silk_Quantization_Offsets_Q10[psEncCtrlC.Sigtype][psEncCtrlC.QuantOffsetType])
	if LSFInterpFactor_Q2 == (1 << 2) {
		LSF_interpolation_flag = 0
	} else {
		LSF_interpolation_flag = 1
	}
	NSQ.SLTP_shp_buf_idx = psEncC.Frame_length
	NSQ.SLTP_buf_idx = psEncC.Frame_length
	pxq = &NSQ.Xq[psEncC.Frame_length]
	for k = 0; k < NB_SUBFR; k++ {
		A_Q12 = &PredCoef_Q12[((k>>1)|(1-LSF_interpolation_flag))*MAX_LPC_ORDER]
		B_Q14 = &LTPCoef_Q14[k*LTP_ORDER]
		AR_shp_Q13 = &AR2_Q13[k*MAX_SHAPE_LPC_ORDER]
		SKP_assert(HarmShapeGain_Q14[k] >= 0)
		HarmShapeFIRPacked_Q14 = (HarmShapeGain_Q14[k]) >> 2
		HarmShapeFIRPacked_Q14 |= ((HarmShapeGain_Q14[k]) >> 1) << 16
		NSQ.Rewhite_flag = 0
		if psEncCtrlC.Sigtype == SIG_TYPE_VOICED {
			lag = psEncCtrlC.PitchL[k]
			if (k & (3 - (LSF_interpolation_flag << 1))) == 0 {
				start_idx = psEncC.Frame_length - lag - psEncC.PredictLPCOrder - LTP_ORDER/2
				SKP_assert(start_idx >= 0)
				SKP_assert(start_idx <= psEncC.Frame_length-psEncC.PredictLPCOrder)
				memset(unsafe.Pointer(&FiltState[0]), 0, size_t(uintptr(psEncC.PredictLPCOrder)*unsafe.Sizeof(int32(0))))
				SKP_Silk_MA_Prediction(([]int16)(&NSQ.Xq[start_idx+k*(psEncC.Frame_length>>2)]), ([]int16)(A_Q12), FiltState[:], ([]int16)(&sLTP[start_idx]), psEncC.Frame_length-start_idx, psEncC.PredictLPCOrder)
				NSQ.Rewhite_flag = 1
				NSQ.SLTP_buf_idx = psEncC.Frame_length
			}
		}
		SKP_Silk_nsq_scale_states(NSQ, x, x_sc_Q10[:], psEncC.Subfr_length, sLTP[:], sLTP_Q16[:], k, LTP_scale_Q14, Gains_Q16, psEncCtrlC.PitchL)
		SKP_Silk_noise_shape_quantizer(NSQ, psEncCtrlC.Sigtype, x_sc_Q10[:], q, ([]int16)(pxq), sLTP_Q16[:], ([]int16)(A_Q12), ([]int16)(B_Q14), ([]int16)(AR_shp_Q13), lag, HarmShapeFIRPacked_Q14, Tilt_Q14[k], LF_shp_Q14[k], Gains_Q16[k], Lambda_Q10, offset_Q10, psEncC.Subfr_length, psEncC.ShapingLPCOrder, psEncC.PredictLPCOrder)
		x += ([]int16)(psEncC.Subfr_length)
		q += ([]int8)(psEncC.Subfr_length)
		pxq = (*int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(int16(0))*uintptr(psEncC.Subfr_length)))
	}
	NSQ.LagPrev = psEncCtrlC.PitchL[NB_SUBFR-1]
	memcpy(unsafe.Pointer(&NSQ.Xq[0]), unsafe.Pointer(&NSQ.Xq[psEncC.Frame_length]), size_t(uintptr(psEncC.Frame_length)*unsafe.Sizeof(int16(0))))
	memcpy(unsafe.Pointer(&NSQ.SLTP_shp_Q10[0]), unsafe.Pointer(&NSQ.SLTP_shp_Q10[psEncC.Frame_length]), size_t(uintptr(psEncC.Frame_length)*unsafe.Sizeof(int32(0))))
}
func SKP_Silk_noise_shape_quantizer(NSQ *SKP_Silk_nsq_state, sigtype int32, x_sc_Q10 []int32, q []int8, xq []int16, sLTP_Q16 []int32, a_Q12 []int16, b_Q14 []int16, AR_shp_Q13 []int16, lag int32, HarmShapeFIRPacked_Q14 int32, Tilt_Q14 int32, LF_shp_Q14 int32, Gain_Q16 int32, Lambda_Q10 int32, offset_Q10 int32, length int32, shapingLPCOrder int32, predictLPCOrder int32) {
	var (
		i              int32
		j              int32
		LTP_pred_Q14   int32
		LPC_pred_Q10   int32
		n_AR_Q10       int32
		n_LTP_Q14      int32
		n_LF_Q10       int32
		r_Q10          int32
		q_Q0           int32
		q_Q10          int32
		thr1_Q10       int32
		thr2_Q10       int32
		thr3_Q10       int32
		dither         int32
		exc_Q10        int32
		LPC_exc_Q10    int32
		xq_Q10         int32
		tmp1           int32
		tmp2           int32
		sLF_AR_shp_Q10 int32
		psLPC_Q14      []int32
		shp_lag_ptr    *int32
		pred_lag_ptr   []int32
	)
	shp_lag_ptr = &NSQ.SLTP_shp_Q10[NSQ.SLTP_shp_buf_idx-lag+HARM_SHAPE_FIR_TAPS/2]
	pred_lag_ptr = ([]int32)(&sLTP_Q16[NSQ.SLTP_buf_idx-lag+LTP_ORDER/2])
	psLPC_Q14 = ([]int32)(&NSQ.SLPC_Q14[DECISION_DELAY-1])
	thr1_Q10 = int32(int64(-1536) - int64(Lambda_Q10>>1))
	thr2_Q10 = int32(int64(-512) - int64(Lambda_Q10>>1))
	thr2_Q10 = thr2_Q10 + (SKP_SMULBB(offset_Q10, Lambda_Q10) >> 10)
	thr3_Q10 = (Lambda_Q10 >> 1) + 512
	for i = 0; i < length; i++ {
		NSQ.Rand_seed = int32((uint32(NSQ.Rand_seed) * 0xBB38435) + 0x3619636B)
		dither = NSQ.Rand_seed >> 31
		SKP_assert((predictLPCOrder & 1) == 0)
		SKP_assert((int64(uintptr(unsafe.Pointer((*int8)(unsafe.Add(unsafe.Pointer((*int8)(unsafe.Pointer(&a_Q12[0]))), 0))))) & 3) == 0)
		SKP_assert(predictLPCOrder >= 10)
		LPC_pred_Q10 = SKP_SMULWB(psLPC_Q14[0], int32(a_Q12[0]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-1], int32(a_Q12[1]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-2], int32(a_Q12[2]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-3], int32(a_Q12[3]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-4], int32(a_Q12[4]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-5], int32(a_Q12[5]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-6], int32(a_Q12[6]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-7], int32(a_Q12[7]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-8], int32(a_Q12[8]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-9], int32(a_Q12[9]))
		for j = 10; j < predictLPCOrder; j++ {
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psLPC_Q14[-j], int32(a_Q12[j]))
		}
		if sigtype == SIG_TYPE_VOICED {
			LTP_pred_Q14 = SKP_SMULWB(pred_lag_ptr[0], int32(b_Q14[0]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, pred_lag_ptr[-1], int32(b_Q14[1]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, pred_lag_ptr[-2], int32(b_Q14[2]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, pred_lag_ptr[-3], int32(b_Q14[3]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, pred_lag_ptr[-4], int32(b_Q14[4]))
			pred_lag_ptr++
		} else {
			LTP_pred_Q14 = 0
		}
		SKP_assert((shapingLPCOrder & 1) == 0)
		tmp2 = psLPC_Q14[0]
		tmp1 = NSQ.SAR2_Q14[0]
		NSQ.SAR2_Q14[0] = tmp2
		n_AR_Q10 = SKP_SMULWB(tmp2, int32(AR_shp_Q13[0]))
		for j = 2; j < shapingLPCOrder; j += 2 {
			tmp2 = NSQ.SAR2_Q14[j-1]
			NSQ.SAR2_Q14[j-1] = tmp1
			n_AR_Q10 = SKP_SMLAWB(n_AR_Q10, tmp1, int32(AR_shp_Q13[j-1]))
			tmp1 = NSQ.SAR2_Q14[j+0]
			NSQ.SAR2_Q14[j+0] = tmp2
			n_AR_Q10 = SKP_SMLAWB(n_AR_Q10, tmp2, int32(AR_shp_Q13[j]))
		}
		NSQ.SAR2_Q14[shapingLPCOrder-1] = tmp1
		n_AR_Q10 = SKP_SMLAWB(n_AR_Q10, tmp1, int32(AR_shp_Q13[shapingLPCOrder-1]))
		n_AR_Q10 = n_AR_Q10 >> 1
		n_AR_Q10 = SKP_SMLAWB(n_AR_Q10, NSQ.SLF_AR_shp_Q12, Tilt_Q14)
		n_LF_Q10 = SKP_SMULWB(NSQ.SLTP_shp_Q10[NSQ.SLTP_shp_buf_idx-1], LF_shp_Q14) << 2
		n_LF_Q10 = SKP_SMLAWT(n_LF_Q10, NSQ.SLF_AR_shp_Q12, LF_shp_Q14)
		SKP_assert(lag > 0 || sigtype == SIG_TYPE_UNVOICED)
		if lag > 0 {
			n_LTP_Q14 = SKP_SMULWB((*(*int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(int32(0))*0)))+(*(*int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(int32(0))*2)))), HarmShapeFIRPacked_Q14)
			n_LTP_Q14 = SKP_SMLAWT(n_LTP_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(int32(0))*1))), HarmShapeFIRPacked_Q14)
			n_LTP_Q14 = n_LTP_Q14 << 6
			shp_lag_ptr = (*int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(int32(0))*1))
		} else {
			n_LTP_Q14 = 0
		}
		tmp1 = LTP_pred_Q14 - n_LTP_Q14
		tmp1 = tmp1 >> 4
		tmp1 = tmp1 + LPC_pred_Q10
		tmp1 = tmp1 - n_AR_Q10
		tmp1 = tmp1 - n_LF_Q10
		r_Q10 = (x_sc_Q10[i]) - tmp1
		r_Q10 = (r_Q10 ^ dither) - dither
		r_Q10 = r_Q10 - offset_Q10
		r_Q10 = SKP_LIMIT_32(r_Q10, int32(-64<<10), 64<<10)
		q_Q0 = 0
		q_Q10 = 0
		if r_Q10 < thr2_Q10 {
			if r_Q10 < thr1_Q10 {
				q_Q0 = SKP_RSHIFT_ROUND(r_Q10+(Lambda_Q10>>1), 10)
				q_Q10 = q_Q0 << 10
			} else {
				q_Q0 = -1
				q_Q10 = -1024
			}
		} else {
			if r_Q10 > thr3_Q10 {
				q_Q0 = SKP_RSHIFT_ROUND(r_Q10-(Lambda_Q10>>1), 10)
				q_Q10 = q_Q0 << 10
			}
		}
		q[i] = int8(q_Q0)
		exc_Q10 = q_Q10 + offset_Q10
		exc_Q10 = (exc_Q10 ^ dither) - dither
		LPC_exc_Q10 = exc_Q10 + SKP_RSHIFT_ROUND(LTP_pred_Q14, 4)
		xq_Q10 = LPC_exc_Q10 + LPC_pred_Q10
		xq[i] = SKP_SAT16(SKP_RSHIFT_ROUND(SKP_SMULWW(xq_Q10, Gain_Q16), 10))
		psLPC_Q14++
		psLPC_Q14[0] = xq_Q10 << 4
		sLF_AR_shp_Q10 = xq_Q10 - n_AR_Q10
		NSQ.SLF_AR_shp_Q12 = sLF_AR_shp_Q10 << 2
		NSQ.SLTP_shp_Q10[NSQ.SLTP_shp_buf_idx] = sLF_AR_shp_Q10 - n_LF_Q10
		sLTP_Q16[NSQ.SLTP_buf_idx] = LPC_exc_Q10 << 6
		NSQ.SLTP_shp_buf_idx++
		NSQ.SLTP_buf_idx++
		NSQ.Rand_seed += int32(q[i])
	}
	memcpy(unsafe.Pointer(&NSQ.SLPC_Q14[0]), unsafe.Pointer(&NSQ.SLPC_Q14[length]), size_t(DECISION_DELAY*unsafe.Sizeof(int32(0))))
}
func SKP_Silk_nsq_scale_states(NSQ *SKP_Silk_nsq_state, x []int16, x_sc_Q10 []int32, subfr_length int32, sLTP []int16, sLTP_Q16 []int32, subfr int32, LTP_scale_Q14 int32, Gains_Q16 [4]int32, pitchL [4]int32) {
	var (
		i            int32
		lag          int32
		inv_gain_Q16 int32
		gain_adj_Q16 int32
		inv_gain_Q32 int32
	)
	inv_gain_Q16 = SKP_INVERSE32_varQ(func() int32 {
		if (Gains_Q16[subfr]) > 1 {
			return Gains_Q16[subfr]
		}
		return 1
	}(), 32)
	if inv_gain_Q16 < SKP_int16_MAX {
		inv_gain_Q16 = inv_gain_Q16
	} else {
		inv_gain_Q16 = SKP_int16_MAX
	}
	lag = pitchL[subfr]
	if NSQ.Rewhite_flag != 0 {
		inv_gain_Q32 = inv_gain_Q16 << 16
		if subfr == 0 {
			inv_gain_Q32 = SKP_SMULWB(inv_gain_Q32, LTP_scale_Q14) << 2
		}
		for i = NSQ.SLTP_buf_idx - lag - LTP_ORDER/2; i < NSQ.SLTP_buf_idx; i++ {
			SKP_assert(i < (FRAME_LENGTH_MS * MAX_FS_KHZ))
			sLTP_Q16[i] = SKP_SMULWB(inv_gain_Q32, int32(sLTP[i]))
		}
	}
	if inv_gain_Q16 != NSQ.Prev_inv_gain_Q16 {
		gain_adj_Q16 = SKP_DIV32_varQ(inv_gain_Q16, NSQ.Prev_inv_gain_Q16, 16)
		for i = NSQ.SLTP_shp_buf_idx - subfr_length*NB_SUBFR; i < NSQ.SLTP_shp_buf_idx; i++ {
			NSQ.SLTP_shp_Q10[i] = SKP_SMULWW(gain_adj_Q16, NSQ.SLTP_shp_Q10[i])
		}
		if NSQ.Rewhite_flag == 0 {
			for i = NSQ.SLTP_buf_idx - lag - LTP_ORDER/2; i < NSQ.SLTP_buf_idx; i++ {
				sLTP_Q16[i] = SKP_SMULWW(gain_adj_Q16, sLTP_Q16[i])
			}
		}
		NSQ.SLF_AR_shp_Q12 = SKP_SMULWW(gain_adj_Q16, NSQ.SLF_AR_shp_Q12)
		for i = 0; i < DECISION_DELAY; i++ {
			NSQ.SLPC_Q14[i] = SKP_SMULWW(gain_adj_Q16, NSQ.SLPC_Q14[i])
		}
		for i = 0; i < MAX_SHAPE_LPC_ORDER; i++ {
			NSQ.SAR2_Q14[i] = SKP_SMULWW(gain_adj_Q16, NSQ.SAR2_Q14[i])
		}
	}
	for i = 0; i < subfr_length; i++ {
		x_sc_Q10[i] = SKP_SMULBB(int32(x[i]), int32(int16(inv_gain_Q16))) >> 6
	}
	SKP_assert(inv_gain_Q16 != 0)
	NSQ.Prev_inv_gain_Q16 = inv_gain_Q16
}
