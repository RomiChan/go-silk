package silk

import "math"

var (
	sigm_LUT_slope_Q10 = [6]int32{237, 153, 73, 30, 12, 7}
	sigm_LUT_pos_Q15   = [6]int32{0x4000, 0x5D93, 0x70BD, 0x79ED, 0x7DB2, 0x7F24}
	sigm_LUT_neg_Q15   = [6]int32{0x4000, 8812, 3906, 1554, 589, 219}
)

func SKP_Silk_sigm_Q15(in_Q5 int32) int32 {
	var ind int32
	if int64(in_Q5) < 0 {
		in_Q5 = -in_Q5
		if int64(in_Q5) >= 6*32 {
			return 0
		} else {
			ind = in_Q5 >> 5
			return int32(int64(sigm_LUT_neg_Q15[ind]) - int64(SKP_SMULBB(sigm_LUT_slope_Q10[ind], int32(int64(in_Q5)&31))))
		}
	} else {
		if int64(in_Q5) >= 6*32 {
			return math.MaxInt16
		} else {
			ind = in_Q5 >> 5
			return int32(int64(sigm_LUT_pos_Q15[ind]) + int64(SKP_SMULBB(sigm_LUT_slope_Q10[ind], int32(int64(in_Q5)&31))))
		}
	}
}
