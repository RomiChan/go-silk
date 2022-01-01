package silk

import "unsafe"

func SKP_Silk_bwexpander_32(ar *int32, d int32, chirp_Q16 int32) {
	var (
		i             int32
		tmp_chirp_Q16 int32
	)
	tmp_chirp_Q16 = chirp_Q16
	for i = 0; i < d-1; i++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_SMULWW(*(*int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int32(0))*uintptr(i))), tmp_chirp_Q16)
		tmp_chirp_Q16 = SKP_SMULWW(chirp_Q16, tmp_chirp_Q16)
	}
	*(*int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int32(0))*uintptr(d-1))) = SKP_SMULWW(*(*int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int32(0))*uintptr(d-1))), tmp_chirp_Q16)
}
