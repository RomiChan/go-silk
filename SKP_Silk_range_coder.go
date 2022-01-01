package silk

import (
	"math"
	"unsafe"
)

func SKP_Silk_range_encoder(psRC *SKP_Silk_range_coder_state, data int32, prob []uint16) {
	var (
		low_Q16   uint32
		high_Q16  uint32
		base_tmp  uint32
		range_Q32 uint32
		base_Q32  uint32 = psRC.Base_Q32
		range_Q16 uint32 = psRC.Range_Q16
		bufferIx  int32  = psRC.BufferIx
		buffer    *uint8 = &psRC.Buffer[0]
	)
	if int64(psRC.Error) != 0 {
		return
	}
	low_Q16 = uint32(prob[data])
	high_Q16 = uint32(prob[int64(data)+1])
	base_tmp = base_Q32
	base_Q32 += uint32(int32(int64(range_Q16) * int64(low_Q16)))
	range_Q32 = uint32(int32(int64(range_Q16) * (int64(high_Q16) - int64(low_Q16))))
	if int64(base_Q32) < int64(base_tmp) {
		var bufferIx_tmp int32 = bufferIx
		for int64(func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(buffer), func() int32 {
				p := &bufferIx_tmp
				*p--
				return *p
			}()))
			*p++
			return *p
		}()) == 0 {
		}
	}
	if int64(range_Q32)&0xFF000000 != 0 {
		range_Q16 = uint32(int32(int64(range_Q32) >> 16))
	} else {
		if int64(range_Q32)&0xFFFF0000 != 0 {
			range_Q16 = uint32(int32(int64(range_Q32) >> 8))
		} else {
			range_Q16 = range_Q32
			if int64(bufferIx) >= int64(psRC.BufferLength) {
				psRC.Error = -1
				return
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(buffer), func() int32 {
				p := &bufferIx
				x := *p
				*p++
				return x
			}())) = uint8(int8(int64(base_Q32) >> 24))
			base_Q32 = uint32(int32(int64(base_Q32) << 8))
		}
		if int64(bufferIx) >= int64(psRC.BufferLength) {
			psRC.Error = -1
			return
		}
		*(*uint8)(unsafe.Add(unsafe.Pointer(buffer), func() int32 {
			p := &bufferIx
			x := *p
			*p++
			return x
		}())) = uint8(int8(int64(base_Q32) >> 24))
		base_Q32 = uint32(int32(int64(base_Q32) << 8))
	}
	psRC.Base_Q32 = base_Q32
	psRC.Range_Q16 = range_Q16
	psRC.BufferIx = bufferIx
}
func SKP_Silk_range_encoder_multi(psRC *SKP_Silk_range_coder_state, data []int32, prob []*uint16, nSymbols int32) {
	var k int32
	for k = 0; int64(k) < int64(nSymbols); k++ {
		SKP_Silk_range_encoder(psRC, data[k], ([]uint16)(prob[k]))
	}
}
func SKP_Silk_range_decoder(data *int32, psRC *SKP_Silk_range_coder_state, prob []uint16, probIx int32) {
	var (
		low_Q16   uint32
		high_Q16  uint32
		base_tmp  uint32
		range_Q32 uint32
		base_Q32  uint32 = psRC.Base_Q32
		range_Q16 uint32 = psRC.Range_Q16
		bufferIx  int32  = psRC.BufferIx
		buffer    *uint8 = &psRC.Buffer[4]
	)
	if int64(psRC.Error) != 0 {
		*data = 0
		return
	}
	high_Q16 = uint32(prob[probIx])
	base_tmp = uint32(int32(int64(range_Q16) * int64(high_Q16)))
	if int64(base_tmp) > int64(base_Q32) {
		for {
			low_Q16 = uint32(prob[func() int32 {
				p := &probIx
				*p--
				return *p
			}()])
			base_tmp = uint32(int32(int64(range_Q16) * int64(low_Q16)))
			if int64(base_tmp) <= int64(base_Q32) {
				break
			}
			high_Q16 = low_Q16
			if int64(high_Q16) == 0 {
				psRC.Error = -2
				*data = 0
				return
			}
		}
	} else {
		for {
			low_Q16 = high_Q16
			high_Q16 = uint32(prob[func() int32 {
				p := &probIx
				*p++
				return *p
			}()])
			base_tmp = uint32(int32(int64(range_Q16) * int64(high_Q16)))
			if int64(base_tmp) > int64(base_Q32) {
				probIx--
				break
			}
			if int64(high_Q16) == math.MaxUint16 {
				psRC.Error = -2
				*data = 0
				return
			}
		}
	}
	*data = probIx
	base_Q32 -= uint32(int32(int64(range_Q16) * int64(low_Q16)))
	range_Q32 = uint32(int32(int64(range_Q16) * (int64(high_Q16) - int64(low_Q16))))
	if int64(range_Q32)&0xFF000000 != 0 {
		range_Q16 = uint32(int32(int64(range_Q32) >> 16))
	} else {
		if int64(range_Q32)&0xFFFF0000 != 0 {
			range_Q16 = uint32(int32(int64(range_Q32) >> 8))
			if (int64(base_Q32) >> 24) != 0 {
				psRC.Error = -3
				*data = 0
				return
			}
		} else {
			range_Q16 = range_Q32
			if (int64(base_Q32) >> 16) != 0 {
				psRC.Error = -3
				*data = 0
				return
			}
			base_Q32 = uint32(int32(int64(base_Q32) << 8))
			if int64(bufferIx) < int64(psRC.BufferLength) {
				base_Q32 |= uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(buffer), func() int32 {
					p := &bufferIx
					x := *p
					*p++
					return x
				}())))
			}
		}
		base_Q32 = uint32(int32(int64(base_Q32) << 8))
		if int64(bufferIx) < int64(psRC.BufferLength) {
			base_Q32 |= uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(buffer), func() int32 {
				p := &bufferIx
				x := *p
				*p++
				return x
			}())))
		}
	}
	if int64(range_Q16) == 0 {
		psRC.Error = -4
		*data = 0
		return
	}
	psRC.Base_Q32 = base_Q32
	psRC.Range_Q16 = range_Q16
	psRC.BufferIx = bufferIx
}
func SKP_Silk_range_decoder_multi(data []int32, psRC *SKP_Silk_range_coder_state, prob []*uint16, probStartIx []int32, nSymbols int32) {
	var k int32
	for k = 0; int64(k) < int64(nSymbols); k++ {
		SKP_Silk_range_decoder(&data[k], psRC, ([]uint16)(prob[k]), probStartIx[k])
	}
}
func SKP_Silk_range_enc_init(psRC *SKP_Silk_range_coder_state) {
	psRC.BufferLength = MAX_ARITHM_BYTES
	psRC.Range_Q16 = math.MaxUint16
	psRC.BufferIx = 0
	psRC.Base_Q32 = 0
	psRC.Error = 0
}
func SKP_Silk_range_dec_init(psRC *SKP_Silk_range_coder_state, buffer []uint8, bufferLength int32) {
	if int64(bufferLength) > MAX_ARITHM_BYTES || int64(bufferLength) < 0 {
		psRC.Error = -8
		return
	}
	memcpy(unsafe.Pointer(&psRC.Buffer[0]), unsafe.Pointer(&buffer[0]), size_t(uintptr(bufferLength)*unsafe.Sizeof(uint8(0))))
	psRC.BufferLength = bufferLength
	psRC.BufferIx = 0
	psRC.Base_Q32 = uint32(int32((int64(uint32(buffer[0])) << 24) | int64(uint32(buffer[1]))<<16 | int64(uint32(buffer[2]))<<8 | int64(buffer[3])))
	psRC.Range_Q16 = math.MaxUint16
	psRC.Error = 0
}
func SKP_Silk_range_coder_get_length(psRC *SKP_Silk_range_coder_state, nBytes *int32) int32 {
	var nBits int32
	nBits = int32((int64(psRC.BufferIx) << 3) + int64(SKP_Silk_CLZ32(int32(int64(psRC.Range_Q16)-1))) - 14)
	*nBytes = int32((int64(nBits) + 7) >> 3)
	return nBits
}
func SKP_Silk_range_enc_wrap_up(psRC *SKP_Silk_range_coder_state) {
	var (
		bufferIx_tmp   int32
		bits_to_store  int32
		bits_in_stream int32
		nBytes         int32
		mask           int32
		base_Q24       uint32
	)
	base_Q24 = uint32(int32(int64(psRC.Base_Q32) >> 8))
	bits_in_stream = SKP_Silk_range_coder_get_length(psRC, &nBytes)
	bits_to_store = int32(int64(bits_in_stream) - (int64(psRC.BufferIx) << 3))
	base_Q24 += uint32(int32(0x800000 >> (int64(bits_to_store) - 1)))
	base_Q24 &= uint32(int32(math.MaxUint32 << (24 - int64(bits_to_store))))
	if int64(base_Q24)&0x1000000 != 0 {
		bufferIx_tmp = psRC.BufferIx
		for int64(func() uint8 {
			p := &(psRC.Buffer[func() int32 {
				p := &bufferIx_tmp
				*p--
				return *p
			}()])
			*p++
			return *p
		}()) == 0 {
		}
	}
	if int64(psRC.BufferIx) < int64(psRC.BufferLength) {
		psRC.Buffer[func() int32 {
			p := &psRC.BufferIx
			x := *p
			*p++
			return x
		}()] = uint8(int8(int64(base_Q24) >> 16))
		if int64(bits_to_store) > 8 {
			if int64(psRC.BufferIx) < int64(psRC.BufferLength) {
				psRC.Buffer[func() int32 {
					p := &psRC.BufferIx
					x := *p
					*p++
					return x
				}()] = uint8(int8(int64(base_Q24) >> 8))
			}
		}
	}
	if int64(bits_in_stream)&7 != 0 {
		mask = int32(math.MaxUint8 >> (int64(bits_in_stream) & 7))
		if int64(nBytes)-1 < int64(psRC.BufferLength) {
			psRC.Buffer[int64(nBytes)-1] |= uint8(int8(mask))
		}
	}
}
func SKP_Silk_range_coder_check_after_decoding(psRC *SKP_Silk_range_coder_state) {
	var (
		bits_in_stream int32
		nBytes         int32
		mask           int32
	)
	bits_in_stream = SKP_Silk_range_coder_get_length(psRC, &nBytes)
	if int64(nBytes)-1 >= int64(psRC.BufferLength) {
		psRC.Error = -5
		return
	}
	if int64(bits_in_stream)&7 != 0 {
		mask = int32(math.MaxUint8 >> (int64(bits_in_stream) & 7))
		if (int64(psRC.Buffer[int64(nBytes)-1]) & int64(mask)) != int64(mask) {
			psRC.Error = -5
			return
		}
	}
}
