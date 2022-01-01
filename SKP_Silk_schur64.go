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
	if int64(c[0]) <= 0 {
		memset(unsafe.Pointer(&rc_Q16[0]), 0, size_t(uintptr(order)*unsafe.Sizeof(int32(0))))
		return 0
	}
	for k = 0; int64(k) < int64(order)+1; k++ {
		C[k][0] = func() int32 {
			p := &C[k][1]
			C[k][1] = c[k]
			return *p
		}()
	}
	for k = 0; int64(k) < int64(order); k++ {
		rc_tmp_Q31 = SKP_DIV32_varQ(-C[int64(k)+1][0], C[0][1], 31)
		rc_Q16[k] = SKP_RSHIFT_ROUND(rc_tmp_Q31, 15)
		for n = 0; int64(n) < int64(order)-int64(k); n++ {
			Ctmp1_Q30 = C[int64(n)+int64(k)+1][0]
			Ctmp2_Q30 = C[n][1]
			C[int64(n)+int64(k)+1][0] = int32(int64(Ctmp1_Q30) + int64(SKP_SMMUL(int32(int64(Ctmp2_Q30)<<1), rc_tmp_Q31)))
			C[n][1] = int32(int64(Ctmp2_Q30) + int64(SKP_SMMUL(int32(int64(Ctmp1_Q30)<<1), rc_tmp_Q31)))
		}
	}
	return C[0][1]
}
