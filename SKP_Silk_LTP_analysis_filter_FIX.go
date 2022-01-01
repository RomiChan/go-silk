package silk

import "unsafe"

func SKP_Silk_LTP_analysis_filter_FIX(LTP_res *int16, x *int16, LTPCoef_Q14 [20]int16, pitchL [4]int32, invGains_Q16 [4]int32, subfr_length int32, pre_length int32) {
	var (
		x_ptr       *int16
		x_lag_ptr   *int16
		Btmp_Q14    [5]int16
		LTP_res_ptr *int16
		k           int32
		i           int32
		j           int32
		LTP_est     int32
	)
	x_ptr = x
	LTP_res_ptr = LTP_res
	for k = 0; int64(k) < NB_SUBFR; k++ {
		x_lag_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(pitchL[k]))))
		for i = 0; int64(i) < LTP_ORDER; i++ {
			Btmp_Q14[i] = LTPCoef_Q14[int64(k)*LTP_ORDER+int64(i)]
		}
		for i = 0; int64(i) < int64(subfr_length)+int64(pre_length); i++ {
			*(*int16)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(int16(0))*uintptr(i))) = *(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(i)))
			LTP_est = SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_lag_ptr), unsafe.Sizeof(int16(0))*uintptr(LTP_ORDER/2)))), int32(Btmp_Q14[0]))
			for j = 1; int64(j) < LTP_ORDER; j++ {
				LTP_est = int32(int64(uint32(LTP_est)) + int64(uint32(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_lag_ptr), unsafe.Sizeof(int16(0))*uintptr(LTP_ORDER/2-int64(j))))), int32(Btmp_Q14[j])))))
			}
			LTP_est = SKP_RSHIFT_ROUND(LTP_est, 14)
			*(*int16)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(int16(0))*uintptr(i))) = SKP_SAT16(int16(int64(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(i)))) - int64(LTP_est)))
			*(*int16)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(SKP_SMULWB(invGains_Q16[k], int32(*(*int16)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(int16(0))*uintptr(i))))))
			x_lag_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_lag_ptr), unsafe.Sizeof(int16(0))*1))
		}
		LTP_res_ptr = (*int16)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)+int64(pre_length))))
		x_ptr = (*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(subfr_length)))
	}
}
