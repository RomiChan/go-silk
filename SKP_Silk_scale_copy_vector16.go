package silk

import "unsafe"

func SKP_Silk_scale_copy_vector16(data_out *int16, data_in *int16, gain_Q16 int32, dataSize int32) {
	var (
		i     int32
		tmp32 int32
	)
	for i = 0; int64(i) < int64(dataSize); i++ {
		tmp32 = SKP_SMULWB(gain_Q16, int32(*(*int16)(unsafe.Add(unsafe.Pointer(data_in), unsafe.Sizeof(int16(0))*uintptr(i)))))
		*(*int16)(unsafe.Add(unsafe.Pointer(data_out), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(tmp32)
	}
}
