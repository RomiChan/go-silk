package silk

import "unsafe"

const MAX_LOOPS = 20

func SKP_Silk_NLSF_stabilize(NLSF_Q15 *int32, NDeltaMin_Q15 *int32, L int32) {
	var (
		center_freq_Q15 int32
		diff_Q15        int32
		min_center_Q15  int32
		max_center_Q15  int32
		min_diff_Q15    int32
		loops           int32
		i               int32
		I               int32 = 0
		k               int32
	)
	for loops = 0; loops < MAX_LOOPS; loops++ {
		min_diff_Q15 = *(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)) - *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*0))
		I = 0
		for i = 1; i <= L-1; i++ {
			diff_Q15 = *(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))) - (*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i-1))) + *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(i))))
			if diff_Q15 < min_diff_Q15 {
				min_diff_Q15 = diff_Q15
				I = i
			}
		}
		diff_Q15 = (1 << 15) - (*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(L-1))) + *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(L))))
		if diff_Q15 < min_diff_Q15 {
			min_diff_Q15 = diff_Q15
			I = L
		}
		if min_diff_Q15 >= 0 {
			return
		}
		if I == 0 {
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)) = *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*0))
		} else if I == L {
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(L-1))) = (1 << 15) - *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(L)))
		} else {
			min_center_Q15 = 0
			for k = 0; k < I; k++ {
				min_center_Q15 += *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(k)))
			}
			min_center_Q15 += (*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))) >> 1
			max_center_Q15 = 1 << 15
			for k = L; k > I; k-- {
				max_center_Q15 -= *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(k)))
			}
			max_center_Q15 -= *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I))) - ((*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))) >> 1)
			center_freq_Q15 = SKP_LIMIT_32(SKP_RSHIFT_ROUND(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I-1)))+*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))), 1), min_center_Q15, max_center_Q15)
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I-1))) = center_freq_Q15 - ((*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))) >> 1)
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))) = *(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I-1))) + *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))
		}
	}
	if loops == MAX_LOOPS {
		SKP_Silk_insertion_sort_increasing_all_values((*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)), L)
		*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)) = SKP_max_int(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)), *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*0)))
		for i = 1; i < L; i++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_max_int(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))), *(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i-1)))+*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(i))))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(L-1))) = SKP_min_int(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(L-1))), (1<<15)-*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(L))))
		for i = L - 2; i >= 0; i-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_min_int(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))), *(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i+1)))-*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(i+1))))
		}
	}
}
