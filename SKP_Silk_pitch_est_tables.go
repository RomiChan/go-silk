package silk

var SKP_Silk_CB_lags_stage2 [4][11]int16 = [4][11]int16{{0, 2, -1, -1, -1, 0, 0, 1, 1, 0, 1}, {0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0}, {0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0}, {0, -1, 2, 1, 0, 1, 1, 0, 0, -1, -1}}
var SKP_Silk_CB_lags_stage3 [4][34]int16 = [4][34]int16{{-9, -7, -6, -5, -5, -4, -4, -3, -3, -2, -2, -2, -1, -1, -1, 0, 0, 0, 1, 1, 0, 1, 2, 2, 2, 3, 3, 4, 4, 5, 6, 5, 6, 8}, {-3, -2, -2, -2, -1, -1, -1, -1, -1, 0, 0, -1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 1, 2, 1, 2, 2, 2, 2, 3}, {3, 3, 2, 2, 2, 2, 1, 2, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, -1, 0, 0, -1, -1, -1, -1, -1, -2, -2, -2}, {9, 8, 6, 5, 6, 5, 4, 4, 3, 3, 2, 2, 2, 1, 0, 1, 1, 0, 0, 0, -1, -1, -1, -2, -2, -2, -3, -3, -4, -4, -5, -5, -6, -7}}
var SKP_Silk_Lag_range_stage3 [3][4][2]int16 = [3][4][2]int16{{{-2, 6}, {-1, 5}, {-1, 5}, {-2, 7}}, {{-4, 8}, {-1, 6}, {-1, 6}, {-4, 9}}, {{-9, 12}, {-3, 7}, {-2, 7}, {-7, 13}}}
var SKP_Silk_cbk_sizes_stage3 [3]int16 = [3]int16{PITCH_EST_NB_CBKS_STAGE3_MIN, PITCH_EST_NB_CBKS_STAGE3_MID, PITCH_EST_NB_CBKS_STAGE3_MAX}
var SKP_Silk_cbk_offsets_stage3 [3]int16 = [3]int16{((PITCH_EST_NB_CBKS_STAGE3_MAX - PITCH_EST_NB_CBKS_STAGE3_MIN) >> 1), ((PITCH_EST_NB_CBKS_STAGE3_MAX - PITCH_EST_NB_CBKS_STAGE3_MID) >> 1), 0}
