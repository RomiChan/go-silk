package silk

// reviewed by wdvxdr1123 2022-02-08

var freq_table_Q16 = [27]int16{
	0x2F4F, 9804, 8235, 7100, 6239, 5565, 5022, 4575, 4202,
	3885, 3612, 3375, 3167, 2984, 2820, 2674, 2542, 2422,
	2313, 2214, 2123, 2038, 1961, 1889, 1822, 1760, 1702,
}

func SKP_Silk_apply_sine_window(
	px_win []int16, /* O    Pointer to windowed signal                  */
	px []int16, /* I    Pointer to input signal                     */
	win_type int, /* I    Selects a window type                       */
	length int, /* I    Window length, multiple of 4                */
) {
	SKP_assert(win_type == 1 || win_type == 2)
	/* Length must be in a range from 16 to 120 and a multiple of 4 */
	SKP_assert(length >= 16 && length <= 120)
	SKP_assert((length & 3) == 0)

	/* Input pointer must be 4-byte aligned */
	// SKP_assert( ( ( SKP_int64 )( ( SKP_int8* )px - ( SKP_int8* )0 ) & 3 ) == 0 );
	var S0_Q16, S1_Q16 int32

	/* Frequency */
	k := (length >> 2) - 4
	SKP_assert(k >= 0 && k <= 26)
	f_Q16 := int32(freq_table_Q16[k])

	/* Factor used for cosine approximation */
	c_Q16 := SKP_SMULWB(f_Q16, -f_Q16)
	SKP_assert(c_Q16 >= -32768)

	/* initialize state */
	if win_type == 1 {
		/* start from 0 */
		S0_Q16 = 0
		/* approximation of sin(f) */
		S1_Q16 = f_Q16 + int32(length>>3)
	} else {
		/* start from 1 */
		S0_Q16 = 1 << 16
		/* approximation of cos(f) */
		S1_Q16 = int32(int(c_Q16>>1) + (1 << 16) + (length >> 4))
	}

	/* Uses the recursive equation:   sin(n*f) = 2 * cos(f) * sin((n-1)*f) - sin((n-2)*f)    */
	/* 4 samples at a time */
	for k := 0; k < length; k += 4 {
		px_win[k] = int16(SKP_SMULWB((S0_Q16+S1_Q16)>>1, int32(px[k])))
		px_win[k+1] = int16(SKP_SMULWB(S1_Q16, int32(px[k+1])))
		S0_Q16 = SKP_SMULWB(S1_Q16, c_Q16) + (S1_Q16 << 1) - S0_Q16 + 1
		if S0_Q16 >= (1 << 16) {
			S0_Q16 = 1 << 16
		}
		px_win[k+2] = int16(SKP_SMULWB((S0_Q16+S1_Q16)>>1, int32(px[k+2])))
		px_win[k+3] = int16(SKP_SMULWB(S0_Q16, int32(px[k+3])))
		S1_Q16 = SKP_SMULWB(S0_Q16, c_Q16) + (S0_Q16 << 1) - S1_Q16
		if S1_Q16 >= (1 << 16) {
			S1_Q16 = 1 << 16
		}
	}
}
