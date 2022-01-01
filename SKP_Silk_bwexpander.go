package silk

func SKP_Silk_bwexpander(ar []int16, d int32, chirp_Q16 int32) {
	chirp_minus_one_Q16 := chirp_Q16 - 0x10000
	for i := int32(0); i < d-1; i++ {
		ar[i] = int16(SKP_RSHIFT_ROUND(int32(int64(chirp_Q16)*int64(ar[i])), 16))
		chirp_Q16 += SKP_RSHIFT_ROUND(chirp_Q16*chirp_minus_one_Q16, 16)
	}
	ar[d-1] = int16(SKP_RSHIFT_ROUND(int32(int64(chirp_Q16)*int64(ar[d-1])), 16))
}

func SKP_Silk_bwexpander_32(ar []int32, d int32, chirp_Q16 int32) {
	tmp_chirp_Q16 := chirp_Q16
	for i := int32(0); i < d-1; i++ {
		ar[i] = SKP_SMULWW(ar[i], tmp_chirp_Q16)
		tmp_chirp_Q16 = SKP_SMULWW(chirp_Q16, tmp_chirp_Q16)
	}
	ar[d-1] = SKP_SMULWW(ar[d-1], tmp_chirp_Q16)
}
