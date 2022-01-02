package silk

func SKP_Silk_VQ_WMat_EC_FIX(ind *int32, rate_dist_Q14 *int32, in_Q14 []int16, W_Q18 []int32, cb_Q14 *int16, cl_Q6 []int16, mu_Q8 int32, L int32) {
	var (
		k          int32
		cb_row_Q14 []int16
		diff_Q14   [5]int16
		sum1_Q14   int32
		sum2_Q16   int32
	)
	*rate_dist_Q14 = SKP_int32_MAX
	cb_row_Q14 = ([]int16)(cb_Q14)
	for k = 0; k < L; k++ {
		diff_Q14[0] = in_Q14[0] - cb_row_Q14[0]
		diff_Q14[1] = in_Q14[1] - cb_row_Q14[1]
		diff_Q14[2] = in_Q14[2] - cb_row_Q14[2]
		diff_Q14[3] = in_Q14[3] - cb_row_Q14[3]
		diff_Q14[4] = in_Q14[4] - cb_row_Q14[4]
		sum1_Q14 = SKP_SMULBB(mu_Q8, int32(cl_Q6[k]))
		sum2_Q16 = SKP_SMULWB(W_Q18[1], int32(diff_Q14[1]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[2], int32(diff_Q14[2]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[3], int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[4], int32(diff_Q14[4]))
		sum2_Q16 = sum2_Q16 << 1
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[0], int32(diff_Q14[0]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[0]))
		sum2_Q16 = SKP_SMULWB(W_Q18[7], int32(diff_Q14[2]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[8], int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[9], int32(diff_Q14[4]))
		sum2_Q16 = sum2_Q16 << 1
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[6], int32(diff_Q14[1]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[1]))
		sum2_Q16 = SKP_SMULWB(W_Q18[13], int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[14], int32(diff_Q14[4]))
		sum2_Q16 = sum2_Q16 << 1
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[12], int32(diff_Q14[2]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[2]))
		sum2_Q16 = SKP_SMULWB(W_Q18[19], int32(diff_Q14[4]))
		sum2_Q16 = sum2_Q16 << 1
		sum2_Q16 = SKP_SMLAWB(sum2_Q16, W_Q18[18], int32(diff_Q14[3]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[3]))
		sum2_Q16 = SKP_SMULWB(W_Q18[24], int32(diff_Q14[4]))
		sum1_Q14 = SKP_SMLAWB(sum1_Q14, sum2_Q16, int32(diff_Q14[4]))
		if sum1_Q14 < *rate_dist_Q14 {
			*rate_dist_Q14 = sum1_Q14
			*ind = k
		}
		cb_row_Q14 += LTP_ORDER
	}
}
