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
	for k = 1; int64(k) < int64(dd); k++ {
		ftmp = *(*int32)(unsafe.Add(unsafe.Pointer(cLSF), unsafe.Sizeof(int32(0))*uintptr(int64(k)*2)))
		*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(int64(k)+1))) = int32((int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(int64(k)-1)))) << 1) - int64(SKP_RSHIFT_ROUND64(int64(ftmp)*int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k)))), 20)))
		for n = k; int64(n) > 1; n-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n))) += int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(int64(n)-2)))) - int64(SKP_RSHIFT_ROUND64(int64(ftmp)*int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(int64(n)-1)))), 20)))
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
	for k = 0; int64(k) < int64(d); k++ {
		f_int = (*(*int32)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(int32(0))*uintptr(k)))) >> (15 - 7)
		f_frac = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(int32(0))*uintptr(k)))) - (int64(f_int) << (15 - 7)))
		cos_val = SKP_Silk_LSFCosTab_FIX_Q12[f_int]
		delta = int32(int64(SKP_Silk_LSFCosTab_FIX_Q12[int64(f_int)+1]) - int64(cos_val))
		cos_LSF_Q20[k] = int32((int64(cos_val) << 8) + int64(delta)*int64(f_frac))
	}
	dd = d >> 1
	SKP_Silk_NLSF2A_find_poly(&P[0], &cos_LSF_Q20[0], dd)
	SKP_Silk_NLSF2A_find_poly(&Q[0], &cos_LSF_Q20[1], dd)
	for k = 0; int64(k) < int64(dd); k++ {
		Ptmp = int32(int64(P[int64(k)+1]) + int64(P[k]))
		Qtmp = int32(int64(Q[int64(k)+1]) - int64(Q[k]))
		a_int32[k] = -SKP_RSHIFT_ROUND(int32(int64(Ptmp)+int64(Qtmp)), 9)
		a_int32[int64(d)-int64(k)-1] = SKP_RSHIFT_ROUND(int32(int64(Qtmp)-int64(Ptmp)), 9)
	}
	for i = 0; int64(i) < 10; i++ {
		maxabs = 0
		for k = 0; int64(k) < int64(d); k++ {
			absval = int32(SKP_abs(int64(a_int32[k])))
			if int64(absval) > int64(maxabs) {
				maxabs = absval
				idx = k
			}
		}
		if int64(maxabs) > SKP_int16_MAX {
			if int64(maxabs) < 98369 {
				maxabs = maxabs
			} else {
				maxabs = 98369
			}
			sc_Q16 = int32(65470 - int64(int32(((int64(maxabs)-SKP_int16_MAX)*(65470>>2))/((int64(maxabs)*(int64(idx)+1))>>2))))
			SKP_Silk_bwexpander_32(&a_int32[0], d, sc_Q16)
		} else {
			break
		}
	}
	if int64(i) == 10 {
		for k = 0; int64(k) < int64(d); k++ {
			a_int32[k] = int32(SKP_SAT16(int16(a_int32[k])))
		}
	}
	for k = 0; int64(k) < int64(d); k++ {
		*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(a_int32[k])
	}
}
