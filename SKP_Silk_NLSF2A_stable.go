package silk

func SKP_Silk_NLSF2A_stable(pAR_Q12 [16]int16, pNLSF [16]int32, LPC_order int32) {
	var (
		i           int32
		invGain_Q30 int32
	)
	SKP_Silk_NLSF2A(&pAR_Q12[0], &pNLSF[0], LPC_order)
	for i = 0; int64(i) < MAX_LPC_STABILIZE_ITERATIONS; i++ {
		if int64(SKP_Silk_LPC_inverse_pred_gain(&invGain_Q30, &pAR_Q12[0], LPC_order)) == 1 {
			SKP_Silk_bwexpander(&pAR_Q12[0], LPC_order, int32(0x10000-int64(SKP_SMULBB(int32(int64(i)+10), i))))
		} else {
			break
		}
	}
	if int64(i) == MAX_LPC_STABILIZE_ITERATIONS {
		for i = 0; int64(i) < int64(LPC_order); i++ {
			pAR_Q12[i] = 0
		}
	}
}
