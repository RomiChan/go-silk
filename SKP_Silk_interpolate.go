package silk

func SKP_Silk_interpolate(xi [16]int32, x0 [16]int32, x1 [16]int32, ifact_Q2 int32, d int32) {
	var i int32
	for i = 0; i < d; i++ {
		xi[i] = x0[i] + (((x1[i] - x0[i]) * ifact_Q2) >> 2)
	}
}
