package silk

func SKP_Silk_regularize_correlations_FIX(XX []int32, xx []int32, noise int32, D int32) {
	var i int32
	for i = 0; i < D; i++ {
		XX[(i*D+i)+0] = (XX[(i*D+i)+0]) + noise
	}
	xx[0] += noise
}
