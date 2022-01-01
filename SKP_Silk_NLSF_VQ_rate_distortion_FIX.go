package silk

import "unsafe"

func SKP_Silk_NLSF_VQ_rate_distortion_FIX(pRD_Q20 *int32, psNLSF_CBS *SKP_Silk_NLSF_CBS, in_Q15 *int32, w_Q6 *int32, rate_acc_Q5 *int32, mu_Q15 int32, N int32, LPC_order int32) {
	var (
		i           int32
		n           int32
		pRD_vec_Q20 *int32
	)
	SKP_Silk_NLSF_VQ_sum_error_FIX(pRD_Q20, in_Q15, w_Q6, psNLSF_CBS.CB_NLSF_Q15, N, psNLSF_CBS.NVectors, LPC_order)
	pRD_vec_Q20 = pRD_Q20
	for n = 0; int64(n) < int64(N); n++ {
		for i = 0; int64(i) < int64(psNLSF_CBS.NVectors); i++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(pRD_vec_Q20), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_SMLABB(*(*int32)(unsafe.Add(unsafe.Pointer(pRD_vec_Q20), unsafe.Sizeof(int32(0))*uintptr(i))), int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(rate_acc_Q5), unsafe.Sizeof(int32(0))*uintptr(n))))+int64(*(*int16)(unsafe.Add(unsafe.Pointer(psNLSF_CBS.Rates_Q5), unsafe.Sizeof(int16(0))*uintptr(i))))), mu_Q15)
		}
		pRD_vec_Q20 = (*int32)(unsafe.Add(unsafe.Pointer(pRD_vec_Q20), unsafe.Sizeof(int32(0))*uintptr(psNLSF_CBS.NVectors)))
	}
}
