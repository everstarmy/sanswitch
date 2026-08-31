package sanswitch

import (
	"context"
	"encoding/xml"
	"net"
	"strings"
)

// fabricSwitchResponse 是 GET /brocade-fabric/fabric-switch 的 XML 响应包装
type fabricSwitchResponse struct {
	XMLName  xml.Name       `xml:"Response"`
	Switches []FabricSwitch `xml:"fabric-switch"`
}

// FabricSwitch 描述 Fabric 中的一台交换机，对应 brocade-fabric/fabric-switch YANG leaf
type FabricSwitch struct {
	XMLName                 xml.Name `xml:"fabric-switch" json:"-"`
	Name                    string   `xml:"name" json:"name"`
	SwitchUserFriendlyName  string   `xml:"switch-user-friendly-name" json:"switch_user_friendly_name"`
	ChassisWWN              string   `xml:"chassis-wwn" json:"chassis_wwn"`
	ChassisUserFriendlyName string   `xml:"chassis-user-friendly-name" json:"chassis_user_friendly_name"`
	DomainID                int      `xml:"domain-id" json:"domain_id"`
	Fcid                    string   `xml:"fcid" json:"fcid"`
	FcidHex                 string   `xml:"fcid-hex" json:"fcid_hex"`
	IPAddress               string   `xml:"ip-address" json:"ip_address"`
	FCIPAddress             string   `xml:"fcip-address" json:"fcip_address"`
	IPv6Address             string   `xml:"ipv6-address" json:"ipv6_address"`
	FirmwareVersion         string   `xml:"firmware-version" json:"firmware_version"`
	SwitchModel             string   `xml:"switch-model" json:"switch_model"`
	SerialNumber            string   `xml:"serial-number" json:"serial_number"`
	Principal               int      `xml:"principal" json:"principal"`
	IsPrincipal             bool     `xml:"is-principal" json:"is_principal"`
	VFID                    int      `xml:"vf-id,omitempty" json:"vf_id,omitempty"`
}

// FabricSwitches 获取 Fabric 中所有交换机的详细信息。
func (c *Client) FabricSwitches(ctx context.Context) ([]FabricSwitch, error) {
	var resp fabricSwitchResponse
	err := c.get(ctx, c.endpoints().FabricSwitches(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Switches, nil
}

// Switch returns summary information for the connected switch.
func (c *Client) Switch(ctx context.Context) (*Switch, error) {
	switches, err := c.FabricSwitches(ctx)
	if err != nil {
		return nil, err
	}

	if len(switches) == 0 {
		return nil, ErrNotFound
	}

	sw := c.localFabricSwitch(switches)
	return &Switch{
		Name:            sw.SwitchUserFriendlyName,
		WWN:             sw.Name,
		DomainID:        sw.DomainID,
		FirmwareVersion: sw.FirmwareVersion,
		ModelName:       sw.SwitchModel,
		SerialNumber:    sw.SerialNumber,
		IPAddress:       sw.IPAddress,
		IPv6Address:     sw.IPv6Address,
		ChassisWWN:      sw.ChassisWWN,
		Fcid:            sw.Fcid,
		FcidHex:         sw.FcidHex,
		Principal:       sw.IsPrincipal || sw.Principal != 0,
	}, nil
}

func (c *Client) localFabricSwitch(switches []FabricSwitch) FabricSwitch {
	host := strings.TrimSpace(c.switchAddress)
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		host = parsedHost
	}
	host = strings.Trim(host, "[]")
	for _, sw := range switches {
		if strings.EqualFold(host, sw.IPAddress) || strings.EqualFold(host, sw.IPv6Address) {
			return sw
		}
	}
	return switches[0]
}
