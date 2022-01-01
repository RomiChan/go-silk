package silk

import "unsafe"

func SKP_Silk_residual_energy16_covar_FIX(c *int16, wXX *int32, wXx *int32, wxx int32, D int32, cQ int32) int32 {
	var (
		i       int32
		j       int32
		lshifts int32
		Qxtra   int32
		c_max   int32
		w_max   int32
		tmp     int32
		tmp2    int32
		nrg     int32
		cn      [16]int32
		pRow    *int32
	)
	lshifts = int32(16 - int64(cQ))
	Qxtra = lshifts
	c_max = 0
	for i = 0; int64(i) < int64(D); i++ {
		c_max = SKP_max_32(c_max, int32(SKP_abs(int64(*(*int16)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(int16(0))*uintptr(i)))))))
	}
	Qxtra = SKP_min_int(Qxtra, int32(int64(SKP_Silk_CLZ32(c_max))-17))
	w_max = SKP_max_32(*(*int32)(unsafe.Add(unsafe.Pointer(wXX), unsafe.Sizeof(int32(0))*0)), *(*int32)(unsafe.Add(unsafe.Pointer(wXX), unsafe.Sizeof(int32(0))*uintptr(int64(D)*int64(D)-1))))
	Qxtra = SKP_min_int(Qxtra, int32(int64(SKP_Silk_CLZ32(int32(int64(D)*(int64(SKP_SMULWB(w_max, c_max))>>4))))-5))
	Qxtra = SKP_max_int(Qxtra, 0)
	for i = 0; int64(i) < int64(D); i++ {
		cn[i] = int32(int64(int32(*(*int16)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(int16(0))*uintptr(i))))) << int64(Qxtra))
	}
	lshifts -= Qxtra
	tmp = 0
	for i = 0; int64(i) < int64(D); i++ {
		tmp = SKP_SMLAWB(tmp, *(*int32)(unsafe.Add(unsafe.Pointer(wXx), unsafe.Sizeof(int32(0))*uintptr(i))), cn[i])
	}
	nrg = int32((int64(wxx) >> (int64(lshifts) + 1)) - int64(tmp))
	tmp2 = 0
	for i = 0; int64(i) < int64(D); i++ {
		tmp = 0
		pRow = (*int32)(unsafe.Add(unsafe.Pointer(wXX), unsafe.Sizeof(int32(0))*uintptr(int64(i)*int64(D))))
		for j = int32(int64(i) + 1); int64(j) < int64(D); j++ {
			tmp = SKP_SMLAWB(tmp, *(*int32)(unsafe.Add(unsafe.Pointer(pRow), unsafe.Sizeof(int32(0))*uintptr(j))), cn[j])
		}
		tmp = SKP_SMLAWB(tmp, int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(pRow), unsafe.Sizeof(int32(0))*uintptr(i))))>>1), cn[i])
		tmp2 = SKP_SMLAWB(tmp2, tmp, cn[i])
	}
	nrg = int32(int64(nrg) + (int64(tmp2) << int64(lshifts)))
	if int64(nrg) < 1 {
		nrg = 1
	} else if int64(nrg) > (SKP_int32_MAX >> (int64(lshifts) + 2)) {
		nrg = SKP_int32_MAX >> 1
	} else {
		nrg = int32(int64(nrg) << (int64(lshifts) + 1))
	}
	return nrg
}
