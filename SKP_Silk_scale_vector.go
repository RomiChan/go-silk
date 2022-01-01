package silk

import "unsafe"

func SKP_Silk_scale_vector32_Q26_lshift_18(data1 *int32, gain_Q26 int32, dataSize int32) {
	var i int32
	for i = 0; int64(i) < int64(dataSize); i++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(int32(0))*uintptr(i))) = int32((int64(*(*int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(int32(0))*uintptr(i)))) * int64(gain_Q26)) >> 8)
	}
}
