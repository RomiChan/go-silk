package silk

// reviewed by wdvxdr1123 2022-02-08

import "math"

// SKP_Silk_log2lin Approximation of 2^() (very close inverse of SKP_Silk_lin2log())
// Convert input to a linear scale
func SKP_Silk_log2lin(inLog_Q7 int32) int32 {
	if inLog_Q7 < 0 {
		return 0
	} else if inLog_Q7 >= (31 << 7) {
		/* Saturate, and prevent wrap-around */
		return SKP_int32_MAX
	}
	out := int32(1 << (inLog_Q7 >> 7))
	frac_Q7 := inLog_Q7 & math.MaxInt8
	if inLog_Q7 < 2048 {
		/* Piece-wise parabolic approximation */
		out = out + ((out * SKP_SMLAWB(frac_Q7, frac_Q7*(128-frac_Q7), -174)) >> 7)
	} else {
		/* Piece-wise parabolic approximation */
		out = out + ((out)>>7)*(SKP_SMLAWB(frac_Q7, frac_Q7*(128-frac_Q7), -174))
	}
	return out
}

// SKP_Silk_lin2log Approximation of 128 * log2() (very close inverse of approx 2^() below)
// Convert input to a log scale
func SKP_Silk_lin2log(inLin int32) int32 {
	var lz, frac_Q7 int32
	SKP_Silk_CLZ_FRAC(inLin, &lz, &frac_Q7)
	/* Piece-wise parabolic approximation */
	return ((31 - lz) << 7) + SKP_SMLAWB(frac_Q7, frac_Q7*(128-frac_Q7), 179)
}
