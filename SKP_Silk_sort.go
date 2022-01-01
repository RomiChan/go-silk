package silk

func SKP_Silk_insertion_sort_increasing(a []int32, index []int32, L int32, K int32) {
	var (
		value int32
		i     int32
		j     int32
	)
	for i = 0; i < K; i++ {
		index[i] = i
	}
	for i = 1; i < K; i++ {
		value = a[i]
		for j = i - 1; j >= 0 && value < a[j]; j-- {
			a[j+1] = a[j]
			index[j+1] = index[j]
		}
		a[j+1] = value
		index[j+1] = i
	}
	for i = K; i < L; i++ {
		value = a[i]
		if value < a[K-1] {
			for j = K - 2; j >= 0 && value < a[j]; j-- {
				a[j+1] = a[j]
				index[j+1] = index[j]
			}
			a[j+1] = value
			index[j+1] = i
		}
	}
}
func SKP_Silk_insertion_sort_decreasing_int16(a []int16, index []int32, L int32, K int32) {
	var (
		i     int32
		j     int32
		value int32
	)
	for i = 0; i < K; i++ {
		index[i] = i
	}
	for i = 1; i < K; i++ {
		value = int32(a[i])
		for j = i - 1; j >= 0 && int64(value) > int64(a[j]); j-- {
			a[j+1] = a[j]
			index[j+1] = index[j]
		}
		a[j+1] = int16(value)
		index[j+1] = i
	}
	for i = K; i < L; i++ {
		value = int32(a[i])
		if int64(value) > int64(a[K-1]) {
			for j = K - 2; j >= 0 && int64(value) > int64(a[j]); j-- {
				a[j+1] = a[j]
				index[j+1] = index[j]
			}
			a[j+1] = int16(value)
			index[j+1] = i
		}
	}
}
func SKP_Silk_insertion_sort_increasing_all_values(a []int32, L int32) {
	var (
		value int32
		i     int32
		j     int32
	)
	for i = 1; i < L; i++ {
		value = a[i]
		for j = i - 1; j >= 0 && value < a[j]; j-- {
			a[j+1] = a[j]
		}
		a[j+1] = value
	}
}
