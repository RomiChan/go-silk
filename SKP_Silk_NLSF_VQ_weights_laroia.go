package silk

import "unsafe"

const Q_OUT = 6
const MIN_NDELTA = 3

func SKP_Silk_NLSF_VQ_weights_laroia(pNLSFW_Q6 *int32, pNLSF_Q15 *int32, D int32) {
	var (
		k        int32
		tmp1_int int32
		tmp2_int int32
	)
	tmp1_int = SKP_max_int(*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*0)), MIN_NDELTA)
	tmp1_int = (1 << (Q_OUT + 15)) / tmp1_int
	tmp2_int = SKP_max_int(*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*1))-*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*0)), MIN_NDELTA)
	tmp2_int = (1 << (Q_OUT + 15)) / tmp2_int
	*(*int32)(unsafe.Add(unsafe.Pointer(pNLSFW_Q6), unsafe.Sizeof(int32(0))*0)) = SKP_min_int(tmp1_int+tmp2_int, SKP_int16_MAX)
	for k = 1; k < D-1; k += 2 {
		tmp1_int = SKP_max_int(*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(k+1)))-*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(k))), MIN_NDELTA)
		tmp1_int = (1 << (Q_OUT + 15)) / tmp1_int
		*(*int32)(unsafe.Add(unsafe.Pointer(pNLSFW_Q6), unsafe.Sizeof(int32(0))*uintptr(k))) = SKP_min_int(tmp1_int+tmp2_int, SKP_int16_MAX)
		tmp2_int = SKP_max_int(*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(k+2)))-*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(k+1))), MIN_NDELTA)
		tmp2_int = (1 << (Q_OUT + 15)) / tmp2_int
		*(*int32)(unsafe.Add(unsafe.Pointer(pNLSFW_Q6), unsafe.Sizeof(int32(0))*uintptr(k+1))) = SKP_min_int(tmp1_int+tmp2_int, SKP_int16_MAX)
	}
	tmp1_int = SKP_max_int((1<<15)-*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(D-1))), MIN_NDELTA)
	tmp1_int = (1 << (Q_OUT + 15)) / tmp1_int
	*(*int32)(unsafe.Add(unsafe.Pointer(pNLSFW_Q6), unsafe.Sizeof(int32(0))*uintptr(D-1))) = SKP_min_int(tmp1_int+tmp2_int, SKP_int16_MAX)
}
