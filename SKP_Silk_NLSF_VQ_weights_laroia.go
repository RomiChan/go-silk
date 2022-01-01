package silk

const Q_OUT = 6
const MIN_NDELTA = 3

func SKP_Silk_NLSF_VQ_weights_laroia(pNLSFW_Q6 []int32, pNLSF_Q15 []int32, D int32) {
	tmp1_int := SKP_max_32(pNLSF_Q15[0], MIN_NDELTA)
	tmp1_int = (1 << (Q_OUT + 15)) / tmp1_int
	tmp2_int := SKP_max_32(pNLSF_Q15[1]-pNLSF_Q15[0], MIN_NDELTA)
	tmp2_int = (1 << (Q_OUT + 15)) / tmp2_int
	pNLSFW_Q6[0] = SKP_min_32(tmp1_int+tmp2_int, SKP_int16_MAX)
	for k := int32(1); k < D-1; k += 2 {
		tmp1_int = SKP_max_32(pNLSF_Q15[k+1]-pNLSF_Q15[k], MIN_NDELTA)
		tmp1_int = (1 << (Q_OUT + 15)) / tmp1_int
		pNLSFW_Q6[k] = SKP_min_32(tmp1_int+tmp2_int, SKP_int16_MAX)
		tmp2_int = SKP_max_32(pNLSF_Q15[k+2]-pNLSF_Q15[k+1], MIN_NDELTA)
		tmp2_int = (1 << (Q_OUT + 15)) / tmp2_int
		pNLSFW_Q6[k+1] = SKP_min_32(tmp1_int+tmp2_int, SKP_int16_MAX)
	}
	tmp1_int = SKP_max_32((1<<15)-pNLSF_Q15[D-1], MIN_NDELTA)
	tmp1_int = (1 << (Q_OUT + 15)) / tmp1_int
	pNLSFW_Q6[D-1] = SKP_min_32(tmp1_int+tmp2_int, SKP_int16_MAX)
}
