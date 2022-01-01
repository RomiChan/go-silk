package silk

import "unsafe"

func SKP_Silk_insertion_sort_increasing(a *int32, index *int32, L int32, K int32) {
	var (
		value int32
		i     int32
		j     int32
	)
	for i = 0; int64(i) < int64(K); i++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(i))) = i
	}
	for i = 1; int64(i) < int64(K); i++ {
		value = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(i)))
		for j = int32(int64(i) - 1); int64(j) >= 0 && int64(value) < int64(*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))); j-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))
			*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(j)))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = value
		*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = i
	}
	for i = K; int64(i) < int64(L); i++ {
		value = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(i)))
		if int64(value) < int64(*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(int64(K)-1)))) {
			for j = int32(int64(K) - 2); int64(j) >= 0 && int64(value) < int64(*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))); j-- {
				*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))
				*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(j)))
			}
			*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = value
			*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = i
		}
	}
}
func SKP_Silk_insertion_sort_decreasing_int16(a *int16, index *int32, L int32, K int32) {
	var (
		i     int32
		j     int32
		value int32
	)
	for i = 0; int64(i) < int64(K); i++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(i))) = i
	}
	for i = 1; int64(i) < int64(K); i++ {
		value = int32(*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(i))))
		for j = int32(int64(i) - 1); int64(j) >= 0 && int64(value) > int64(*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(j)))); j-- {
			*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(int64(j)+1))) = *(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(j)))
			*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(j)))
		}
		*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(int64(j)+1))) = int16(value)
		*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = i
	}
	for i = K; int64(i) < int64(L); i++ {
		value = int32(*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(i))))
		if int64(value) > int64(*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(int64(K)-1)))) {
			for j = int32(int64(K) - 2); int64(j) >= 0 && int64(value) > int64(*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(j)))); j-- {
				*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(int64(j)+1))) = *(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(j)))
				*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(j)))
			}
			*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(int64(j)+1))) = int16(value)
			*(*int32)(unsafe.Add(unsafe.Pointer(index), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = i
		}
	}
}
func SKP_Silk_insertion_sort_increasing_all_values(a *int32, L int32) {
	var (
		value int32
		i     int32
		j     int32
	)
	for i = 1; int64(i) < int64(L); i++ {
		value = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(i)))
		for j = int32(int64(i) - 1); int64(j) >= 0 && int64(value) < int64(*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))); j-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(int64(j)+1))) = value
	}
}
