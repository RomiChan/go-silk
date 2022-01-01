package silk

import "unsafe"

const MAX_FRAME_SIZE = 544
const MAX_NB_SUBFR = 4
const N_BITS_HEAD_ROOM = 2
const MIN_RSHIFTS = -16
const MAX_RSHIFTS = 7

func SKP_Silk_burg_modified(res_nrg *int32, res_nrg_Q *int32, A_Q16 []int32, x []int16, subfr_length int32, nb_subfr int32, WhiteNoiseFrac_Q32 int32, D int32) {
	var (
		k             int32
		n             int32
		s             int32
		lz            int32
		rshifts       int32
		rshifts_extra int32
		C0            int32
		num           int32
		nrg           int32
		rc_Q31        int32
		Atmp_25       int32
		Atmp1         int32
		tmp1          int32
		tmp2          int32
		x1            int32
		x2            int32
		x_ptr         *int16
		C_first_row   [16]int32
		C_last_row    [16]int32
		Af_25         [16]int32
		CAf           [17]int32
		CAb           [17]int32
	)
	SKP_Silk_sum_sqr_shift(&C0, &rshifts, &x[0], int32(int64(nb_subfr)*int64(subfr_length)))
	if int64(rshifts) > (32 - 25) {
		C0 = int32(int64(C0) << (int64(rshifts) - (32 - 25)))
		rshifts = 32 - 25
	} else {
		lz = int32(int64(SKP_Silk_CLZ32(C0)) - 1)
		rshifts_extra = int32(N_BITS_HEAD_ROOM - int64(lz))
		if int64(rshifts_extra) > 0 {
			if int64(rshifts_extra) < ((32 - 25) - int64(rshifts)) {
				rshifts_extra = rshifts_extra
			} else {
				rshifts_extra = int32((32 - 25) - int64(rshifts))
			}
			C0 = C0 >> int64(rshifts_extra)
		} else {
			if int64(rshifts_extra) > (int64(-16 - int64(rshifts))) {
				rshifts_extra = rshifts_extra
			} else {
				rshifts_extra = int32(-16 - int64(rshifts))
			}
			C0 = int32(int64(C0) << int64(-rshifts_extra))
		}
		rshifts += rshifts_extra
	}
	memset(unsafe.Pointer(&C_first_row[0]), 0, size_t(SKP_Silk_MAX_ORDER_LPC*unsafe.Sizeof(int32(0))))
	if int64(rshifts) > 0 {
		for s = 0; int64(s) < int64(nb_subfr); s++ {
			x_ptr = &x[int64(s)*int64(subfr_length)]
			for n = 1; int64(n) < int64(D)+1; n++ {
				C_first_row[int64(n)-1] += int32(SKP_Silk_inner_prod16_aligned_64(([]int16)(x_ptr), ([]int16)((*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(n)))), int32(int64(subfr_length)-int64(n))) >> int64(rshifts))
			}
		}
	} else {
		for s = 0; int64(s) < int64(nb_subfr); s++ {
			x_ptr = &x[int64(s)*int64(subfr_length)]
			for n = 1; int64(n) < int64(D)+1; n++ {
				C_first_row[int64(n)-1] += int32(int64(SKP_Silk_inner_prod_aligned(([]int16)(x_ptr), ([]int16)((*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(n)))), int32(int64(subfr_length)-int64(n)))) << int64(-rshifts))
			}
		}
	}
	memcpy(unsafe.Pointer(&C_last_row[0]), unsafe.Pointer(&C_first_row[0]), size_t(SKP_Silk_MAX_ORDER_LPC*unsafe.Sizeof(int32(0))))
	CAb[0] = func() int32 {
		p := &CAf[0]
		CAf[0] = int32(int64(C0) + int64(SKP_SMMUL(WhiteNoiseFrac_Q32, C0)) + 1)
		return *p
	}()
	for n = 0; int64(n) < int64(D); n++ {
		if int64(rshifts) > -2 {
			for s = 0; int64(s) < int64(nb_subfr); s++ {
				x_ptr = &x[int64(s)*int64(subfr_length)]
				x1 = int32(-(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(n))))) << (16 - int64(rshifts))))
				x2 = int32(-(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)-1))))) << (16 - int64(rshifts))))
				tmp1 = int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(n))))) << (25 - 16))
				tmp2 = int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)-1))))) << (25 - 16))
				for k = 0; int64(k) < int64(n); k++ {
					C_first_row[k] = SKP_SMLAWB(C_first_row[k], x1, int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(n)-int64(k)-1)))))
					C_last_row[k] = SKP_SMLAWB(C_last_row[k], x2, int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)+int64(k))))))
					Atmp_25 = Af_25[k]
					tmp1 = SKP_SMLAWB(tmp1, Atmp_25, int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(n)-int64(k)-1)))))
					tmp2 = SKP_SMLAWB(tmp2, Atmp_25, int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)+int64(k))))))
				}
				tmp1 = int32(int64(-tmp1) << (32 - 25 - int64(rshifts)))
				tmp2 = int32(int64(-tmp2) << (32 - 25 - int64(rshifts)))
				for k = 0; int64(k) <= int64(n); k++ {
					CAf[k] = SKP_SMLAWB(CAf[k], tmp1, int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(n)-int64(k))))))
					CAb[k] = SKP_SMLAWB(CAb[k], tmp2, int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)+int64(k)-1)))))
				}
			}
		} else {
			for s = 0; int64(s) < int64(nb_subfr); s++ {
				x_ptr = &x[int64(s)*int64(subfr_length)]
				x1 = int32(-(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(n))))) << int64(-rshifts)))
				x2 = int32(-(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)-1))))) << int64(-rshifts)))
				tmp1 = int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(n))))) << 17)
				tmp2 = int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)-1))))) << 17)
				for k = 0; int64(k) < int64(n); k++ {
					C_first_row[k] = int32(int64(C_first_row[k]) + int64(x1)*int64(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(n)-int64(k)-1)))))
					C_last_row[k] = int32(int64(C_last_row[k]) + int64(x2)*int64(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)+int64(k))))))
					Atmp1 = SKP_RSHIFT_ROUND(Af_25[k], 25-17)
					tmp1 = int32(int64(tmp1) + int64(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(n)-int64(k)-1))))*int64(Atmp1))
					tmp2 = int32(int64(tmp2) + int64(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)+int64(k)))))*int64(Atmp1))
				}
				tmp1 = -tmp1
				tmp2 = -tmp2
				for k = 0; int64(k) <= int64(n); k++ {
					CAf[k] = SKP_SMLAWW(CAf[k], tmp1, int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(n)-int64(k))))))<<(int64(-rshifts)-1)))
					CAb[k] = SKP_SMLAWW(CAb[k], tmp2, int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(int16(0))*uintptr(int64(subfr_length)-int64(n)+int64(k)-1)))))<<(int64(-rshifts)-1)))
				}
			}
		}
		tmp1 = C_first_row[n]
		tmp2 = C_last_row[n]
		num = 0
		nrg = int32(int64(CAb[0]) + int64(CAf[0]))
		for k = 0; int64(k) < int64(n); k++ {
			Atmp_25 = Af_25[k]
			lz = int32(int64(SKP_Silk_CLZ32(int32(SKP_abs(int64(Atmp_25))))) - 1)
			if (32 - 25) < int64(lz) {
				lz = 32 - 25
			} else {
				lz = lz
			}
			Atmp1 = int32(int64(Atmp_25) << int64(lz))
			tmp1 = int32(int64(tmp1) + (int64(SKP_SMMUL(C_last_row[int64(n)-int64(k)-1], Atmp1)) << (32 - 25 - int64(lz))))
			tmp2 = int32(int64(tmp2) + (int64(SKP_SMMUL(C_first_row[int64(n)-int64(k)-1], Atmp1)) << (32 - 25 - int64(lz))))
			num = int32(int64(num) + (int64(SKP_SMMUL(CAb[int64(n)-int64(k)], Atmp1)) << (32 - 25 - int64(lz))))
			nrg = int32(int64(nrg) + (int64(SKP_SMMUL(int32(int64(CAb[int64(k)+1])+int64(CAf[int64(k)+1])), Atmp1)) << (32 - 25 - int64(lz))))
		}
		CAf[int64(n)+1] = tmp1
		CAb[int64(n)+1] = tmp2
		num = int32(int64(num) + int64(tmp2))
		num = int32(int64(-num) << 1)
		if SKP_abs(int64(num)) < int64(nrg) {
			rc_Q31 = SKP_DIV32_varQ(num, nrg, 31)
		} else {
			memset(unsafe.Pointer(&Af_25[n]), 0, size_t((int64(D)-int64(n))*int64(unsafe.Sizeof(int32(0)))))
			break
		}
		for k = 0; int64(k) < (int64(n)+1)>>1; k++ {
			tmp1 = Af_25[k]
			tmp2 = Af_25[int64(n)-int64(k)-1]
			Af_25[k] = int32(int64(tmp1) + (int64(SKP_SMMUL(tmp2, rc_Q31)) << 1))
			Af_25[int64(n)-int64(k)-1] = int32(int64(tmp2) + (int64(SKP_SMMUL(tmp1, rc_Q31)) << 1))
		}
		Af_25[n] = rc_Q31 >> (31 - 25)
		for k = 0; int64(k) <= int64(n)+1; k++ {
			tmp1 = CAf[k]
			tmp2 = CAb[int64(n)-int64(k)+1]
			CAf[k] = int32(int64(tmp1) + (int64(SKP_SMMUL(tmp2, rc_Q31)) << 1))
			CAb[int64(n)-int64(k)+1] = int32(int64(tmp2) + (int64(SKP_SMMUL(tmp1, rc_Q31)) << 1))
		}
	}
	nrg = CAf[0]
	tmp1 = 1 << 16
	for k = 0; int64(k) < int64(D); k++ {
		Atmp1 = SKP_RSHIFT_ROUND(Af_25[k], 25-16)
		nrg = SKP_SMLAWW(nrg, CAf[int64(k)+1], Atmp1)
		tmp1 = SKP_SMLAWW(tmp1, Atmp1, Atmp1)
		A_Q16[k] = -Atmp1
	}
	*res_nrg = SKP_SMLAWW(nrg, SKP_SMMUL(WhiteNoiseFrac_Q32, C0), -tmp1)
	*res_nrg_Q = -rshifts
}
