package silk

const NB_THRESHOLDS = 11

var LTPScaleThresholds_Q15 = [11]int16{0x7999, 0x6666, 0x4000, 0x3333, 9830, 6554, 4915, 3276, 2621, 2458, 0}

func SKP_Silk_LTP_scale_ctrl_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX) {
	var (
		round_loss        int32
		frames_per_packet int32
		g_out_Q5          int32
		g_limit_Q15       int32
		thrld1_Q15        int32
		thrld2_Q15        int32
	)
	psEnc.HPLTPredCodGain_Q7 = int32(SKP_max_int(int(psEncCtrl.LTPredCodGain_Q7-psEnc.PrevLTPredCodGain_Q7), 0)) + SKP_RSHIFT_ROUND(psEnc.HPLTPredCodGain_Q7, 1)
	psEnc.PrevLTPredCodGain_Q7 = psEncCtrl.LTPredCodGain_Q7
	g_out_Q5 = SKP_RSHIFT_ROUND((psEncCtrl.LTPredCodGain_Q7>>1)+(psEnc.HPLTPredCodGain_Q7>>1), 3)
	g_limit_Q15 = SKP_Silk_sigm_Q15(g_out_Q5 - (3 << 5))
	psEncCtrl.SCmn.LTP_scaleIndex = 0
	round_loss = psEnc.SCmn.PacketLoss_perc
	if psEnc.SCmn.NFramesInPayloadBuf == 0 {
		frames_per_packet = psEnc.SCmn.PacketSize_ms / FRAME_LENGTH_MS
		round_loss += frames_per_packet - 1
		thrld1_Q15 = int32(LTPScaleThresholds_Q15[SKP_min_int(int(round_loss), NB_THRESHOLDS-1)])
		thrld2_Q15 = int32(LTPScaleThresholds_Q15[SKP_min_int(int(round_loss+1), NB_THRESHOLDS-1)])
		if g_limit_Q15 > thrld1_Q15 {
			psEncCtrl.SCmn.LTP_scaleIndex = 2
		} else if g_limit_Q15 > thrld2_Q15 {
			psEncCtrl.SCmn.LTP_scaleIndex = 1
		}
	}
	psEncCtrl.LTP_scale_Q14 = int32(SKP_Silk_LTPScales_table_Q14[psEncCtrl.SCmn.LTP_scaleIndex])
}
