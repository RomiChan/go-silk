package silk

import "unsafe"

func SKP_Silk_quant_LTP_gains_FIX(B_Q14 []int16, cbk_index []int32, periodicity_index *int32, W_Q18 []int32, mu_Q8 int32, lowComplexity int32) {
	var (
		j               int32
		k               int32
		temp_idx        [4]int32
		cbk_size        int32
		cl_ptr          *int16
		cbk_ptr_Q14     *int16
		b_Q14_ptr       *int16
		W_Q18_ptr       *int32
		rate_dist_subfr int32
		rate_dist       int32
		min_rate_dist   int32
	)
	min_rate_dist = SKP_int32_MAX
	for k = 0; int64(k) < 3; k++ {
		cl_ptr = SKP_Silk_LTP_gain_BITS_Q6_ptrs[k]
		cbk_ptr_Q14 = SKP_Silk_LTP_vq_ptrs_Q14[k]
		cbk_size = SKP_Silk_LTP_vq_sizes[k]
		W_Q18_ptr = &W_Q18[0]
		b_Q14_ptr = &B_Q14[0]
		rate_dist = 0
		for j = 0; int64(j) < NB_SUBFR; j++ {
			SKP_Silk_VQ_WMat_EC_FIX(&temp_idx[j], &rate_dist_subfr, b_Q14_ptr, W_Q18_ptr, cbk_ptr_Q14, cl_ptr, mu_Q8, cbk_size)
			rate_dist = SKP_ADD_POS_SAT32(rate_dist, rate_dist_subfr)
			b_Q14_ptr = (*int16)(unsafe.Add(unsafe.Pointer(b_Q14_ptr), unsafe.Sizeof(int16(0))*uintptr(LTP_ORDER)))
			W_Q18_ptr = (*int32)(unsafe.Add(unsafe.Pointer(W_Q18_ptr), unsafe.Sizeof(int32(0))*uintptr(LTP_ORDER*LTP_ORDER)))
		}
		if (SKP_int32_MAX - 1) < int64(rate_dist) {
			rate_dist = SKP_int32_MAX - 1
		} else {
			rate_dist = rate_dist
		}
		if int64(rate_dist) < int64(min_rate_dist) {
			min_rate_dist = rate_dist
			memcpy(unsafe.Pointer(&cbk_index[0]), unsafe.Pointer(&temp_idx[0]), size_t(NB_SUBFR*unsafe.Sizeof(int32(0))))
			*periodicity_index = k
		}
		if int64(lowComplexity) != 0 && int64(rate_dist) < int64(SKP_Silk_LTP_gain_middle_avg_RD_Q14) {
			break
		}
	}
	cbk_ptr_Q14 = SKP_Silk_LTP_vq_ptrs_Q14[*periodicity_index]
	for j = 0; int64(j) < NB_SUBFR; j++ {
		for k = 0; int64(k) < LTP_ORDER; k++ {
			B_Q14[int64(j)*LTP_ORDER+int64(k)] = *(*int16)(unsafe.Add(unsafe.Pointer(cbk_ptr_Q14), unsafe.Sizeof(int16(0))*uintptr(int64(k)+int64(cbk_index[j])*LTP_ORDER)))
		}
	}
}
