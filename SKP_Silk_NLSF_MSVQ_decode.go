package silk

import "unsafe"

func SKP_Silk_NLSF_MSVQ_decode(pNLSF_Q15 []int32, psNLSF_CB *SKP_Silk_NLSF_CB_struct, NLSFIndices []int32, LPC_order int32) {
	var (
		pCB_element []int16
		s           int32
		i           int32
	)
	pCB_element = ([]int16)((*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*0))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr((NLSFIndices[0])*LPC_order))))
	for i = 0; i < LPC_order; i++ {
		pNLSF_Q15[i] = int32(pCB_element[i])
	}
	for s = 1; s < psNLSF_CB.NStages; s++ {
		if LPC_order == 16 {
			pCB_element = ([]int16)((*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*uintptr(s)))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr((NLSFIndices[s])<<4))))

			pNLSF_Q15[0] += int32(pCB_element[0])
			pNLSF_Q15[1] += int32(pCB_element[1])
			pNLSF_Q15[2] += int32(pCB_element[2])
			pNLSF_Q15[3] += int32(pCB_element[3])
			pNLSF_Q15[4] += int32(pCB_element[4])
			pNLSF_Q15[5] += int32(pCB_element[5])
			pNLSF_Q15[6] += int32(pCB_element[6])
			pNLSF_Q15[7] += int32(pCB_element[7])
			pNLSF_Q15[8] += int32(pCB_element[8])
			pNLSF_Q15[9] += int32(pCB_element[9])
			pNLSF_Q15[10] += int32(pCB_element[10])
			pNLSF_Q15[11] += int32(pCB_element[11])
			pNLSF_Q15[12] += int32(pCB_element[12])
			pNLSF_Q15[13] += int32(pCB_element[13])
			pNLSF_Q15[14] += int32(pCB_element[14])
			pNLSF_Q15[15] += int32(pCB_element[15])
		} else {
			pCB_element = ([]int16)((*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*uintptr(s)))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(SKP_SMULBB(NLSFIndices[s], LPC_order)))))
			for i = 0; i < LPC_order; i++ {
				pNLSF_Q15[i] += int32(pCB_element[i])
			}
		}
	}
	SKP_Silk_NLSF_stabilize(pNLSF_Q15, ([]int32)(psNLSF_CB.NDeltaMin_Q15), LPC_order)
}
