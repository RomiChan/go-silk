package silk

var freq_table_Q16 [27]int16 = [27]int16{0x2F4F, 9804, 8235, 7100, 6239, 5565, 5022, 4575, 4202, 3885, 3612, 3375, 3167, 2984, 2820, 2674, 2542, 2422, 2313, 2214, 2123, 2038, 1961, 1889, 1822, 1760, 1702}

func SKP_Silk_apply_sine_window(px_win []int16, px []int16, win_type int32, length int32) {
	var (
		k      int32
		f_Q16  int32
		c_Q16  int32
		S0_Q16 int32
		S1_Q16 int32
	)
	k = int32((int64(length) >> 2) - 4)
	f_Q16 = int32(freq_table_Q16[k])
	c_Q16 = SKP_SMULWB(f_Q16, -f_Q16)
	if int64(win_type) == 1 {
		S0_Q16 = 0
		S1_Q16 = int32(int64(f_Q16) + (int64(length) >> 3))
	} else {
		S0_Q16 = 1 << 16
		S1_Q16 = int32((int64(c_Q16) >> 1) + (1 << 16) + (int64(length) >> 4))
	}
	for k = 0; int64(k) < int64(length); k += 4 {
		px_win[k] = int16(SKP_SMULWB(int32((int64(S0_Q16)+int64(S1_Q16))>>1), int32(px[k])))
		px_win[int64(k)+1] = int16(SKP_SMULWB(S1_Q16, int32(px[int64(k)+1])))
		S0_Q16 = int32(int64(SKP_SMULWB(S1_Q16, c_Q16)) + (int64(S1_Q16) << 1) - int64(S0_Q16) + 1)
		if int64(S0_Q16) < (1 << 16) {
			S0_Q16 = S0_Q16
		} else {
			S0_Q16 = 1 << 16
		}
		px_win[int64(k)+2] = int16(SKP_SMULWB(int32((int64(S0_Q16)+int64(S1_Q16))>>1), int32(px[int64(k)+2])))
		px_win[int64(k)+3] = int16(SKP_SMULWB(S0_Q16, int32(px[int64(k)+3])))
		S1_Q16 = int32(int64(SKP_SMULWB(S0_Q16, c_Q16)) + (int64(S0_Q16) << 1) - int64(S1_Q16))
		if int64(S1_Q16) < (1 << 16) {
			S1_Q16 = S1_Q16
		} else {
			S1_Q16 = 1 << 16
		}
	}
}
