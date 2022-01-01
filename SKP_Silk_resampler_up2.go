package silk

func SKP_Silk_resampler_up2(S []int32, out []int16, in []int16, len_ int32) {
	var (
		k     int32
		in32  int32
		out32 int32
		Y     int32
		X     int32
	)
	for k = 0; k < len_; k++ {
		in32 = (int32(in[k])) << 10
		Y = in32 - (S[0])
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_up2_lq_0))
		out32 = (S[0]) + X
		S[0] = in32 + X
		out[k*2] = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 10))
		Y = in32 - (S[1])
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_up2_lq_1))
		out32 = (S[1]) + X
		S[1] = in32 + X
		out[k*2+1] = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 10))
	}
}
