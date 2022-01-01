package silk

import "unsafe"

func SKP_Silk_corrVector_FIX(x *int16, t *int16, L int32, order int32, Xt *int32, rshifts int32) {
	var (
		lag        int32
		i          int32
		ptr1       *int16
		ptr2       *int16
		inner_prod int32
	)
	ptr1 = (*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(int64(order)-1)))
	ptr2 = t
	if int64(rshifts) > 0 {
		for lag = 0; int64(lag) < int64(order); lag++ {
			inner_prod = 0
			for i = 0; int64(i) < int64(L); i++ {
				inner_prod += SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(i)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int16(0))*uintptr(i))))) >> int64(rshifts)
			}
			*(*int32)(unsafe.Add(unsafe.Pointer(Xt), unsafe.Sizeof(int32(0))*uintptr(lag))) = inner_prod
			ptr1 = (*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*1)))
		}
	} else {
		for lag = 0; int64(lag) < int64(order); lag++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(Xt), unsafe.Sizeof(int32(0))*uintptr(lag))) = SKP_Silk_inner_prod_aligned(([]int16)(ptr1), ([]int16)(ptr2), L)
			ptr1 = (*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*1)))
		}
	}
}
func SKP_Silk_corrMatrix_FIX(x *int16, L int32, order int32, head_room int32, XX *int32, rshifts *int32) {
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
	SKP_Silk_sum_sqr_shift(&energy, &rshifts_local, x, int32(int64(L)+int64(order)-1))
	if (int64(head_room) - int64(SKP_Silk_CLZ32(energy))) > 0 {
		head_room_rshifts = int32(int64(head_room) - int64(SKP_Silk_CLZ32(energy)))
	} else {
		head_room_rshifts = 0
	}
	energy = energy >> int64(head_room_rshifts)
	rshifts_local += head_room_rshifts
	for i = 0; int64(i) < int64(order)-1; i++ {
		energy -= SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(i)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(i))))) >> int64(rshifts_local)
	}
	if int64(rshifts_local) < int64(*rshifts) {
		energy = energy >> (int64(*rshifts) - int64(rshifts_local))
		rshifts_local = *rshifts
	}
	*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr(int64(order)*0+0)))) = energy
	ptr1 = (*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(int64(order)-1)))
	for j = 1; int64(j) < int64(order); j++ {
		energy = int32(int64(energy) - (int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(int64(L)-int64(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(int64(L)-int64(j))))))) >> int64(rshifts_local)))
		energy = int32(int64(energy) + (int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*uintptr(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*uintptr(j))))))) >> int64(rshifts_local)))
		*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr(int64(j)*int64(order)+int64(j))))) = energy
	}
	ptr2 = (*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(int64(order)-2)))
	if int64(rshifts_local) > 0 {
		for lag = 1; int64(lag) < int64(order); lag++ {
			energy = 0
			for i = 0; int64(i) < int64(L); i++ {
				energy += SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(i)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int16(0))*uintptr(i))))) >> int64(rshifts_local)
			}
			*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr(int64(lag)*int64(order)+0)))) = energy
			*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr(int64(order)*0+int64(lag))))) = energy
			for j = 1; int64(j) < (int64(order) - int64(lag)); j++ {
				energy = int32(int64(energy) - (int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(int64(L)-int64(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int16(0))*uintptr(int64(L)-int64(j))))))) >> int64(rshifts_local)))
				energy = int32(int64(energy) + (int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*uintptr(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(int16(0))*uintptr(j))))))) >> int64(rshifts_local)))
				*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr((int64(lag)+int64(j))*int64(order)+int64(j))))) = energy
				*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr(int64(j)*int64(order)+(int64(lag)+int64(j)))))) = energy
			}
			ptr2 = (*int16)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(int16(0))*1)))
		}
	} else {
		for lag = 1; int64(lag) < int64(order); lag++ {
			energy = SKP_Silk_inner_prod_aligned(([]int16)(ptr1), ([]int16)(ptr2), L)
			*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr(int64(lag)*int64(order)+0)))) = energy
			*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr(int64(order)*0+int64(lag))))) = energy
			for j = 1; int64(j) < (int64(order) - int64(lag)); j++ {
				energy = int32(int64(energy) - int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int16(0))*uintptr(int64(L)-int64(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int16(0))*uintptr(int64(L)-int64(j))))))))
				energy = SKP_SMLABB(energy, int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(int16(0))*uintptr(j))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(int16(0))*uintptr(j))))))
				*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr((int64(lag)+int64(j))*int64(order)+int64(j))))) = energy
				*((*int32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(int32(0))*uintptr(int64(j)*int64(order)+(int64(lag)+int64(j)))))) = energy
			}
			ptr2 = (*int16)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(int16(0))*1)))
		}
	}
	*rshifts = rshifts_local
}
