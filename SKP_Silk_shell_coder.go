package silk

import "unsafe"

func combine_pulses(out []int32, in []int32, len_ int32) {
	var k int32
	for k = 0; k < len_; k++ {
		out[k] = in[k*2] + in[k*2+1]
	}
}
func encode_split(sRC *SKP_Silk_range_coder_state, p_child1 int32, p int32, shell_table []uint16) {
	var cdf []uint16
	if p > 0 {
		cdf = ([]uint16)(&shell_table[SKP_Silk_shell_code_table_offsets[p]])
		SKP_Silk_range_encoder(sRC, p_child1, cdf)
	}
}
func decode_split(p_child1 []int32, p_child2 []int32, sRC *SKP_Silk_range_coder_state, p int32, shell_table *uint16) {
	var (
		cdf_middle int32
		cdf        []uint16
	)
	if p > 0 {
		cdf_middle = p >> 1
		cdf = ([]uint16)((*uint16)(unsafe.Add(unsafe.Pointer(shell_table), unsafe.Sizeof(uint16(0))*uintptr(SKP_Silk_shell_code_table_offsets[p]))))
		SKP_Silk_range_decoder(&p_child1[0], sRC, cdf, cdf_middle)
		p_child2[0] = p - p_child1[0]
	} else {
		p_child1[0] = 0
		p_child2[0] = 0
	}
}
func SKP_Silk_shell_encoder(sRC *SKP_Silk_range_coder_state, pulses0 []int32) {
	var (
		pulses1 [8]int32
		pulses2 [4]int32
		pulses3 [2]int32
		pulses4 [1]int32
	)
	combine_pulses(pulses1[:], pulses0, 8)
	combine_pulses(pulses2[:], pulses1[:], 4)
	combine_pulses(pulses3[:], pulses2[:], 2)
	combine_pulses(pulses4[:], pulses3[:], 1)
	encode_split(sRC, pulses3[0], pulses4[0], SKP_Silk_shell_code_table3[:])
	encode_split(sRC, pulses2[0], pulses3[0], SKP_Silk_shell_code_table2[:])
	encode_split(sRC, pulses1[0], pulses2[0], SKP_Silk_shell_code_table1[:])
	encode_split(sRC, pulses0[0], pulses1[0], SKP_Silk_shell_code_table0[:])
	encode_split(sRC, pulses0[2], pulses1[1], SKP_Silk_shell_code_table0[:])
	encode_split(sRC, pulses1[2], pulses2[1], SKP_Silk_shell_code_table1[:])
	encode_split(sRC, pulses0[4], pulses1[2], SKP_Silk_shell_code_table0[:])
	encode_split(sRC, pulses0[6], pulses1[3], SKP_Silk_shell_code_table0[:])
	encode_split(sRC, pulses2[2], pulses3[1], SKP_Silk_shell_code_table2[:])
	encode_split(sRC, pulses1[4], pulses2[2], SKP_Silk_shell_code_table1[:])
	encode_split(sRC, pulses0[8], pulses1[4], SKP_Silk_shell_code_table0[:])
	encode_split(sRC, pulses0[10], pulses1[5], SKP_Silk_shell_code_table0[:])
	encode_split(sRC, pulses1[6], pulses2[3], SKP_Silk_shell_code_table1[:])
	encode_split(sRC, pulses0[12], pulses1[6], SKP_Silk_shell_code_table0[:])
	encode_split(sRC, pulses0[14], pulses1[7], SKP_Silk_shell_code_table0[:])
}
func SKP_Silk_shell_decoder(pulses0 []int32, sRC *SKP_Silk_range_coder_state, pulses4 int32) {
	var (
		pulses3 [2]int32
		pulses2 [4]int32
		pulses1 [8]int32
	)
	decode_split(pulses3[:], pulses3[1:], sRC, pulses4, &SKP_Silk_shell_code_table3[0])
	decode_split(pulses2[:], pulses2[1:], sRC, pulses3[0], &SKP_Silk_shell_code_table2[0])
	decode_split(pulses1[:], pulses1[1:], sRC, pulses2[0], &SKP_Silk_shell_code_table1[0])
	decode_split(pulses0, pulses0[1:], sRC, pulses1[0], &SKP_Silk_shell_code_table0[0])
	decode_split(pulses0[2:], pulses0[3:], sRC, pulses1[1], &SKP_Silk_shell_code_table0[0])
	decode_split(pulses1[2:], pulses1[3:], sRC, pulses2[1], &SKP_Silk_shell_code_table1[0])
	decode_split(pulses0[4:], pulses0[5:], sRC, pulses1[2], &SKP_Silk_shell_code_table0[0])
	decode_split(pulses0[6:], pulses0[7:], sRC, pulses1[3], &SKP_Silk_shell_code_table0[0])
	decode_split(pulses2[2:], pulses2[3:], sRC, pulses3[1], &SKP_Silk_shell_code_table2[0])
	decode_split(pulses1[4:], pulses1[5:], sRC, pulses2[2], &SKP_Silk_shell_code_table1[0])
	decode_split(pulses0[8:], pulses0[9:], sRC, pulses1[4], &SKP_Silk_shell_code_table0[0])
	decode_split(pulses0[10:], pulses0[11:], sRC, pulses1[5], &SKP_Silk_shell_code_table0[0])
	decode_split(pulses1[6:], pulses1[7:], sRC, pulses2[3], &SKP_Silk_shell_code_table1[0])
	decode_split(pulses0[12:], pulses0[13:], sRC, pulses1[6], &SKP_Silk_shell_code_table0[0])
	decode_split(pulses0[14:], pulses0[15:], sRC, pulses1[7], &SKP_Silk_shell_code_table0[0])
}
