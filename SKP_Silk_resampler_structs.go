package silk

import "unsafe"

const (
	RESAMPLER_SUPPORT_ABOVE_48KHZ    = 1
	SKP_Silk_RESAMPLER_MAX_FIR_ORDER = 16
	SKP_Silk_RESAMPLER_MAX_IIR_ORDER = 6
)

type _SKP_Silk_resampler_state_struct struct {
	SIIR               [6]int32
	SFIR               [16]int32
	SDown2             [2]int32
	Resampler_function func(unsafe.Pointer, *int16, *int16, int32)
	Up2_function       func(*int32, *int16, *int16, int32)
	BatchSize          int32
	InvRatio_Q16       int32
	FIR_Fracs          int32
	Input2x            int32
	Coefs              *int16
	SDownPre           [2]int32
	SUpPost            [2]int32
	Down_pre_function  func(*int32, *int16, *int16, int32)
	Up_post_function   func(*int32, *int16, *int16, int32)
	BatchSizePrePost   int32
	Ratio_Q16          int32
	NPreDownsamplers   int32
	NPostUpsamplers    int32
	Magic_number       int32
}
type SKP_Silk_resampler_state_struct _SKP_Silk_resampler_state_struct
