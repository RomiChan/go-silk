package silk

// reviewed by wdvxdr1123 2022-02-08

// SKP_Silk_LBRR_reset Resets LBRR buffer, used if packet size changes
func SKP_Silk_LBRR_reset(psEncC *SKP_Silk_encoder_state) {
	for i := 0; i < MAX_LBRR_DELAY; i++ {
		psEncC.LBRR_buffer[i].Usage = SKP_SILK_NO_LBRR
	}
}
