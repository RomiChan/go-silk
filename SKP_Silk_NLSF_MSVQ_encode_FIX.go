package silk

import "unsafe"

func SKP_Silk_NLSF_MSVQ_encode_FIX(NLSFIndices *int32, pNLSF_Q15 *int32, psNLSF_CB *SKP_Silk_NLSF_CB_struct, pNLSF_q_Q15_prev *int32, pW_Q6 *int32, NLSF_mu_Q15 int32, NLSF_mu_fluc_red_Q16 int32, NLSF_MSVQ_Survivors int32, LPC_order int32, deactivate_fluc_red int32) {
	var (
		i                     int32
		s                     int32
		k                     int32
		cur_survivors         int32 = 0
		prev_survivors        int32
		min_survivors         int32
		input_index           int32
		cb_index              int32
		bestIndex             int32
		rateDistThreshold_Q18 int32
		se_Q15                int32
		wsse_Q20              int32
		bestRateDist_Q20      int32
		pRateDist_Q18         [256]int32
		pRate_Q5              [16]int32
		pRate_new_Q5          [16]int32
		pTempIndices          [16]int32
		pPath                 [160]int32
		pPath_new             [160]int32
		pRes_Q15              [256]int32
		pRes_new_Q15          [256]int32
		pConstInt             *int32
		pInt                  *int32
		pCB_element           *int16
		pCurrentCBStage       *SKP_Silk_NLSF_CBS
	)
	memset(unsafe.Pointer(&pRate_Q5[0]), 0, size_t(uintptr(NLSF_MSVQ_Survivors)*unsafe.Sizeof(int32(0))))
	for i = 0; i < LPC_order; i++ {
		pRes_Q15[i] = *(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i)))
	}
	prev_survivors = 1
	min_survivors = NLSF_MSVQ_Survivors / 2
	for s = 0; s < psNLSF_CB.NStages; s++ {
		pCurrentCBStage = (*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*uintptr(s)))
		cur_survivors = SKP_min_32(NLSF_MSVQ_Survivors, SKP_SMULBB(prev_survivors, pCurrentCBStage.NVectors))
		SKP_Silk_NLSF_VQ_rate_distortion_FIX(&pRateDist_Q18[0], pCurrentCBStage, &pRes_Q15[0], pW_Q6, &pRate_Q5[0], NLSF_mu_Q15, prev_survivors, LPC_order)
		SKP_Silk_insertion_sort_increasing(&pRateDist_Q18[0], &pTempIndices[0], prev_survivors*pCurrentCBStage.NVectors, cur_survivors)
		if pRateDist_Q18[0] < SKP_int32_MAX/MAX_NLSF_MSVQ_SURVIVORS {
			rateDistThreshold_Q18 = SKP_SMLAWB(pRateDist_Q18[0], NLSF_MSVQ_Survivors*(pRateDist_Q18[0]), SKP_FIX_CONST(0.1, 16))
			for pRateDist_Q18[cur_survivors-1] > rateDistThreshold_Q18 && cur_survivors > min_survivors {
				cur_survivors--
			}
		}
		for k = 0; k < cur_survivors; k++ {
			if s > 0 {
				if pCurrentCBStage.NVectors == 8 {
					input_index = (pTempIndices[k]) >> 3
					cb_index = pTempIndices[k] & 7
				} else {
					input_index = (pTempIndices[k]) / pCurrentCBStage.NVectors
					cb_index = pTempIndices[k] - SKP_SMULBB(input_index, pCurrentCBStage.NVectors)
				}
			} else {
				input_index = 0
				cb_index = pTempIndices[k]
			}
			pConstInt = &pRes_Q15[SKP_SMULBB(input_index, LPC_order)]
			pCB_element = (*int16)(unsafe.Add(unsafe.Pointer(pCurrentCBStage.CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(SKP_SMULBB(cb_index, LPC_order))))
			pInt = &pRes_new_Q15[SKP_SMULBB(k, LPC_order)]
			for i = 0; i < LPC_order; i++ {
				*(*int32)(unsafe.Add(unsafe.Pointer(pInt), unsafe.Sizeof(int32(0))*uintptr(i))) = *(*int32)(unsafe.Add(unsafe.Pointer(pConstInt), unsafe.Sizeof(int32(0))*uintptr(i))) - int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*uintptr(i))))
			}
			pRate_new_Q5[k] = int32(int64(pRate_Q5[input_index]) + int64(*(*int16)(unsafe.Add(unsafe.Pointer(pCurrentCBStage.Rates_Q5), unsafe.Sizeof(int16(0))*uintptr(cb_index)))))
			pConstInt = &pPath[SKP_SMULBB(input_index, psNLSF_CB.NStages)]
			pInt = &pPath_new[SKP_SMULBB(k, psNLSF_CB.NStages)]
			for i = 0; i < s; i++ {
				*(*int32)(unsafe.Add(unsafe.Pointer(pInt), unsafe.Sizeof(int32(0))*uintptr(i))) = *(*int32)(unsafe.Add(unsafe.Pointer(pConstInt), unsafe.Sizeof(int32(0))*uintptr(i)))
			}
			*(*int32)(unsafe.Add(unsafe.Pointer(pInt), unsafe.Sizeof(int32(0))*uintptr(s))) = cb_index
		}
		if s < psNLSF_CB.NStages-1 {
			memcpy(unsafe.Pointer(&pRes_Q15[0]), unsafe.Pointer(&pRes_new_Q15[0]), size_t(uintptr(SKP_SMULBB(cur_survivors, LPC_order))*unsafe.Sizeof(int32(0))))
			memcpy(unsafe.Pointer(&pRate_Q5[0]), unsafe.Pointer(&pRate_new_Q5[0]), size_t(uintptr(cur_survivors)*unsafe.Sizeof(int32(0))))
			memcpy(unsafe.Pointer(&pPath[0]), unsafe.Pointer(&pPath_new[0]), size_t(uintptr(SKP_SMULBB(cur_survivors, psNLSF_CB.NStages))*unsafe.Sizeof(int32(0))))
		}
		prev_survivors = cur_survivors
	}
	bestIndex = 0
	if deactivate_fluc_red != 1 {
		bestRateDist_Q20 = SKP_int32_MAX
		for s = 0; s < cur_survivors; s++ {
			SKP_Silk_NLSF_MSVQ_decode(pNLSF_Q15, psNLSF_CB, &pPath_new[SKP_SMULBB(s, psNLSF_CB.NStages)], LPC_order)
			wsse_Q20 = 0
			for i = 0; i < LPC_order; i += 2 {
				se_Q15 = *(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))) - *(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_q_Q15_prev), unsafe.Sizeof(int32(0))*uintptr(i)))
				wsse_Q20 = SKP_SMLAWB(wsse_Q20, SKP_SMULBB(se_Q15, se_Q15), *(*int32)(unsafe.Add(unsafe.Pointer(pW_Q6), unsafe.Sizeof(int32(0))*uintptr(i))))
				se_Q15 = *(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i+1))) - *(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_q_Q15_prev), unsafe.Sizeof(int32(0))*uintptr(i+1)))
				wsse_Q20 = SKP_SMLAWB(wsse_Q20, SKP_SMULBB(se_Q15, se_Q15), *(*int32)(unsafe.Add(unsafe.Pointer(pW_Q6), unsafe.Sizeof(int32(0))*uintptr(i+1))))
			}
			wsse_Q20 = SKP_ADD_POS_SAT32(pRateDist_Q18[s], SKP_SMULWB(wsse_Q20, NLSF_mu_fluc_red_Q16))
			if wsse_Q20 < bestRateDist_Q20 {
				bestRateDist_Q20 = wsse_Q20
				bestIndex = s
			}
		}
	}
	memcpy(unsafe.Pointer(NLSFIndices), unsafe.Pointer(&pPath_new[SKP_SMULBB(bestIndex, psNLSF_CB.NStages)]), size_t(uintptr(psNLSF_CB.NStages)*unsafe.Sizeof(int32(0))))
	SKP_Silk_NLSF_MSVQ_decode(pNLSF_Q15, psNLSF_CB, NLSFIndices, LPC_order)
}
