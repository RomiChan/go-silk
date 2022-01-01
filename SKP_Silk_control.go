package silk

type SKP_SILK_SDK_EncControlStruct struct {
	API_sampleRate        int32
	MaxInternalSampleRate int32
	PacketSize            int32
	BitRate               int32
	PacketLossPercentage  int32
	Complexity            int32
	UseInBandFEC          int32
	UseDTX                int32
}
type SKP_SILK_SDK_DecControlStruct struct {
	API_sampleRate            int32
	FrameSize                 int32
	FramesPerPacket           int32
	MoreInternalDecoderFrames int32
	InBandFECOffset           int32
}
