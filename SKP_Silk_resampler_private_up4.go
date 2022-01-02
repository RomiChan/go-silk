package silk

func SKP_Silk_resampler_private_up4(S []int32, out []int16, in []int16, len_ int32) {
	var (
		k     int32
		in32  int32
		out32 int32
		Y     int32
		X     int32
		out16 int16
	)
	SKP_assert(SKP_Silk_resampler_up2_lq_0 > 0)
	SKP_assert(SKP_Silk_resampler_up2_lq_1 < 0)
	for k = 0; k < len_; k++ {
		in32 = (int32(in[k])) << 10
		Y = in32 - (S[0])
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_up2_lq_0))
		out32 = (S[0]) + X
		S[0] = in32 + X
		out16 = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 10))
		out[k*4] = out16
		out[k*4+1] = out16
		Y = in32 - (S[1])
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_up2_lq_1))
		out32 = (S[1]) + X
		S[1] = in32 + X
		out16 = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 10))
		out[k*4+2] = out16
		out[k*4+3] = out16
	}
}
