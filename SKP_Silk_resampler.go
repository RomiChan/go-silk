package silk

import "unsafe"

func gcd(a int32, b int32) int32 {
	var tmp int32
	for b > 0 {
		tmp = a - b*(a/b)
		a = b
		b = tmp
	}
	return a
}
func SKP_Silk_resampler_init(S *SKP_Silk_resampler_state_struct, Fs_Hz_in int32, Fs_Hz_out int32) int32 {
	var (
		cycleLen       int32
		cyclesPerBatch int32
		up2            int32 = 0
		down2          int32 = 0
	)
	memset(unsafe.Pointer(S), 0, size_t(unsafe.Sizeof(SKP_Silk_resampler_state_struct{})))
	if Fs_Hz_in < 8000 || Fs_Hz_in > 192000 || Fs_Hz_out < 8000 || Fs_Hz_out > 192000 {
		return -1
	}
	if Fs_Hz_in > 96000 {
		S.NPreDownsamplers = 2
		S.Down_pre_function = SKP_Silk_resampler_private_down4
	} else if Fs_Hz_in > 48000 {
		S.NPreDownsamplers = 1
		S.Down_pre_function = SKP_Silk_resampler_down2
	} else {
		S.NPreDownsamplers = 0
		S.Down_pre_function = 0
	}
	if Fs_Hz_out > 96000 {
		S.NPostUpsamplers = 2
		S.Up_post_function = SKP_Silk_resampler_private_up4
	} else if Fs_Hz_out > 48000 {
		S.NPostUpsamplers = 1
		S.Up_post_function = SKP_Silk_resampler_up2
	} else {
		S.NPostUpsamplers = 0
		S.Up_post_function = 0
	}
	if S.NPreDownsamplers+S.NPostUpsamplers > 0 {
		S.Ratio_Q16 = ((Fs_Hz_out << 13) / Fs_Hz_in) << 3
		for SKP_SMULWW(S.Ratio_Q16, Fs_Hz_in) < Fs_Hz_out {
			S.Ratio_Q16++
		}
		S.BatchSizePrePost = Fs_Hz_in / 100
		Fs_Hz_in = Fs_Hz_in >> S.NPreDownsamplers
		Fs_Hz_out = Fs_Hz_out >> S.NPostUpsamplers
	}
	S.BatchSize = Fs_Hz_in / 100
	if (S.BatchSize*100) != Fs_Hz_in || Fs_Hz_in%100 != 0 {
		cycleLen = Fs_Hz_in / gcd(Fs_Hz_in, Fs_Hz_out)
		cyclesPerBatch = RESAMPLER_MAX_BATCH_SIZE_IN / cycleLen
		if cyclesPerBatch == 0 {
			S.BatchSize = RESAMPLER_MAX_BATCH_SIZE_IN
		} else {
			S.BatchSize = cyclesPerBatch * cycleLen
		}
	}
	if Fs_Hz_out > Fs_Hz_in {
		if Fs_Hz_out == (Fs_Hz_in * 2) {
			S.Resampler_function = SKP_Silk_resampler_private_up2_HQ_wrapper
		} else {
			S.Resampler_function = SKP_Silk_resampler_private_IIR_FIR

			up2 = 1
			if Fs_Hz_in > 24000 {
				S.Up2_function = SKP_Silk_resampler_up2
			} else {
				S.Up2_function = SKP_Silk_resampler_private_up2_HQ
			}
		}
	} else if Fs_Hz_out < Fs_Hz_in {
		if (Fs_Hz_out * 4) == (Fs_Hz_in * 3) {
			S.FIR_Fracs = 3
			S.Coefs = &SKP_Silk_Resampler_3_4_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_down_FIR
		} else if (Fs_Hz_out * 3) == (Fs_Hz_in * 2) {
			S.FIR_Fracs = 2
			S.Coefs = &SKP_Silk_Resampler_2_3_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_down_FIR
		} else if (Fs_Hz_out * 2) == Fs_Hz_in {
			S.FIR_Fracs = 1
			S.Coefs = &SKP_Silk_Resampler_1_2_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_down_FIR
		} else if (Fs_Hz_out * 8) == (Fs_Hz_in * 3) {
			S.FIR_Fracs = 3
			S.Coefs = &SKP_Silk_Resampler_3_8_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_down_FIR
		} else if (Fs_Hz_out * 3) == Fs_Hz_in {
			S.FIR_Fracs = 1
			S.Coefs = &SKP_Silk_Resampler_1_3_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_down_FIR
		} else if (Fs_Hz_out * 4) == Fs_Hz_in {
			S.FIR_Fracs = 1
			down2 = 1
			S.Coefs = &SKP_Silk_Resampler_1_2_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_down_FIR
		} else if (Fs_Hz_out * 6) == Fs_Hz_in {
			S.FIR_Fracs = 1
			down2 = 1
			S.Coefs = &SKP_Silk_Resampler_1_3_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_down_FIR
		} else if (Fs_Hz_out * 441) == (Fs_Hz_in * 80) {
			S.Coefs = &SKP_Silk_Resampler_80_441_ARMA4_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_IIR_FIR
		} else if (Fs_Hz_out * 441) == (Fs_Hz_in * 120) {
			S.Coefs = &SKP_Silk_Resampler_120_441_ARMA4_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_IIR_FIR
		} else if (Fs_Hz_out * 441) == (Fs_Hz_in * 160) {
			S.Coefs = &SKP_Silk_Resampler_160_441_ARMA4_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_IIR_FIR
		} else if (Fs_Hz_out * 441) == (Fs_Hz_in * 240) {
			S.Coefs = &SKP_Silk_Resampler_240_441_ARMA4_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_IIR_FIR
		} else if (Fs_Hz_out * 441) == (Fs_Hz_in * 320) {
			S.Coefs = &SKP_Silk_Resampler_320_441_ARMA4_COEFS[0]
			S.Resampler_function = SKP_Silk_resampler_private_IIR_FIR
		} else {
			S.Resampler_function = SKP_Silk_resampler_private_IIR_FIR
			up2 = 1
			if Fs_Hz_in > 24000 {
				S.Up2_function = SKP_Silk_resampler_up2
			} else {
				S.Up2_function = SKP_Silk_resampler_private_up2_HQ
			}
		}
	} else {
		S.Resampler_function = SKP_Silk_resampler_private_copy
	}
	S.Input2x = up2 | down2
	S.InvRatio_Q16 = ((Fs_Hz_in << (up2 + 14 - down2)) / Fs_Hz_out) << 2
	for SKP_SMULWW(S.InvRatio_Q16, Fs_Hz_out<<down2) < (Fs_Hz_in << up2) {
		S.InvRatio_Q16++
	}
	S.Magic_number = 0x75BCD15
	return 0
}
func SKP_Silk_resampler_clear(S *SKP_Silk_resampler_state_struct) int32 {
	memset(unsafe.Pointer(&S.SDown2[0]), 0, size_t(unsafe.Sizeof([2]int32{})))
	memset(unsafe.Pointer(&S.SIIR[0]), 0, size_t(unsafe.Sizeof([6]int32{})))
	memset(unsafe.Pointer(&S.SFIR[0]), 0, size_t(unsafe.Sizeof([16]int32{})))
	memset(unsafe.Pointer(&S.SDownPre[0]), 0, size_t(unsafe.Sizeof([2]int32{})))
	memset(unsafe.Pointer(&S.SUpPost[0]), 0, size_t(unsafe.Sizeof([2]int32{})))
	return 0
}
func SKP_Silk_resampler(S *SKP_Silk_resampler_state_struct, out []int16, in []int16, inLen int) int32 {
	if S.Magic_number != 0x75BCD15 {
		return -1
	}
	if S.NPreDownsamplers+S.NPostUpsamplers > 0 {
		var (
			nSamplesIn  int32
			nSamplesOut int32
			in_buf      [480]int16
			out_buf     [480]int16
		)
		for inLen > 0 {
			if inLen < S.BatchSizePrePost {
				nSamplesIn = inLen
			} else {
				nSamplesIn = S.BatchSizePrePost
			}
			nSamplesOut = SKP_SMULWB(S.Ratio_Q16, nSamplesIn)
			if S.NPreDownsamplers > 0 {
				S.Down_pre_function(&S.SDownPre[0], &in_buf[0], &in[0], nSamplesIn)
				if S.NPostUpsamplers > 0 {
					S.Resampler_function(S, out_buf, in_buf, nSamplesIn>>S.NPreDownsamplers)
					S.Up_post_function(&S.SUpPost[0], &out[0], &out_buf[0], nSamplesOut>>S.NPostUpsamplers)
				} else {
					S.Resampler_function(S, out, in_buf, nSamplesIn>>S.NPreDownsamplers)
				}
			} else {
				S.Resampler_function(S, out_buf, in, nSamplesIn>>S.NPreDownsamplers)
				S.Up_post_function(&S.SUpPost[0], out, out_buf, nSamplesOut>>S.NPostUpsamplers)
			}
			in += ([]int16)(nSamplesIn)
			out += ([]int16)(nSamplesOut)
			inLen -= nSamplesIn
		}
	} else {
		S.Resampler_function(S, out, in, inLen)
	}
	return 0
}
