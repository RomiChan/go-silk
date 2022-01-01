package silk

import (
	"math"
	"unsafe"
)

func SKP_Silk_resampler_private_down_FIR_INTERPOL0(out *int16, buf2 *int32, FIR_Coefs *int16, max_index_Q16 int32, index_increment_Q16 int32) *int16 {
	var (
		index_Q16 int32
		res_Q6    int32
		buf_ptr   *int32
	)
	for index_Q16 = 0; int64(index_Q16) < int64(max_index_Q16); index_Q16 += index_increment_Q16 {
		buf_ptr = (*int32)(unsafe.Add(unsafe.Pointer(buf2), unsafe.Sizeof(int32(0))*uintptr(int64(index_Q16)>>16)))
		res_Q6 = SKP_SMULWB(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*0)))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*11)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*0))))
		res_Q6 = SKP_SMLAWB(res_Q6, int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*1)))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*10)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*1))))
		res_Q6 = SKP_SMLAWB(res_Q6, int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*2)))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*9)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*2))))
		res_Q6 = SKP_SMLAWB(res_Q6, int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*3)))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*8)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*3))))
		res_Q6 = SKP_SMLAWB(res_Q6, int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*4)))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*7)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*4))))
		res_Q6 = SKP_SMLAWB(res_Q6, int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*5)))+int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*6)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*5))))
		*func() *int16 {
			p := &out
			x := *p
			*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
			return x
		}() = SKP_SAT16(int16(SKP_RSHIFT_ROUND(res_Q6, 6)))
	}
	return out
}
func SKP_Silk_resampler_private_down_FIR_INTERPOL1(out *int16, buf2 *int32, FIR_Coefs *int16, max_index_Q16 int32, index_increment_Q16 int32, FIR_Fracs int32) *int16 {
	var (
		index_Q16    int32
		res_Q6       int32
		buf_ptr      *int32
		interpol_ind int32
		interpol_ptr *int16
	)
	for index_Q16 = 0; int64(index_Q16) < int64(max_index_Q16); index_Q16 += index_increment_Q16 {
		buf_ptr = (*int32)(unsafe.Add(unsafe.Pointer(buf2), unsafe.Sizeof(int32(0))*uintptr(int64(index_Q16)>>16)))
		interpol_ind = SKP_SMULWB(int32(int64(index_Q16)&math.MaxUint16), FIR_Fracs)
		interpol_ptr = (*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*uintptr(RESAMPLER_DOWN_ORDER_FIR/2*int64(interpol_ind))))
		res_Q6 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*0)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*0))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*1)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*1))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*2)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*2))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*3)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*3))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*4)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*4))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*5)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*5))))
		interpol_ptr = (*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*uintptr(RESAMPLER_DOWN_ORDER_FIR/2*(int64(FIR_Fracs)-1-int64(interpol_ind)))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*11)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*0))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*10)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*1))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*9)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*2))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*8)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*3))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*7)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*4))))
		res_Q6 = SKP_SMLAWB(res_Q6, *(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*6)), int32(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*5))))
		*func() *int16 {
			p := &out
			x := *p
			*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
			return x
		}() = SKP_SAT16(int16(SKP_RSHIFT_ROUND(res_Q6, 6)))
	}
	return out
}
func SKP_Silk_resampler_private_down_FIR(SS *SKP_Silk_resampler_state_struct, out []int16, in []int16, inLen int) {
	var (
		S                   *SKP_Silk_resampler_state_struct = SS
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
		if int64(inLen) < int64(S.BatchSize) {
			nSamplesIn = inLen
		} else {
			nSamplesIn = S.BatchSize
		}
		if int64(S.Input2x) == 1 {
			SKP_Silk_resampler_down2(&S.SDown2[0], &buf1[0], &in[0], nSamplesIn)
			nSamplesIn = nSamplesIn >> 1
			SKP_Silk_resampler_private_AR2(S.SIIR[:], ([]int32)(&buf2[RESAMPLER_DOWN_ORDER_FIR]), buf1[:], ([]int16)(S.Coefs), nSamplesIn)
		} else {
			SKP_Silk_resampler_private_AR2(S.SIIR[:], ([]int32)(&buf2[RESAMPLER_DOWN_ORDER_FIR]), in, ([]int16)(S.Coefs), nSamplesIn)
		}
		max_index_Q16 = int32(int64(nSamplesIn) << 16)
		if int64(S.FIR_Fracs) == 1 {
			out = ([]int16)(SKP_Silk_resampler_private_down_FIR_INTERPOL0(&out[0], &buf2[0], FIR_Coefs, max_index_Q16, index_increment_Q16))
		} else {
			out = ([]int16)(SKP_Silk_resampler_private_down_FIR_INTERPOL1(&out[0], &buf2[0], FIR_Coefs, max_index_Q16, index_increment_Q16, S.FIR_Fracs))
		}
		in += ([]int16)(int64(nSamplesIn) << int64(S.Input2x))
		inLen -= int32(int64(nSamplesIn) << int64(S.Input2x))
		if int64(inLen) > int64(S.Input2x) {
			memcpy(unsafe.Pointer(&buf2[0]), unsafe.Pointer(&buf2[nSamplesIn]), size_t(RESAMPLER_DOWN_ORDER_FIR*unsafe.Sizeof(int32(0))))
		} else {
			break
		}
	}
	memcpy(unsafe.Pointer(&S.SFIR[0]), unsafe.Pointer(&buf2[nSamplesIn]), size_t(RESAMPLER_DOWN_ORDER_FIR*unsafe.Sizeof(int32(0))))
}
