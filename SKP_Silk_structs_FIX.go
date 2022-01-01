package silk

type SKP_Silk_shape_state_FIX struct {
	LastGainIndex          int32
	HarmBoost_smth_Q16     int32
	HarmShapeGain_smth_Q16 int32
	Tilt_smth_Q16          int32
}
type SKP_Silk_prefilter_state_FIX struct {
	SLTP_shp         [512]int16
	SAR_shp          [17]int32
	SLTP_shp_buf_idx int32
	SLF_AR_shp_Q12   int32
	SLF_MA_shp_Q12   int32
	SHarmHP          int32
	Rand_seed        int32
	LagPrev          int32
}
type SKP_Silk_predict_state_FIX struct {
	Pitch_LPC_win_length int32
	Min_pitch_lag        int32
	Max_pitch_lag        int32
	Prev_NLSFq_Q15       [16]int32
}
type SKP_Silk_encoder_state_FIX struct {
	SCmn                           SKP_Silk_encoder_state
	Variable_HP_smth1_Q15          int32
	Variable_HP_smth2_Q15          int32
	SShape                         SKP_Silk_shape_state_FIX
	SPrefilt                       SKP_Silk_prefilter_state_FIX
	SPred                          SKP_Silk_predict_state_FIX
	X_buf                          [1080]int16
	LTPCorr_Q15                    int32
	Mu_LTP_Q8                      int32
	SNR_dB_Q7                      int32
	AvgGain_Q16                    int32
	AvgGain_Q16_one_bit_per_sample int32
	BufferedInChannel_ms           int32
	Speech_activity_Q8             int32
	PrevLTPredCodGain_Q7           int32
	HPLTPredCodGain_Q7             int32
	InBandFEC_SNR_comp_Q8          int32
}
type SKP_Silk_encoder_control_FIX struct {
	SCmn                    SKP_Silk_encoder_control
	Gains_Q16               [4]int32
	PredCoef_Q12            [2][16]int16
	LTPCoef_Q14             [20]int16
	LTP_scale_Q14           int32
	AR1_Q13                 [64]int16
	AR2_Q13                 [64]int16
	LF_shp_Q14              [4]int32
	GainsPre_Q14            [4]int32
	HarmBoost_Q14           [4]int32
	Tilt_Q14                [4]int32
	HarmShapeGain_Q14       [4]int32
	Lambda_Q10              int32
	Input_quality_Q14       int32
	Coding_quality_Q14      int32
	Pitch_freq_low_Hz       int32
	Current_SNR_dB_Q7       int32
	Sparseness_Q8           int32
	PredGain_Q16            int32
	LTPredCodGain_Q7        int32
	Input_quality_bands_Q15 [4]int32
	Input_tilt_Q15          int32
	ResNrg                  [4]int32
	ResNrgQ                 [4]int32
}
