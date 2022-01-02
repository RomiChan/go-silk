package silk

func SKP_Silk_resampler_private_down4(S []int32, out []int16, in []int16, inLen int32) {
	var (
		k     int32
		len4  int32 = (inLen >> 2)
		in32  int32
		out32 int32
		Y     int32
		X     int32
	)
	SKP_assert(SKP_Silk_resampler_down2_0 > 0)
	SKP_assert(SKP_Silk_resampler_down2_1 < 0)
	for k = 0; k < len4; k++ {
		in32 = ((int32(in[k*4])) + (int32(in[k*4+1]))) << 9
		Y = in32 - (S[0])
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_down2_1))
		out32 = (S[0]) + X
		S[0] = in32 + X
		in32 = ((int32(in[k*4+2])) + (int32(in[k*4+3]))) << 9
		Y = in32 - (S[1])
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_down2_0))
		out32 = out32 + (S[1])
		out32 = out32 + X
		S[1] = in32 + X
		out[k] = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 11))
	}
}
