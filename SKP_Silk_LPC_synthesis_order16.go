package silk

// reviewed by wdvxdr1123 2022-02-08

// SKP_Silk_LPC_synthesis_order16 16th order AR filter
func SKP_Silk_LPC_synthesis_order16(in []int16, A_Q12 []int16, Gain_Q26 int32, S []int32, out []int16, len_ int32) {
	for k := int32(0); k < len_; k++ {
		/* unrolled loop: prolog */
		/* multiply-add two prediction coefficients per iteration */
		SA := S[15]
		SB := S[14]
		S[14] = SA
		out32_Q10 := SKP_SMULWB(SA, int32(A_Q12[0]))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(A_Q12[1]))))
		SA = S[13]
		S[13] = SB

		/* unrolled loop: main loop */
		SB = S[12]
		S[12] = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(A_Q12[2]))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(A_Q12[3]))))
		SA = S[11]
		S[11] = SB

		SB = S[10]
		S[10] = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(A_Q12[4]))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(A_Q12[5]))))
		SA = S[9]
		S[9] = SB

		SB = S[8]
		S[8] = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(A_Q12[6]))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(A_Q12[7]))))
		SA = S[7]
		S[7] = SB

		SB = S[6]
		S[6] = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(A_Q12[8]))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(A_Q12[9]))))
		SA = S[5]
		S[5] = SB

		SB = S[4]
		S[4] = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(A_Q12[10]))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(A_Q12[11]))))
		SA = S[3]
		S[3] = SB

		SB = S[2]
		S[2] = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(A_Q12[12]))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(A_Q12[13]))))
		SA = S[1]
		S[1] = SB

		/* unrolled loop: epilog */
		SB = S[0]
		S[0] = SA
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SA, int32(A_Q12[14]))))
		out32_Q10 = int32(uint32(out32_Q10) + uint32(SKP_SMULWB(SB, int32(A_Q12[15]))))

		/* unrolled loop: end */
		/* apply gain to excitation signal and add to prediction */
		out32_Q10 = SKP_ADD_SAT32(out32_Q10, SKP_SMULWB(Gain_Q26, int32(in[k])))

		/* scale to Q0 */
		out32 := SKP_RSHIFT_ROUND(out32_Q10, 10)

		/* saturate output */
		out[k] = SKP_SAT16(out32)

		/* move result into delay line */
		S[15] = SKP_LSHIFT_SAT32(out32_Q10, 4)
	}
}
