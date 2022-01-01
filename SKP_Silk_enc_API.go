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
	encStatus.PacketSize = int32((int64(psEnc.SCmn.API_fs_Hz) * int64(psEnc.SCmn.PacketSize_ms)) / 1000)
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
	if int64(func() int32 {
		ret += SKP_Silk_init_encoder_FIX(psEnc)
		return ret
	}()) != 0 {
	}
	if int64(func() int32 {
		ret += SKP_Silk_SDK_QueryEncoder(encState, encStatus)
		return ret
	}()) != 0 {
	}
	return ret
}
func SKP_Silk_SDK_Encode(encState unsafe.Pointer, encControl *SKP_SILK_SDK_EncControlStruct, samplesIn []int16, nSamplesIn int32, outData *uint8, nBytesOut []int16) int32 {
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
		psEnc               *SKP_Silk_encoder_state_FIX = (*SKP_Silk_encoder_state_FIX)(encState)
	)
	if int64(encControl.API_sampleRate) != 8000 && int64(encControl.API_sampleRate) != 12000 && int64(encControl.API_sampleRate) != 16000 && int64(encControl.API_sampleRate) != 24000 && int64(encControl.API_sampleRate) != 32000 && int64(encControl.API_sampleRate) != 44100 && int64(encControl.API_sampleRate) != 48000 || int64(encControl.MaxInternalSampleRate) != 8000 && int64(encControl.MaxInternalSampleRate) != 12000 && int64(encControl.MaxInternalSampleRate) != 16000 && int64(encControl.MaxInternalSampleRate) != 24000 {
		ret = -2
		return ret
	}
	API_fs_Hz = encControl.API_sampleRate
	max_internal_fs_kHz = int32(int64(int32(int64(encControl.MaxInternalSampleRate)>>10)) + 1)
	PacketSize_ms = int32((int64(encControl.PacketSize) * 1000) / int64(API_fs_Hz))
	TargetRate_bps = encControl.BitRate
	PacketLoss_perc = encControl.PacketLossPercentage
	UseInBandFEC = encControl.UseInBandFEC
	Complexity = encControl.Complexity
	UseDTX = encControl.UseDTX
	psEnc.SCmn.API_fs_Hz = API_fs_Hz
	psEnc.SCmn.MaxInternal_fs_kHz = max_internal_fs_kHz
	psEnc.SCmn.UseInBandFEC = UseInBandFEC
	input_10ms = int32((int64(nSamplesIn) * 100) / int64(API_fs_Hz))
	if int64(input_10ms)*int64(API_fs_Hz) != int64(nSamplesIn)*100 || int64(nSamplesIn) < 0 {
		ret = -1
		return ret
	}
	if MIN_TARGET_RATE_BPS > MAX_TARGET_RATE_BPS {
		if int64(TargetRate_bps) > MIN_TARGET_RATE_BPS {
			TargetRate_bps = MIN_TARGET_RATE_BPS
		} else if int64(TargetRate_bps) < MAX_TARGET_RATE_BPS {
			TargetRate_bps = MAX_TARGET_RATE_BPS
		} else {
			TargetRate_bps = TargetRate_bps
		}
	} else if int64(TargetRate_bps) > MAX_TARGET_RATE_BPS {
		TargetRate_bps = MAX_TARGET_RATE_BPS
	} else if int64(TargetRate_bps) < MIN_TARGET_RATE_BPS {
		TargetRate_bps = MIN_TARGET_RATE_BPS
	} else {
		TargetRate_bps = TargetRate_bps
	}
	if int64(func() int32 {
		ret = SKP_Silk_control_encoder_FIX(psEnc, PacketSize_ms, TargetRate_bps, PacketLoss_perc, UseDTX, Complexity)
		return ret
	}()) != 0 {
		return ret
	}
	if int64(nSamplesIn)*1000 > int64(psEnc.SCmn.PacketSize_ms)*int64(API_fs_Hz) {
		ret = -1
		return ret
	}
	if (func() int64 {
		if int64(API_fs_Hz) < (int64(max_internal_fs_kHz) * 1000) {
			return int64(API_fs_Hz)
		}
		return int64(max_internal_fs_kHz) * 1000
	}()) == 24000 && int64(psEnc.SCmn.SSWBdetect.SWB_detected) == 0 && int64(psEnc.SCmn.SSWBdetect.WB_detected) == 0 {
		SKP_Silk_detect_SWB_input(&psEnc.SCmn.SSWBdetect, ([]int16)(samplesIn), nSamplesIn)
	}
	MaxBytesOut = 0
	for {
		nSamplesToBuffer = int32(int64(psEnc.SCmn.Frame_length) - int64(psEnc.SCmn.InputBufIx))
		if int64(API_fs_Hz) == int64(SKP_SMULBB(1000, psEnc.SCmn.Fs_kHz)) {
			nSamplesToBuffer = SKP_min_int(nSamplesToBuffer, nSamplesIn)
			nSamplesFromInput = nSamplesToBuffer
			memcpy(unsafe.Pointer(&psEnc.SCmn.InputBuf[psEnc.SCmn.InputBufIx]), unsafe.Pointer(samplesIn), size_t(uintptr(nSamplesFromInput)*unsafe.Sizeof(int16(0))))
		} else {
			if int64(nSamplesToBuffer) < (int64(input_10ms) * 10 * int64(psEnc.SCmn.Fs_kHz)) {
				nSamplesToBuffer = nSamplesToBuffer
			} else {
				nSamplesToBuffer = int32(int64(input_10ms) * 10 * int64(psEnc.SCmn.Fs_kHz))
			}
			nSamplesFromInput = int32((int64(nSamplesToBuffer) * int64(API_fs_Hz)) / (int64(psEnc.SCmn.Fs_kHz) * 1000))
			ret += SKP_Silk_resampler(&psEnc.SCmn.Resampler_state, ([]int16)(&psEnc.SCmn.InputBuf[psEnc.SCmn.InputBufIx]), ([]int16)(samplesIn), nSamplesFromInput)
		}
		samplesIn = (*int16)(unsafe.Add(unsafe.Pointer(samplesIn), unsafe.Sizeof(int16(0))*uintptr(nSamplesFromInput)))
		nSamplesIn -= nSamplesFromInput
		psEnc.SCmn.InputBufIx += nSamplesToBuffer
		if int64(psEnc.SCmn.InputBufIx) >= int64(psEnc.SCmn.Frame_length) {
			if int64(MaxBytesOut) == 0 {
				MaxBytesOut = *nBytesOut
				ret = SKP_Silk_encode_frame_FIX(psEnc, outData, &MaxBytesOut, &psEnc.SCmn.InputBuf[0])
				if int64(ret) != 0 {
				}
			} else {
				ret = SKP_Silk_encode_frame_FIX(psEnc, outData, nBytesOut, &psEnc.SCmn.InputBuf[0])
				if int64(ret) != 0 {
				}
			}
			psEnc.SCmn.InputBufIx = 0
			psEnc.SCmn.Controlled_since_last_payload = 0
			if int64(nSamplesIn) == 0 {
				break
			}
		} else {
			break
		}
	}
	*nBytesOut = MaxBytesOut
	if int64(psEnc.SCmn.UseDTX) != 0 && int64(psEnc.SCmn.InDTX) != 0 {
		*nBytesOut = 0
	}
	return ret
}
