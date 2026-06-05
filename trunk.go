package san

import "encoding/xml"

// ==================== Trunk Members ====================

// TrunkResponse 对应 GET /brocade-fibrechannel-trunk/trunk
type TrunkResponse struct {
	XMLName xml.Name    `xml:"Response"`
	Trunks  []TrunkInfo `xml:"trunk"`
}

// TrunkInfo 描述一条 ISL Trunk 链路的成员信息
type TrunkInfo struct {
	XMLName            xml.Name `xml:"trunk" json:"-"`
	Group              uint32   `xml:"group" json:"group"`
	SourcePort         uint32   `xml:"source-port" json:"source_port"`
	DestinationPort    uint32   `xml:"destination-port" json:"destination_port"`
	NeighborWWN        string   `xml:"neighbor-wwn" json:"neighbor_wwn"`
	NeighborSwitchName string   `xml:"neighbor-switch-name" json:"neighbor_switch_name"`
	NeighborDomainID   uint32   `xml:"neighbor-domain-id" json:"neighbor_domain_id"`
	Deskew             uint32   `xml:"deskew" json:"deskew"`
	Master             bool     `xml:"master" json:"master"`
	TrunkType          string   `xml:"trunk-type" json:"trunk_type"`
}

// ==================== Trunk Performance ====================

// TrunkPerformanceResponse 对应 GET /brocade-fibrechannel-trunk/performance
type TrunkPerformanceResponse struct {
	XMLName      xml.Name               `xml:"Response"`
	Performances []TrunkPerformanceInfo `xml:"performance"`
}

// TrunkPerformanceInfo 描述一条 Trunk 的性能统计
type TrunkPerformanceInfo struct {
	XMLName         xml.Name `xml:"performance" json:"-"`
	Group           uint32   `xml:"group" json:"group"`
	TxBandwidth     uint32   `xml:"tx-bandwidth" json:"tx_bandwidth"`
	TxCapacity      float64  `xml:"tx-capacity" json:"tx_capacity"`
	TxUtilization   uint64   `xml:"tx-utilization" json:"tx_utilization"`
	TxThroughput    uint64   `xml:"tx-throughput" json:"tx_throughput"` // Deprecated: use TxUtilization
	TxPercentage    float64  `xml:"tx-percentage" json:"tx_percentage"`
	RxBandwidth     uint32   `xml:"rx-bandwidth" json:"rx_bandwidth"`
	RxCapacity      float64  `xml:"rx-capacity" json:"rx_capacity"`
	RxUtilization   uint64   `xml:"rx-utilization" json:"rx_utilization"`
	RxThroughput    uint64   `xml:"rx-throughput" json:"rx_throughput"` // Deprecated: use RxUtilization
	RxPercentage    float64  `xml:"rx-percentage" json:"rx_percentage"`
	TxRxBandwidth   uint32   `xml:"txrx-bandwidth" json:"txrx_bandwidth"`
	TxRxCapacity    float64  `xml:"txrx-capacity" json:"txrx_capacity"`
	TxRxUtilization uint64   `xml:"txrx-utilization" json:"txrx_utilization"`
	TxRxThroughput  uint64   `xml:"txrx-throughput" json:"txrx_throughput"` // Deprecated: use TxRxUtilization
	TxRxPercentage  float64  `xml:"txrx-percentage" json:"txrx_percentage"`
}

// ==================== Trunk Area ====================

// TrunkAreaResponse 对应 GET /brocade-fibrechannel-trunk/trunk-area
type TrunkAreaResponse struct {
	XMLName    xml.Name        `xml:"Response"`
	TrunkAreas []TrunkAreaInfo `xml:"trunk-area"`
}

// TrunkAreaInfo 描述一个 Trunk Area 组
type TrunkAreaInfo struct {
	XMLName      xml.Name `xml:"trunk-area" json:"-"`
	TrunkIndex   uint32   `xml:"trunk-index" json:"trunk_index"`
	MasterPort   string   `xml:"master-port" json:"master_port"`
	TrunkActive  bool     `xml:"trunk-active" json:"trunk_active"`
	TrunkMembers []string `xml:"trunk-members>trunk-member" json:"trunk_members"`
}

// ==================== Client Methods ====================

// GetTrunks 获取交换机上所有 Trunk 成员列表
func (c *Client) GetTrunks() ([]TrunkInfo, error) {
	var resp TrunkResponse
	err := c.Get(c.endpoints().Trunks(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Trunks, nil
}

// GetTrunkPerformances 获取所有 Trunk 的性能统计
func (c *Client) GetTrunkPerformances() ([]TrunkPerformanceInfo, error) {
	var resp TrunkPerformanceResponse
	err := c.Get(c.endpoints().TrunkPerformances(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Performances, nil
}

// GetTrunkAreas 获取所有 Trunk Area 组信息
func (c *Client) GetTrunkAreas() ([]TrunkAreaInfo, error) {
	var resp TrunkAreaResponse
	err := c.Get(c.endpoints().TrunkAreas(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.TrunkAreas, nil
}
