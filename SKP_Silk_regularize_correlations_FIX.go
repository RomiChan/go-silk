package silk

import "unsafe"

func SKP_Silk_regularize_correlations_FIX(XX *int32, xx *int32, noise int32, D int32) {
	var i int32
	for i = 0; int64(i) < int64(D); i++ {
		*((*int32)(unsafe.Add(unsafe.Pointer((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*0))), unsafe.Sizeof(int32(0))*uintptr(int64(i)*int64(D)+int64(i))))) = int32(int64(*((*int32)(unsafe.Add(unsafe.Pointer((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*0))), unsafe.Sizeof(int32(0))*uintptr(int64(i)*int64(D)+int64(i)))))) + int64(noise))
	}
	*(*int32)(unsafe.Add(unsafe.Pointer(xx), unsafe.Sizeof(int32(0))*0)) += noise
}
