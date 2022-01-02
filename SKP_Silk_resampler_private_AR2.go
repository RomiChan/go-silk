package silk

// Reviewed by wdvxdr1123 2022-01-03

// SKP_Silk_resampler_private_AR2
// Second order AR filter with single delay elements
func SKP_Silk_resampler_private_AR2(
	S []int32, // I/O: State vector [ 2 ]
	out_Q8 []int32, // O:	Output signal
	in []int16, // I:	Input signal
	A_Q14 []int16, // I:	AR coefficients, Q14
	len_ int32, // I:	Signal length
) {
	for k := int32(0); k < len_; k++ {
		out32 := S[0] + (int32(in[k]) << 8)
		out_Q8[k] = out32
		out32 = out32 << 2
		S[0] = SKP_SMLAWB(S[1], out32, int32(A_Q14[0]))
		S[1] = SKP_SMULWB(out32, int32(A_Q14[1]))
	}
}
