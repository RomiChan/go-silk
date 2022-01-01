package silk

import "unsafe"

func LPC_inverse_pred_gain_16(invGain_Q30 *int32, A_16 [2][16]int32, order int32) int32 {
	var (
		k            int32
		n            int32
		headrm       int32
		rc_Q31       int32
		rc_mult1_Q30 int32
		rc_mult2_Q16 int32
		tmp_16       int32
		Aold_16      *int32
		Anew_16      *int32
	)
	Anew_16 = &A_16[int64(order)&1][0]
	*invGain_Q30 = 1 << 30
	for k = int32(int64(order) - 1); int64(k) > 0; k-- {
		if int64(*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*uintptr(k)))) > int64(SKP_FIX_CONST(0.99975, 16)) || int64(*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*uintptr(k)))) < int64(-SKP_FIX_CONST(0.99975, 16)) {
			return 1
		}
		rc_Q31 = int32(-(int64(*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*uintptr(k)))) << (31 - 16)))
		rc_mult1_Q30 = int32((SKP_int32_MAX >> 1) - int64(SKP_SMMUL(rc_Q31, rc_Q31)))
		rc_mult2_Q16 = SKP_INVERSE32_varQ(rc_mult1_Q30, 46)
		*invGain_Q30 = int32(int64(SKP_SMMUL(*invGain_Q30, rc_mult1_Q30)) << 2)
		Aold_16 = Anew_16
		Anew_16 = &A_16[int64(k)&1][0]
		headrm = int32(int64(SKP_Silk_CLZ32(rc_mult2_Q16)) - 1)
		rc_mult2_Q16 = int32(int64(rc_mult2_Q16) << int64(headrm))
		for n = 0; int64(n) < int64(k); n++ {
			tmp_16 = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(Aold_16), unsafe.Sizeof(int32(0))*uintptr(n)))) - (int64(SKP_SMMUL(*(*int32)(unsafe.Add(unsafe.Pointer(Aold_16), unsafe.Sizeof(int32(0))*uintptr(int64(k)-int64(n)-1))), rc_Q31)) << 1))
			*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*uintptr(n))) = int32(int64(SKP_SMMUL(tmp_16, rc_mult2_Q16)) << (16 - int64(headrm)))
		}
	}
	if int64(*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*0))) > int64(SKP_FIX_CONST(0.99975, 16)) || int64(*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*0))) < int64(-SKP_FIX_CONST(0.99975, 16)) {
		return 1
	}
	rc_Q31 = int32(-(int64(*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*0))) << (31 - 16)))
	rc_mult1_Q30 = int32((SKP_int32_MAX >> 1) - int64(SKP_SMMUL(rc_Q31, rc_Q31)))
	*invGain_Q30 = int32(int64(SKP_SMMUL(*invGain_Q30, rc_mult1_Q30)) << 2)
	return 0
}
func SKP_Silk_LPC_inverse_pred_gain(invGain_Q30 *int32, A_Q12 *int16, order int32) int32 {
	var (
		k       int32
		Atmp_16 [2][16]int32
		Anew_16 *int32
	)
	Anew_16 = &Atmp_16[int64(order)&1][0]
	for k = 0; int64(k) < int64(order); k++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*uintptr(k))) = int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(int16(0))*uintptr(k))))) << (16 - 12))
	}
	return LPC_inverse_pred_gain_16(invGain_Q30, Atmp_16, order)
}
func SKP_Silk_LPC_inverse_pred_gain_Q24(invGain_Q30 *int32, A_Q24 *int32, order int32) int32 {
	var (
		k       int32
		Atmp_16 [2][16]int32
		Anew_16 *int32
	)
	Anew_16 = &Atmp_16[int64(order)&1][0]
	for k = 0; int64(k) < int64(order); k++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(Anew_16), unsafe.Sizeof(int32(0))*uintptr(k))) = SKP_RSHIFT_ROUND(*(*int32)(unsafe.Add(unsafe.Pointer(A_Q24), unsafe.Sizeof(int32(0))*uintptr(k))), 24-16)
	}
	return LPC_inverse_pred_gain_16(invGain_Q30, Atmp_16, order)
}
