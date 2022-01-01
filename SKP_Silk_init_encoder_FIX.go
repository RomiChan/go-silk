package silk

import "unsafe"

func SKP_Silk_init_encoder_FIX(psEnc *SKP_Silk_encoder_state_FIX) int32 {
	var ret int32 = 0
	memset(unsafe.Pointer(psEnc), 0, unsafe.Sizeof(SKP_Silk_encoder_state_FIX{}))
	psEnc.Variable_HP_smth1_Q15 = 0x3108C
	psEnc.Variable_HP_smth2_Q15 = 0x3108C
	psEnc.SCmn.First_frame_after_reset = 1
	ret += SKP_Silk_VAD_Init(&psEnc.SCmn.SVAD)
	psEnc.SCmn.SNSQ.Prev_inv_gain_Q16 = 0x10000
	psEnc.SCmn.SNSQ_LBRR.Prev_inv_gain_Q16 = 0x10000
	return ret
}
