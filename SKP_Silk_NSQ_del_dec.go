package silk

import "unsafe"

type NSQ_del_dec_struct struct {
	RandState [32]int32
	Q_Q10     [32]int32
	Xq_Q10    [32]int32
	Pred_Q16  [32]int32
	Shape_Q10 [32]int32
	Gain_Q16  [32]int32
	SAR2_Q14  [16]int32
	SLPC_Q14  [152]int32
	LF_AR_Q12 int32
	Seed      int32
	SeedInit  int32
	RD_Q10    int32
}
type NSQ_sample_struct struct {
	Q_Q10        int32
	RD_Q10       int32
	Xq_Q14       int32
	LF_AR_Q12    int32
	SLTP_shp_Q10 int32
	LPC_exc_Q16  int32
}

func SKP_Silk_NSQ_del_dec(psEncC *SKP_Silk_encoder_state, psEncCtrlC *SKP_Silk_encoder_control, NSQ *SKP_Silk_nsq_state, x []int16, q []int8, LSFInterpFactor_Q2 int32, PredCoef_Q12 [32]int16, LTPCoef_Q14 [20]int16, AR2_Q13 [64]int16, HarmShapeGain_Q14 [4]int32, Tilt_Q14 [4]int32, LF_shp_Q14 [4]int32, Gains_Q16 [4]int32, Lambda_Q10 int32, LTP_scale_Q14 int32) {
	var (
		i                      int32
		k                      int32
		lag                    int32
		start_idx              int32
		LSF_interpolation_flag int32
		Winner_ind             int32
		subfr                  int32
		last_smple_idx         int32
		smpl_buf_idx           int32
		decisionDelay          int32
		subfr_length           int32
		A_Q12                  *int16
		B_Q14                  *int16
		AR_shp_Q13             *int16
		pxq                    *int16
		sLTP_Q16               [960]int32
		sLTP                   [960]int16
		HarmShapeFIRPacked_Q14 int32
		offset_Q10             int32
		FiltState              [16]int32
		RDmin_Q10              int32
		x_sc_Q10               [120]int32
		psDelDec               [4]NSQ_del_dec_struct
		psDD                   *NSQ_del_dec_struct
	)
	subfr_length = psEncC.Frame_length / NB_SUBFR
	lag = NSQ.LagPrev
	memset(unsafe.Pointer(&psDelDec[0]), 0, uintptr(psEncC.NStatesDelayedDecision)*unsafe.Sizeof(NSQ_del_dec_struct{}))
	for k = 0; k < psEncC.NStatesDelayedDecision; k++ {
		psDD = &psDelDec[k]
		psDD.Seed = (k + psEncCtrlC.Seed) & 3
		psDD.SeedInit = psDD.Seed
		psDD.RD_Q10 = 0
		psDD.LF_AR_Q12 = NSQ.SLF_AR_shp_Q12
		psDD.Shape_Q10[0] = NSQ.SLTP_shp_Q10[psEncC.Frame_length-1]
		memcpy(unsafe.Pointer(&psDD.SLPC_Q14[0]), unsafe.Pointer(&NSQ.SLPC_Q14[0]), size_t(DECISION_DELAY*unsafe.Sizeof(int32(0))))
		memcpy(unsafe.Pointer(&psDD.SAR2_Q14[0]), unsafe.Pointer(&NSQ.SAR2_Q14[0]), size_t(unsafe.Sizeof([16]int32{})))
	}
	offset_Q10 = int32(SKP_Silk_Quantization_Offsets_Q10[psEncCtrlC.Sigtype][psEncCtrlC.QuantOffsetType])
	smpl_buf_idx = 0
	decisionDelay = SKP_min_int(DECISION_DELAY, subfr_length)
	if psEncCtrlC.Sigtype == SIG_TYPE_VOICED {
		for k = 0; k < NB_SUBFR; k++ {
			decisionDelay = SKP_min_int(decisionDelay, psEncCtrlC.PitchL[k]-LTP_ORDER/2-1)
		}
	} else {
		if lag > 0 {
			decisionDelay = SKP_min_int(decisionDelay, lag-LTP_ORDER/2-1)
		}
	}
	if LSFInterpFactor_Q2 == (1 << 2) {
		LSF_interpolation_flag = 0
	} else {
		LSF_interpolation_flag = 1
	}
	pxq = &NSQ.Xq[psEncC.Frame_length]
	NSQ.SLTP_shp_buf_idx = psEncC.Frame_length
	NSQ.SLTP_buf_idx = psEncC.Frame_length
	subfr = 0
	for k = 0; k < NB_SUBFR; k++ {
		A_Q12 = &PredCoef_Q12[((k>>1)|(1-LSF_interpolation_flag))*MAX_LPC_ORDER]
		B_Q14 = &LTPCoef_Q14[k*LTP_ORDER]
		AR_shp_Q13 = &AR2_Q13[k*MAX_SHAPE_LPC_ORDER]
		HarmShapeFIRPacked_Q14 = (HarmShapeGain_Q14[k]) >> 2
		HarmShapeFIRPacked_Q14 |= ((HarmShapeGain_Q14[k]) >> 1) << 16
		NSQ.Rewhite_flag = 0
		if psEncCtrlC.Sigtype == SIG_TYPE_VOICED {
			lag = psEncCtrlC.PitchL[k]
			if (k & (3 - (LSF_interpolation_flag << 1))) == 0 {
				if k == 2 {
					RDmin_Q10 = psDelDec[0].RD_Q10
					Winner_ind = 0
					for i = 1; i < psEncC.NStatesDelayedDecision; i++ {
						if psDelDec[i].RD_Q10 < RDmin_Q10 {
							RDmin_Q10 = psDelDec[i].RD_Q10
							Winner_ind = i
						}
					}
					for i = 0; i < psEncC.NStatesDelayedDecision; i++ {
						if i != Winner_ind {
							psDelDec[i].RD_Q10 += SKP_int32_MAX >> 4
						}
					}
					psDD = &psDelDec[Winner_ind]
					last_smple_idx = smpl_buf_idx + decisionDelay
					for i = 0; i < decisionDelay; i++ {
						last_smple_idx = (last_smple_idx - 1) & (DECISION_DELAY - 1)
						q[i-decisionDelay] = int8((psDD.Q_Q10[last_smple_idx]) >> 10)
						*(*int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(int16(0))*uintptr(i-decisionDelay))) = SKP_SAT16(SKP_RSHIFT_ROUND(SKP_SMULWW(psDD.Xq_Q10[last_smple_idx], psDD.Gain_Q16[last_smple_idx]), 10))
						NSQ.SLTP_shp_Q10[NSQ.SLTP_shp_buf_idx-decisionDelay+i] = psDD.Shape_Q10[last_smple_idx]
					}
					subfr = 0
				}
				start_idx = psEncC.Frame_length - lag - psEncC.PredictLPCOrder - LTP_ORDER/2
				memset(unsafe.Pointer(&FiltState[0]), 0, size_t(uintptr(psEncC.PredictLPCOrder)*unsafe.Sizeof(int32(0))))
				SKP_Silk_MA_Prediction(&NSQ.Xq[start_idx+k*psEncC.Subfr_length], A_Q12, &FiltState[0], &sLTP[start_idx], psEncC.Frame_length-start_idx, psEncC.PredictLPCOrder)
				NSQ.SLTP_buf_idx = psEncC.Frame_length
				NSQ.Rewhite_flag = 1
			}
		}
		SKP_Silk_nsq_del_dec_scale_states(NSQ, psDelDec[:], x, x_sc_Q10[:], subfr_length, sLTP[:], sLTP_Q16[:], k, psEncC.NStatesDelayedDecision, smpl_buf_idx, LTP_scale_Q14, Gains_Q16, psEncCtrlC.PitchL)
		SKP_Silk_noise_shape_quantizer_del_dec(NSQ, psDelDec[:], psEncCtrlC.Sigtype, x_sc_Q10[:], q, ([]int16)(pxq), sLTP_Q16[:], ([]int16)(A_Q12), ([]int16)(B_Q14), ([]int16)(AR_shp_Q13), lag, HarmShapeFIRPacked_Q14, Tilt_Q14[k], LF_shp_Q14[k], Gains_Q16[k], Lambda_Q10, offset_Q10, psEncC.Subfr_length, func() int32 {
			p := &subfr
			x := *p
			*p++
			return x
		}(), psEncC.ShapingLPCOrder, psEncC.PredictLPCOrder, psEncC.Warping_Q16, psEncC.NStatesDelayedDecision, &smpl_buf_idx, decisionDelay)
		x += ([]int16)(psEncC.Subfr_length)
		q += ([]int8)(psEncC.Subfr_length)
		pxq = (*int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(int16(0))*uintptr(psEncC.Subfr_length)))
	}
	RDmin_Q10 = psDelDec[0].RD_Q10
	Winner_ind = 0
	for k = 1; k < psEncC.NStatesDelayedDecision; k++ {
		if psDelDec[k].RD_Q10 < RDmin_Q10 {
			RDmin_Q10 = psDelDec[k].RD_Q10
			Winner_ind = k
		}
	}
	psDD = &psDelDec[Winner_ind]
	psEncCtrlC.Seed = psDD.SeedInit
	last_smple_idx = smpl_buf_idx + decisionDelay
	for i = 0; i < decisionDelay; i++ {
		last_smple_idx = (last_smple_idx - 1) & (DECISION_DELAY - 1)
		q[i-decisionDelay] = int8((psDD.Q_Q10[last_smple_idx]) >> 10)
		*(*int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(int16(0))*uintptr(i-decisionDelay))) = SKP_SAT16(SKP_RSHIFT_ROUND(SKP_SMULWW(psDD.Xq_Q10[last_smple_idx], psDD.Gain_Q16[last_smple_idx]), 10))
		NSQ.SLTP_shp_Q10[NSQ.SLTP_shp_buf_idx-decisionDelay+i] = psDD.Shape_Q10[last_smple_idx]
		sLTP_Q16[NSQ.SLTP_buf_idx-decisionDelay+i] = psDD.Pred_Q16[last_smple_idx]
	}
	memcpy(unsafe.Pointer(&NSQ.SLPC_Q14[0]), unsafe.Pointer(&psDD.SLPC_Q14[psEncC.Subfr_length]), size_t(DECISION_DELAY*unsafe.Sizeof(int32(0))))
	memcpy(unsafe.Pointer(&NSQ.SAR2_Q14[0]), unsafe.Pointer(&psDD.SAR2_Q14[0]), size_t(unsafe.Sizeof([16]int32{})))
	NSQ.SLF_AR_shp_Q12 = psDD.LF_AR_Q12
	NSQ.LagPrev = psEncCtrlC.PitchL[NB_SUBFR-1]
	memcpy(unsafe.Pointer(&NSQ.Xq[0]), unsafe.Pointer(&NSQ.Xq[psEncC.Frame_length]), size_t(uintptr(psEncC.Frame_length)*unsafe.Sizeof(int16(0))))
	memcpy(unsafe.Pointer(&NSQ.SLTP_shp_Q10[0]), unsafe.Pointer(&NSQ.SLTP_shp_Q10[psEncC.Frame_length]), size_t(uintptr(psEncC.Frame_length)*unsafe.Sizeof(int32(0))))
}
func SKP_Silk_noise_shape_quantizer_del_dec(NSQ *SKP_Silk_nsq_state, psDelDec []NSQ_del_dec_struct, sigtype int32, x_Q10 []int32, q []int8, xq []int16, sLTP_Q16 []int32, a_Q12 []int16, b_Q14 []int16, AR_shp_Q13 []int16, lag int32, HarmShapeFIRPacked_Q14 int32, Tilt_Q14 int32, LF_shp_Q14 int32, Gain_Q16 int32, Lambda_Q10 int32, offset_Q10 int32, length int32, subfr int32, shapingLPCOrder int32, predictLPCOrder int32, warping_Q16 int32, nStatesDelayedDecision int32, smpl_buf_idx *int32, decisionDelay int32) {
	var (
		i                 int32
		j                 int32
		k                 int32
		Winner_ind        int32
		RDmin_ind         int32
		RDmax_ind         int32
		last_smple_idx    int32
		Winner_rand_state int32
		LTP_pred_Q14      int32
		LPC_pred_Q10      int32
		n_AR_Q10          int32
		n_LTP_Q14         int32
		n_LF_Q10          int32
		r_Q10             int32
		rr_Q20            int32
		rd1_Q10           int32
		rd2_Q10           int32
		RDmin_Q10         int32
		RDmax_Q10         int32
		q1_Q10            int32
		q2_Q10            int32
		dither            int32
		exc_Q10           int32
		LPC_exc_Q10       int32
		xq_Q10            int32
		tmp1              int32
		tmp2              int32
		sLF_AR_shp_Q10    int32
		pred_lag_ptr      *int32
		shp_lag_ptr       *int32
		psLPC_Q14         *int32
		psSampleState     [4][2]NSQ_sample_struct
		psDD              *NSQ_del_dec_struct
		psSS              *NSQ_sample_struct
	)
	shp_lag_ptr = &NSQ.SLTP_shp_Q10[NSQ.SLTP_shp_buf_idx-lag+HARM_SHAPE_FIR_TAPS/2]
	pred_lag_ptr = &sLTP_Q16[NSQ.SLTP_buf_idx-lag+LTP_ORDER/2]
	for i = 0; i < length; i++ {
		if sigtype == SIG_TYPE_VOICED {
			LTP_pred_Q14 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(int32(0))*0)), int32(b_Q14[0]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*1))), int32(b_Q14[1]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*2))), int32(b_Q14[2]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*3))), int32(b_Q14[3]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*4))), int32(b_Q14[4]))
			pred_lag_ptr = (*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(int32(0))*1))
		} else {
			LTP_pred_Q14 = 0
		}
		if lag > 0 {
			n_LTP_Q14 = SKP_SMULWB((*(*int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(int32(0))*0)))+(*(*int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(int32(0))*2)))), HarmShapeFIRPacked_Q14)
			n_LTP_Q14 = SKP_SMLAWT(n_LTP_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(int32(0))*1))), HarmShapeFIRPacked_Q14)
			n_LTP_Q14 = n_LTP_Q14 << 6
			shp_lag_ptr = (*int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(int32(0))*1))
		} else {
			n_LTP_Q14 = 0
		}
		for k = 0; k < nStatesDelayedDecision; k++ {
			psDD = &psDelDec[k]
			psSS = &psSampleState[k][0]
			psDD.Seed = int32((uint32(psDD.Seed) * 0xBB38435) + 0x3619636B)
			dither = psDD.Seed >> 31
			psLPC_Q14 = &psDD.SLPC_Q14[DECISION_DELAY-1+i]
			LPC_pred_Q10 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), unsafe.Sizeof(int32(0))*0)), int32(a_Q12[0]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*1))), int32(a_Q12[1]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*2))), int32(a_Q12[2]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*3))), int32(a_Q12[3]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*4))), int32(a_Q12[4]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*5))), int32(a_Q12[5]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*6))), int32(a_Q12[6]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*7))), int32(a_Q12[7]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*8))), int32(a_Q12[8]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*9))), int32(a_Q12[9]))
			for j = 10; j < predictLPCOrder; j++ {
				LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, *(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), -int(unsafe.Sizeof(int32(0))*uintptr(j)))), int32(a_Q12[j]))
			}
			tmp2 = SKP_SMLAWB(*(*int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), unsafe.Sizeof(int32(0))*0)), psDD.SAR2_Q14[0], warping_Q16)
			tmp1 = SKP_SMLAWB(psDD.SAR2_Q14[0], psDD.SAR2_Q14[1]-tmp2, warping_Q16)
			psDD.SAR2_Q14[0] = tmp2
			n_AR_Q10 = SKP_SMULWB(tmp2, int32(AR_shp_Q13[0]))
			for j = 2; j < shapingLPCOrder; j += 2 {
				tmp2 = SKP_SMLAWB(psDD.SAR2_Q14[j-1], psDD.SAR2_Q14[j+0]-tmp1, warping_Q16)
				psDD.SAR2_Q14[j-1] = tmp1
				n_AR_Q10 = SKP_SMLAWB(n_AR_Q10, tmp1, int32(AR_shp_Q13[j-1]))
				tmp1 = SKP_SMLAWB(psDD.SAR2_Q14[j+0], psDD.SAR2_Q14[j+1]-tmp2, warping_Q16)
				psDD.SAR2_Q14[j+0] = tmp2
				n_AR_Q10 = SKP_SMLAWB(n_AR_Q10, tmp2, int32(AR_shp_Q13[j]))
			}
			psDD.SAR2_Q14[shapingLPCOrder-1] = tmp1
			n_AR_Q10 = SKP_SMLAWB(n_AR_Q10, tmp1, int32(AR_shp_Q13[shapingLPCOrder-1]))
			n_AR_Q10 = n_AR_Q10 >> 1
			n_AR_Q10 = SKP_SMLAWB(n_AR_Q10, psDD.LF_AR_Q12, Tilt_Q14)
			n_LF_Q10 = SKP_SMULWB(psDD.Shape_Q10[*smpl_buf_idx], LF_shp_Q14) << 2
			n_LF_Q10 = SKP_SMLAWT(n_LF_Q10, psDD.LF_AR_Q12, LF_shp_Q14)
			tmp1 = LTP_pred_Q14 - n_LTP_Q14
			tmp1 = tmp1 >> 4
			tmp1 = tmp1 + LPC_pred_Q10
			tmp1 = tmp1 - n_AR_Q10
			tmp1 = tmp1 - n_LF_Q10
			r_Q10 = (x_Q10[i]) - tmp1
			r_Q10 = (r_Q10 ^ dither) - dither
			r_Q10 = r_Q10 - offset_Q10
			r_Q10 = SKP_LIMIT_32(r_Q10, int32(-64<<10), 64<<10)
			if int64(r_Q10) < -1536 {
				q1_Q10 = SKP_RSHIFT_ROUND(r_Q10, 10) << 10
				r_Q10 = r_Q10 - q1_Q10
				rd1_Q10 = SKP_SMLABB((-(q1_Q10+offset_Q10))*Lambda_Q10, r_Q10, r_Q10) >> 10
				rd2_Q10 = rd1_Q10 + 1024
				rd2_Q10 = rd2_Q10 - (Lambda_Q10 + (r_Q10 << 1))
				q2_Q10 = q1_Q10 + 1024
			} else if r_Q10 > 512 {
				q1_Q10 = SKP_RSHIFT_ROUND(r_Q10, 10) << 10
				r_Q10 = r_Q10 - q1_Q10
				rd1_Q10 = SKP_SMLABB((q1_Q10+offset_Q10)*Lambda_Q10, r_Q10, r_Q10) >> 10
				rd2_Q10 = rd1_Q10 + 1024
				rd2_Q10 = rd2_Q10 - (Lambda_Q10 - (r_Q10 << 1))
				q2_Q10 = q1_Q10 - 1024
			} else {
				rr_Q20 = SKP_SMULBB(offset_Q10, Lambda_Q10)
				rd2_Q10 = SKP_SMLABB(rr_Q20, r_Q10, r_Q10) >> 10
				rd1_Q10 = rd2_Q10 + 1024
				rd1_Q10 = rd1_Q10 + ((Lambda_Q10 + (r_Q10 << 1)) - (rr_Q20 >> 9))
				q1_Q10 = -1024
				q2_Q10 = 0
			}
			if rd1_Q10 < rd2_Q10 {
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).RD_Q10 = psDD.RD_Q10 + rd1_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).RD_Q10 = psDD.RD_Q10 + rd2_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Q_Q10 = q1_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Q_Q10 = q2_Q10
			} else {
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).RD_Q10 = psDD.RD_Q10 + rd2_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).RD_Q10 = psDD.RD_Q10 + rd1_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Q_Q10 = q2_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Q_Q10 = q1_Q10
			}
			exc_Q10 = offset_Q10 + (*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Q_Q10
			exc_Q10 = (exc_Q10 ^ dither) - dither
			LPC_exc_Q10 = exc_Q10 + SKP_RSHIFT_ROUND(LTP_pred_Q14, 4)
			xq_Q10 = LPC_exc_Q10 + LPC_pred_Q10
			sLF_AR_shp_Q10 = xq_Q10 - n_AR_Q10
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).SLTP_shp_Q10 = sLF_AR_shp_Q10 - n_LF_Q10
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).LF_AR_Q12 = sLF_AR_shp_Q10 << 2
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Xq_Q14 = xq_Q10 << 4
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).LPC_exc_Q16 = LPC_exc_Q10 << 6
			exc_Q10 = offset_Q10 + (*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Q_Q10
			exc_Q10 = (exc_Q10 ^ dither) - dither
			LPC_exc_Q10 = exc_Q10 + SKP_RSHIFT_ROUND(LTP_pred_Q14, 4)
			xq_Q10 = LPC_exc_Q10 + LPC_pred_Q10
			sLF_AR_shp_Q10 = xq_Q10 - n_AR_Q10
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).SLTP_shp_Q10 = sLF_AR_shp_Q10 - n_LF_Q10
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).LF_AR_Q12 = sLF_AR_shp_Q10 << 2
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Xq_Q14 = xq_Q10 << 4
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).LPC_exc_Q16 = LPC_exc_Q10 << 6
		}
		*smpl_buf_idx = (*smpl_buf_idx - 1) & (DECISION_DELAY - 1)
		last_smple_idx = (*smpl_buf_idx + decisionDelay) & (DECISION_DELAY - 1)
		RDmin_Q10 = psSampleState[0][0].RD_Q10
		Winner_ind = 0
		for k = 1; k < nStatesDelayedDecision; k++ {
			if psSampleState[k][0].RD_Q10 < RDmin_Q10 {
				RDmin_Q10 = psSampleState[k][0].RD_Q10
				Winner_ind = k
			}
		}
		Winner_rand_state = psDelDec[Winner_ind].RandState[last_smple_idx]
		for k = 0; k < nStatesDelayedDecision; k++ {
			if psDelDec[k].RandState[last_smple_idx] != Winner_rand_state {
				psSampleState[k][0].RD_Q10 = psSampleState[k][0].RD_Q10 + (SKP_int32_MAX >> 4)
				psSampleState[k][1].RD_Q10 = psSampleState[k][1].RD_Q10 + (SKP_int32_MAX >> 4)
			}
		}
		RDmax_Q10 = psSampleState[0][0].RD_Q10
		RDmin_Q10 = psSampleState[0][1].RD_Q10
		RDmax_ind = 0
		RDmin_ind = 0
		for k = 1; k < nStatesDelayedDecision; k++ {
			if psSampleState[k][0].RD_Q10 > RDmax_Q10 {
				RDmax_Q10 = psSampleState[k][0].RD_Q10
				RDmax_ind = k
			}
			if psSampleState[k][1].RD_Q10 < RDmin_Q10 {
				RDmin_Q10 = psSampleState[k][1].RD_Q10
				RDmin_ind = k
			}
		}
		if RDmin_Q10 < RDmax_Q10 {
			SKP_Silk_copy_del_dec_state(&psDelDec[RDmax_ind], &psDelDec[RDmin_ind], i)
			memcpy(unsafe.Pointer(&psSampleState[RDmax_ind][0]), unsafe.Pointer(&psSampleState[RDmin_ind][1]), size_t(unsafe.Sizeof(NSQ_sample_struct{})))
		}
		psDD = &psDelDec[Winner_ind]
		if subfr > 0 || i >= decisionDelay {
			q[i-decisionDelay] = int8((psDD.Q_Q10[last_smple_idx]) >> 10)
			xq[i-decisionDelay] = SKP_SAT16(SKP_RSHIFT_ROUND(SKP_SMULWW(psDD.Xq_Q10[last_smple_idx], psDD.Gain_Q16[last_smple_idx]), 10))
			NSQ.SLTP_shp_Q10[NSQ.SLTP_shp_buf_idx-decisionDelay] = psDD.Shape_Q10[last_smple_idx]
			sLTP_Q16[NSQ.SLTP_buf_idx-decisionDelay] = psDD.Pred_Q16[last_smple_idx]
		}
		NSQ.SLTP_shp_buf_idx++
		NSQ.SLTP_buf_idx++
		for k = 0; k < nStatesDelayedDecision; k++ {
			psDD = &psDelDec[k]
			psSS = &psSampleState[k][0]
			psDD.LF_AR_Q12 = psSS.LF_AR_Q12
			psDD.SLPC_Q14[DECISION_DELAY+i] = psSS.Xq_Q14
			psDD.Xq_Q10[*smpl_buf_idx] = psSS.Xq_Q14 >> 4
			psDD.Q_Q10[*smpl_buf_idx] = psSS.Q_Q10
			psDD.Pred_Q16[*smpl_buf_idx] = psSS.LPC_exc_Q16
			psDD.Shape_Q10[*smpl_buf_idx] = psSS.SLTP_shp_Q10
			psDD.Seed = psDD.Seed + (psSS.Q_Q10 >> 10)
			psDD.RandState[*smpl_buf_idx] = psDD.Seed
			psDD.RD_Q10 = psSS.RD_Q10
			psDD.Gain_Q16[*smpl_buf_idx] = Gain_Q16
		}
	}
	for k = 0; k < nStatesDelayedDecision; k++ {
		psDD = &psDelDec[k]
		memcpy(unsafe.Pointer(&psDD.SLPC_Q14[0]), unsafe.Pointer(&psDD.SLPC_Q14[length]), size_t(DECISION_DELAY*unsafe.Sizeof(int32(0))))
	}
}
func SKP_Silk_nsq_del_dec_scale_states(NSQ *SKP_Silk_nsq_state, psDelDec []NSQ_del_dec_struct, x []int16, x_sc_Q10 []int32, subfr_length int32, sLTP []int16, sLTP_Q16 []int32, subfr int32, nStatesDelayedDecision int32, smpl_buf_idx int32, LTP_scale_Q14 int32, Gains_Q16 [4]int32, pitchL [4]int32) {
	var (
		i            int32
		k            int32
		lag          int32
		inv_gain_Q16 int32
		gain_adj_Q16 int32
		inv_gain_Q32 int32
		psDD         *NSQ_del_dec_struct
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
		for k = 0; k < nStatesDelayedDecision; k++ {
			psDD = &psDelDec[k]
			psDD.LF_AR_Q12 = SKP_SMULWW(gain_adj_Q16, psDD.LF_AR_Q12)
			for i = 0; i < DECISION_DELAY; i++ {
				psDD.SLPC_Q14[i] = SKP_SMULWW(gain_adj_Q16, psDD.SLPC_Q14[i])
			}
			for i = 0; i < MAX_SHAPE_LPC_ORDER; i++ {
				psDD.SAR2_Q14[i] = SKP_SMULWW(gain_adj_Q16, psDD.SAR2_Q14[i])
			}
			for i = 0; i < DECISION_DELAY; i++ {
				psDD.Pred_Q16[i] = SKP_SMULWW(gain_adj_Q16, psDD.Pred_Q16[i])
				psDD.Shape_Q10[i] = SKP_SMULWW(gain_adj_Q16, psDD.Shape_Q10[i])
			}
		}
	}
	for i = 0; i < subfr_length; i++ {
		x_sc_Q10[i] = SKP_SMULBB(int32(x[i]), int32(int16(inv_gain_Q16))) >> 6
	}
	NSQ.Prev_inv_gain_Q16 = inv_gain_Q16
}
func SKP_Silk_copy_del_dec_state(DD_dst *NSQ_del_dec_struct, DD_src *NSQ_del_dec_struct, LPC_state_idx int32) {
	memcpy(unsafe.Pointer(&DD_dst.RandState[0]), unsafe.Pointer(&DD_src.RandState[0]), size_t(unsafe.Sizeof([32]int32{})))
	memcpy(unsafe.Pointer(&DD_dst.Q_Q10[0]), unsafe.Pointer(&DD_src.Q_Q10[0]), size_t(unsafe.Sizeof([32]int32{})))
	memcpy(unsafe.Pointer(&DD_dst.Pred_Q16[0]), unsafe.Pointer(&DD_src.Pred_Q16[0]), size_t(unsafe.Sizeof([32]int32{})))
	memcpy(unsafe.Pointer(&DD_dst.Shape_Q10[0]), unsafe.Pointer(&DD_src.Shape_Q10[0]), size_t(unsafe.Sizeof([32]int32{})))
	memcpy(unsafe.Pointer(&DD_dst.Xq_Q10[0]), unsafe.Pointer(&DD_src.Xq_Q10[0]), size_t(unsafe.Sizeof([32]int32{})))
	memcpy(unsafe.Pointer(&DD_dst.SAR2_Q14[0]), unsafe.Pointer(&DD_src.SAR2_Q14[0]), size_t(unsafe.Sizeof([16]int32{})))
	memcpy(unsafe.Pointer(&DD_dst.SLPC_Q14[LPC_state_idx]), unsafe.Pointer(&DD_src.SLPC_Q14[LPC_state_idx]), size_t(DECISION_DELAY*unsafe.Sizeof(int32(0))))
	DD_dst.LF_AR_Q12 = DD_src.LF_AR_Q12
	DD_dst.Seed = DD_src.Seed
	DD_dst.SeedInit = DD_src.SeedInit
	DD_dst.RD_Q10 = DD_src.RD_Q10
}
