package silk

// reviewed by wdvxdr1123 2022-02-08

func SKP_Silk_k2a(A_Q24 []int32, rc_Q15 []int16, order int32) {
	var Atmp [16]int32
	for k := 0; k < int(order); k++ {
		for n := 0; n < k; n++ {
			Atmp[n] = A_Q24[n]
		}
		for n := 0; n < k; n++ {
			A_Q24[n] = SKP_SMLAWB(A_Q24[n], Atmp[int64(k)-int64(n)-1]<<1, int32(rc_Q15[k]))
		}
		A_Q24[k] = int32(-(rc_Q15[k] << 9))
	}
}

func SKP_Silk_k2a_Q16(A_Q24 []int32, rc_Q16 []int32, order int32) {
	var Atmp [16]int32
	for k := 0; k < int(order); k++ {
		for n := 0; n < k; n++ {
			Atmp[n] = A_Q24[n]
		}
		for n := 0; n < k; n++ {
			A_Q24[n] = SKP_SMLAWW(A_Q24[n], Atmp[int64(k)-int64(n)-1], rc_Q16[k])
		}
		A_Q24[k] = -(rc_Q16[k] << 8)
	}
}
