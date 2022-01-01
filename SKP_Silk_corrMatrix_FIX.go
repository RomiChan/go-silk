package silk

import "unsafe"

func SKP_Silk_corrVector_FIX(
	x []int16, /* I    x vector [L + order - 1] used to form data matrix X */
	t []int16, /* I    target vector [L]*/
	L int32,
	order int32,
	Xt []int32, /* O    Pointer to X'*t correlation vector [order]          */
	rshifts int32,
) {
	var (
		lag        int32
		i          int32
		ptr1       *int16
		ptr2       *int16
		inner_prod int32
	)
	ptr1 = &x[order-1]
	ptr2 = &t[0]
	if rshifts > 0 {
		for lag = 0; lag < order; lag++ {
			inner_prod = 0
			for i = 0; i < L; i++ {
				inner_prod += SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(i)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int16(0))*uintptr(i))))) >> rshifts
			}
			Xt[lag] = inner_prod
			ptr1 = (*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*1)))
		}
	} else {
		for lag = 0; lag < order; lag++ {
			Xt[lag] = SKP_Silk_inner_prod_aligned(([]int16)(ptr1), ([]int16)(ptr2), L)
			ptr1 = (*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*1)))
		}
	}
}

// SKP_Silk_corrMatrix_FIX .
// see SKP_Silk_corrMatrix_FIX.c
func SKP_Silk_corrMatrix_FIX(
	x []int16, // [L+order-1]
	L int32,
	order int32,
	head_room int32,
	XX []int32, // [order * order]
	rshifts *int32,
) {
	var (
		i                 int32
		j                 int32
		lag               int32
		rshifts_local     int32
		head_room_rshifts int32
		energy            int32
		ptr1              *int16
		ptr2              *int16
	)
	SKP_Silk_sum_sqr_shift(&energy, &rshifts_local, x, L+order-1)
	if (head_room - SKP_Silk_CLZ32(energy)) > 0 {
		head_room_rshifts = head_room - SKP_Silk_CLZ32(energy)
	} else {
		head_room_rshifts = 0
	}
	energy = energy >> head_room_rshifts
	rshifts_local += head_room_rshifts
	for i = 0; i < order-1; i++ {
		energy -= SKP_SMULBB(int32(x[i]), int32(x[i])) >> rshifts_local
	}
	if rshifts_local < *rshifts {
		energy = energy >> (*rshifts - rshifts_local)
		rshifts_local = *rshifts
	}
	XX[order*0+0] = energy
	ptr1 = &x[order-1]
	for j = 1; j < order; j++ {
		energy = energy - (SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(L-j)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(L-j))))) >> rshifts_local)
		energy = energy + (SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*uintptr(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*uintptr(j)))))) >> rshifts_local)
		XX[j*order+j] = energy
	}
	ptr2 = &x[order-2]
	if rshifts_local > 0 {
		for lag = 1; lag < order; lag++ {
			energy = 0
			for i = 0; i < L; i++ {
				energy += SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(i)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int16(0))*uintptr(i))))) >> rshifts_local
			}
			XX[lag*order+0] = energy
			XX[order*0+lag] = energy
			for j = 1; j < (order - lag); j++ {
				energy = energy - (SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(L-j)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int16(0))*uintptr(L-j))))) >> rshifts_local)
				energy = energy + (SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*uintptr(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(int16(0))*uintptr(j)))))) >> rshifts_local)
				XX[(lag+j)*order+j] = energy
				XX[j*order+(lag+j)] = energy
			}
			ptr2 = (*int16)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(int16(0))*1)))
		}
	} else {
		for lag = 1; lag < order; lag++ {
			energy = SKP_Silk_inner_prod_aligned(([]int16)(ptr1), ([]int16)(ptr2), L)
			XX[lag*order+0] = energy
			XX[order*0+lag] = energy
			for j = 1; j < (order - lag); j++ {
				energy = energy - SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(L-j)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int16(0))*uintptr(L-j)))))
				energy = SKP_SMLABB(energy, int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*uintptr(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(int16(0))*uintptr(j))))))
				XX[(lag+j)*order+j] = energy
				XX[j*order+(lag+j)] = energy
			}
			ptr2 = (*int16)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(int16(0))*1)))
		}
	}
	*rshifts = rshifts_local
}
