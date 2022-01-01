package silk

import "unsafe"

func SKP_Silk_init_decoder(psDec *SKP_Silk_decoder_state) int32 {
	memset(unsafe.Pointer(psDec), 0, size_t(unsafe.Sizeof(SKP_Silk_decoder_state{})))
	SKP_Silk_decoder_set_fs(psDec, 24)
	psDec.First_frame_after_reset = 1
	psDec.Prev_inv_gain_Q16 = 0x10000
	SKP_Silk_CNG_Reset(psDec)
	SKP_Silk_PLC_Reset(psDec)
	return 0
}
