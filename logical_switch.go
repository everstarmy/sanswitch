package san

import "encoding/xml"

// LogicalSwitchResponse 是 GET /brocade-fibrechannel-logical-switch/fibrechannel-logical-switch 的 XML 响应包装
type LogicalSwitchResponse struct {
	XMLName         xml.Name            `xml:"Response"`
	LogicalSwitches []LogicalSwitchInfo `xml:"fibrechannel-logical-switch"`
}

// LogicalSwitchInfo 描述一个逻辑交换机的配置信息，包含 Fabric ID、WWN、
// 基础交换机标志、FICON 模式、端口成员列表等字段。
// 对应 YANG 模型: brocade-fibrechannel-logical-switch/fibrechannel-logical-switch
type LogicalSwitchInfo struct {
	XMLName                xml.Name `xml:"fibrechannel-logical-switch" json:"-"`
	FabricID               uint32   `xml:"fabric-id" json:"fabric_id"`                                 // Fabric ID
	SwitchWWN              string   `xml:"switch-wwn" json:"switch_wwn"`                               // 交换机 WWN
	BaseSwitchEnabled      bool     `xml:"base-switch-enabled-v2" json:"base_switch_enabled"`          // 是否为基础交换机
	DefaultSwitch          bool     `xml:"default-switch" json:"default_switch"`                       // 是否为默认交换机
	LogicalISLEnabled      bool     `xml:"logical-isl-enabled-v2" json:"logical_isl_enabled"`          // 逻辑 ISL 是否启用
	FiconModeEnabled       bool     `xml:"ficon-mode-enabled-v2" json:"ficon_mode_enabled"`            // FICON 模式是否启用
	IPStorageSwitchEnabled bool     `xml:"ip-storage-switch-enabled" json:"ip_storage_switch_enabled"` // IP 存储交换机是否启用
	PortMembers            []string `xml:"port-member-list>port-member" json:"port_members"`           // FC 端口成员列表
	GePortMembers          []string `xml:"ge-port-member-list>port-member" json:"ge_port_members"`     // GE 端口成员列表
	PortIndexMembers       []uint32 `xml:"port-index-members>port-index" json:"port_index_members"`    // 端口索引成员列表
}

// GetLogicalSwitches 获取交换机上所有逻辑交换机的配置信息。
// 对应 API: GET /brocade-fibrechannel-logical-switch/fibrechannel-logical-switch
func (c *Client) GetLogicalSwitches() ([]LogicalSwitchInfo, error) {
	var resp LogicalSwitchResponse
	err := c.Get("/brocade-fibrechannel-logical-switch/fibrechannel-logical-switch", &resp)
	if err != nil {
		return nil, err
	}

	return resp.LogicalSwitches, nil
}
