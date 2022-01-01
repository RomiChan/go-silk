package silk

import (
	"math"
	"unsafe"
)

func SKP_Silk_MA_Prediction(in *int16, B *int16, S *int32, out *int16, len_ int32, order int32) {
	var (
		k     int32
		d     int32
		in16  int32
		out32 int32
	)
	for k = 0; int64(k) < int64(len_); k++ {
		in16 = int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))
		out32 = int32((int64(in16) << 12) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))))
		out32 = SKP_RSHIFT_ROUND(out32, 12)
		for d = 0; int64(d) < int64(order)-1; d++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(d))) = int32(int64(uint32(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(int64(d)+1))))) + int64(uint32(SKP_SMULBB(in16, int32(*(*int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(int16(0))*uintptr(d))))))))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(int64(order)-1))) = SKP_SMULBB(in16, int32(*(*int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(int16(0))*uintptr(int64(order)-1)))))
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = SKP_SAT16(int16(out32))
	}
}
func SKP_Silk_LPC_analysis_filter(in *int16, B *int16, S *int16, out *int16, len_ int32, Order int32) {
	var (
		k          int32
		j          int32
		idx        int32
		Order_half int32 = int32(int64(Order) >> 1)
		out32_Q12  int32
		out32      int32
		SA         int16
		SB         int16
	)
	for k = 0; int64(k) < int64(len_); k++ {
		SA = *(*int16)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int16(0))*0))
		out32_Q12 = 0
		for j = 0; int64(j) < (int64(Order_half) - 1); j++ {
			idx = int32(int64(SKP_SMULBB(2, j)) + 1)
			SB = *(*int16)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int16(0))*uintptr(idx)))
			*(*int16)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int16(0))*uintptr(idx))) = SA
			out32_Q12 = SKP_SMLABB(out32_Q12, int32(SA), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(int16(0))*uintptr(int64(idx)-1)))))
			out32_Q12 = SKP_SMLABB(out32_Q12, int32(SB), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(int16(0))*uintptr(idx)))))
			SA = *(*int16)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int16(0))*uintptr(int64(idx)+1)))
			*(*int16)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int16(0))*uintptr(int64(idx)+1))) = SB
		}
		SB = *(*int16)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int16(0))*uintptr(int64(Order)-1)))
		*(*int16)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int16(0))*uintptr(int64(Order)-1))) = SA
		out32_Q12 = SKP_SMLABB(out32_Q12, int32(SA), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(int16(0))*uintptr(int64(Order)-2)))))
		out32_Q12 = SKP_SMLABB(out32_Q12, int32(SB), int32(*(*int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(int16(0))*uintptr(int64(Order)-1)))))
		if (((int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))) << 12) - int64(out32_Q12)) & 0x80000000) == 0 {
			if ((int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))) << 12) & (int64(out32_Q12) ^ 0x80000000) & 0x80000000) != 0 {
				out32_Q12 = math.MinInt32
			} else {
				out32_Q12 = int32((int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))) << 12) - int64(out32_Q12))
			}
		} else if (((int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))) << 12) ^ 0x80000000) & int64(out32_Q12) & 0x80000000) != 0 {
			out32_Q12 = SKP_int32_MAX
		} else {
			out32_Q12 = int32((int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))) << 12) - int64(out32_Q12))
		}
		out32 = SKP_RSHIFT_ROUND(out32_Q12, 12)
		*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = SKP_SAT16(int16(out32))
		*(*int16)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int16(0))*0)) = *(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))
	}
}
