package silk

const FIND_PITCH_WHITE_NOISE_FRACTION = 0.001
const FIND_PITCH_BANDWITH_EXPANSION = 0.99
const FIND_PITCH_CORRELATION_THRESHOLD_HC_MODE = 0.7
const FIND_PITCH_CORRELATION_THRESHOLD_MC_MODE = 0.75
const FIND_PITCH_CORRELATION_THRESHOLD_LC_MODE = 0.8
const FIND_LPC_COND_FAC = 2.5e-05
const FIND_LPC_CHIRP = 0.99995
const FIND_LTP_COND_FAC = 1e-05
const MU_LTP_QUANT_NB = 0.03
const MU_LTP_QUANT_MB = 0.025
const MU_LTP_QUANT_WB = 0.02
const MU_LTP_QUANT_SWB = 0.016
const VARIABLE_HP_SMTH_COEF1 = 0.1
const VARIABLE_HP_SMTH_COEF2 = 0.015
const VARIABLE_HP_MIN_FREQ = 80.0
const VARIABLE_HP_MAX_FREQ = 150.0
const VARIABLE_HP_MAX_DELTA_FREQ = 0.4
const WB_DETECT_ACTIVE_SPEECH_LEVEL_THRES = 0.7
const SPEECH_ACTIVITY_DTX_THRES = 0.1
const LBRR_SPEECH_ACTIVITY_THRES = 0.5
const BG_SNR_DECR_dB = 4.0
const HARM_SNR_INCR_dB = 2.0
const SPARSE_SNR_INCR_dB = 2.0
const SPARSENESS_THRESHOLD_QNT_OFFSET = 0.75
const WARPING_MULTIPLIER = 0.015
const SHAPE_WHITE_NOISE_FRACTION = 1e-05
const BANDWIDTH_EXPANSION = 0.95
const LOW_RATE_BANDWIDTH_EXPANSION_DELTA = 0.01
const DE_ESSER_COEF_SWB_dB = 2.0
const DE_ESSER_COEF_WB_dB = 1.0
const LOW_RATE_HARMONIC_BOOST = 0.1
const LOW_INPUT_QUALITY_HARMONIC_BOOST = 0.1
const HARMONIC_SHAPING = 0.3
const HIGH_RATE_OR_LOW_QUALITY_HARMONIC_SHAPING = 0.2
const HP_NOISE_COEF = 0.3
const HARM_HP_NOISE_COEF = 0.35
const INPUT_TILT = 0.05
const HIGH_RATE_INPUT_TILT = 0.1
const LOW_FREQ_SHAPING = 3.0
const LOW_QUALITY_LOW_FREQ_SHAPING_DECR = 0.5
const NOISE_FLOOR_dB = 4.0
const RELATIVE_MIN_GAIN_dB = 0
const GAIN_SMOOTHING_COEF = 0.001
const SUBFR_SMTH_COEF = 0.4
const LAMBDA_OFFSET = 1.2
const LAMBDA_SPEECH_ACT = 0
const LAMBDA_DELAYED_DECISIONS = 0
const LAMBDA_INPUT_QUALITY = 0
const LAMBDA_CODING_QUALITY = 0
const LAMBDA_QUANT_OFFSET = 1.5