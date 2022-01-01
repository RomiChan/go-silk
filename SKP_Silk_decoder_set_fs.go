package silk

import "unsafe"

func SKP_Silk_decoder_set_fs(psDec *SKP_Silk_decoder_state, fs_kHz int32) {
	if int64(psDec.Fs_kHz) != int64(fs_kHz) {
		psDec.Fs_kHz = fs_kHz
		psDec.Frame_length = SKP_SMULBB(FRAME_LENGTH_MS, fs_kHz)
		psDec.Subfr_length = SKP_SMULBB(FRAME_LENGTH_MS/NB_SUBFR, fs_kHz)
		if int64(psDec.Fs_kHz) == 8 {
			psDec.LPC_order = MIN_LPC_ORDER
			psDec.PsNLSF_CB[0] = &SKP_Silk_NLSF_CB0_10
			psDec.PsNLSF_CB[1] = &SKP_Silk_NLSF_CB1_10
		} else {
			psDec.LPC_order = MAX_LPC_ORDER
			psDec.PsNLSF_CB[0] = &SKP_Silk_NLSF_CB0_16
			psDec.PsNLSF_CB[1] = &SKP_Silk_NLSF_CB1_16
		}
		memset(unsafe.Pointer(&psDec.SLPC_Q14[0]), 0, MAX_LPC_ORDER*unsafe.Sizeof(int32(0)))
		memset(unsafe.Pointer(&psDec.OutBuf[0]), 0, FRAME_LENGTH_MS*MAX_FS_KHZ*unsafe.Sizeof(int16(0)))
		memset(unsafe.Pointer(&psDec.PrevNLSF_Q15[0]), 0, MAX_LPC_ORDER*unsafe.Sizeof(int32(0)))
		psDec.LagPrev = 100
		psDec.LastGainIndex = 1
		psDec.Prev_sigtype = 0
		psDec.First_frame_after_reset = 1
		if int64(fs_kHz) == 24 {
			psDec.HP_A = &SKP_Silk_Dec_A_HP_24[0]
			psDec.HP_B = &SKP_Silk_Dec_B_HP_24[0]
		} else if int64(fs_kHz) == 16 {
			psDec.HP_A = &SKP_Silk_Dec_A_HP_16[0]
			psDec.HP_B = &SKP_Silk_Dec_B_HP_16[0]
		} else if int64(fs_kHz) == 12 {
			psDec.HP_A = &SKP_Silk_Dec_A_HP_12[0]
			psDec.HP_B = &SKP_Silk_Dec_B_HP_12[0]
		} else if int64(fs_kHz) == 8 {
			psDec.HP_A = &SKP_Silk_Dec_A_HP_8[0]
			psDec.HP_B = &SKP_Silk_Dec_B_HP_8[0]
		}
	}
}
