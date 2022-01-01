package silk

import "unsafe"

func gcd(a int32, b int32) int32 {
	var tmp int32
	for int64(b) > 0 {
		tmp = int32(int64(a) - int64(b)*int64(int32(int64(a)/int64(b))))
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
	if int64(Fs_Hz_in) < 8000 || int64(Fs_Hz_in) > 192000 || int64(Fs_Hz_out) < 8000 || int64(Fs_Hz_out) > 192000 {
		return -1
	}
	if int64(Fs_Hz_in) > 96000 {
		S.NPreDownsamplers = 2
		S.Down_pre_function = SKP_Silk_resampler_private_down4
	} else if int64(Fs_Hz_in) > 48000 {
		S.NPreDownsamplers = 1
		S.Down_pre_function = SKP_Silk_resampler_down2
	} else {
		S.NPreDownsamplers = 0
		S.Down_pre_function = nil
	}
	if int64(Fs_Hz_out) > 96000 {
		S.NPostUpsamplers = 2
		S.Up_post_function = SKP_Silk_resampler_private_up4
	} else if int64(Fs_Hz_out) > 48000 {
		S.NPostUpsamplers = 1
		S.Up_post_function = SKP_Silk_resampler_up2
	} else {
		S.NPostUpsamplers = 0
		S.Up_post_function = nil
	}
	if int64(S.NPreDownsamplers)+int64(S.NPostUpsamplers) > 0 {
		S.Ratio_Q16 = int32(int64(int32((int64(Fs_Hz_out)<<13)/int64(Fs_Hz_in))) << 3)
		for int64(SKP_SMULWW(S.Ratio_Q16, Fs_Hz_in)) < int64(Fs_Hz_out) {
			S.Ratio_Q16++
		}
		S.BatchSizePrePost = int32(int64(Fs_Hz_in) / 100)
		Fs_Hz_in = Fs_Hz_in >> int64(S.NPreDownsamplers)
		Fs_Hz_out = Fs_Hz_out >> int64(S.NPostUpsamplers)
	}
	S.BatchSize = int32(int64(Fs_Hz_in) / 100)
	if (int64(S.BatchSize)*100) != int64(Fs_Hz_in) || int64(Fs_Hz_in)%100 != 0 {
		cycleLen = int32(int64(Fs_Hz_in) / int64(gcd(Fs_Hz_in, Fs_Hz_out)))
		cyclesPerBatch = int32(RESAMPLER_MAX_BATCH_SIZE_IN / int64(cycleLen))
		if int64(cyclesPerBatch) == 0 {
			S.BatchSize = RESAMPLER_MAX_BATCH_SIZE_IN
		} else {
			S.BatchSize = int32(int64(cyclesPerBatch) * int64(cycleLen))
		}
	}
	if int64(Fs_Hz_out) > int64(Fs_Hz_in) {
		if int64(Fs_Hz_out) == (int64(Fs_Hz_in) * 2) {
			S.Resampler_function = SKP_Silk_resampler_private_up2_HQ_wrapper
		} else {
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_IIR_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
			up2 = 1
			if int64(Fs_Hz_in) > 24000 {
				S.Up2_function = SKP_Silk_resampler_up2
			} else {
				S.Up2_function = SKP_Silk_resampler_private_up2_HQ
			}
		}
	} else if int64(Fs_Hz_out) < int64(Fs_Hz_in) {
		if (int64(Fs_Hz_out) * 4) == (int64(Fs_Hz_in) * 3) {
			S.FIR_Fracs = 3
			S.Coefs = &SKP_Silk_Resampler_3_4_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_down_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 3) == (int64(Fs_Hz_in) * 2) {
			S.FIR_Fracs = 2
			S.Coefs = &SKP_Silk_Resampler_2_3_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_down_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 2) == int64(Fs_Hz_in) {
			S.FIR_Fracs = 1
			S.Coefs = &SKP_Silk_Resampler_1_2_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_down_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 8) == (int64(Fs_Hz_in) * 3) {
			S.FIR_Fracs = 3
			S.Coefs = &SKP_Silk_Resampler_3_8_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_down_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 3) == int64(Fs_Hz_in) {
			S.FIR_Fracs = 1
			S.Coefs = &SKP_Silk_Resampler_1_3_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_down_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 4) == int64(Fs_Hz_in) {
			S.FIR_Fracs = 1
			down2 = 1
			S.Coefs = &SKP_Silk_Resampler_1_2_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_down_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 6) == int64(Fs_Hz_in) {
			S.FIR_Fracs = 1
			down2 = 1
			S.Coefs = &SKP_Silk_Resampler_1_3_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_down_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 441) == (int64(Fs_Hz_in) * 80) {
			S.Coefs = &SKP_Silk_Resampler_80_441_ARMA4_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_IIR_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 441) == (int64(Fs_Hz_in) * 120) {
			S.Coefs = &SKP_Silk_Resampler_120_441_ARMA4_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_IIR_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 441) == (int64(Fs_Hz_in) * 160) {
			S.Coefs = &SKP_Silk_Resampler_160_441_ARMA4_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_IIR_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 441) == (int64(Fs_Hz_in) * 240) {
			S.Coefs = &SKP_Silk_Resampler_240_441_ARMA4_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_IIR_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else if (int64(Fs_Hz_out) * 441) == (int64(Fs_Hz_in) * 320) {
			S.Coefs = &SKP_Silk_Resampler_320_441_ARMA4_COEFS[0]
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_IIR_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
		} else {
			S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
				SKP_Silk_resampler_private_IIR_FIR(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
			}
			up2 = 1
			if int64(Fs_Hz_in) > 24000 {
				S.Up2_function = SKP_Silk_resampler_up2
			} else {
				S.Up2_function = SKP_Silk_resampler_private_up2_HQ
			}
		}
	} else {
		S.Resampler_function = func(arg1 unsafe.Pointer, arg2 *int16, arg3 *int16, arg4 int32) {
			SKP_Silk_resampler_private_copy(arg1, ([]int16)(arg2), ([]int16)(arg3), arg4)
		}
	}
	S.Input2x = int32(int64(up2) | int64(down2))
	S.InvRatio_Q16 = int32(int64(int32((int64(Fs_Hz_in)<<(int64(up2)+14-int64(down2)))/int64(Fs_Hz_out))) << 2)
	for int64(SKP_SMULWW(S.InvRatio_Q16, int32(int64(Fs_Hz_out)<<int64(down2)))) < (int64(Fs_Hz_in) << int64(up2)) {
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
func SKP_Silk_resampler(S *SKP_Silk_resampler_state_struct, out []int16, in []int16, inLen int32) int32 {
	if int64(S.Magic_number) != 0x75BCD15 {
		return -1
	}
	if int64(S.NPreDownsamplers)+int64(S.NPostUpsamplers) > 0 {
		var (
			nSamplesIn  int32
			nSamplesOut int32
			in_buf      [480]int16
			out_buf     [480]int16
		)
		for int64(inLen) > 0 {
			if int64(inLen) < int64(S.BatchSizePrePost) {
				nSamplesIn = inLen
			} else {
				nSamplesIn = S.BatchSizePrePost
			}
			nSamplesOut = SKP_SMULWB(S.Ratio_Q16, nSamplesIn)
			if int64(S.NPreDownsamplers) > 0 {
				S.Down_pre_function(&S.SDownPre[0], &in_buf[0], &in[0], nSamplesIn)
				if int64(S.NPostUpsamplers) > 0 {
					S.Resampler_function(unsafe.Pointer(S), &out_buf[0], &in_buf[0], int32(int64(nSamplesIn)>>int64(S.NPreDownsamplers)))
					S.Up_post_function(&S.SUpPost[0], &out[0], &out_buf[0], int32(int64(nSamplesOut)>>int64(S.NPostUpsamplers)))
				} else {
					S.Resampler_function(unsafe.Pointer(S), &out[0], &in_buf[0], int32(int64(nSamplesIn)>>int64(S.NPreDownsamplers)))
				}
			} else {
				S.Resampler_function(unsafe.Pointer(S), &out_buf[0], &in[0], int32(int64(nSamplesIn)>>int64(S.NPreDownsamplers)))
				S.Up_post_function(&S.SUpPost[0], &out[0], &out_buf[0], int32(int64(nSamplesOut)>>int64(S.NPostUpsamplers)))
			}
			in += ([]int16)(nSamplesIn)
			out += ([]int16)(nSamplesOut)
			inLen -= nSamplesIn
		}
	} else {
		S.Resampler_function(unsafe.Pointer(S), &out[0], &in[0], inLen)
	}
	return 0
}
