package silk

import "unsafe"

const QC = 10
const QS = 14

func SKP_Silk_warped_autocorrelation_FIX(corr *int32, scale *int32, input *int16, warping_Q16 int16, length int32, order int32) {
	var (
		n        int32
		i        int32
		lsh      int32
		tmp1_QS  int32
		tmp2_QS  int32
		state_QS [17]int32 = [17]int32{}
		corr_QC  [17]int64 = [17]int64{}
	)
	for n = 0; n < length; n++ {
		tmp1_QS = (int32(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(n))))) << QS
		for i = 0; i < order; i += 2 {
			tmp2_QS = SKP_SMLAWB(state_QS[i], state_QS[i+1]-tmp1_QS, int32(warping_Q16))
			state_QS[i] = tmp1_QS
			corr_QC[i] += (int64(tmp1_QS) * int64(state_QS[0])) >> (QS*2 - QC)
			tmp1_QS = SKP_SMLAWB(state_QS[i+1], state_QS[i+2]-tmp2_QS, int32(warping_Q16))
			state_QS[i+1] = tmp2_QS
			corr_QC[i+1] += (int64(tmp2_QS) * int64(state_QS[0])) >> (QS*2 - QC)
		}
		state_QS[order] = tmp1_QS
		corr_QC[order] += (int64(tmp1_QS) * int64(state_QS[0])) >> (QS*2 - QC)
	}
	lsh = SKP_Silk_CLZ64(corr_QC[0]) - 35
	if (-12 - QC) > (30 - QC) {
		if int64(lsh) > int64(-12-QC) {
			lsh = int32(-12 - QC)
		} else if lsh < (30 - QC) {
			lsh = 30 - QC
		} else {
			lsh = lsh
		}
	} else if lsh > (30 - QC) {
		lsh = 30 - QC
	} else if int64(lsh) < int64(-12-QC) {
		lsh = int32(-12 - QC)
	} else {
		lsh = lsh
	}
	*scale = -(QC + lsh)
	if lsh >= 0 {
		for i = 0; i < order+1; i++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(corr), unsafe.Sizeof(int32(0))*uintptr(i))) = int32((corr_QC[i]) << int64(lsh))
		}
	} else {
		for i = 0; i < order+1; i++ {
			*(*int32)(unsafe.Add(unsafe.Pointer(corr), unsafe.Sizeof(int32(0))*uintptr(i))) = int32((corr_QC[i]) >> int64(-lsh))
		}
	}
}
