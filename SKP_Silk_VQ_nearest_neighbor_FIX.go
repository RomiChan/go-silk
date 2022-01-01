package silk

import "unsafe"

func SKP_Silk_VQ_WMat_EC_FIX(ind *int32, rate_dist_Q14 *int32, in_Q14 *int16, W_Q18 *int32, cb_Q14 *int16, cl_Q6 *int16, mu_Q8 int32, L int32) {
	var (
		k          int32
		cb_row_Q14 *int16
		diff_Q14   [5]int16
		sum1_Q14   int32
		sum2_Q16   int32
	)
	*rate_dist_Q14 = SKP_int32_MAX
	cb_row_Q14 = cb_Q14
	for k = 0; int64(k) < int64(L); k++ {
		diff_Q14[0] = int16(int64(*(*int16)(unsafe.Add(unsafe.Pointer(in_Q14), unsafe.Sizeof(int16(0))*0))) - int64(*(*int16)(unsafe.Add(unsafe.Pointer(cb_row_Q14), unsafe.Sizeof(int16(0))*0))))
		diff_Q14[1] = int16(int64(*(*int16)(unsafe.Add(unsafe.Pointer(in_Q14), unsafe.Sizeof(int16(0))*1))) - int64(*(*int16)(unsafe.Add(unsafe.Pointer(cb_row_Q14), unsafe.Sizeof(int16(0))*1))))
		diff_Q14[2] = int16(int64(*(*int16)(unsafe.Add(unsafe.Pointer(in_Q14), unsafe.Sizeof(int16(0))*2))) - int64(*(*int16)(unsafe.Add(unsafe.Pointer(cb_row_Q14), unsafe.Sizeof(int16(0))*2))))
		diff_Q14[3] = int16(int64(*(*int16)(unsafe.Add(unsafe.Pointer(in_Q14), unsafe.Sizeof(int16(0))*3))) - int64(*(*int16)(unsafe.Add(unsafe.Pointer(cb_row_Q14), unsafe.Sizeof(int16(0))*3))))
		diff_Q14[4] = int16(int64(*(*int16)(unsafe.Add(unsafe.Pointer(in_Q14), unsafe.Sizeof(int16(0))*4))) - int64(*(*int16)(unsafe.Add(unsafe.Pointer(cb_row_Q14), unsafe.Sizeof(int16(0))*4))))
		sum1_Q14 = SKP_SMULBB(mu_Q8, int32(*(*int16)(unsafe.Add(unsafe.Pointer(cl_Q6), unsafe.Sizeof(int16(0))*uintptr(k)))))
		sum2_Q16 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*1)), int32(diff_Q14[1]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*2)), int32(diff_Q14[2]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*3)), int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*4)), int32(diff_Q14[4]))
		sum2_Q16 = int32(int64(sum2_Q16) << 1)
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*0)), int32(diff_Q14[0]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[0]))
		sum2_Q16 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*7)), int32(diff_Q14[2]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*8)), int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*9)), int32(diff_Q14[4]))
		sum2_Q16 = int32(int64(sum2_Q16) << 1)
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*6)), int32(diff_Q14[1]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[1]))
		sum2_Q16 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*13)), int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*14)), int32(diff_Q14[4]))
		sum2_Q16 = int32(int64(sum2_Q16) << 1)
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*12)), int32(diff_Q14[2]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[2]))
		sum2_Q16 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*19)), int32(diff_Q14[4]))
		sum2_Q16 = int32(int64(sum2_Q16) << 1)
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, *(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*18)), int32(diff_Q14[3]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(W_Q18), unsafe.Sizeof(int32(0))*24)), int32(diff_Q14[4]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[4]))
		if int64(sum1_Q14) < int64(*rate_dist_Q14) {
			*rate_dist_Q14 = sum1_Q14
			*ind = k
		}
		cb_row_Q14 = (*int16)(unsafe.Add(unsafe.Pointer(cb_row_Q14), unsafe.Sizeof(int16(0))*uintptr(LTP_ORDER)))
	}
}
