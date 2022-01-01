package silk

func SKP_Silk_setup_complexity(psEncC *SKP_Silk_encoder_state, Complexity int32) int32 {
	var ret int32 = SKP_SILK_NO_ERROR
	if LOW_COMPLEXITY_ONLY != 0 && int64(Complexity) != 0 {
		ret = -6
	}
	if int64(Complexity) == 0 || LOW_COMPLEXITY_ONLY != 0 {
		psEncC.Complexity = 0
		psEncC.PitchEstimationComplexity = SKP_Silk_PITCH_EST_MIN_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = SKP_FIX_CONST(0.8, 16)
		psEncC.PitchEstimationLPCOrder = 6
		psEncC.ShapingLPCOrder = 8
		psEncC.La_shape = int32(int64(psEncC.Fs_kHz) * 3)
		psEncC.NStatesDelayedDecision = 1
		psEncC.UseInterpolatedNLSFs = 0
		psEncC.LTPQuantLowComplexity = 1
		psEncC.NLSF_MSVQ_Survivors = MAX_NLSF_MSVQ_SURVIVORS_LC_MODE
		psEncC.Warping_Q16 = 0
	} else if int64(Complexity) == 1 {
		psEncC.Complexity = 1
		psEncC.PitchEstimationComplexity = SKP_Silk_PITCH_EST_MID_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = SKP_FIX_CONST(0.75, 16)
		psEncC.PitchEstimationLPCOrder = 12
		psEncC.ShapingLPCOrder = 12
		psEncC.La_shape = int32(int64(psEncC.Fs_kHz) * 5)
		psEncC.NStatesDelayedDecision = 2
		psEncC.UseInterpolatedNLSFs = 0
		psEncC.LTPQuantLowComplexity = 0
		psEncC.NLSF_MSVQ_Survivors = MAX_NLSF_MSVQ_SURVIVORS_MC_MODE
		psEncC.Warping_Q16 = int32(int64(psEncC.Fs_kHz) * int64(SKP_FIX_CONST(0.015, 16)))
	} else if int64(Complexity) == 2 {
		psEncC.Complexity = 2
		psEncC.PitchEstimationComplexity = SKP_Silk_PITCH_EST_MAX_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = SKP_FIX_CONST(0.7, 16)
		psEncC.PitchEstimationLPCOrder = 16
		psEncC.ShapingLPCOrder = 16
		psEncC.La_shape = int32(int64(psEncC.Fs_kHz) * 5)
		psEncC.NStatesDelayedDecision = MAX_DEL_DEC_STATES
		psEncC.UseInterpolatedNLSFs = 1
		psEncC.LTPQuantLowComplexity = 0
		psEncC.NLSF_MSVQ_Survivors = MAX_NLSF_MSVQ_SURVIVORS
		psEncC.Warping_Q16 = int32(int64(psEncC.Fs_kHz) * int64(SKP_FIX_CONST(0.015, 16)))
	} else {
		ret = -6
	}
	psEncC.PitchEstimationLPCOrder = SKP_min_int(psEncC.PitchEstimationLPCOrder, psEncC.PredictLPCOrder)
	psEncC.ShapeWinLength = int32(int64(psEncC.Fs_kHz)*5 + int64(psEncC.La_shape)*2)
	return ret
}
