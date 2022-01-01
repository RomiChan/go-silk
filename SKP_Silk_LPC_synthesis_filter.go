package silk

import (
	"math"
	"unsafe"
)

func SKP_Silk_LPC_synthesis_filter(in *int16, A_Q12 *int16, Gain_Q26 int32, S *int32, out *int16, len_ int32, Order int32) {
	var (
		k          int32
		j          int32
		idx        int32
		Order_half int32 = int32(int64(Order) >> 1)
		SA         int32
		SB         int32
		out32_Q10  int32
		out32      int32
	)
	for k = 0; int64(k) < int64(len_); k++ {
		SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(int64(Order)-1)))
		out32_Q10 = 0
		for j = 0; int64(j) < (int64(Order_half) - 1); j++ {
			idx = int32(int64(SKP_SMULBB(2, j)) + 1)
			SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(int64(Order)-1-int64(idx))))
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(int64(Order)-1-int64(idx)))) = SA
			out32_Q10 = SKP_SMLAWB(out32_Q10, SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*uintptr(int64(j)<<1)))))
			out32_Q10 = SKP_SMLAWB(out32_Q10, SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*uintptr((int64(j)<<1)+1)))))
			SA = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(int64(Order)-2-int64(idx))))
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(int64(Order)-2-int64(idx)))) = SB
		}
		SB = *(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = SA
		out32_Q10 = SKP_SMLAWB(out32_Q10, SA, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*uintptr(int64(Order)-2)))))
		out32_Q10 = SKP_SMLAWB(out32_Q10, SB, int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*uintptr(int64(Order)-1)))))
		if ((int64(out32_Q10) + int64(SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))))) & 0x80000000) == 0 {
			if ((int64(out32_Q10) & int64(SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))))) & 0x80000000) != 0 {
				out32_Q10 = math.MinInt32
			} else {
				out32_Q10 = int32(int64(out32_Q10) + int64(SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))))))
			}
		} else if ((int64(out32_Q10) | int64(SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))))) & 0x80000000) == 0 {
			out32_Q10 = SKP_int32_MAX
		} else {
			out32_Q10 = int32(int64(out32_Q10) + int64(SKP_SMULWB(Gain_Q26, int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))))))
		}
		out32 = SKP_RSHIFT_ROUND(out32_Q10, 10)
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = SKP_SAT16(int16(out32))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(int64(Order)-1))) = int32((func() int64 {
			if (int64(math.MinInt32) >> 4) > (SKP_int32_MAX >> 4) {
				if int64(out32_Q10) > (int64(math.MinInt32) >> 4) {
					return int64(math.MinInt32) >> 4
				}
				if int64(out32_Q10) < (SKP_int32_MAX >> 4) {
					return SKP_int32_MAX >> 4
				}
				return int64(out32_Q10)
			}
			if int64(out32_Q10) > (SKP_int32_MAX >> 4) {
				return SKP_int32_MAX >> 4
			}
			if int64(out32_Q10) < (int64(math.MinInt32) >> 4) {
				return int64(math.MinInt32) >> 4
			}
			return int64(out32_Q10)
		}()) << 4)
	}
}
