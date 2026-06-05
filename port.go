package san

import "encoding/xml"

// FibreChannelResponse 是 GET /brocade-interface/fibrechannel 的 XML 响应包装
type FibreChannelResponse struct {
	XMLName xml.Name   `xml:"Response"`
	Ports   []PortInfo `xml:"fibrechannel"`
}

// PortInfo 描述一个 FC 端口的完整状态，包含运行状态、速率、WWN、邻居、缓冲区、协议能力等字段
type PortInfo struct {
	XMLName                             xml.Name `xml:"fibrechannel" json:"-"`
	Name                                string   `xml:"name" json:"name"`
	WWN                                 string   `xml:"wwn" json:"wwn"`
	OperationalStatus                   uint32   `xml:"operational-status" json:"operational_status"`
	OperationalStatusString             string   `xml:"operational-status-string" json:"operational_status_string"`
	EnabledState                        bool     `xml:"is-enabled-state" json:"enabled_state"`
	UserFriendlyName                    string   `xml:"user-friendly-name" json:"user_friendly_name"`
	Speed                               string   `xml:"protocol-speed" json:"speed"`
	MaxSpeed                            string   `xml:"max-protocol-speed" json:"max_speed"`
	AutoNegotiate                       bool     `xml:"auto-negotiate-v2" json:"auto_negotiate"`
	PortType                            string   `xml:"port-type-string" json:"port_type"`
	FCID                                string   `xml:"fcid" json:"fcid"`
	FCIDHex                             string   `xml:"fcid-hex" json:"fcid_hex"`
	NPIVEnabled                         bool     `xml:"npiv-enabled-v2" json:"npiv_enabled"`
	NPIVPPLimit                         int16    `xml:"npiv-pp-limit" json:"npiv_pp_limit"`
	NPIVFLOGILogoutEnabled              bool     `xml:"npiv-flogi-logout-enabled-v2" json:"npiv_flogi_logout_enabled"`
	LongDistance                        string   `xml:"long-distance-string" json:"long_distance"`
	VCLinkInit                          bool     `xml:"vc-link-init-v2" json:"vc_link_init"`
	ISLReadyModeEnabled                 bool     `xml:"isl-ready-mode-enabled-v2" json:"isl_ready_mode_enabled"`
	TrunkPortEnabled                    bool     `xml:"trunk-port-enabled-v2" json:"trunk_port_enabled"`
	DiagnosticsStatus                   string   `xml:"diagnostics-status" json:"diagnostics_status"`
	LockedEPortEnabled                  bool     `xml:"locked-e-port-enabled" json:"locked_e_port_enabled"`
	BladePortNumber                     uint32   `xml:"blade-port-number" json:"blade_port_number"`
	NeighborWWNs                        []string `xml:"neighbor>wwn" json:"neighbor_wwns"`
	NeighborSlotPort                    string   `xml:"neighbor-slot-port" json:"neighbor_slot_port"`
	NeighborNodeWWN                     string   `xml:"neighbor-node-wwn" json:"neighbor_node_wwn"`
	NeighborPortIndex                   uint32   `xml:"neighbor-port-index" json:"neighbor_port_index"`
	NeighborSwitchName                  string   `xml:"neighbor-switch-user-friendly-name" json:"neighbor_switch_name"`
	NeighborSwitchIPv4Addr              string   `xml:"neighbor-switch-ipv4-address" json:"neighbor_switch_ipv4_address"`
	NeighborSwitchIPv6Addr              string   `xml:"neighbor-switch-ipv6-address" json:"neighbor_switch_ipv6_address"`
	NeighborFabricName                  string   `xml:"neighbor-fabric-name" json:"neighbor_fabric_name"`
	DestinationFabricPrincipalSwitchWWN string   `xml:"destination-fabric-principal-switch-wwn" json:"destination_fabric_principal_switch_wwn"`
	CompressionActive                   bool     `xml:"compression-active-v2" json:"compression_active"`
	EncryptionActive                    bool     `xml:"encryption-active-v2" json:"encryption_active"`
	FECActive                           bool     `xml:"fec-active-v2" json:"fec_active"`
	ReservedBuffers                     uint16   `xml:"reserved-buffers" json:"reserved_buffers"`
	AverageTransmitBufferUsage          uint16   `xml:"average-transmit-buffer-usage" json:"average_transmit_buffer_usage"`
	AverageTransmitFrameSize            uint16   `xml:"average-transmit-frame-size" json:"average_transmit_frame_size"`
	AverageReceiveBufferUsage           uint16   `xml:"average-receive-buffer-usage" json:"average_receive_buffer_usage"`
	AverageReceiveFrameSize             uint16   `xml:"average-receive-frame-size" json:"average_receive_frame_size"`
	CurrentBufferUsage                  uint16   `xml:"current-buffer-usage" json:"current_buffer_usage"`
	RecommendedBuffers                  uint16   `xml:"recommended-buffers" json:"recommended_buffers"`
	MeasuredLinkDistance                string   `xml:"measured-link-distance" json:"measured_link_distance"`
	ChipInstance                        uint16   `xml:"chip-instance" json:"chip_instance"`
	ChipBuffersAvailable                uint16   `xml:"chip-buffers-available" json:"chip_buffers_available"`
	PortHealth                          string   `xml:"port-health" json:"port_health"`
	AuthenticationProtocol              string   `xml:"authentication-protocol" json:"authentication_protocol"`
	DisableReason                       string   `xml:"disable-reason" json:"disable_reason"`
	LEDomain                            uint32   `xml:"le-domain" json:"le_domain"`
	PortPeerBeaconEnabled               bool     `xml:"port-peer-beacon-enabled" json:"port_peer_beacon_enabled"`
	CleanAddressEnabled                 bool     `xml:"clean-address-enabled" json:"clean_address_enabled"`
	CongestionSignalEnabled             bool     `xml:"congestion-signal-enabled" json:"congestion_signal_enabled"`
	SegmentationReason                  string   `xml:"segmentation-reason" json:"segmentation_reason"`
	LAGName                             string   `xml:"lag-name" json:"lag_name"`
	LAGMemberLinkStatus                 string   `xml:"lag-member-link-status" json:"lag_member_link_status"`
	PortGenerationNumber                uint32   `xml:"port-generation-number" json:"port_generation_number"`
	CableDistance                       uint32   `xml:"cable-distance" json:"cable_distance"`
	Index                               uint32   `xml:"index" json:"index"`
	DefaultIndex                        uint32   `xml:"default-index" json:"default_index"`
	PhysicalState                       string   `xml:"physical-state" json:"physical_state"`
	PODLicenseStatus                    bool     `xml:"pod-license-status" json:"pod_license_status"`
	Areas                               []uint32 `xml:"areas>area" json:"areas"`
	BoundAddresses                      []string `xml:"bound-address-list>bound-address" json:"bound_addresses"`
	UserBoundEnabled                    bool     `xml:"user-bound-enabled" json:"user_bound_enabled"`
	AddressBindingMode                  string   `xml:"address-binding-mode" json:"address_binding_mode"`
	EXPortEnabled                       bool     `xml:"ex-port-enabled-v2" json:"ex_port_enabled"`
	FCRouterPortCost                    uint32   `xml:"fc-router-port-cost" json:"fc_router_port_cost"`
	EdgeFabricID                        uint32   `xml:"edge-fabric-id" json:"edge_fabric_id"`
	PreferredFrontDomainID              uint32   `xml:"preferred-front-domain-id" json:"preferred_front_domain_id"`
	GPortLocked                         uint8    `xml:"g-port-locked" json:"g_port_locked"`
	EPortDisable                        bool     `xml:"e-port-disable-v2" json:"e_port_disable"`
	EPortDisableV3                      string   `xml:"e-port-disable-v3" json:"e_port_disable_v3"`
	NPortEnabled                        bool     `xml:"n-port-enabled-v2" json:"n_port_enabled"`
	DPortEnable                         bool     `xml:"d-port-enable-v2" json:"d_port_enable"`
	PersistentDisable                   bool     `xml:"persistent-disable-v2" json:"persistent_disable"`
	ApplicationHeaderEnabled            bool     `xml:"application-header-enabled" json:"application_header_enabled"`
	DWDMLOSyncEnabled                   bool     `xml:"dwdm-losync-enabled" json:"dwdm_losync_enabled"`
	DPortDWDMEnabled                    bool     `xml:"d-port-dwdm-enabled" json:"d_port_dwdm_enabled"`
	QoSEnabled                          bool     `xml:"qos-enabled-v2" json:"qos_enabled"`
	CompressionConfigured               bool     `xml:"compression-configured-v2" json:"compression_configured"`
	EncryptionEnabled                   bool     `xml:"encryption-enabled-v2" json:"encryption_enabled"`
	TargetDrivenZoningEnable            bool     `xml:"target-driven-zoning-enable-v2" json:"target_driven_zoning_enable"`
	SimPortEnabled                      bool     `xml:"sim-port-enabled-v2" json:"sim_port_enabled"`
	MirrorPortEnabled                   bool     `xml:"mirror-port-enabled-v2" json:"mirror_port_enabled"`
	CreditRecoveryEnabled               bool     `xml:"credit-recovery-enabled-v2" json:"credit_recovery_enabled"`
	CreditRecoveryActive                bool     `xml:"credit-recovery-active-v2" json:"credit_recovery_active"`
	FECEnabled                          bool     `xml:"fec-enabled-v2" json:"fec_enabled"`
	ViaTTSFECEnabled                    bool     `xml:"via-tts-fec-enabled-v2" json:"via_tts_fec_enabled"`
	PortAutodisableEnabled              bool     `xml:"port-autodisable-enabled-v2" json:"port_autodisable_enabled"`
	CSCTLModeEnabled                    bool     `xml:"csctl-mode-enabled-v2" json:"csctl_mode_enabled"`
	FaultDelayEnabled                   bool     `xml:"fault-delay-enabled-v2" json:"fault_delay_enabled"`
	OctetSpeedComboString               string   `xml:"octet-speed-combo-string" json:"octet_speed_combo_string"`
	RSCNsuppressionEnabled              bool     `xml:"rscn-suppression-enabled-v2" json:"rscn_suppression_enabled"`
	LOStoTOVModeEnabledString           string   `xml:"los-tov-mode-enabled-string" json:"los_to_tov_mode_enabled_string"`
	FlexportProtocol                    string   `xml:"flexport-protocol" json:"flexport_protocol"`
	EthernetLAGName                     string   `xml:"ethernet-lag-name" json:"ethernet_lag_name"`
	EthernetLAGMemberTimeout            uint32   `xml:"ethernet-lag-member-timeout" json:"ethernet_lag_member_timeout"`
	EthernetLLDPEenabled                bool     `xml:"ethernet-lldp-enabled" json:"ethernet_lldp_enabled"`
	EthernetLLDPProfileName             string   `xml:"ethernet-lldp-profile-name" json:"ethernet_lldp_profile_name"`
	ICLInterconnectStatus               bool     `xml:"icl-interconnect-status" json:"icl_interconnect_status"`
	LinkLatencyDeterminationEnabled     bool     `xml:"link-latency-determination-enabled" json:"link_latency_determination_enabled"`
	PodLicenseState                     string   `xml:"pod-license-state" json:"pod_license_state"`
	PortSCN                             string   `xml:"port-scn" json:"port_scn"`
}

// GetPorts 获取交换机上所有 FC 端口的详细信息
// 对应 API: GET /brocade-interface/fibrechannel
func (c *Client) GetPorts() ([]PortInfo, error) {
	var resp FibreChannelResponse
	err := c.Get(c.endpoints().FibreChannelPorts(), &resp)
	if err != nil {
		return nil, err
	}

	return resp.Ports, nil
}
