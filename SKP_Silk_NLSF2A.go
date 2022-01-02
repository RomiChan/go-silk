package silk

import "math"

func SKP_Silk_NLSF2A_find_poly(out []int32, cLSF []int32, dd int32) {
	var (
		k    int32
		n    int32
		ftmp int32
	)
	out[0] = 1 << 20
	out[1] = -cLSF[0]
	for k = 1; k < dd; k++ {
		ftmp = cLSF[k*2]
		out[k+1] = ((out[k-1]) << 1) - SKP_RSHIFT_ROUND64(int64(ftmp)*int64(out[k]), 20)
		for n = k; n > 1; n-- {
			out[n] += out[n-2] - SKP_RSHIFT_ROUND64(int64(ftmp)*int64(out[n-1]), 20)
		}
		out[1] -= ftmp
	}
}
func SKP_Silk_NLSF2A(a []int16, NLSF []int32, d int32) {
	var (
		k           int32
		i           int32
		dd          int32
		cos_LSF_Q20 [16]int32
		P           [9]int32
		Q           [9]int32
		Ptmp        int32
		Qtmp        int32
		f_int       int32
		f_frac      int32
		cos_val     int32
		delta       int32
		a_int32     [16]int32
		maxabs      int32
		absval      int32
		idx         int32 = 0
		sc_Q16      int32
	)
	SKP_assert(LSF_COS_TAB_SZ_FIX == 128)
	for k = 0; k < d; k++ {
		SKP_assert(NLSF[k] >= 0)
		SKP_assert(NLSF[k] <= math.MaxInt16)
		f_int = (NLSF[k]) >> (15 - 7)
		f_frac = NLSF[k] - (f_int << (15 - 7))
		SKP_assert(f_int >= 0)
		SKP_assert(f_int < LSF_COS_TAB_SZ_FIX)
		cos_val = SKP_Silk_LSFCosTab_FIX_Q12[f_int]
		delta = SKP_Silk_LSFCosTab_FIX_Q12[f_int+1] - cos_val
		cos_LSF_Q20[k] = (cos_val << 8) + delta*f_frac
	}
	dd = d >> 1
	SKP_Silk_NLSF2A_find_poly(P[:], ([]int32)(&cos_LSF_Q20[0]), dd)
	SKP_Silk_NLSF2A_find_poly(Q[:], ([]int32)(&cos_LSF_Q20[1]), dd)
	for k = 0; k < dd; k++ {
		Ptmp = P[k+1] + P[k]
		Qtmp = Q[k+1] - Q[k]
		a_int32[k] = -SKP_RSHIFT_ROUND(Ptmp+Qtmp, 9)
		a_int32[d-k-1] = SKP_RSHIFT_ROUND(Qtmp-Ptmp, 9)
	}
	for i = 0; i < 10; i++ {
		maxabs = 0
		for k = 0; k < d; k++ {
			absval = int32(SKP_abs(int64(a_int32[k])))
			if absval > maxabs {
				maxabs = absval
				idx = k
			}
		}
		if maxabs > SKP_int16_MAX {
			if maxabs < 98369 {
				maxabs = maxabs
			} else {
				maxabs = 98369
			}
			sc_Q16 = 65470 - ((maxabs-SKP_int16_MAX)*(65470>>2))/((maxabs*(idx+1))>>2)
			SKP_Silk_bwexpander_32(a_int32[:], d, sc_Q16)
		} else {
			break
		}
	}
	if i == 10 {
		SKP_assert(0)
		for k = 0; k < d; k++ {
			a_int32[k] = int32(SKP_SAT16(a_int32[k]))
		}
	}
	for k = 0; k < d; k++ {
		a[k] = int16(a_int32[k])
	}
}
