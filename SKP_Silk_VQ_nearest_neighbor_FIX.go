package silk

// Reviewed by wdvxdr1123 2022-01-02

// SKP_Silk_VQ_WMat_EC_FIX
// Entropy constrained MATRIX-weighted VQ, hard-coded to 5-element vectors, for a single input data vector
func SKP_Silk_VQ_WMat_EC_FIX(
	ind *int32, /* O    index of best codebook vector               */
	rate_dist_Q14 *int32, /* O    best weighted quantization error + mu * rate*/
	in_Q14 []int16, /* I    input vector to be quantized                */
	W_Q18 []int32, /* I    weighting matrix                            */
	cb_Q14 []int16, /* I    codebook                                    */
	cl_Q6 []int16, /* I    code length for each codebook vector        */
	mu_Q8 int32, /* I    tradeoff between weighted error and rate    */
	L int32, /* I    number of vectors in codebook               */
) {
	var diff_Q14 [5]int16

	_, _ = in_Q14[4], W_Q18[24] // early bounds check

	/* Loop over codebook */
	*rate_dist_Q14 = SKP_int32_MAX
	cb_row_Q14 := cb_Q14
	for k := int32(0); k < L; k++ {
		_ = cb_row_Q14[4] // early bounds check
		diff_Q14[0] = in_Q14[0] - cb_row_Q14[0]
		diff_Q14[1] = in_Q14[1] - cb_row_Q14[1]
		diff_Q14[2] = in_Q14[2] - cb_row_Q14[2]
		diff_Q14[3] = in_Q14[3] - cb_row_Q14[3]
		diff_Q14[4] = in_Q14[4] - cb_row_Q14[4]

		/* Weighted rate */
		sum1_Q14 := SKP_SMULBB(mu_Q8, int32(cl_Q6[k]))

		SKP_assert(sum1_Q14 >= 0)

		/* first row of W_Q18 */
		sum2_Q16 := SKP_SMULWB(W_Q18[1], int32(diff_Q14[1]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[2], int32(diff_Q14[2]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[3], int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[4], int32(diff_Q14[4]))
		sum2_Q16 = sum2_Q16 << 1
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[0], int32(diff_Q14[0]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[0]))

		/* second row of W_Q18 */
		sum2_Q16 = SKP_SMULWB(W_Q18[7], int32(diff_Q14[2]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[8], int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[9], int32(diff_Q14[4]))
		sum2_Q16 = sum2_Q16 << 1
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[6], int32(diff_Q14[1]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[1]))

		/* third row of W_Q18 */
		sum2_Q16 = SKP_SMULWB(W_Q18[13], int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[14], int32(diff_Q14[4]))
		sum2_Q16 = sum2_Q16 << 1
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[12], int32(diff_Q14[2]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[2]))

		/* fourth row of W_Q18 */
		sum2_Q16 = SKP_SMULWB(W_Q18[19], int32(diff_Q14[4]))
		sum2_Q16 = sum2_Q16 << 1
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[18], int32(diff_Q14[3]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[3]))

		/* last row of W_Q18 */
		sum2_Q16 = SKP_SMULWB(W_Q18[24], int32(diff_Q14[4]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[4]))
		SKP_assert(sum1_Q14 >= 0)

		/* find best */
		if sum1_Q14 < *rate_dist_Q14 {
			*rate_dist_Q14 = sum1_Q14
			*ind = k
		}
		/* Go to next cbk vector */
		cb_row_Q14 = cb_row_Q14[LTP_ORDER:]
	}
}
