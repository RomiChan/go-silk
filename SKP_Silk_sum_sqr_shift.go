package silk

import "unsafe"

func SKP_Silk_sum_sqr_shift(energy *int32, shift *int32, x *int16, len_ int32) {
	var (
		i       int32
		shft    int32
		in32    int32
		nrg_tmp int32
		nrg     int32
	)
	if int32(int64(uintptr(unsafe.Pointer(x)))&2) != 0 {
		nrg = SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*0))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*0))))
		i = 1
	} else {
		nrg = 0
		i = 0
	}
	shft = 0
	len_--
	for i < len_ {
		in32 = *((*int32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(i))))
		nrg = int32(uint32(nrg) + uint32(SKP_SMULBB(in32, in32)))
		nrg = int32(uint32(nrg) + uint32(SKP_SMULTT(in32, in32)))
		i += 2
		if nrg < 0 {
			nrg = int32((uint32(nrg)) >> 2)
			shft = 2
			break
		}
	}
	for ; i < len_; i += 2 {
		in32 = *((*int32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(i))))
		nrg_tmp = SKP_SMULBB(in32, in32)
		nrg_tmp = int32(uint32(nrg_tmp) + uint32(SKP_SMULTT(in32, in32)))
		nrg = int32(int64(nrg) + (int64(uint32(nrg_tmp)) >> int64(shft)))
		if nrg < 0 {
			nrg = int32((uint32(nrg)) >> 2)
			shft += 2
		}
	}
	if i == len_ {
		nrg_tmp = SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(i)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(int16(0))*uintptr(i)))))
		nrg = nrg + (nrg_tmp >> shft)
	}
	if int(nrg)&0xC0000000 != 0 {
		nrg = int32((uint32(nrg)) >> 2)
		shft += 2
	}
	*shift = shft
	*energy = nrg
}
