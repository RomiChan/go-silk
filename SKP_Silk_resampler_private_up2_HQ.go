package silk

import "unsafe"

func SKP_Silk_resampler_private_up2_HQ(S *int32, out *int16, in *int16, len_ int32) {
	var (
		k       int32
		in32    int32
		out32_1 int32
		out32_2 int32
		Y       int32
		X       int32
	)
	for k = 0; int64(k) < int64(len_); k++ {
		in32 = int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))) << 10)
		Y = int32(int64(in32) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))))
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_up2_hq_0[0]))
		out32_1 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + int64(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int64(in32) + int64(X))
		Y = int32(int64(out32_1) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))))
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_up2_hq_0[1]))
		out32_2 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + int64(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int64(out32_1) + int64(X))
		out32_2 = SKP_SMLAWB(out32_2, *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5)), int32(SKP_Silk_resampler_up2_hq_notch[2]))
		out32_2 = SKP_SMLAWB(out32_2, *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4)), int32(SKP_Silk_resampler_up2_hq_notch[1]))
		out32_1 = SKP_SMLAWB(out32_2, *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4)), int32(SKP_Silk_resampler_up2_hq_notch[0]))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5)) = int32(int64(out32_2) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5))))
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int64(k)*2))) = SKP_SAT16(int16(int64(SKP_SMLAWB(256, out32_1, int32(SKP_Silk_resampler_up2_hq_notch[3]))) >> 9))
		Y = int32(int64(in32) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2))))
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_up2_hq_1[0]))
		out32_1 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2))) + int64(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2)) = int32(int64(in32) + int64(X))
		Y = int32(int64(out32_1) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3))))
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_up2_hq_1[1]))
		out32_2 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3))) + int64(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3)) = int32(int64(out32_1) + int64(X))
		out32_2 = SKP_SMLAWB(out32_2, *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4)), int32(SKP_Silk_resampler_up2_hq_notch[2]))
		out32_2 = SKP_SMLAWB(out32_2, *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5)), int32(SKP_Silk_resampler_up2_hq_notch[1]))
		out32_1 = SKP_SMLAWB(out32_2, *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5)), int32(SKP_Silk_resampler_up2_hq_notch[0]))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4)) = int32(int64(out32_2) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4))))
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int64(k)*2+1))) = SKP_SAT16(int16(int64(SKP_SMLAWB(256, out32_1, int32(SKP_Silk_resampler_up2_hq_notch[3]))) >> 9))
	}
}
func SKP_Silk_resampler_private_up2_HQ_wrapper(SS unsafe.Pointer, out *int16, in *int16, len_ int32) {
	var S *SKP_Silk_resampler_state_struct = (*SKP_Silk_resampler_state_struct)(SS)
	SKP_Silk_resampler_private_up2_HQ(&S.SIIR[0], out, in, len_)
}