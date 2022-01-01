package silk

func SKP_Silk_resampler_private_ARMA4(S []int32, out []int16, in []int16, Coef []int16, len_ int32) {
	var (
		k       int32
		in_Q8   int32
		out1_Q8 int32
		out2_Q8 int32
		X       int32
	)
	for k = 0; k < len_; k++ {
		in_Q8 = (int32(in[k])) << 8
		out1_Q8 = in_Q8 + ((S[0]) << 2)
		out2_Q8 = out1_Q8 + ((S[2]) << 2)
		X = SKP_SMLAWB(S[1], in_Q8, int32(Coef[0]))
		S[0] = SKP_SMLAWB(X, out1_Q8, int32(Coef[2]))
		X = SKP_SMLAWB(S[3], out1_Q8, int32(Coef[1]))
		S[2] = SKP_SMLAWB(X, out2_Q8, int32(Coef[4]))
		S[1] = SKP_SMLAWB(in_Q8>>2, out1_Q8, int32(Coef[3]))
		S[3] = SKP_SMLAWB(out1_Q8>>2, out2_Q8, int32(Coef[5]))
		out[k] = SKP_SAT16(SKP_SMLAWB(128, out2_Q8, int32(Coef[6])) >> 8)
	}
}
