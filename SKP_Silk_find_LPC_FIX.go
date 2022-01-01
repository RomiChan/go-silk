package silk

import "unsafe"

func SKP_Silk_find_LPC_FIX(NLSF_Q15 []int32, interpIndex *int32, prev_NLSFq_Q15 []int32, useInterpolatedNLSFs int32, LPC_order int32, x []int16, subfr_length int32) {
	var (
		k                int32
		a_Q16            [16]int32
		isInterpLower    int32
		shift            int32
		S                [16]int16
		res_nrg0         int32
		res_nrg1         int32
		rshift0          int32
		rshift1          int32
		a_tmp_Q16        [16]int32
		res_nrg_interp   int32
		res_nrg          int32
		res_tmp_nrg      int32
		res_nrg_interp_Q int32
		res_nrg_Q        int32
		res_tmp_nrg_Q    int32
		a_tmp_Q12        [16]int16
		NLSF0_Q15        [16]int32
		LPC_res          [272]int16
	)
	*interpIndex = 4
	SKP_Silk_burg_modified(&res_nrg, &res_nrg_Q, a_Q16[:], x, subfr_length, NB_SUBFR, SKP_FIX_CONST(FIND_LPC_COND_FAC, 32), LPC_order)
	SKP_Silk_bwexpander_32(&a_Q16[0], LPC_order, SKP_FIX_CONST(FIND_LPC_CHIRP, 16))
	if int64(useInterpolatedNLSFs) == 1 {
		SKP_Silk_burg_modified(&res_tmp_nrg, &res_tmp_nrg_Q, a_tmp_Q16[:], ([]int16)(&x[(NB_SUBFR>>1)*int64(subfr_length)]), subfr_length, NB_SUBFR>>1, SKP_FIX_CONST(FIND_LPC_COND_FAC, 32), LPC_order)
		SKP_Silk_bwexpander_32(&a_tmp_Q16[0], LPC_order, SKP_FIX_CONST(FIND_LPC_CHIRP, 16))
		shift = int32(int64(res_tmp_nrg_Q) - int64(res_nrg_Q))
		if int64(shift) >= 0 {
			if int64(shift) < 32 {
				res_nrg = int32(int64(res_nrg) - (int64(res_tmp_nrg) >> int64(shift)))
			}
		} else {
			res_nrg = int32((int64(res_nrg) >> int64(-shift)) - int64(res_tmp_nrg))
			res_nrg_Q = res_tmp_nrg_Q
		}
		SKP_Silk_A2NLSF(NLSF_Q15, a_tmp_Q16[:], LPC_order)
		for k = 3; int64(k) >= 0; k-- {
			SKP_Silk_interpolate(NLSF0_Q15, ([16]int32)(prev_NLSFq_Q15), ([16]int32)(NLSF_Q15), k, LPC_order)
			SKP_Silk_NLSF2A_stable(a_tmp_Q12, NLSF0_Q15, LPC_order)
			memset(unsafe.Pointer(&S[0]), 0, size_t(uintptr(LPC_order)*unsafe.Sizeof(int16(0))))
			SKP_Silk_LPC_analysis_filter(&x[0], &a_tmp_Q12[0], &S[0], &LPC_res[0], int32(int64(subfr_length)*2), LPC_order)
			SKP_Silk_sum_sqr_shift(&res_nrg0, &rshift0, &LPC_res[LPC_order], int32(int64(subfr_length)-int64(LPC_order)))
			SKP_Silk_sum_sqr_shift(&res_nrg1, &rshift1, &LPC_res[int64(LPC_order)+int64(subfr_length)], int32(int64(subfr_length)-int64(LPC_order)))
			shift = int32(int64(rshift0) - int64(rshift1))
			if int64(shift) >= 0 {
				res_nrg1 = res_nrg1 >> int64(shift)
				res_nrg_interp_Q = -rshift0
			} else {
				res_nrg0 = res_nrg0 >> int64(-shift)
				res_nrg_interp_Q = -rshift1
			}
			res_nrg_interp = int32(int64(res_nrg0) + int64(res_nrg1))
			shift = int32(int64(res_nrg_interp_Q) - int64(res_nrg_Q))
			if int64(shift) >= 0 {
				if (int64(res_nrg_interp) >> int64(shift)) < int64(res_nrg) {
					isInterpLower = SKP_TRUE
				} else {
					isInterpLower = SKP_FALSE
				}
			} else {
				if int64(-shift) < 32 {
					if int64(res_nrg_interp) < (int64(res_nrg) >> int64(-shift)) {
						isInterpLower = SKP_TRUE
					} else {
						isInterpLower = SKP_FALSE
					}
				} else {
					isInterpLower = SKP_FALSE
				}
			}
			if int64(isInterpLower) == SKP_TRUE {
				res_nrg = res_nrg_interp
				res_nrg_Q = res_nrg_interp_Q
				*interpIndex = k
			}
		}
	}
	if int64(*interpIndex) == 4 {
		SKP_Silk_A2NLSF(NLSF_Q15, a_Q16[:], LPC_order)
	}
}
