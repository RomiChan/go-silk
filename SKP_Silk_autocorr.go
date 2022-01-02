package silk

func SKP_Silk_autocorr(results []int32, scale *int, inputData []int16, inputDataSize int, correlationCount int) {
	corrCount := SKP_min_int(inputDataSize, correlationCount)
	corr64 := SKP_Silk_inner_prod16_aligned_64(inputData, inputData, int32(inputDataSize))

	/* deal with all-zero input data */
	corr64 += 1

	lz := SKP_Silk_CLZ64(corr64)

	nRightShifts := 35 - lz
	*scale = int(nRightShifts)

	if int64(nRightShifts) <= 0 {
		results[0] = int32(corr64) << -nRightShifts
		for i := 1; i < corrCount; i++ {
			results[i] = SKP_Silk_inner_prod_aligned(inputData[:inputDataSize-i], inputData[i:], int32(inputDataSize-i)) << -nRightShifts
		}
	} else {
		results[0] = int32(corr64 >> nRightShifts)
		for i := 1; i < corrCount; i++ {
			results[i] = int32(SKP_Silk_inner_prod16_aligned_64(inputData[:inputDataSize-i], inputData[i:], int32(inputDataSize-i))) >> nRightShifts
		}
	}
}
