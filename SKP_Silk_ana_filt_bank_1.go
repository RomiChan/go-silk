package silk

var (
	A_fb1_20 = [1]int16{5394 << 1}
	A_fb1_21 = [1]int16{int16(int32(0x508F << 1))}
)

// see SKP_Silk_ana_filt_bank_1.c
// S -> [2]int32
func ana_filt_bank_1(in []int16, S []int32, outL []int16, outH []int16) {
	N2 := len(in) >> 1
	for k := 0; k < N2; k++ {
		/* Convert to Q10 */
		in32 := int32(in[k*2] << 10)

		Y := int32(int64(in32) - int64(S[0]))
		X := SKP_SMLAWB(Y, Y, int32(A_fb1_21[0]))
		out_1 := int32(int64(S[0]) + int64(X))
		S[0] = int32(int64(in32) + int64(X))
		in32 = int32(in[int64(k)*2+1] << 10)
		Y = int32(int64(in32) - int64(S[1]))
		X = SKP_SMULWB(Y, int32(A_fb1_20[0]))
		out_2 := int32(int64(S[1]) + int64(X))
		S[1] = int32(int64(in32) + int64(X))
		outL[k] = SKP_SAT16(SKP_RSHIFT_ROUND(int32(int64(out_2)+int64(out_1)), 11))
		outH[k] = SKP_SAT16(SKP_RSHIFT_ROUND(int32(int64(out_2)-int64(out_1)), 11))
	}
}
