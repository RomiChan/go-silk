package silk

import "unsafe"

func SKP_Silk_decode_parameters(psDec *SKP_Silk_decoder_state, psDecCtrl *SKP_Silk_decoder_control, q []int32, fullDecoding int32) {
	var (
		i            int32
		k            int32
		Ix           int32
		fs_kHz_dec   int32
		nBytesUsed   int32
		Ixs          [4]int32
		GainsIndices [4]int32
		NLSFIndices  [10]int32
		pNLSF_Q15    [16]int32
		pNLSF0_Q15   [16]int32
		cbk_ptr_Q14  *int16
		psNLSF_CB    *SKP_Silk_NLSF_CB_struct
		psRC         = &psDec.SRC
	)
	if psDec.NFramesDecoded == 0 {
		SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_SamplingRates_CDF[:], SKP_Silk_SamplingRates_offset)
		if Ix < 0 || Ix > 3 {
			psRC.Error = -7
			return
		}
		fs_kHz_dec = SKP_Silk_SamplingRates_table[Ix]
		SKP_Silk_decoder_set_fs(psDec, fs_kHz_dec)
	}
	if psDec.NFramesDecoded == 0 {
		SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_type_offset_CDF[:], SKP_Silk_type_offset_CDF_offset)
	} else {
		SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_type_offset_joint_CDF[psDec.TypeOffsetPrev][:], SKP_Silk_type_offset_CDF_offset)
	}
	psDecCtrl.Sigtype = Ix >> 1
	psDecCtrl.QuantOffsetType = Ix & 1
	psDec.TypeOffsetPrev = Ix
	if psDec.NFramesDecoded == 0 {
		SKP_Silk_range_decoder(&GainsIndices[0], psRC, SKP_Silk_gain_CDF[psDecCtrl.Sigtype][:], SKP_Silk_gain_CDF_offset)
	} else {
		SKP_Silk_range_decoder(&GainsIndices[0], psRC, SKP_Silk_delta_gain_CDF[:], SKP_Silk_delta_gain_CDF_offset)
	}
	for i = 1; i < NB_SUBFR; i++ {
		SKP_Silk_range_decoder(&GainsIndices[i], psRC, SKP_Silk_delta_gain_CDF[:], SKP_Silk_delta_gain_CDF_offset)
	}
	SKP_Silk_gains_dequant(psDecCtrl.Gains_Q16, GainsIndices, &psDec.LastGainIndex, psDec.NFramesDecoded)
	psNLSF_CB = psDec.PsNLSF_CB[psDecCtrl.Sigtype]
	SKP_Silk_range_decoder_multi(NLSFIndices[:], psRC, ([]*uint16)(psNLSF_CB.StartPtr), ([]int32)(psNLSF_CB.MiddleIx), psNLSF_CB.NStages)
	SKP_Silk_NLSF_MSVQ_decode(pNLSF_Q15[:], psNLSF_CB, NLSFIndices[:], psDec.LPC_order)
	SKP_Silk_range_decoder(&psDecCtrl.NLSFInterpCoef_Q2, psRC, SKP_Silk_NLSF_interpolation_factor_CDF[:], SKP_Silk_NLSF_interpolation_factor_offset)
	if psDec.First_frame_after_reset == 1 {
		psDecCtrl.NLSFInterpCoef_Q2 = 4
	}
	if fullDecoding != 0 {
		SKP_Silk_NLSF2A_stable(psDecCtrl.PredCoef_Q12[1], pNLSF_Q15, psDec.LPC_order)
		if psDecCtrl.NLSFInterpCoef_Q2 < 4 {
			for i = 0; i < psDec.LPC_order; i++ {
				pNLSF0_Q15[i] = psDec.PrevNLSF_Q15[i] + ((psDecCtrl.NLSFInterpCoef_Q2 * (pNLSF_Q15[i] - psDec.PrevNLSF_Q15[i])) >> 2)
			}
			SKP_Silk_NLSF2A_stable(psDecCtrl.PredCoef_Q12[0], pNLSF0_Q15, psDec.LPC_order)
		} else {
			memcpy(unsafe.Pointer(&(psDecCtrl.PredCoef_Q12[0])[0]), unsafe.Pointer(&(psDecCtrl.PredCoef_Q12[1])[0]), uintptr(psDec.LPC_order)*unsafe.Sizeof(int16(0)))
		}
	}
	memcpy(unsafe.Pointer(&psDec.PrevNLSF_Q15[0]), unsafe.Pointer(&pNLSF_Q15[0]), uintptr(psDec.LPC_order)*unsafe.Sizeof(int32(0)))
	if psDec.LossCnt != 0 {
		SKP_Silk_bwexpander(psDecCtrl.PredCoef_Q12[0][:], psDec.LPC_order, BWE_AFTER_LOSS_Q16)
		SKP_Silk_bwexpander(psDecCtrl.PredCoef_Q12[1][:], psDec.LPC_order, BWE_AFTER_LOSS_Q16)
	}
	if psDecCtrl.Sigtype == SIG_TYPE_VOICED {
		if psDec.Fs_kHz == 8 {
			SKP_Silk_range_decoder(&Ixs[0], psRC, SKP_Silk_pitch_lag_NB_CDF[:], SKP_Silk_pitch_lag_NB_CDF_offset)
		} else if psDec.Fs_kHz == 12 {
			SKP_Silk_range_decoder(&Ixs[0], psRC, SKP_Silk_pitch_lag_MB_CDF[:], SKP_Silk_pitch_lag_MB_CDF_offset)
		} else if psDec.Fs_kHz == 16 {
			SKP_Silk_range_decoder(&Ixs[0], psRC, SKP_Silk_pitch_lag_WB_CDF[:], SKP_Silk_pitch_lag_WB_CDF_offset)
		} else {
			SKP_Silk_range_decoder(&Ixs[0], psRC, SKP_Silk_pitch_lag_SWB_CDF[:], SKP_Silk_pitch_lag_SWB_CDF_offset)
		}
		if psDec.Fs_kHz == 8 {
			SKP_Silk_range_decoder(&Ixs[1], psRC, SKP_Silk_pitch_contour_NB_CDF[:], SKP_Silk_pitch_contour_NB_CDF_offset)
		} else {
			SKP_Silk_range_decoder(&Ixs[1], psRC, SKP_Silk_pitch_contour_CDF[:], SKP_Silk_pitch_contour_CDF_offset)
		}
		SKP_Silk_decode_pitch(Ixs[0], Ixs[1], psDecCtrl.PitchL[:], psDec.Fs_kHz)
		SKP_Silk_range_decoder(&psDecCtrl.PERIndex, psRC, SKP_Silk_LTP_per_index_CDF[:], SKP_Silk_LTP_per_index_CDF_offset)
		cbk_ptr_Q14 = SKP_Silk_LTP_vq_ptrs_Q14[psDecCtrl.PERIndex]
		for k = 0; k < NB_SUBFR; k++ {
			SKP_Silk_range_decoder(&Ix, psRC, ([]uint16)(SKP_Silk_LTP_gain_CDF_ptrs[psDecCtrl.PERIndex]), SKP_Silk_LTP_gain_CDF_offsets[psDecCtrl.PERIndex])
			for i = 0; i < LTP_ORDER; i++ {
				psDecCtrl.LTPCoef_Q14[k*LTP_ORDER+i] = *(*int16)(unsafe.Add(unsafe.Pointer(cbk_ptr_Q14), unsafe.Sizeof(int16(0))*uintptr(Ix*LTP_ORDER+i)))
			}
		}
		SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_LTPscale_CDF[:], SKP_Silk_LTPscale_offset)
		psDecCtrl.LTP_scale_Q14 = int32(SKP_Silk_LTPScales_table_Q14[Ix])
	} else {
		memset(unsafe.Pointer(&psDecCtrl.PitchL[0]), 0, NB_SUBFR*unsafe.Sizeof(int32(0)))
		memset(unsafe.Pointer(&psDecCtrl.LTPCoef_Q14[0]), 0, LTP_ORDER*NB_SUBFR*unsafe.Sizeof(int16(0)))
		psDecCtrl.PERIndex = 0
		psDecCtrl.LTP_scale_Q14 = 0
	}
	SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_Seed_CDF[:], SKP_Silk_Seed_offset)
	psDecCtrl.Seed = Ix
	SKP_Silk_decode_pulses(psRC, psDecCtrl, q, psDec.Frame_length)
	SKP_Silk_range_decoder(&psDec.VadFlag, psRC, SKP_Silk_vadflag_CDF[:], SKP_Silk_vadflag_offset)
	SKP_Silk_range_decoder(&psDec.FrameTermination, psRC, SKP_Silk_FrameTermination_CDF[:], SKP_Silk_FrameTermination_offset)
	SKP_Silk_range_coder_get_length(psRC, &nBytesUsed)
	psDec.NBytesLeft = psRC.BufferLength - nBytesUsed
	if psDec.NBytesLeft < 0 {
		psRC.Error = -6
	}
	if psDec.NBytesLeft == 0 {
		SKP_Silk_range_coder_check_after_decoding(psRC)
	}
}
