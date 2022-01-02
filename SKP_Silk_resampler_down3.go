package silk

// Reviewed by wdvxdr1123 2022-01-02

// SKP_Silk_resampler_down3
// Downsample by a factor 3, low quality
func SKP_Silk_resampler_down3(
	S []int32, //   I/O: State vector [ 8 ]
	out []int16, // O:   Output signal [ floor(inLen/3) ]
	in []int16, //  I:   Input signal [ inLen ]
	inLen int32, // I:   Number of input samples
) {
	var buf [505]int32 // todo: compute size
	var nSamplesIn int32
	out_ind := 0

	/* Copy buffered samples to start of buffer */
	copy(buf[:25], S[:])

	/* Iterate over blocks of frameSizeIn input samples */
	for {
		if inLen < RESAMPLER_MAX_BATCH_SIZE_IN {
			nSamplesIn = inLen
		} else {
			nSamplesIn = RESAMPLER_MAX_BATCH_SIZE_IN
		}

		/* Second-order AR filter (output in Q8) */
		SKP_Silk_resampler_private_AR2(S[25:], buf[25:], in,
			SKP_Silk_Resampler_1_3_COEFS_LQ[:], nSamplesIn)

		/* Interpolate filtered signal */
		buf_ptr := buf[:]
		counter := nSamplesIn
		for counter > 2 {
			/* Inner product */
			res_Q6 := SKP_SMULWB(buf_ptr[0]+buf_ptr[5], int32(SKP_Silk_Resampler_1_3_COEFS_LQ[2]))
			res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[1]+buf_ptr[4], int32(SKP_Silk_Resampler_1_3_COEFS_LQ[3]))
			res_Q6 = SKP_SMLAWB(res_Q6, buf_ptr[2]+buf_ptr[3], int32(SKP_Silk_Resampler_1_3_COEFS_LQ[4]))

			/* Scale down, saturate and store in output array */
			out[out_ind] = SKP_SAT16(SKP_RSHIFT_ROUND(res_Q6, 6))
			out_ind++

			buf_ptr = buf_ptr[3:]
			counter -= 3
		}

		in = in[nSamplesIn:]
		inLen -= nSamplesIn

		if inLen > 0 {
			/* More iterations to do; copy last part of filtered signal to beginning of buffer */
			copy(buf[:25], buf[nSamplesIn:])
		} else {
			break
		}
	}

	/* Copy last part of filtered signal to the state for the next call */
	copy(buf[:25], buf[nSamplesIn:])
}
