package silk

const MAX_LOOPS = 20

func SKP_Silk_NLSF_stabilize(NLSF_Q15 []int32, NDeltaMin_Q15 []int32, L int32) {
	var (
		center_freq_Q15 int32
		diff_Q15        int32
		min_center_Q15  int32
		max_center_Q15  int32
		min_diff_Q15    int32
		loops           int32
		i               int32
		I               int32 = 0
		k               int32
	)
	for loops = 0; loops < MAX_LOOPS; loops++ {
		min_diff_Q15 = NLSF_Q15[0] - NDeltaMin_Q15[0]
		I = 0
		for i = 1; i <= L-1; i++ {
			diff_Q15 = NLSF_Q15[i] - (NLSF_Q15[i-1] + NDeltaMin_Q15[i])
			if diff_Q15 < min_diff_Q15 {
				min_diff_Q15 = diff_Q15
				I = i
			}
		}
		diff_Q15 = (1 << 15) - (NLSF_Q15[L-1] + NDeltaMin_Q15[L])
		if diff_Q15 < min_diff_Q15 {
			min_diff_Q15 = diff_Q15
			I = L
		}
		if min_diff_Q15 >= 0 {
			return
		}
		if I == 0 {
			NLSF_Q15[0] = NDeltaMin_Q15[0]
		} else if I == L {
			NLSF_Q15[L-1] = (1 << 15) - NDeltaMin_Q15[L]
		} else {
			min_center_Q15 = 0
			for k = 0; k < I; k++ {
				min_center_Q15 += NDeltaMin_Q15[k]
			}
			min_center_Q15 += (NDeltaMin_Q15[I]) >> 1
			max_center_Q15 = 1 << 15
			for k = L; k > I; k-- {
				max_center_Q15 -= NDeltaMin_Q15[k]
			}
			max_center_Q15 -= NDeltaMin_Q15[I] - ((NDeltaMin_Q15[I]) >> 1)
			center_freq_Q15 = SKP_LIMIT_32(SKP_RSHIFT_ROUND(NLSF_Q15[I-1]+NLSF_Q15[I], 1), min_center_Q15, max_center_Q15)
			NLSF_Q15[I-1] = center_freq_Q15 - ((NDeltaMin_Q15[I]) >> 1)
			NLSF_Q15[I] = NLSF_Q15[I-1] + NDeltaMin_Q15[I]
		}
	}
	if loops == MAX_LOOPS {
		SKP_Silk_insertion_sort_increasing_all_values(NLSF_Q15, L)
		NLSF_Q15[0] = SKP_max_int(NLSF_Q15[0], NDeltaMin_Q15[0])
		for i = 1; i < L; i++ {
			NLSF_Q15[i] = SKP_max_int(NLSF_Q15[i], NLSF_Q15[i-1]+NDeltaMin_Q15[i])
		}
		NLSF_Q15[L-1] = SKP_min_int(NLSF_Q15[L-1], (1<<15)-NDeltaMin_Q15[L])
		for i = L - 2; i >= 0; i-- {
			NLSF_Q15[i] = SKP_min_int(NLSF_Q15[i], NLSF_Q15[i+1]-NDeltaMin_Q15[i+1])
		}
	}
}
