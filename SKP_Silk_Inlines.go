package silk

import "math"

const (
	SKP_SIN_APPROX_CONST0 = 1073735400
	SKP_SIN_APPROX_CONST1 = -82778932
	SKP_SIN_APPROX_CONST2 = 0x102AF9
	SKP_SIN_APPROX_CONST3 = -5013
)

func SKP_Silk_CLZ_FRAC(in int32, lz *int32, frac_Q7 *int32) {
	var lzeros int32 = SKP_Silk_CLZ32(in)
	*lz = lzeros
	*frac_Q7 = int32(int64(SKP_ROR32(in, int32(24-int64(lzeros)))) & math.MaxInt8)
}

func SKP_Silk_SQRT_APPROX(x int32) int32 {
	var (
		y       int32
		lz      int32
		frac_Q7 int32
	)
	if int64(x) <= 0 {
		return 0
	}
	SKP_Silk_CLZ_FRAC(x, &lz, &frac_Q7)
	if int64(lz)&1 != 0 {
		y = 0x8000
	} else {
		y = 0xB486
	}
	y >>= lz >> 1
	y = SKP_SMLAWB(y, y, SKP_SMULBB(213, frac_Q7))
	return y
}

func SKP_Silk_norm16(a int16) int32 {
	var a32 int32
	if (int64(a) << 1) == 0 {
		return 0
	}
	a32 = int32(a)
	a32 ^= a32 >> 31
	return int32(int64(SKP_Silk_CLZ32(a32)) - 17)
}

func SKP_Silk_norm32(a int32) int32 {
	if (int64(a) << 1) == 0 {
		return 0
	}
	a ^= a >> 31
	return int32(int64(SKP_Silk_CLZ32(a)) - 1)
}

func SKP_DIV32_varQ(a32 int32, b32 int32, Qres int32) int32 {
	var (
		a_headrm int32
		b_headrm int32
		lshift   int32
		b32_inv  int32
		a32_nrm  int32
		b32_nrm  int32
		result   int32
	)
	a_headrm = int32(int64(SKP_Silk_CLZ32(int32(SKP_abs(int64(a32))))) - 1)
	a32_nrm = int32(int64(a32) << int64(a_headrm))
	b_headrm = int32(int64(SKP_Silk_CLZ32(int32(SKP_abs(int64(b32))))) - 1)
	b32_nrm = int32(int64(b32) << int64(b_headrm))
	b32_inv = int32((SKP_int32_MAX >> 2) / (int64(b32_nrm) >> 16))
	result = SKP_SMULWB(a32_nrm, b32_inv)
	a32_nrm -= int32(int64(SKP_SMMUL(b32_nrm, result)) << 3)
	result = SKP_SMLAWB(result, a32_nrm, b32_inv)
	lshift = int32(int64(a_headrm) + 29 - int64(b_headrm) - int64(Qres))
	if int64(lshift) <= 0 {
		return int32((func() int64 {
			if (int64(math.MinInt32) >> int64(-lshift)) > (SKP_int32_MAX >> int64(-lshift)) {
				if int64(result) > (int64(math.MinInt32) >> int64(-lshift)) {
					return int64(math.MinInt32) >> int64(-lshift)
				}
				if int64(result) < (SKP_int32_MAX >> int64(-lshift)) {
					return SKP_int32_MAX >> int64(-lshift)
				}
				return int64(result)
			}
			if int64(result) > (SKP_int32_MAX >> int64(-lshift)) {
				return SKP_int32_MAX >> int64(-lshift)
			}
			if int64(result) < (int64(math.MinInt32) >> int64(-lshift)) {
				return int64(math.MinInt32) >> int64(-lshift)
			}
			return int64(result)
		}()) << int64(-lshift))
	} else {
		if int64(lshift) < 32 {
			return result >> int64(lshift)
		} else {
			return 0
		}
	}
}

func SKP_INVERSE32_varQ(b32 int32, Qres int32) int32 {
	var (
		b_headrm int32
		lshift   int32
		b32_inv  int32
		b32_nrm  int32
		err_Q32  int32
		result   int32
	)
	b_headrm = int32(int64(SKP_Silk_CLZ32(int32(SKP_abs(int64(b32))))) - 1)
	b32_nrm = int32(int64(b32) << int64(b_headrm))
	b32_inv = int32((SKP_int32_MAX >> 2) / (int64(b32_nrm) >> 16))
	result = int32(int64(b32_inv) << 16)
	err_Q32 = int32(int64(-SKP_SMULWB(b32_nrm, b32_inv)) << 3)
	result = SKP_SMLAWW(result, err_Q32, b32_inv)
	lshift = int32(61 - int64(b_headrm) - int64(Qres))
	if int64(lshift) <= 0 {
		return int32((func() int64 {
			if (int64(math.MinInt32) >> int64(-lshift)) > (SKP_int32_MAX >> int64(-lshift)) {
				if int64(result) > (int64(math.MinInt32) >> int64(-lshift)) {
					return int64(math.MinInt32) >> int64(-lshift)
				}
				if int64(result) < (SKP_int32_MAX >> int64(-lshift)) {
					return SKP_int32_MAX >> int64(-lshift)
				}
				return int64(result)
			}
			if int64(result) > (SKP_int32_MAX >> int64(-lshift)) {
				return SKP_int32_MAX >> int64(-lshift)
			}
			if int64(result) < (int64(math.MinInt32) >> int64(-lshift)) {
				return int64(math.MinInt32) >> int64(-lshift)
			}
			return int64(result)
		}()) << int64(-lshift))
	} else {
		if int64(lshift) < 32 {
			return result >> int64(lshift)
		} else {
			return 0
		}
	}
}

func SKP_Silk_SIN_APPROX_Q24(x int32) int32 {
	var y_Q30 int32
	x &= math.MaxUint16
	if int64(x) <= 0x8000 {
		if int64(x) < 0x4000 {
			x = int32(0x4000 - int64(x))
		} else {
			x -= 0x4000
		}
		if int64(x) < 1100 {
			return SKP_SMLAWB(1<<24, int32(int64(x)*int64(x)), -5053)
		}
		x = SKP_SMULWB(int32(int64(x)<<8), x)
		y_Q30 = SKP_SMLAWB(0x102AF9, x, int32(-5013))
		y_Q30 = SKP_SMLAWW(-82778932, x, y_Q30)
		y_Q30 = SKP_SMLAWW(1073735400+66, x, y_Q30)
	} else {
		if int64(x) < 0xC000 {
			x = int32(0xC000 - int64(x))
		} else {
			x -= 0xC000
		}
		if int64(x) < 1100 {
			return SKP_SMLAWB(int32(-1<<24), int32(int64(x)*int64(x)), 5053)
		}
		x = SKP_SMULWB(int32(int64(x)<<8), x)
		y_Q30 = SKP_SMLAWB(-1059577, x, int32(-(-5013)))
		y_Q30 = SKP_SMLAWW(-(-82778932), x, y_Q30)
		y_Q30 = SKP_SMLAWW(-1073735400, x, y_Q30)
	}
	return SKP_RSHIFT_ROUND(y_Q30, 6)
}

func SKP_Silk_COS_APPROX_Q24(x int32) int32 {
	return SKP_Silk_SIN_APPROX_Q24(int32(int64(x) + 0x4000))
}
