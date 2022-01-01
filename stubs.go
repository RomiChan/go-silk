package silk

import (
	"bytes"
	"math/bits"
	"unsafe"
)

type size_t = uintptr

func SKP_Silk_CLZ64(a int64) int32 {
	return int32(bits.LeadingZeros64(uint64(a)))
}

func SKP_Silk_CLZ32(a int32) int32 {
	return int32(bits.LeadingZeros32(uint32(a)))
}

func SKP_Silk_CLZ16(a int16) int32 {
	return int32(bits.LeadingZeros16(uint16(a)))
}

func memset(p unsafe.Pointer, ch byte, sz uintptr) unsafe.Pointer {
	b := unsafe.Slice((*byte)(p), sz)
	if ch == 0 {
		copy(b, make([]byte, len(b)))
	} else {
		copy(b, bytes.Repeat([]byte{ch}, len(b)))
	}
	return p
}

func memcpy(dst, src unsafe.Pointer, sz uintptr) unsafe.Pointer {
	if dst == nil {
		panic("nil destination")
	}
	if sz == 0 || src == nil {
		return dst
	}
	bdst := unsafe.Slice(dst, sz)
	bsrc := unsafe.Slice(src, sz)
	copy(bdst, bsrc)
	return dst
}

func memmove(dst, src unsafe.Pointer, sz uintptr) unsafe.Pointer {
	if sz == 0 {
		return dst
	}
	return memcpy(dst, src, sz)
}
