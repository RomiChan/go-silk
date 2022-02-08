package silk

// reviewed by wdvxdr 2022-02-08

func SKP_Silk_interpolate(xi [16]int32, x0 [16]int32, x1 [16]int32, ifact_Q2 int32, d int32) {
	SKP_assert(ifact_Q2 >= 0)
	SKP_assert(ifact_Q2 <= (1 << 2))
	for i := int32(0); i < d; i++ {
		xi[i] = x0[i] + (((x1[i] - x0[i]) * ifact_Q2) >> 2)
	}
}
