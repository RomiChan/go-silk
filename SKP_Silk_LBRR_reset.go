package silk

func SKP_Silk_LBRR_reset(psEncC *SKP_Silk_encoder_state) {
	var i int32
	for i = 0; i < MAX_LBRR_DELAY; i++ {
		psEncC.LBRR_buffer[i].Usage = SKP_SILK_NO_LBRR
	}
}
