package silk

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func SKP_Silk_process_NLSFs_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX, pNLSF_Q15 *int32) {
	var (
		doInterpolate        int32
		pNLSFW_Q6            [16]int32
		NLSF_mu_Q15          int32
		NLSF_mu_fluc_red_Q16 int32
		i_sqr_Q15            int32
		psNLSF_CB            *SKP_Silk_NLSF_CB_struct
		pNLSF0_temp_Q15      [16]int32
		pNLSFW0_temp_Q6      [16]int32
		i                    int32
	)
	if int64(psEncCtrl.SCmn.Sigtype) == SIG_TYPE_VOICED {
		NLSF_mu_Q15 = SKP_SMLAWB(66, -8388, psEnc.Speech_activity_Q8)
		NLSF_mu_fluc_red_Q16 = SKP_SMLAWB(6554, -838848, psEnc.Speech_activity_Q8)
	} else {
		NLSF_mu_Q15 = SKP_SMLAWB(164, -33554, psEnc.Speech_activity_Q8)
		NLSF_mu_fluc_red_Q16 = SKP_SMLAWB(0x3333, -1677696, int32(int64(psEnc.Speech_activity_Q8)+int64(psEncCtrl.Sparseness_Q8)))
	}
	if int64(NLSF_mu_Q15) > 1 {
		NLSF_mu_Q15 = NLSF_mu_Q15
	} else {
		NLSF_mu_Q15 = 1
	}
	SKP_Silk_NLSF_VQ_weights_laroia(&pNLSFW_Q6[0], pNLSF_Q15, psEnc.SCmn.PredictLPCOrder)
	doInterpolate = libc.BoolToInt(int64(psEnc.SCmn.UseInterpolatedNLSFs) == 1 && int64(psEncCtrl.SCmn.NLSFInterpCoef_Q2) < (1<<2))
	if int64(doInterpolate) != 0 {
		SKP_Silk_interpolate(pNLSF0_temp_Q15, psEnc.SPred.Prev_NLSFq_Q15, ([16]int32)(pNLSF_Q15), psEncCtrl.SCmn.NLSFInterpCoef_Q2, psEnc.SCmn.PredictLPCOrder)
		SKP_Silk_NLSF_VQ_weights_laroia(&pNLSFW0_temp_Q6[0], &pNLSF0_temp_Q15[0], psEnc.SCmn.PredictLPCOrder)
		i_sqr_Q15 = int32(int64(SKP_SMULBB(psEncCtrl.SCmn.NLSFInterpCoef_Q2, psEncCtrl.SCmn.NLSFInterpCoef_Q2)) << 11)
		for i = 0; int64(i) < int64(psEnc.SCmn.PredictLPCOrder); i++ {
			pNLSFW_Q6[i] = SKP_SMLAWB(int32(int64(pNLSFW_Q6[i])>>1), pNLSFW0_temp_Q6[i], i_sqr_Q15)
		}
	}
	psNLSF_CB = psEnc.SCmn.PsNLSF_CB[psEncCtrl.SCmn.Sigtype]
	SKP_Silk_NLSF_MSVQ_encode_FIX(&psEncCtrl.SCmn.NLSFIndices[0], pNLSF_Q15, psNLSF_CB, &psEnc.SPred.Prev_NLSFq_Q15[0], &pNLSFW_Q6[0], NLSF_mu_Q15, NLSF_mu_fluc_red_Q16, psEnc.SCmn.NLSF_MSVQ_Survivors, psEnc.SCmn.PredictLPCOrder, psEnc.SCmn.First_frame_after_reset)
	SKP_Silk_NLSF2A_stable(psEncCtrl.PredCoef_Q12[1], ([16]int32)(pNLSF_Q15), psEnc.SCmn.PredictLPCOrder)
	if int64(doInterpolate) != 0 {
		SKP_Silk_interpolate(pNLSF0_temp_Q15, psEnc.SPred.Prev_NLSFq_Q15, ([16]int32)(pNLSF_Q15), psEncCtrl.SCmn.NLSFInterpCoef_Q2, psEnc.SCmn.PredictLPCOrder)
		SKP_Silk_NLSF2A_stable(psEncCtrl.PredCoef_Q12[0], pNLSF0_temp_Q15, psEnc.SCmn.PredictLPCOrder)
	} else {
		memcpy(unsafe.Pointer(&(psEncCtrl.PredCoef_Q12[0])[0]), unsafe.Pointer(&(psEncCtrl.PredCoef_Q12[1])[0]), size_t(uintptr(psEnc.SCmn.PredictLPCOrder)*unsafe.Sizeof(int16(0))))
	}
}
