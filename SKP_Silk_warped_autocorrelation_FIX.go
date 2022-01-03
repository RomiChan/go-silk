package silk

// Reviewed by wdvxdr1123 2022-01-03

const QC = 10
const QS = 14

// SKP_Silk_warped_autocorrelation_FIX
// Autocorrelations for a warped frequency axis
func SKP_Silk_warped_autocorrelation_FIX(
	corr []int32, /* O    Result [order + 1]                      */
	scale *int32, /* O    Scaling of the correlation vector       */
	input []int16, /* I    Input data to correlate                 */
	warping_Q16 int16, /* I    Warping coefficient                     */
	length int32, /* I    Length of input                         */
	order int32, /* I    Correlation order (even)                */
) {
	var n, i int32
	var state_QS [MAX_SHAPE_LPC_ORDER + 1]int32
	var corr_QC [MAX_SHAPE_LPC_ORDER + 1]int64

	/* Order must be even */
	SKP_assert((order & 1) == 0)
	SKP_assert(QS*2-QC >= 0)

	/* Loop over samples */
	for n = 0; n < length; n++ {
		tmp1_QS := (int32(input[n])) << QS
		/* Loop over allpass sections */
		for i = 0; i < order; i += 2 {
			/* Output of allpass section */
			tmp2_QS := SKP_SMLAWB(state_QS[i], state_QS[i+1]-tmp1_QS, int32(warping_Q16))
			state_QS[i] = tmp1_QS
			corr_QC[i] += (int64(tmp1_QS) * int64(state_QS[0])) >> (QS*2 - QC)
			/* Output of allpass section */
			tmp1_QS = SKP_SMLAWB(state_QS[i+1], state_QS[i+2]-tmp2_QS, int32(warping_Q16))
			state_QS[i+1] = tmp2_QS
			corr_QC[i+1] += (int64(tmp2_QS) * int64(state_QS[0])) >> (QS*2 - QC)
		}
		state_QS[order] = tmp1_QS
		corr_QC[order] += (int64(tmp1_QS) * int64(state_QS[0])) >> (QS*2 - QC)
	}
	lsh := SKP_Silk_CLZ64(corr_QC[0]) - 35
	if lsh > (30 - QC) {
		lsh = 30 - QC
	} else if lsh < -12-QC {
		lsh = int32(-12 - QC)
	}
	*scale = -(QC + lsh)
	SKP_assert(int64(*scale) >= -30 && *scale <= 12)
	if lsh >= 0 {
		for i = 0; i < order+1; i++ {
			corr[i] = int32((corr_QC[i]) << lsh)
		}
	} else {
		for i = 0; i < order+1; i++ {
			corr[i] = int32((corr_QC[i]) >> -lsh)
		}
	}
	SKP_assert(corr_QC[0] >= 0) // If breaking, decrease QC
}
