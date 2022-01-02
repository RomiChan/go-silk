package silk

func LPC_inverse_pred_gain_16(invGain_Q30 *int32, A_16 [2][16]int32, order int32) int32 {
	var (
		k            int32
		n            int32
		headrm       int32
		rc_Q31       int32
		rc_mult1_Q30 int32
		rc_mult2_Q16 int32
		tmp_16       int32
		Aold_16      []int32
		Anew_16      []int32
	)
	Anew_16 = ([]int32)(A_16[order&1][:])
	*invGain_Q30 = 1 << 30
	for k = order - 1; k > 0; k-- {
		if Anew_16[k] > SKP_FIX_CONST(0.99975, 16) || Anew_16[k] < -SKP_FIX_CONST(0.99975, 16) {
			return 1
		}
		rc_Q31 = -((Anew_16[k]) << (31 - 16))
		rc_mult1_Q30 = (SKP_int32_MAX >> 1) - SKP_SMMUL(rc_Q31, rc_Q31)
		SKP_assert(rc_mult1_Q30 > (1 << 15))
		SKP_assert(rc_mult1_Q30 < (1 << 30))
		rc_mult2_Q16 = SKP_INVERSE32_varQ(rc_mult1_Q30, 46)
		*invGain_Q30 = SKP_SMMUL(*invGain_Q30, rc_mult1_Q30) << 2
		SKP_assert(*invGain_Q30 >= 0)
		SKP_assert(*invGain_Q30 <= (1 << 30))
		Aold_16 = Anew_16
		Anew_16 = ([]int32)(A_16[k&1][:])
		headrm = SKP_Silk_CLZ32(rc_mult2_Q16) - 1
		rc_mult2_Q16 = rc_mult2_Q16 << headrm
		for n = 0; n < k; n++ {
			tmp_16 = Aold_16[n] - (SKP_SMMUL(Aold_16[k-n-1], rc_Q31) << 1)
			Anew_16[n] = SKP_SMMUL(tmp_16, rc_mult2_Q16) << (16 - headrm)
		}
	}
	if Anew_16[0] > SKP_FIX_CONST(0.99975, 16) || Anew_16[0] < -SKP_FIX_CONST(0.99975, 16) {
		return 1
	}
	rc_Q31 = -((Anew_16[0]) << (31 - 16))
	rc_mult1_Q30 = (SKP_int32_MAX >> 1) - SKP_SMMUL(rc_Q31, rc_Q31)
	*invGain_Q30 = SKP_SMMUL(*invGain_Q30, rc_mult1_Q30) << 2
	SKP_assert(*invGain_Q30 >= 0)
	SKP_assert(*invGain_Q30 <= 1<<30)
	return 0
}
func SKP_Silk_LPC_inverse_pred_gain(invGain_Q30 *int32, A_Q12 []int16, order int32) int32 {
	var (
		k       int32
		Atmp_16 [2][16]int32
		Anew_16 []int32
	)
	_ = Anew_16
	Anew_16 = ([]int32)(Atmp_16[order&1][:])
	for k = 0; k < order; k++ {
		Anew_16[k] = (int32(A_Q12[k])) << (16 - 12)
	}
	return LPC_inverse_pred_gain_16(invGain_Q30, Atmp_16, order)
}
func SKP_Silk_LPC_inverse_pred_gain_Q24(invGain_Q30 *int32, A_Q24 []int32, order int32) int32 {
	var (
		k       int32
		Atmp_16 [2][16]int32
		Anew_16 []int32
	)
	_ = Anew_16
	Anew_16 = ([]int32)(Atmp_16[order&1][:])
	for k = 0; k < order; k++ {
		Anew_16[k] = SKP_RSHIFT_ROUND(A_Q24[k], 24-16)
	}
	return LPC_inverse_pred_gain_16(invGain_Q30, Atmp_16, order)
}
