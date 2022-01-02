package silk

func SKP_Silk_inner_prod_aligned(inVec1 []int16, inVec2 []int16, _len int32) int32 {
	var sum int32
	for i := int32(0); i < _len; i++ {
		sum = SKP_SMLABB(sum, int32(inVec1[i]), int32(inVec2[i]))
	}
	return sum
}

func SKP_Silk_inner_prod16_aligned_64(inVec1 []int16, inVec2 []int16, _len int32) int64 {
	var sum int64
	for i := int32(0); i < _len; i++ {
		sum = sum + int64(inVec1[i])*int64(inVec2[i])
	}
	return sum
}
