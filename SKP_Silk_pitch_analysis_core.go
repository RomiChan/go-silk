package silk

import "unsafe"

const SCRATCH_SIZE = 22

func SKP_Silk_pitch_analysis_core(signal []int16, pitch_out []int32, lagIndex *int32, contourIndex *int32, LTPCorr_Q15 *int32, prevLag int32, search_thres1_Q16 int32, search_thres2_Q15 int32, Fs_kHz int32, complexity int32, forLJC int32) int32 {
	var (
		signal_8kHz           [480]int16
		signal_4kHz           [240]int16
		scratch_mem           [2880]int32
		input_signal_ptr      []int16
		filt_state            [7]int32
		i                     int32
		k                     int32
		d                     int32
		j                     int32
		C                     [4][221]int16
		target_ptr            []int16
		basis_ptr             []int16
		cross_corr            int32
		normalizer            int32
		energy                int32
		shift                 int32
		energy_basis          int32
		energy_target         int32
		d_srch                [24]int32
		d_comp                [221]int16
		Cmax                  int32
		length_d_srch         int32
		length_d_comp         int32
		sum                   int32
		threshold             int32
		temp32                int32
		CBimax                int32
		CBimax_new            int32
		CBimax_old            int32
		lag                   int32
		start_lag             int32
		end_lag               int32
		lag_new               int32
		CC                    [11]int32
		CCmax                 int32
		CCmax_b               int32
		CCmax_new_b           int32
		CCmax_new             int32
		energies_st3          [4][34][5]int32
		crosscorr_st3         [4][34][5]int32
		lag_counter           int32
		frame_length          int32
		frame_length_8kHz     int32
		frame_length_4kHz     int32
		max_sum_sq_length     int32
		sf_length             int32
		sf_length_8kHz        int32
		min_lag               int32
		min_lag_8kHz          int32
		min_lag_4kHz          int32
		max_lag               int32
		max_lag_8kHz          int32
		max_lag_4kHz          int32
		contour_bias          int32
		diff                  int32
		lz                    int32
		lshift                int32
		cbk_offset            int32
		cbk_size              int32
		nb_cbks_stage2        int32
		delta_lag_log2_sqr_Q7 int32
		lag_log2_Q7           int32
		prevLag_log2_Q7       int32
		prev_lag_bias_Q15     int32
		corr_thres_Q15        int32
	)
	frame_length = PITCH_EST_FRAME_LENGTH_MS * Fs_kHz
	frame_length_4kHz = PITCH_EST_FRAME_LENGTH_MS * 4
	frame_length_8kHz = PITCH_EST_FRAME_LENGTH_MS * 8
	sf_length = frame_length >> 3
	sf_length_8kHz = frame_length_8kHz >> 3
	min_lag = Fs_kHz * 2
	min_lag_4kHz = 2 * 4
	min_lag_8kHz = 2 * 8
	max_lag = Fs_kHz * 18
	max_lag_4kHz = 18 * 4
	max_lag_8kHz = 18 * 8
	memset(unsafe.Pointer(&C[0][0]), 0, size_t(unsafe.Sizeof(int16(0))*PITCH_EST_NB_SUBFR*(((PITCH_EST_MAX_FS_KHZ*18)>>1)+5)))
	if Fs_kHz == 16 {
		memset(unsafe.Pointer(&filt_state[0]), 0, size_t(unsafe.Sizeof(int32(0))*2))
		SKP_Silk_resampler_down2(filt_state[:], signal_8kHz[:], signal, frame_length)
	} else if Fs_kHz == 12 {
		var R23 [6]int32
		memset(unsafe.Pointer(&R23[0]), 0, size_t(unsafe.Sizeof(int32(0))*6))
		SKP_Silk_resampler_down2_3(R23[:], signal_8kHz[:], signal, PITCH_EST_FRAME_LENGTH_MS*12)
	} else if Fs_kHz == 24 {
		var filt_state_fix [8]int32
		memset(unsafe.Pointer(&filt_state_fix[0]), 0, size_t(unsafe.Sizeof(int32(0))*8))
		SKP_Silk_resampler_down3(filt_state_fix[:], signal_8kHz[:], signal, PITCH_EST_FRAME_LENGTH_MS*24)
	} else {
		memcpy(unsafe.Pointer(&signal_8kHz[0]), unsafe.Pointer(&signal[0]), size_t(uintptr(frame_length_8kHz)*unsafe.Sizeof(int16(0))))
	}
	memset(unsafe.Pointer(&filt_state[0]), 0, size_t(unsafe.Sizeof(int32(0))*2))
	SKP_Silk_resampler_down2(filt_state[:], signal_4kHz[:], signal_8kHz[:], frame_length_8kHz)
	for i = frame_length_4kHz - 1; i > 0; i-- {
		signal_4kHz[i] = SKP_SAT16(int32(int64(int32(signal_4kHz[i])) + int64(signal_4kHz[i-1])))
	}
	max_sum_sq_length = SKP_max_32(sf_length_8kHz, frame_length_4kHz>>1)
	shift = SKP_FIX_P_Ana_find_scaling(&signal_4kHz[0], frame_length_4kHz, max_sum_sq_length)
	if shift > 0 {
		for i = 0; i < frame_length_4kHz; i++ {
			signal_4kHz[i] = (signal_4kHz[i]) >> int64(shift)
		}
	}
	target_ptr = ([]int16)(&signal_4kHz[frame_length_4kHz>>1])
	for k = 0; k < 2; k++ {
		basis_ptr = ([]int16)((*int16)(unsafe.Add(unsafe.Pointer(&target_ptr[0]), -int(unsafe.Sizeof(int16(0))*uintptr(min_lag_4kHz)))))
		normalizer = 0
		cross_corr = 0
		cross_corr = SKP_Silk_inner_prod_aligned(target_ptr, basis_ptr, sf_length_8kHz)
		normalizer = SKP_Silk_inner_prod_aligned(basis_ptr, basis_ptr, sf_length_8kHz)
		normalizer = SKP_ADD_SAT32(normalizer, SKP_SMULBB(sf_length_8kHz, 4000))
		temp32 = cross_corr / (SKP_Silk_SQRT_APPROX(normalizer) + 1)
		C[k][min_lag_4kHz] = SKP_SAT16(temp32)
		for d = min_lag_4kHz + 1; d <= max_lag_4kHz; d++ {
			basis_ptr--
			cross_corr = SKP_Silk_inner_prod_aligned(target_ptr, basis_ptr, sf_length_8kHz)
			normalizer += SKP_SMULBB(int32(basis_ptr[0]), int32(basis_ptr[0])) - SKP_SMULBB(int32(basis_ptr[sf_length_8kHz]), int32(basis_ptr[sf_length_8kHz]))
			temp32 = cross_corr / (SKP_Silk_SQRT_APPROX(normalizer) + 1)
			C[k][d] = SKP_SAT16(temp32)
		}
		target_ptr += ([]int16)(sf_length_8kHz)
	}
	for i = max_lag_4kHz; i >= min_lag_4kHz; i-- {
		sum = int32(C[0][i]) + int32(C[1][i])
		sum = sum >> 1
		sum = SKP_SMLAWB(sum, sum, (-i)<<4)
		C[0][i] = int16(sum)
	}
	length_d_srch = complexity*2 + 4
	SKP_Silk_insertion_sort_decreasing_int16(([]int16)(&C[0][min_lag_4kHz]), d_srch[:], max_lag_4kHz-min_lag_4kHz+1, length_d_srch)
	target_ptr = ([]int16)(&signal_4kHz[frame_length_4kHz>>1])
	energy = SKP_Silk_inner_prod_aligned(target_ptr, target_ptr, frame_length_4kHz>>1)
	energy = SKP_ADD_POS_SAT32(energy, 1000)
	Cmax = int32(C[0][min_lag_4kHz])
	threshold = SKP_SMULBB(Cmax, Cmax)
	if (energy >> (4 + 2)) > threshold {
		memset(unsafe.Pointer(&pitch_out[0]), 0, size_t(PITCH_EST_NB_SUBFR*unsafe.Sizeof(int32(0))))
		*LTPCorr_Q15 = 0
		*lagIndex = 0
		*contourIndex = 0
		return 1
	}
	threshold = SKP_SMULWB(search_thres1_Q16, Cmax)
	for i = 0; i < length_d_srch; i++ {
		if int64(C[0][min_lag_4kHz+i]) > int64(threshold) {
			d_srch[i] = (d_srch[i] + min_lag_4kHz) << 1
		} else {
			length_d_srch = i
			break
		}
	}
	for i = min_lag_8kHz - 5; i < max_lag_8kHz+5; i++ {
		d_comp[i] = 0
	}
	for i = 0; i < length_d_srch; i++ {
		d_comp[d_srch[i]] = 1
	}
	for i = max_lag_8kHz + 3; i >= min_lag_8kHz; i-- {
		d_comp[i] += d_comp[i-1] + d_comp[i-2]
	}
	length_d_srch = 0
	for i = min_lag_8kHz; i < max_lag_8kHz+1; i++ {
		if d_comp[i+1] > 0 {
			d_srch[length_d_srch] = i
			length_d_srch++
		}
	}
	for i = max_lag_8kHz + 3; i >= min_lag_8kHz; i-- {
		d_comp[i] += d_comp[i-1] + d_comp[i-2] + d_comp[i-3]
	}
	length_d_comp = 0
	for i = min_lag_8kHz; i < max_lag_8kHz+4; i++ {
		if d_comp[i] > 0 {
			d_comp[length_d_comp] = int16(i - 2)
			length_d_comp++
		}
	}
	shift = SKP_FIX_P_Ana_find_scaling(&signal_8kHz[0], frame_length_8kHz, sf_length_8kHz)
	if shift > 0 {
		for i = 0; i < frame_length_8kHz; i++ {
			signal_8kHz[i] = (signal_8kHz[i]) >> int64(shift)
		}
	}
	memset(unsafe.Pointer(&C[0][0]), 0, size_t(PITCH_EST_NB_SUBFR*(((PITCH_EST_MAX_FS_KHZ*18)>>1)+5)*unsafe.Sizeof(int16(0))))
	target_ptr = ([]int16)(&signal_8kHz[frame_length_4kHz])
	for k = 0; k < PITCH_EST_NB_SUBFR; k++ {
		energy_target = SKP_Silk_inner_prod_aligned(target_ptr, target_ptr, sf_length_8kHz)
		for j = 0; j < length_d_comp; j++ {
			d = int32(d_comp[j])
			basis_ptr = ([]int16)((*int16)(unsafe.Add(unsafe.Pointer(&target_ptr[0]), -int(unsafe.Sizeof(int16(0))*uintptr(d)))))
			cross_corr = SKP_Silk_inner_prod_aligned(target_ptr, basis_ptr, sf_length_8kHz)
			energy_basis = SKP_Silk_inner_prod_aligned(basis_ptr, basis_ptr, sf_length_8kHz)
			if cross_corr > 0 {
				if energy_target > energy_basis {
					energy = energy_target
				} else {
					energy = energy_basis
				}
				lz = SKP_Silk_CLZ32(cross_corr)
				lshift = SKP_LIMIT_32(lz-1, 0, 15)
				temp32 = (cross_corr << lshift) / ((energy >> (15 - lshift)) + 1)
				temp32 = SKP_SMULWB(cross_corr, temp32)
				temp32 = SKP_ADD_SAT32(temp32, temp32)
				lz = SKP_Silk_CLZ32(temp32)
				lshift = SKP_LIMIT_32(lz-1, 0, 15)
				if energy_target < energy_basis {
					energy = energy_target
				} else {
					energy = energy_basis
				}
				C[k][d] = int16((temp32 << lshift) / ((energy >> (15 - lshift)) + 1))
			} else {
				C[k][d] = 0
			}
		}
		target_ptr += ([]int16)(sf_length_8kHz)
	}
	CCmax = 0x80000000
	CCmax_b = 0x80000000
	CBimax = 0
	lag = -1
	if prevLag > 0 {
		if Fs_kHz == 12 {
			prevLag = (prevLag << 1) / 3
		} else if Fs_kHz == 16 {
			prevLag = prevLag >> 1
		} else if Fs_kHz == 24 {
			prevLag = prevLag / 3
		}
		prevLag_log2_Q7 = SKP_Silk_lin2log(prevLag)
	} else {
		prevLag_log2_Q7 = 0
	}
	corr_thres_Q15 = SKP_SMULBB(search_thres2_Q15, search_thres2_Q15) >> 13
	if Fs_kHz == 8 && complexity > SKP_Silk_PITCH_EST_MIN_COMPLEX {
		nb_cbks_stage2 = PITCH_EST_NB_CBKS_STAGE2_EXT
	} else {
		nb_cbks_stage2 = PITCH_EST_NB_CBKS_STAGE2
	}
	for k = 0; k < length_d_srch; k++ {
		d = d_srch[k]
		for j = 0; j < nb_cbks_stage2; j++ {
			CC[j] = 0
			for i = 0; i < PITCH_EST_NB_SUBFR; i++ {
				CC[j] = CC[j] + int32(C[i][int64(d)+int64(SKP_Silk_CB_lags_stage2[i][j])])
			}
		}
		CCmax_new = 0x80000000
		CBimax_new = 0
		for i = 0; i < nb_cbks_stage2; i++ {
			if CC[i] > CCmax_new {
				CCmax_new = CC[i]
				CBimax_new = i
			}
		}
		lag_log2_Q7 = SKP_Silk_lin2log(d)
		if forLJC != 0 {
			CCmax_new_b = CCmax_new
		} else {
			CCmax_new_b = CCmax_new - (SKP_SMULBB(PITCH_EST_NB_SUBFR*PITCH_EST_SHORTLAG_BIAS_Q15, lag_log2_Q7) >> 7)
		}
		if prevLag > 0 {
			delta_lag_log2_sqr_Q7 = lag_log2_Q7 - prevLag_log2_Q7
			delta_lag_log2_sqr_Q7 = SKP_SMULBB(delta_lag_log2_sqr_Q7, delta_lag_log2_sqr_Q7) >> 7
			prev_lag_bias_Q15 = SKP_SMULBB(PITCH_EST_NB_SUBFR*PITCH_EST_PREVLAG_BIAS_Q15, *LTPCorr_Q15) >> 15
			prev_lag_bias_Q15 = (prev_lag_bias_Q15 * delta_lag_log2_sqr_Q7) / (delta_lag_log2_sqr_Q7 + (1 << 6))
			CCmax_new_b -= prev_lag_bias_Q15
		}
		if CCmax_new_b > CCmax_b && CCmax_new > corr_thres_Q15 && int64(SKP_Silk_CB_lags_stage2[0][CBimax_new]) <= int64(min_lag_8kHz) {
			CCmax_b = CCmax_new_b
			CCmax = CCmax_new
			lag = d
			CBimax = CBimax_new
		}
	}
	if int64(lag) == -1 {
		memset(unsafe.Pointer(&pitch_out[0]), 0, size_t(PITCH_EST_NB_SUBFR*unsafe.Sizeof(int32(0))))
		*LTPCorr_Q15 = 0
		*lagIndex = 0
		*contourIndex = 0
		return 1
	}
	if Fs_kHz > 8 {
		shift = SKP_FIX_P_Ana_find_scaling(&signal[0], frame_length, sf_length)
		if shift > 0 {
			input_signal_ptr = ([]int16)((*int16)(unsafe.Pointer(&scratch_mem[0])))
			for i = 0; i < frame_length; i++ {
				input_signal_ptr[i] = (signal[i]) >> int64(shift)
			}
		} else {
			input_signal_ptr = signal
		}
		CBimax_old = CBimax
		if Fs_kHz == 12 {
			lag = SKP_SMULBB(lag, 3) >> 1
		} else if Fs_kHz == 16 {
			lag = lag << 1
		} else {
			lag = SKP_SMULBB(lag, 3)
		}
		lag = SKP_LIMIT_int(lag, min_lag, max_lag)
		start_lag = SKP_max_int(lag-2, min_lag)
		end_lag = SKP_min_int(lag+2, max_lag)
		lag_new = lag
		CBimax = 0
		*LTPCorr_Q15 = SKP_Silk_SQRT_APPROX(CCmax << 13)
		CCmax = 0x80000000
		for k = 0; k < PITCH_EST_NB_SUBFR; k++ {
			pitch_out[k] = int32(int64(lag) + int64(SKP_Silk_CB_lags_stage2[k][CBimax_old]*2))
		}
		SKP_FIX_P_Ana_calc_corr_st3(crosscorr_st3, input_signal_ptr, start_lag, sf_length, complexity)
		SKP_FIX_P_Ana_calc_energy_st3(energies_st3, input_signal_ptr, start_lag, sf_length, complexity)
		lag_counter = 0
		contour_bias = PITCH_EST_FLATCONTOUR_BIAS_Q20 / lag
		cbk_size = int32(SKP_Silk_cbk_sizes_stage3[complexity])
		cbk_offset = int32(SKP_Silk_cbk_offsets_stage3[complexity])
		for d = start_lag; d <= end_lag; d++ {
			for j = cbk_offset; j < (cbk_offset + cbk_size); j++ {
				cross_corr = 0
				energy = 0
				for k = 0; k < PITCH_EST_NB_SUBFR; k++ {
					energy += (energies_st3[k][j][lag_counter]) >> 2
					cross_corr += (crosscorr_st3[k][j][lag_counter]) >> 2
				}
				if cross_corr > 0 {
					lz = SKP_Silk_CLZ32(cross_corr)
					lshift = SKP_LIMIT_32(lz-1, 0, 13)
					CCmax_new = (cross_corr << lshift) / ((energy >> (13 - lshift)) + 1)
					CCmax_new = int32(SKP_SAT16(CCmax_new))
					CCmax_new = SKP_SMULWB(cross_corr, CCmax_new)
					if CCmax_new > (SKP_int32_MAX >> 3) {
						CCmax_new = SKP_int32_MAX
					} else {
						CCmax_new = CCmax_new << 3
					}
					diff = j - (PITCH_EST_NB_CBKS_STAGE3_MAX >> 1)
					diff = diff * diff
					diff = SKP_int16_MAX - ((contour_bias * diff) >> 5)
					CCmax_new = SKP_SMULWB(CCmax_new, diff) << 1
				} else {
					CCmax_new = 0
				}
				if CCmax_new > CCmax && (d+int32(SKP_Silk_CB_lags_stage3[0][j])) <= max_lag {
					CCmax = CCmax_new
					lag_new = d
					CBimax = j
				}
			}
			lag_counter++
		}
		for k = 0; k < PITCH_EST_NB_SUBFR; k++ {
			pitch_out[k] = int32(int64(lag_new) + int64(SKP_Silk_CB_lags_stage3[k][CBimax]))
		}
		*lagIndex = lag_new - min_lag
		*contourIndex = CBimax
	} else {
		if CCmax > 0 {
			CCmax = CCmax
		} else {
			CCmax = 0
		}
		*LTPCorr_Q15 = SKP_Silk_SQRT_APPROX(CCmax << 13)
		for k = 0; k < PITCH_EST_NB_SUBFR; k++ {
			pitch_out[k] = int32(int64(lag) + int64(SKP_Silk_CB_lags_stage2[k][CBimax]))
		}
		*lagIndex = lag - min_lag_8kHz
		*contourIndex = CBimax
	}
	return 0
}
func SKP_FIX_P_Ana_calc_corr_st3(cross_corr_st3 [4][34][5]int32, signal []int16, start_lag int32, sf_length int32, complexity int32) {
	var (
		target_ptr  []int16
		basis_ptr   []int16
		cross_corr  int32
		i           int32
		j           int32
		k           int32
		lag_counter int32
		cbk_offset  int32
		cbk_size    int32
		delta       int32
		idx         int32
		scratch_mem [22]int32
	)
	cbk_offset = int32(SKP_Silk_cbk_offsets_stage3[complexity])
	cbk_size = int32(SKP_Silk_cbk_sizes_stage3[complexity])
	target_ptr = ([]int16)(&signal[sf_length<<2])
	for k = 0; k < PITCH_EST_NB_SUBFR; k++ {
		lag_counter = 0
		for j = int32(SKP_Silk_Lag_range_stage3[complexity][k][0]); int64(j) <= int64(SKP_Silk_Lag_range_stage3[complexity][k][1]); j++ {
			basis_ptr = ([]int16)((*int16)(unsafe.Add(unsafe.Pointer(&target_ptr[0]), -int(unsafe.Sizeof(int16(0))*uintptr(start_lag+j)))))
			cross_corr = SKP_Silk_inner_prod_aligned(target_ptr, basis_ptr, sf_length)
			scratch_mem[lag_counter] = cross_corr
			lag_counter++
		}
		delta = int32(SKP_Silk_Lag_range_stage3[complexity][k][0])
		for i = cbk_offset; i < (cbk_offset + cbk_size); i++ {
			idx = int32(int64(SKP_Silk_CB_lags_stage3[k][i]) - int64(delta))
			for j = 0; j < PITCH_EST_NB_STAGE3_LAGS; j++ {
				cross_corr_st3[k][i][j] = scratch_mem[idx+j]
			}
		}
		target_ptr += ([]int16)(sf_length)
	}
}
func SKP_FIX_P_Ana_calc_energy_st3(energies_st3 [4][34][5]int32, signal []int16, start_lag int32, sf_length int32, complexity int32) {
	var (
		target_ptr  []int16
		basis_ptr   []int16
		energy      int32
		k           int32
		i           int32
		j           int32
		lag_counter int32
		cbk_offset  int32
		cbk_size    int32
		delta       int32
		idx         int32
		scratch_mem [22]int32
	)
	cbk_offset = int32(SKP_Silk_cbk_offsets_stage3[complexity])
	cbk_size = int32(SKP_Silk_cbk_sizes_stage3[complexity])
	target_ptr = ([]int16)(&signal[sf_length<<2])
	for k = 0; k < PITCH_EST_NB_SUBFR; k++ {
		lag_counter = 0
		basis_ptr = ([]int16)((*int16)(unsafe.Add(unsafe.Pointer(&target_ptr[0]), -int(unsafe.Sizeof(int16(0))*uintptr(int64(start_lag)+int64(SKP_Silk_Lag_range_stage3[complexity][k][0]))))))
		energy = SKP_Silk_inner_prod_aligned(basis_ptr, basis_ptr, sf_length)
		scratch_mem[lag_counter] = energy
		lag_counter++
		for i = 1; int64(i) < int64(SKP_Silk_Lag_range_stage3[complexity][k][1]-SKP_Silk_Lag_range_stage3[complexity][k][0]+1); i++ {
			energy -= SKP_SMULBB(int32(basis_ptr[sf_length-i]), int32(basis_ptr[sf_length-i]))
			energy = SKP_ADD_SAT32(energy, SKP_SMULBB(int32(basis_ptr[-i]), int32(basis_ptr[-i])))
			scratch_mem[lag_counter] = energy
			lag_counter++
		}
		delta = int32(SKP_Silk_Lag_range_stage3[complexity][k][0])
		for i = cbk_offset; i < (cbk_offset + cbk_size); i++ {
			idx = int32(int64(SKP_Silk_CB_lags_stage3[k][i]) - int64(delta))
			for j = 0; j < PITCH_EST_NB_STAGE3_LAGS; j++ {
				energies_st3[k][i][j] = scratch_mem[idx+j]
			}
		}
		target_ptr += ([]int16)(sf_length)
	}
}
func SKP_FIX_P_Ana_find_scaling(signal *int16, signal_length int32, sum_sqr_len int32) int32 {
	var (
		nbits int32
		x_max int32
	)
	x_max = int32(SKP_Silk_int16_array_maxabs(([]int16)(signal), signal_length))
	if x_max < SKP_int16_MAX {
		nbits = 32 - SKP_Silk_CLZ32(SKP_SMULBB(x_max, x_max))
	} else {
		nbits = 30
	}
	nbits += 17 - SKP_Silk_CLZ16(int16(sum_sqr_len))
	if nbits < 31 {
		return 0
	} else {
		return nbits - 30
	}
}
