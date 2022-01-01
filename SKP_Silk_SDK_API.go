package silk

const SILK_MAX_FRAMES_PER_PACKET = 5

type SKP_Silk_TOC_struct struct {
	FramesInPacket int32
	Fs_kHz         int32
	InbandLBRR     int32
	Corrupt        int32
	VadFlags       [5]int32
	SigtypeFlags   [5]int32
}
