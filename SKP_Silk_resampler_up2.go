package silk

// SKP_Silk_resampler_up2 Upsample by a factor 2, low quality
// see SKP_Silk_resampler_up2.c
func SKP_Silk_resampler_up2(S []int32, out []int16, in []int16, len_ int32) {
	/* Internal variables and state are in Q10 format */
	for k := int32(0); k < len_; k++ {
		/* Convert to Q10 */
		in32 := (int32(in[k])) << 10

		/* All-pass section for even output sample */
		Y := in32 - (S[0])
		X := SKP_SMULWB(Y, int32(SKP_Silk_resampler_up2_lq_0))
		out32 := (S[0]) + X
		S[0] = in32 + X

		/* Convert back to int16 and store to output */
		out[k*2] = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 10))

		/* All-pass section for odd output sample */
		Y = in32 - (S[1])
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_up2_lq_1))
		out32 = (S[1]) + X
		S[1] = in32 + X

		/* Convert back to int16 and store to output */
		out[k*2+1] = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 10))
	}
}
