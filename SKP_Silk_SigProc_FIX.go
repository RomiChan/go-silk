package silk

import "math"

const SKP_Silk_MAX_ORDER_LPC = 16
const SKP_Silk_MAX_CORRELATION_LENGTH = 640
const SKP_Silk_PITCH_EST_MIN_COMPLEX = 0
const SKP_Silk_PITCH_EST_MID_COMPLEX = 1
const SKP_Silk_PITCH_EST_MAX_COMPLEX = 2
const LSF_COS_TAB_SZ_FIX = 128

func SKP_ROR32(a32 int32, rot int32) int32 {
	var (
		x uint32 = uint32(a32)
		r uint32 = uint32(rot)
		m uint32 = uint32(-rot)
	)
	if rot <= 0 {
		return int32((x << m) | x>>(32-m))
	} else {
		return int32((x << (32 - r)) | x>>r)
	}
}

func SKP_SMULTT(a32 int32, b32 int32) int32 {
	return (a32 >> 16) * (b32 >> 16)
}

func SKP_SMLATT(a32 int32, b32 int32, c32 int32) int32 {
	return a32 + (b32>>16)*(c32>>16)
}

func SKP_SAT16(a int32) int16 {
	if a > SKP_int16_MAX {
		return SKP_int16_MAX
	}
	if a < 0x8000 {
		return math.MinInt16
	}
	return int16(a)
}

func SKP_ADD_POS_SAT32(a int32, b int32) int32 {
	if ((int64(a) + int64(b)) & int64(math.MinInt32)) != 0 {
		return SKP_int32_MAX
	}
	return a + b
}

func SKP_LIMIT_32(a int32, limit1 int32, limit2 int32) int32 {
	if limit1 > limit2 {
		if a > limit1 {
			return limit1
		}
		if a < limit2 {
			return limit2
		}
		return a
	}
	if a > limit2 {
		return limit2
	}
	if a < limit1 {
		return limit1
	}
	return a
}

func SKP_LSHIFT_SAT32(a int32, shift int32) int32 {
	return SKP_LIMIT_32(a, math.MinInt32>>shift, SKP_int32_MAX>>shift) << shift
}

func SKP_RSHIFT_ROUND(a int32, shift int32) int32 {
	if shift == 1 {
		return (a >> 1) + (a & 1)
	}
	return ((a >> (shift - 1)) + 1) >> 1
}

func SKP_RSHIFT_ROUND64(a int64, shift int32) int32 {
	if shift == 1 {
		return int32((a >> 1) + (a & 1))
	}
	return int32(((a >> int64(shift-1)) + 1) >> 1)
}

func SKP_FIX_CONST(C float64, Q int32) int32 {
	return int32(C*float64(1<<Q) + 0.5)
}

func SKP_min_int(a int32, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func SKP_min_32(a int32, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func SKP_max_int(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func SKP_max_16(a int16, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func SKP_max_32(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func SKP_LIMIT_int(a int32, limit1 int32, limit2 int32) int32 {
	if limit1 > limit2 {
		if a > limit1 {
			return limit1
		}
		if a < limit2 {
			return limit2
		}
		return a
	}
	if a > limit2 {
		return limit2
	}
	if a < limit1 {
		return limit1
	}
	return a
}
func SKP_abs(a int64) int64 {
	if a > 0 {
		return a
	}
	return -a
}
