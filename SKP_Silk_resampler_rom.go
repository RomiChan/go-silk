package silk

import "math"

const RESAMPLER_DOWN_ORDER_FIR = 12
const RESAMPLER_ORDER_FIR_144 = 6

var SKP_Silk_resampler_down2_0 int16 = 9872
var SKP_Silk_resampler_down2_1 int16 = 0x9B81 - 0x10000
var SKP_Silk_resampler_up2_lq_0 int16 = 8102
var SKP_Silk_resampler_up2_lq_1 int16 = 0x8FAF - 0x10000
var SKP_Silk_resampler_up2_hq_0 [2]int16 = [2]int16{4280, 0x83BF - 0x10000}
var SKP_Silk_resampler_up2_hq_1 [2]int16 = [2]int16{0x3FA7, 0xD2FF - 0x10000}
var SKP_Silk_resampler_up2_hq_notch [4]int16 = [4]int16{7864, -3604, 0x3333, 0x6F5C}
var SKP_Silk_Resampler_3_4_COEFS [20]int16 = [20]int16{-18249, -12532, -97, 284, -495, 309, 0x281C, 0x4F5D, -94, 156, -48, -720, 5984, 0x4766, -45, -4, 237, -847, 2540, 0x3946}
var SKP_Silk_Resampler_2_3_COEFS [14]int16 = [14]int16{-11891, -12486, 20, 211, -657, 688, 8423, 0x3E27, -44, 197, -152, -653, 3855, 0x32D7}
var SKP_Silk_Resampler_1_2_COEFS [8]int16 = [8]int16{2415, -13101, 158, -295, -400, 1265, 4832, 7968}
var SKP_Silk_Resampler_3_8_COEFS [20]int16 = [20]int16{13270, -13738, -294, -123, 747, 2043, 3339, 3995, -151, -311, 414, 1583, 2947, 3877, -33, -389, 143, 1141, 2503, 3653}
var SKP_Silk_Resampler_1_3_COEFS [8]int16 = [8]int16{0x4103, -14000, -331, 19, 581, 1421, 2290, 2845}
var SKP_Silk_Resampler_2_3_COEFS_LQ [6]int16 = [6]int16{-2797, -6507, 4697, 0x29F3, 1567, 8276}
var SKP_Silk_Resampler_1_3_COEFS_LQ [5]int16 = [5]int16{0x4189, -9792, 890, 1614, 2148}
var SKP_Silk_Resampler_320_441_ARMA4_COEFS [7]int16 = [7]int16{0x7ADE, 0x60AA, -9706, -3386, -17911, -13243, 0x60DD}
var SKP_Silk_Resampler_240_441_ARMA4_COEFS [7]int16 = [7]int16{0x7031, 0x2BF6, 3189, -2546, -1495, -12618, 0x2D2A}
var SKP_Silk_Resampler_160_441_ARMA4_COEFS [7]int16 = [7]int16{0x5BC4, -6457, 0x3816, -4856, 0x393E, -13008, 4456}
var SKP_Silk_Resampler_120_441_ARMA4_COEFS [7]int16 = [7]int16{0x4B6F, -15569, 0x4C21, -6950, 0x53C1, -13559, 2370}
var SKP_Silk_Resampler_80_441_ARMA4_COEFS [7]int16 = [7]int16{0x33C0, -23849, 0x5E3E, -9486, 0x68B6, -14286, 1065}
var SKP_Silk_resampler_frac_FIR_144 [144][3]int16 = [144][3]int16{{-647, 1884, 0x757E}, {-625, 1736, 0x755C}, {-603, 1591, 0x7535}, {-581, 1448, 0x750B}, {-559, 1308, 0x74DD}, {-537, 1169, 0x74AB}, {-515, 1032, 0x7475}, {-494, 898, 0x743B}, {-473, 766, 0x73FD}, {-452, 636, 0x73BB}, {-431, 508, 0x7376}, {-410, 383, 0x732C}, {-390, 260, 0x72DF}, {-369, 139, 0x728F}, {-349, 20, 0x723A}, {-330, -97, 0x71E2}, {-310, -211, 0x7186}, {-291, -324, 0x7127}, {-271, -434, 0x70C4}, {-253, -542, 0x705D}, {-234, -647, 0x6FF3}, {-215, -751, 28550}, {-197, -852, 0x6F14}, {-179, -951, 0x6EA0}, {-162, -1048, 28200}, {-144, -1143, 0x6DAD}, {-127, -1235, 27950}, {-110, -1326, 27820}, {-94, -1414, 0x6C27}, {-77, -1500, 27550}, {-61, -1584, 27410}, {-45, -1665, 0x6A84}, {-30, -1745, 0x69F2}, {-15, -1822, 0x695C}, {0, -1897, 26820}, {15, -1970, 0x6829}, {29, -2041, 0x678B}, {44, -2110, 0x66EA}, {57, -2177, 0x6646}, {71, -2242, 0x659F}, {84, -2305, 0x64F5}, {97, -2365, 0x6449}, {110, -2424, 0x639A}, {122, -2480, 25320}, {134, -2534, 25140}, {146, -2587, 0x617C}, {157, -2637, 0x60C3}, {168, -2685, 0x6007}, {179, -2732, 0x5F48}, {190, -2776, 0x5E87}, {200, -2819, 0x5DC3}, {210, -2859, 0x5CFD}, {220, -2898, 0x5C35}, {229, -2934, 0x5B6B}, {238, -2969, 0x5A9E}, {247, -3002, 0x59D0}, {math.MaxUint8, -3033, 0x58FF}, {263, -3062, 0x582C}, {271, -3089, 0x5757}, {279, -3114, 0x5680}, {286, -3138, 0x55A7}, {293, -3160, 0x54CD}, {300, -3180, 0x53F0}, {306, -3198, 0x5312}, {312, -3215, 0x5232}, {318, -3229, 0x5150}, {323, -3242, 0x506D}, {328, -3254, 20360}, {333, -3263, 20130}, {338, -3272, 0x4DBA}, {342, -3278, 0x4CD1}, {346, -3283, 19430}, {350, -3286, 0x4AFA}, {353, -3288, 0x4A0D}, {356, -3288, 0x491E}, {359, -3286, 0x482E}, {362, -3283, 0x473E}, {364, -3279, 0x464C}, {366, -3273, 0x4559}, {368, -3266, 0x4465}, {369, -3257, 0x4370}, {371, -3247, 0x427A}, {372, -3235, 0x4184}, {372, -3222, 0x408D}, {373, -3208, 0x3F95}, {373, -3192, 0x3E9C}, {373, -3175, 0x3DA3}, {373, -3157, 0x3CA9}, {372, -3138, 0x3BAF}, {371, -3117, 0x3AB4}, {370, -3095, 0x39B9}, {369, -3072, 0x38BE}, {368, -3048, 0x37C2}, {366, -3022, 0x36C6}, {364, -2996, 13770}, {362, -2968, 0x34CD}, {359, -2940, 0x33D1}, {357, -2910, 0x32D4}, {354, -2880, 12760}, {351, -2848, 0x30DC}, {348, -2815, 0x2FDF}, {344, -2782, 0x2EE3}, {341, -2747, 0x2DE7}, {337, -2712, 11500}, {333, -2676, 0x2BF0}, {328, -2639, 0x2AF5}, {324, -2601, 0x29FB}, {320, -2562, 0x2901}, {315, -2523, 0x2807}, {310, -2482, 9998}, {305, -2442, 9750}, {300, -2400, 9502}, {294, -2358, 9255}, {289, -2315, 9009}, {283, -2271, 8763}, {277, -2227, 8519}, {271, -2182, 8275}, {265, -2137, 8032}, {259, -2091, 7791}, {252, -2045, 7550}, {246, -1998, 7311}, {239, -1951, 7072}, {232, -1904, 6835}, {226, -1856, 6599}, {219, -1807, 6364}, {212, -1758, 6131}, {204, -1709, 5899}, {197, -1660, 5668}, {190, -1611, 5439}, {183, -1561, 5212}, {175, -1511, 4986}, {168, -1460, 4761}, {160, -1410, 4538}, {152, -1359, 4317}, {145, -1309, 4098}, {137, -1258, 3880}, {129, -1207, 3664}, {121, -1156, 3450}, {113, -1105, 3238}, {105, -1054, 3028}, {97, -1003, 2820}, {89, -952, 2614}, {81, -901, 2409}, {73, -851, 2207}}
