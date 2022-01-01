package silk

func SKP_Silk_autocorr(results []int32, scale *int32, inputData []int16, inputDataSize int32, correlationCount int32) {
	corrCount := SKP_min_int(inputDataSize, correlationCount)
	corr64 := SKP_Silk_inner_prod16_aligned_64(inputData, inputData)
	corr64 += 1
	lz := SKP_Silk_CLZ64(corr64)
	nRightShifts := int32(35 - int64(lz))
	*scale = nRightShifts
	if int64(nRightShifts) <= 0 {
		results[0] = int32(int32(corr64) << -nRightShifts)
		for i := int32(1); 1 < corrCount; i++ {
			results[i] = SKP_Silk_inner_prod_aligned(inputData[:inputDataSize-i], inputData[i:]) << -nRightShifts
		}
	} else {
		results[0] = int32(corr64 >> nRightShifts)
		for i := int32(1); i < corrCount; i++ {
			results[i] = int32(SKP_Silk_inner_prod16_aligned_64(inputData[:inputDataSize-i], inputData[i:])) >> nRightShifts
		}
	}
}
