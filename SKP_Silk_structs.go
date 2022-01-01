package silk

type SKP_Silk_nsq_state struct {
	Xq                [960]int16
	SLTP_shp_Q10      [960]int32
	SLPC_Q14          [152]int32
	SAR2_Q14          [16]int32
	SLF_AR_shp_Q12    int32
	LagPrev           int32
	SLTP_buf_idx      int32
	SLTP_shp_buf_idx  int32
	Rand_seed         int32
	Prev_inv_gain_Q16 int32
	Rewhite_flag      int32
}
type SKP_SILK_LBRR_struct struct {
	Payload [1024]uint8
	NBytes  int32
	Usage   int32
}
type SKP_Silk_VAD_state struct {
	AnaState        [2]int32
	AnaState1       [2]int32
	AnaState2       [2]int32
	XnrgSubfr       [4]int32
	NrgRatioSmth_Q8 [4]int32
	HPstate         int16
	NL              [4]int32
	Inv_NL          [4]int32
	NoiseLevelBias  [4]int32
	Counter         int32
}
type SKP_Silk_range_coder_state struct {
	BufferLength int32
	BufferIx     int32
	Base_Q32     uint32
	Range_Q16    uint32
	Error        int32
	Buffer       [1024]uint8
}
type SKP_Silk_detect_SWB_state struct {
	S_HP_8_kHz            [3][2]int32
	ConsecSmplsAboveThres int32
	ActiveSpeech_ms       int32
	SWB_detected          int32
	WB_detected           int32
}
type SKP_Silk_LP_state struct {
	In_LP_State         [2]int32
	Transition_frame_no int32
	Mode                int32
}
type SKP_Silk_NLSF_CBS struct {
	NVectors    int32
	CB_NLSF_Q15 *int16
	Rates_Q5    *int16
}
type SKP_Silk_NLSF_CB_struct struct {
	NStages       int32
	CBStages      *SKP_Silk_NLSF_CBS
	NDeltaMin_Q15 *int32
	CDF           *uint16
	StartPtr      **uint16
	MiddleIx      *int32
}
type SKP_Silk_encoder_state struct {
	SRC                           SKP_Silk_range_coder_state
	SRC_LBRR                      SKP_Silk_range_coder_state
	SNSQ                          SKP_Silk_nsq_state
	SNSQ_LBRR                     SKP_Silk_nsq_state
	In_HP_State                   [2]int32
	SLP                           SKP_Silk_LP_state
	SVAD                          SKP_Silk_VAD_state
	LBRRprevLastGainIndex         int32
	Prev_sigtype                  int32
	TypeOffsetPrev                int32
	PrevLag                       int32
	Prev_lagIndex                 int32
	API_fs_Hz                     int32
	Prev_API_fs_Hz                int32
	MaxInternal_fs_kHz            int32
	Fs_kHz                        int32
	Fs_kHz_changed                int32
	Frame_length                  int32
	Subfr_length                  int32
	La_pitch                      int32
	La_shape                      int32
	ShapeWinLength                int32
	TargetRate_bps                int32
	PacketSize_ms                 int32
	PacketLoss_perc               int32
	FrameCounter                  int32
	Complexity                    int32
	NStatesDelayedDecision        int32
	UseInterpolatedNLSFs          int32
	ShapingLPCOrder               int32
	PredictLPCOrder               int32
	PitchEstimationComplexity     int32
	PitchEstimationLPCOrder       int32
	PitchEstimationThreshold_Q16  int32
	LTPQuantLowComplexity         int32
	NLSF_MSVQ_Survivors           int32
	First_frame_after_reset       int32
	Controlled_since_last_payload int32
	Warping_Q16                   int32
	InputBuf                      [480]int16
	InputBufIx                    int32
	NFramesInPayloadBuf           int32
	NBytesInPayloadBuf            int32
	Frames_since_onset            int32
	PsNLSF_CB                     [2]*SKP_Silk_NLSF_CB_struct
	LBRR_buffer                   [2]SKP_SILK_LBRR_struct
	Oldest_LBRR_idx               int32
	UseInBandFEC                  int32
	LBRR_enabled                  int32
	LBRR_GainIncreases            int32
	BitrateDiff                   int32
	Bitrate_threshold_up          int32
	Bitrate_threshold_down        int32
	Resampler_state               SKP_Silk_resampler_state_struct
	NoSpeechCounter               int32
	UseDTX                        int32
	InDTX                         int32
	VadFlag                       int32
	SSWBdetect                    SKP_Silk_detect_SWB_state
	Q                             [480]int8
	Q_LBRR                        [480]int8
}
type SKP_Silk_encoder_control struct {
	LagIndex          int32
	ContourIndex      int32
	PERIndex          int32
	LTPIndex          [4]int32
	NLSFIndices       [10]int32
	NLSFInterpCoef_Q2 int32
	GainsIndices      [4]int32
	Seed              int32
	LTP_scaleIndex    int32
	RateLevelIndex    int32
	QuantOffsetType   int32
	Sigtype           int32
	PitchL            [4]int32
	LBRR_usage        int32
}
type SKP_Silk_PLC_struct struct {
	PitchL_Q8         int32
	LTPCoef_Q14       [5]int16
	PrevLPC_Q12       [16]int16
	Last_frame_lost   int32
	Rand_seed         int32
	RandScale_Q14     int16
	Conc_energy       int32
	Conc_energy_shift int32
	PrevLTP_scale_Q14 int16
	PrevGain_Q16      [4]int32
	Fs_kHz            int32
}
type SKP_Silk_CNG_struct struct {
	CNG_exc_buf_Q10   [480]int32
	CNG_smth_NLSF_Q15 [16]int32
	CNG_synth_state   [16]int32
	CNG_smth_Gain_Q16 int32
	Rand_seed         int32
	Fs_kHz            int32
}
type SKP_Silk_decoder_state struct {
	SRC                       SKP_Silk_range_coder_state
	Prev_inv_gain_Q16         int32
	SLTP_Q16                  [960]int32
	SLPC_Q14                  [136]int32
	Exc_Q10                   [480]int32
	Res_Q10                   [480]int32
	OutBuf                    [960]int16
	LagPrev                   int32
	LastGainIndex             int32
	LastGainIndex_EnhLayer    int32
	TypeOffsetPrev            int32
	HPState                   [2]int32
	HP_A                      *int16
	HP_B                      *int16
	Fs_kHz                    int32
	Prev_API_sampleRate       int32
	Frame_length              int32
	Subfr_length              int32
	LPC_order                 int32
	PrevNLSF_Q15              [16]int32
	First_frame_after_reset   int32
	NBytesLeft                int32
	NFramesDecoded            int32
	NFramesInPacket           int32
	MoreInternalDecoderFrames int32
	FrameTermination          int32
	Resampler_state           SKP_Silk_resampler_state_struct
	PsNLSF_CB                 [2]*SKP_Silk_NLSF_CB_struct
	VadFlag                   int32
	No_FEC_counter            int32
	Inband_FEC_offset         int32
	SCNG                      SKP_Silk_CNG_struct
	LossCnt                   int32
	Prev_sigtype              int32
	SPLC                      SKP_Silk_PLC_struct
}
type SKP_Silk_decoder_control struct {
	PitchL            [4]int32
	Gains_Q16         [4]int32
	Seed              int32
	PredCoef_Q12      [2][16]int16
	LTPCoef_Q14       [20]int16
	LTP_scale_Q14     int32
	PERIndex          int32
	RateLevelIndex    int32
	QuantOffsetType   int32
	Sigtype           int32
	NLSFInterpCoef_Q2 int32
}
