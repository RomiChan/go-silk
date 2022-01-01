package silk

// see SKP_Silk_array_maxabs.c
// len_ -> len(vec)
func int16_array_maxabs(vec []int16) int16 {
	if len(vec) == 0 {
		return 0
	}
	ind := len(vec) - 1
	max := SKP_SMULBB(int32(vec[ind]), int32(vec[ind]))
	for i := len(vec) - 2; i >= 0; i-- {
		lvl := SKP_SMULBB(int32(vec[i]), int32(vec[i]))
		if lvl > max {
			max = lvl
			ind = i
		}
	}
	if int64(max) >= 0x3FFF0001 {
		return SKP_int16_MAX
	} else {
		if int64(vec[ind]) < 0 {
			return -vec[ind]
		} else {
			return vec[ind]
		}
	}
}
