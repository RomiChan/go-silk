package silk

import "math"

func SKP_Silk_encode_signs(sRC *SKP_Silk_range_coder_state, q []int8, length int32, sigtype int32, QuantOffsetType int32, RateLevelIndex int32) {
	var (
		i      int32
		inData int32
		cdf    [3]uint16
	)
	i = SKP_SMULBB(N_RATE_LEVELS-1, (sigtype<<1)+QuantOffsetType) + RateLevelIndex
	cdf[0] = 0
	cdf[1] = SKP_Silk_sign_CDF[i]
	cdf[2] = math.MaxUint16
	for i = 0; i < length; i++ {
		if q[i] != 0 {
			inData = int32(((q[i]) >> 15) + 1)
			SKP_Silk_range_encoder(sRC, inData, cdf[:])
		}
	}
}
func SKP_Silk_decode_signs(sRC *SKP_Silk_range_coder_state, q []int32, length int32, sigtype int32, QuantOffsetType int32, RateLevelIndex int32) {
	var (
		i    int32
		data int32
		cdf  [3]uint16
	)
	i = SKP_SMULBB(N_RATE_LEVELS-1, (sigtype<<1)+QuantOffsetType) + RateLevelIndex
	cdf[0] = 0
	cdf[1] = SKP_Silk_sign_CDF[i]
	cdf[2] = math.MaxUint16
	for i = 0; i < length; i++ {
		if q[i] > 0 {
			SKP_Silk_range_decoder(&data, sRC, cdf[:], 1)
			q[i] *= (data << 1) - 1
		}
	}
}
