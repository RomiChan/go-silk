package silk

// Reviewed by wdvxdr1123 2022-01-03

// SKP_Silk_resampler_private_down4 Downsample by a factor 4.
// Note: very low quality, only use with input sampling rates above 96 kHz.
func SKP_Silk_resampler_private_down4(
	S []int32, /* I/O: State vector [ 2 ]                      */
	out []int16, /* O:   Output signal [ floor(len/2) ]          */
	in []int16, /* I:   Input signal [ len ]                    */
	inLen int32, /* I:   Number of input samples                 */
) {
	len4 := inLen >> 2
	SKP_assert(SKP_Silk_resampler_down2_0 > 0)
	SKP_assert(SKP_Silk_resampler_down2_1 < 0)

	/* Internal variables and state are in Q10 format */
	for k := int32(0); k < len4; k++ {
		/* Add two input samples and convert to Q10 */
		in32 := (int32(in[k*4]) + int32(in[k*4+1])) << 9

		/* All-pass section for even input sample */
		Y := in32 - S[0]
		X := SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_down2_1))
		out32 := S[0] + X
		S[0] = in32 + X

		/* Add two input samples and convert to Q10 */
		in32 = (int32(in[k*4+2]) + int32(in[k*4+3])) << 9

		/* All-pass section for odd input sample */
		Y = in32 - S[1]
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_down2_0))
		out32 = out32 + S[1]
		out32 = out32 + X
		S[1] = in32 + X

		/* Add, convert back to int16 and store to output */
		out[k] = SKP_SAT16(SKP_RSHIFT_ROUND(out32, 11))
	}
}
