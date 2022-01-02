package silk

// Reviewed by wdvxdr1123 2022-01-02

// SKP_Silk_resampler_down2
// Downsample by a factor 2, mediocre quality
func SKP_Silk_resampler_down2(
	S []int32, //   I/O: State vector [ 2 ]
	out []int16, // O:   Output signal [ len ]
	in []int16, //  I:   Input signal [ floor(len/2) ]
	inLen int32, // I:   Number of input samples
) {
	len2 := inLen >> 1
	SKP_assert(SKP_Silk_resampler_down2_0 > 0)
	SKP_assert(SKP_Silk_resampler_down2_1 < 0)

	/* Internal variables and state are in Q10 format */
	for k := int32(0); k < len2; k++ {
		/* Convert to Q10 */
		in32 := int32(in[k*2]) << 10

		/* All-pass section for even input sample */
		Y := in32 - S[0]
		X := SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_down2_1))
		out32 := S[0] + X
		S[0] = in32 + X

		/* Convert to Q10 */
		in32 = (int32(in[k*2+1])) << 10

		/* All-pass section for odd input sample, and add to output of previous section */
		Y = in32 - S[1]
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_down2_0))
		out32 = out32 + S[1]
		out32 = out32 + X
		S[1] = in32 + X

		/* Add, convert back to int16 and store to output */
		out[k] = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 11))
	}
}
