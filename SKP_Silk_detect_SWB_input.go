package silk

func SKP_Silk_detect_SWB_input(psSWBdetect *SKP_Silk_detect_SWB_state, samplesIn []int16, nSamplesIn int32) {
	var (
		HP_8_kHz_len int32
		i            int32
		shift        int32
		in_HP_8_kHz  [480]int16
		energy_32    int32
	)
	HP_8_kHz_len = SKP_min_int(nSamplesIn, FRAME_LENGTH_MS*MAX_FS_KHZ)
	HP_8_kHz_len = SKP_max_int(HP_8_kHz_len, 0)
	SKP_Silk_biquad(samplesIn, SKP_Silk_SWB_detect_B_HP_Q13[0], SKP_Silk_SWB_detect_A_HP_Q13[0], psSWBdetect.S_HP_8_kHz[0][:], in_HP_8_kHz[:], HP_8_kHz_len)
	for i = 1; i < NB_SOS; i++ {
		SKP_Silk_biquad(in_HP_8_kHz[:], SKP_Silk_SWB_detect_B_HP_Q13[i], SKP_Silk_SWB_detect_A_HP_Q13[i], psSWBdetect.S_HP_8_kHz[i][:], in_HP_8_kHz[:], HP_8_kHz_len)
	}
	SKP_Silk_sum_sqr_shift(&energy_32, &shift, &in_HP_8_kHz[0], HP_8_kHz_len)
	if energy_32 > (SKP_SMULBB(HP_8_KHZ_THRES, HP_8_kHz_len) >> shift) {
		psSWBdetect.ConsecSmplsAboveThres += nSamplesIn
		if psSWBdetect.ConsecSmplsAboveThres > 480*15 {
			psSWBdetect.SWB_detected = 1
		}
	} else {
		psSWBdetect.ConsecSmplsAboveThres -= nSamplesIn
		if psSWBdetect.ConsecSmplsAboveThres > 0 {
			psSWBdetect.ConsecSmplsAboveThres = psSWBdetect.ConsecSmplsAboveThres
		} else {
			psSWBdetect.ConsecSmplsAboveThres = 0
		}
	}
	if psSWBdetect.ActiveSpeech_ms > WB_DETECT_ACTIVE_SPEECH_MS_THRES && psSWBdetect.SWB_detected == 0 {
		psSWBdetect.WB_detected = 1
	}
}
