package silk

import "unsafe"

func SKP_Silk_NLSF2A_find_poly(out *int32, cLSF *int32, dd int32) {
	var (
		k    int32
		n    int32
		ftmp int32
	)
	*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*0)) = 1 << 20
	*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*1)) = -*(*int32)(unsafe.Add(unsafe.Pointer(cLSF), unsafe.Sizeof(int32(0))*0))
	for k = 1; k < dd; k++ {
		ftmp = *(*int32)(unsafe.Add(unsafe.Pointer(cLSF), unsafe.Sizeof(int32(0))*uintptr(k*2)))
		*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k+1))) = ((*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k-1)))) << 1) - SKP_RSHIFT_ROUND64(int64(ftmp)*int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k)))), 20)
		for n = k; n > 1; n-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n))) += *(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n-2))) - SKP_RSHIFT_ROUND64(int64(ftmp)*int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n-1)))), 20)
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*1)) -= ftmp
	}
}
func SKP_Silk_NLSF2A(a *int16, NLSF *int32, d int32) {
	var (
		k           int32
		i           int32
		dd          int32
		cos_LSF_Q20 [16]int32
		P           [9]int32
		Q           [9]int32
		Ptmp        int32
		Qtmp        int32
		f_int       int32
		f_frac      int32
		cos_val     int32
		delta       int32
		a_int32     [16]int32
		maxabs      int32
		absval      int32
		idx         int32 = 0
		sc_Q16      int32
	)
	for k = 0; k < d; k++ {
		f_int = (*(*int32)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(int32(0))*uintptr(k)))) >> (15 - 7)
		f_frac = *(*int32)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(int32(0))*uintptr(k))) - (f_int << (15 - 7))
		cos_val = SKP_Silk_LSFCosTab_FIX_Q12[f_int]
		delta = SKP_Silk_LSFCosTab_FIX_Q12[f_int+1] - cos_val
		cos_LSF_Q20[k] = (cos_val << 8) + delta*f_frac
	}
	dd = d >> 1
	SKP_Silk_NLSF2A_find_poly(&P[0], &cos_LSF_Q20[0], dd)
	SKP_Silk_NLSF2A_find_poly(&Q[0], &cos_LSF_Q20[1], dd)
	for k = 0; k < dd; k++ {
		Ptmp = P[k+1] + P[k]
		Qtmp = Q[k+1] - Q[k]
		a_int32[k] = -SKP_RSHIFT_ROUND(Ptmp+Qtmp, 9)
		a_int32[d-k-1] = SKP_RSHIFT_ROUND(Qtmp-Ptmp, 9)
	}
	for i = 0; i < 10; i++ {
		maxabs = 0
		for k = 0; k < d; k++ {
			absval = int32(SKP_abs(int64(a_int32[k])))
			if absval > maxabs {
				maxabs = absval
				idx = k
			}
		}
		if maxabs > SKP_int16_MAX {
			if maxabs < 98369 {
				maxabs = maxabs
			} else {
				maxabs = 98369
			}
			sc_Q16 = 65470 - ((maxabs-SKP_int16_MAX)*(65470>>2))/((maxabs*(idx+1))>>2)
			SKP_Silk_bwexpander_32(&a_int32[0], d, sc_Q16)
		} else {
			break
		}
	}
	if i == 10 {
		for k = 0; k < d; k++ {
			a_int32[k] = int32(SKP_SAT16(a_int32[k]))
		}
	}
	for k = 0; k < d; k++ {
		*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(a_int32[k])
	}
}
