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
	for loops = 0; int64(loops) < MAX_LOOPS; loops++ {
		min_diff_Q15 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0))) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*0))))
		I = 0
		for i = 1; int64(i) <= int64(L)-1; i++ {
			diff_Q15 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i)))) - (int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(i)-1)))) + int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(i))))))
			if int64(diff_Q15) < int64(min_diff_Q15) {
				min_diff_Q15 = diff_Q15
				I = i
			}
		}
		diff_Q15 = int32((1 << 15) - (int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(L)-1)))) + int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(L))))))
		if int64(diff_Q15) < int64(min_diff_Q15) {
			min_diff_Q15 = diff_Q15
			I = L
		}
		if int64(min_diff_Q15) >= 0 {
			return
		}
		if int64(I) == 0 {
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)) = *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*0))
		} else if int64(I) == int64(L) {
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(L)-1))) = int32((1 << 15) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(L)))))
		} else {
			min_center_Q15 = 0
			for k = 0; int64(k) < int64(I); k++ {
				min_center_Q15 += *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(k)))
			}
			min_center_Q15 += (*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))) >> 1
			max_center_Q15 = 1 << 15
			for k = L; int64(k) > int64(I); k-- {
				max_center_Q15 -= *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(k)))
			}
			max_center_Q15 -= int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))) - (int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))) >> 1))
			if int64(min_center_Q15) > int64(max_center_Q15) {
				if int64(SKP_RSHIFT_ROUND(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(I)-1))))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))))), 1)) > int64(min_center_Q15) {
					center_freq_Q15 = min_center_Q15
				} else if int64(SKP_RSHIFT_ROUND(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(I)-1))))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))))), 1)) < int64(max_center_Q15) {
					center_freq_Q15 = max_center_Q15
				} else {
					center_freq_Q15 = SKP_RSHIFT_ROUND(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(I)-1))))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))))), 1)
				}
			} else if int64(SKP_RSHIFT_ROUND(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(I)-1))))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))))), 1)) > int64(max_center_Q15) {
				center_freq_Q15 = max_center_Q15
			} else if int64(SKP_RSHIFT_ROUND(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(I)-1))))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))))), 1)) < int64(min_center_Q15) {
				center_freq_Q15 = min_center_Q15
			} else {
				center_freq_Q15 = SKP_RSHIFT_ROUND(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(I)-1))))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))))), 1)
			}
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(I)-1))) = int32(int64(center_freq_Q15) - (int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))) >> 1))
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(I))) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(I)-1)))) + int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(I)))))
		}
	}
	if int64(loops) == MAX_LOOPS {
		SKP_Silk_insertion_sort_increasing_all_values((*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)), L)
		*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)) = SKP_max_int(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*0)), *(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*0)))
		for i = 1; int64(i) < int64(L); i++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_max_int(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))), int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(i)-1))))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(i))))))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(L)-1))) = SKP_min_int(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(L)-1))), int32((1<<15)-int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(L))))))
		for i = int32(int64(L) - 2); int64(i) >= 0; i-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))) = SKP_min_int(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))), int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(i)+1))))-int64(*(*int32)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int32(0))*uintptr(int64(i)+1))))))
		}
	}
}
