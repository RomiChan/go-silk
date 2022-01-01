package silk

import "unsafe"

func SKP_Silk_SDK_Get_Encoder_Size(encSizeBytes *int32) int32 {
	var ret int32 = 0
	*encSizeBytes = int32(uint32(unsafe.Sizeof(SKP_Silk_encoder_state_FIX{})))
	return ret
}
func SKP_Silk_SDK_QueryEncoder(encState unsafe.Pointer, encStatus *SKP_SILK_SDK_EncControlStruct) int32 {
	var (
		psEnc *SKP_Silk_encoder_state_FIX
		ret   int32 = 0
	)
	psEnc = (*SKP_Silk_encoder_state_FIX)(encState)
	encStatus.API_sampleRate = psEnc.SCmn.API_fs_Hz
	encStatus.MaxInternalSampleRate = SKP_SMULBB(psEnc.SCmn.MaxInternal_fs_kHz, 1000)
	encStatus.PacketSize = (psEnc.SCmn.API_fs_Hz * psEnc.SCmn.PacketSize_ms) / 1000
	encStatus.BitRate = psEnc.SCmn.TargetRate_bps
	encStatus.PacketLossPercentage = psEnc.SCmn.PacketLoss_perc
	encStatus.Complexity = psEnc.SCmn.Complexity
	encStatus.UseInBandFEC = psEnc.SCmn.UseInBandFEC
	encStatus.UseDTX = psEnc.SCmn.UseDTX
	return ret
}
func SKP_Silk_SDK_InitEncoder(encState unsafe.Pointer, encStatus *SKP_SILK_SDK_EncControlStruct) int32 {
	var (
		psEnc *SKP_Silk_encoder_state_FIX
		ret   int32 = 0
	)
	psEnc = (*SKP_Silk_encoder_state_FIX)(encState)
	if func() int32 {
		ret += SKP_Silk_init_encoder_FIX(psEnc)
		return ret
	}() != 0 {
	}
	if func() int32 {
		ret += SKP_Silk_SDK_QueryEncoder(encState, encStatus)
		return ret
	}() != 0 {
	}
	return ret
}
func SKP_Silk_SDK_Encode(encState unsafe.Pointer, encControl *SKP_SILK_SDK_EncControlStruct, samplesIn []int16, nSamplesIn int32, outData []uint8, nBytesOut *int16) int32 {
	var (
		max_internal_fs_kHz int32
		PacketSize_ms       int32
		PacketLoss_perc     int32
		UseInBandFEC        int32
		UseDTX              int32
		ret                 int32 = 0
		nSamplesToBuffer    int32
		Complexity          int32
		input_10ms          int32
		nSamplesFromInput   int32 = 0
		TargetRate_bps      int32
		API_fs_Hz           int32
		MaxBytesOut         int16
		psEnc               = (*SKP_Silk_encoder_state_FIX)(encState)
	)
	if encControl.API_sampleRate != 8000 && encControl.API_sampleRate != 12000 && encControl.API_sampleRate != 16000 && encControl.API_sampleRate != 24000 && encControl.API_sampleRate != 32000 && encControl.API_sampleRate != 44100 && encControl.API_sampleRate != 48000 || encControl.MaxInternalSampleRate != 8000 && encControl.MaxInternalSampleRate != 12000 && encControl.MaxInternalSampleRate != 16000 && encControl.MaxInternalSampleRate != 24000 {
		ret = -2
		return ret
	}
	API_fs_Hz = encControl.API_sampleRate
	max_internal_fs_kHz = (encControl.MaxInternalSampleRate >> 10) + 1
	PacketSize_ms = (encControl.PacketSize * 1000) / API_fs_Hz
	TargetRate_bps = encControl.BitRate
	PacketLoss_perc = encControl.PacketLossPercentage
	UseInBandFEC = encControl.UseInBandFEC
	Complexity = encControl.Complexity
	UseDTX = encControl.UseDTX
	psEnc.SCmn.API_fs_Hz = API_fs_Hz
	psEnc.SCmn.MaxInternal_fs_kHz = max_internal_fs_kHz
	psEnc.SCmn.UseInBandFEC = UseInBandFEC
	input_10ms = (nSamplesIn * 100) / API_fs_Hz
	if input_10ms*API_fs_Hz != nSamplesIn*100 || nSamplesIn < 0 {
		ret = -1
		return ret
	}
	if MIN_TARGET_RATE_BPS > MAX_TARGET_RATE_BPS {
		if TargetRate_bps > MIN_TARGET_RATE_BPS {
			TargetRate_bps = MIN_TARGET_RATE_BPS
		} else if TargetRate_bps < MAX_TARGET_RATE_BPS {
			TargetRate_bps = MAX_TARGET_RATE_BPS
		} else {
			TargetRate_bps = TargetRate_bps
		}
	} else if TargetRate_bps > MAX_TARGET_RATE_BPS {
		TargetRate_bps = MAX_TARGET_RATE_BPS
	} else if TargetRate_bps < MIN_TARGET_RATE_BPS {
		TargetRate_bps = MIN_TARGET_RATE_BPS
	} else {
		TargetRate_bps = TargetRate_bps
	}
	if (func() int32 {
		ret = SKP_Silk_control_encoder_FIX(psEnc, PacketSize_ms, TargetRate_bps, PacketLoss_perc, UseDTX, Complexity)
		return ret
	}()) != 0 {
		return ret
	}
	if nSamplesIn*1000 > psEnc.SCmn.PacketSize_ms*API_fs_Hz {
		ret = -1
		return ret
	}
	if (func() int32 {
		if API_fs_Hz < (max_internal_fs_kHz * 1000) {
			return API_fs_Hz
		}
		return max_internal_fs_kHz * 1000
	}()) == 24000 && psEnc.SCmn.SSWBdetect.SWB_detected == 0 && psEnc.SCmn.SSWBdetect.WB_detected == 0 {
		SKP_Silk_detect_SWB_input(&psEnc.SCmn.SSWBdetect, samplesIn, nSamplesIn)
	}
	MaxBytesOut = 0
	for {
		nSamplesToBuffer = psEnc.SCmn.Frame_length - psEnc.SCmn.InputBufIx
		if API_fs_Hz == SKP_SMULBB(1000, psEnc.SCmn.Fs_kHz) {
			nSamplesToBuffer = SKP_min_int(nSamplesToBuffer, nSamplesIn)
			nSamplesFromInput = nSamplesToBuffer
			memcpy(unsafe.Pointer(&psEnc.SCmn.InputBuf[psEnc.SCmn.InputBufIx]), unsafe.Pointer(&samplesIn[0]), size_t(uintptr(nSamplesFromInput)*unsafe.Sizeof(int16(0))))
		} else {
			if nSamplesToBuffer < (input_10ms * 10 * psEnc.SCmn.Fs_kHz) {
				nSamplesToBuffer = nSamplesToBuffer
			} else {
				nSamplesToBuffer = input_10ms * 10 * psEnc.SCmn.Fs_kHz
			}
			nSamplesFromInput = (nSamplesToBuffer * API_fs_Hz) / (psEnc.SCmn.Fs_kHz * 1000)
			ret += SKP_Silk_resampler(&psEnc.SCmn.Resampler_state, ([]int16)(&psEnc.SCmn.InputBuf[psEnc.SCmn.InputBufIx]), samplesIn, nSamplesFromInput)
		}
		samplesIn += ([]int16)(nSamplesFromInput)
		nSamplesIn -= nSamplesFromInput
		psEnc.SCmn.InputBufIx += nSamplesToBuffer
		if psEnc.SCmn.InputBufIx >= psEnc.SCmn.Frame_length {
			if MaxBytesOut == 0 {
				MaxBytesOut = *nBytesOut
				ret = SKP_Silk_encode_frame_FIX(psEnc, &outData[0], &MaxBytesOut, &psEnc.SCmn.InputBuf[0])
				if ret != 0 {
				}
			} else {
				ret = SKP_Silk_encode_frame_FIX(psEnc, &outData[0], nBytesOut, &psEnc.SCmn.InputBuf[0])
				if ret != 0 {
				}
			}
			psEnc.SCmn.InputBufIx = 0
			psEnc.SCmn.Controlled_since_last_payload = 0
			if nSamplesIn == 0 {
				break
			}
		} else {
			break
		}
	}
	*nBytesOut = MaxBytesOut
	if psEnc.SCmn.UseDTX != 0 && psEnc.SCmn.InDTX != 0 {
		*nBytesOut = 0
	}
	return ret
}
