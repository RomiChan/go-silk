package silk

import "math"

const NLSF_MSVQ_CB0_16_STAGES = 10
const NLSF_MSVQ_CB0_16_VECTORS = 216

var SKP_Silk_NLSF_MSVQ_CB0_16_CDF [226]uint16 = [226]uint16{0, 1449, 2749, 4022, 5267, 6434, 7600, 8647, 9695, 0x29F6, 0x2DA1, 0x3139, 0x3484, 0x37AB, 0x3AA0, 0x3D94, 0x4089, 0x436D, 0x4652, 18710, 0x4BDB, 0x4EA0, 0x5165, 0x541B, 0x56D1, 0x5987, 0x5C2E, 0x5ED5, 0x617C, 25620, 0x6690, 0x68F1, 0x6B53, 0x6DA7, 0x6FFB, 0x724F, 0x74A3, 0x76EB, 0x7932, 0x7B6D, 0x7DA8, 0x7FD7, 0x81F9, 0x8410, 0x861C, 0x8827, 0x8A33, 0x8C3E, 0x8E3F, 0x903F, 0x923F, 0x9435, 0x962A, 0x9814, 0x99FF, 0x9BEA, 0x9DD4, 0x9FB4, 0xA194, 0xA374, 0xA554, 0xA735, 0xA915, 0xAAEA, 0xACC0, 0xAE96, 0xB06C, 0xB237, 0xB403, 0xB5CF, 0xB791, 0xB953, 0xBB14, 0xBCCD, 0xBE85, 0xC03D, 0xC1F5, 0xC3A4, 0xC553, 0xC702, 0xC8B1, 0xCA56, 0xCBF3, 0xCD86, 0xCF1A, 0xD0AE, 0xD239, 0xD3C4, 0xD54F, 0xD6DA, 0xD85C, 0xD9DF, 0xDB62, 0xDCE4, 56910, 0xDFB9, 0xE123, 0xE28D, 0xE3F0, 0xE54A, 0xE69E, 59370, 0xE936, 0xEA6E, 0xEBA5, 60630, 0xEE06, 0xEF37, 0xF061, 0xF17E, 0xF284, 0xF38A, 0xF48A, 0xF575, 0xF660, 0xF746, 0xF81D, 0xF8F4, 0xF9C3, 0xFA91, 0xFB56, 0xFC10, 0xFCBD, 0xFD63, 0xFDFF, 0xFE7F, 0xFEFF, 0xFF7F, math.MaxUint16, 0, 5099, 9982, 14760, 0x4C52, 0x5E95, 0x6FB3, 0x80D0, 0x9082, 0xA034, 0xAF90, 0xBE97, 0xCD4D, 0xDAC9, 0xE7FC, 0xF4E6, math.MaxUint16, 0, 9955, 0x4CF1, 0x7099, 0x8FEA, 0xAE8E, 0xCBE6, 0xE63B, math.MaxUint16, 0, 8949, 0x43B7, 25720, 0x8486, 0xA3E5, 0xC343, 0xE1F5, math.MaxUint16, 0, 9724, 0x48D2, 0x6976, 0x8A1B, 0xAA0C, 0xC94E, 0xE7E5, math.MaxUint16, 0, 8750, 0x445B, 0x6689, 0x8690, 0xA5E7, 0xC53E, 0xE342, math.MaxUint16, 0, 8730, 0x4379, 0x64D8, 0x8580, 0xA628, 0xC41B, 0xE20D, math.MaxUint16, 0, 8769, 0x4482, 0x66C3, 0x86DD, 0xA6F6, 0xC660, 0xE3CF, math.MaxUint16, 0, 8736, 0x42CD, 0x637A, 0x8375, 0xA36F, 0xC369, 0xE208, math.MaxUint16, 0, 4368, 8735, 0x3276, 17100, 0x5323, 0x6379, 0x7376, 0x8373, 0x9370, 0xA36C, 0xB369, 0xC366, 0xD30B, 0xE25B, 0xF156, math.MaxUint16}
var SKP_Silk_NLSF_MSVQ_CB0_16_CDF_start_ptr [10]*uint16 = [10]*uint16{&SKP_Silk_NLSF_MSVQ_CB0_16_CDF[0], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[129], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[146], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[155], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[164], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[173], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[182], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[191], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[200], &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[209]}
var SKP_Silk_NLSF_MSVQ_CB0_16_CDF_middle_idx [10]int32 = [10]int32{42, 8, 4, 5, 5, 5, 5, 5, 5, 9}
var SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5 [216]int16 = [216]int16{176, 181, 182, 183, 186, 186, 191, 191, 191, 196, 197, 201, 203, 206, 206, 206, 207, 207, 209, 209, 209, 209, 210, 210, 210, 211, 211, 211, 212, 214, 216, 216, 217, 217, 217, 217, 218, 218, 219, 219, 220, 221, 222, 223, 223, 223, 223, 224, 224, 224, 225, 225, 226, 226, 226, 226, 227, 227, 227, 227, 227, 227, 228, 228, 228, 228, 229, 229, 229, 230, 230, 230, 231, 231, 231, 231, 232, 232, 232, 232, 233, 234, 235, 235, 235, 236, 236, 236, 236, 237, 237, 237, 237, 240, 240, 240, 240, 241, 242, 243, 244, 244, 247, 247, 248, 248, 248, 249, 251, math.MaxUint8, math.MaxUint8, 256, 260, 260, 261, 264, 264, 266, 266, 268, 271, 274, 276, 279, 288, 288, 288, 288, 118, 120, 121, 121, 122, 125, 125, 129, 129, 130, 131, 132, 136, 137, 138, 145, 87, 88, 91, 97, 98, 100, 105, 106, 92, 95, 95, 96, 97, 97, 98, 99, 88, 92, 95, 95, 96, 97, 98, 109, 93, 93, 93, 96, 97, 97, 99, 101, 93, 94, 94, 95, 95, 99, 99, 99, 93, 93, 93, 96, 96, 97, 100, 102, 93, 95, 95, 96, 96, 96, 98, 99, 125, 125, math.MaxInt8, math.MaxInt8, math.MaxInt8, math.MaxInt8, 128, 128, 128, 128, 128, 128, 129, 130, 131, 132}
var SKP_Silk_NLSF_MSVQ_CB0_16_ndelta_min_Q15 [17]int32 = [17]int32{266, 3, 40, 3, 3, 16, 78, 89, 107, 141, 188, 146, 272, 240, 235, 215, 632}
var SKP_Silk_NLSF_MSVQ_CB0_16_Q15 [3456]int16 = [3456]int16{1170, 2278, 3658, 5374, 7666, 9113, 0x2C22, 0x33F8, 0x3C0B, 0x448D, 0x4C83, 0x53EF, 0x5CF6, 0x65B6, 0x6E9E, 0x75F9, 1628, 2334, 4115, 6036, 7818, 9544, 0x2E01, 0x36C5, 0x3DAB, 0x4400, 0x4C0A, 0x530D, 0x5966, 0x5FF5, 0x685A, 0x6D9B, 1724, 2670, 4056, 6532, 8357, 0x2787, 0x2F3D, 0x36ED, 0x406B, 0x496B, 0x4FC1, 0x5782, 0x5EBB, 0x6670, 28410, 0x7504, 1493, 3427, 4789, 6399, 8435, 0x27B8, 12000, 0x36F2, 0x3F65, 18210, 20040, 0x5652, 0x5E59, 0x65EF, 0x6E17, 0x75A9, 1119, 2089, 4295, 6245, 8691, 0x29F5, 0x3190, 0x3AD1, 0x4284, 0x4968, 0x50ED, 0x57F2, 0x5FB1, 0x67B4, 0x6FCB, 30630, 1363, 2417, 3927, 5556, 7422, 9315, 0x2E67, 0x35C7, 0x3F0F, 18520, 0x4FEA, 0x5832, 0x5FDB, 0x6744, 0x6E9E, 0x766E, 1122, 2503, 5216, 7148, 9310, 0x2B46, 0x3377, 14800, 0x41E0, 18700, 0x4FD4, 0x57D8, 0x5FFC, 0x67EA, 0x6F8B, 0x76DA, 600, 1317, 2970, 5609, 7694, 9784, 0x2F89, 0x3707, 0x3FFB, 0x47CA, 0x5047, 0x589E, 0x60A3, 0x6849, 0x6FE6, 0x7693, 941, 1882, 4274, 5540, 8482, 9858, 11940, 0x37CF, 0x3EDB, 0x4845, 0x4F66, 0x5854, 0x6087, 0x680E, 0x708E, 30430, 635, 1699, 4376, 5948, 8097, 0x2783, 0x2FF2, 0x3762, 0x3EEF, 0x4595, 0x4CEF, 0x550D, 0x5D77, 0x650A, 0x6D76, 0x75B6, 1408, 2222, 3524, 5615, 7345, 8849, 0x2AED, 0x31E4, 0x3BF8, 0x4282, 0x49E7, 0x5246, 0x5B21, 0x627F, 0x6A49, 0x715F, 701, 1307, 3548, 6301, 7744, 9574, 0x2BDB, 0x32B2, 15170, 0x449D, 0x4D3F, 0x5651, 24230, 0x66DF, 0x6ED9, 0x7617, 1752, 2364, 4879, 6569, 7813, 9796, 0x2BBF, 14290, 0x3DB3, 18000, 0x4FAC, 0x5791, 0x5EF4, 0x660C, 28360, 0x77A9, 901, 1629, 3356, 4635, 7256, 8767, 9971, 0x2D26, 0x3B6F, 0x4488, 0x4C43, 0x555C, 23900, 0x657A, 0x6DE5, 0x75E8, 981, 1669, 3323, 4693, 6213, 8692, 0x2976, 0x329C, 0x3B6B, 0x452F, 0x4D90, 0x566A, 0x5F18, 0x67E0, 0x7033, 0x7711, 1607, 2577, 4220, 5512, 8532, 0x2894, 0x2D6B, 0x3567, 0x3D88, 0x432F, 0x4D80, 0x5563, 0x5BC6, 0x64BA, 0x6DBB, 0x75B3, 811, 1471, 3144, 5041, 7430, 9389, 0x2BA6, 0x33C7, 0x3B35, 0x4165, 0x4C7F, 0x5697, 0x5E33, 0x661E, 0x6EDF, 0x76BB, 1543, 2144, 3629, 6347, 7333, 9339, 10710, 0x351C, 0x3AFB, 17340, 0x4E86, 0x557E, 0x5CB4, 0x6425, 0x6CAA, 0x74DD, 492, 1185, 2940, 5488, 7095, 8751, 0x2D4C, 0x350B, 0x3EAD, 0x465F, 0x4ED2, 0x566F, 0x5EC9, 0x6726, 0x6F44, 0x7695, 1547, 2282, 3693, 6341, 7758, 9607, 0x2E48, 0x33B4, 0x40B4, 0x4695, 0x4D2F, 0x539C, 24110, 0x67EE, 0x7072, 0x77BF, 685, 1338, 3409, 5262, 6950, 9222, 0x2C96, 0x38BB, 0x3FD1, 0x45E5, 0x4BEC, 0x5332, 0x5AFD, 0x625D, 0x6D45, 30520, 887, 1581, 3057, 4318, 7192, 8617, 0x273F, 0x3332, 0x3F89, 0x45E5, 0x4F09, 22350, 0x5F3B, 0x6710, 0x6E9A, 0x75ED, 2285, 3745, 5662, 7576, 9323, 11320, 0x33B7, 0x3B57, 0x4317, 0x4B19, 0x5274, 0x59BC, 0x60F5, 0x681F, 0x6F91, 30460, 1496, 2108, 3448, 6898, 8328, 9656, 0x2BF4, 0x3217, 0x3A83, 0x4062, 18180, 0x4E75, 0x59B2, 25160, 0x6C39, 0x73BD, 575, 1261, 3861, 6627, 8294, 0x2A39, 0x31A1, 0x39B0, 0x42B4, 0x4A67, 0x51F2, 0x5A0F, 0x618C, 0x684F, 0x7030, 0x7689, 1682, 2213, 3882, 6238, 7208, 9646, 0x2A7D, 0x3477, 0x39D5, 0x3F55, 0x4615, 0x5189, 23550, 0x64A5, 0x6C6C, 0x7315, 888, 1616, 3924, 5195, 7206, 8647, 9842, 0x2CD1, 0x3EC3, 0x472D, 0x4F77, 0x58F6, 0x5FB7, 0x672C, 0x6D96, 0x7423, 805, 1454, 2683, 4472, 7936, 9360, 0x2C86, 0x3809, 0x3F4D, 0x45A8, 0x4BFD, 0x548E, 0x5D5B, 0x6548, 0x6EE3, 0x76FF, 1640, 2383, 3484, 5082, 6032, 8606, 11640, 0x32A6, 0x3DE2, 0x43D8, 0x4B92, 0x52BE, 0x5C56, 0x6521, 0x6ED0, 0x765B, 1632, 2204, 4510, 7580, 8718, 0x2910, 0x2EBA, 0x3710, 15640, 0x432A, 0x4AC7, 0x56E7, 0x5FF3, 0x67C1, 0x6FBC, 0x772D, 2043, 2612, 3985, 6851, 8038, 9514, 0x2AE3, 0x31F5, 0x3C42, 0x4158, 0x49D3, 0x4F35, 0x5976, 0x6661, 0x7027, 0x779A, 2224, 2798, 4465, 5320, 7108, 9436, 0x2AEA, 0x33A6, 0x3907, 0x478D, 0x4EAD, 0x5553, 0x5C31, 25700, 0x6E18, 0x7776, 835, 1541, 4083, 5769, 7386, 9399, 0x2ADB, 0x30A8, 0x3AAD, 0x48D2, 0x516B, 23100, 0x62CC, 0x6956, 0x7118, 0x76D6, 1795, 2343, 4809, 5896, 7178, 8545, 0x27EF, 13370, 0x390E, 0x4055, 0x4761, 0x5100, 0x5C5D, 0x6691, 0x6E40, 30390, 1734, 2254, 4031, 5188, 6506, 7872, 9651, 0x32E1, 0x3853, 0x4399, 0x4C27, 22190, 0x5F53, 0x66BE, 0x6E23, 0x75E1, 1841, 2349, 3968, 4764, 6376, 9825, 0x2B28, 0x3421, 0x395A, 0x3F7C, 0x4707, 0x5373, 0x5D6E, 0x662C, 0x6D7F, 0x74EF, 1432, 2047, 5631, 6927, 8198, 9675, 0x2C5E, 0x34C2, 0x39D2, 0x4023, 0x47A3, 0x5603, 0x5E3C, 0x6641, 28130, 0x777A, 1730, 2320, 3744, 4808, 6007, 9666, 0x2AF5, 0x3536, 0x3B82, 0x4457, 0x4E78, 0x55F2, 0x5C33, 25400, 0x6AF3, 0x7246, 1267, 1915, 5483, 6812, 8229, 9919, 0x2D45, 0x3419, 0x399B, 0x462D, 0x5048, 0x5697, 0x5FC7, 0x68C3, 0x70D3, 0x77B2, 1526, 2229, 4240, 7388, 8953, 10450, 0x2E7B, 0x3596, 0x41DD, 0x4793, 0x4F9B, 0x5890, 0x60DD, 0x691A, 0x70EA, 0x779E, 2175, 2791, 4104, 6875, 8612, 9798, 0x2F78, 0x34E0, 0x3D07, 0x4512, 0x4B0D, 21060, 0x5F3E, 26760, 0x6FD9, 0x7628, 454, 1231, 4339, 5738, 7550, 9006, 0x2850, 0x34D5, 0x3E85, 0x45B9, 0x4E67, 0x55E8, 0x5D8D, 0x65BB, 0x6E55, 0x75DF, 2250, 2791, 4230, 5283, 6762, 0x296F, 0x2E67, 0x35FD, 0x3DB5, 0x4370, 0x4E3D, 0x56FA, 0x600C, 0x6745, 0x6E54, 0x76D3, 1696, 2216, 4308, 8385, 9766, 11030, 0x310C, 0x3713, 0x3FC2, 17640, 0x4ADE, 20590, 0x5D9F, 0x68EA, 0x707E, 0x7762, 2452, 3236, 4369, 6118, 7156, 9003, 0x2CF5, 0x31FC, 0x3D85, 0x438B, 0x4C23, 0x56E1, 24530, 0x676A, 0x6E71, 0x7579, 1811, 2541, 3555, 5480, 9123, 0x291F, 0x2E76, 0x355B, 0x3B9E, 0x4203, 0x4BA6, 0x524D, 0x58A6, 0x5EFA, 0x6A78, 0x751F, 1553, 2246, 4559, 5500, 6754, 7874, 0x2DDB, 0x3503, 0x3B54, 0x45D7, 0x4F39, 22510, 0x6026, 0x6819, 0x7072, 0x7823, 1982, 2768, 3834, 5964, 8732, 9908, 0x2E15, 0x39DD, 0x3FB7, 0x461A, 0x5269, 0x5943, 0x5F88, 0x66C0, 0x6E06, 0x743B, 1824, 2529, 3817, 5449, 6854, 8714, 0x288D, 0x2FFE, 0x3772, 0x3D9E, 0x4C44, 0x537E, 0x5C8F, 0x65D5, 0x6DC0, 0x7604, 2212, 2854, 3947, 5898, 9930, 0x2D24, 0x3236, 0x39C4, 0x3FC8, 17700, 0x4F61, 0x5652, 0x5C78, 0x62CB, 0x6960, 0x6FAA, 2023, 2599, 4024, 4916, 6613, 0x2B8D, 0x30A9, 0x3922, 0x3FC0, 0x459E, 0x4CD9, 0x52B4, 0x5A4B, 0x65C3, 0x7099, 0x7826, 1628, 2206, 3467, 4364, 8679, 0x27BD, 0x2E58, 0x356F, 0x3A96, 0x422A, 0x4B07, 0x5374, 23850, 0x6603, 0x6DDC, 0x7641, 2014, 2603, 4114, 7254, 8516, 0x273B, 0x2E2E, 0x34BF, 0x3FC9, 0x45A2, 0x4CF1, 0x5320, 0x5A6F, 0x6055, 0x68B7, 0x75D1, 2376, 2980, 4422, 5770, 7016, 9723, 0x2B75, 0x34CC, 0x3C7D, 0x4259, 19160, 0x506B, 0x5F51, 27180, 0x7176, 0x77B7, 2454, 3502, 4624, 6019, 7632, 8849, 0x2A28, 0x368C, 0x3CA3, 0x42BD, 0x4C9B, 0x52F6, 0x5948, 0x6214, 0x6DCA, 29890, 1573, 2274, 3308, 5999, 8977, 0x2778, 0x30A9, 0x37B2, 0x3D85, 18180, 0x4E06, 0x5305, 0x5A05, 0x61E2, 0x6C5D, 0x766B, 1943, 2730, 4140, 6160, 7491, 8986, 0x2C2D, 0x31E7, 14820, 0x40AE, 0x45F5, 0x4D2D, 0x5408, 0x5C35, 0x6A8A, 0x7357, 2021, 2582, 4494, 5835, 6993, 8245, 9827, 0x398D, 0x404E, 0x45E6, 0x4CBF, 0x525B, 0x5CD4, 0x682B, 0x7190, 30990, 1052, 1775, 3218, 4378, 7666, 9403, 0x2BF0, 0x340F, 0x3A7C, 0x462A, 0x5116, 0x5752, 0x61EF, 0x6A49, 0x7149, 0x7791, 2218, 2866, 4223, 5352, 6581, 9980, 0x2D43, 0x3341, 0x3B59, 0x40C7, 0x47D2, 0x4E70, 0x55FD, 0x62E5, 0x6DDF, 29880, 2146, 2840, 4397, 5840, 7449, 8721, 0x2910, 0x2EA0, 0x351B, 0x4365, 19310, 0x519B, 0x5B79, 0x641B, 0x6C65, 0x7617, 1972, 2619, 3756, 6367, 7641, 8814, 0x2FFE, 0x35C8, 0x3BCD, 0x4674, 0x4C65, 0x51A8, 0x5836, 0x612C, 27800, 30440, 2005, 2577, 4272, 7373, 8558, 0x27EF, 11770, 0x345A, 0x4076, 18000, 0x4CBD, 0x5270, 22990, 0x68B6, 0x7341, 0x78DE, 1153, 1822, 3724, 5443, 6990, 8702, 0x2831, 0x2E7B, 0x3620, 0x3BD3, 0x44C1, 0x5248, 0x5C8C, 0x65E3, 0x6FAA, 0x77AF, 1304, 1869, 3318, 7195, 9613, 0x29ED, 0x3069, 0x35A0, 0x3DCE, 0x4442, 0x49C2, 0x50D4, 0x5A4A, 25540, 0x6C24, 0x723C, 2093, 2691, 4018, 6658, 7947, 9147, 0x2901, 0x2E69, 0x3E10, 0x459D, 0x4B85, 0x52F1, 0x5B4B, 0x6292, 0x6BA1, 0x752E, 575, 1331, 5304, 6910, 8425, 0x2766, 0x2D39, 0x34BA, 0x403C, 0x485F, 0x5055, 0x593F, 0x6152, 0x6844, 0x7057, 0x75CD, 1435, 2024, 3283, 4156, 7611, 0x2960, 0x2F11, 0x3667, 0x3C63, 0x47ED, 0x500F, 22270, 0x5E9E, 0x65ED, 0x6DA1, 0x7593, 1632, 2168, 5540, 7478, 8630, 0x2897, 0x2D7C, 0x37F1, 0x3D7D, 0x43CD, 0x4944, 0x4FD2, 0x590F, 26060, 0x6F7E, 0x77E8, 1407, 2245, 3405, 5639, 9419, 0x29BD, 0x2F48, 0x34B7, 0x3CAF, 0x47B5, 0x4E1C, 0x54B9, 0x5F1F, 26550, 0x70B5, 0x7764, 1675, 2226, 4005, 8223, 9975, 0x2B93, 0x3216, 0x37EC, 0x4078, 0x46D9, 0x4C76, 21050, 0x58E7, 0x6150, 0x6E88, 0x77AA, 1080, 1614, 3622, 7565, 8748, 0x283F, 0x2DC1, 0x3618, 0x3D11, 0x441A, 0x4D31, 0x5541, 0x5C13, 0x6331, 0x6B0E, 0x7187, 1693, 2229, 3456, 4354, 5670, 10890, 0x3113, 0x3757, 0x3E07, 0x43E1, 0x4D69, 0x55D3, 0x5E1E, 0x6613, 0x6E8A, 0x7593, 2042, 2959, 4195, 5740, 7106, 8267, 0x2B76, 0x3A7D, 0x4212, 0x4777, 0x5034, 0x55DE, 0x5C9F, 0x64A9, 0x6BD9, 0x72A7, 984, 1612, 3808, 5265, 6885, 8411, 9547, 0x2A89, 0x30EA, 16520, 0x4C5D, 0x5487, 0x5CC2, 0x65CA, 28310, 0x76A6, 2036, 2538, 4166, 7761, 9146, 0x28AC, 0x2F70, 0x3529, 0x3CE4, 0x4311, 0x487F, 0x4E91, 21820, 0x5EF9, 0x6D7D, 0x7794, 1871, 2355, 4061, 5143, 7464, 0x2791, 0x2EA5, 0x3A99, 16680, 0x47B2, 0x4DF5, 0x5707, 0x611D, 0x68F8, 0x713C, 0x7797, 2566, 3161, 4643, 6227, 7406, 9970, 0x2D62, 0x3468, 0x3E11, 0x43D4, 0x4AB1, 0x5151, 0x5840, 0x6090, 0x703D, 0x796A, 1700, 2327, 4828, 5939, 7567, 9154, 0x2B4F, 0x31E3, 0x3781, 0x3EF9, 0x4EFE, 0x588F, 0x6048, 0x6820, 0x7018, 0x7819, 3169, 3873, 5046, 6868, 8184, 9480, 0x302F, 0x36F4, 0x3D9E, 0x4633, 0x4F07, 0x54CF, 0x5BE0, 0x629D, 0x6992, 28730, 1564, 2391, 4229, 6730, 8905, 0x28DB, 0x32E2, 0x3AB9, 0x4371, 0x4D61, 0x5559, 0x5CBD, 25490, 0x6AB0, 0x7185, 0x773F, 2864, 3559, 4719, 6441, 9592, 0x2B2F, 0x31DB, 0x39C0, 0x402C, 0x46F4, 0x5006, 0x56F6, 0x5E77, 0x6697, 0x6EDF, 0x7610, 2673, 3449, 4581, 5983, 6863, 8311, 0x30B0, 0x3657, 0x3D7A, 0x457F, 0x4BD8, 0x52BE, 0x5DD9, 0x67C1, 0x7033, 30440, 2419, 3049, 4274, 6384, 8564, 9661, 0x2C18, 0x3184, 0x386F, 0x44AA, 0x4D68, 0x52EF, 0x5A3B, 25270, 0x6913, 0x70FE, 1278, 2001, 3000, 5353, 9995, 0x2E01, 0x32DA, 14570, 16050, 0x4562, 0x4E0E, 0x5471, 0x5B4B, 0x61FB, 0x6C08, 0x75DC, 932, 1624, 2798, 4570, 8592, 9988, 0x2D20, 13050, 0x4219, 0x48F5, 0x4FBF, 22810, 0x60F1, 0x68C3, 0x7084, 0x76B1, 2324, 2973, 4156, 5702, 6919, 8806, 0x2813, 0x30D7, 0x3AA7, 0x40B7, 0x4BDA, 0x537F, 0x599F, 24550, 0x6990, 0x7499, 1564, 2373, 3455, 4907, 5975, 7436, 0x2E0A, 0x38A9, 0x3EEB, 0x46E4, 0x4E33, 0x5495, 23740, 0x64D6, 0x6FA2, 0x76A4, 3025, 3729, 4866, 6520, 9487, 0x2ABF, 0x3046, 0x37B2, 0x3F2E, 0x445D, 0x4C14, 0x53A0, 0x5ABB, 0x614A, 0x6AD3, 0x72DF, 1270, 1965, 6802, 7995, 9204, 0x2A4C, 0x30DB, 14230, 0x3D8F, 17860, 0x4F91, 0x57E6, 0x6039, 0x6792, 0x6F77, 0x773D, 2210, 2749, 4266, 7487, 9878, 0x2B0A, 0x3217, 0x385F, 0x3F77, 0x48C2, 20450, 0x5626, 0x5CBB, 0x62CB, 0x69C2, 0x71F1, 1275, 1926, 4330, 6573, 8441, 10920, 13260, 0x3AA0, 0x421F, 0x488D, 0x50A4, 0x56C9, 0x5DAF, 0x6382, 0x6AEC, 0x6FE5, 3015, 3670, 5086, 6372, 7888, 9309, 0x2AD6, 0x3162, 0x389F, 0x3F2C, 0x46A0, 0x4E04, 0x57B6, 0x6143, 0x6AE2, 0x7517, 2882, 3733, 5113, 6482, 8125, 9685, 0x2D4E, 0x33E8, 0x3C2D, 0x4328, 0x4ED2, 0x579A, 0x60E1, 0x6986, 0x721C, 0x785B, 2300, 2968, 4101, 5442, 6327, 7910, 0x30A7, 0x3626, 0x3D83, 0x4461, 0x4A6D, 0x50C7, 0x5857, 0x6052, 0x6B6B, 0x7571, 2257, 2940, 4430, 5991, 7042, 8364, 9414, 0x2BD8, 0x3D6B, 17420, 0x4B35, 0x53DD, 0x5D6B, 0x65C5, 28430, 0x76B0, 1227, 2045, 3818, 5011, 6990, 9231, 0x2B10, 0x32D3, 0x43BD, 0x4A49, 0x5067, 0x590F, 0x626B, 0x68FC, 0x72A7, 0x7855, 1354, 1924, 3789, 8077, 0x28D5, 0x2D77, 0x3428, 0x39E1, 0x4167, 0x470D, 0x4E7F, 0x55FE, 0x6011, 0x6835, 0x6FE7, 0x7630, 3142, 4049, 6197, 7417, 8753, 0x27AC, 0x2D0D, 0x337D, 0x3E4B, 0x44F7, 0x4C96, 0x539A, 0x5BBF, 0x643B, 0x6DDB, 0x7660, 1317, 2263, 4725, 7611, 9667, 0x2D72, 0x373F, 0x3F82, 0x4924, 0x50DA, 0x576B, 0x5DC7, 0x64AF, 0x6A73, 28930, 0x7781, 1570, 2323, 3818, 6215, 9893, 0x2D24, 13070, 0x3927, 0x3F18, 18290, 0x538A, 0x5B32, 0x621A, 0x692B, 0x7028, 0x75D8, 2297, 3905, 6287, 8558, 0x29AC, 0x31DE, 0x3AAB, 0x42CE, 0x4A5C, 0x50C5, 0x5745, 0x5D3F, 0x6386, 0x69CD, 0x70B3, 30520, 1915, 2507, 4033, 5749, 7059, 8871, 0x29A3, 0x2FA6, 0x3671, 0x3C17, 0x41E5, 0x4913, 0x5A87, 0x64DA, 0x6F62, 0x7725, 2404, 2918, 5190, 6252, 7426, 9887, 0x3063, 0x39CB, 0x4172, 0x47C0, 0x4F72, 0x55F3, 0x5EAC, 0x6758, 28490, 0x76BD, 1621, 2227, 3479, 5085, 9425, 0x325C, 0x37A6, 0x3D24, 0x4335, 0x48F2, 0x4FDE, 0x56C1, 0x5CE2, 0x650B, 0x6D1B, 0x758D, 1869, 2390, 4105, 7021, 0x2BD5, 0x31E7, 0x36EB, 15590, 0x4280, 0x48B0, 0x5073, 0x563B, 0x5C61, 0x6242, 0x6922, 0x6FFF, 2551, 3252, 4688, 6562, 7869, 9125, 0x28EB, 11800, 0x3C2A, 18780, 0x5200, 0x581B, 0x5EE1, 0x6570, 0x6B49, 0x7230, 2705, 3493, 4735, 6360, 7905, 9352, 0x2D12, 13430, 0x3B87, 0x4217, 0x48BB, 0x4E7E, 21800, 0x5B2E, 25200, 0x7249, 2166, 2791, 4011, 5081, 5896, 9038, 0x345F, 0x396F, 0x409F, 0x470D, 0x4DB8, 0x5561, 0x6128, 0x695B, 0x711B, 0x7732, 1865, 3021, 4696, 6534, 8343, 9914, 0x31F5, 0x3717, 0x4095, 0x4541, 21340, 0x57A7, 0x6129, 26330, 0x6F0C, 0x75CA, 3369, 4345, 6573, 8763, 0x2845, 0x2DC1, 0x3437, 0x39C0, 0x4063, 0x46E1, 0x4D7F, 0x52FF, 0x5AFC, 0x6385, 0x6BA3, 0x7307, 1265, 2184, 5443, 7893, 0x295F, 0x3353, 0x3B01, 0x40FF, 0x47E2, 0x4D72, 0x53AB, 0x59D3, 0x608F, 0x6745, 0x6ECB, 0x75AD, 1584, 2004, 3535, 4450, 8662, 0x2A0C, 0x3220, 0x3A82, 0x424C, 0x496A, 0x51C4, 0x5813, 0x603C, 0x6799, 0x701D, 0x7767, 3419, 4528, 6602, 7890, 9508, 0x2A7B, 0x31E3, 0x3815, 0x3EB3, 18330, 20630, 22490, 25070, 0x6938, 0x7112, 0x774E, 1726, 2252, 4597, 6950, 8379, 9823, 0x2C63, 0x31FA, 0x37E2, 0x3C74, 0x419E, 0x4662, 0x54A7, 25550, 0x6DF4, 0x769F, 3385, 3870, 5307, 6388, 7141, 8684, 0x3197, 0x3A5B, 0x4060, 0x4765, 0x5039, 0x5620, 0x5D8B, 0x656D, 0x6E36, 0x7504, 2771, 3306, 4450, 5560, 6453, 9493, 0x34EC, 0x39A2, 0x4167, 0x480F, 0x4E3C, 0x54E8, 0x5CC2, 0x6309, 0x6A05, 0x718A, 3028, 3900, 6617, 7893, 9211, 0x28F0, 0x2F0F, 0x350F, 0x3B4E, 0x4116, 0x4846, 0x4E7C, 22190, 0x5F26, 0x66BE, 0x711D, 2000, 2550, 4067, 6837, 9628, 0x2AFA, 0x3132, 0x3712, 0x3CE5, 0x432B, 0x48F7, 0x4E83, 21530, 0x5A2D, 0x6041, 0x715E, 2844, 3302, 5103, 6107, 6911, 8598, 0x3080, 0x36E6, 0x3E9A, 0x4887, 0x50C0, 22270, 0x5D90, 0x64AB, 0x6C0A, 0x754A, 4043, 5150, 7268, 9056, 0x2AA4, 0x315E, 0x38CF, 0x3F38, 0x461C, 0x4CEB, 0x536D, 0x59C5, 0x60F9, 0x67DF, 0x6F3F, 0x7619, 2109, 2625, 4320, 5525, 7454, 10220, 12980, 0x396A, 0x44DB, 0x4B3F, 0x5005, 0x576D, 0x5ED7, 0x64B1, 0x6CC7, 0x76FA, 1550, 2667, 6473, 9496, 0x2AE9, 0x3040, 0x35E3, 0x3B81, 0x42CB, 0x48D2, 0x4FED, 0x5664, 0x5E85, 0x66B3, 0x6EF3, 0x75B4, 2411, 3084, 4145, 5394, 6367, 8154, 0x3345, 0x3EB1, 0x4499, 0x4AB5, 0x530A, 0x58EA, 0x5F8B, 0x66CD, 0x6E5F, 0x7406, 4159, 4516, 5956, 7635, 8254, 8980, 0x2BC8, 0x3735, 16210, 0x45D3, 0x4EE4, 0x5568, 0x5D20, 0x6493, 0x6D9A, 0x753C, 2026, 2431, 2845, 3618, 7950, 9802, 0x31B1, 14460, 0x40C0, 0x4A28, 0x5380, 0x5B17, 0x6181, 0x685E, 0x712B, 0x77B0, 3429, 3833, 4472, 4912, 7723, 0x2892, 0x32B5, 0x3BDA, 0x413B, 0x4977, 0x512A, 0x5817, 0x6033, 0x677E, 0x6EAE, 0x7712, 4740, 5169, 5796, 6485, 6998, 8830, 0x2E01, 0x384E, 0x41BF, 0x47ED, 0x5135, 0x5761, 0x5EAC, 0x64EB, 0x6C9F, 0x7545, 150, 168, -17, -107, -142, -229, -320, -406, -503, -620, -867, -935, -902, -680, -398, -114, -398, -355, 49, math.MaxUint8, 114, 260, 399, 264, 317, 431, 514, 531, 435, 356, 238, 106, -43, -36, -169, -224, -391, -633, -776, -970, -844, -455, -181, -12, 85, 85, 164, 195, 122, 85, -158, -640, -903, 9, 7, -124, 149, 32, 220, 369, 242, 115, 79, 84, -146, -216, -70, 1024, 751, 574, 440, 377, 352, 203, 30, 16, -3, 81, 161, 100, -148, -176, 933, 750, 404, 171, -2, -146, -411, -442, -541, -552, -442, -269, -240, -52, 603, 635, 405, 178, 215, 19, -153, -167, -290, -219, 151, 271, 151, 119, 303, 266, 100, 69, -293, -657, 939, 659, 442, 351, 132, 98, -16, -1, -135, -200, -223, -89, 167, 154, 172, 237, -45, -183, -228, -486, 263, 608, 158, -125, -390, -227, -118, 43, -457, -392, -769, -840, 20, -117, -194, -189, -173, -173, -33, 32, 174, 144, 115, 167, 57, 44, 14, 147, 96, -54, -142, -129, -254, -331, 304, 310, -52, -419, -846, -1060, -88, -123, -202, -343, -554, -961, -951, 327, 159, 81, math.MaxUint8, 227, 120, 203, 256, 192, 164, 224, 290, 195, 216, 209, 128, 832, 1028, 889, 698, 504, 408, 355, 218, 32, -115, -84, -276, -100, -312, -484, 899, 682, 465, 456, 241, -12, -275, -425, -461, -367, -33, -28, -102, -194, -527, 863, 906, 463, 245, 13, -212, -305, -105, 163, 279, 176, 93, 67, 115, 192, 61, -50, -132, -175, -224, -271, -629, -252, 1158, 972, 638, 280, 300, 326, 143, -152, -214, -287, 53, -42, -236, -352, -423, -248, -129, -163, -178, -119, 85, 57, 514, 382, 374, 402, 424, 423, 271, 197, 97, 40, 39, -97, -191, -164, -230, -256, -410, 396, 327, math.MaxInt8, 10, -119, -167, -291, -274, -141, -99, -226, -218, -139, -224, -209, -268, -442, -413, 222, 58, 521, 344, 258, 76, -42, -142, -165, -123, -92, 47, 8, -3, -191, -11, -164, -167, -351, -740, 311, 538, 291, 184, 29, -105, 9, -30, -54, -17, -77, -271, -412, -622, -648, 476, 186, -66, -197, -73, -94, -15, 47, 28, 112, -58, -33, 65, 19, 84, 86, 276, 114, 472, 786, 799, 625, 415, 178, -35, -26, 5, 9, 83, 39, 37, 39, -184, -374, -265, -362, -501, 337, 716, 478, -60, -125, -163, 362, 17, -122, -233, 279, 138, 157, 318, 193, 189, 209, 266, 252, -46, -56, -277, -429, 464, 386, 142, 44, -43, 66, 264, 182, 47, 14, -26, -79, 49, 15, math.MinInt8, -203, -400, -478, 325, 27, 234, 411, 205, 129, 12, 58, 123, 57, 171, 137, 96, 128, -32, 134, -12, 57, 119, 26, -22, -165, -500, -701, -528, -116, 64, -8, 97, -9, -162, -66, -156, -194, -303, -546, -341, 546, 358, 95, 45, 76, 270, 403, 205, 100, 123, 50, -53, -144, -110, -13, 32, -228, -130, 353, 296, 56, -372, -253, 365, 73, 10, -34, -139, -191, -96, 5, 44, -85, -179, -129, -192, -246, -85, -110, -155, -44, -27, 145, 138, 79, 32, -148, -577, -634, 191, 94, -9, -35, -77, -84, -56, -171, -298, -271, -243, -156, -328, -235, -76, math.MinInt8, -121, 129, 13, -22, 32, 45, -248, -65, 193, -81, 299, 57, -147, 192, -165, -354, -334, -106, -156, -40, -3, -68, 124, -257, 78, 124, 170, 412, 227, 105, -104, 12, 154, 250, 274, 258, 4, -27, 235, 152, 51, 338, 300, 7, -314, -411, 215, 170, -9, -93, -77, 76, 67, 54, 200, 315, 163, 72, -91, -402, 158, 187, -156, -91, 290, 267, 167, 91, 140, 171, 112, 9, -42, -177, -440, 385, 80, 15, 172, 129, 41, -129, -372, -24, -75, -30, -170, 10, -118, 57, 78, -101, 232, 161, 123, 256, 277, 101, -192, -629, -100, -60, -232, 66, 13, -13, -80, -239, 239, 37, 32, 89, -319, -579, 450, 360, 3, -29, -299, -89, -54, -110, -246, -164, 6, -188, 338, 176, -92, 197, 137, 134, 12, -2, 56, -183, 114, -36, -131, -204, 75, -25, -174, 191, -15, -290, -429, -267, 79, 37, 106, 23, -384, 425, 70, -14, 212, 105, 15, -2, -42, -37, -123, 108, 28, -48, 193, 197, 173, -33, 37, 73, -57, 256, 137, -58, -430, -228, 217, -51, -10, -58, -6, 22, 104, 61, -119, 169, 144, 16, -46, -394, 60, 454, -80, -298, -65, 25, 0, -24, -65, -417, 465, 276, -3, -194, -13, 130, 19, -6, -21, -24, -180, -53, -85, 20, 118, 147, 113, -75, -289, 226, -122, 227, 270, 125, 109, 197, 125, 138, 44, 60, 25, -55, -167, -32, -139, -193, -173, -316, 287, -208, 253, 239, 27, -80, -188, -28, -182, -235, 156, -117, 128, -48, -58, -226, 172, 181, 167, 19, 62, 10, 2, 181, 151, 108, -16, -11, -78, -331, 411, 133, 17, 104, 64, -184, 24, -30, -3, -283, 121, 204, -8, -199, -21, -80, -169, -157, -191, -136, 81, 155, 14, -131, 244, 74, -57, -47, -280, 347, 111, -77, math.MinInt8, -142, -194, -125, -6, -68, 91, 1, 23, 14, -154, -34, 23, -38, -343, 503, 146, -38, -46, -41, 58, 31, 63, -48, -117, 45, 28, 1, -89, -5, -44, -29, -448, 487, 204, 81, 46, -106, -302, 380, 120, -38, -12, -39, 70, -3, 25, -65, 30, -11, 34, -15, 22, -115, 0, -79, -83, 45, 114, 43, 150, 36, 233, 149, 195, 5, 25, -52, -475, 274, 28, -39, -8, -66, -255, 258, 56, 143, -45, -190, 165, -60, 20, 2, 125, -129, 51, -8, -335, 288, 38, 59, 25, -42, 23, -118, -112, 11, -55, -133, -109, 24, -105, 78, -64, -245, 202, -65, -127, 162, 40, -94, 89, -85, -119, -103, 97, 9, -70, -28, 194, 86, -112, -92, -114, 74, -49, 46, -84, -178, 113, 52, -205, 333, 88, 222, 56, -55, 13, 86, 4, -77, 224, 114, -105, 112, 125, -29, -18, -144, 22, -58, -99, 28, 114, -66, -32, -169, -314, 285, 72, -74, 179, 28, -79, -182, 13, -55, 147, 13, 12, -54, 31, -84, -17, -75, -228, 83, -375, 436, 110, -63, -27, -136, 169, -56, -8, -171, 184, -42, 148, 68, 204, 235, 110, -229, 91, 171, -43, -3, -26, -99, -111, 71, -170, 202, -67, 181, -37, 109, -120, 3, -55, -260, -16, 152, 91, 142, 42, 44, 134, 47, 17, -35, 22, 79, -169, 41, 46, 277, -93, -49, -126, 37, -103, -34, -22, -90, -134, -205, 92, -9, 1, -195, -239, 45, 54, 18, -23, -1, -80, -98, -20, -261, 306, 72, 20, -89, -217, 11, 6, -82, 89, 13, -129, -89, 83, -71, -55, 130, -98, -146, -27, -57, 53, 275, 17, 170, -5, -54, 132, -64, 72, 160, -125, -168, 72, 40, 170, 78, 248, 116, 20, 84, 31, -34, 190, 38, 13, -106, 225, 27, -168, 24, -157, -122, 165, 11, -161, -213, -12, -51, -101, 42, 101, 27, 55, 111, 75, 71, -96, -1, 65, -277, 393, -26, -44, -68, -84, -66, -95, 235, 179, -25, -41, 27, -91, math.MinInt8, -222, 146, -72, -30, -24, 55, -126, -68, -58, -127, 13, -97, -106, 174, -100, 155, 101, -146, -21, 261, 22, 38, -66, 65, 4, 70, 64, 144, 59, 213, 71, -337, 303, -52, 51, -56, 1, 10, -15, -5, 34, 52, 228, 131, 161, -127, -214, 238, 123, 64, -147, -50, -34, -127, 204, 162, 85, 41, 5, -140, 73, -150, 56, -96, -66, -20, 2, -235, 59, -22, -107, 150, -16, -47, -4, 81, -67, 167, 149, 149, -157, 288, -156, -27, -8, 18, 83, -24, -41, -167, 158, -100, 93, 53, 201, 15, 42, 266, 278, -12, -6, -37, 85, 6, 20, -188, -271, 107, -13, -80, 51, 202, 173, -69, 78, -188, 46, 4, 153, 12, -138, 169, 5, -58, -123, -108, -243, 150, 10, -191, 246, -15, 38, 25, -10, 14, 61, 50, -206, -215, -220, 90, 5, -149, -219, 56, 142, 24, -376, 77, -80, 75, 6, 42, -101, 16, 56, 14, -57, 3, -17, 80, 57, -36, 88, -59, -97, -19, -148, 46, -219, 226, 114, -4, -72, -15, 37, -49, -28, 247, 44, 123, 47, -122, -38, 17, 4, -113, -32, -224, 154, -134, 196, 71, -267, -85, 28, -70, 89, -120, 99, -2, 64, 76, -166, -48, 189, -35, -92, -169, -123, 339, 38, -25, 38, -35, 225, -139, -50, -63, 246, 60, -185, -109, -49, -53, -167, 51, 149, 60, -101, -33, 25, -76, 120, 32, -30, -83, 102, 91, -186, -261, 131, -197}
var SKP_Silk_NLSF_CB0_16_Stage_info [10]SKP_Silk_NLSF_CBS = [10]SKP_Silk_NLSF_CBS{{NVectors: 128, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*0], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[0]}, {NVectors: 16, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*128], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[128]}, {NVectors: 8, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*144], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[144]}, {NVectors: 8, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*152], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[152]}, {NVectors: 8, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*160], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[160]}, {NVectors: 8, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*168], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[168]}, {NVectors: 8, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*176], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[176]}, {NVectors: 8, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*184], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[184]}, {NVectors: 8, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*192], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[192]}, {NVectors: 16, CB_NLSF_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_Q15[16*200], Rates_Q5: &SKP_Silk_NLSF_MSVQ_CB0_16_rates_Q5[200]}}
var SKP_Silk_NLSF_CB0_16 SKP_Silk_NLSF_CB_struct = SKP_Silk_NLSF_CB_struct{NStages: NLSF_MSVQ_CB0_16_STAGES, CBStages: &SKP_Silk_NLSF_CB0_16_Stage_info[0], NDeltaMin_Q15: &SKP_Silk_NLSF_MSVQ_CB0_16_ndelta_min_Q15[0], CDF: &SKP_Silk_NLSF_MSVQ_CB0_16_CDF[0], StartPtr: &SKP_Silk_NLSF_MSVQ_CB0_16_CDF_start_ptr[0], MiddleIx: &SKP_Silk_NLSF_MSVQ_CB0_16_CDF_middle_idx[0]}