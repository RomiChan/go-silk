package silk

import "unsafe"

func SKP_Silk_NLSF_MSVQ_decode(pNLSF_Q15 *int32, psNLSF_CB *SKP_Silk_NLSF_CB_struct, NLSFIndices *int32, LPC_order int32) {
	var (
		pCB_element *int16
		s           int32
		i           int32
	)
	pCB_element = (*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*0))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSFIndices), unsafe.Sizeof(int32(0))*0)))*int64(LPC_order))))
	for i = 0; int64(i) < int64(LPC_order); i++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))) = int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*uintptr(i))))
	}
	for s = 1; int64(s) < int64(psNLSF_CB.NStages); s++ {
		if int64(LPC_order) == 16 {
			pCB_element = (*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*uintptr(s)))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(int64(*(*int32)(unsafe.Add(unsafe.Pointer(NLSFIndices), unsafe.Sizeof(int32(0))*uintptr(s))))<<4)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*0)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*0)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*1)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*1)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*2)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*2)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*3)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*3)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*4)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*4)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*5)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*5)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*6)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*6)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*7)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*7)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*8)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*8)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*9)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*9)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*10)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*10)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*11)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*11)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*12)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*12)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*13)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*13)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*14)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*14)))
			*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*15)) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*15)))
		} else {
			pCB_element = (*int16)(unsafe.Add(unsafe.Pointer((*(*SKP_Silk_NLSF_CBS)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CBStages), unsafe.Sizeof(SKP_Silk_NLSF_CBS{})*uintptr(s)))).CB_NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(SKP_SMULBB(*(*int32)(unsafe.Add(unsafe.Pointer(NLSFIndices), unsafe.Sizeof(int32(0))*uintptr(s))), LPC_order))))
			for i = 0; int64(i) < int64(LPC_order); i++ {
				*(*int32)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int32(0))*uintptr(i))) += int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_element), unsafe.Sizeof(int16(0))*uintptr(i))))
			}
		}
	}
	SKP_Silk_NLSF_stabilize(pNLSF_Q15, psNLSF_CB.NDeltaMin_Q15, LPC_order)
}
