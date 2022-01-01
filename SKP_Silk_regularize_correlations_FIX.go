package silk

import "unsafe"

func SKP_Silk_regularize_correlations_FIX(XX *int32, xx *int32, noise int32, D int32) {
	var i int32
	for i = 0; i < D; i++ {
		*((*int32)(unsafe.Add(unsafe.Pointer((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*0))), unsafe.Sizeof(int32(0))*uintptr(i*D+i)))) = (*((*int32)(unsafe.Add(unsafe.Pointer((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*0))), unsafe.Sizeof(int32(0))*uintptr(i*D+i))))) + noise
	}
	*(*int32)(unsafe.Add(unsafe.Pointer(xx), unsafe.Sizeof(int32(0))*0)) += noise
}
