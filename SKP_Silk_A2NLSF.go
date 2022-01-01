package silk

import "math"

const BIN_DIV_STEPS_A2NLSF_FIX = 3
const QPoly = 16
const MAX_ITERATIONS_A2NLSF_FIX = 30
const OVERSAMPLE_COSINE_TABLE = false

func SKP_Silk_A2NLSF_trans_poly(p []int32, dd int) {
	for k := 2; k < dd; k++ {
		for n := dd; n > k; n-- {
			p[n-2] -= p[n]
		}
		p[k-2] -= p[k] << 1
	}
}
func SKP_Silk_A2NLSF_eval_poly(p []int32, x int32, dd int32) int32 {
	y32 := p[dd]
	x_Q16 := x << 4
	for n := dd - 1; n >= 0; n-- {
		y32 = SKP_SMLAWW(p[n], y32, x_Q16)
	}
	return y32
}

func SKP_Silk_A2NLSF_init(a_Q16 []int32, P []int32, Q []int32, dd int) {
	/* Convert filter coefs to even and odd polynomials */
	P[dd], Q[dd] = 1<<QPoly, 1<<QPoly
	for k := 0; k < dd; k++ {
		if QPoly < 16 {
			if 16-QPoly == 1 {
				P[k] = ((-a_Q16[dd-k-1] - a_Q16[dd+k]) >> 1) + ((-a_Q16[dd-k-1]-a_Q16[dd+k])&1)>>1
			} else {
				P[k] = (((-a_Q16[dd-k-1]-a_Q16[dd+k])>>(16-QPoly) - 1) + 1) >> 1
			}
		} else if QPoly == 16 {
			P[k], Q[k] = -a_Q16[dd-k-1]-a_Q16[dd+k], -a_Q16[dd-k-1]+a_Q16[dd+k] // QPoly
		} else {
			P[k] = (-a_Q16[dd-k-1] - a_Q16[dd+k]) << (QPoly - 16) /* QPoly */
			Q[k] = (-a_Q16[dd-k-1] + a_Q16[dd+k]) << (QPoly - 16) /* QPoly */
		}
	}
	/* Divide out zeros as we have that for even filter orders, */
	/* z =  1 is always a root in Q, and                        */
	/* z = -1 is always a root in P                             */
	for k := dd; k > 0; k-- {
		P[k-1] -= P[k]
		Q[k-1] += Q[k]
	}

	SKP_Silk_A2NLSF_trans_poly(P, dd)
	SKP_Silk_A2NLSF_trans_poly(Q, dd)
}

func SKP_Silk_A2NLSF(NLSF []int, a_Q16 []int32, d int) {
	var root_ix, ffrac int
	var xlo, xhi, xmid, ylo, yhi, ymid int32
	var nom, den int32
	P := make([]int32, MAX_ORDER_LPC/2+1)
	Q := make([]int32, MAX_ORDER_LPC/2+1)
	PQ := []*[]int32{&P, &Q}

	p := &P

	dd := d >> 1

	SKP_Silk_A2NLSF_init(a_Q16, P[:], Q[:], dd)

	xlo = int32(SKP_Silk_LSFCosTab_FIX_Q12[0])
	ylo = SKP_Silk_A2NLSF_eval_poly(*p, int32(xlo), dd)

	if ylo < 0 {
		NLSF[0] = 0
		p = &Q
		ylo = SKP_Silk_A2NLSF_eval_poly(*p, int32(xlo), dd)
		root_ix = 1
	} else {
		root_ix = 0
	}

	k := 0
	i := 0
	for {
		if OVERSAMPLE_COSINE_TABLE {
			xhi = SKP_Silk_LSFCosTab_FIX_Q12[k>>1] +
				((SKP_Silk_LSFCosTab_FIX_Q12[(k+1)>>1] -
					SKP_Silk_LSFCosTab_FIX_Q12[k>>1]) >> 1) /* Q12 */
		} else {
			xhi = SKP_Silk_LSFCosTab_FIX_Q12[k]
		}
		yhi = SKP_Silk_A2NLSF_eval_poly(*p, xhi, dd)

		/* Detect zero crossing */
		if (ylo <= 0 && yhi >= 0) || (ylo >= 0 && yhi <= 0) {
			/* Binary division */
			if OVERSAMPLE_COSINE_TABLE {
				ffrac = -128
			} else {
				ffrac = -256
			}

			for m := 0; m < BIN_DIV_STEPS_A2NLSF_FIX; m++ {
				/* Evaluate polynomial */
				xmid = SKP_RSHIFT_ROUND(xlo+xhi, 1)
				ymid = SKP_Silk_A2NLSF_eval_poly(*p, xmid, dd)
				/* Detect zero crossing */
				if (ylo <= 0 && ymid >= 0) || ylo >= 0 && ymid <= 0 {
					/* Reduce frequency */
					xhi = xmid
					yhi = ymid
				} else {
					/* Increase frequency */
					xlo = xmid
					ylo = ymid

					if OVERSAMPLE_COSINE_TABLE {
						ffrac = ffrac + (64 >> m)
					} else {
						ffrac = ffrac + (128 >> m)
					}
				}
			}

			if SKP_abs(int64(ylo)) < 0x10000 {
				/* Avoid dividing by zero */
				den = ylo - yhi
				nom = (ylo << (8 - BIN_DIV_STEPA_A2NLSF_FIX)) + (den >> 1)
				if den != 0 {
					ffrac += int(nom / den)
				}
			} else {
				/* No risk of dividing by zero because abs(ylo - yhi) >= abs(ylo) >= 65536 */
				ffrac += int(ylo / ((ylo - yhi) >> (8 - BIN_DIV_STEPA_A2NLSF_FIX)))
			}
			if OVERSAMPLE_COSINE_TABLE {
				NLSF[root_ix] = int(SKP_min_32((int32(k)<<7)+int32(ffrac), math.MaxInt16))
			} else {
				NLSF[root_ix] = int(SKP_min_32((int32(k)<<8)+int32(ffrac), math.MaxInt16))
			}

			/* Next root */
			root_ix++
			if root_ix >= d {
				/* Found all roots */
				break
			}
			/* Alternate pointer to polynomial */
			p = PQ[root_ix&1]

			/* Evaluate polynomial */
			if OVERSAMPLE_COSINE_TABLE {
				xlo = int32(SKP_Silk_LSFCosTab_FIX_Q12[(k-1)>>1] +
					((SKP_Silk_LSFCosTab_FIX_Q12[k>>1] -
						SKP_Silk_LSFCosTab_FIX_Q12[(k-1)>>1]) >> 1)) // Q12
			} else {
				xlo = int32(SKP_Silk_LSFCosTab_FIX_Q12[k-1]) // Q12
			}
			ylo = int32((1 - (root_ix & 2)) << 12)
		} else {
			/* Increment loop counter */
			k++
			xlo = xhi
			ylo = yhi

			if k > LSF_COS_TAB_SZ_FIX {
				i++
				if i > MAX_ITERATIONS_A2NLSF_FIX {
					NLSF[0] = int(int32((1 << 15) / (d + 1)))
					for k = 1; k < d; k++ {
						NLSF[k] = int(SKP_SMULBB(int32(int64(k)+1), int32(NLSF[0])))
					}
					return
				}

				SKP_Silk_bwexpander_32(a_Q16, d, int32(0x10000-int64(SKP_SMULBB(int32(i+10), int32(i)))))

				SKP_Silk_A2NLSF_init(a_Q16, P, Q, dd)
				p = &P
				xlo = SKP_Silk_LSFCosTab_FIX_Q12[0]
				ylo = SKP_Silk_A2NLSF_eval_poly(*p, xlo, dd)
				if int64(ylo) < 0 {
					/* Set the first NLSF to zero and move on to the next */
					NLSF[0] = 0
					p = &Q
					ylo = SKP_Silk_A2NLSF_eval_poly(*p, xlo, dd)
					root_ix = 1
				} else {
					root_ix = 0
				}
				k = 1
			}
		}
	}
}
