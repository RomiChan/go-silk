package silk

import "unsafe"

func combine_and_check(pulses_comb []int32, pulses_in []int32, max_pulses int32, len_ int32) int32 {
	var (
		k   int32
		sum int32
	)
	for k = 0; k < len_; k++ {
		sum = pulses_in[k*2] + pulses_in[k*2+1]
		if sum > max_pulses {
			return 1
		}
		pulses_comb[k] = sum
	}
	return 0
}
func SKP_Silk_encode_pulses(psRC *SKP_Silk_range_coder_state, sigtype int32, QuantOffsetType int32, q []int8, frame_length int32) {
	var (
		i              int32
		k              int32
		j              int32
		iter           int32
		bit            int32
		nLS            int32
		scale_down     int32
		RateLevelIndex int32 = 0
		abs_q          int32
		minSumBits_Q6  int32
		sumBits_Q6     int32
		abs_pulses     [480]int32
		sum_pulses     [30]int32
		nRshifts       [30]int32
		pulses_comb    [8]int32
		abs_pulses_ptr []int32
		pulses_ptr     []int8
		cdf_ptr        []uint16
		nBits_ptr      []int16
	)
	memset(unsafe.Pointer(&pulses_comb[0]), 0, size_t(unsafe.Sizeof(int32(0))*8))
	iter = frame_length / SHELL_CODEC_FRAME_LENGTH
	for i = 0; i < frame_length; i += 4 {
		abs_pulses[i+0] = int32(SKP_abs(int64(q[i+0])))
		abs_pulses[i+1] = int32(SKP_abs(int64(q[i+1])))
		abs_pulses[i+2] = int32(SKP_abs(int64(q[i+2])))
		abs_pulses[i+3] = int32(SKP_abs(int64(q[i+3])))
	}
	abs_pulses_ptr = ([]int32)(abs_pulses[:])
	for i = 0; i < iter; i++ {
		nRshifts[i] = 0
		for {
			scale_down = combine_and_check(pulses_comb[:], abs_pulses_ptr, SKP_Silk_max_pulses_table[0], 8)
			scale_down += combine_and_check(pulses_comb[:], pulses_comb[:], SKP_Silk_max_pulses_table[1], 4)
			scale_down += combine_and_check(pulses_comb[:], pulses_comb[:], SKP_Silk_max_pulses_table[2], 2)
			sum_pulses[i] = pulses_comb[0] + pulses_comb[1]
			if sum_pulses[i] > SKP_Silk_max_pulses_table[3] {
				scale_down++
			}
			if scale_down != 0 {
				nRshifts[i]++
				for k = 0; k < SHELL_CODEC_FRAME_LENGTH; k++ {
					abs_pulses_ptr[k] = (abs_pulses_ptr[k]) >> 1
				}
			} else {
				break
			}
		}
		abs_pulses_ptr += SHELL_CODEC_FRAME_LENGTH
	}
	minSumBits_Q6 = SKP_int32_MAX
	for k = 0; k < N_RATE_LEVELS-1; k++ {
		nBits_ptr = ([]int16)(SKP_Silk_pulses_per_block_BITS_Q6[k][:])
		sumBits_Q6 = int32(SKP_Silk_rate_levels_BITS_Q6[sigtype][k])
		for i = 0; i < iter; i++ {
			if nRshifts[i] > 0 {
				sumBits_Q6 += int32(nBits_ptr[MAX_PULSES+1])
			} else {
				sumBits_Q6 += int32(nBits_ptr[sum_pulses[i]])
			}
		}
		if sumBits_Q6 < minSumBits_Q6 {
			minSumBits_Q6 = sumBits_Q6
			RateLevelIndex = k
		}
	}
	SKP_Silk_range_encoder(psRC, RateLevelIndex, SKP_Silk_rate_levels_CDF[sigtype][:])
	cdf_ptr = ([]uint16)(SKP_Silk_pulses_per_block_CDF[RateLevelIndex][:])
	for i = 0; i < iter; i++ {
		if nRshifts[i] == 0 {
			SKP_Silk_range_encoder(psRC, sum_pulses[i], cdf_ptr)
		} else {
			SKP_Silk_range_encoder(psRC, MAX_PULSES+1, cdf_ptr)
			for k = 0; k < nRshifts[i]-1; k++ {
				SKP_Silk_range_encoder(psRC, MAX_PULSES+1, SKP_Silk_pulses_per_block_CDF[N_RATE_LEVELS-1][:])
			}
			SKP_Silk_range_encoder(psRC, sum_pulses[i], SKP_Silk_pulses_per_block_CDF[N_RATE_LEVELS-1][:])
		}
	}
	for i = 0; i < iter; i++ {
		if sum_pulses[i] > 0 {
			SKP_Silk_shell_encoder(psRC, ([]int32)(&abs_pulses[i*SHELL_CODEC_FRAME_LENGTH]))
		}
	}
	for i = 0; i < iter; i++ {
		if nRshifts[i] > 0 {
			pulses_ptr = ([]int8)(&q[i*SHELL_CODEC_FRAME_LENGTH])
			nLS = nRshifts[i] - 1
			for k = 0; k < SHELL_CODEC_FRAME_LENGTH; k++ {
				abs_q = int32(int8(SKP_abs(int64(pulses_ptr[k]))))
				for j = nLS; j > 0; j-- {
					bit = (abs_q >> j) & 1
					SKP_Silk_range_encoder(psRC, bit, SKP_Silk_lsb_CDF[:])
				}
				bit = abs_q & 1
				SKP_Silk_range_encoder(psRC, bit, SKP_Silk_lsb_CDF[:])
			}
		}
	}
	SKP_Silk_encode_signs(psRC, q, frame_length, sigtype, QuantOffsetType, RateLevelIndex)
}
