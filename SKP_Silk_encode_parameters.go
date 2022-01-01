package silk

func SKP_Silk_encode_parameters(psEncC *SKP_Silk_encoder_state, psEncCtrlC *SKP_Silk_encoder_control, psRC *SKP_Silk_range_coder_state, q *int8) {
	var (
		i          int32
		k          int32
		typeOffset int32
		psNLSF_CB  *SKP_Silk_NLSF_CB_struct
	)
	if int64(psEncC.NFramesInPayloadBuf) == 0 {
		for i = 0; int64(i) < 3; i++ {
			if int64(SKP_Silk_SamplingRates_table[i]) == int64(psEncC.Fs_kHz) {
				break
			}
		}
		SKP_Silk_range_encoder(psRC, i, SKP_Silk_SamplingRates_CDF[:])
	}
	typeOffset = int32(int64(psEncCtrlC.Sigtype)*2 + int64(psEncCtrlC.QuantOffsetType))
	if int64(psEncC.NFramesInPayloadBuf) == 0 {
		SKP_Silk_range_encoder(psRC, typeOffset, SKP_Silk_type_offset_CDF[:])
	} else {
		SKP_Silk_range_encoder(psRC, typeOffset, SKP_Silk_type_offset_joint_CDF[psEncC.TypeOffsetPrev][:])
	}
	psEncC.TypeOffsetPrev = typeOffset
	if int64(psEncC.NFramesInPayloadBuf) == 0 {
		SKP_Silk_range_encoder(psRC, psEncCtrlC.GainsIndices[0], SKP_Silk_gain_CDF[psEncCtrlC.Sigtype][:])
	} else {
		SKP_Silk_range_encoder(psRC, psEncCtrlC.GainsIndices[0], SKP_Silk_delta_gain_CDF[:])
	}
	for i = 1; int64(i) < NB_SUBFR; i++ {
		SKP_Silk_range_encoder(psRC, psEncCtrlC.GainsIndices[i], SKP_Silk_delta_gain_CDF[:])
	}
	psNLSF_CB = psEncC.PsNLSF_CB[psEncCtrlC.Sigtype]
	SKP_Silk_range_encoder_multi(psRC, psEncCtrlC.NLSFIndices[:], ([]*uint16)(psNLSF_CB.StartPtr), psNLSF_CB.NStages)
	SKP_Silk_range_encoder(psRC, psEncCtrlC.NLSFInterpCoef_Q2, SKP_Silk_NLSF_interpolation_factor_CDF[:])
	if int64(psEncCtrlC.Sigtype) == SIG_TYPE_VOICED {
		if int64(psEncC.Fs_kHz) == 8 {
			SKP_Silk_range_encoder(psRC, psEncCtrlC.LagIndex, SKP_Silk_pitch_lag_NB_CDF[:])
		} else if int64(psEncC.Fs_kHz) == 12 {
			SKP_Silk_range_encoder(psRC, psEncCtrlC.LagIndex, SKP_Silk_pitch_lag_MB_CDF[:])
		} else if int64(psEncC.Fs_kHz) == 16 {
			SKP_Silk_range_encoder(psRC, psEncCtrlC.LagIndex, SKP_Silk_pitch_lag_WB_CDF[:])
		} else {
			SKP_Silk_range_encoder(psRC, psEncCtrlC.LagIndex, SKP_Silk_pitch_lag_SWB_CDF[:])
		}
		if int64(psEncC.Fs_kHz) == 8 {
			SKP_Silk_range_encoder(psRC, psEncCtrlC.ContourIndex, SKP_Silk_pitch_contour_NB_CDF[:])
		} else {
			SKP_Silk_range_encoder(psRC, psEncCtrlC.ContourIndex, SKP_Silk_pitch_contour_CDF[:])
		}
		SKP_Silk_range_encoder(psRC, psEncCtrlC.PERIndex, SKP_Silk_LTP_per_index_CDF[:])
		for k = 0; int64(k) < NB_SUBFR; k++ {
			SKP_Silk_range_encoder(psRC, psEncCtrlC.LTPIndex[k], ([]uint16)(SKP_Silk_LTP_gain_CDF_ptrs[psEncCtrlC.PERIndex]))
		}
		SKP_Silk_range_encoder(psRC, psEncCtrlC.LTP_scaleIndex, SKP_Silk_LTPscale_CDF[:])
	}
	SKP_Silk_range_encoder(psRC, psEncCtrlC.Seed, SKP_Silk_Seed_CDF[:])
	SKP_Silk_encode_pulses(psRC, psEncCtrlC.Sigtype, psEncCtrlC.QuantOffsetType, ([]int8)(q), psEncC.Frame_length)
	SKP_Silk_range_encoder(psRC, psEncC.VadFlag, SKP_Silk_vadflag_CDF[:])
}
