package silk

import "unsafe"

func SKP_Silk_resampler_private_down4(S *int32, out *int16, in *int16, inLen int32) {
	var (
		k     int32
		len4  int32 = int32(int64(inLen) >> 2)
		in32  int32
		out32 int32
		Y     int32
		X     int32
	)
	for k = 0; int64(k) < int64(len4); k++ {
		in32 = int32((int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(int64(k)*4))))) + int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(int64(k)*4+1)))))) << 9)
		Y = int32(int64(in32) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))))
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_down2_1))
		out32 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + int64(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int64(in32) + int64(X))
		in32 = int32((int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(int64(k)*4+2))))) + int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(int64(k)*4+3)))))) << 9)
		Y = int32(int64(in32) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))))
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_down2_0))
		out32 = int32(int64(out32) + int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))))
		out32 = int32(int64(out32) + int64(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int64(in32) + int64(X))
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = SKP_SAT16(int16(SKP_RSHIFT_ROUND(out32, 11)))
	}
}