package silk

func SKP_Silk_MA_Prediction(in []int16, B []int16, S []int32, out []int16, len_ int32, order int32) {
	var (
		k     int32
		d     int32
		in16  int32
		out32 int32
	)
	for k = 0; k < len_; k++ {
		in16 = int32(in[k])
		out32 = (in16 << 12) - S[0]
		out32 = SKP_RSHIFT_ROUND(out32, 12)
		for d = 0; d < order-1; d++ {
			S[d] = int32(uint32(S[d+1]) + uint32(SKP_SMULBB(in16, int32(B[d]))))
		}
		S[order-1] = SKP_SMULBB(in16, int32(B[order-1]))
		out[k] = SKP_SAT16(out32)
	}
}
func SKP_Silk_LPC_analysis_filter(in []int16, B []int16, S []int16, out []int16, len_ int32, Order int32) {
	var (
		k          int32
		j          int32
		idx        int32
		Order_half = Order >> 1
		out32_Q12  int32
		out32      int32
		SA         int16
		SB         int16
	)
	for k = 0; k < len_; k++ {
		SA = S[0]
		out32_Q12 = 0
		for j = 0; j < (Order_half - 1); j++ {
			idx = SKP_SMULBB(2, j) + 1
			SB = S[idx]
			S[idx] = SA
			out32_Q12 = SKP_SMLABB(out32_Q12, int32(SA), int32(B[idx-1]))
			out32_Q12 = SKP_SMLABB(out32_Q12, int32(SB), int32(B[idx]))
			SA = S[idx+1]
			S[idx+1] = SB
		}
		SB = S[Order-1]
		S[Order-1] = SA
		out32_Q12 = SKP_SMLABB(out32_Q12, int32(SA), int32(B[Order-2]))
		out32_Q12 = SKP_SMLABB(out32_Q12, int32(SB), int32(B[Order-1]))
		out32_Q12 = SKP_SUB_SAT32((int32(in[k]))<<12, out32_Q12)
		out32 = SKP_RSHIFT_ROUND(out32_Q12, 12)
		out[k] = SKP_SAT16(out32)
		S[0] = in[k]
	}
}
