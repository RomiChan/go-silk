package silk

// SKP_Silk_decode_pitch SKP_Silk_decode_pitch.c
// len(pitch_lags) must be 4
func SKP_Silk_decode_pitch(lagIndex int32, contourIndex int32, pitch_lags []int32, Fs_kHz int32) {
	min_lag := SKP_SMULBB(2, Fs_kHz)
	lag := int32(int64(min_lag) + int64(lagIndex))
	if int64(Fs_kHz) == 8 {
		for i := 0; i < PITCH_EST_NB_SUBFR; i++ {
			pitch_lags[i] = int32(int64(lag) + int64(SKP_Silk_CB_lags_stage2[i][contourIndex]))
		}
	} else {
		for i := 0; i < PITCH_EST_NB_SUBFR; i++ {
			pitch_lags[i] = int32(int64(lag) + int64(SKP_Silk_CB_lags_stage3[i][contourIndex]))
		}
	}
}
