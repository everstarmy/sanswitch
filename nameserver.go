package sanswitch

import (
	"context"
	"encoding/xml"
)

// fibreChannelNameServerResponse 是 GET /brocade-name-server/fibrechannel-name-server 的 XML 响应包装
type fibreChannelNameServerResponse struct {
	XMLName     xml.Name          `xml:"Response"`
	NameServers []NameServerEntry `xml:"fibrechannel-name-server"`
}

// NameServerEntry 描述 Fabric Name Server 中注册的一个设备条目，
// 包含端口 WWN、节点 WWN、端口类型、FC4 类型、速率、协议等字段。
// 对应 YANG 模型: brocade-name-server/fibrechannel-name-server
type NameServerEntry struct {
	XMLName                    xml.Name `xml:"fibrechannel-name-server" json:"-"`
	PortID                     string   `xml:"port-id" json:"port_id"`
	PortName                   string   `xml:"port-name" json:"port_name"`
	PortSymbolicName           string   `xml:"port-symbolic-name" json:"port_symbolic_name"`
	FabricPortName             string   `xml:"fabric-port-name" json:"fabric_port_name"`
	PermanentPortName          string   `xml:"permanent-port-name" json:"permanent_port_name"`
	NodeName                   string   `xml:"node-name" json:"node_name"`
	NodeSymbolicName           string   `xml:"node-symbolic-name" json:"node_symbolic_name"`
	ClassOfService             string   `xml:"class-of-service" json:"class_of_service"`
	FC4Type                    string   `xml:"fc4-type" json:"fc4_type"`
	FC4Features                string   `xml:"fc4-features" json:"fc4_features"`
	PortType                   string   `xml:"port-type" json:"port_type"`
	StateChangeRegistration    string   `xml:"state-change-registration" json:"state_change_registration"`
	NameServerDeviceType       string   `xml:"name-server-device-type" json:"name_server_device_type"`
	PortIndex                  uint32   `xml:"port-index" json:"port_index"`
	ShareArea                  string   `xml:"share-area" json:"share_area"`
	FrameRedirection           string   `xml:"frame-redirection" json:"frame_redirection"`
	Partial                    string   `xml:"partial" json:"partial"`
	LSAN                       string   `xml:"lsan" json:"lsan"`
	ProtocolSpeed              string   `xml:"protocol-speed" json:"protocol_speed"`
	PortProperties             string   `xml:"port-properties" json:"port_properties"`
	CascadedAG                 string   `xml:"cascaded-ag" json:"cascaded_ag"`
	ConnectedThroughAG         string   `xml:"connected-through-ag" json:"connected_through_ag"`
	RealDeviceBehindAG         string   `xml:"real-device-behind-ag" json:"real_device_behind_ag"`
	FCoEDevice                 string   `xml:"fcoe-device" json:"fcoe_device"`
	SlowDrainDeviceQuarantine  string   `xml:"slow-drain-device-quarantine" json:"slow_drain_device_quarantine"`
	ConnectedThroughFCLAG      bool     `xml:"connected-through-fc-lag" json:"connected_through_fc_lag"`
	PlatformNameID             string   `xml:"platform-name-id" json:"platform_name_id"`
	PlatformNameIDSymbolicName string   `xml:"platform-name-id-symbolic-name" json:"platform_name_id_symbolic_name"`
}

// FibreChannelNameServers 获取 Name Server 条目并允许取消请求。
func (c *Client) FibreChannelNameServers(ctx context.Context) ([]NameServerEntry, error) {
	var resp fibreChannelNameServerResponse
	err := c.get(ctx, c.endpoints().NameServer(), &resp)
	if err != nil {
		return nil, err
	}

	return resp.NameServers, nil
}
