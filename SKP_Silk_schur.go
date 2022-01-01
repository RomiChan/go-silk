package silk

import "unsafe"

func SKP_Silk_schur(rc_Q15 *int16, c *int32, order int32) int32 {
	var (
		k          int32
		n          int32
		lz         int32
		C          [17][2]int32
		Ctmp1      int32
		Ctmp2      int32
		rc_tmp_Q15 int32
	)
	lz = SKP_Silk_CLZ32(*(*int32)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(int32(0))*0)))
	if lz < 2 {
		for k = 0; k < order+1; k++ {
			C[k][0] = func() int32 {
				p := &C[k][1]
				C[k][1] = (*(*int32)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(int32(0))*uintptr(k)))) >> 1
				return *p
			}()
		}
	} else if lz > 2 {
		lz -= 2
		for k = 0; k < order+1; k++ {
			C[k][0] = func() int32 {
				p := &C[k][1]
				C[k][1] = (*(*int32)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(int32(0))*uintptr(k)))) << lz
				return *p
			}()
		}
	} else {
		for k = 0; k < order+1; k++ {
			C[k][0] = func() int32 {
				p := &C[k][1]
				C[k][1] = *(*int32)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(int32(0))*uintptr(k)))
				return *p
			}()
		}
	}
	for k = 0; k < order; k++ {
		rc_tmp_Q15 = -((C[k+1][0]) / SKP_max_32((C[0][1])>>15, 1))
		rc_tmp_Q15 = int32(SKP_SAT16(rc_tmp_Q15))
		*(*int16)(unsafe.Add(unsafe.Pointer(rc_Q15), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(rc_tmp_Q15)
		for n = 0; n < order-k; n++ {
			Ctmp1 = C[n+k+1][0]
			Ctmp2 = C[n][1]
			C[n+k+1][0] = SKP_SMLAWB(Ctmp1, Ctmp2<<1, rc_tmp_Q15)
			C[n][1] = SKP_SMLAWB(Ctmp2, Ctmp1<<1, rc_tmp_Q15)
		}
	}
	return C[0][1]
}
