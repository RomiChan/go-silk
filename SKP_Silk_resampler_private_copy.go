package silk

import "unsafe"

func SKP_Silk_resampler_private_copy(SS unsafe.Pointer, out []int16, in []int16, inLen int32) {
	memcpy(unsafe.Pointer(&out[0]), unsafe.Pointer(&in[0]), size_t(uintptr(inLen)*unsafe.Sizeof(int16(0))))
}
