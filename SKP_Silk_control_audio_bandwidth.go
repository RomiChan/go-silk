package silk

import "unsafe"

func SKP_Silk_control_audio_bandwidth(psEncC *SKP_Silk_encoder_state, TargetRate_bps int32) int32 {
	var fs_kHz int32
	fs_kHz = psEncC.Fs_kHz
	if fs_kHz == 0 {
		if TargetRate_bps >= SWB2WB_BITRATE_BPS {
			fs_kHz = 24
		} else if TargetRate_bps >= WB2MB_BITRATE_BPS {
			fs_kHz = 16
		} else if TargetRate_bps >= MB2NB_BITRATE_BPS {
			fs_kHz = 12
		} else {
			fs_kHz = 8
		}
		if fs_kHz < (psEncC.API_fs_Hz / 1000) {
			fs_kHz = fs_kHz
		} else {
			fs_kHz = psEncC.API_fs_Hz / 1000
		}
		if fs_kHz < psEncC.MaxInternal_fs_kHz {
			fs_kHz = fs_kHz
		} else {
			fs_kHz = psEncC.MaxInternal_fs_kHz
		}
	} else if SKP_SMULBB(fs_kHz, 1000) > psEncC.API_fs_Hz || fs_kHz > psEncC.MaxInternal_fs_kHz {
		fs_kHz = psEncC.API_fs_Hz / 1000
		if fs_kHz < psEncC.MaxInternal_fs_kHz {
			fs_kHz = fs_kHz
		} else {
			fs_kHz = psEncC.MaxInternal_fs_kHz
		}
	} else {
		if psEncC.API_fs_Hz > 8000 {
			psEncC.BitrateDiff += psEncC.PacketSize_ms * (TargetRate_bps - psEncC.Bitrate_threshold_down)
			if psEncC.BitrateDiff < 0 {
				psEncC.BitrateDiff = psEncC.BitrateDiff
			} else {
				psEncC.BitrateDiff = 0
			}
			if psEncC.VadFlag == NO_VOICE_ACTIVITY {
				if psEncC.SLP.Transition_frame_no == 0 && (psEncC.BitrateDiff <= -ACCUM_BITS_DIFF_THRESHOLD || psEncC.SSWBdetect.WB_detected*psEncC.Fs_kHz == 24) {
					psEncC.SLP.Transition_frame_no = 1
					psEncC.SLP.Mode = 0
				} else if psEncC.SLP.Transition_frame_no >= (TRANSITION_TIME_DOWN_MS/FRAME_LENGTH_MS) && psEncC.SLP.Mode == 0 {
					psEncC.SLP.Transition_frame_no = 0
					psEncC.BitrateDiff = 0
					if psEncC.Fs_kHz == 24 {
						fs_kHz = 16
					} else if psEncC.Fs_kHz == 16 {
						fs_kHz = 12
					} else {
						fs_kHz = 8
					}
				}
				if psEncC.Fs_kHz*1000 < psEncC.API_fs_Hz && TargetRate_bps >= psEncC.Bitrate_threshold_up && psEncC.SSWBdetect.WB_detected*psEncC.Fs_kHz < 16 && (psEncC.Fs_kHz == 16 && psEncC.MaxInternal_fs_kHz >= 24 || psEncC.Fs_kHz == 12 && psEncC.MaxInternal_fs_kHz >= 16 || psEncC.Fs_kHz == 8 && psEncC.MaxInternal_fs_kHz >= 12) && psEncC.SLP.Transition_frame_no == 0 {
					psEncC.SLP.Mode = 1
					psEncC.BitrateDiff = 0
					if psEncC.Fs_kHz == 8 {
						fs_kHz = 12
					} else if psEncC.Fs_kHz == 12 {
						fs_kHz = 16
					} else {
						fs_kHz = 24
					}
				}
			}
		}
		if psEncC.SLP.Mode == 1 && psEncC.SLP.Transition_frame_no >= (TRANSITION_TIME_UP_MS/FRAME_LENGTH_MS) && psEncC.VadFlag == NO_VOICE_ACTIVITY {
			psEncC.SLP.Transition_frame_no = 0
			memset(unsafe.Pointer(&psEncC.SLP.In_LP_State[0]), 0, size_t(unsafe.Sizeof(int32(0))*2))
		}
	}
	return fs_kHz
}
