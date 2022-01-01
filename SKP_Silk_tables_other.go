package silk

import "math"

var TargetRate_table_NB [8]int32 = [8]int32{0, 8000, 9000, 11000, 13000, 16000, 22000, MAX_TARGET_RATE_BPS}
var TargetRate_table_MB [8]int32 = [8]int32{0, 10000, 12000, 14000, 17000, 21000, 28000, MAX_TARGET_RATE_BPS}
var TargetRate_table_WB [8]int32 = [8]int32{0, 11000, 14000, 17000, 21000, 26000, 36000, MAX_TARGET_RATE_BPS}
var TargetRate_table_SWB [8]int32 = [8]int32{0, 13000, 16000, 19000, 25000, 32000, 46000, MAX_TARGET_RATE_BPS}
var SNR_table_Q1 [8]int32 = [8]int32{19, 31, 35, 39, 43, 47, 54, 64}
var SNR_table_one_bit_per_sample_Q7 [4]int32 = [4]int32{1984, 2240, 2408, 2708}
var SKP_Silk_SWB_detect_B_HP_Q13 [3][3]int16 = [3][3]int16{{575, -948, 575}, {575, -221, 575}, {575, 104, 575}}
var SKP_Silk_SWB_detect_A_HP_Q13 [3][2]int16 = [3][2]int16{{0x3915, 6868}, {0x3253, 7337}, {0x2D42, 7911}}
var SKP_Silk_Dec_A_HP_24 [2]int16 = [2]int16{-16220, 8030}
var SKP_Silk_Dec_B_HP_24 [3]int16 = [3]int16{8000, -16000, 8000}
var SKP_Silk_Dec_A_HP_16 [2]int16 = [2]int16{-16127, 7940}
var SKP_Silk_Dec_B_HP_16 [3]int16 = [3]int16{8000, -16000, 8000}
var SKP_Silk_Dec_A_HP_12 [2]int16 = [2]int16{-16043, 7859}
var SKP_Silk_Dec_B_HP_12 [3]int16 = [3]int16{8000, -16000, 8000}
var SKP_Silk_Dec_A_HP_8 [2]int16 = [2]int16{-15885, 7710}
var SKP_Silk_Dec_B_HP_8 [3]int16 = [3]int16{8000, -16000, 8000}
var SKP_Silk_lsb_CDF [3]uint16 = [3]uint16{0, 40000, math.MaxUint16}
var SKP_Silk_LTPscale_CDF [4]uint16 = [4]uint16{0, 32000, 48000, math.MaxUint16}
var SKP_Silk_LTPscale_offset int32 = 2
var SKP_Silk_vadflag_CDF [3]uint16 = [3]uint16{0, 22000, math.MaxUint16}
var SKP_Silk_vadflag_offset int32 = 1
var SKP_Silk_SamplingRates_table [4]int32 = [4]int32{8, 12, 16, 24}
var SKP_Silk_SamplingRates_CDF [5]uint16 = [5]uint16{0, 16000, 32000, 48000, math.MaxUint16}
var SKP_Silk_SamplingRates_offset int32 = 2
var SKP_Silk_NLSF_interpolation_factor_CDF [6]uint16 = [6]uint16{0, 3706, 8703, 0x4B1A, 0x78CE, math.MaxUint16}
var SKP_Silk_NLSF_interpolation_factor_offset int32 = 4
var SKP_Silk_FrameTermination_CDF [5]uint16 = [5]uint16{0, 20000, 45000, 56000, math.MaxUint16}
var SKP_Silk_FrameTermination_offset int32 = 2
var SKP_Silk_Seed_CDF [5]uint16 = [5]uint16{0, 0x4000, 0x8000, 0xC000, math.MaxUint16}
var SKP_Silk_Seed_offset int32 = 2
var SKP_Silk_Quantization_Offsets_Q10 [2][2]int16 = [2][2]int16{{OFFSET_VL_Q10, OFFSET_VH_Q10}, {OFFSET_UVL_Q10, OFFSET_UVH_Q10}}
var SKP_Silk_LTPScales_table_Q14 [3]int16 = [3]int16{0x3CCD, 0x2CCD, 8192}
var SKP_Silk_Transition_LP_B_Q28 [5][3]int32 = [5][3]int32{{0xEF2670A, 0x1DE4CD56, 0xEF2670A}, {0xC825275, 0x19049A59, 0xC825275}, {0xA311146, 0x146203ED, 0xA311146}, {0x7D702DA, 0xFADC6F9, 0x7D702DA}, {0x552B622, 0xAA4FADA, 0x552B622}}
var SKP_Silk_Transition_LP_A_Q28 [5][2]int32 = [5][2]int32{{0x1E2EF346, 0xE4BE32B}, {0x1880661F, 0xA1D2C1C}, {306733530, 0x6F49CED}, {0xB1330EC, 0x4A590E3}, {0x21DA4ED, 0x36BDF0A}}
