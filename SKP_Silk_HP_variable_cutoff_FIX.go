package silk

import "unsafe"

const SKP_RADIANS_CONSTANT_Q19 = 1482
const SKP_LOG2_VARIABLE_HP_MIN_FREQ_Q7 = 809

func SKP_Silk_HP_variable_cutoff_FIX(psEnc *SKP_Silk_encoder_state_FIX, psEncCtrl *SKP_Silk_encoder_control_FIX, out *int16, in *int16) {
	var (
		quality_Q15       int32
		B_Q28             [3]int32
		A_Q28             [2]int32
		Fc_Q19            int32
		r_Q28             int32
		r_Q22             int32
		pitch_freq_Hz_Q16 int32
		pitch_freq_log_Q7 int32
		delta_freq_Q7     int32
	)
	if int64(psEnc.SCmn.Prev_sigtype) == SIG_TYPE_VOICED {
		pitch_freq_Hz_Q16 = int32(((int64(psEnc.SCmn.Fs_kHz) * 1000) << 16) / int64(psEnc.SCmn.PrevLag))
		pitch_freq_log_Q7 = int32(int64(SKP_Silk_lin2log(pitch_freq_Hz_Q16)) - (16 << 7))
		quality_Q15 = psEncCtrl.Input_quality_bands_Q15[0]
		pitch_freq_log_Q7 = int32(int64(pitch_freq_log_Q7) - int64(SKP_SMULWB(SKP_SMULWB(quality_Q15<<2, quality_Q15), int32(int64(pitch_freq_log_Q7)-SKP_LOG2_VARIABLE_HP_MIN_FREQ_Q7))))
		pitch_freq_log_Q7 = int32(int64(pitch_freq_log_Q7) + ((int64(SKP_FIX_CONST(0.6, 15)) - int64(quality_Q15)) >> 9))
		delta_freq_Q7 = int32(int64(pitch_freq_log_Q7) - (int64(psEnc.Variable_HP_smth1_Q15) >> 8))
		if int64(delta_freq_Q7) < 0 {
			delta_freq_Q7 = int32(int64(delta_freq_Q7) * 3)
		}
		if int64(-SKP_FIX_CONST(0.4, 7)) > int64(SKP_FIX_CONST(0.4, 7)) {
			if int64(delta_freq_Q7) > int64(-SKP_FIX_CONST(0.4, 7)) {
				delta_freq_Q7 = -SKP_FIX_CONST(0.4, 7)
			} else if int64(delta_freq_Q7) < int64(SKP_FIX_CONST(0.4, 7)) {
				delta_freq_Q7 = SKP_FIX_CONST(0.4, 7)
			}
		} else if int64(delta_freq_Q7) > int64(SKP_FIX_CONST(0.4, 7)) {
			delta_freq_Q7 = SKP_FIX_CONST(0.4, 7)
		} else if int64(delta_freq_Q7) < int64(-SKP_FIX_CONST(0.4, 7)) {
			delta_freq_Q7 = -SKP_FIX_CONST(0.4, 7)
		}
		psEnc.Variable_HP_smth1_Q15 = SKP_SMLAWB(psEnc.Variable_HP_smth1_Q15, int32((int64(psEnc.Speech_activity_Q8)<<1)*int64(delta_freq_Q7)), SKP_FIX_CONST(0.1, 16))
	}
	psEnc.Variable_HP_smth2_Q15 = SKP_SMLAWB(psEnc.Variable_HP_smth2_Q15, int32(int64(psEnc.Variable_HP_smth1_Q15)-int64(psEnc.Variable_HP_smth2_Q15)), SKP_FIX_CONST(0.015, 16))
	psEncCtrl.Pitch_freq_low_Hz = SKP_Silk_log2lin(psEnc.Variable_HP_smth2_Q15 >> 8)
	if int64(SKP_FIX_CONST(80.0, 0)) > int64(SKP_FIX_CONST(150.0, 0)) {
		if int64(psEncCtrl.Pitch_freq_low_Hz) > int64(SKP_FIX_CONST(80.0, 0)) {
			psEncCtrl.Pitch_freq_low_Hz = SKP_FIX_CONST(80.0, 0)
		} else if int64(psEncCtrl.Pitch_freq_low_Hz) < int64(SKP_FIX_CONST(150.0, 0)) {
			psEncCtrl.Pitch_freq_low_Hz = SKP_FIX_CONST(150.0, 0)
		}
	} else if int64(psEncCtrl.Pitch_freq_low_Hz) > int64(SKP_FIX_CONST(150.0, 0)) {
		psEncCtrl.Pitch_freq_low_Hz = SKP_FIX_CONST(150.0, 0)
	} else if int64(psEncCtrl.Pitch_freq_low_Hz) < int64(SKP_FIX_CONST(80.0, 0)) {
		psEncCtrl.Pitch_freq_low_Hz = SKP_FIX_CONST(80.0, 0)
	}
	Fc_Q19 = int32(int64(SKP_SMULBB(SKP_RADIANS_CONSTANT_Q19, psEncCtrl.Pitch_freq_low_Hz)) / int64(psEnc.SCmn.Fs_kHz))
	r_Q28 = int32(int64(SKP_FIX_CONST(1.0, 28)) - int64(SKP_FIX_CONST(0.92, 9))*int64(Fc_Q19))
	B_Q28[0] = r_Q28
	B_Q28[1] = -r_Q28 << 1
	B_Q28[2] = r_Q28
	r_Q22 = r_Q28 >> 6
	A_Q28[0] = SKP_SMULWW(r_Q22, int32(int64(SKP_SMULWW(Fc_Q19, Fc_Q19))-int64(SKP_FIX_CONST(2.0, 22))))
	A_Q28[1] = SKP_SMULWW(r_Q22, r_Q22)
	length := psEnc.SCmn.Frame_length
	SKP_Silk_biquad_alt(unsafe.Slice(in, length), B_Q28, A_Q28, psEnc.SCmn.In_HP_State[:], unsafe.Slice(out, length))
}
