package silk

// Reviewed by wdvxdr1123 2022-01-02

/***********************************************************
* Pitch analyser function
********************************************************** */

// SKP_Silk_decode_pitch SKP_Silk_decode_pitch.c
// len(pitch_lags) must be 4
func SKP_Silk_decode_pitch(
	lagIndex int32, /* I                             */
	contourIndex int32, /* O                             */
	pitch_lags []int32, /* O 4 pitch values              */
	Fs_kHz int32, /* I sampling frequency (kHz)    */
) {
	min_lag := SKP_SMULBB(2, Fs_kHz)

	/* Only for 24 / 16 kHz version for now */
	lag := min_lag + lagIndex
	if int64(Fs_kHz) == 8 {
		/* Only a small codebook for 8 khz */
		for i := 0; i < PITCH_EST_NB_SUBFR; i++ {
			pitch_lags[i] = lag + int32(SKP_Silk_CB_lags_stage2[i][contourIndex])
		}
	} else {
		for i := 0; i < PITCH_EST_NB_SUBFR; i++ {
			pitch_lags[i] = lag + int32(SKP_Silk_CB_lags_stage3[i][contourIndex])
		}
	}
}
