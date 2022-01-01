package silk

import (
	"math"
	"unsafe"
)

func SKP_Silk_resampler_private_IIR_FIR_INTERPOL(out *int16, buf *int16, max_index_Q16 int32, index_increment_Q16 int32) *int16 {
	var (
		index_Q16   int32
		res_Q15     int32
		buf_ptr     *int16
		table_index int32
	)
	for index_Q16 = 0; index_Q16 < max_index_Q16; index_Q16 += index_increment_Q16 {
		table_index = SKP_SMULWB(index_Q16&math.MaxUint16, 144)
		buf_ptr = (*int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int16(0))*uintptr(index_Q16>>16)))
		res_Q15 = SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*0))), int32(SKP_Silk_resampler_frac_FIR_144[table_index][0]))
		res_Q15 = SKP_SMLABB(res_Q15, int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*1))), int32(SKP_Silk_resampler_frac_FIR_144[table_index][1]))
		res_Q15 = SKP_SMLABB(res_Q15, int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*2))), int32(SKP_Silk_resampler_frac_FIR_144[table_index][2]))
		res_Q15 = SKP_SMLABB(res_Q15, int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*3))), int32(SKP_Silk_resampler_frac_FIR_144[143-table_index][2]))
		res_Q15 = SKP_SMLABB(res_Q15, int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*4))), int32(SKP_Silk_resampler_frac_FIR_144[143-table_index][1]))
		res_Q15 = SKP_SMLABB(res_Q15, int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*5))), int32(SKP_Silk_resampler_frac_FIR_144[143-table_index][0]))
		*func() *int16 {
			p := &out
			x := *p
			*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
			return x
		}() = SKP_SAT16(SKP_RSHIFT_ROUND(res_Q15, 15))
	}
	return out
}
func SKP_Silk_resampler_private_IIR_FIR(SS *SKP_Silk_resampler_state_struct, out []int16, in []int16, inLen int) {
	var (
		S                   *SKP_Silk_resampler_state_struct = SS
		nSamplesIn          int32
		max_index_Q16       int32
		index_increment_Q16 int32
		buf                 [966]int16
	)
	memcpy(unsafe.Pointer(&buf[0]), unsafe.Pointer(&S.SFIR[0]), size_t(RESAMPLER_ORDER_FIR_144*unsafe.Sizeof(int32(0))))
	index_increment_Q16 = S.InvRatio_Q16
	for {
		if inLen < S.BatchSize {
			nSamplesIn = inLen
		} else {
			nSamplesIn = S.BatchSize
		}
		if S.Input2x == 1 {
			S.Up2_function(&S.SIIR[0], &buf[RESAMPLER_ORDER_FIR_144], &in[0], nSamplesIn)
		} else {
			SKP_Silk_resampler_private_ARMA4(S.SIIR[:], ([]int16)(&buf[RESAMPLER_ORDER_FIR_144]), in, ([]int16)(S.Coefs), nSamplesIn)
		}
		max_index_Q16 = nSamplesIn << (S.Input2x + 16)
		out = ([]int16)(SKP_Silk_resampler_private_IIR_FIR_INTERPOL(&out[0], &buf[0], max_index_Q16, index_increment_Q16))
		in += ([]int16)(nSamplesIn)
		inLen -= nSamplesIn
		if inLen > 0 {
			memcpy(unsafe.Pointer(&buf[0]), unsafe.Pointer(&buf[nSamplesIn<<S.Input2x]), size_t(RESAMPLER_ORDER_FIR_144*unsafe.Sizeof(int32(0))))
		} else {
			break
		}
	}
	memcpy(unsafe.Pointer(&S.SFIR[0]), unsafe.Pointer(&buf[nSamplesIn<<S.Input2x]), size_t(RESAMPLER_ORDER_FIR_144*unsafe.Sizeof(int32(0))))
}
