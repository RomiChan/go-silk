package silk

import (
	"math"
	"unsafe"
)

func SKP_Silk_resampler_private_down_FIR_INTERPOL0(out []int16, buf2 []int32, FIR_Coefs []int16, max_index_Q16 int32, index_increment_Q16 int32) *int16 {
	var (
		index_Q16 int32
		res_Q6    int32
		buf_ptr   []int32
	)
	for index_Q16 = 0; index_Q16 < max_index_Q16; index_Q16 += index_increment_Q16 {
		buf_ptr = ([]int32)(&buf2[index_Q16>>16])
		res_Q6 = SKP_SMULWB((buf_ptr[0])+(buf_ptr[11]), int32(FIR_Coefs[0]))
		res_Q6 = SKP_SMLAWB(res_Q6, (buf_ptr[1])+(buf_ptr[10]), int32(FIR_Coefs[1]))
		res_Q6 = SKP_SMLAWB(res_Q6, (buf_ptr[2])+(buf_ptr[9]), int32(FIR_Coefs[2]))
		res_Q6 = SKP_SMLAWB(res_Q6, (buf_ptr[3])+(buf_ptr[8]), int32(FIR_Coefs[3]))
		res_Q6 = SKP_SMLAWB(res_Q6, (buf_ptr[4])+(buf_ptr[7]), int32(FIR_Coefs[4]))
		res_Q6 = SKP_SMLAWB(res_Q6, (buf_ptr[5])+(buf_ptr[6]), int32(FIR_Coefs[5]))
		func() []int16 {
			p := &out[0]
			x := *p
			*p++
			return x
		}()[0] = SKP_SAT16(SKP_RSHIFT_ROUND(res_Q6, 6))
	}
	return &out[0]
}
func SKP_Silk_resampler_private_down_FIR_INTERPOL1(out []int16, buf2 []int32, FIR_Coefs []int16, max_index_Q16 int32, index_increment_Q16 int32, FIR_Fracs int32) *int16 {
	var (
		index_Q16    int32
		res_Q6       int32
		buf_ptr      []int32
		interpol_ind int32
		interpol_ptr []int16
	)
	for index_Q16 = 0; index_Q16 < max_index_Q16; index_Q16 += index_increment_Q16 {
		buf_ptr = ([]int32)(&buf2[index_Q16>>16])
		interpol_ind = SKP_SMULWB(index_Q16&math.MaxUint16, FIR_Fracs)
		interpol_ptr = ([]int16)(&FIR_Coefs[RESAMPLER_DOWN_ORDER_FIR/2*interpol_ind])
		res_Q6 = SKP_SMULWB(buf_ptr[0], int32(interpol_ptr[0]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[1], int32(interpol_ptr[1]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[2], int32(interpol_ptr[2]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[3], int32(interpol_ptr[3]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[4], int32(interpol_ptr[4]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[5], int32(interpol_ptr[5]))
		interpol_ptr = ([]int16)(&FIR_Coefs[RESAMPLER_DOWN_ORDER_FIR/2*(FIR_Fracs-1-interpol_ind)])
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[11], int32(interpol_ptr[0]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[10], int32(interpol_ptr[1]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[9], int32(interpol_ptr[2]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[8], int32(interpol_ptr[3]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[7], int32(interpol_ptr[4]))
		res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[6], int32(interpol_ptr[5]))
		func() []int16 {
			p := &out[0]
			x := *p
			*p++
			return x
		}()[0] = SKP_SAT16(SKP_RSHIFT_ROUND(res_Q6, 6))
	}
	return &out[0]
}
func SKP_Silk_resampler_private_down_FIR(SS unsafe.Pointer, out []int16, in []int16, inLen int32) {
	var (
		S                   *SKP_Silk_resampler_state_struct = (*SKP_Silk_resampler_state_struct)(SS)
		nSamplesIn          int32
		max_index_Q16       int32
		index_increment_Q16 int32
		buf1                [240]int16
		buf2                [492]int32
		FIR_Coefs           *int16
	)
	memcpy(unsafe.Pointer(&buf2[0]), unsafe.Pointer(&S.SFIR[0]), size_t(RESAMPLER_DOWN_ORDER_FIR*unsafe.Sizeof(int32(0))))
	FIR_Coefs = (*int16)(unsafe.Add(unsafe.Pointer(S.Coefs), unsafe.Sizeof(int16(0))*2))
	index_increment_Q16 = S.InvRatio_Q16
	for {
		if inLen < S.BatchSize {
			nSamplesIn = inLen
		} else {
			nSamplesIn = S.BatchSize
		}
		if S.Input2x == 1 {
			SKP_Silk_resampler_down2(S.SDown2[:], buf1[:], in, nSamplesIn)
			nSamplesIn = nSamplesIn >> 1
			SKP_Silk_resampler_private_AR2(S.SIIR[:], ([]int32)(&buf2[RESAMPLER_DOWN_ORDER_FIR]), buf1[:], ([]int16)(S.Coefs), nSamplesIn)
		} else {
			SKP_Silk_resampler_private_AR2(S.SIIR[:], ([]int32)(&buf2[RESAMPLER_DOWN_ORDER_FIR]), in, ([]int16)(S.Coefs), nSamplesIn)
		}
		max_index_Q16 = nSamplesIn << 16
		if S.FIR_Fracs == 1 {
			out = ([]int16)(SKP_Silk_resampler_private_down_FIR_INTERPOL0(out, buf2[:], ([]int16)(FIR_Coefs), max_index_Q16, index_increment_Q16))
		} else {
			out = ([]int16)(SKP_Silk_resampler_private_down_FIR_INTERPOL1(out, buf2[:], ([]int16)(FIR_Coefs), max_index_Q16, index_increment_Q16, S.FIR_Fracs))
		}
		in += ([]int16)(nSamplesIn << S.Input2x)
		inLen -= nSamplesIn << S.Input2x
		if inLen > S.Input2x {
			memcpy(unsafe.Pointer(&buf2[0]), unsafe.Pointer(&buf2[nSamplesIn]), size_t(RESAMPLER_DOWN_ORDER_FIR*unsafe.Sizeof(int32(0))))
		} else {
			break
		}
	}
	memcpy(unsafe.Pointer(&S.SFIR[0]), unsafe.Pointer(&buf2[nSamplesIn]), size_t(RESAMPLER_DOWN_ORDER_FIR*unsafe.Sizeof(int32(0))))
}
