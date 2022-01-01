package silk

import "unsafe"

func SKP_Silk_resampler_down3(S *int32, out *int16, in *int16, inLen int32) {
	var (
		nSamplesIn int32
		counter    int32
		res_Q6     int32
		buf        [505]int32
		buf_ptr    []int32
	)
	memcpy(unsafe.Pointer(&buf[0]), unsafe.Pointer(S), size_t(unsafe.Sizeof(int32(0))*25))
	for {
		if inLen < RESAMPLER_MAX_BATCH_SIZE_IN {
			nSamplesIn = inLen
		} else {
			nSamplesIn = RESAMPLER_MAX_BATCH_SIZE_IN
		}
		SKP_Silk_resampler_private_AR2(([]int32)((*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*25))), ([]int32)(&buf[25]), ([]int16)(in), SKP_Silk_Resampler_1_3_COEFS_LQ[:], nSamplesIn)
		buf_ptr = ([]int32)(buf[:])
		counter = nSamplesIn
		for counter > 2 {
			res_Q6 = SKP_SMULWB((buf_ptr[0])+(buf_ptr[5]), int32(SKP_Silk_Resampler_1_3_COEFS_LQ[2]))
			res_Q6 = SKP_SMLAWB(res_Q6, (buf_ptr[1])+(buf_ptr[4]), int32(SKP_Silk_Resampler_1_3_COEFS_LQ[3]))
			res_Q6 = SKP_SMLAWB(res_Q6, (buf_ptr[2])+(buf_ptr[3]), int32(SKP_Silk_Resampler_1_3_COEFS_LQ[4]))
			*func() *int16 {
				p := &out
				x := *p
				*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
				return x
			}() = SKP_SAT16(SKP_RSHIFT_ROUND(res_Q6, 6))
			buf_ptr += 3
			counter -= 3
		}
		in = (*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(nSamplesIn)))
		inLen -= nSamplesIn
		if inLen > 0 {
			memcpy(unsafe.Pointer(&buf[0]), unsafe.Pointer(&buf[nSamplesIn]), size_t(unsafe.Sizeof(int32(0))*25))
		} else {
			break
		}
	}
	memcpy(unsafe.Pointer(S), unsafe.Pointer(&buf[nSamplesIn]), size_t(unsafe.Sizeof(int32(0))*25))
}
