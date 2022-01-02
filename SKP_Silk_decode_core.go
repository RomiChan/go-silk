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
		B_Q14                   []int16
		pxq                     []int16
	)
	_ = pxq
	var A_Q12_tmp [16]int16
	var sLTP [480]int16
	var LTP_pred_Q14 int32
	var Gain_Q16 int32
	var inv_gain_Q16 int32
	var inv_gain_Q32 int32
	var gain_adj_Q16 int32
	var rand_seed int32
	var offset_Q10 int32
	var dither int32
	var pred_lag_ptr []int32
	var pexc_Q10 []int32
	var pres_Q10 []int32
	var vec_Q10 [120]int32
	var FiltState [16]int32
	offset_Q10 = int32(SKP_Silk_Quantization_Offsets_Q10[psDecCtrl.Sigtype][psDecCtrl.QuantOffsetType])
	if psDecCtrl.NLSFInterpCoef_Q2 < (1 << 2) {
		NLSF_interpolation_flag = 1
	} else {
		NLSF_interpolation_flag = 0
	}
	rand_seed = psDecCtrl.Seed
	for i = 0; i < psDec.Frame_length; i++ {
		rand_seed = int32((uint32(rand_seed) * 0xBB38435) + 0x3619636B)
		dither = rand_seed >> 31
		psDec.Exc_Q10[i] = ((q[i]) << 10) + offset_Q10
		psDec.Exc_Q10[i] = (psDec.Exc_Q10[i] ^ dither) - dither
		rand_seed += q[i]
	}
	pexc_Q10 = ([]int32)(psDec.Exc_Q10[:])
	pres_Q10 = ([]int32)(psDec.Res_Q10[:])
	pxq = ([]int16)(&psDec.OutBuf[psDec.Frame_length])
	sLTP_buf_idx = psDec.Frame_length
	for k = 0; k < NB_SUBFR; k++ {
		A_Q12 = &psDecCtrl.PredCoef_Q12[k>>1][0]
		memcpy(unsafe.Pointer(&A_Q12_tmp[0]), unsafe.Pointer(A_Q12), size_t(uintptr(psDec.LPC_order)*unsafe.Sizeof(int16(0))))
		B_Q14 = ([]int16)(&psDecCtrl.LTPCoef_Q14[k*LTP_ORDER])
		Gain_Q16 = psDecCtrl.Gains_Q16[k]
		sigtype = psDecCtrl.Sigtype
		inv_gain_Q16 = SKP_INVERSE32_varQ(func() int32 {
			if Gain_Q16 > 1 {
				return Gain_Q16
			}
			return 1
		}(), 32)
		if inv_gain_Q16 < SKP_int16_MAX {
			inv_gain_Q16 = inv_gain_Q16
		} else {
			inv_gain_Q16 = SKP_int16_MAX
		}
		gain_adj_Q16 = 1 << 16
		if inv_gain_Q16 != psDec.Prev_inv_gain_Q16 {
			gain_adj_Q16 = SKP_DIV32_varQ(inv_gain_Q16, psDec.Prev_inv_gain_Q16, 16)
		}
		if psDec.LossCnt != 0 && psDec.Prev_sigtype == SIG_TYPE_VOICED && psDecCtrl.Sigtype == SIG_TYPE_UNVOICED && k < (NB_SUBFR>>1) {
			memset(unsafe.Pointer(&B_Q14[0]), 0, size_t(LTP_ORDER*unsafe.Sizeof(int16(0))))
			B_Q14[LTP_ORDER/2] = 1 << 12
			sigtype = SIG_TYPE_VOICED
			psDecCtrl.PitchL[k] = psDec.LagPrev
		}
		if sigtype == SIG_TYPE_VOICED {
			lag = psDecCtrl.PitchL[k]
			if (k & (3 - (NLSF_interpolation_flag << 1))) == 0 {
				start_idx = psDec.Frame_length - lag - psDec.LPC_order - LTP_ORDER/2
				memset(unsafe.Pointer(&FiltState[0]), 0, size_t(uintptr(psDec.LPC_order)*unsafe.Sizeof(int32(0))))
				SKP_Silk_MA_Prediction(([]int16)(&psDec.OutBuf[start_idx+k*(psDec.Frame_length>>2)]), ([]int16)(A_Q12), FiltState[:], ([]int16)(&sLTP[start_idx]), psDec.Frame_length-start_idx, psDec.LPC_order)
				inv_gain_Q32 = inv_gain_Q16 << 16
				if k == 0 {
					inv_gain_Q32 = SKP_SMULWB(inv_gain_Q32, psDecCtrl.LTP_scale_Q14) << 2
				}
				for i = 0; i < (lag + LTP_ORDER/2); i++ {
					psDec.SLTP_Q16[sLTP_buf_idx-i-1] = SKP_SMULWB(inv_gain_Q32, int32(sLTP[psDec.Frame_length-i-1]))
				}
			} else {
				if gain_adj_Q16 != 1<<16 {
					for i = 0; i < (lag + LTP_ORDER/2); i++ {
						psDec.SLTP_Q16[sLTP_buf_idx-i-1] = SKP_SMULWW(gain_adj_Q16, psDec.SLTP_Q16[sLTP_buf_idx-i-1])
					}
				}
			}
		}
		for i = 0; i < MAX_LPC_ORDER; i++ {
			psDec.SLPC_Q14[i] = SKP_SMULWW(gain_adj_Q16, psDec.SLPC_Q14[i])
		}
		psDec.Prev_inv_gain_Q16 = inv_gain_Q16
		if sigtype == SIG_TYPE_VOICED {
			pred_lag_ptr = ([]int32)(&psDec.SLTP_Q16[sLTP_buf_idx-lag+LTP_ORDER/2])
			for i = 0; i < psDec.Subfr_length; i++ {
				LTP_pred_Q14 = SKP_SMULWB(pred_lag_ptr[0], int32(B_Q14[0]))
				LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, pred_lag_ptr[-1], int32(B_Q14[1]))
				LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, pred_lag_ptr[-2], int32(B_Q14[2]))
				LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, pred_lag_ptr[-3], int32(B_Q14[3]))
				LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, pred_lag_ptr[-4], int32(B_Q14[4]))
				pred_lag_ptr++
				pres_Q10[i] = (pexc_Q10[i]) + SKP_RSHIFT_ROUND(LTP_pred_Q14, 4)
				psDec.SLTP_Q16[sLTP_buf_idx] = (pres_Q10[i]) << 6
				sLTP_buf_idx++
			}
		} else {
			memcpy(unsafe.Pointer(&pres_Q10[0]), unsafe.Pointer(&pexc_Q10[0]), size_t(uintptr(psDec.Subfr_length)*unsafe.Sizeof(int32(0))))
		}
		SKP_Silk_decode_short_term_prediction(vec_Q10[:], pres_Q10, psDec.SLPC_Q14[:], A_Q12_tmp[:], psDec.LPC_order, psDec.Subfr_length)
		for i = 0; i < psDec.Subfr_length; i++ {
			pxq[i] = SKP_SAT16(SKP_RSHIFT_ROUND(SKP_SMULWW(vec_Q10[i], Gain_Q16), 10))
		}
		memcpy(unsafe.Pointer(&psDec.SLPC_Q14[0]), unsafe.Pointer(&psDec.SLPC_Q14[psDec.Subfr_length]), size_t(MAX_LPC_ORDER*unsafe.Sizeof(int32(0))))
		pexc_Q10 += ([]int32)(psDec.Subfr_length)
		pres_Q10 += ([]int32)(psDec.Subfr_length)
		pxq += ([]int16)(psDec.Subfr_length)
	}
	memcpy(unsafe.Pointer(&xq[0]), unsafe.Pointer(&psDec.OutBuf[psDec.Frame_length]), size_t(uintptr(psDec.Frame_length)*unsafe.Sizeof(int16(0))))
}
func SKP_Silk_decode_short_term_prediction(vec_Q10 []int32, pres_Q10 []int32, sLPC_Q14 []int32, A_Q12_tmp []int16, LPC_order int32, subfr_length int32) {
	var (
		i            int32
		LPC_pred_Q10 int32
		j            int32
	)
	for i = 0; i < subfr_length; i++ {
		LPC_pred_Q10 = SKP_SMULWB(sLPC_Q14[MAX_LPC_ORDER+i-1], int32(A_Q12_tmp[0]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-2], int32(A_Q12_tmp[1]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-3], int32(A_Q12_tmp[2]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-4], int32(A_Q12_tmp[3]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-5], int32(A_Q12_tmp[4]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-6], int32(A_Q12_tmp[5]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-7], int32(A_Q12_tmp[6]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-8], int32(A_Q12_tmp[7]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-9], int32(A_Q12_tmp[8]))
		LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-10], int32(A_Q12_tmp[9]))
		for j = 10; j < LPC_order; j++ {
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, sLPC_Q14[MAX_LPC_ORDER+i-j-1], int32(A_Q12_tmp[j]))
		}
		vec_Q10[i] = (pres_Q10[i]) + LPC_pred_Q10
		sLPC_Q14[MAX_LPC_ORDER+i] = (vec_Q10[i]) << 4
	}
}
