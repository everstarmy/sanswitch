package san

import "encoding/xml"

// FibreChannelStatisticsResponse 是 GET /brocade-interface/fibrechannel-statistics 的 XML 响应包装
type FibreChannelStatisticsResponse struct {
	XMLName    xml.Name                     `xml:"Response"`
	Statistics []FibreChannelStatisticsInfo `xml:"fibrechannel-statistics"`
}

// FibreChannelStatisticsInfo 描述一个 FC 端口的性能统计计数器，
// 包含流量、错误、链路状态、缓冲区、FEC 等 120+ 项指标
type FibreChannelStatisticsInfo struct {
	XMLName                              xml.Name `xml:"fibrechannel-statistics" json:"-"`
	Name                                 string   `xml:"name" json:"name"`
	SamplingInterval                     uint16   `xml:"sampling-interval" json:"sampling_interval"`
	TimeGenerated                        string   `xml:"time-generated" json:"time_generated"`
	TimeRefreshed                        string   `xml:"time-refreshed" json:"time_refreshed"`
	ResetStatistics                      uint8    `xml:"reset-statistics" json:"reset_statistics"`
	InOctets                             uint64   `xml:"in-octets" json:"in_octets"`
	OutOctets                            uint64   `xml:"out-octets" json:"out_octets"`
	InMulticastPkts                      uint64   `xml:"in-multicast-pkts" json:"in_multicast_pkts"`
	OutMulticastPkts                     uint64   `xml:"out-multicast-pkts" json:"out_multicast_pkts"`
	InLinkResets                         uint64   `xml:"in-link-resets" json:"in_link_resets"`
	OutLinkResets                        uint64   `xml:"out-link-resets" json:"out_link_resets"`
	InOfflineSequences                   uint64   `xml:"in-offline-sequences" json:"in_offline_sequences"`
	OutOfflineSequences                  uint64   `xml:"out-offline-sequences" json:"out_offline_sequences"`
	InvalidOrderedSets                   uint64   `xml:"invalid-ordered-sets" json:"invalid_ordered_sets"`
	FramesTooLong                        uint64   `xml:"frames-too-long" json:"frames_too_long"`
	TruncatedFrames                      uint64   `xml:"truncated-frames" json:"truncated_frames"`
	AddressErrors                        uint64   `xml:"address-errors" json:"address_errors"`
	DelimiterErrors                      uint64   `xml:"delimiter-errors" json:"delimiter_errors"`
	EncodingDisparityErrors              uint64   `xml:"encoding-disparity-errors" json:"encoding_disparity_errors"`
	TooManyRdys                          uint64   `xml:"too-many-rdys" json:"too_many_rdys"`
	InCRCErrors                          uint64   `xml:"in-crc-errors" json:"in_crc_errors"`
	CRCErrors                            uint64   `xml:"crc-errors" json:"crc_errors"`
	BadEOFsReceived                      uint64   `xml:"bad-eofs-received" json:"bad_eofs_received"`
	EncodingErrorsOutsideFrame           uint64   `xml:"encoding-errors-outside-frame" json:"encoding_errors_outside_frame"`
	MulticastTimeouts                    uint64   `xml:"multicast-timeouts" json:"multicast_timeouts"`
	InLCs                                uint64   `xml:"in-lcs" json:"in_lcs"`
	InFrameRate                          uint64   `xml:"in-frame-rate" json:"in_frame_rate"`
	OutFrameRate                         uint64   `xml:"out-frame-rate" json:"out_frame_rate"`
	InMaxFrameRate                       uint64   `xml:"in-max-frame-rate" json:"in_max_frame_rate"`
	OutMaxFrameRate                      uint64   `xml:"out-max-frame-rate" json:"out_max_frame_rate"`
	InRate                               uint64   `xml:"in-rate" json:"in_rate"`
	OutRate                              uint64   `xml:"out-rate" json:"out_rate"`
	InPeakRate                           uint64   `xml:"in-peak-rate" json:"in_peak_rate"`
	OutPeakRate                          uint64   `xml:"out-peak-rate" json:"out_peak_rate"`
	InFrames                             uint64   `xml:"in-frames" json:"in_frames"`
	OutFrames                            uint64   `xml:"out-frames" json:"out_frames"`
	BB_CreditZero                        uint64   `xml:"bb-credit-zero" json:"bb_credit_zero"`
	InputBufferFull                      uint64   `xml:"input-buffer-full" json:"input_buffer_full"`
	FBusyFrames                          uint64   `xml:"f-busy-frames" json:"f_busy_frames"`
	PBusyFrames                          uint64   `xml:"p-busy-frames" json:"p_busy_frames"`
	FRJTFrames                           uint64   `xml:"f-rjt-frames" json:"f_rjt_frames"`
	PRJTFrames                           uint64   `xml:"p-rjt-frames" json:"p_rjt_frames"`
	Class1Frames                         uint64   `xml:"class-1-frames" json:"class_1_frames"`
	Class2Frames                         uint64   `xml:"class-2-frames" json:"class_2_frames"`
	Class3Frames                         uint64   `xml:"class-3-frames" json:"class_3_frames"`
	Class3Discards                       uint64   `xml:"class-3-discards" json:"class_3_discards"`
	LinkFailures                         uint64   `xml:"link-failures" json:"link_failures"`
	InvalidTransmissionWords             uint64   `xml:"invalid-transmission-words" json:"invalid_transmission_words"`
	PrimitiveSequenceProtocolError       uint64   `xml:"primitive-sequence-protocol-error" json:"primitive_sequence_protocol_error"`
	LossOfSignal                         uint64   `xml:"loss-of-signal" json:"loss_of_signal"`
	LossOfSync                           uint64   `xml:"loss-of-sync" json:"loss_of_sync"`
	Class3InDiscards                     uint64   `xml:"class3-in-discards" json:"class3_in_discards"`
	Class3OutDiscards                    uint64   `xml:"class3-out-discards" json:"class3_out_discards"`
	PCSBlockErrors                       uint64   `xml:"pcs-block-errors" json:"pcs_block_errors"`
	RemoteLinkFailures                   uint64   `xml:"remote-link-failures" json:"remote_link_failures"`
	RemoteInvalidTransmissionWords       uint64   `xml:"remote-invalid-transmission-words" json:"remote_invalid_transmission_words"`
	RemotePrimitiveSequenceProtocolError uint64   `xml:"remote-primitive-sequence-protocol-error" json:"remote_primitive_sequence_protocol_error"`
	RemoteLossOfSignal                   uint64   `xml:"remote-loss-of-signal" json:"remote_loss_of_signal"`
	RemoteLossOfSync                     uint64   `xml:"remote-loss-of-sync" json:"remote_loss_of_sync"`
	RemoteCRCErrors                      uint64   `xml:"remote-crc-errors" json:"remote_crc_errors"`
	RemoteFECUncorrected                 uint64   `xml:"remote-fec-uncorrected" json:"remote_fec_uncorrected"`
	LinkLevelInterrupts                  uint64   `xml:"link-level-interrupts" json:"link_level_interrupts"`
	FramesProcessingRequired             uint64   `xml:"frames-processing-required" json:"frames_processing_required"`
	FramesTimedOut                       uint64   `xml:"frames-timed-out" json:"frames_timed_out"`
	FramesTransmitterUnavailableErrors   uint64   `xml:"frames-transmitter-unavailable-errors" json:"frames_transmitter_unavailable_errors"`
	NonOperationalSequencesIn            uint64   `xml:"non-operational-sequences-in" json:"non_operational_sequences_in"`
	NonOperationalSequencesOut           uint64   `xml:"non-operational-sequences-out" json:"non_operational_sequences_out"`
	FECUncorrected                       uint64   `xml:"fec-uncorrected" json:"fec_uncorrected"`
	TotalUpTime                          float64  `xml:"total-up-time" json:"total_up_time"`
	TotalDownTime                        float64  `xml:"total-down-time" json:"total_down_time"`
	DownOccurrence                       uint32   `xml:"down-occurrence" json:"down_occurrence"`
	TotalOfflineTime                     float64  `xml:"total-offline-time" json:"total_offline_time"`
	EncodingErrorIn                      uint64   `xml:"encoding-error-in" json:"encoding_error_in"`
	TxPeakFrame                          uint64   `xml:"tx-peak-frame" json:"tx_peak_frame"`
	RxPeakFrame                          uint64   `xml:"rx-peak-frame" json:"rx_peak_frame"`
	StarvationTenancyStop                uint64   `xml:"starvation-tenancy-stop" json:"starvation_tenancy_stop"`
	InvalidArbitrateLoop                 uint64   `xml:"invalid-arbitrate-loop" json:"invalid_arbitrate_loop"`
	FTBType1Miss                         uint64   `xml:"ftb-type-1-miss" json:"ftb_type_1_miss"`
	FTBType2Miss                         uint64   `xml:"ftb-type-2-miss" json:"ftb_type_2_miss"`
	FTBType6Miss                         uint64   `xml:"ftb-type-6-miss" json:"ftb_type_6_miss"`
	HardZoneMiss                         uint64   `xml:"hard-zone-miss" json:"hard_zone_miss"`
	LUNZoneMiss                          uint64   `xml:"lun-zone-miss" json:"lun_zone_miss"`
	StateTransitionCount                 uint32   `xml:"state-transition-count" json:"state_transition_count"`
	InterruptCount                       uint32   `xml:"interrupt-count" json:"interrupt_count"`
	UnknownInterruptCount                uint32   `xml:"unknown-interrupt-count" json:"unknown_interrupt_count"`
	CongestionPrimitiveIn                uint32   `xml:"congestion-primitive-in" json:"congestion_primitive_in"`
	OutTotalPkts                         uint64   `xml:"out-total-pkts" json:"out_total_pkts"`
	OutTotalOctets                       uint64   `xml:"out-total-octets" json:"out_total_octets"`
	OutPausePkts                         uint64   `xml:"out-pause-pkts" json:"out_pause_pkts"`
	Out64bPkts                           uint64   `xml:"out-64b-pkts" json:"out_64b_pkts"`
	Out65b127bPkts                       uint64   `xml:"out-65b-127b-pkts" json:"out_65b_127b_pkts"`
	Out128b255bPkts                      uint64   `xml:"out-128b-255b-pkts" json:"out_128b_255b_pkts"`
	Out256b511bPkts                      uint64   `xml:"out-256b-511b-pkts" json:"out_256b_511b_pkts"`
	Out512b1023bPkts                     uint64   `xml:"out-512b-1023b-pkts" json:"out_512b_1023b_pkts"`
	Out1024b1518bPkts                    uint64   `xml:"out-1024b-1518b-pkts" json:"out_1024b_1518b_pkts"`
	OutLargePkts                         uint64   `xml:"out-large-pkts" json:"out_large_pkts"`
	InTotalPkts                          uint64   `xml:"in-total-pkts" json:"in_total_pkts"`
	InTotalOctets                        uint64   `xml:"in-total-octets" json:"in_total_octets"`
	InPausePkts                          uint64   `xml:"in-pause-pkts" json:"in_pause_pkts"`
	InGoodPkts                           uint64   `xml:"in-good-pkts" json:"in_good_pkts"`
	In64bPkts                            uint64   `xml:"in-64b-pkts" json:"in_64b_pkts"`
	In65b127bPkts                        uint64   `xml:"in-65b-127b-pkts" json:"in_65b_127b_pkts"`
	In128b255bPkts                       uint64   `xml:"in-128b-255b-pkts" json:"in_128b_255b_pkts"`
	In256b511bPkts                       uint64   `xml:"in-256b-511b-pkts" json:"in_256b_511b_pkts"`
	In512b1023bPkts                      uint64   `xml:"in-512b-1023b-pkts" json:"in_512b_1023b_pkts"`
	In1024b1518bPkts                     uint64   `xml:"in-1024b-1518b-pkts" json:"in_1024b_1518b_pkts"`
	InLargePkts                          uint64   `xml:"in-large-pkts" json:"in_large_pkts"`
	InOversizedPkts                      uint64   `xml:"in-oversized-pkts" json:"in_oversized_pkts"`
	InRuntPkts                           uint64   `xml:"in-runt-pkts" json:"in_runt_pkts"`
	InRuntBadCRC                         uint64   `xml:"in-runt-bad-crc" json:"in_runt_bad_crc"`
	InBadTermination                     uint64   `xml:"in-bad-termination" json:"in_bad_termination"`
	InCRCAlignmentError                  uint64   `xml:"in-crc-alignment-error" json:"in_crc_alignment_error"`
	InCRCStomp                           uint64   `xml:"in-crc-stomp" json:"in_crc_stomp"`
	InSymbolError                        uint64   `xml:"in-symbol-error" json:"in_symbol_error"`
	InIFGViolation                       uint64   `xml:"in-ifg-violation" json:"in_ifg_violation"`
	InEthernetDiscards                   uint32   `xml:"in-ethernet-discards" json:"in_ethernet_discards"`
	MulticastBroadcastDiscardPkts        uint32   `xml:"multicast-broadcast-discard-pkts" json:"multicast_broadcast_discard_pkts"`
	IPTTLDiscards                        uint32   `xml:"ip-ttl-discards" json:"ip_ttl_discards"`
	RemoteBBCredit                       uint64   `xml:"remote-buffer-credit-info>bb-credit" json:"remote_bb_credit"`
	RemotePeerBBCredit                   uint64   `xml:"remote-buffer-credit-info>peer-bb-credit" json:"remote_peer_bb_credit"`
}

// GetFibreChannelStatistics 获取所有 FC 端口的性能统计计数器
// 对应 API: GET /brocade-interface/fibrechannel-statistics
func (c *Client) GetFibreChannelStatistics() ([]FibreChannelStatisticsInfo, error) {
	var resp FibreChannelStatisticsResponse
	err := c.Get(c.endpoints().FibreChannelStatistics(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Statistics, nil
}
