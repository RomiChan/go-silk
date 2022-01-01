package silk

import (
	"math"
	"unsafe"
)

const NB_ATT = 2
const BWE_COEF_Q16 = 0xFD70
const V_PITCH_GAIN_START_MIN_Q14 = 0x2CCD
const V_PITCH_GAIN_START_MAX_Q14 = 0x3CCD
const MAX_PITCH_LAG_MS = 18
const SA_THRES_Q8 = 50
const USE_SINGLE_TAP = 1
const RAND_BUF_SIZE = 128
const RAND_BUF_MASK uint8 = math.MaxInt8
const LOG2_INV_LPC_GAIN_HIGH_THRES = 3
const LOG2_INV_LPC_GAIN_LOW_THRES = 8
const PITCH_DRIFT_FAC_Q16 = 655

var HARM_ATT_Q15 [2]int16 = [2]int16{32440, 31130}
var PLC_RAND_ATTENUATE_V_Q15 [2]int16 = [2]int16{31130, 0x6666}
var PLC_RAND_ATTENUATE_UV_Q15 [2]int16 = [2]int16{32440, 0x7333}

func SKP_Silk_PLC_Reset(psDec *SKP_Silk_decoder_state) {
	psDec.SPLC.PitchL_Q8 = psDec.Frame_length >> 1
}
func SKP_Silk_PLC(psDec *SKP_Silk_decoder_state, psDecCtrl *SKP_Silk_decoder_control, signal []int16, length int32, lost int32) {
	if psDec.Fs_kHz != psDec.SPLC.Fs_kHz {
		SKP_Silk_PLC_Reset(psDec)
		psDec.SPLC.Fs_kHz = psDec.Fs_kHz
	}
	if lost != 0 {
		SKP_Silk_PLC_conceal(psDec, psDecCtrl, signal, length)
		psDec.LossCnt++
	} else {
		SKP_Silk_PLC_update(psDec, psDecCtrl, signal, length)
	}
}
func SKP_Silk_PLC_update(psDec *SKP_Silk_decoder_state, psDecCtrl *SKP_Silk_decoder_control, signal []int16, length int32) {
	var (
		LTP_Gain_Q14      int32
		temp_LTP_Gain_Q14 int32
		i                 int32
		j                 int32
		psPLC             *SKP_Silk_PLC_struct
	)
	psPLC = &psDec.SPLC
	psDec.Prev_sigtype = psDecCtrl.Sigtype
	LTP_Gain_Q14 = 0
	if psDecCtrl.Sigtype == SIG_TYPE_VOICED {
		for j = 0; j*psDec.Subfr_length < psDecCtrl.PitchL[NB_SUBFR-1]; j++ {
			temp_LTP_Gain_Q14 = 0
			for i = 0; i < LTP_ORDER; i++ {
				temp_LTP_Gain_Q14 += int32(psDecCtrl.LTPCoef_Q14[(NB_SUBFR-1-j)*LTP_ORDER+i])
			}
			if temp_LTP_Gain_Q14 > LTP_Gain_Q14 {
				LTP_Gain_Q14 = temp_LTP_Gain_Q14
				memcpy(unsafe.Pointer(&psPLC.LTPCoef_Q14[0]), unsafe.Pointer(&psDecCtrl.LTPCoef_Q14[SKP_SMULBB(NB_SUBFR-1-j, LTP_ORDER)]), size_t(LTP_ORDER*unsafe.Sizeof(int16(0))))
				psPLC.PitchL_Q8 = (psDecCtrl.PitchL[NB_SUBFR-1-j]) << 8
			}
		}
		memset(unsafe.Pointer(&psPLC.LTPCoef_Q14[0]), 0, size_t(LTP_ORDER*unsafe.Sizeof(int16(0))))
		psPLC.LTPCoef_Q14[LTP_ORDER/2] = int16(LTP_Gain_Q14)
		if LTP_Gain_Q14 < V_PITCH_GAIN_START_MIN_Q14 {
			var (
				scale_Q10 int32
				tmp       int32
			)
			tmp = V_PITCH_GAIN_START_MIN_Q14 << 10
			scale_Q10 = tmp / (func() int32 {
				if LTP_Gain_Q14 > 1 {
					return LTP_Gain_Q14
				}
				return 1
			}())
			for i = 0; i < LTP_ORDER; i++ {
				psPLC.LTPCoef_Q14[i] = int16(SKP_SMULBB(int32(psPLC.LTPCoef_Q14[i]), scale_Q10) >> 10)
			}
		} else if LTP_Gain_Q14 > V_PITCH_GAIN_START_MAX_Q14 {
			var (
				scale_Q14 int32
				tmp       int32
			)
			tmp = V_PITCH_GAIN_START_MAX_Q14 << 14
			scale_Q14 = tmp / (func() int32 {
				if LTP_Gain_Q14 > 1 {
					return LTP_Gain_Q14
				}
				return 1
			}())
			for i = 0; i < LTP_ORDER; i++ {
				psPLC.LTPCoef_Q14[i] = int16(SKP_SMULBB(int32(psPLC.LTPCoef_Q14[i]), scale_Q14) >> 14)
			}
		}
	} else {
		psPLC.PitchL_Q8 = SKP_SMULBB(psDec.Fs_kHz, 18) << 8
		memset(unsafe.Pointer(&psPLC.LTPCoef_Q14[0]), 0, size_t(LTP_ORDER*unsafe.Sizeof(int16(0))))
	}
	memcpy(unsafe.Pointer(&psPLC.PrevLPC_Q12[0]), unsafe.Pointer(&(psDecCtrl.PredCoef_Q12[1])[0]), size_t(uintptr(psDec.LPC_order)*unsafe.Sizeof(int16(0))))
	psPLC.PrevLTP_scale_Q14 = int16(psDecCtrl.LTP_scale_Q14)
	memcpy(unsafe.Pointer(&psPLC.PrevGain_Q16[0]), unsafe.Pointer(&psDecCtrl.Gains_Q16[0]), size_t(NB_SUBFR*unsafe.Sizeof(int32(0))))
}
func SKP_Silk_PLC_conceal(psDec *SKP_Silk_decoder_state, psDecCtrl *SKP_Silk_decoder_control, signal []int16, length int32) {
	var (
		i           int32
		j           int32
		k           int32
		B_Q14       []int16
		exc_buf     [480]int16
		exc_buf_ptr []int16
	)
	_ = exc_buf_ptr
	var rand_scale_Q14 int16
	var A_Q12_tmp struct {
		// union
		As_int16 [16]int16
		As_int32 [8]int32
	}
	var rand_seed int32
	var harm_Gain_Q15 int32
	var rand_Gain_Q15 int32
	var lag int32
	var idx int32
	var sLTP_buf_idx int32
	var shift1 int32
	var shift2 int32
	var energy1 int32
	var energy2 int32
	var rand_ptr *int32
	var pred_lag_ptr *int32
	var sig_Q10 [480]int32
	var sig_Q10_ptr []int32
	var LPC_exc_Q10 int32
	var LPC_pred_Q10 int32
	var LTP_pred_Q14 int32
	var psPLC *SKP_Silk_PLC_struct
	psPLC = &psDec.SPLC
	memcpy(unsafe.Pointer(&psDec.SLTP_Q16[0]), unsafe.Pointer(&psDec.SLTP_Q16[psDec.Frame_length]), size_t(uintptr(psDec.Frame_length)*unsafe.Sizeof(int32(0))))
	SKP_Silk_bwexpander(psPLC.PrevLPC_Q12[:], psDec.LPC_order, BWE_COEF_Q16)
	exc_buf_ptr = ([]int16)(exc_buf[:])
	for k = NB_SUBFR >> 1; k < NB_SUBFR; k++ {
		for i = 0; i < psDec.Subfr_length; i++ {
			exc_buf_ptr[i] = int16(SKP_SMULWW(psDec.Exc_Q10[i+k*psDec.Subfr_length], psPLC.PrevGain_Q16[k]) >> 10)
		}
		exc_buf_ptr += ([]int16)(psDec.Subfr_length)
	}
	SKP_Silk_sum_sqr_shift(&energy1, &shift1, exc_buf[:], psDec.Subfr_length)
	SKP_Silk_sum_sqr_shift(&energy2, &shift2, ([]int16)(&exc_buf[psDec.Subfr_length]), psDec.Subfr_length)
	if (energy1 >> shift2) < (energy2 >> shift1) {
		rand_ptr = &psDec.Exc_Q10[SKP_max_int(0, psDec.Subfr_length*3-RAND_BUF_SIZE)]
	} else {
		rand_ptr = &psDec.Exc_Q10[SKP_max_int(0, psDec.Frame_length-RAND_BUF_SIZE)]
	}
	B_Q14 = ([]int16)(psPLC.LTPCoef_Q14[:])
	rand_scale_Q14 = psPLC.RandScale_Q14
	harm_Gain_Q15 = int32(HARM_ATT_Q15[SKP_min_int(NB_ATT-1, psDec.LossCnt)])
	if psDec.Prev_sigtype == SIG_TYPE_VOICED {
		rand_Gain_Q15 = int32(PLC_RAND_ATTENUATE_V_Q15[SKP_min_int(NB_ATT-1, psDec.LossCnt)])
	} else {
		rand_Gain_Q15 = int32(PLC_RAND_ATTENUATE_UV_Q15[SKP_min_int(NB_ATT-1, psDec.LossCnt)])
	}
	if psDec.LossCnt == 0 {
		rand_scale_Q14 = 1 << 14
		if psDec.Prev_sigtype == SIG_TYPE_VOICED {
			for i = 0; i < LTP_ORDER; i++ {
				rand_scale_Q14 -= B_Q14[i]
			}
			rand_scale_Q14 = SKP_max_16(3277, rand_scale_Q14)
			rand_scale_Q14 = int16(SKP_SMULBB(int32(rand_scale_Q14), int32(psPLC.PrevLTP_scale_Q14)) >> 14)
		}
		if psDec.Prev_sigtype == SIG_TYPE_UNVOICED {
			var (
				invGain_Q30    int32
				down_scale_Q30 int32
			)
			SKP_Silk_LPC_inverse_pred_gain(&invGain_Q30, psPLC.PrevLPC_Q12[:], psDec.LPC_order)
			down_scale_Q30 = SKP_min_32((1<<30)>>LOG2_INV_LPC_GAIN_HIGH_THRES, invGain_Q30)
			down_scale_Q30 = SKP_max_32((1<<30)>>LOG2_INV_LPC_GAIN_LOW_THRES, down_scale_Q30)
			down_scale_Q30 = down_scale_Q30 << LOG2_INV_LPC_GAIN_HIGH_THRES
			rand_Gain_Q15 = SKP_SMULWB(down_scale_Q30, rand_Gain_Q15) >> 14
		}
	}
	rand_seed = psPLC.Rand_seed
	lag = SKP_RSHIFT_ROUND(psPLC.PitchL_Q8, 8)
	sLTP_buf_idx = psDec.Frame_length
	sig_Q10_ptr = ([]int32)(sig_Q10[:])
	for k = 0; k < NB_SUBFR; k++ {
		pred_lag_ptr = &psDec.SLTP_Q16[sLTP_buf_idx-lag+LTP_ORDER/2]
		for i = 0; i < psDec.Subfr_length; i++ {
			rand_seed = int32((uint32(rand_seed) * 0xBB38435) + 0x3619636B)
			idx = (rand_seed >> 25) & (RAND_BUF_SIZE - 1)
			LTP_pred_Q14 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(int32(0))*0)), int32(B_Q14[0]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*1))), int32(B_Q14[1]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*2))), int32(B_Q14[2]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*3))), int32(B_Q14[3]))
			LTP_pred_Q14 = SKP_SMLAWB(LTP_pred_Q14, *(*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(int32(0))*4))), int32(B_Q14[4]))
			pred_lag_ptr = (*int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(int32(0))*1))
			LPC_exc_Q10 = SKP_SMULWB(*(*int32)(unsafe.Add(unsafe.Pointer(rand_ptr), unsafe.Sizeof(int32(0))*uintptr(idx))), int32(rand_scale_Q14)) << 2
			LPC_exc_Q10 = LPC_exc_Q10 + SKP_RSHIFT_ROUND(LTP_pred_Q14, 4)
			psDec.SLTP_Q16[sLTP_buf_idx] = LPC_exc_Q10 << 6
			sLTP_buf_idx++
			sig_Q10_ptr[i] = LPC_exc_Q10
		}
		sig_Q10_ptr += ([]int32)(psDec.Subfr_length)
		for j = 0; j < LTP_ORDER; j++ {
			B_Q14[j] = int16(SKP_SMULBB(harm_Gain_Q15, int32(B_Q14[j])) >> 15)
		}
		rand_scale_Q14 = int16(SKP_SMULBB(int32(rand_scale_Q14), rand_Gain_Q15) >> 15)
		psPLC.PitchL_Q8 += SKP_SMULWB(psPLC.PitchL_Q8, PITCH_DRIFT_FAC_Q16)
		psPLC.PitchL_Q8 = SKP_min_32(psPLC.PitchL_Q8, SKP_SMULBB(MAX_PITCH_LAG_MS, psDec.Fs_kHz)<<8)
		lag = SKP_RSHIFT_ROUND(psPLC.PitchL_Q8, 8)
	}
	sig_Q10_ptr = ([]int32)(sig_Q10[:])
	memcpy(unsafe.Pointer(&A_Q12_tmp.As_int16[0]), unsafe.Pointer(&psPLC.PrevLPC_Q12[0]), size_t(uintptr(psDec.LPC_order)*unsafe.Sizeof(int16(0))))
	for k = 0; k < NB_SUBFR; k++ {
		for i = 0; i < psDec.Subfr_length; i++ {
			LPC_pred_Q10 = SKP_SMULWB(psDec.SLPC_Q14[MAX_LPC_ORDER+i-1], int32(A_Q12_tmp.As_int16[0]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-2], int32(A_Q12_tmp.As_int16[1]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-3], int32(A_Q12_tmp.As_int16[2]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-4], int32(A_Q12_tmp.As_int16[3]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-5], int32(A_Q12_tmp.As_int16[4]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-6], int32(A_Q12_tmp.As_int16[5]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-7], int32(A_Q12_tmp.As_int16[6]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-8], int32(A_Q12_tmp.As_int16[7]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-9], int32(A_Q12_tmp.As_int16[8]))
			LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-10], int32(A_Q12_tmp.As_int16[9]))
			for j = 10; j < psDec.LPC_order; j++ {
				LPC_pred_Q10 = SKP_SMLAWB(LPC_pred_Q10, psDec.SLPC_Q14[MAX_LPC_ORDER+i-j-1], int32(A_Q12_tmp.As_int16[j]))
			}
			sig_Q10_ptr[i] = (sig_Q10_ptr[i]) + LPC_pred_Q10
			psDec.SLPC_Q14[MAX_LPC_ORDER+i] = (sig_Q10_ptr[i]) << 4
		}
		sig_Q10_ptr += ([]int32)(psDec.Subfr_length)
		memcpy(unsafe.Pointer(&psDec.SLPC_Q14[0]), unsafe.Pointer(&psDec.SLPC_Q14[psDec.Subfr_length]), size_t(MAX_LPC_ORDER*unsafe.Sizeof(int32(0))))
	}
	for i = 0; i < psDec.Frame_length; i++ {
		signal[i] = SKP_SAT16(SKP_RSHIFT_ROUND(SKP_SMULWW(sig_Q10[i], psPLC.PrevGain_Q16[NB_SUBFR-1]), 10))
	}
	psPLC.Rand_seed = rand_seed
	psPLC.RandScale_Q14 = rand_scale_Q14
	for i = 0; i < NB_SUBFR; i++ {
		psDecCtrl.PitchL[i] = lag
	}
}
func SKP_Silk_PLC_glue_frames(psDec *SKP_Silk_decoder_state, psDecCtrl *SKP_Silk_decoder_control, signal []int16, length int32) {
	var (
		i            int32
		energy_shift int32
		energy       int32
		psPLC        *SKP_Silk_PLC_struct
	)
	psPLC = &psDec.SPLC
	if psDec.LossCnt != 0 {
		SKP_Silk_sum_sqr_shift(&psPLC.Conc_energy, &psPLC.Conc_energy_shift, signal, length)
		psPLC.Last_frame_lost = 1
	} else {
		if psDec.SPLC.Last_frame_lost != 0 {
			SKP_Silk_sum_sqr_shift(&energy, &energy_shift, signal, length)
			if energy_shift > psPLC.Conc_energy_shift {
				psPLC.Conc_energy = psPLC.Conc_energy >> (energy_shift - psPLC.Conc_energy_shift)
			} else if energy_shift < psPLC.Conc_energy_shift {
				energy = energy >> (psPLC.Conc_energy_shift - energy_shift)
			}
			if energy > psPLC.Conc_energy {
				var (
					frac_Q24  int32
					LZ        int32
					gain_Q12  int32
					slope_Q12 int32
				)
				LZ = SKP_Silk_CLZ32(psPLC.Conc_energy)
				LZ = LZ - 1
				psPLC.Conc_energy = psPLC.Conc_energy << LZ
				energy = energy >> SKP_max_32(24-LZ, 0)
				frac_Q24 = psPLC.Conc_energy / (func() int32 {
					if energy > 1 {
						return energy
					}
					return 1
				}())
				gain_Q12 = SKP_Silk_SQRT_APPROX(frac_Q24)
				slope_Q12 = ((1 << 12) - gain_Q12) / length
				for i = 0; i < length; i++ {
					signal[i] = int16((int64(gain_Q12) * int64(signal[i])) >> 12)
					gain_Q12 += slope_Q12
					if gain_Q12 < (1 << 12) {
						gain_Q12 = gain_Q12
					} else {
						gain_Q12 = 1 << 12
					}
				}
			}
		}
		psPLC.Last_frame_lost = 0
	}
}
