package silk

import "unsafe"

func SKP_Silk_LPC_synthesis_order16(in *int16, A_Q12 *int16, Gain_Q26 int32, S *int32, out *int16, len_ int32) {
	var (
		k         int32
		SA        int32
		SB        int32
		out32_Q10 int32
		out32     int32
	)
	for k = 0; k < len_; k++ {
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*15))
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*14))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*14)) = SA
		out32_Q10 = SKP_SMULWB(SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*0))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*1))))))
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*13))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*13)) = SB
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*12))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*12)) = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*2))))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*3))))))
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*11))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*11)) = SB
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*10))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*10)) = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*4))))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*5))))))
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*9))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*9)) = SB
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*8))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*8)) = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*6))))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*7))))))
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*7))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*7)) = SB
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*6))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*6)) = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*8))))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*9))))))
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5)) = SB
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4)) = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*10))))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*11))))))
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3)) = SB
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2)) = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*12))))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*13))))))
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = SB
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*14))))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*15))))))
		if ((out32_Q10 + SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))))) & math.MinInt32) == 0 {
			if ((out32_Q10 & SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))))) & math.MinInt32) != 0 {
				out32_Q10 = math.MinInt32
			} else {
				out32_Q10 = out32_Q10 + SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))))
			}
		} else if ((out32_Q10 | SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))))) & math.MinInt32) == 0 {
			out32_Q10 = SKP_int32_MAX
		} else {
			out32_Q10 = out32_Q10 + SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))))
		}
		out32 = SKP_RSHIFT_ROUND(out32_Q10, 10)
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = SKP_SAT16(out32)
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*15)) = SKP_LSHIFT_SAT32(out32_Q10, 4)
	}
}
