package silk

import "math"

func SKP_Silk_log2lin(inLog_Q7 int32) int32 {
	if int64(inLog_Q7) < 0 {
		return 0
	} else if int64(inLog_Q7) >= (31 << 7) {
		return SKP_int32_MAX
	}
	out := int32(1 << (inLog_Q7 >> 7))
	frac_Q7 := int32(int64(inLog_Q7) & math.MaxInt8)
	if int64(inLog_Q7) < 2048 {
		out = int32(int64(out) + ((int64(out) * int64(SKP_SMLAWB(frac_Q7, int32(int64(frac_Q7)*(128-int64(frac_Q7))), -174))) >> 7))
	} else {
		out = int32(int64(out) + (int64(out)>>7)*int64(SKP_SMLAWB(frac_Q7, int32(int64(frac_Q7)*(128-int64(frac_Q7))), -174)))
	}
	return out
}

func SKP_Silk_lin2log(inLin int32) int32 {
	var lz, frac_Q7 int32
	SKP_Silk_CLZ_FRAC(inLin, &lz, &frac_Q7)
	return int32(((31 - int64(lz)) << 7) + int64(SKP_SMLAWB(frac_Q7, int32(int64(frac_Q7)*(128-int64(frac_Q7))), 179)))
}
