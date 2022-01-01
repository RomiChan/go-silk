package silk

import "unsafe"

func SKP_Silk_resampler_private_up2_HQ(S []int32, out []int16, in []int16, len_ int32) {
	var (
		k       int32
		in32    int32
		out32_1 int32
		out32_2 int32
		Y       int32
		X       int32
	)
	for k = 0; k < len_; k++ {
		in32 = (int32(in[k])) << 10
		Y = in32 - (S[0])
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_up2_hq_0[0]))
		out32_1 = (S[0]) + X
		S[0] = in32 + X
		Y = out32_1 - (S[1])
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_up2_hq_0[1]))
		out32_2 = (S[1]) + X
		S[1] = out32_1 + X
		out32_2 = SKP_SMLAWB(out32_2, S[5], int32(SKP_Silk_resampler_up2_hq_notch[2]))
		out32_2 = SKP_SMLAWB(out32_2, S[4], int32(SKP_Silk_resampler_up2_hq_notch[1]))
		out32_1 = SKP_SMLAWB(out32_2, S[4], int32(SKP_Silk_resampler_up2_hq_notch[0]))
		S[5] = out32_2 - (S[5])
		out[k*2] = SKP_SAT16(SKP_SMLAWB(256, out32_1, int32(SKP_Silk_resampler_up2_hq_notch[3])) >> 9)
		Y = in32 - (S[2])
		X = SKP_SMULWB(Y, int32(SKP_Silk_resampler_up2_hq_1[0]))
		out32_1 = (S[2]) + X
		S[2] = in32 + X
		Y = out32_1 - (S[3])
		X = SKP_SMLAWB(Y, Y, int32(SKP_Silk_resampler_up2_hq_1[1]))
		out32_2 = (S[3]) + X
		S[3] = out32_1 + X
		out32_2 = SKP_SMLAWB(out32_2, S[4], int32(SKP_Silk_resampler_up2_hq_notch[2]))
		out32_2 = SKP_SMLAWB(out32_2, S[5], int32(SKP_Silk_resampler_up2_hq_notch[1]))
		out32_1 = SKP_SMLAWB(out32_2, S[5], int32(SKP_Silk_resampler_up2_hq_notch[0]))
		S[4] = out32_2 - (S[4])
		out[k*2+1] = SKP_SAT16(SKP_SMLAWB(256, out32_1, int32(SKP_Silk_resampler_up2_hq_notch[3])) >> 9)
	}
}
func SKP_Silk_resampler_private_up2_HQ_wrapper(SS unsafe.Pointer, out []int16, in []int16, len_ int32) {
	var S *SKP_Silk_resampler_state_struct = (*SKP_Silk_resampler_state_struct)(SS)
	SKP_Silk_resampler_private_up2_HQ(S.SIIR[:], out, in, len_)
}
