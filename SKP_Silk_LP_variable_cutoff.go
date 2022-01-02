package silk

import "unsafe"

func SKP_Silk_LP_interpolate_filter_taps(B_Q28 []int32, A_Q28 []int32, ind int32, fac_Q16 int32) {
	var (
		nb int32
		na int32
	)
	if ind < TRANSITION_INT_NUM-1 {
		if fac_Q16 > 0 {
			if int64(fac_Q16) == int64(SKP_SAT16(fac_Q16)) {
				for nb = 0; nb < TRANSITION_NB; nb++ {
					B_Q28[nb] = SKP_SMLAWB(SKP_Silk_Transition_LP_B_Q28[ind][nb], SKP_Silk_Transition_LP_B_Q28[ind+1][nb]-SKP_Silk_Transition_LP_B_Q28[ind][nb], fac_Q16)
				}
				for na = 0; na < TRANSITION_NA; na++ {
					A_Q28[na] = SKP_SMLAWB(SKP_Silk_Transition_LP_A_Q28[ind][na], SKP_Silk_Transition_LP_A_Q28[ind+1][na]-SKP_Silk_Transition_LP_A_Q28[ind][na], fac_Q16)
				}
			} else if fac_Q16 == (1 << 15) {
				for nb = 0; nb < TRANSITION_NB; nb++ {
					B_Q28[nb] = (SKP_Silk_Transition_LP_B_Q28[ind][nb] + SKP_Silk_Transition_LP_B_Q28[ind+1][nb]) >> 1
				}
				for na = 0; na < TRANSITION_NA; na++ {
					A_Q28[na] = (SKP_Silk_Transition_LP_A_Q28[ind][na] + SKP_Silk_Transition_LP_A_Q28[ind+1][na]) >> 1
				}
			} else {
				SKP_assert(int64((1<<16)-fac_Q16) == int64(SKP_SAT16((1<<16)-fac_Q16)))
				for nb = 0; nb < TRANSITION_NB; nb++ {
					B_Q28[nb] = SKP_SMLAWB(SKP_Silk_Transition_LP_B_Q28[ind+1][nb], SKP_Silk_Transition_LP_B_Q28[ind][nb]-SKP_Silk_Transition_LP_B_Q28[ind+1][nb], (1<<16)-fac_Q16)
				}
				for na = 0; na < TRANSITION_NA; na++ {
					A_Q28[na] = SKP_SMLAWB(SKP_Silk_Transition_LP_A_Q28[ind+1][na], SKP_Silk_Transition_LP_A_Q28[ind][na]-SKP_Silk_Transition_LP_A_Q28[ind+1][na], (1<<16)-fac_Q16)
				}
			}
		} else {
			memcpy(unsafe.Pointer(&B_Q28[0]), unsafe.Pointer(&(SKP_Silk_Transition_LP_B_Q28[ind])[0]), size_t(TRANSITION_NB*unsafe.Sizeof(int32(0))))
			memcpy(unsafe.Pointer(&A_Q28[0]), unsafe.Pointer(&(SKP_Silk_Transition_LP_A_Q28[ind])[0]), size_t(TRANSITION_NA*unsafe.Sizeof(int32(0))))
		}
	} else {
		memcpy(unsafe.Pointer(&B_Q28[0]), unsafe.Pointer(&(SKP_Silk_Transition_LP_B_Q28[TRANSITION_INT_NUM-1])[0]), size_t(TRANSITION_NB*unsafe.Sizeof(int32(0))))
		memcpy(unsafe.Pointer(&A_Q28[0]), unsafe.Pointer(&(SKP_Silk_Transition_LP_A_Q28[TRANSITION_INT_NUM-1])[0]), size_t(TRANSITION_NA*unsafe.Sizeof(int32(0))))
	}
}
func SKP_Silk_LP_variable_cutoff(psLP *SKP_Silk_LP_state, out *int16, in *int16, frame_length int32) {
	var (
		B_Q28   [3]int32
		A_Q28   [2]int32
		fac_Q16 int32 = 0
		ind     int32 = 0
	)
	SKP_assert(psLP.Transition_frame_no >= 0)
	SKP_assert(psLP.Transition_frame_no <= (TRANSITION_TIME_DOWN_MS/FRAME_LENGTH_MS) && psLP.Mode == 0 || psLP.Transition_frame_no <= (TRANSITION_TIME_UP_MS/FRAME_LENGTH_MS) && psLP.Mode == 1)
	if psLP.Transition_frame_no > 0 {
		if psLP.Mode == 0 {
			if psLP.Transition_frame_no < (TRANSITION_TIME_DOWN_MS / FRAME_LENGTH_MS) {
				fac_Q16 = psLP.Transition_frame_no << (16 - 5)
				ind = fac_Q16 >> 16
				fac_Q16 -= ind << 16
				SKP_assert(ind >= 0)
				SKP_assert(ind < TRANSITION_INT_NUM)
				SKP_Silk_LP_interpolate_filter_taps(B_Q28[:], A_Q28[:], ind, fac_Q16)
				psLP.Transition_frame_no++
			} else {
				SKP_assert(psLP.Transition_frame_no == (TRANSITION_TIME_DOWN_MS / FRAME_LENGTH_MS))
				SKP_Silk_LP_interpolate_filter_taps(B_Q28[:], A_Q28[:], TRANSITION_INT_NUM-1, 0)
			}
		} else {
			SKP_assert(psLP.Mode == 1)
			if psLP.Transition_frame_no < (TRANSITION_TIME_UP_MS / FRAME_LENGTH_MS) {
				fac_Q16 = ((TRANSITION_TIME_UP_MS / FRAME_LENGTH_MS) - psLP.Transition_frame_no) << (16 - 6)
				ind = fac_Q16 >> 16
				fac_Q16 -= ind << 16
				SKP_assert(ind >= 0)
				SKP_assert(ind < TRANSITION_INT_NUM)
				SKP_Silk_LP_interpolate_filter_taps(B_Q28[:], A_Q28[:], ind, fac_Q16)
				psLP.Transition_frame_no++
			} else {
				SKP_assert(psLP.Transition_frame_no == (TRANSITION_TIME_UP_MS / FRAME_LENGTH_MS))
				SKP_Silk_LP_interpolate_filter_taps(B_Q28[:], A_Q28[:], 0, 0)
			}
		}
	}
	if psLP.Transition_frame_no > 0 {
		SKP_assert(TRANSITION_NB == 3 && TRANSITION_NA == 2)
		SKP_Silk_biquad_alt(([]int16)(in), B_Q28, A_Q28, psLP.In_LP_State[:], ([]int16)(out), frame_length)
	} else {
		memcpy(unsafe.Pointer(out), unsafe.Pointer(in), size_t(uintptr(frame_length)*unsafe.Sizeof(int16(0))))
	}
}
