package silk

import "unsafe"

func SKP_Silk_NLSF_MSVQ_decode(pNLSF_Q15 []int32, psNLSF_CB *SKP_Silk_NLSF_CB_struct, NLSFIndices []int32, LPC_order int32) {
	var (
		pCB_element *int16
		s           int32
		i           int32
	)
	pCB_element = (*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*0))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr((NLSFIndices[0])*LPC_order)))
	for i = 0; i < LPC_order; i++ {
		pNLSF_Q15[i] = int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*uintptr(i))))
	}
	for s = 1; s < psNLSF_CB.NStages; s++ {
		if LPC_order == 16 {
			pCB_element = (*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*uintptr(s)))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr((NLSFIndices[s])<<4)))
			pNLSF_Q15[0] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*0)))
			pNLSF_Q15[1] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*1)))
			pNLSF_Q15[2] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*2)))
			pNLSF_Q15[3] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*3)))
			pNLSF_Q15[4] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*4)))
			pNLSF_Q15[5] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*5)))
			pNLSF_Q15[6] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*6)))
			pNLSF_Q15[7] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*7)))
			pNLSF_Q15[8] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*8)))
			pNLSF_Q15[9] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*9)))
			pNLSF_Q15[10] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*10)))
			pNLSF_Q15[11] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*11)))
			pNLSF_Q15[12] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*12)))
			pNLSF_Q15[13] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*13)))
			pNLSF_Q15[14] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*14)))
			pNLSF_Q15[15] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*15)))
		} else {
			pCB_element = (*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*uintptr(s)))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(SKP_SMULBB(NLSFIndices[s], LPC_order))))
			for i = 0; i < LPC_order; i++ {
				pNLSF_Q15[i] += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*uintptr(i))))
			}
		}
	}
	SKP_Silk_NLSF_stabilize(&pNLSF_Q15[0], psNLSF_CB.NDeltaMin_Q15, LPC_order)
}
