package silk

import "unsafe"

type inv_D_t struct {
	Q36_part int32
	Q48_part int32
}

func SKP_Silk_solve_LDL_FIX(A *int32, M int32, b *int32, x_Q16 *int32) {
	var (
		L_Q16 [256]int32
		Y     [16]int32
		inv_D [16]inv_D_t
	)
	SKP_Silk_LDL_factorize_FIX(A, M, &L_Q16[0], &inv_D[0])
	SKP_Silk_LS_SolveFirst_FIX(&L_Q16[0], M, b, &Y[0])
	SKP_Silk_LS_divide_Q16_FIX(Y[:], &inv_D[0], M)
	SKP_Silk_LS_SolveLast_FIX(&L_Q16[0], M, &Y[0], x_Q16)
}
func SKP_Silk_LDL_factorize_FIX(A *int32, M int32, L_Q16 *int32, inv_D *inv_D_t) {
	var (
		i                int32
		j                int32
		k                int32
		status           int32
		loop_count       int32
		ptr1             *int32
		ptr2             *int32
		diag_min_value   int32
		tmp_32           int32
		err              int32
		v_Q0             [16]int32
		D_Q0             [16]int32
		one_div_diag_Q36 int32
		one_div_diag_Q40 int32
		one_div_diag_Q48 int32
	)
	status = 1
	diag_min_value = SKP_max_32(SKP_SMMUL(func() int32 {
		if (((*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*0))) + (*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(SKP_SMULBB(M, M)-1))))) & math.MinInt32) == 0 {
			if (((*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*0))) & (*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(SKP_SMULBB(M, M)-1))))) & math.MinInt32) != 0 {
				return math.MinInt32
			}
			return (*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*0))) + (*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(SKP_SMULBB(M, M)-1))))
		}
		if (((*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*0))) | (*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(SKP_SMULBB(M, M)-1))))) & math.MinInt32) == 0 {
			return SKP_int32_MAX
		}
		return (*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*0))) + (*(*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(SKP_SMULBB(M, M)-1))))
	}(), SKP_FIX_CONST(1e-05, 31)), 1<<9)
	for loop_count = 0; loop_count < M && status == 1; loop_count++ {
		status = 0
		for j = 0; j < M; j++ {
			ptr1 = (*int32)(unsafe.Add(unsafe.Pointer(L_Q16), unsafe.Sizeof(int32(0))*uintptr(j*M+0)))
			tmp_32 = 0
			for i = 0; i < j; i++ {
				v_Q0[i] = SKP_SMULWW(D_Q0[i], *(*int32)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int32(0))*uintptr(i))))
				tmp_32 = SKP_SMLAWW(tmp_32, v_Q0[i], *(*int32)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int32(0))*uintptr(i))))
			}
			tmp_32 = (*((*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(j*M+j))))) - tmp_32
			if tmp_32 < diag_min_value {
				tmp_32 = SKP_SMULBB(loop_count+1, diag_min_value) - tmp_32
				for i = 0; i < M; i++ {
					*((*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(i*M+i)))) = (*((*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(i*M+i))))) + tmp_32
				}
				status = 1
				break
			}
			D_Q0[j] = tmp_32
			one_div_diag_Q36 = SKP_INVERSE32_varQ(tmp_32, 36)
			one_div_diag_Q40 = one_div_diag_Q36 << 4
			err = (1 << 24) - SKP_SMULWW(tmp_32, one_div_diag_Q40)
			one_div_diag_Q48 = SKP_SMULWW(err, one_div_diag_Q40)
			(*(*inv_D_t)(unsafe.Add(unsafe.Pointer(inv_D), unsafe.Sizeof(inv_D_t{})*uintptr(j)))).Q36_part = one_div_diag_Q36
			(*(*inv_D_t)(unsafe.Add(unsafe.Pointer(inv_D), unsafe.Sizeof(inv_D_t{})*uintptr(j)))).Q48_part = one_div_diag_Q48
			*((*int32)(unsafe.Add(unsafe.Pointer(L_Q16), unsafe.Sizeof(int32(0))*uintptr(j*M+j)))) = 0x10000
			ptr1 = (*int32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(int32(0))*uintptr(j*M+0)))
			ptr2 = (*int32)(unsafe.Add(unsafe.Pointer(L_Q16), unsafe.Sizeof(int32(0))*uintptr((j+1)*M+0)))
			for i = j + 1; i < M; i++ {
				tmp_32 = 0
				for k = 0; k < j; k++ {
					tmp_32 = SKP_SMLAWW(tmp_32, v_Q0[k], *(*int32)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int32(0))*uintptr(k))))
				}
				tmp_32 = (*(*int32)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(int32(0))*uintptr(i)))) - tmp_32
				*((*int32)(unsafe.Add(unsafe.Pointer(L_Q16), unsafe.Sizeof(int32(0))*uintptr(i*M+j)))) = SKP_SMMUL(tmp_32, one_div_diag_Q48) + (SKP_SMULWW(tmp_32, one_div_diag_Q36) >> 4)
				ptr2 = (*int32)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(int32(0))*uintptr(M)))
			}
		}
	}
}
func SKP_Silk_LS_divide_Q16_FIX(T []int32, inv_D *inv_D_t, M int32) {
	var (
		i                int32
		tmp_32           int32
		one_div_diag_Q36 int32
		one_div_diag_Q48 int32
	)
	for i = 0; i < M; i++ {
		one_div_diag_Q36 = (*(*inv_D_t)(unsafe.Add(unsafe.Pointer(inv_D), unsafe.Sizeof(inv_D_t{})*uintptr(i)))).Q36_part
		one_div_diag_Q48 = (*(*inv_D_t)(unsafe.Add(unsafe.Pointer(inv_D), unsafe.Sizeof(inv_D_t{})*uintptr(i)))).Q48_part
		tmp_32 = T[i]
		T[i] = SKP_SMMUL(tmp_32, one_div_diag_Q48) + (SKP_SMULWW(tmp_32, one_div_diag_Q36) >> 4)
	}
}
func SKP_Silk_LS_SolveFirst_FIX(L_Q16 *int32, M int32, b *int32, x_Q16 *int32) {
	var (
		i      int32
		j      int32
		ptr32  *int32
		tmp_32 int32
	)
	for i = 0; i < M; i++ {
		ptr32 = (*int32)(unsafe.Add(unsafe.Pointer(L_Q16), unsafe.Sizeof(int32(0))*uintptr(i*M+0)))
		tmp_32 = 0
		for j = 0; j < i; j++ {
			tmp_32 = SKP_SMLAWW(tmp_32, *(*int32)(unsafe.Add(unsafe.Pointer(ptr32), unsafe.Sizeof(int32(0))*uintptr(j))), *(*int32)(unsafe.Add(unsafe.Pointer(x_Q16), unsafe.Sizeof(int32(0))*uintptr(j))))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(x_Q16), unsafe.Sizeof(int32(0))*uintptr(i))) = (*(*int32)(unsafe.Add(unsafe.Pointer(b), unsafe.Sizeof(int32(0))*uintptr(i)))) - tmp_32
	}
}
func SKP_Silk_LS_SolveLast_FIX(L_Q16 *int32, M int32, b *int32, x_Q16 *int32) {
	var (
		i      int32
		j      int32
		ptr32  *int32
		tmp_32 int32
	)
	for i = M - 1; i >= 0; i-- {
		ptr32 = (*int32)(unsafe.Add(unsafe.Pointer(L_Q16), unsafe.Sizeof(int32(0))*uintptr(M*0+i)))
		tmp_32 = 0
		for j = M - 1; j > i; j-- {
			tmp_32 = SKP_SMLAWW(tmp_32, *(*int32)(unsafe.Add(unsafe.Pointer(ptr32), unsafe.Sizeof(int32(0))*uintptr(SKP_SMULBB(j, M)))), *(*int32)(unsafe.Add(unsafe.Pointer(x_Q16), unsafe.Sizeof(int32(0))*uintptr(j))))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(x_Q16), unsafe.Sizeof(int32(0))*uintptr(i))) = (*(*int32)(unsafe.Add(unsafe.Pointer(b), unsafe.Sizeof(int32(0))*uintptr(i)))) - tmp_32
	}
}
