package silk

import "unsafe"

func SKP_Silk_control_encoder_FIX(psEnc *SKP_Silk_encoder_state_FIX, PacketSize_ms int32, TargetRate_bps int32, PacketLoss_perc int32, DTX_enabled int32, Complexity int32) int32 {
	var (
		fs_kHz int32
		ret    int32 = 0
	)
	if int64(psEnc.SCmn.Controlled_since_last_payload) != 0 {
		if int64(psEnc.SCmn.API_fs_Hz) != int64(psEnc.SCmn.Prev_API_fs_Hz) && int64(psEnc.SCmn.Fs_kHz) > 0 {
			ret += SKP_Silk_setup_resamplers_FIX(psEnc, psEnc.SCmn.Fs_kHz)
		}
		return ret
	}
	fs_kHz = SKP_Silk_control_audio_bandwidth(&psEnc.SCmn, TargetRate_bps)
	ret += SKP_Silk_setup_resamplers_FIX(psEnc, fs_kHz)
	ret += SKP_Silk_setup_packetsize_FIX(psEnc, PacketSize_ms)
	ret += SKP_Silk_setup_fs_FIX(psEnc, fs_kHz)
	ret += SKP_Silk_setup_complexity(&psEnc.SCmn, Complexity)
	ret += SKP_Silk_setup_rate_FIX(psEnc, TargetRate_bps)
	if int64(PacketLoss_perc) < 0 || int64(PacketLoss_perc) > 100 {
		ret = -5
	}
	psEnc.SCmn.PacketLoss_perc = PacketLoss_perc
	ret += SKP_Silk_setup_LBRR_FIX(psEnc)
	if int64(DTX_enabled) < 0 || int64(DTX_enabled) > 1 {
		ret = -8
	}
	psEnc.SCmn.UseDTX = DTX_enabled
	psEnc.SCmn.Controlled_since_last_payload = 1
	return ret
}
func SKP_Silk_LBRR_ctrl_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrlC *SKP_Silk_encoder_control) {
	var LBRR_usage int32
	if int64(psEnc.SCmn.LBRR_enabled) != 0 {
		LBRR_usage = SKP_SILK_NO_LBRR
		if int64(psEnc.Speech_activity_Q8) > int64(SKP_FIX_CONST(0.5, 8)) && int64(psEnc.SCmn.PacketLoss_perc) > LBRR_LOSS_THRES {
			LBRR_usage = SKP_SILK_ADD_LBRR_TO_PLUS1
		}
		psEncCtrlC.LBRR_usage = LBRR_usage
	} else {
		psEncCtrlC.LBRR_usage = SKP_SILK_NO_LBRR
	}
}
func SKP_Silk_setup_resamplers_FIX(psEnc *SKP_Silk_encoder_state_FIX, fs_kHz int32) int32 {
	var ret int32 = SKP_SILK_NO_ERROR
	if int64(psEnc.SCmn.Fs_kHz) != int64(fs_kHz) || int64(psEnc.SCmn.Prev_API_fs_Hz) != int64(psEnc.SCmn.API_fs_Hz) {
		if int64(psEnc.SCmn.Fs_kHz) == 0 {
			ret += SKP_Silk_resampler_init(&psEnc.SCmn.Resampler_state, psEnc.SCmn.API_fs_Hz, int32(int64(fs_kHz)*1000))
		} else {
			var (
				x_buf_API_fs_Hz [6480]int16
				nSamples_temp   int32 = int32((int64(psEnc.SCmn.Frame_length) << 1) + LA_SHAPE_MS*int64(psEnc.SCmn.Fs_kHz))
			)
			if int64(SKP_SMULBB(fs_kHz, 1000)) < int64(psEnc.SCmn.API_fs_Hz) && int64(psEnc.SCmn.Fs_kHz) != 0 {
				var temp_resampler_state SKP_Silk_resampler_state_struct
				ret += SKP_Silk_resampler_init(&temp_resampler_state, SKP_SMULBB(psEnc.SCmn.Fs_kHz, 1000), psEnc.SCmn.API_fs_Hz)
				ret += SKP_Silk_resampler(&temp_resampler_state, x_buf_API_fs_Hz[:], psEnc.X_buf[:], nSamples_temp)
				nSamples_temp = int32((int64(nSamples_temp) * int64(psEnc.SCmn.API_fs_Hz)) / int64(SKP_SMULBB(psEnc.SCmn.Fs_kHz, 1000)))
				ret += SKP_Silk_resampler_init(&psEnc.SCmn.Resampler_state, psEnc.SCmn.API_fs_Hz, SKP_SMULBB(fs_kHz, 1000))
			} else {
				memcpy(unsafe.Pointer(&x_buf_API_fs_Hz[0]), unsafe.Pointer(&psEnc.X_buf[0]), size_t(uintptr(nSamples_temp)*unsafe.Sizeof(int16(0))))
			}
			if int64(fs_kHz)*1000 != int64(psEnc.SCmn.API_fs_Hz) {
				ret += SKP_Silk_resampler(&psEnc.SCmn.Resampler_state, psEnc.X_buf[:], x_buf_API_fs_Hz[:], nSamples_temp)
			}
		}
	}
	psEnc.SCmn.Prev_API_fs_Hz = psEnc.SCmn.API_fs_Hz
	return ret
}
func SKP_Silk_setup_packetsize_FIX(psEnc *SKP_Silk_encoder_state_FIX, PacketSize_ms int32) int32 {
	var ret int32 = SKP_SILK_NO_ERROR
	if int64(PacketSize_ms) != 20 && int64(PacketSize_ms) != 40 && int64(PacketSize_ms) != 60 && int64(PacketSize_ms) != 80 && int64(PacketSize_ms) != 100 {
		ret = -3
	} else {
		if int64(PacketSize_ms) != int64(psEnc.SCmn.PacketSize_ms) {
			psEnc.SCmn.PacketSize_ms = PacketSize_ms
			SKP_Silk_LBRR_reset(&psEnc.SCmn)
		}
	}
	return ret
}
func SKP_Silk_setup_fs_FIX(psEnc *SKP_Silk_encoder_state_FIX, fs_kHz int32) int32 {
	var ret int32 = SKP_SILK_NO_ERROR
	if int64(psEnc.SCmn.Fs_kHz) != int64(fs_kHz) {
		memset(unsafe.Pointer(&psEnc.SShape), 0, size_t(unsafe.Sizeof(SKP_Silk_shape_state_FIX{})))
		memset(unsafe.Pointer(&psEnc.SPrefilt), 0, size_t(unsafe.Sizeof(SKP_Silk_prefilter_state_FIX{})))
		memset(unsafe.Pointer(&psEnc.SPred), 0, size_t(unsafe.Sizeof(SKP_Silk_predict_state_FIX{})))
		memset(unsafe.Pointer(&psEnc.SCmn.SNSQ), 0, size_t(unsafe.Sizeof(SKP_Silk_nsq_state{})))
		memset(unsafe.Pointer(&psEnc.SCmn.SNSQ_LBRR.Xq[0]), 0, size_t(uintptr((FRAME_LENGTH_MS*MAX_FS_KHZ)*2)*unsafe.Sizeof(int16(0))))
		memset(unsafe.Pointer(&psEnc.SCmn.LBRR_buffer[0]), 0, size_t(MAX_LBRR_DELAY*unsafe.Sizeof(SKP_SILK_LBRR_struct{})))
		memset(unsafe.Pointer(&psEnc.SCmn.SLP.In_LP_State[0]), 0, size_t(unsafe.Sizeof(int32(0))*2))
		if int64(psEnc.SCmn.SLP.Mode) == 1 {
			psEnc.SCmn.SLP.Transition_frame_no = 1
		} else {
			psEnc.SCmn.SLP.Transition_frame_no = 0
		}
		psEnc.SCmn.InputBufIx = 0
		psEnc.SCmn.NFramesInPayloadBuf = 0
		psEnc.SCmn.NBytesInPayloadBuf = 0
		psEnc.SCmn.Oldest_LBRR_idx = 0
		psEnc.SCmn.TargetRate_bps = 0
		memset(unsafe.Pointer(&psEnc.SPred.Prev_NLSFq_Q15[0]), 0, size_t(MAX_LPC_ORDER*unsafe.Sizeof(int32(0))))
		psEnc.SCmn.PrevLag = 100
		psEnc.SCmn.Prev_sigtype = SIG_TYPE_UNVOICED
		psEnc.SCmn.First_frame_after_reset = 1
		psEnc.SPrefilt.LagPrev = 100
		psEnc.SShape.LastGainIndex = 1
		psEnc.SCmn.SNSQ.LagPrev = 100
		psEnc.SCmn.SNSQ.Prev_inv_gain_Q16 = 0x10000
		psEnc.SCmn.SNSQ_LBRR.Prev_inv_gain_Q16 = 0x10000
		psEnc.SCmn.Fs_kHz = fs_kHz
		if int64(psEnc.SCmn.Fs_kHz) == 8 {
			psEnc.SCmn.PredictLPCOrder = MIN_LPC_ORDER
			psEnc.SCmn.PsNLSF_CB[0] = &SKP_Silk_NLSF_CB0_10
			psEnc.SCmn.PsNLSF_CB[1] = &SKP_Silk_NLSF_CB1_10
		} else {
			psEnc.SCmn.PredictLPCOrder = MAX_LPC_ORDER
			psEnc.SCmn.PsNLSF_CB[0] = &SKP_Silk_NLSF_CB0_16
			psEnc.SCmn.PsNLSF_CB[1] = &SKP_Silk_NLSF_CB1_16
		}
		psEnc.SCmn.Frame_length = SKP_SMULBB(FRAME_LENGTH_MS, fs_kHz)
		psEnc.SCmn.Subfr_length = int32(int64(psEnc.SCmn.Frame_length) / NB_SUBFR)
		psEnc.SCmn.La_pitch = SKP_SMULBB(LA_PITCH_MS, fs_kHz)
		psEnc.SPred.Min_pitch_lag = SKP_SMULBB(3, fs_kHz)
		psEnc.SPred.Max_pitch_lag = SKP_SMULBB(18, fs_kHz)
		psEnc.SPred.Pitch_LPC_win_length = SKP_SMULBB((LA_PITCH_MS<<1)+20, fs_kHz)
		if int64(psEnc.SCmn.Fs_kHz) == 24 {
			psEnc.Mu_LTP_Q8 = SKP_FIX_CONST(0.016, 8)
			psEnc.SCmn.Bitrate_threshold_up = SKP_int32_MAX
			psEnc.SCmn.Bitrate_threshold_down = SWB2WB_BITRATE_BPS
		} else if int64(psEnc.SCmn.Fs_kHz) == 16 {
			psEnc.Mu_LTP_Q8 = SKP_FIX_CONST(0.02, 8)
			psEnc.SCmn.Bitrate_threshold_up = WB2SWB_BITRATE_BPS
			psEnc.SCmn.Bitrate_threshold_down = WB2MB_BITRATE_BPS
		} else if int64(psEnc.SCmn.Fs_kHz) == 12 {
			psEnc.Mu_LTP_Q8 = SKP_FIX_CONST(0.025, 8)
			psEnc.SCmn.Bitrate_threshold_up = MB2WB_BITRATE_BPS
			psEnc.SCmn.Bitrate_threshold_down = MB2NB_BITRATE_BPS
		} else {
			psEnc.Mu_LTP_Q8 = SKP_FIX_CONST(0.03, 8)
			psEnc.SCmn.Bitrate_threshold_up = NB2MB_BITRATE_BPS
			psEnc.SCmn.Bitrate_threshold_down = 0
		}
		psEnc.SCmn.Fs_kHz_changed = 1
	}
	return ret
}
func SKP_Silk_setup_rate_FIX(psEnc *SKP_Silk_encoder_state_FIX, TargetRate_bps int32) int32 {
	var (
		k         int32
		ret       int32 = SKP_SILK_NO_ERROR
		frac_Q6   int32
		rateTable *int32
	)
	if int64(TargetRate_bps) != int64(psEnc.SCmn.TargetRate_bps) {
		psEnc.SCmn.TargetRate_bps = TargetRate_bps
		if int64(psEnc.SCmn.Fs_kHz) == 8 {
			rateTable = &TargetRate_table_NB[0]
		} else if int64(psEnc.SCmn.Fs_kHz) == 12 {
			rateTable = &TargetRate_table_MB[0]
		} else if int64(psEnc.SCmn.Fs_kHz) == 16 {
			rateTable = &TargetRate_table_WB[0]
		} else {
			rateTable = &TargetRate_table_SWB[0]
		}
		for k = 1; int64(k) < TARGET_RATE_TAB_SZ; k++ {
			if int64(TargetRate_bps) <= int64(*(*int32)(unsafe.Add(unsafe.Pointer(rateTable), unsafe.Sizeof(int32(0))*uintptr(k)))) {
				frac_Q6 = int32(((int64(TargetRate_bps) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(rateTable), unsafe.Sizeof(int32(0))*uintptr(int64(k)-1))))) << 6) / (int64(*(*int32)(unsafe.Add(unsafe.Pointer(rateTable), unsafe.Sizeof(int32(0))*uintptr(k)))) - int64(*(*int32)(unsafe.Add(unsafe.Pointer(rateTable), unsafe.Sizeof(int32(0))*uintptr(int64(k)-1))))))
				psEnc.SNR_dB_Q7 = int32((int64(SNR_table_Q1[int64(k)-1]) << 6) + int64(frac_Q6)*(int64(SNR_table_Q1[k])-int64(SNR_table_Q1[int64(k)-1])))
				break
			}
		}
	}
	return ret
}
func SKP_Silk_setup_LBRR_FIX(psEnc *SKP_Silk_encoder_state_FIX) int32 {
	var (
		ret                int32 = SKP_SILK_NO_ERROR
		LBRRRate_thres_bps int32
	)
	if int64(psEnc.SCmn.UseInBandFEC) < 0 || int64(psEnc.SCmn.UseInBandFEC) > 1 {
		ret = -7
	}
	psEnc.SCmn.LBRR_enabled = psEnc.SCmn.UseInBandFEC
	if int64(psEnc.SCmn.Fs_kHz) == 8 {
		LBRRRate_thres_bps = INBAND_FEC_MIN_RATE_BPS - 9000
	} else if int64(psEnc.SCmn.Fs_kHz) == 12 {
		LBRRRate_thres_bps = INBAND_FEC_MIN_RATE_BPS - 6000
	} else if int64(psEnc.SCmn.Fs_kHz) == 16 {
		LBRRRate_thres_bps = INBAND_FEC_MIN_RATE_BPS - 3000
	} else {
		LBRRRate_thres_bps = INBAND_FEC_MIN_RATE_BPS
	}
	if int64(psEnc.SCmn.TargetRate_bps) >= int64(LBRRRate_thres_bps) {
		psEnc.SCmn.LBRR_GainIncreases = SKP_max_int(int32(8-(int64(psEnc.SCmn.PacketLoss_perc)>>1)), 0)
		if int64(psEnc.SCmn.LBRR_enabled) != 0 && int64(psEnc.SCmn.PacketLoss_perc) > LBRR_LOSS_THRES {
			psEnc.InBandFEC_SNR_comp_Q8 = int32(int64(SKP_FIX_CONST(6.0, 8)) - (int64(psEnc.SCmn.LBRR_GainIncreases) << 7))
		} else {
			psEnc.InBandFEC_SNR_comp_Q8 = 0
			psEnc.SCmn.LBRR_enabled = 0
		}
	} else {
		psEnc.InBandFEC_SNR_comp_Q8 = 0
		psEnc.SCmn.LBRR_enabled = 0
	}
	return ret
}
