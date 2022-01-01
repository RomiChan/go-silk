package silk

import "math"

func SKP_SMULWB(a32 int32, b32 int32) int32 {
	return int32(((int64(a32) >> 16) * int64(int16(b32))) + (((int64(a32) & math.MaxUint16) * int64(int16(b32))) >> 16))
}
func SKP_SMLAWB(a32 int32, b32 int32, c32 int32) int32 {
	return int32(int64(a32) + (((int64(b32) >> 16) * int64(int16(c32))) + (((int64(b32) & math.MaxUint16) * int64(int16(c32))) >> 16)))
}
func SKP_SMULWT(a32 int32, b32 int32) int32 {
	return int32((int64(a32)>>16)*(int64(b32)>>16) + (((int64(a32) & math.MaxUint16) * (int64(b32) >> 16)) >> 16))
}
func SKP_SMLAWT(a32 int32, b32 int32, c32 int32) int32 {
	return int32(int64(a32) + (int64(b32)>>16)*(int64(c32)>>16) + (((int64(b32) & math.MaxUint16) * (int64(c32) >> 16)) >> 16))
}
func SKP_SMULBB(a32 int32, b32 int32) int32 {
	return int32(int64(int16(a32)) * int64(int16(b32)))
}
func SKP_SMLABB(a32 int32, b32 int32, c32 int32) int32 {
	return int32(int64(a32) + int64(int32(int16(b32)))*int64(int16(c32)))
}
func SKP_SMULBT(a32 int32, b32 int32) int32 {
	return int32(int64(int16(a32)) * (int64(b32) >> 16))
}
func SKP_SMLABT(a32 int32, b32 int32, c32 int32) int32 {
	return int32(int64(a32) + int64(int32(int16(b32)))*(int64(c32)>>16))
}
func SKP_SMLAL(a64 int64, b32 int32, c32 int32) int64 {
	return a64 + int64(b32)*int64(c32)
}
func SKP_SMULWW(a32 int32, b32 int32) int32 {
	return int32(int64(SKP_SMULWB(a32, b32)) + int64(a32)*int64(SKP_RSHIFT_ROUND(b32, 16)))
}
func SKP_SMLAWW(a32 int32, b32 int32, c32 int32) int32 {
	return int32(int64(SKP_SMLAWB(a32, b32, c32)) + int64(b32)*int64(SKP_RSHIFT_ROUND(c32, 16)))
}
func SKP_SMMUL(a32 int32, b32 int32) int32 {
	return int32((int64(a32) * int64(b32)) >> 32)
}
