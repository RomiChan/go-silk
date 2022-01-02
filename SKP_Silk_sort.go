package silk

/* Insertion sort (fast for already almost sorted arrays):   */
/* Best case:  O(n)   for an already sorted array            */
/* Worst case: O(n^2) for an inversely sorted array          */

func SKP_Silk_insertion_sort_increasing(
	a []int32, /* I/O:  Unsorted / Sorted vector                     */
	index []int32, /* O:    Index vector for the sorted elements         */
	L int32, /* I:    Vector length                                */
	K int32, /* I:    Number of correctly sorted output positions  */
) {
	var i, j, value int32

	/* Write start indices in index vector */
	for i = 0; i < K; i++ {
		index[i] = i
	}

	/* Sort vector elements by value, increasing order */
	for i = 1; i < K; i++ {
		value = a[i]
		for j = i - 1; j >= 0 && value < a[j]; j-- {
			a[j+1] = a[j]         /* Shift value */
			index[j+1] = index[j] /* Shift index */
		}
		a[j+1] = value /* Write value */
		index[j+1] = i /* Write index */
	}

	/* If less than L values are asked for, check the remaining values, */
	/* but only spend CPU to ensure that the K first values are correct */
	for i = K; i < L; i++ {
		value = a[i]
		if value < a[K-1] {
			for j = K - 2; j >= 0 && value < a[j]; j-- {
				a[j+1] = a[j]         /* Shift value */
				index[j+1] = index[j] /* Shift index */
			}
			a[j+1] = value /* Write value */
			index[j+1] = i /* Write index */
		}
	}
}

func SKP_Silk_insertion_sort_decreasing_int16(
	a []int16, // I/O:  Unsorted / Sorted vector                     */
	index []int32, // O:    Index vector for the sorted elements         */
	L int32, // I:    Vector length                                */
	K int32, // I:    Number of correctly sorted output positions  */
) {
	var i, j, value int32

	/* Write start indices in index vector */
	for i = 0; i < K; i++ {
		index[i] = i
	}

	/* Sort vector elements by value, decreasing order */
	for i = 1; i < K; i++ {
		value = int32(a[i])
		for j = i - 1; j >= 0 && int64(value) > int64(a[j]); j-- {
			a[j+1] = a[j]         /* Shift value */
			index[j+1] = index[j] /* Shift index */
		}
		a[j+1] = int16(value) /* Write value */
		index[j+1] = i        /* Write index */
	}

	/* If less than L values are asked for, check the remaining values, */
	/* but only spend CPU to ensure that the K first values are correct */
	for i = K; i < L; i++ {
		value = int32(a[i])
		if int64(value) > int64(a[K-1]) {
			for j = K - 2; j >= 0 && int64(value) > int64(a[j]); j-- {
				a[j+1] = a[j]         /* Shift value */
				index[j+1] = index[j] /* Shift index */
			}
			a[j+1] = int16(value) /* Write value */
			index[j+1] = i        /* Write index */
		}
	}
}

func SKP_Silk_insertion_sort_increasing_all_values(
	a []int32, /* I/O: Unsorted / Sorted vector                */
	L int32, /* I:   Vector length                           */
) {
	var i, j, value int32

	/* Sort vector elements by value, increasing order */
	for i = 1; i < L; i++ {
		value = a[i]
		for j = i - 1; j >= 0 && value < a[j]; j-- {
			a[j+1] = a[j] /* Shift value */
		}
		a[j+1] = value /* Write value */
	}
}
