package silk

func SKP_Silk_resampler_private_AR2(S []int32, out_Q8 []int32, in []int16, A_Q14 []int16, len_ int32) {
	var (
		k     int32
		out32 int32
	)
	for k = 0; k < len_; k++ {
		out32 = (S[0]) + ((int32(in[k])) << 8)
		out_Q8[k] = out32
		out32 = out32 << 2
		S[0] = SKP_SMLAWB(S[1], out32, int32(A_Q14[0]))
		S[1] = SKP_SMULWB(out32, int32(A_Q14[1]))
	}
}
