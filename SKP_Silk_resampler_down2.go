package silk

import "unsafe"

func SKP_Silk_resampler_down2(S *int32, out *int16, in *int16, inLen int32) {
	var (
		k     int32
		len2  = inLen >> 1
		in32  int32
		out32 int32
		Y     int32
		X     int32
	)
	for k = 0; k < len2; k++ {
		in32 = (int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2))))) << 10
		Y = in32 - (*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)))
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_down2_1))
		out32 = (*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + X
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = in32 + X
		in32 = (int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2+1))))) << 10
		Y = in32 - (*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)))
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_down2_0))
		out32 = out32 + (*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)))
		out32 = out32 + X
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = in32 + X
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 11))
	}
}
