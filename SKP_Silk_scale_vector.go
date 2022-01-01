package silk

// SKP_Silk_scale_copy_vector16
// see SKP_Silk_scale_copy_vector16.c
func SKP_Silk_scale_copy_vector16(data_out []int16, data_in []int16, gain_Q16 int32, dataSize int32) {
	for i := int32(0); i < dataSize; i++ {
		tmp32 := SKP_SMULWB(gain_Q16, int32(data_in[i]))
		data_out[i] = int16(tmp32)
	}
}

func SKP_Silk_scale_vector32_Q26_lshift_18(data1 []int32, gain_Q26 int32, dataSize int32) {
	for i := int32(0); i < dataSize; i++ {
		data1[i] = int32((int64(data1[i]) * int64(gain_Q26)) >> 8)
	}
}
