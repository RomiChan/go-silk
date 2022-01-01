package silk

import "unsafe"

func SKP_Silk_resampler_down2_3(S []int32, out []int16, in []int16, inLen int32) {
	var (
		nSamplesIn int32
		counter    int32
		res_Q6     int32
		buf        [505]int32
		buf_ptr    []int32
	)
	memcpy(unsafe.Pointer(&buf[0]), unsafe.Pointer(&S[0]), size_t(unsafe.Sizeof(int32(0))*25))
	for {
		if inLen < RESAMPLER_MAX_BATCH_SIZE_IN {
			nSamplesIn = inLen
		} else {
			nSamplesIn = RESAMPLER_MAX_BATCH_SIZE_IN
		}
		SKP_Silk_resampler_private_AR2(([]int32)(&S[25]), ([]int32)(&buf[25]), in, SKP_Silk_Resampler_2_3_COEFS_LQ[:], nSamplesIn)
		buf_ptr = ([]int32)(buf[:])
		counter = nSamplesIn
		for counter > 2 {
			res_Q6 = SKP_SMULWB(buf_ptr[0], int32(SKP_Silk_Resampler_2_3_COEFS_LQ[2]))
			res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[1], int32(SKP_Silk_Resampler_2_3_COEFS_LQ[3]))
			res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[2], int32(SKP_Silk_Resampler_2_3_COEFS_LQ[5]))
			res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[3], int32(SKP_Silk_Resampler_2_3_COEFS_LQ[4]))
			func() []int16 {
				p := &out[0]
				x := *p
				*p++
				return x
			}()[0] = SKP_SAT16(SKP_RSHIFT_ROUND(res_Q6, 6))
			res_Q6 = SKP_SMULWB(buf_ptr[1], int32(SKP_Silk_Resampler_2_3_COEFS_LQ[4]))
			res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[2], int32(SKP_Silk_Resampler_2_3_COEFS_LQ[5]))
			res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[3], int32(SKP_Silk_Resampler_2_3_COEFS_LQ[3]))
			res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[4], int32(SKP_Silk_Resampler_2_3_COEFS_LQ[2]))
			func() []int16 {
				p := &out[0]
				x := *p
				*p++
				return x
			}()[0] = SKP_SAT16(SKP_RSHIFT_ROUND(res_Q6, 6))
			buf_ptr += 3
			counter -= 3
		}
		in += ([]int16)(nSamplesIn)
		inLen -= nSamplesIn
		if inLen > 0 {
			memcpy(unsafe.Pointer(&buf[0]), unsafe.Pointer(&buf[nSamplesIn]), size_t(unsafe.Sizeof(int32(0))*25))
		} else {
			break
		}
	}
	memcpy(unsafe.Pointer(&S[0]), unsafe.Pointer(&buf[nSamplesIn]), size_t(unsafe.Sizeof(int32(0))*25))
}
