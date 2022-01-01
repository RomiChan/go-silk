package silk

import "math"

func SKP_Silk_LPC_synthesis_filter(in []int16, A_Q12 []int16, Gain_Q26 int32, S []int32, out []int16, len_ int32, Order int32) {
	var (
		k          int32
		j          int32
		idx        int32
		Order_half = Order >> 1
		SA         int32
		SB         int32
		out32_Q10  int32
		out32      int32
	)
	for k = 0; k < len_; k++ {
		SA = S[Order-1]
		out32_Q10 = 0
		for j = 0; j < (Order_half - 1); j++ {
			idx = SKP_SMULBB(2, j) + 1
			SB = S[Order-1-idx]
			S[Order-1-idx] = SA
			out32_Q10 = SKP_SMLAWB(out32_Q10, SA, int32(A_Q12[j<<1]))
			out32_Q10 = SKP_SMLAWB(out32_Q10, SB, int32(A_Q12[(j<<1)+1]))
			SA = S[Order-2-idx]
			S[Order-2-idx] = SB
		}
		SB = S[0]
		S[0] = SA
		out32_Q10 = SKP_SMLAWB(out32_Q10, SA, int32(A_Q12[Order-2]))
		out32_Q10 = SKP_SMLAWB(out32_Q10, SB, int32(A_Q12[Order-1]))
		if ((out32_Q10 + SKP_SMULWB(Gain_Q26, int32(in[k]))) & math.MinInt32) == 0 {
			if ((out32_Q10 & SKP_SMULWB(Gain_Q26, int32(in[k]))) & math.MinInt32) != 0 {
				out32_Q10 = math.MinInt32
			} else {
				out32_Q10 = out32_Q10 + SKP_SMULWB(Gain_Q26, int32(in[k]))
			}
		} else if ((out32_Q10 | SKP_SMULWB(Gain_Q26, int32(in[k]))) & math.MinInt32) == 0 {
			out32_Q10 = SKP_int32_MAX
		} else {
			out32_Q10 = out32_Q10 + SKP_SMULWB(Gain_Q26, int32(in[k]))
		}
		out32 = SKP_RSHIFT_ROUND(out32_Q10, 10)
		out[k] = SKP_SAT16(out32)
		S[Order-1] = SKP_LSHIFT_SAT32(out32_Q10, 4)
	}
}
