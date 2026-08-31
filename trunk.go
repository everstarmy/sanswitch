package sanswitch

import (
	"context"
	"encoding/xml"
)

// ==================== Trunk Members ====================

// trunkResponse 对应 GET /brocade-fibrechannel-trunk/trunk
type trunkResponse struct {
	XMLName xml.Name `xml:"Response"`
	Trunks  []Trunk  `xml:"trunk"`
}

// Trunk 描述一条 ISL Trunk 链路的成员信息
type Trunk struct {
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

// trunkPerformanceResponse 对应 GET /brocade-fibrechannel-trunk/performance
type trunkPerformanceResponse struct {
	XMLName      xml.Name           `xml:"Response"`
	Performances []TrunkPerformance `xml:"performance"`
}

// TrunkPerformance 描述一条 Trunk 的性能统计
type TrunkPerformance struct {
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

// trunkAreaResponse 对应 GET /brocade-fibrechannel-trunk/trunk-area
type trunkAreaResponse struct {
	XMLName    xml.Name    `xml:"Response"`
	TrunkAreas []TrunkArea `xml:"trunk-area"`
}

// TrunkArea 描述一个 Trunk Area 组
type TrunkArea struct {
	XMLName      xml.Name `xml:"trunk-area" json:"-"`
	TrunkIndex   uint32   `xml:"trunk-index" json:"trunk_index"`
	MasterPort   string   `xml:"master-port" json:"master_port"`
	TrunkActive  bool     `xml:"trunk-active" json:"trunk_active"`
	TrunkMembers []string `xml:"trunk-members>trunk-member" json:"trunk_members"`
}

// ==================== Client Methods ====================

// Trunks 获取交换机上所有 Trunk 成员并允许取消请求。
func (c *Client) Trunks(ctx context.Context) ([]Trunk, error) {
	var resp trunkResponse
	err := c.get(ctx, c.endpoints().Trunks(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Trunks, nil
}

// TrunkPerformances 获取 Trunk 性能统计并允许取消请求。
func (c *Client) TrunkPerformances(ctx context.Context) ([]TrunkPerformance, error) {
	var resp trunkPerformanceResponse
	err := c.get(ctx, c.endpoints().TrunkPerformances(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Performances, nil
}

// TrunkAreas 获取 Trunk Area 信息并允许取消请求。
func (c *Client) TrunkAreas(ctx context.Context) ([]TrunkArea, error) {
	var resp trunkAreaResponse
	err := c.get(ctx, c.endpoints().TrunkAreas(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.TrunkAreas, nil
}
