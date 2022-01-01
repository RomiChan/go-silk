package silk

import "unsafe"

func SKP_Silk_SDK_Get_Decoder_Size(decSizeBytes *int32) int32 {
	var ret int32 = 0
	*decSizeBytes = int32(uint32(unsafe.Sizeof(SKP_Silk_decoder_state{})))
	return ret
}
func SKP_Silk_SDK_InitDecoder(decState unsafe.Pointer) int32 {
	var (
		ret   int32 = 0
		struc *SKP_Silk_decoder_state
	)
	struc = (*SKP_Silk_decoder_state)(decState)
	ret = SKP_Silk_init_decoder(struc)
	return ret
}
func SKP_Silk_SDK_Decode(decState unsafe.Pointer, decControl *SKP_SILK_SDK_DecControlStruct, lostFlag int32, inData *uint8, nBytesIn int32, samplesOut *int16, nSamplesOut *int16) int32 {
	var (
		ret                 int32 = 0
		used_bytes          int32
		prev_fs_kHz         int32
		psDec               *SKP_Silk_decoder_state
		samplesOutInternal  [960]int16
		pSamplesOutInternal *int16
	)
	psDec = (*SKP_Silk_decoder_state)(decState)
	pSamplesOutInternal = samplesOut
	if psDec.Fs_kHz*1000 > decControl.API_sampleRate {
		pSamplesOutInternal = &samplesOutInternal[0]
	}
	if psDec.MoreInternalDecoderFrames == 0 {
		psDec.NFramesDecoded = 0
	}
	if psDec.MoreInternalDecoderFrames == 0 && lostFlag == 0 && nBytesIn > MAX_ARITHM_BYTES {
		lostFlag = 1
		ret = -11
	}
	prev_fs_kHz = psDec.Fs_kHz
	ret += SKP_Silk_decode_frame(psDec, ([]int16)(pSamplesOutInternal), nSamplesOut, ([]uint8)(inData), nBytesIn, lostFlag, &used_bytes)
	if used_bytes != 0 {
		if psDec.NBytesLeft > 0 && psDec.FrameTermination == SKP_SILK_MORE_FRAMES && psDec.NFramesDecoded < 5 {
			psDec.MoreInternalDecoderFrames = 1
		} else {
			psDec.MoreInternalDecoderFrames = 0
			psDec.NFramesInPacket = psDec.NFramesDecoded
			if psDec.VadFlag == VOICE_ACTIVITY {
				if psDec.FrameTermination == SKP_SILK_LAST_FRAME {
					psDec.No_FEC_counter++
					if psDec.No_FEC_counter > NO_LBRR_THRES {
						psDec.Inband_FEC_offset = 0
					}
				} else if psDec.FrameTermination == SKP_SILK_LBRR_VER1 {
					psDec.Inband_FEC_offset = 1
					psDec.No_FEC_counter = 0
				} else if psDec.FrameTermination == SKP_SILK_LBRR_VER2 {
					psDec.Inband_FEC_offset = 2
					psDec.No_FEC_counter = 0
				}
			}
		}
	}
	if MAX_API_FS_KHZ*1000 < decControl.API_sampleRate || 8000 > decControl.API_sampleRate {
		ret = -10
		return ret
	}
	if psDec.Fs_kHz*1000 != decControl.API_sampleRate {
		var samplesOut_tmp [960]int16
		memcpy(unsafe.Pointer(&samplesOut_tmp[0]), unsafe.Pointer(pSamplesOutInternal), size_t(uintptr(*nSamplesOut)*unsafe.Sizeof(int16(0))))
		if prev_fs_kHz != psDec.Fs_kHz || psDec.Prev_API_sampleRate != decControl.API_sampleRate {
			ret = SKP_Silk_resampler_init(&psDec.Resampler_state, SKP_SMULBB(psDec.Fs_kHz, 1000), decControl.API_sampleRate)
		}
		ret += SKP_Silk_resampler(&psDec.Resampler_state, ([]int16)(samplesOut), samplesOut_tmp[:], int32(*nSamplesOut))
		*nSamplesOut = int16((int32(*nSamplesOut) * decControl.API_sampleRate) / (psDec.Fs_kHz * 1000))
	} else if prev_fs_kHz*1000 > decControl.API_sampleRate {
		memcpy(unsafe.Pointer(samplesOut), unsafe.Pointer(pSamplesOutInternal), size_t(uintptr(*nSamplesOut)*unsafe.Sizeof(int16(0))))
	}
	psDec.Prev_API_sampleRate = decControl.API_sampleRate
	decControl.FrameSize = int32(uint16(int16(decControl.API_sampleRate / 50)))
	decControl.FramesPerPacket = psDec.NFramesInPacket
	decControl.InBandFECOffset = psDec.Inband_FEC_offset
	decControl.MoreInternalDecoderFrames = psDec.MoreInternalDecoderFrames
	return ret
}
func SKP_Silk_SDK_search_for_LBRR(inData *uint8, nBytesIn int32, lost_offset int32, LBRRData *uint8, nLBRRBytes *int16) {
	var (
		sDec     SKP_Silk_decoder_state
		sDecCtrl SKP_Silk_decoder_control
		TempQ    [480]int32
	)
	if lost_offset < 1 || lost_offset > MAX_LBRR_DELAY {
		*nLBRRBytes = 0
		return
	}
	sDec.NFramesDecoded = 0
	sDec.Fs_kHz = 0
	sDec.LossCnt = 0
	memset(unsafe.Pointer(&sDec.PrevNLSF_Q15[0]), 0, size_t(MAX_LPC_ORDER*unsafe.Sizeof(int32(0))))
	SKP_Silk_range_dec_init(&sDec.SRC, ([]uint8)(inData), nBytesIn)
	for {
		SKP_Silk_decode_parameters(&sDec, &sDecCtrl, TempQ[:], 0)
		if sDec.SRC.Error != 0 {
			*nLBRRBytes = 0
			return
		}
		if (sDec.FrameTermination-1)&lost_offset != 0 && sDec.FrameTermination > 0 && sDec.NBytesLeft >= 0 {
			*nLBRRBytes = int16(sDec.NBytesLeft)
			memcpy(unsafe.Pointer(LBRRData), unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(inData), nBytesIn-sDec.NBytesLeft))), size_t(uintptr(sDec.NBytesLeft)*unsafe.Sizeof(uint8(0))))
			break
		}
		if sDec.NBytesLeft > 0 && sDec.FrameTermination == SKP_SILK_MORE_FRAMES {
			sDec.NFramesDecoded++
		} else {
			LBRRData = 0
			*nLBRRBytes = 0
			break
		}
	}
}
func SKP_Silk_SDK_get_TOC(inData *uint8, nBytesIn int32, Silk_TOC *SKP_Silk_TOC_struct) {
	var (
		sDec     SKP_Silk_decoder_state
		sDecCtrl SKP_Silk_decoder_control
		TempQ    [480]int32
	)
	sDec.NFramesDecoded = 0
	sDec.Fs_kHz = 0
	SKP_Silk_range_dec_init(&sDec.SRC, ([]uint8)(inData), nBytesIn)
	Silk_TOC.Corrupt = 0
	for {
		SKP_Silk_decode_parameters(&sDec, &sDecCtrl, TempQ[:], 0)
		Silk_TOC.VadFlags[sDec.NFramesDecoded] = sDec.VadFlag
		Silk_TOC.SigtypeFlags[sDec.NFramesDecoded] = sDecCtrl.Sigtype
		if sDec.SRC.Error != 0 {
			Silk_TOC.Corrupt = 1
			break
		}
		if sDec.NBytesLeft > 0 && sDec.FrameTermination == SKP_SILK_MORE_FRAMES {
			sDec.NFramesDecoded++
		} else {
			break
		}
	}
	if Silk_TOC.Corrupt != 0 || sDec.FrameTermination == SKP_SILK_MORE_FRAMES || sDec.NFramesInPacket > SILK_MAX_FRAMES_PER_PACKET {
		memset(unsafe.Pointer(Silk_TOC), 0, size_t(unsafe.Sizeof(SKP_Silk_TOC_struct{})))
		Silk_TOC.Corrupt = 1
	} else {
		Silk_TOC.FramesInPacket = sDec.NFramesDecoded + 1
		Silk_TOC.Fs_kHz = sDec.Fs_kHz
		if sDec.FrameTermination == SKP_SILK_LAST_FRAME {
			Silk_TOC.InbandLBRR = sDec.FrameTermination
		} else {
			Silk_TOC.InbandLBRR = sDec.FrameTermination - 1
		}
	}
}
func SKP_Silk_SDK_get_version() *byte {
	var version [6]byte = func() [6]byte {
		var t [6]byte
		copy(t[:], ([]byte)("1.0.9"))
		return t
	}()
	return &version[0]
}
