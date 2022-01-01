package silk

const BIN_DIV_STEPS_A2NLSF_FIX = 3
const QPoly = 16
const MAX_ITERATIONS_A2NLSF_FIX = 30
const OVERSAMPLE_COSINE_TABLE = 0

func SKP_Silk_A2NLSF_trans_poly(p []int32, dd int32) {
	var (
		k int32
		n int32
	)
	for k = 2; int64(k) <= int64(dd); k++ {
		for n = dd; int64(n) > int64(k); n-- {
			p[int64(n)-2] -= p[n]
		}
		p[int64(k)-2] -= int32(int64(p[k]) << 1)
	}
}
func SKP_Silk_A2NLSF_eval_poly(p []int32, x int32, dd int32) int32 {
	var (
		n     int32
		x_Q16 int32
		y32   int32
	)
	y32 = p[dd]
	x_Q16 = int32(int64(x) << 4)
	for n = int32(int64(dd) - 1); int64(n) >= 0; n-- {
		y32 = SKP_SMLAWW(p[n], y32, x_Q16)
	}
	return y32
}
func SKP_Silk_A2NLSF_init(a_Q16 []int32, P []int32, Q []int32, dd int32) {
	var k int32
	P[dd] = 1 << QPoly
	Q[dd] = 1 << QPoly
	for k = 0; int64(k) < int64(dd); k++ {
		P[k] = int32(int64(-a_Q16[int64(dd)-int64(k)-1]) - int64(a_Q16[int64(dd)+int64(k)]))
		Q[k] = int32(int64(-a_Q16[int64(dd)-int64(k)-1]) + int64(a_Q16[int64(dd)+int64(k)]))
	}
	for k = dd; int64(k) > 0; k-- {
		P[int64(k)-1] -= P[k]
		Q[int64(k)-1] += Q[k]
	}
	SKP_Silk_A2NLSF_trans_poly(P, dd)
	SKP_Silk_A2NLSF_trans_poly(Q, dd)
}
func SKP_Silk_A2NLSF(NLSF []int32, a_Q16 []int32, d int32) {
	var (
		i       int32
		k       int32
		m       int32
		dd      int32
		root_ix int32
		ffrac   int32
		xlo     int32
		xhi     int32
		xmid    int32
		ylo     int32
		yhi     int32
		ymid    int32
		nom     int32
		den     int32
		P       [9]int32
		Q       [9]int32
		PQ      [2]*int32
		p       *int32
	)
	PQ[0] = &P[0]
	PQ[1] = &Q[0]
	dd = d >> 1
	SKP_Silk_A2NLSF_init(a_Q16, P[:], Q[:], dd)
	p = &P[0]
	xlo = SKP_Silk_LSFCosTab_FIX_Q12[0]
	ylo = SKP_Silk_A2NLSF_eval_poly(([]int32)(p), xlo, dd)
	if int64(ylo) < 0 {
		NLSF[0] = 0
		p = &Q[0]
		ylo = SKP_Silk_A2NLSF_eval_poly(([]int32)(p), xlo, dd)
		root_ix = 1
	} else {
		root_ix = 0
	}
	k = 1
	i = 0
	for {
		xhi = SKP_Silk_LSFCosTab_FIX_Q12[k]
		yhi = SKP_Silk_A2NLSF_eval_poly(([]int32)(p), xhi, dd)
		if int64(ylo) <= 0 && int64(yhi) >= 0 || int64(ylo) >= 0 && int64(yhi) <= 0 {
			ffrac = -256
			for m = 0; int64(m) < BIN_DIV_STEPS_A2NLSF_FIX; m++ {
				xmid = SKP_RSHIFT_ROUND(int32(int64(xlo)+int64(xhi)), 1)
				ymid = SKP_Silk_A2NLSF_eval_poly(([]int32)(p), xmid, dd)
				if int64(ylo) <= 0 && int64(ymid) >= 0 || int64(ylo) >= 0 && int64(ymid) <= 0 {
					xhi = xmid
					yhi = ymid
				} else {
					xlo = xmid
					ylo = ymid
					ffrac = int32(int64(ffrac) + (128 >> int64(m)))
				}
			}
			if SKP_abs(int64(ylo)) < 0x10000 {
				den = int32(int64(ylo) - int64(yhi))
				nom = int32((int64(ylo) << (8 - BIN_DIV_STEPS_A2NLSF_FIX)) + (int64(den) >> 1))
				if int64(den) != 0 {
					ffrac += int32(int64(nom) / int64(den))
				}
			} else {
				ffrac += int32(int64(ylo) / ((int64(ylo) - int64(yhi)) >> (8 - BIN_DIV_STEPS_A2NLSF_FIX)))
			}
			NLSF[root_ix] = SKP_min_32(int32((int64(k)<<8)+int64(ffrac)), SKP_int16_MAX)
			root_ix++
			if int64(root_ix) >= int64(d) {
				break
			}
			p = PQ[int64(root_ix)&1]
			xlo = SKP_Silk_LSFCosTab_FIX_Q12[int64(k)-1]
			ylo = int32((1 - (int64(root_ix) & 2)) << 12)
		} else {
			k++
			xlo = xhi
			ylo = yhi
			if int64(k) > LSF_COS_TAB_SZ_FIX {
				i++
				if int64(i) > MAX_ITERATIONS_A2NLSF_FIX {
					NLSF[0] = int32((1 << 15) / (int64(d) + 1))
					for k = 1; int64(k) < int64(d); k++ {
						NLSF[k] = SKP_SMULBB(int32(int64(k)+1), NLSF[0])
					}
					return
				}
				SKP_Silk_bwexpander_32(&a_Q16[0], d, int32(0x10000-int64(SKP_SMULBB(int32(int64(i)+10), i))))
				SKP_Silk_A2NLSF_init(a_Q16, P[:], Q[:], dd)
				p = &P[0]
				xlo = SKP_Silk_LSFCosTab_FIX_Q12[0]
				ylo = SKP_Silk_A2NLSF_eval_poly(([]int32)(p), xlo, dd)
				if int64(ylo) < 0 {
					NLSF[0] = 0
					p = &Q[0]
					ylo = SKP_Silk_A2NLSF_eval_poly(([]int32)(p), xlo, dd)
					root_ix = 1
				} else {
					root_ix = 0
				}
				k = 1
			}
		}
	}
}
