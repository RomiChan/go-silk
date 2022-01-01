package silk

import "unsafe"

func SKP_Silk_schur64(rc_Q16 []int32, c []int32, order int32) int32 {
	var (
		k          int32
		n          int32
		C          [17][2]int32
		Ctmp1_Q30  int32
		Ctmp2_Q30  int32
		rc_tmp_Q31 int32
	)
	if c[0] <= 0 {
		memset(unsafe.Pointer(&rc_Q16[0]), 0, size_t(uintptr(order)*unsafe.Sizeof(int32(0))))
		return 0
	}
	for k = 0; k < order+1; k++ {
		C[k][0] = func() int32 {
			p := &C[k][1]
			C[k][1] = c[k]
			return *p
		}()
	}
	for k = 0; k < order; k++ {
		rc_tmp_Q31 = SKP_DIV32_varQ(-C[k+1][0], C[0][1], 31)
		rc_Q16[k] = SKP_RSHIFT_ROUND(rc_tmp_Q31, 15)
		for n = 0; n < order-k; n++ {
			Ctmp1_Q30 = C[n+k+1][0]
			Ctmp2_Q30 = C[n][1]
			C[n+k+1][0] = Ctmp1_Q30 + SKP_SMMUL(Ctmp2_Q30<<1, rc_tmp_Q31)
			C[n][1] = Ctmp2_Q30 + SKP_SMMUL(Ctmp1_Q30<<1, rc_tmp_Q31)
		}
	}
	return C[0][1]
}
