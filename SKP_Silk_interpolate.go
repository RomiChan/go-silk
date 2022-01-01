package silk

func SKP_Silk_interpolate(xi [16]int32, x0 [16]int32, x1 [16]int32, ifact_Q2 int32, d int32) {
	var i int32
	for i = 0; int64(i) < int64(d); i++ {
		xi[i] = int32(int64(x0[i]) + (((int64(x1[i]) - int64(x0[i])) * int64(ifact_Q2)) >> 2))
	}
}
