package silk

import "unsafe"

func SKP_Silk_control_audio_bandwidth(psEncC *SKP_Silk_encoder_state, TargetRate_bps int32) int32 {
	var fs_kHz int32
	fs_kHz = psEncC.Fs_kHz
	if int64(fs_kHz) == 0 {
		if int64(TargetRate_bps) >= SWB2WB_BITRATE_BPS {
			fs_kHz = 24
		} else if int64(TargetRate_bps) >= WB2MB_BITRATE_BPS {
			fs_kHz = 16
		} else if int64(TargetRate_bps) >= MB2NB_BITRATE_BPS {
			fs_kHz = 12
		} else {
			fs_kHz = 8
		}
		if int64(fs_kHz) < int64(int32(int64(psEncC.API_fs_Hz)/1000)) {
			fs_kHz = fs_kHz
		} else {
			fs_kHz = int32(int64(psEncC.API_fs_Hz) / 1000)
		}
		if int64(fs_kHz) < int64(psEncC.MaxInternal_fs_kHz) {
			fs_kHz = fs_kHz
		} else {
			fs_kHz = psEncC.MaxInternal_fs_kHz
		}
	} else if int64(SKP_SMULBB(fs_kHz, 1000)) > int64(psEncC.API_fs_Hz) || int64(fs_kHz) > int64(psEncC.MaxInternal_fs_kHz) {
		fs_kHz = int32(int64(psEncC.API_fs_Hz) / 1000)
		if int64(fs_kHz) < int64(psEncC.MaxInternal_fs_kHz) {
			fs_kHz = fs_kHz
		} else {
			fs_kHz = psEncC.MaxInternal_fs_kHz
		}
	} else {
		if int64(psEncC.API_fs_Hz) > 8000 {
			psEncC.BitrateDiff += int32(int64(psEncC.PacketSize_ms) * (int64(TargetRate_bps) - int64(psEncC.Bitrate_threshold_down)))
			if int64(psEncC.BitrateDiff) < 0 {
				psEncC.BitrateDiff = psEncC.BitrateDiff
			} else {
				psEncC.BitrateDiff = 0
			}
			if int64(psEncC.VadFlag) == NO_VOICE_ACTIVITY {
				if int64(psEncC.SLP.Transition_frame_no) == 0 && (int64(psEncC.BitrateDiff) <= -ACCUM_BITS_DIFF_THRESHOLD || int64(psEncC.SSWBdetect.WB_detected)*int64(psEncC.Fs_kHz) == 24) {
					psEncC.SLP.Transition_frame_no = 1
					psEncC.SLP.Mode = 0
				} else if int64(psEncC.SLP.Transition_frame_no) >= (TRANSITION_TIME_DOWN_MS/FRAME_LENGTH_MS) && int64(psEncC.SLP.Mode) == 0 {
					psEncC.SLP.Transition_frame_no = 0
					psEncC.BitrateDiff = 0
					if int64(psEncC.Fs_kHz) == 24 {
						fs_kHz = 16
					} else if int64(psEncC.Fs_kHz) == 16 {
						fs_kHz = 12
					} else {
						fs_kHz = 8
					}
				}
				if int64(psEncC.Fs_kHz)*1000 < int64(psEncC.API_fs_Hz) && int64(TargetRate_bps) >= int64(psEncC.Bitrate_threshold_up) && int64(psEncC.SSWBdetect.WB_detected)*int64(psEncC.Fs_kHz) < 16 && (int64(psEncC.Fs_kHz) == 16 && int64(psEncC.MaxInternal_fs_kHz) >= 24 || int64(psEncC.Fs_kHz) == 12 && int64(psEncC.MaxInternal_fs_kHz) >= 16 || int64(psEncC.Fs_kHz) == 8 && int64(psEncC.MaxInternal_fs_kHz) >= 12) && int64(psEncC.SLP.Transition_frame_no) == 0 {
					psEncC.SLP.Mode = 1
					psEncC.BitrateDiff = 0
					if int64(psEncC.Fs_kHz) == 8 {
						fs_kHz = 12
					} else if int64(psEncC.Fs_kHz) == 12 {
						fs_kHz = 16
					} else {
						fs_kHz = 24
					}
				}
			}
		}
		if int64(psEncC.SLP.Mode) == 1 && int64(psEncC.SLP.Transition_frame_no) >= (TRANSITION_TIME_UP_MS/FRAME_LENGTH_MS) && int64(psEncC.VadFlag) == NO_VOICE_ACTIVITY {
			psEncC.SLP.Transition_frame_no = 0
			memset(unsafe.Pointer(&psEncC.SLP.In_LP_State[0]), 0, size_t(unsafe.Sizeof(int32(0))*2))
		}
	}
	return fs_kHz
}
