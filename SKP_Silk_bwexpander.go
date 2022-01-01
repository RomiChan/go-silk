package silk

import "unsafe"

func SKP_Silk_bwexpander(ar *int16, d int32, chirp_Q16 int32) {
	var (
		i                   int32
		chirp_minus_one_Q16 int32
	)
	chirp_minus_one_Q16 = chirp_Q16 - 0x10000
	for i = 0; i < d-1; i++ {
		*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(SKP_RSHIFT_ROUND(int32(int64(chirp_Q16)*int64(*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(i))))), 16))
		chirp_Q16 += SKP_RSHIFT_ROUND(chirp_Q16*chirp_minus_one_Q16, 16)
	}
	*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(d-1))) = int16(SKP_RSHIFT_ROUND(int32(int64(chirp_Q16)*int64(*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(d-1))))), 16))
}
