package silk

// reviewed by wdvxdr1123 2022-02-08

import "math"

const SKP_SIN_APPROX_CONST0 = 1073735400
const SKP_SIN_APPROX_CONST1 = -82778932
const SKP_SIN_APPROX_CONST2 = 0x102AF9
const SKP_SIN_APPROX_CONST3 = -5013

func SKP_Silk_CLZ_FRAC(in int32, lz *int32, frac_Q7 *int32) {
	lzeros := SKP_Silk_CLZ32(in)
	*lz = lzeros
	*frac_Q7 = SKP_ROR32(in, 24-lzeros) & math.MaxInt8
}

func SKP_Silk_SQRT_APPROX(x int32) int32 {
	var y, lz, frac_Q7 int32
	if x <= 0 {
		return 0
	}
	SKP_Silk_CLZ_FRAC(x, &lz, &frac_Q7)
	if lz&1 != 0 {
		y = 0x8000
	} else {
		y = 0xB486
	}
	y >>= lz >> 1
	return SKP_SMLAWB(y, y, SKP_SMULBB(213, frac_Q7))
}

func SKP_Silk_norm16(a int16) int32 {
	if (a << 1) == 0 {
		return 0
	}
	a32 := int32(a)
	a32 ^= a32 >> 31
	return SKP_Silk_CLZ32(a32) - 17
}

func SKP_Silk_norm32(a int32) int32 {
	if (a << 1) == 0 {
		return 0
	}
	a ^= a >> 31
	return SKP_Silk_CLZ32(a) - 1
}

func SKP_DIV32_varQ(a32 int32, b32 int32, Qres int32) int32 {
	SKP_assert(b32 != 0)
	SKP_assert(Qres >= 0)
	a_headrm := SKP_Silk_CLZ32(int32(SKP_abs(int64(a32)))) - 1
	a32_nrm := a32 << a_headrm
	b_headrm := SKP_Silk_CLZ32(int32(SKP_abs(int64(b32)))) - 1
	b32_nrm := b32 << b_headrm
	b32_inv := (SKP_int32_MAX >> 2) / (b32_nrm >> 16)
	result := SKP_SMULWB(a32_nrm, b32_inv)
	a32_nrm -= SKP_SMMUL(b32_nrm, result) << 3
	result = SKP_SMLAWB(result, a32_nrm, b32_inv)
	lshift := a_headrm + 29 - b_headrm - Qres
	if lshift <= 0 {
		return SKP_LSHIFT_SAT32(result, -lshift)
	} else {
		if lshift < 32 {
			return result >> lshift
		} else {
			return 0
		}
	}
}
func SKP_INVERSE32_varQ(b32 int32, Qres int32) int32 {
	SKP_assert(b32 != 0)
	// SKP_assert(b32 != 0x80000000)
	SKP_assert(Qres > 0)
	b_headrm := SKP_Silk_CLZ32(int32(SKP_abs(int64(b32)))) - 1
	b32_nrm := b32 << b_headrm
	b32_inv := (SKP_int32_MAX >> 2) / (b32_nrm >> 16)
	result := b32_inv << 16
	err_Q32 := (-SKP_SMULWB(b32_nrm, b32_inv)) << 3
	result = SKP_SMLAWW(result, err_Q32, b32_inv)
	lshift := 61 - b_headrm - Qres
	if lshift <= 0 {
		return SKP_LSHIFT_SAT32(result, -lshift)
	} else {
		if lshift < 32 {
			return result >> lshift
		} else {
			return 0
		}
	}
}
func SKP_Silk_SIN_APPROX_Q24(x int32) int32 {
	var y_Q30 int32
	x &= math.MaxUint16
	if x <= 0x8000 {
		if x < 0x4000 {
			x = 0x4000 - x
		} else {
			x -= 0x4000
		}
		if x < 1100 {
			return SKP_SMLAWB(1<<24, x*x, -5053)
		}
		x = SKP_SMULWB(x<<8, x)
		y_Q30 = SKP_SMLAWB(0x102AF9, x, int32(-5013))
		y_Q30 = SKP_SMLAWW(-82778932, x, y_Q30)
		y_Q30 = SKP_SMLAWW(1073735400+66, x, y_Q30)
	} else {
		if x < 0xC000 {
			x = 0xC000 - x
		} else {
			x -= 0xC000
		}
		if x < 1100 {
			return SKP_SMLAWB(int32(-1<<24), x*x, 5053)
		}
		x = SKP_SMULWB(x<<8, x)
		y_Q30 = SKP_SMLAWB(-1059577, x, int32(-(-5013)))
		y_Q30 = SKP_SMLAWW(-(-82778932), x, y_Q30)
		y_Q30 = SKP_SMLAWW(-1073735400, x, y_Q30)
	}
	return SKP_RSHIFT_ROUND(y_Q30, 6)
}

func SKP_Silk_COS_APPROX_Q24(x int32) int32 {
	return SKP_Silk_SIN_APPROX_Q24(x + 0x4000)
}
