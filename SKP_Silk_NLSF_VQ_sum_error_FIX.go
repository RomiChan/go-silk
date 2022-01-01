package silk

import "unsafe"

func SKP_Silk_NLSF_VQ_sum_error_FIX(err_Q20 []int32, in_Q15 []int32, w_Q6 []int32, pCB_Q15 []int16, N int32, K int32, LPC_order int32) {
	var (
		i          int32
		n          int32
		m          int32
		diff_Q15   int32
		sum_error  int32
		Wtmp_Q6    int32
		Wcpy_Q6    [8]int32
		cb_vec_Q15 *int16
	)
	for m = 0; m < (LPC_order >> 1); m++ {
		Wcpy_Q6[m] = w_Q6[m*2] | (w_Q6[m*2+1])<<16
	}
	for n = 0; n < N; n++ {
		cb_vec_Q15 = &pCB_Q15[0]
		for i = 0; i < K; i++ {
			sum_error = 0
			for m = 0; m < LPC_order; m += 2 {
				Wtmp_Q6 = Wcpy_Q6[m>>1]
				diff_Q15 = int32(int64(in_Q15[m]) - int64(*func() *int16 {
					p := &cb_vec_Q15
					x := *p
					*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
					return x
				}()))
				sum_error = SKP_SMLAWB(sum_error, SKP_SMULBB(diff_Q15, diff_Q15), Wtmp_Q6)
				diff_Q15 = int32(int64(in_Q15[m+1]) - int64(*func() *int16 {
					p := &cb_vec_Q15
					x := *p
					*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
					return x
				}()))
				sum_error = SKP_SMLAWT(sum_error, SKP_SMULBB(diff_Q15, diff_Q15), Wtmp_Q6)
			}
			err_Q20[i] = sum_error
		}
		err_Q20 += ([]int32)(K)
		in_Q15 += ([]int32)(LPC_order)
	}
}
