package silk

import "unsafe"

func SKP_Silk_decode_pulses(psRC *SKP_Silk_range_coder_state, psDecCtrl *SKP_Silk_decoder_control, q []int32, frame_length int32) {
	var (
		i          int32
		j          int32
		k          int32
		iter       int32
		abs_q      int32
		nLS        int32
		bit        int32
		sum_pulses [30]int32
		nLshifts   [30]int32
		pulses_ptr *int32
		cdf_ptr    *uint16
	)
	SKP_Silk_range_decoder(&psDecCtrl.RateLevelIndex, psRC, SKP_Silk_rate_levels_CDF[psDecCtrl.Sigtype][:], SKP_Silk_rate_levels_CDF_offset)
	iter = frame_length / SHELL_CODEC_FRAME_LENGTH
	cdf_ptr = &SKP_Silk_pulses_per_block_CDF[psDecCtrl.RateLevelIndex][0]
	for i = 0; i < iter; i++ {
		nLshifts[i] = 0
		SKP_Silk_range_decoder(&sum_pulses[i], psRC, ([]uint16)(cdf_ptr), SKP_Silk_pulses_per_block_CDF_offset)
		for sum_pulses[i] == (MAX_PULSES + 1) {
			nLshifts[i]++
			SKP_Silk_range_decoder(&sum_pulses[i], psRC, SKP_Silk_pulses_per_block_CDF[N_RATE_LEVELS-1][:], SKP_Silk_pulses_per_block_CDF_offset)
		}
	}
	for i = 0; i < iter; i++ {
		if sum_pulses[i] > 0 {
			SKP_Silk_shell_decoder(([]int32)(&q[SKP_SMULBB(i, SHELL_CODEC_FRAME_LENGTH)]), psRC, sum_pulses[i])
		} else {
			memset(unsafe.Pointer(&q[SKP_SMULBB(i, SHELL_CODEC_FRAME_LENGTH)]), 0, size_t(SHELL_CODEC_FRAME_LENGTH*unsafe.Sizeof(int32(0))))
		}
	}
	for i = 0; i < iter; i++ {
		if nLshifts[i] > 0 {
			nLS = nLshifts[i]
			pulses_ptr = &q[SKP_SMULBB(i, SHELL_CODEC_FRAME_LENGTH)]
			for k = 0; k < SHELL_CODEC_FRAME_LENGTH; k++ {
				abs_q = *(*int32)(unsafe.Add(unsafe.Pointer(pulses_ptr), unsafe.Sizeof(int32(0))*uintptr(k)))
				for j = 0; j < nLS; j++ {
					abs_q = abs_q << 1
					SKP_Silk_range_decoder(&bit, psRC, SKP_Silk_lsb_CDF[:], 1)
					abs_q += bit
				}
				*(*int32)(unsafe.Add(unsafe.Pointer(pulses_ptr), unsafe.Sizeof(int32(0))*uintptr(k))) = abs_q
			}
		}
	}
	SKP_Silk_decode_signs(psRC, q, frame_length, psDecCtrl.Sigtype, psDecCtrl.QuantOffsetType, psDecCtrl.RateLevelIndex)
}
