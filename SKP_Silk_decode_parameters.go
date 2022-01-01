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
		psNLSF_CB    *SKP_Silk_NLSF_CB_struct    = (*SKP_Silk_NLSF_CB_struct)(nil)
		psRC         *SKP_Silk_range_coder_state = &psDec.SRC
	)
	if int64(psDec.NFramesDecoded) == 0 {
		SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_SamplingRates_CDF[:], SKP_Silk_SamplingRates_offset)
		if int64(Ix) < 0 || int64(Ix) > 3 {
			psRC.Error = -7
			return
		}
		fs_kHz_dec = SKP_Silk_SamplingRates_table[Ix]
		SKP_Silk_decoder_set_fs(psDec, fs_kHz_dec)
	}
	if int64(psDec.NFramesDecoded) == 0 {
		SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_type_offset_CDF[:], SKP_Silk_type_offset_CDF_offset)
	} else {
		SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_type_offset_joint_CDF[psDec.TypeOffsetPrev][:], SKP_Silk_type_offset_CDF_offset)
	}
	psDecCtrl.Sigtype = Ix >> 1
	psDecCtrl.QuantOffsetType = int32(int64(Ix) & 1)
	psDec.TypeOffsetPrev = Ix
	if int64(psDec.NFramesDecoded) == 0 {
		SKP_Silk_range_decoder(&GainsIndices[0], psRC, SKP_Silk_gain_CDF[psDecCtrl.Sigtype][:], SKP_Silk_gain_CDF_offset)
	} else {
		SKP_Silk_range_decoder(&GainsIndices[0], psRC, SKP_Silk_delta_gain_CDF[:], SKP_Silk_delta_gain_CDF_offset)
	}
	for i = 1; int64(i) < NB_SUBFR; i++ {
		SKP_Silk_range_decoder(&GainsIndices[i], psRC, SKP_Silk_delta_gain_CDF[:], SKP_Silk_delta_gain_CDF_offset)
	}
	SKP_Silk_gains_dequant(psDecCtrl.Gains_Q16, GainsIndices, &psDec.LastGainIndex, psDec.NFramesDecoded)
	psNLSF_CB = psDec.PsNLSF_CB[psDecCtrl.Sigtype]
	SKP_Silk_range_decoder_multi(NLSFIndices[:], psRC, ([]*uint16)(psNLSF_CB.StartPtr), ([]int32)(psNLSF_CB.MiddleIx), psNLSF_CB.NStages)
	SKP_Silk_NLSF_MSVQ_decode(&pNLSF_Q15[0], psNLSF_CB, &NLSFIndices[0], psDec.LPC_order)
	SKP_Silk_range_decoder(&psDecCtrl.NLSFInterpCoef_Q2, psRC, SKP_Silk_NLSF_interpolation_factor_CDF[:], SKP_Silk_NLSF_interpolation_factor_offset)
	if int64(psDec.First_frame_after_reset) == 1 {
		psDecCtrl.NLSFInterpCoef_Q2 = 4
	}
	if int64(fullDecoding) != 0 {
		SKP_Silk_NLSF2A_stable(psDecCtrl.PredCoef_Q12[1], pNLSF_Q15, psDec.LPC_order)
		if int64(psDecCtrl.NLSFInterpCoef_Q2) < 4 {
			for i = 0; int64(i) < int64(psDec.LPC_order); i++ {
				pNLSF0_Q15[i] = int32(int64(psDec.PrevNLSF_Q15[i]) + ((int64(psDecCtrl.NLSFInterpCoef_Q2) * (int64(pNLSF_Q15[i]) - int64(psDec.PrevNLSF_Q15[i]))) >> 2))
			}
			SKP_Silk_NLSF2A_stable(psDecCtrl.PredCoef_Q12[0], pNLSF0_Q15, psDec.LPC_order)
		} else {
			memcpy(unsafe.Pointer(&(psDecCtrl.PredCoef_Q12[0])[0]), unsafe.Pointer(&(psDecCtrl.PredCoef_Q12[1])[0]), size_t(uintptr(psDec.LPC_order)*unsafe.Sizeof(int16(0))))
		}
	}
	memcpy(unsafe.Pointer(&psDec.PrevNLSF_Q15[0]), unsafe.Pointer(&pNLSF_Q15[0]), size_t(uintptr(psDec.LPC_order)*unsafe.Sizeof(int32(0))))
	if int64(psDec.LossCnt) != 0 {
		SKP_Silk_bwexpander(&psDecCtrl.PredCoef_Q12[0][0], psDec.LPC_order, BWE_AFTER_LOSS_Q16)
		SKP_Silk_bwexpander(&psDecCtrl.PredCoef_Q12[1][0], psDec.LPC_order, BWE_AFTER_LOSS_Q16)
	}
	if int64(psDecCtrl.Sigtype) == SIG_TYPE_VOICED {
		if int64(psDec.Fs_kHz) == 8 {
			SKP_Silk_range_decoder(&Ixs[0], psRC, SKP_Silk_pitch_lag_NB_CDF[:], SKP_Silk_pitch_lag_NB_CDF_offset)
		} else if int64(psDec.Fs_kHz) == 12 {
			SKP_Silk_range_decoder(&Ixs[0], psRC, SKP_Silk_pitch_lag_MB_CDF[:], SKP_Silk_pitch_lag_MB_CDF_offset)
		} else if int64(psDec.Fs_kHz) == 16 {
			SKP_Silk_range_decoder(&Ixs[0], psRC, SKP_Silk_pitch_lag_WB_CDF[:], SKP_Silk_pitch_lag_WB_CDF_offset)
		} else {
			SKP_Silk_range_decoder(&Ixs[0], psRC, SKP_Silk_pitch_lag_SWB_CDF[:], SKP_Silk_pitch_lag_SWB_CDF_offset)
		}
		if int64(psDec.Fs_kHz) == 8 {
			SKP_Silk_range_decoder(&Ixs[1], psRC, SKP_Silk_pitch_contour_NB_CDF[:], SKP_Silk_pitch_contour_NB_CDF_offset)
		} else {
			SKP_Silk_range_decoder(&Ixs[1], psRC, SKP_Silk_pitch_contour_CDF[:], SKP_Silk_pitch_contour_CDF_offset)
		}
		SKP_Silk_decode_pitch(Ixs[0], Ixs[1], psDecCtrl.PitchL[:], psDec.Fs_kHz)
		SKP_Silk_range_decoder(&psDecCtrl.PERIndex, psRC, SKP_Silk_LTP_per_index_CDF[:], SKP_Silk_LTP_per_index_CDF_offset)
		cbk_ptr_Q14 = SKP_Silk_LTP_vq_ptrs_Q14[psDecCtrl.PERIndex]
		for k = 0; int64(k) < NB_SUBFR; k++ {
			SKP_Silk_range_decoder(&Ix, psRC, ([]uint16)(SKP_Silk_LTP_gain_CDF_ptrs[psDecCtrl.PERIndex]), SKP_Silk_LTP_gain_CDF_offsets[psDecCtrl.PERIndex])
			for i = 0; int64(i) < LTP_ORDER; i++ {
				psDecCtrl.LTPCoef_Q14[int64(k)*LTP_ORDER+int64(i)] = *(*int16)(unsafe.Add(unsafe.Pointer(cbk_ptr_Q14), unsafe.Sizeof(int16(0))*uintptr(int64(Ix)*LTP_ORDER+int64(i))))
			}
		}
		SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_LTPscale_CDF[:], SKP_Silk_LTPscale_offset)
		psDecCtrl.LTP_scale_Q14 = int32(SKP_Silk_LTPScales_table_Q14[Ix])
	} else {
		memset(unsafe.Pointer(&psDecCtrl.PitchL[0]), 0, size_t(NB_SUBFR*unsafe.Sizeof(int32(0))))
		memset(unsafe.Pointer(&psDecCtrl.LTPCoef_Q14[0]), 0, size_t(uintptr(LTP_ORDER*NB_SUBFR)*unsafe.Sizeof(int16(0))))
		psDecCtrl.PERIndex = 0
		psDecCtrl.LTP_scale_Q14 = 0
	}
	SKP_Silk_range_decoder(&Ix, psRC, SKP_Silk_Seed_CDF[:], SKP_Silk_Seed_offset)
	psDecCtrl.Seed = Ix
	SKP_Silk_decode_pulses(psRC, psDecCtrl, q, psDec.Frame_length)
	SKP_Silk_range_decoder(&psDec.VadFlag, psRC, SKP_Silk_vadflag_CDF[:], SKP_Silk_vadflag_offset)
	SKP_Silk_range_decoder(&psDec.FrameTermination, psRC, SKP_Silk_FrameTermination_CDF[:], SKP_Silk_FrameTermination_offset)
	SKP_Silk_range_coder_get_length(psRC, &nBytesUsed)
	psDec.NBytesLeft = int32(int64(psRC.BufferLength) - int64(nBytesUsed))
	if int64(psDec.NBytesLeft) < 0 {
		psRC.Error = -6
	}
	if int64(psDec.NBytesLeft) == 0 {
		SKP_Silk_range_coder_check_after_decoding(psRC)
	}
}
