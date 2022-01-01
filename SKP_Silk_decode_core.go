package silk

import "unsafe"

func SKP_Silk_decode_core(psDec *SKP_Silk_decoder_state, psDecCtrl *SKP_Silk_decoder_control, xq []int16, q [480]int32) {
	var (
		i                       int32
		k                       int32
		lag                     int32 = 0
		start_idx               int32
		sLTP_buf_idx            int32
		NLSF_interpolation_flag int32
		sigtype                 int32
		A_Q12                   *int16
		B_Q14                   *int16
		pxq                     *int16
		A_Q12_tmp               [16]int16
		sLTP                    [480]int16
		LTP_pred_Q14            int32
		Gain_Q16                int32
		inv_gain_Q16            int32
		inv_gain_Q32            int32
		gain_adj_Q16            int32
		rand_seed               int32
		offset_Q10              int32
		dither                  int32
		pred_lag_ptr            *int32
		pexc_Q10                *int32
		pres_Q10                *int32
		vec_Q10                 [120]int32
		FiltState               [16]int32
	)
	offset_Q10 = int32(SKP_Silk_Quantization_Offsets_Q10[psDecCtrl.Sigtype][psDecCtrl.QuantOffsetType])
	if int64(psDecCtrl.NLSFInterpCoef_Q2) < (1 << 2) {
		NLSF_interpolation_flag = 1
	} else {
		NLSF_interpolation_flag = 0
	}
	rand_seed = psDecCtrl.Seed
	for i = 0; int64(i) < int64(psDec.Frame_length); i++ {
		rand_seed = int32(int64(uint32(int32(int64(uint32(rand_seed))*0xBB38435))) + 0x3619636B)
		dither = rand_seed >> 31
		psDec.Exc_Q10[i] = int32((int64(q[i]) << 10) + int64(offset_Q10))
		psDec.Exc_Q10[i] = int32((int64(psDec.Exc_Q10[i]) ^ int64(dither)) - int64(dither))
		rand_seed += q[i]
	}
	pexc_Q10 = &psDec.Exc_Q10[0]
	pres_Q10 = &psDec.Res_Q10[0]
	pxq = &psDec.OutBuf[psDec.Frame_length]
	sLTP_buf_idx = psDec.Frame_length
	for k = 0; int64(k) < NB_SUBFR; k++ {
		A_Q12 = &psDecCtrl.PredCoef_Q12[k>>1][0]
		memcpy(unsafe.Pointer(&A_Q12_tmp[0]), unsafe.Pointer(A_Q12), uintptr(psDec.LPC_order)*unsafe.Sizeof(int16(0)))
		B_Q14 = &psDecCtrl.LTPCoef_Q14[int64(k)*LTP_ORDER]
		Gain_Q16 = psDecCtrl.Gains_Q16[k]
		sigtype = psDecCtrl.Sigtype
		inv_gain_Q16 = SKP_INVERSE32_varQ(int32(func() int64 {
			if int64(Gain_Q16) > 1 {
				return int64(Gain_Q16)
			}
			return 1
		}()), 32)
		if int64(inv_gain_Q16) < SKP_int16_MAX {
		} else {
			inv_gain_Q16 = SKP_int16_MAX
		}
		gain_adj_Q16 = 1 << 16
		if int64(inv_gain_Q16) != int64(psDec.Prev_inv_gain_Q16) {
			gain_adj_Q16 = SKP_DIV32_varQ(inv_gain_Q16, psDec.Prev_inv_gain_Q16, 16)
		}
		if int64(psDec.LossCnt) != 0 && int64(psDec.Prev_sigtype) == SIG_TYPE_VOICED && int64(psDecCtrl.Sigtype) == SIG_TYPE_UNVOICED && int64(k) < (NB_SUBFR>>1) {
			memset(unsafe.Pointer(B_Q14), 0, LTP_ORDER*unsafe.Sizeof(int16(0)))
			*(*int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(int16(0))*uintptr(LTP_ORDER/2))) = 1 << 12
			sigtype = SIG_TYPE_VOICED
			psDecCtrl.PitchL[k] = psDec.LagPrev
		}
		if int64(sigtype) == SIG_TYPE_VOICED {
			lag = psDecCtrl.PitchL[k]
			if (int64(k) & (3 - (int64(NLSF_interpolation_flag) << 1))) == 0 {
				start_idx = int32(int64(psDec.Frame_length) - int64(lag) - int64(psDec.LPC_order) - LTP_ORDER/2)
				memset(unsafe.Pointer(&FiltState[0]), 0, uintptr(psDec.LPC_order)*unsafe.Sizeof(int32(0)))
				SKP_Silk_MA_Prediction(&psDec.OutBuf[int64(start_idx)+int64(k)*(int64(psDec.Frame_length)>>2)], A_Q12, &FiltState[0], &sLTP[start_idx], int32(int64(psDec.Frame_length)-int64(start_idx)), psDec.LPC_order)
				inv_gain_Q32 = inv_gain_Q16 << 16
				if k == 0 {
					inv_gain_Q32 = SKP_SMULWB(inv_gain_Q32, psDecCtrl.LTP_scale_Q14) << 2
				}
				for i = 0; int64(i) < (int64(lag) + LTP_ORDER/2); i++ {
					psDec.SLTP_Q16[int64(sLTP_buf_idx)-int64(i)-1] = SKP_SMULWB(inv_gain_Q32, int32(sLTP[int64(psDec.Frame_length)-int64(i)-1]))
				}
			} else {
				if int64(gain_adj_Q16) != 1<<16 {
					for i = 0; int64(i) < (int64(lag) + LTP_ORDER/2); i++ {
						psDec.SLTP_Q16[int64(sLTP_buf_idx)-int64(i)-1] = SKP_SMULWW(gain_adj_Q16, psDec.SLTP_Q16[int64(sLTP_buf_idx)-int64(i)-1])
					}
				}
			}
		}
		for i = 0; int64(i) < MAX_LPC_ORDER; i++ {
			psDec.SLPC_Q14[i] = SKP_SMULWW(gain_adj_Q16, psDec.SLPC_Q14[i])
		}
		psDec.Prev_inv_gain_Q16 = inv_gain_Q16
		if int64(sigtype) == SIG_TYPE_VOICED {
			pred_lag_ptr = &psDec.SLTP_Q16[int64(sLTP_buf_idx)-int64(lag)+LTP_ORDER/2]
			for i = 0; int64(i) < int64(psDec.Subfr_length); i++ {
				LTP_pred_Q14 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(int32(0))*0)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(int16(0))*0))))
				LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*1))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(int16(0))*1))))
				LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*2))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(int16(0))*2))))
				LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*3))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(int16(0))*3))))
				LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*4))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(int16(0))*4))))
				pred_lag_ptr = (*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(int32(0))*1))
				*(*int32)(unsafe.Add(unsafe.Pointer(pres_Q10), unsafe.Sizeof(int32(0))*uintptr(i))) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(pexc_Q10), unsafe.Sizeof(int32(0))*uintptr(i)))) + int64(SKP_RSHIFT_ROUND(LTP_pred_Q14, 4)))
				psDec.SLTP_Q16[sLTP_buf_idx] = *(*int32)(unsafe.Add(unsafe.Pointer(pres_Q10), unsafe.Sizeof(int32(0))*uintptr(i))) << 6
				sLTP_buf_idx++
			}
		} else {
			memcpy(unsafe.Pointer(pres_Q10), unsafe.Pointer(pexc_Q10), uintptr(psDec.Subfr_length)*unsafe.Sizeof(int32(0)))
		}
		length := psDec.Subfr_length
		SKP_Silk_decode_short_term_prediction(vec_Q10[:length], unsafe.Slice(pres_Q10, length), psDec.SLPC_Q14[:], A_Q12_tmp[:], psDec.LPC_order)
		for i = 0; int64(i) < int64(psDec.Subfr_length); i++ {
			*(*int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(int16(0))*uintptr(i))) = SKP_SAT16(SKP_RSHIFT_ROUND(SKP_SMULWW(vec_Q10[i], Gain_Q16), 10))
		}
		memcpy(unsafe.Pointer(&psDec.SLPC_Q14[0]), unsafe.Pointer(&psDec.SLPC_Q14[psDec.Subfr_length]), MAX_LPC_ORDER*unsafe.Sizeof(int32(0)))
		pexc_Q10 = (*int32)(unsafe.Add(unsafe.Pointer(pexc_Q10), unsafe.Sizeof(int32(0))*uintptr(psDec.Subfr_length)))
		pres_Q10 = (*int32)(unsafe.Add(unsafe.Pointer(pres_Q10), unsafe.Sizeof(int32(0))*uintptr(psDec.Subfr_length)))
		pxq = (*int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(int16(0))*uintptr(psDec.Subfr_length)))
	}
	memcpy(unsafe.Pointer(&xq[0]), unsafe.Pointer(&psDec.OutBuf[psDec.Frame_length]), uintptr(psDec.Frame_length)*unsafe.Sizeof(int16(0)))
}

func SKP_Silk_decode_short_term_prediction(vec_Q10 []int32, pres_Q10 []int32, sLPC_Q14 []int32, A_Q12_tmp []int16, LPC_order int32) {
	for i := range vec_Q10 {
		LPC_pred_Q10 := SKP_SMULWB(sLPC_Q14[MAX_LPC_ORDER+int64(i)-1], int32(A_Q12_tmp[0]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-2], int32(A_Q12_tmp[1]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-3], int32(A_Q12_tmp[2]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-4], int32(A_Q12_tmp[3]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-5], int32(A_Q12_tmp[4]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-6], int32(A_Q12_tmp[5]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-7], int32(A_Q12_tmp[6]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-8], int32(A_Q12_tmp[7]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-9], int32(A_Q12_tmp[8]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-10], int32(A_Q12_tmp[9]))
		for j := int32(10); j < LPC_order; j++ {
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+int64(i)-int64(j)-1], int32(A_Q12_tmp[j]))
		}
		vec_Q10[i] = int32(int64(pres_Q10[i]) + int64(LPC_pred_Q10))
		sLPC_Q14[MAX_LPC_ORDER+int64(i)] = vec_Q10[i] << 4
	}
}
