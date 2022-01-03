package silk

import "math"

// Reviewed by wdvxdr1123 2022-01-03

// #define SKP_enc_map(a)                ((a) > 0 ? 1 : 0)
// #define SKP_dec_map(a)                ((a) > 0 ? 1 : -1)
// shifting avoids if-statement
// #define SKP_enc_map(a)                  ( SKP_RSHIFT( (a), 15 ) + 1 )
// #define SKP_dec_map(a)                  ( SKP_LSHIFT( (a),  1 ) - 1 )

func SKP_enc_map(a int32) int32 { return (a >> 15) + 1 }
func SKP_dec_map(a int32) int32 { return (a << 1) - 1 }

// SKP_Silk_encode_signs Encodes signs of excitation
func SKP_Silk_encode_signs(
	sRC *SKP_Silk_range_coder_state, // I/O  Range coder state
	q []int8, // I    Pulse signal
	length int32, // I    Length of input
	sigtype int32, // I    Signal type
	QuantOffsetType int32, // I    Quantization offset type
	RateLevelIndex int32, // I    Rate level index
) {
	i := SKP_SMULBB(N_RATE_LEVELS-1, (sigtype<<1)+QuantOffsetType) + RateLevelIndex
	cdf := [3]uint16{0, SKP_Silk_sign_CDF[i], math.MaxUint16}

	var inData int32
	for i = 0; i < length; i++ {
		if q[i] != 0 {
			inData = SKP_enc_map(int32(q[i])) /* - = 0, + = 1 */
			SKP_Silk_range_encoder(sRC, inData, cdf[:])
		}
	}
}

// SKP_Silk_decode_signs Decodes signs of excitation
func SKP_Silk_decode_signs(
	sRC *SKP_Silk_range_coder_state, // I/O  Range coder state
	q []int32, // I/O  pulse signal
	length int32, // I    length of output
	sigtype int32, // I    Signal type
	QuantOffsetType int32, // I    Quantization offset type
	RateLevelIndex int32, // I    Rate Level Index
) {
	i := SKP_SMULBB(N_RATE_LEVELS-1, (sigtype<<1)+QuantOffsetType) + RateLevelIndex
	cdf := [3]uint16{0, SKP_Silk_sign_CDF[i], math.MaxUint16}

	var data int32
	for i = 0; i < length; i++ {
		if q[i] > 0 {
			SKP_Silk_range_decoder(&data, sRC, cdf[:], 1)
			/* attach sign */
			/* implementation with shift, subtraction, multiplication */
			q[i] *= SKP_dec_map(data)
		}
	}
}
