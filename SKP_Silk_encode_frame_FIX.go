package silk

import "unsafe"

func SKP_Silk_encode_frame_FIX(psEnc *SKP_Silk_encoder_state_FIX, pCode *uint8, pnBytesOut *int16, pIn *int16) int32 {
	var (
		sEncCtrl             SKP_Silk_encoder_control_FIX
		nBytes               int32
		ret                  int32 = 0
		x_frame              *int16
		res_pitch_frame      *int16
		xfw                  [480]int16
		pIn_HP               [480]int16
		res_pitch            [1008]int16
		LBRR_idx             int32
		frame_terminator     int32
		SNR_dB_Q7            int32
		FrameTermination_CDF *uint16
		LBRRpayload          [1024]uint8
		nBytesLBRR           int16
	)
	sEncCtrl.SCmn.Seed = func() int32 {
		p := &psEnc.SCmn.FrameCounter
		x := *p
		*p++
		return x
	}() & 3
	x_frame = &psEnc.X_buf[psEnc.SCmn.Frame_length]
	res_pitch_frame = &res_pitch[psEnc.SCmn.Frame_length]
	ret = SKP_Silk_VAD_GetSA_Q8(&psEnc.SCmn.SVAD, &psEnc.Speech_activity_Q8, &SNR_dB_Q7, sEncCtrl.Input_quality_bands_Q15, &sEncCtrl.Input_tilt_Q15, ([]int16)(pIn), psEnc.SCmn.Frame_length)
	SKP_Silk_HP_variable_cutoff_FIX(psEnc, &sEncCtrl, &pIn_HP[0], pIn)
	SKP_Silk_LP_variable_cutoff(&psEnc.SCmn.SLP, (*int16)(unsafe.Add(unsafe.Pointer(x_frame), unsafe.Sizeof(int16(0))*uintptr(LA_SHAPE_MS*psEnc.SCmn.Fs_kHz))), &pIn_HP[0], psEnc.SCmn.Frame_length)
	SKP_Silk_find_pitch_lags_FIX(psEnc, &sEncCtrl, res_pitch[:], ([]int16)(x_frame))
	SKP_Silk_noise_shape_analysis_FIX(psEnc, &sEncCtrl, res_pitch_frame, x_frame)
	SKP_Silk_prefilter_FIX(psEnc, &sEncCtrl, xfw[:], ([]int16)(x_frame))
	SKP_Silk_find_pred_coefs_FIX(psEnc, &sEncCtrl, res_pitch[:])
	SKP_Silk_process_gains_FIX(psEnc, &sEncCtrl)
	nBytesLBRR = MAX_ARITHM_BYTES
	SKP_Silk_LBRR_encode_FIX(psEnc, &sEncCtrl, &LBRRpayload[0], &nBytesLBRR, xfw[:])
	if psEnc.SCmn.NStatesDelayedDecision > 1 || psEnc.SCmn.Warping_Q16 > 0 {
		SKP_Silk_NSQ_del_dec(&psEnc.SCmn, &sEncCtrl.SCmn, &psEnc.SCmn.SNSQ, xfw[:], psEnc.SCmn.Q[:], sEncCtrl.SCmn.NLSFInterpCoef_Q2, ([32]int16)(sEncCtrl.PredCoef_Q12[0]), sEncCtrl.LTPCoef_Q14, sEncCtrl.AR2_Q13, sEncCtrl.HarmShapeGain_Q14, sEncCtrl.Tilt_Q14, sEncCtrl.LF_shp_Q14, sEncCtrl.Gains_Q16, sEncCtrl.Lambda_Q10, sEncCtrl.LTP_scale_Q14)
	} else {
		SKP_Silk_NSQ(&psEnc.SCmn, &sEncCtrl.SCmn, &psEnc.SCmn.SNSQ, xfw[:], psEnc.SCmn.Q[:], sEncCtrl.SCmn.NLSFInterpCoef_Q2, ([32]int16)(sEncCtrl.PredCoef_Q12[0]), sEncCtrl.LTPCoef_Q14, sEncCtrl.AR2_Q13, sEncCtrl.HarmShapeGain_Q14, sEncCtrl.Tilt_Q14, sEncCtrl.LF_shp_Q14, sEncCtrl.Gains_Q16, sEncCtrl.Lambda_Q10, sEncCtrl.LTP_scale_Q14)
	}
	if psEnc.Speech_activity_Q8 < SKP_FIX_CONST(SPEECH_ACTIVITY_DTX_THRES, 8) {
		psEnc.SCmn.VadFlag = NO_VOICE_ACTIVITY
		psEnc.SCmn.NoSpeechCounter++
		if psEnc.SCmn.NoSpeechCounter > NO_SPEECH_FRAMES_BEFORE_DTX {
			psEnc.SCmn.InDTX = 1
		}
		if psEnc.SCmn.NoSpeechCounter > MAX_CONSECUTIVE_DTX+NO_SPEECH_FRAMES_BEFORE_DTX {
			psEnc.SCmn.NoSpeechCounter = NO_SPEECH_FRAMES_BEFORE_DTX
			psEnc.SCmn.InDTX = 0
		}
	} else {
		psEnc.SCmn.NoSpeechCounter = 0
		psEnc.SCmn.InDTX = 0
		psEnc.SCmn.VadFlag = VOICE_ACTIVITY
	}
	if psEnc.SCmn.NFramesInPayloadBuf == 0 {
		SKP_Silk_range_enc_init(&psEnc.SCmn.SRC)
		psEnc.SCmn.NBytesInPayloadBuf = 0
	}
	SKP_Silk_encode_parameters(&psEnc.SCmn, &sEncCtrl.SCmn, &psEnc.SCmn.SRC, &psEnc.SCmn.Q[0])
	FrameTermination_CDF = &SKP_Silk_FrameTermination_CDF[0]
	memmove(unsafe.Pointer(&psEnc.X_buf[0]), unsafe.Pointer(&psEnc.X_buf[psEnc.SCmn.Frame_length]), size_t(uintptr(psEnc.SCmn.Frame_length+LA_SHAPE_MS*psEnc.SCmn.Fs_kHz)*unsafe.Sizeof(int16(0))))
	psEnc.SCmn.Prev_sigtype = sEncCtrl.SCmn.Sigtype
	psEnc.SCmn.PrevLag = sEncCtrl.SCmn.PitchL[NB_SUBFR-1]
	psEnc.SCmn.First_frame_after_reset = 0
	if psEnc.SCmn.SRC.Error != 0 {
		psEnc.SCmn.NFramesInPayloadBuf = 0
	} else {
		psEnc.SCmn.NFramesInPayloadBuf++
	}
	if psEnc.SCmn.NFramesInPayloadBuf*FRAME_LENGTH_MS >= psEnc.SCmn.PacketSize_ms {
		LBRR_idx = (psEnc.SCmn.Oldest_LBRR_idx + 1) & LBRR_IDX_MASK
		frame_terminator = SKP_SILK_LAST_FRAME
		if psEnc.SCmn.LBRR_buffer[LBRR_idx].Usage == SKP_SILK_ADD_LBRR_TO_PLUS1 {
			frame_terminator = SKP_SILK_LBRR_VER1
		}
		if psEnc.SCmn.LBRR_buffer[psEnc.SCmn.Oldest_LBRR_idx].Usage == SKP_SILK_ADD_LBRR_TO_PLUS2 {
			frame_terminator = SKP_SILK_LBRR_VER2
			LBRR_idx = psEnc.SCmn.Oldest_LBRR_idx
		}
		SKP_Silk_range_encoder(&psEnc.SCmn.SRC, frame_terminator, ([]uint16)(FrameTermination_CDF))
		SKP_Silk_range_coder_get_length(&psEnc.SCmn.SRC, &nBytes)
		if int64(*pnBytesOut) >= int64(nBytes) {
			SKP_Silk_range_enc_wrap_up(&psEnc.SCmn.SRC)
			memcpy(unsafe.Pointer(pCode), unsafe.Pointer(&psEnc.SCmn.SRC.Buffer[0]), size_t(uintptr(nBytes)*unsafe.Sizeof(uint8(0))))
			if frame_terminator > SKP_SILK_MORE_FRAMES && int64(*pnBytesOut) >= int64(nBytes+psEnc.SCmn.LBRR_buffer[LBRR_idx].NBytes) {
				memcpy(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(pCode), nBytes))), unsafe.Pointer(&psEnc.SCmn.LBRR_buffer[LBRR_idx].Payload[0]), size_t(uintptr(psEnc.SCmn.LBRR_buffer[LBRR_idx].NBytes)*unsafe.Sizeof(uint8(0))))
				nBytes += psEnc.SCmn.LBRR_buffer[LBRR_idx].NBytes
			}
			*pnBytesOut = int16(nBytes)
			memcpy(unsafe.Pointer(&psEnc.SCmn.LBRR_buffer[psEnc.SCmn.Oldest_LBRR_idx].Payload[0]), unsafe.Pointer(&LBRRpayload[0]), size_t(uintptr(nBytesLBRR)*unsafe.Sizeof(uint8(0))))
			psEnc.SCmn.LBRR_buffer[psEnc.SCmn.Oldest_LBRR_idx].NBytes = int32(nBytesLBRR)
			psEnc.SCmn.LBRR_buffer[psEnc.SCmn.Oldest_LBRR_idx].Usage = sEncCtrl.SCmn.LBRR_usage
			psEnc.SCmn.Oldest_LBRR_idx = (psEnc.SCmn.Oldest_LBRR_idx + 1) & LBRR_IDX_MASK
		} else {
			*pnBytesOut = 0
			nBytes = 0
			ret = -4
		}
		psEnc.SCmn.NFramesInPayloadBuf = 0
	} else {
		*pnBytesOut = 0
		frame_terminator = SKP_SILK_MORE_FRAMES
		SKP_Silk_range_encoder(&psEnc.SCmn.SRC, frame_terminator, ([]uint16)(FrameTermination_CDF))
		SKP_Silk_range_coder_get_length(&psEnc.SCmn.SRC, &nBytes)
	}
	if psEnc.SCmn.SRC.Error != 0 {
		ret = -9
	}
	psEnc.BufferedInChannel_ms += ((nBytes - psEnc.SCmn.NBytesInPayloadBuf) * (8 * 1000)) / psEnc.SCmn.TargetRate_bps
	psEnc.BufferedInChannel_ms -= FRAME_LENGTH_MS
	psEnc.BufferedInChannel_ms = SKP_LIMIT_int(psEnc.BufferedInChannel_ms, 0, 100)
	psEnc.SCmn.NBytesInPayloadBuf = nBytes
	if psEnc.Speech_activity_Q8 > SKP_FIX_CONST(WB_DETECT_ACTIVE_SPEECH_LEVEL_THRES, 8) {
		psEnc.SCmn.SSWBdetect.ActiveSpeech_ms = SKP_ADD_POS_SAT32(psEnc.SCmn.SSWBdetect.ActiveSpeech_ms, FRAME_LENGTH_MS)
	}
	return ret
}
func SKP_Silk_LBRR_encode_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX, pCode *uint8, pnBytesOut *int16, xfw []int16) {
	var (
		TempGainsIndices     [4]int32
		frame_terminator     int32
		nBytes               int32
		nFramesInPayloadBuf  int32
		TempGains_Q16        [4]int32
		typeOffset           int32
		LTP_scaleIndex       int32
		Rate_only_parameters int32 = 0
	)
	SKP_Silk_LBRR_ctrl_FIX(psEnc, &psEncCtrl.SCmn)
	if psEnc.SCmn.LBRR_enabled != 0 {
		memcpy(unsafe.Pointer(&TempGainsIndices[0]), unsafe.Pointer(&psEncCtrl.SCmn.GainsIndices[0]), size_t(NB_SUBFR*unsafe.Sizeof(int32(0))))
		memcpy(unsafe.Pointer(&TempGains_Q16[0]), unsafe.Pointer(&psEncCtrl.Gains_Q16[0]), size_t(NB_SUBFR*unsafe.Sizeof(int32(0))))
		typeOffset = psEnc.SCmn.TypeOffsetPrev
		LTP_scaleIndex = psEncCtrl.SCmn.LTP_scaleIndex
		if psEnc.SCmn.Fs_kHz == 8 {
			Rate_only_parameters = 13500
		} else if psEnc.SCmn.Fs_kHz == 12 {
			Rate_only_parameters = 15500
		} else if psEnc.SCmn.Fs_kHz == 16 {
			Rate_only_parameters = 17500
		} else if psEnc.SCmn.Fs_kHz == 24 {
			Rate_only_parameters = 19500
		} else {
		}
		if psEnc.SCmn.Complexity > 0 && psEnc.SCmn.TargetRate_bps > Rate_only_parameters {
			if psEnc.SCmn.NFramesInPayloadBuf == 0 {
				memcpy(unsafe.Pointer(&psEnc.SCmn.SNSQ_LBRR), unsafe.Pointer(&psEnc.SCmn.SNSQ), size_t(unsafe.Sizeof(SKP_Silk_nsq_state{})))
				psEnc.SCmn.LBRRprevLastGainIndex = psEnc.SShape.LastGainIndex
				psEncCtrl.SCmn.GainsIndices[0] = psEncCtrl.SCmn.GainsIndices[0] + psEnc.SCmn.LBRR_GainIncreases
				psEncCtrl.SCmn.GainsIndices[0] = SKP_LIMIT_int(psEncCtrl.SCmn.GainsIndices[0], 0, N_LEVELS_QGAIN-1)
			}
			SKP_Silk_gains_dequant(psEncCtrl.Gains_Q16, psEncCtrl.SCmn.GainsIndices, &psEnc.SCmn.LBRRprevLastGainIndex, psEnc.SCmn.NFramesInPayloadBuf)
			if psEnc.SCmn.NStatesDelayedDecision > 1 || psEnc.SCmn.Warping_Q16 > 0 {
				SKP_Silk_NSQ_del_dec(&psEnc.SCmn, &psEncCtrl.SCmn, &psEnc.SCmn.SNSQ_LBRR, xfw, psEnc.SCmn.Q_LBRR[:], psEncCtrl.SCmn.NLSFInterpCoef_Q2, ([32]int16)(psEncCtrl.PredCoef_Q12[0]), psEncCtrl.LTPCoef_Q14, psEncCtrl.AR2_Q13, psEncCtrl.HarmShapeGain_Q14, psEncCtrl.Tilt_Q14, psEncCtrl.LF_shp_Q14, psEncCtrl.Gains_Q16, psEncCtrl.Lambda_Q10, psEncCtrl.LTP_scale_Q14)
			} else {
				SKP_Silk_NSQ(&psEnc.SCmn, &psEncCtrl.SCmn, &psEnc.SCmn.SNSQ_LBRR, xfw, psEnc.SCmn.Q_LBRR[:], psEncCtrl.SCmn.NLSFInterpCoef_Q2, ([32]int16)(psEncCtrl.PredCoef_Q12[0]), psEncCtrl.LTPCoef_Q14, psEncCtrl.AR2_Q13, psEncCtrl.HarmShapeGain_Q14, psEncCtrl.Tilt_Q14, psEncCtrl.LF_shp_Q14, psEncCtrl.Gains_Q16, psEncCtrl.Lambda_Q10, psEncCtrl.LTP_scale_Q14)
			}
		} else {
			memset(unsafe.Pointer(&psEnc.SCmn.Q_LBRR[0]), 0, size_t(uintptr(psEnc.SCmn.Frame_length)*unsafe.Sizeof(int8(0))))
			psEncCtrl.SCmn.LTP_scaleIndex = 0
		}
		if psEnc.SCmn.NFramesInPayloadBuf == 0 {
			SKP_Silk_range_enc_init(&psEnc.SCmn.SRC_LBRR)
			psEnc.SCmn.NBytesInPayloadBuf = 0
		}
		SKP_Silk_encode_parameters(&psEnc.SCmn, &psEncCtrl.SCmn, &psEnc.SCmn.SRC_LBRR, &psEnc.SCmn.Q_LBRR[0])
		if psEnc.SCmn.SRC_LBRR.Error != 0 {
			nFramesInPayloadBuf = 0
		} else {
			nFramesInPayloadBuf = psEnc.SCmn.NFramesInPayloadBuf + 1
		}
		if SKP_SMULBB(nFramesInPayloadBuf, FRAME_LENGTH_MS) >= psEnc.SCmn.PacketSize_ms {
			frame_terminator = SKP_SILK_LAST_FRAME
			SKP_Silk_range_encoder(&psEnc.SCmn.SRC_LBRR, frame_terminator, SKP_Silk_FrameTermination_CDF[:])
			SKP_Silk_range_coder_get_length(&psEnc.SCmn.SRC_LBRR, &nBytes)
			if int64(*pnBytesOut) >= int64(nBytes) {
				SKP_Silk_range_enc_wrap_up(&psEnc.SCmn.SRC_LBRR)
				memcpy(unsafe.Pointer(pCode), unsafe.Pointer(&psEnc.SCmn.SRC_LBRR.Buffer[0]), size_t(uintptr(nBytes)*unsafe.Sizeof(uint8(0))))
				*pnBytesOut = int16(nBytes)
			} else {
				*pnBytesOut = 0
			}
		} else {
			*pnBytesOut = 0
			frame_terminator = SKP_SILK_MORE_FRAMES
			SKP_Silk_range_encoder(&psEnc.SCmn.SRC_LBRR, frame_terminator, SKP_Silk_FrameTermination_CDF[:])
		}
		memcpy(unsafe.Pointer(&psEncCtrl.SCmn.GainsIndices[0]), unsafe.Pointer(&TempGainsIndices[0]), size_t(NB_SUBFR*unsafe.Sizeof(int32(0))))
		memcpy(unsafe.Pointer(&psEncCtrl.Gains_Q16[0]), unsafe.Pointer(&TempGains_Q16[0]), size_t(NB_SUBFR*unsafe.Sizeof(int32(0))))
		psEncCtrl.SCmn.LTP_scaleIndex = LTP_scaleIndex
		psEnc.SCmn.TypeOffsetPrev = typeOffset
	}
}
