package silk

func SKP_Silk_gains_quant(ind []int32, gain_Q16 []int32, prev_ind *int32, conditional int32) {
	for k := 0; k < NB_SUBFR; k++ {
		ind[k] = SKP_SMULWB(((N_LEVELS_QGAIN-1)*0x10000)/(((MAX_QGAIN_DB-MIN_QGAIN_DB)*128)/6), int32(int64(SKP_Silk_lin2log(gain_Q16[k]))-((MIN_QGAIN_DB*128)/6+16*128)))
		if int64(ind[k]) < int64(*prev_ind) {
			ind[k]++
		}
		if k == 0 && int64(conditional) == 0 {
			ind[k] = SKP_LIMIT_int(ind[k], 0, N_LEVELS_QGAIN-1)
			ind[k] = SKP_max_int(ind[k], int32(int64(*prev_ind)+(-4)))
			*prev_ind = ind[k]
		} else {
			ind[k] = SKP_LIMIT_int(int32(int64(ind[k])-int64(*prev_ind)), -4, MAX_DELTA_GAIN_QUANT)
			*prev_ind += ind[k]
			ind[k] -= -4
		}
		gain_Q16[k] = SKP_Silk_log2lin(SKP_min_32(int32(int64(SKP_SMULWB(((((MAX_QGAIN_DB-MIN_QGAIN_DB)*128)/6)*0x10000)/(N_LEVELS_QGAIN-1), *prev_ind))+((MIN_QGAIN_DB*128)/6+16*128)), 3967))
	}
}

func SKP_Silk_gains_dequant(gain_Q16 [4]int32, ind [4]int32, prev_ind *int32, conditional int32) {
	for k := 0; k < NB_SUBFR; k++ {
		if k == 0 && int64(conditional) == 0 {
			*prev_ind = ind[k]
		} else {
			*prev_ind += int32(int64(ind[k]) + (-4))
		}
		gain_Q16[k] = SKP_Silk_log2lin(SKP_min_32(int32(int64(SKP_SMULWB(((((MAX_QGAIN_DB-MIN_QGAIN_DB)*128)/6)*0x10000)/(N_LEVELS_QGAIN-1), *prev_ind))+((MIN_QGAIN_DB*128)/6+16*128)), 3967))
	}
}
