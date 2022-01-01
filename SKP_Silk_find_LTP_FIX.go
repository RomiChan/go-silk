package silk

import "unsafe"

const LTP_CORRS_HEAD_ROOM = 2

func SKP_Silk_find_LTP_FIX(b_Q14 [20]int16, WLTP [100]int32, LTPredCodGain_Q7 *int32, r_first []int16, r_last []int16, lag [4]int32, Wght_Q15 [4]int32, subfr_length int32, mem_offset int32, corr_rshifts [4]int32) {
	var (
		i                int32
		k                int32
		lshift           int32
		r_ptr            *int16
		lag_ptr          *int16
		b_Q14_ptr        *int16
		regu             int32
		WLTP_ptr         *int32
		b_Q16            [5]int32
		delta_b_Q14      [5]int32
		d_Q14            [4]int32
		nrg              [4]int32
		g_Q26            int32
		w                [4]int32
		WLTP_max         int32
		max_abs_d_Q14    int32
		max_w_bits       int32
		temp32           int32
		denom32          int32
		extra_shifts     int32
		rr_shifts        int32
		maxRshifts       int32
		maxRshifts_wxtra int32
		LZs              int32
		LPC_res_nrg      int32
		LPC_LTP_res_nrg  int32
		div_Q16          int32
		Rr               [5]int32
		rr               [4]int32
		wd               int32
		m_Q12            int32
	)
	b_Q14_ptr = &b_Q14[0]
	WLTP_ptr = &WLTP[0]
	r_ptr = &r_first[mem_offset]
	for k = 0; k < NB_SUBFR; k++ {
		if k == (NB_SUBFR >> 1) {
			r_ptr = &r_last[mem_offset]
		}
		lag_ptr = (*int16)(unsafe.Add(unsafe.Pointer(r_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(lag[k]+LTP_ORDER/2))))
		SKP_Silk_sum_sqr_shift(&rr[k], &rr_shifts, r_ptr, subfr_length)
		LZs = SKP_Silk_CLZ32(rr[k])
		if LZs < LTP_CORRS_HEAD_ROOM {
			rr[k] = SKP_RSHIFT_ROUND(rr[k], LTP_CORRS_HEAD_ROOM-LZs)
			rr_shifts += LTP_CORRS_HEAD_ROOM - LZs
		}
		corr_rshifts[k] = rr_shifts
		SKP_Silk_corrMatrix_FIX(lag_ptr, subfr_length, LTP_ORDER, LTP_CORRS_HEAD_ROOM, WLTP_ptr, &corr_rshifts[k])
		SKP_Silk_corrVector_FIX(lag_ptr, r_ptr, subfr_length, LTP_ORDER, &Rr[0], corr_rshifts[k])
		if corr_rshifts[k] > rr_shifts {
			rr[k] = (rr[k]) >> (corr_rshifts[k] - rr_shifts)
		}
		regu = 1
		regu = SKP_SMLAWB(regu, rr[k], SKP_FIX_CONST(0.01/3, 16))
		regu = SKP_SMLAWB(regu, *((*int32)(unsafe.Add(unsafe.Pointer(WLTP_ptr), unsafe.Sizeof(int32(0))*uintptr(LTP_ORDER*0+0)))), SKP_FIX_CONST(0.01/3, 16))
		regu = SKP_SMLAWB(regu, *((*int32)(unsafe.Add(unsafe.Pointer(WLTP_ptr), unsafe.Sizeof(int32(0))*uintptr((LTP_ORDER-1)*LTP_ORDER+(LTP_ORDER-1))))), SKP_FIX_CONST(0.01/3, 16))
		SKP_Silk_regularize_correlations_FIX(WLTP_ptr, &rr[k], regu, LTP_ORDER)
		SKP_Silk_solve_LDL_FIX(WLTP_ptr, LTP_ORDER, &Rr[0], &b_Q16[0])
		SKP_Silk_fit_LTP(b_Q16[:], ([]int16)(b_Q14_ptr))
		nrg[k] = SKP_Silk_residual_energy16_covar_FIX(b_Q14_ptr, WLTP_ptr, &Rr[0], rr[k], LTP_ORDER, 14)
		extra_shifts = SKP_min_int(corr_rshifts[k], LTP_CORRS_HEAD_ROOM)
		denom32 = SKP_LSHIFT_SAT32(SKP_SMULWB(nrg[k], Wght_Q15[k]), extra_shifts+1) + (SKP_SMULWB(subfr_length, 655) >> (corr_rshifts[k] - extra_shifts))
		if denom32 > 1 {
			denom32 = denom32
		} else {
			denom32 = 1
		}
		temp32 = ((Wght_Q15[k]) << 16) / denom32
		temp32 = temp32 >> (corr_rshifts[k] + 31 - extra_shifts - 26)
		WLTP_max = 0
		for i = 0; i < LTP_ORDER*LTP_ORDER; i++ {
			if (*(*int32)(unsafe.Add(unsafe.Pointer(WLTP_ptr), unsafe.Sizeof(int32(0))*uintptr(i)))) > WLTP_max {
				WLTP_max = *(*int32)(unsafe.Add(unsafe.Pointer(WLTP_ptr), unsafe.Sizeof(int32(0))*uintptr(i)))
			} else {
				WLTP_max = WLTP_max
			}
		}
		lshift = SKP_Silk_CLZ32(WLTP_max) - 1 - 3
		if lshift+(26-18) < 31 {
			temp32 = SKP_min_32(temp32, 1<<(lshift+(26-18)))
		}
		SKP_Silk_scale_vector32_Q26_lshift_18(WLTP_ptr, temp32, LTP_ORDER*LTP_ORDER)
		w[k] = *((*int32)(unsafe.Add(unsafe.Pointer(WLTP_ptr), unsafe.Sizeof(int32(0))*uintptr((LTP_ORDER>>1)*LTP_ORDER+(LTP_ORDER>>1)))))
		r_ptr = (*int16)(unsafe.Add(unsafe.Pointer(r_ptr), unsafe.Sizeof(int16(0))*uintptr(subfr_length)))
		b_Q14_ptr = (*int16)(unsafe.Add(unsafe.Pointer(b_Q14_ptr), unsafe.Sizeof(int16(0))*uintptr(LTP_ORDER)))
		WLTP_ptr = (*int32)(unsafe.Add(unsafe.Pointer(WLTP_ptr), unsafe.Sizeof(int32(0))*uintptr(LTP_ORDER*LTP_ORDER)))
	}
	maxRshifts = 0
	for k = 0; k < NB_SUBFR; k++ {
		maxRshifts = SKP_max_int(corr_rshifts[k], maxRshifts)
	}
	if LTPredCodGain_Q7 != nil {
		LPC_LTP_res_nrg = 0
		LPC_res_nrg = 0
		for k = 0; k < NB_SUBFR; k++ {
			LPC_res_nrg = LPC_res_nrg + ((SKP_SMULWB(rr[k], Wght_Q15[k]) + 1) >> ((maxRshifts - corr_rshifts[k]) + 1))
			LPC_LTP_res_nrg = LPC_LTP_res_nrg + ((SKP_SMULWB(nrg[k], Wght_Q15[k]) + 1) >> ((maxRshifts - corr_rshifts[k]) + 1))
		}
		if LPC_LTP_res_nrg > 1 {
			LPC_LTP_res_nrg = LPC_LTP_res_nrg
		} else {
			LPC_LTP_res_nrg = 1
		}
		div_Q16 = SKP_DIV32_varQ(LPC_res_nrg, LPC_LTP_res_nrg, 16)
		*LTPredCodGain_Q7 = SKP_SMULBB(3, SKP_Silk_lin2log(div_Q16)-(16<<7))
	}
	b_Q14_ptr = &b_Q14[0]
	for k = 0; k < NB_SUBFR; k++ {
		d_Q14[k] = 0
		for i = 0; i < LTP_ORDER; i++ {
			d_Q14[k] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(b_Q14_ptr), unsafe.Sizeof(int16(0))*uintptr(i))))
		}
		b_Q14_ptr = (*int16)(unsafe.Add(unsafe.Pointer(b_Q14_ptr), unsafe.Sizeof(int16(0))*uintptr(LTP_ORDER)))
	}
	max_abs_d_Q14 = 0
	max_w_bits = 0
	for k = 0; k < NB_SUBFR; k++ {
		max_abs_d_Q14 = SKP_max_32(max_abs_d_Q14, int32(SKP_abs(int64(d_Q14[k]))))
		max_w_bits = SKP_max_32(max_w_bits, 32-SKP_Silk_CLZ32(w[k])+corr_rshifts[k]-maxRshifts)
	}
	extra_shifts = max_w_bits + 32 - SKP_Silk_CLZ32(max_abs_d_Q14) - 14
	extra_shifts -= maxRshifts + (32 - 1 - 2)
	extra_shifts = SKP_max_int(extra_shifts, 0)
	maxRshifts_wxtra = maxRshifts + extra_shifts
	temp32 = (262 >> (maxRshifts + extra_shifts)) + 1
	wd = 0
	for k = 0; k < NB_SUBFR; k++ {
		temp32 = temp32 + ((w[k]) >> (maxRshifts_wxtra - corr_rshifts[k]))
		wd = wd + (SKP_SMULWW((w[k])>>(maxRshifts_wxtra-corr_rshifts[k]), d_Q14[k]) << 2)
	}
	m_Q12 = SKP_DIV32_varQ(wd, temp32, 12)
	b_Q14_ptr = &b_Q14[0]
	for k = 0; k < NB_SUBFR; k++ {
		if 2-corr_rshifts[k] > 0 {
			temp32 = (w[k]) >> (2 - corr_rshifts[k])
		} else {
			temp32 = SKP_LSHIFT_SAT32(w[k], corr_rshifts[k]-2)
		}
		g_Q26 = (SKP_FIX_CONST(0.1, 26) / ((SKP_FIX_CONST(0.1, 26) >> 10) + temp32)) * SKP_LSHIFT_SAT32(func() int32 {
			if ((m_Q12 - ((d_Q14[k]) >> 2)) & math.MinInt32) == 0 {
				if (m_Q12 & (((d_Q14[k]) >> 2) ^ math.MinInt32) & math.MinInt32) != 0 {
					return math.MinInt32
				}
				return m_Q12 - ((d_Q14[k]) >> 2)
			}
			if ((m_Q12 ^ math.MinInt32) & ((d_Q14[k]) >> 2) & math.MinInt32) != 0 {
				return SKP_int32_MAX
			}
			return m_Q12 - ((d_Q14[k]) >> 2)
		}(), 4)
		temp32 = 0
		for i = 0; i < LTP_ORDER; i++ {
			delta_b_Q14[i] = int32(SKP_max_16(*(*int16)(unsafe.Add(unsafe.Pointer(b_Q14_ptr), unsafe.Sizeof(int16(0))*uintptr(i))), 1638))
			temp32 += delta_b_Q14[i]
		}
		temp32 = g_Q26 / temp32
		for i = 0; i < LTP_ORDER; i++ {
			*(*int16)(unsafe.Add(unsafe.Pointer(b_Q14_ptr), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(SKP_LIMIT_32(int32(*(*int16)(unsafe.Add(unsafe.Pointer(b_Q14_ptr), unsafe.Sizeof(int16(0))*uintptr(i))))+SKP_SMULWB(SKP_LSHIFT_SAT32(temp32, 4), delta_b_Q14[i]), -16000, 28000))
		}
		b_Q14_ptr = (*int16)(unsafe.Add(unsafe.Pointer(b_Q14_ptr), unsafe.Sizeof(int16(0))*uintptr(LTP_ORDER)))
	}
}
func SKP_Silk_fit_LTP(LTP_coefs_Q16 []int32, LTP_coefs_Q14 []int16) {
	var i int32
	for i = 0; i < LTP_ORDER; i++ {
		LTP_coefs_Q14[i] = SKP_SAT16(SKP_RSHIFT_ROUND(LTP_coefs_Q16[i], 2))
	}
}
