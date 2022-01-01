package silk

import "unsafe"

func SKP_Silk_residual_energy_FIX(nrgs [4]int32, nrgsQ [4]int32, x []int16, a_Q12 [2][16]int16, gains [4]int32, subfr_length int32, LPC_order int32) {
	var (
		offset      int32
		i           int32
		j           int32
		rshift      int32
		lz1         int32
		lz2         int32
		LPC_res_ptr *int16
		LPC_res     [272]int16
		x_ptr       *int16
		S           [16]int16
		tmp32       int32
	)
	x_ptr = &x[0]
	offset = int32(int64(LPC_order) + int64(subfr_length))
	for i = 0; int64(i) < 2; i++ {
		memset(unsafe.Pointer(&S[0]), 0, size_t(uintptr(LPC_order)*unsafe.Sizeof(int16(0))))
		SKP_Silk_LPC_analysis_filter(x_ptr, &a_Q12[i][0], &S[0], &LPC_res[0], int32((NB_SUBFR>>1)*int64(offset)), LPC_order)
		LPC_res_ptr = &LPC_res[LPC_order]
		for j = 0; int64(j) < (NB_SUBFR >> 1); j++ {
			SKP_Silk_sum_sqr_shift(&nrgs[int64(i)*(NB_SUBFR>>1)+int64(j)], &rshift, LPC_res_ptr, subfr_length)
			nrgsQ[int64(i)*(NB_SUBFR>>1)+int64(j)] = -rshift
			LPC_res_ptr = (*int16)(unsafe.Add(unsafe.Pointer(LPC_res_ptr), unsafe.Sizeof(int16(0))*uintptr(offset)))
		}
		x_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr((NB_SUBFR>>1)*int64(offset))))
	}
	for i = 0; int64(i) < NB_SUBFR; i++ {
		lz1 = int32(int64(SKP_Silk_CLZ32(nrgs[i])) - 1)
		lz2 = int32(int64(SKP_Silk_CLZ32(gains[i])) - 1)
		tmp32 = int32(int64(gains[i]) << int64(lz2))
		tmp32 = SKP_SMMUL(tmp32, tmp32)
		nrgs[i] = SKP_SMMUL(tmp32, int32(int64(nrgs[i])<<int64(lz1)))
		nrgsQ[i] += int32(int64(lz1) + int64(lz2)*2 - 32 - 32)
	}
}
