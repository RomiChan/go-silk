package silk

import "unsafe"

func SKP_Silk_decode_frame(psDec *SKP_Silk_decoder_state, pOut []int16, pN *int16, pCode []uint8, nBytes int32, action int32, decBytes *int32) int32 {
	var (
		sDecCtrl SKP_Silk_decoder_control
		ret      int32
		Pulses   [480]int32
	)
	L := psDec.Frame_length
	sDecCtrl.LTP_scale_Q14 = 0
	SKP_assert(L > 0 && L <= (FRAME_LENGTH_MS*MAX_FS_KHZ))
	*decBytes = 0
	if int64(action) == 0 {
		fs_Khz_old := psDec.Fs_kHz
		if int64(psDec.NFramesDecoded) == 0 {
			SKP_Silk_range_dec_init(&psDec.SRC, pCode, nBytes)
		}
		SKP_Silk_decode_parameters(psDec, &sDecCtrl, Pulses[:], 1)
		if psDec.SRC.Error != 0 {
			psDec.NBytesLeft = 0
			action = 1
			SKP_Silk_decoder_set_fs(psDec, fs_Khz_old)
			*decBytes = psDec.SRC.BufferLength
			if int64(psDec.SRC.Error) == -8 {
				ret = -11
			} else {
				ret = -12
			}
		} else {
			*decBytes = psDec.SRC.BufferLength - psDec.NBytesLeft
			psDec.NFramesDecoded++
			L = psDec.Frame_length
			SKP_Silk_decode_core(psDec, &sDecCtrl, pOut, Pulses)
			SKP_Silk_PLC(psDec, &sDecCtrl, pOut, L, action)
			psDec.LossCnt = 0
			psDec.Prev_sigtype = sDecCtrl.Sigtype
			psDec.First_frame_after_reset = 0
		}
	}
	if action == 1 {
		SKP_Silk_PLC(psDec, &sDecCtrl, pOut, L, action)
	}
	memcpy(unsafe.Pointer(&psDec.OutBuf[0]), unsafe.Pointer(&pOut[0]), uintptr(L)*unsafe.Sizeof(int16(0)))
	SKP_Silk_PLC_glue_frames(psDec, &sDecCtrl, pOut, L)
	SKP_Silk_CNG(psDec, &sDecCtrl, pOut, L)
	SKP_assert(psDec.Fs_kHz == 12 && (L%3) == 0 || psDec.Fs_kHz != 12 && (L%2) == 0)
	SKP_Silk_biquad(pOut, ([3]int16)(psDec.HP_B), ([2]int16)(psDec.HP_A), psDec.HPState[:], pOut, L)
	*pN = int16(L)
	psDec.LagPrev = sDecCtrl.PitchL[NB_SUBFR-1]
	return ret
}
