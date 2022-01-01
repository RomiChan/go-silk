package silk

import "unsafe"

func SKP_Silk_CNG_exc(residual []int16, exc_buf_Q10 []int32, Gain_Q16 int32, rand_seed *int32) {
	exc_mask := CNG_BUF_MASK_MAX
	for exc_mask > len(residual) {
		exc_mask = exc_mask >> 1
	}
	seed := *rand_seed
	for i := 0; i < len(residual); i++ {
		seed = int32(int64(uint32(int32(int64(uint32(seed))*0xBB38435))) + 0x3619636B)
		idx := int32((int64(seed) >> 24) & int64(exc_mask))
		residual[i] = SKP_SAT16(SKP_RSHIFT_ROUND(SKP_SMULWW(exc_buf_Q10[idx], Gain_Q16), 10))
	}
	*rand_seed = seed
}

func SKP_Silk_CNG_Reset(psDec *SKP_Silk_decoder_state) {
	NLSF_step_Q15 := int32(SKP_int16_MAX / (int64(psDec.LPC_order) + 1))
	NLSF_acc_Q15 := int32(0)
	for i := int32(0); i < psDec.LPC_order; i++ {
		NLSF_acc_Q15 += NLSF_step_Q15
		psDec.SCNG.CNG_smth_NLSF_Q15[i] = NLSF_acc_Q15
	}
	psDec.SCNG.CNG_smth_Gain_Q16 = 0
	psDec.SCNG.Rand_seed = 0x307880
}

func SKP_Silk_CNG(psDec *SKP_Silk_decoder_state, psDecCtrl *SKP_Silk_decoder_control, signal []int16, length int32) {
	var (
		LPC_buf [16]int16
		CNG_sig = make([]int16, length)
	)
	psCNG := &psDec.SCNG
	if int64(psDec.Fs_kHz) != int64(psCNG.Fs_kHz) {
		SKP_Silk_CNG_Reset(psDec)
		psCNG.Fs_kHz = psDec.Fs_kHz
	}
	if int64(psDec.LossCnt) == 0 && int64(psDec.VadFlag) == NO_VOICE_ACTIVITY {
		for i := int32(0); i < psDec.LPC_order; i++ {
			psCNG.CNG_smth_NLSF_Q15[i] += SKP_SMULWB(int32(int64(psDec.PrevNLSF_Q15[i])-int64(psCNG.CNG_smth_NLSF_Q15[i])), CNG_NLSF_SMTH_Q16)
		}
		max_Gain_Q16 := 0
		subfr := 0
		for i := 0; i < NB_SUBFR; i++ {
			if int64(psDecCtrl.Gains_Q16[i]) > int64(max_Gain_Q16) {
				max_Gain_Q16 = int(psDecCtrl.Gains_Q16[i])
				subfr = i
			}
		}
		memmove(unsafe.Pointer(&psCNG.CNG_exc_buf_Q10[psDec.Subfr_length]), unsafe.Pointer(&psCNG.CNG_exc_buf_Q10[0]), size_t((NB_SUBFR-1)*int64(psDec.Subfr_length)*int64(unsafe.Sizeof(int32(0)))))
		memcpy(unsafe.Pointer(&psCNG.CNG_exc_buf_Q10[0]), unsafe.Pointer(&psDec.Exc_Q10[int64(subfr)*int64(psDec.Subfr_length)]), uintptr(psDec.Subfr_length)*unsafe.Sizeof(int32(0)))
		for i := 0; i < NB_SUBFR; i++ {
			psCNG.CNG_smth_Gain_Q16 += SKP_SMULWB(int32(int64(psDecCtrl.Gains_Q16[i])-int64(psCNG.CNG_smth_Gain_Q16)), CNG_GAIN_SMTH_Q16)
		}
	}
	if int64(psDec.LossCnt) != 0 {
		SKP_Silk_CNG_exc(CNG_sig[:length], psCNG.CNG_exc_buf_Q10[:], psCNG.CNG_smth_Gain_Q16, &psCNG.Rand_seed)
		SKP_Silk_NLSF2A_stable(LPC_buf, psCNG.CNG_smth_NLSF_Q15, psDec.LPC_order)
		Gain_Q26 := int32(1 << 26)
		if int64(psDec.LPC_order) == 16 {
			SKP_Silk_LPC_synthesis_order16(&CNG_sig[0], &LPC_buf[0], Gain_Q26, &psCNG.CNG_synth_state[0], &CNG_sig[0], length)
		} else {
			SKP_Silk_LPC_synthesis_filter(&CNG_sig[0], &LPC_buf[0], Gain_Q26, &psCNG.CNG_synth_state[0], &CNG_sig[0], length, psDec.LPC_order)
		}
		for i, v := range CNG_sig {
			tmp_32 := int32(int64(signal[i]) + int64(v))
			signal[i] = SKP_SAT16(tmp_32)
		}
	} else {
		memset(unsafe.Pointer(&psCNG.CNG_synth_state[0]), 0, uintptr(psDec.LPC_order)*unsafe.Sizeof(int32(0)))
	}
}
