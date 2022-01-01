package silk

import (
	"math"
	"unsafe"
)

const SCRATCH_SIZE = 22

func SKP_Silk_pitch_analysis_core(signal *int16, pitch_out *int32, lagIndex *int32, contourIndex *int32, LTPCorr_Q15 *int32, prevLag int32, search_thres1_Q16 int32, search_thres2_Q15 int32, Fs_kHz int32, complexity int32, forLJC int32) int32 {
	var (
		signal_8kHz           [480]int16
		signal_4kHz           [240]int16
		scratch_mem           [2880]int32
		input_signal_ptr      *int16
		filt_state            [7]int32
		i                     int32
		k                     int32
		d                     int32
		j                     int32
		C                     [4][221]int16
		target_ptr            *int16
		basis_ptr             *int16
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
	frame_length = int32(PITCH_EST_FRAME_LENGTH_MS * int64(Fs_kHz))
	frame_length_4kHz = PITCH_EST_FRAME_LENGTH_MS * 4
	frame_length_8kHz = PITCH_EST_FRAME_LENGTH_MS * 8
	sf_length = frame_length >> 3
	sf_length_8kHz = frame_length_8kHz >> 3
	min_lag = int32(int64(Fs_kHz) * 2)
	min_lag_4kHz = 2 * 4
	min_lag_8kHz = 2 * 8
	max_lag = int32(int64(Fs_kHz) * 18)
	max_lag_4kHz = 18 * 4
	max_lag_8kHz = 18 * 8
	memset(unsafe.Pointer(&C[0][0]), 0, size_t(unsafe.Sizeof(int16(0))*PITCH_EST_NB_SUBFR*uintptr(((PITCH_EST_MAX_FS_KHZ*18)>>1)+5)))
	if int64(Fs_kHz) == 16 {
		memset(unsafe.Pointer(&filt_state[0]), 0, size_t(unsafe.Sizeof(int32(0))*2))
		SKP_Silk_resampler_down2(&filt_state[0], &signal_8kHz[0], signal, frame_length)
	} else if int64(Fs_kHz) == 12 {
		var R23 [6]int32
		memset(unsafe.Pointer(&R23[0]), 0, size_t(unsafe.Sizeof(int32(0))*6))
		SKP_Silk_resampler_down2_3(&R23[0], &signal_8kHz[0], signal, PITCH_EST_FRAME_LENGTH_MS*12)
	} else if int64(Fs_kHz) == 24 {
		var filt_state_fix [8]int32
		memset(unsafe.Pointer(&filt_state_fix[0]), 0, size_t(unsafe.Sizeof(int32(0))*8))
		SKP_Silk_resampler_down3(&filt_state_fix[0], &signal_8kHz[0], signal, PITCH_EST_FRAME_LENGTH_MS*24)
	} else {
		memcpy(unsafe.Pointer(&signal_8kHz[0]), unsafe.Pointer(signal), size_t(uintptr(frame_length_8kHz)*unsafe.Sizeof(int16(0))))
	}
	memset(unsafe.Pointer(&filt_state[0]), 0, size_t(unsafe.Sizeof(int32(0))*2))
	SKP_Silk_resampler_down2(&filt_state[0], &signal_4kHz[0], &signal_8kHz[0], frame_length_8kHz)
	for i = int32(int64(frame_length_4kHz) - 1); int64(i) > 0; i-- {
		signal_4kHz[i] = SKP_SAT16(int16(int64(int32(signal_4kHz[i])) + int64(signal_4kHz[int64(i)-1])))
	}
	max_sum_sq_length = SKP_max_32(sf_length_8kHz, int32(int64(frame_length_4kHz)>>1))
	shift = SKP_FIX_P_Ana_find_scaling(&signal_4kHz[0], frame_length_4kHz, max_sum_sq_length)
	if int64(shift) > 0 {
		for i = 0; int64(i) < int64(frame_length_4kHz); i++ {
			signal_4kHz[i] = (signal_4kHz[i]) >> int64(shift)
		}
	}
	target_ptr = &signal_4kHz[int64(frame_length_4kHz)>>1]
	for k = 0; int64(k) < 2; k++ {
		basis_ptr = (*int16)(unsafe.Add(unsafe.Pointer(target_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(min_lag_4kHz))))
		normalizer = 0
		cross_corr = 0
		cross_corr = SKP_Silk_inner_prod_aligned(([]int16)(target_ptr), ([]int16)(basis_ptr), sf_length_8kHz)
		normalizer = SKP_Silk_inner_prod_aligned(([]int16)(basis_ptr), ([]int16)(basis_ptr), sf_length_8kHz)
		if ((int64(normalizer) + int64(SKP_SMULBB(sf_length_8kHz, 4000))) & 0x80000000) == 0 {
			if ((int64(normalizer) & int64(SKP_SMULBB(sf_length_8kHz, 4000))) & 0x80000000) != 0 {
				normalizer = math.MinInt32
			} else {
				normalizer = int32(int64(normalizer) + int64(SKP_SMULBB(sf_length_8kHz, 4000)))
			}
		} else if ((int64(normalizer) | int64(SKP_SMULBB(sf_length_8kHz, 4000))) & 0x80000000) == 0 {
			normalizer = SKP_int32_MAX
		} else {
			normalizer = int32(int64(normalizer) + int64(SKP_SMULBB(sf_length_8kHz, 4000)))
		}
		temp32 = int32(int64(cross_corr) / (int64(SKP_Silk_SQRT_APPROX(normalizer)) + 1))
		C[k][min_lag_4kHz] = SKP_SAT16(int16(temp32))
		for d = int32(int64(min_lag_4kHz) + 1); int64(d) <= int64(max_lag_4kHz); d++ {
			basis_ptr = (*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*1)))
			cross_corr = SKP_Silk_inner_prod_aligned(([]int16)(target_ptr), ([]int16)(basis_ptr), sf_length_8kHz)
			normalizer += int32(int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), unsafe.Sizeof(int16(0))*0))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), unsafe.Sizeof(int16(0))*0))))) - int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), unsafe.Sizeof(int16(0))*uintptr(sf_length_8kHz)))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), unsafe.Sizeof(int16(0))*uintptr(sf_length_8kHz)))))))
			temp32 = int32(int64(cross_corr) / (int64(SKP_Silk_SQRT_APPROX(normalizer)) + 1))
			C[k][d] = SKP_SAT16(int16(temp32))
		}
		target_ptr = (*int16)(unsafe.Add(unsafe.Pointer(target_ptr), unsafe.Sizeof(int16(0))*uintptr(sf_length_8kHz)))
	}
	for i = max_lag_4kHz; int64(i) >= int64(min_lag_4kHz); i-- {
		sum = int32(int64(C[0][i]) + int64(C[1][i]))
		sum = sum >> 1
		sum = SKP_SMLAWB(sum, sum, int32(int64(-i)<<4))
		C[0][i] = int16(sum)
	}
	length_d_srch = int32(int64(complexity)*2 + 4)
	SKP_Silk_insertion_sort_decreasing_int16(&C[0][min_lag_4kHz], &d_srch[0], int32(int64(max_lag_4kHz)-int64(min_lag_4kHz)+1), length_d_srch)
	target_ptr = &signal_4kHz[int64(frame_length_4kHz)>>1]
	energy = SKP_Silk_inner_prod_aligned(([]int16)(target_ptr), ([]int16)(target_ptr), int32(int64(frame_length_4kHz)>>1))
	energy = SKP_ADD_POS_SAT32(energy, 1000)
	Cmax = int32(C[0][min_lag_4kHz])
	threshold = SKP_SMULBB(Cmax, Cmax)
	if (int64(energy) >> (4 + 2)) > int64(threshold) {
		memset(unsafe.Pointer(pitch_out), 0, size_t(PITCH_EST_NB_SUBFR*unsafe.Sizeof(int32(0))))
		*LTPCorr_Q15 = 0
		*lagIndex = 0
		*contourIndex = 0
		return 1
	}
	threshold = SKP_SMULWB(search_thres1_Q16, Cmax)
	for i = 0; int64(i) < int64(length_d_srch); i++ {
		if int64(C[0][int64(min_lag_4kHz)+int64(i)]) > int64(threshold) {
			d_srch[i] = int32((int64(d_srch[i]) + int64(min_lag_4kHz)) << 1)
		} else {
			length_d_srch = i
			break
		}
	}
	for i = int32(int64(min_lag_8kHz) - 5); int64(i) < int64(max_lag_8kHz)+5; i++ {
		d_comp[i] = 0
	}
	for i = 0; int64(i) < int64(length_d_srch); i++ {
		d_comp[d_srch[i]] = 1
	}
	for i = int32(int64(max_lag_8kHz) + 3); int64(i) >= int64(min_lag_8kHz); i-- {
		d_comp[i] += int16(int64(d_comp[int64(i)-1]) + int64(d_comp[int64(i)-2]))
	}
	length_d_srch = 0
	for i = min_lag_8kHz; int64(i) < int64(max_lag_8kHz)+1; i++ {
		if int64(d_comp[int64(i)+1]) > 0 {
			d_srch[length_d_srch] = i
			length_d_srch++
		}
	}
	for i = int32(int64(max_lag_8kHz) + 3); int64(i) >= int64(min_lag_8kHz); i-- {
		d_comp[i] += int16(int64(d_comp[int64(i)-1]) + int64(d_comp[int64(i)-2]) + int64(d_comp[int64(i)-3]))
	}
	length_d_comp = 0
	for i = min_lag_8kHz; int64(i) < int64(max_lag_8kHz)+4; i++ {
		if int64(d_comp[i]) > 0 {
			d_comp[length_d_comp] = int16(int64(i) - 2)
			length_d_comp++
		}
	}
	shift = SKP_FIX_P_Ana_find_scaling(&signal_8kHz[0], frame_length_8kHz, sf_length_8kHz)
	if int64(shift) > 0 {
		for i = 0; int64(i) < int64(frame_length_8kHz); i++ {
			signal_8kHz[i] = (signal_8kHz[i]) >> int64(shift)
		}
	}
	memset(unsafe.Pointer(&C[0][0]), 0, size_t(uintptr(PITCH_EST_NB_SUBFR*(((PITCH_EST_MAX_FS_KHZ*18)>>1)+5))*unsafe.Sizeof(int16(0))))
	target_ptr = &signal_8kHz[frame_length_4kHz]
	for k = 0; int64(k) < PITCH_EST_NB_SUBFR; k++ {
		energy_target = SKP_Silk_inner_prod_aligned(([]int16)(target_ptr), ([]int16)(target_ptr), sf_length_8kHz)
		for j = 0; int64(j) < int64(length_d_comp); j++ {
			d = int32(d_comp[j])
			basis_ptr = (*int16)(unsafe.Add(unsafe.Pointer(target_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(d))))
			cross_corr = SKP_Silk_inner_prod_aligned(([]int16)(target_ptr), ([]int16)(basis_ptr), sf_length_8kHz)
			energy_basis = SKP_Silk_inner_prod_aligned(([]int16)(basis_ptr), ([]int16)(basis_ptr), sf_length_8kHz)
			if int64(cross_corr) > 0 {
				if int64(energy_target) > int64(energy_basis) {
					energy = energy_target
				} else {
					energy = energy_basis
				}
				lz = SKP_Silk_CLZ32(cross_corr)
				if (int64(lz) - 1) > 15 {
					lshift = 15
				} else if (int64(lz) - 1) < 0 {
					lshift = 0
				} else {
					lshift = int32(int64(lz) - 1)
				}
				temp32 = int32((int64(cross_corr) << int64(lshift)) / ((int64(energy) >> (15 - int64(lshift))) + 1))
				temp32 = SKP_SMULWB(cross_corr, temp32)
				if ((int64(temp32) + int64(temp32)) & 0x80000000) == 0 {
					if ((int64(temp32) & int64(temp32)) & 0x80000000) != 0 {
						temp32 = math.MinInt32
					} else {
						temp32 = int32(int64(temp32) + int64(temp32))
					}
				} else if ((int64(temp32) | int64(temp32)) & 0x80000000) == 0 {
					temp32 = SKP_int32_MAX
				} else {
					temp32 = int32(int64(temp32) + int64(temp32))
				}
				lz = SKP_Silk_CLZ32(temp32)
				if (int64(lz) - 1) > 15 {
					lshift = 15
				} else if (int64(lz) - 1) < 0 {
					lshift = 0
				} else {
					lshift = int32(int64(lz) - 1)
				}
				if int64(energy_target) < int64(energy_basis) {
					energy = energy_target
				} else {
					energy = energy_basis
				}
				C[k][d] = int16(int32((int64(temp32) << int64(lshift)) / ((int64(energy) >> (15 - int64(lshift))) + 1)))
			} else {
				C[k][d] = 0
			}
		}
		target_ptr = (*int16)(unsafe.Add(unsafe.Pointer(target_ptr), unsafe.Sizeof(int16(0))*uintptr(sf_length_8kHz)))
	}
	CCmax = math.MinInt32
	CCmax_b = math.MinInt32
	CBimax = 0
	lag = -1
	if int64(prevLag) > 0 {
		if int64(Fs_kHz) == 12 {
			prevLag = int32((int64(prevLag) << 1) / 3)
		} else if int64(Fs_kHz) == 16 {
			prevLag = prevLag >> 1
		} else if int64(Fs_kHz) == 24 {
			prevLag = int32(int64(prevLag) / 3)
		}
		prevLag_log2_Q7 = SKP_Silk_lin2log(prevLag)
	} else {
		prevLag_log2_Q7 = 0
	}
	corr_thres_Q15 = SKP_SMULBB(search_thres2_Q15, search_thres2_Q15) >> 13
	if int64(Fs_kHz) == 8 && int64(complexity) > SKP_Silk_PITCH_EST_MIN_COMPLEX {
		nb_cbks_stage2 = PITCH_EST_NB_CBKS_STAGE2_EXT
	} else {
		nb_cbks_stage2 = PITCH_EST_NB_CBKS_STAGE2
	}
	for k = 0; int64(k) < int64(length_d_srch); k++ {
		d = d_srch[k]
		for j = 0; int64(j) < int64(nb_cbks_stage2); j++ {
			CC[j] = 0
			for i = 0; int64(i) < PITCH_EST_NB_SUBFR; i++ {
				CC[j] = int32(int64(CC[j]) + int64(C[i][int64(d)+int64(SKP_Silk_CB_lags_stage2[i][j])]))
			}
		}
		CCmax_new = math.MinInt32
		CBimax_new = 0
		for i = 0; int64(i) < int64(nb_cbks_stage2); i++ {
			if int64(CC[i]) > int64(CCmax_new) {
				CCmax_new = CC[i]
				CBimax_new = i
			}
		}
		lag_log2_Q7 = SKP_Silk_lin2log(d)
		if int64(forLJC) != 0 {
			CCmax_new_b = CCmax_new
		} else {
			CCmax_new_b = int32(int64(CCmax_new) - (int64(SKP_SMULBB(PITCH_EST_NB_SUBFR*PITCH_EST_SHORTLAG_BIAS_Q15, lag_log2_Q7)) >> 7))
		}
		if int64(prevLag) > 0 {
			delta_lag_log2_sqr_Q7 = int32(int64(lag_log2_Q7) - int64(prevLag_log2_Q7))
			delta_lag_log2_sqr_Q7 = SKP_SMULBB(delta_lag_log2_sqr_Q7, delta_lag_log2_sqr_Q7) >> 7
			prev_lag_bias_Q15 = SKP_SMULBB(PITCH_EST_NB_SUBFR*PITCH_EST_PREVLAG_BIAS_Q15, *LTPCorr_Q15) >> 15
			prev_lag_bias_Q15 = int32((int64(prev_lag_bias_Q15) * int64(delta_lag_log2_sqr_Q7)) / (int64(delta_lag_log2_sqr_Q7) + (1 << 6)))
			CCmax_new_b -= prev_lag_bias_Q15
		}
		if int64(CCmax_new_b) > int64(CCmax_b) && int64(CCmax_new) > int64(corr_thres_Q15) && int64(SKP_Silk_CB_lags_stage2[0][CBimax_new]) <= int64(min_lag_8kHz) {
			CCmax_b = CCmax_new_b
			CCmax = CCmax_new
			lag = d
			CBimax = CBimax_new
		}
	}
	if int64(lag) == -1 {
		memset(unsafe.Pointer(pitch_out), 0, size_t(PITCH_EST_NB_SUBFR*unsafe.Sizeof(int32(0))))
		*LTPCorr_Q15 = 0
		*lagIndex = 0
		*contourIndex = 0
		return 1
	}
	if int64(Fs_kHz) > 8 {
		shift = SKP_FIX_P_Ana_find_scaling(signal, frame_length, sf_length)
		if int64(shift) > 0 {
			input_signal_ptr = (*int16)(unsafe.Pointer(&scratch_mem[0]))
			for i = 0; int64(i) < int64(frame_length); i++ {
				*(*int16)(unsafe.Add(unsafe.Pointer(input_signal_ptr), unsafe.Sizeof(int16(0))*uintptr(i))) = (*(*int16)(unsafe.Add(unsafe.Pointer(signal), unsafe.Sizeof(int16(0))*uintptr(i)))) >> int64(shift)
			}
		} else {
			input_signal_ptr = signal
		}
		CBimax_old = CBimax
		if int64(Fs_kHz) == 12 {
			lag = SKP_SMULBB(lag, 3) >> 1
		} else if int64(Fs_kHz) == 16 {
			lag = int32(int64(lag) << 1)
		} else {
			lag = SKP_SMULBB(lag, 3)
		}
		lag = SKP_LIMIT_int(lag, min_lag, max_lag)
		start_lag = SKP_max_int(int32(int64(lag)-2), min_lag)
		end_lag = SKP_min_int(int32(int64(lag)+2), max_lag)
		lag_new = lag
		CBimax = 0
		*LTPCorr_Q15 = SKP_Silk_SQRT_APPROX(int32(int64(CCmax) << 13))
		CCmax = math.MinInt32
		for k = 0; int64(k) < PITCH_EST_NB_SUBFR; k++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(pitch_out), unsafe.Sizeof(int32(0))*uintptr(k))) = int32(int64(lag) + int64(SKP_Silk_CB_lags_stage2[k][CBimax_old])*2)
		}
		SKP_FIX_P_Ana_calc_corr_st3(crosscorr_st3, ([]int16)(input_signal_ptr), start_lag, sf_length, complexity)
		SKP_FIX_P_Ana_calc_energy_st3(energies_st3, ([]int16)(input_signal_ptr), start_lag, sf_length, complexity)
		lag_counter = 0
		contour_bias = int32(PITCH_EST_FLATCONTOUR_BIAS_Q20 / int64(lag))
		cbk_size = int32(SKP_Silk_cbk_sizes_stage3[complexity])
		cbk_offset = int32(SKP_Silk_cbk_offsets_stage3[complexity])
		for d = start_lag; int64(d) <= int64(end_lag); d++ {
			for j = cbk_offset; int64(j) < (int64(cbk_offset) + int64(cbk_size)); j++ {
				cross_corr = 0
				energy = 0
				for k = 0; int64(k) < PITCH_EST_NB_SUBFR; k++ {
					energy += (energies_st3[k][j][lag_counter]) >> 2
					cross_corr += (crosscorr_st3[k][j][lag_counter]) >> 2
				}
				if int64(cross_corr) > 0 {
					lz = SKP_Silk_CLZ32(cross_corr)
					if (int64(lz) - 1) > 13 {
						lshift = 13
					} else if (int64(lz) - 1) < 0 {
						lshift = 0
					} else {
						lshift = int32(int64(lz) - 1)
					}
					CCmax_new = int32((int64(cross_corr) << int64(lshift)) / ((int64(energy) >> (13 - int64(lshift))) + 1))
					CCmax_new = int32(SKP_SAT16(int16(CCmax_new)))
					CCmax_new = SKP_SMULWB(cross_corr, CCmax_new)
					if int64(CCmax_new) > (SKP_int32_MAX >> 3) {
						CCmax_new = SKP_int32_MAX
					} else {
						CCmax_new = int32(int64(CCmax_new) << 3)
					}
					diff = int32(int64(j) - (PITCH_EST_NB_CBKS_STAGE3_MAX >> 1))
					diff = int32(int64(diff) * int64(diff))
					diff = int32(SKP_int16_MAX - ((int64(contour_bias) * int64(diff)) >> 5))
					CCmax_new = int32(int64(SKP_SMULWB(CCmax_new, diff)) << 1)
				} else {
					CCmax_new = 0
				}
				if int64(CCmax_new) > int64(CCmax) && (int64(d)+int64(SKP_Silk_CB_lags_stage3[0][j])) <= int64(max_lag) {
					CCmax = CCmax_new
					lag_new = d
					CBimax = j
				}
			}
			lag_counter++
		}
		for k = 0; int64(k) < PITCH_EST_NB_SUBFR; k++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(pitch_out), unsafe.Sizeof(int32(0))*uintptr(k))) = int32(int64(lag_new) + int64(SKP_Silk_CB_lags_stage3[k][CBimax]))
		}
		*lagIndex = int32(int64(lag_new) - int64(min_lag))
		*contourIndex = CBimax
	} else {
		if int64(CCmax) > 0 {
			CCmax = CCmax
		} else {
			CCmax = 0
		}
		*LTPCorr_Q15 = SKP_Silk_SQRT_APPROX(int32(int64(CCmax) << 13))
		for k = 0; int64(k) < PITCH_EST_NB_SUBFR; k++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(pitch_out), unsafe.Sizeof(int32(0))*uintptr(k))) = int32(int64(lag) + int64(SKP_Silk_CB_lags_stage2[k][CBimax]))
		}
		*lagIndex = int32(int64(lag) - int64(min_lag_8kHz))
		*contourIndex = CBimax
	}
	return 0
}
func SKP_FIX_P_Ana_calc_corr_st3(cross_corr_st3 [4][34][5]int32, signal []int16, start_lag int32, sf_length int32, complexity int32) {
	var (
		target_ptr  *int16
		basis_ptr   *int16
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
	target_ptr = &signal[int64(sf_length)<<2]
	for k = 0; int64(k) < PITCH_EST_NB_SUBFR; k++ {
		lag_counter = 0
		for j = int32(SKP_Silk_Lag_range_stage3[complexity][k][0]); int64(j) <= int64(SKP_Silk_Lag_range_stage3[complexity][k][1]); j++ {
			basis_ptr = (*int16)(unsafe.Add(unsafe.Pointer(target_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(int64(start_lag)+int64(j)))))
			cross_corr = SKP_Silk_inner_prod_aligned(([]int16)(target_ptr), ([]int16)(basis_ptr), sf_length)
			scratch_mem[lag_counter] = cross_corr
			lag_counter++
		}
		delta = int32(SKP_Silk_Lag_range_stage3[complexity][k][0])
		for i = cbk_offset; int64(i) < (int64(cbk_offset) + int64(cbk_size)); i++ {
			idx = int32(int64(SKP_Silk_CB_lags_stage3[k][i]) - int64(delta))
			for j = 0; int64(j) < PITCH_EST_NB_STAGE3_LAGS; j++ {
				cross_corr_st3[k][i][j] = scratch_mem[int64(idx)+int64(j)]
			}
		}
		target_ptr = (*int16)(unsafe.Add(unsafe.Pointer(target_ptr), unsafe.Sizeof(int16(0))*uintptr(sf_length)))
	}
}
func SKP_FIX_P_Ana_calc_energy_st3(energies_st3 [4][34][5]int32, signal []int16, start_lag int32, sf_length int32, complexity int32) {
	var (
		target_ptr  *int16
		basis_ptr   *int16
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
	target_ptr = &signal[int64(sf_length)<<2]
	for k = 0; int64(k) < PITCH_EST_NB_SUBFR; k++ {
		lag_counter = 0
		basis_ptr = (*int16)(unsafe.Add(unsafe.Pointer(target_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(int64(start_lag)+int64(SKP_Silk_Lag_range_stage3[complexity][k][0])))))
		energy = SKP_Silk_inner_prod_aligned(([]int16)(basis_ptr), ([]int16)(basis_ptr), sf_length)
		scratch_mem[lag_counter] = energy
		lag_counter++
		for i = 1; int64(i) < (int64(SKP_Silk_Lag_range_stage3[complexity][k][1]) - int64(SKP_Silk_Lag_range_stage3[complexity][k][0]) + 1); i++ {
			energy -= SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(sf_length)-int64(i))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(sf_length)-int64(i))))))
			if ((int64(energy) + int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i)))))))) & 0x80000000) == 0 {
				if ((int64(energy) & int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i)))))))) & 0x80000000) != 0 {
					energy = math.MinInt32
				} else {
					energy = int32(int64(energy) + int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i))))))))
				}
			} else if ((int64(energy) | int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i)))))))) & 0x80000000) == 0 {
				energy = SKP_int32_MAX
			} else {
				energy = int32(int64(energy) + int64(SKP_SMULBB(int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i))))), int32(*(*int16)(unsafe.Add(unsafe.Pointer(basis_ptr), -int(unsafe.Sizeof(int16(0))*uintptr(i))))))))
			}
			scratch_mem[lag_counter] = energy
			lag_counter++
		}
		delta = int32(SKP_Silk_Lag_range_stage3[complexity][k][0])
		for i = cbk_offset; int64(i) < (int64(cbk_offset) + int64(cbk_size)); i++ {
			idx = int32(int64(SKP_Silk_CB_lags_stage3[k][i]) - int64(delta))
			for j = 0; int64(j) < PITCH_EST_NB_STAGE3_LAGS; j++ {
				energies_st3[k][i][j] = scratch_mem[int64(idx)+int64(j)]
			}
		}
		target_ptr = (*int16)(unsafe.Add(unsafe.Pointer(target_ptr), unsafe.Sizeof(int16(0))*uintptr(sf_length)))
	}
}
func SKP_FIX_P_Ana_find_scaling(signal *int16, signal_length int32, sum_sqr_len int32) int32 {
	var (
		nbits int32
		x_max int32
	)
	x_max = int32(SKP_Silk_int16_array_maxabs(([]int16)(signal), signal_length))
	if int64(x_max) < SKP_int16_MAX {
		nbits = int32(32 - int64(SKP_Silk_CLZ32(SKP_SMULBB(x_max, x_max))))
	} else {
		nbits = 30
	}
	nbits += int32(17 - int64(SKP_Silk_CLZ16(int16(sum_sqr_len))))
	if int64(nbits) < 31 {
		return 0
	} else {
		return int32(int64(nbits) - 30)
	}
}
