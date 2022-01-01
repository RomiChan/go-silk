package silk

import "unsafe"

func SKP_Silk_resampler_up2(S *int32, out *int16, in *int16, len_ int32) {
	var (
		k     int32
		in32  int32
		out32 int32
		Y     int32
		X     int32
	)
	for k = 0; int64(k) < int64(len_); k++ {
		in32 = int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))) << 10)
		Y = int32(int64(in32) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))))
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_up2_lq_0))
		out32 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + int64(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int64(in32) + int64(X))
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int64(k)*2))) = SKP_SAT16(int16(SKP_RSHIFT_ROUND(out32, 10)))
		Y = int32(int64(in32) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))))
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_up2_lq_1))
		out32 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + int64(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int64(in32) + int64(X))
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int64(k)*2+1))) = SKP_SAT16(int16(SKP_RSHIFT_ROUND(out32, 10)))
	}
}
