package silk

func SKP_Silk_resampler_down2(S []int32, out []int16, in []int16, inLen int32) {
	var (
		k     int32
		len2  int32 = (inLen >> 1)
		in32  int32
		out32 int32
		Y     int32
		X     int32
	)
	for k = 0; k < len2; k++ {
		in32 = (int32(in[k*2])) << 10
		Y = in32 - (S[0])
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_down2_1))
		out32 = (S[0]) + X
		S[0] = in32 + X
		in32 = (int32(in[k*2+1])) << 10
		Y = in32 - (S[1])
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_down2_0))
		out32 = out32 + (S[1])
		out32 = out32 + X
		S[1] = in32 + X
		out[k] = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 11))
	}
}
