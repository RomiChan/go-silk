package silk

// SKP_Silk_biquad see biquad.c
func SKP_Silk_biquad(in []int16, B [3]int16, A [2]int16, S []int32, out []int16) {
	S0 := S[0]
	S1 := S[1]
	A0_neg := int32(A[0])
	A1_neg := int32(A[1])
	for k := range in {
		in16 := int32(in[k])
		out32 := SKP_SMLABB(S0, in16, int32(B[0]))
		S0 = SKP_SMLABB(S1, in16, int32(B[1]))
		S0 += SKP_SMULWB(out32, A0_neg) << 3
		S1 = SKP_SMULWB(out32, A1_neg) << 3
		S1 = SKP_SMLABB(S1, in16, int32(B[2]))
		tmp32 := int32(int64(SKP_RSHIFT_ROUND(out32, 13)) + 1)
		out[k] = SKP_SAT16(tmp32)
	}
	S[0] = S0
	S[1] = S1
}

// SKP_Silk_biquad_alt see biquad_alt.c
func SKP_Silk_biquad_alt(in []int16, B_Q28 [3]int32, A_Q28 [2]int32, S []int32, out []int16) {
	A0_L_Q28 := int32(int64(A_Q28[0]) & 0x3FFF)
	A0_U_Q28 := (A_Q28[0]) >> 14
	A1_L_Q28 := int32(int64(A_Q28[1]) & 0x3FFF)
	A1_U_Q28 := (A_Q28[1]) >> 14
	for k := range in {
		inval := int32(in[k])
		out32_Q14 := SKP_SMLAWB(S[0], B_Q28[0], inval) << 2
		S[0] = int32(int64(S[1]) + int64(SKP_RSHIFT_ROUND(SKP_SMULWB(out32_Q14, A0_L_Q28), 14)))
		S[0] = SKP_SMLAWB(S[0], out32_Q14, A0_U_Q28)
		S[0] = SKP_SMLAWB(S[0], B_Q28[1], inval)
		S[1] = SKP_RSHIFT_ROUND(SKP_SMULWB(out32_Q14, A1_L_Q28), 14)
		S[1] = SKP_SMLAWB(S[1], out32_Q14, A1_U_Q28)
		S[1] = SKP_SMLAWB(S[1], B_Q28[2], inval)
		out[k] = SKP_SAT16(int32(int64(out32_Q14)+(1<<14)-1) >> 14)
	}
}
