package sanswitch

import (
	"context"
	"encoding/xml"
)

// ==================== SNMP System ====================

// sNMPSystemResponse 对应 GET /brocade-snmp/system
type sNMPSystemResponse struct {
	XMLName xml.Name   `xml:"Response"`
	System  SNMPSystem `xml:"system"`
}

// SNMPSystem 描述 SNMP 系统级配置
type SNMPSystem struct {
	XMLName                xml.Name `xml:"system" json:"-"`
	Description            string   `xml:"description" json:"description"`
	Location               string   `xml:"location" json:"location"`
	Contact                string   `xml:"contact" json:"contact"`
	InformsEnabled         bool     `xml:"informs-enabled" json:"informs_enabled"`
	EncryptionEnabled      bool     `xml:"encryption-enabled" json:"encryption_enabled"`
	AuditInterval          uint16   `xml:"audit-interval" json:"audit_interval"`
	SecurityGetLevelString string   `xml:"security-get-level-string" json:"security_get_level_string"`
	SecuritySetLevelString string   `xml:"security-set-level-string" json:"security_set_level_string"`
}

// ==================== SNMPv1 Account ====================

// sNMPv1AccountResponse 对应 GET /brocade-snmp/v1-account
type sNMPv1AccountResponse struct {
	XMLName  xml.Name        `xml:"Response"`
	Accounts []SNMPv1Account `xml:"v1-account"`
}

// SNMPv1Account 描述一个 SNMPv1 社区字符串账户
type SNMPv1Account struct {
	XMLName        xml.Name `xml:"v1-account" json:"-"`
	Index          uint16   `xml:"index" json:"index"`
	CommunityName  string   `xml:"community-name" json:"-"`
	CommunityGroup string   `xml:"community-group" json:"community_group"`
}

// ==================== SNMPv1 Trap ====================

// sNMPv1TrapResponse 对应 GET /brocade-snmp/v1-trap
type sNMPv1TrapResponse struct {
	XMLName xml.Name     `xml:"Response"`
	Traps   []SNMPv1Trap `xml:"v1-trap"`
}

// SNMPv1Trap 描述一个 SNMPv1 Trap 接收者配置
type SNMPv1Trap struct {
	XMLName           xml.Name `xml:"v1-trap" json:"-"`
	Index             uint16   `xml:"index" json:"index"`
	Host              string   `xml:"host" json:"host"`
	TrapSeverityLevel string   `xml:"trap-severity-level" json:"trap_severity_level"`
	PortNumber        uint16   `xml:"port-number" json:"port_number"`
}

// ==================== SNMPv3 Account ====================

// sNMPv3AccountResponse 对应 GET /brocade-snmp/v3-account
type sNMPv3AccountResponse struct {
	XMLName  xml.Name        `xml:"Response"`
	Accounts []SNMPv3Account `xml:"v3-account"`
}

// SNMPv3Account 描述一个 SNMPv3 用户账户
type SNMPv3Account struct {
	XMLName                xml.Name `xml:"v3-account" json:"-"`
	Index                  uint16   `xml:"index" json:"index"`
	UserName               string   `xml:"user-name" json:"user_name"`
	UserGroup              string   `xml:"user-group" json:"user_group"`
	AuthenticationProtocol string   `xml:"authentication-protocol" json:"authentication_protocol"`
	PrivacyProtocol        string   `xml:"privacy-protocol" json:"privacy_protocol"`
	AuthenticationPassword string   `xml:"authentication-password" json:"-"`
	PrivacyPassword        string   `xml:"privacy-password" json:"-"`
	ManagerEngineID        string   `xml:"manager-engine-id" json:"manager_engine_id"`
}

// ==================== SNMPv3 Trap ====================

// sNMPv3TrapResponse 对应 GET /brocade-snmp/v3-trap
type sNMPv3TrapResponse struct {
	XMLName xml.Name     `xml:"Response"`
	Traps   []SNMPv3Trap `xml:"v3-trap"`
}

// SNMPv3Trap 描述一个 SNMPv3 Trap 接收者配置
type SNMPv3Trap struct {
	XMLName           xml.Name `xml:"v3-trap" json:"-"`
	TrapIndex         uint16   `xml:"trap-index" json:"trap_index"`
	USMIndex          uint16   `xml:"usm-index" json:"usm_index"`
	Host              string   `xml:"host" json:"host"`
	TrapSeverityLevel string   `xml:"trap-severity-level" json:"trap_severity_level"`
	PortNumber        uint16   `xml:"port-number" json:"port_number"`
	InformsEnabled    bool     `xml:"informs-enabled" json:"informs_enabled"`
}

// ==================== Client Methods ====================

// SNMPSystem 获取 SNMP 系统级配置并允许取消请求。
func (c *Client) SNMPSystem(ctx context.Context) (*SNMPSystem, error) {
	var resp sNMPSystemResponse
	err := c.get(ctx, c.endpoints().SNMPSystem(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.System, nil
}

// SNMPv1Accounts 获取 SNMPv1 社区字符串账户并允许取消请求。
func (c *Client) SNMPv1Accounts(ctx context.Context) ([]SNMPv1Account, error) {
	var resp sNMPv1AccountResponse
	err := c.get(ctx, c.endpoints().SNMPv1Accounts(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}

// SNMPv1Traps 获取 SNMPv1 Trap 接收者并允许取消请求。
func (c *Client) SNMPv1Traps(ctx context.Context) ([]SNMPv1Trap, error) {
	var resp sNMPv1TrapResponse
	err := c.get(ctx, c.endpoints().SNMPv1Traps(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Traps, nil
}

// SNMPv3Accounts 获取 SNMPv3 用户账户并允许取消请求。
func (c *Client) SNMPv3Accounts(ctx context.Context) ([]SNMPv3Account, error) {
	var resp sNMPv3AccountResponse
	err := c.get(ctx, c.endpoints().SNMPv3Accounts(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}

// SNMPv3Traps 获取 SNMPv3 Trap 接收者并允许取消请求。
func (c *Client) SNMPv3Traps(ctx context.Context) ([]SNMPv3Trap, error) {
	var resp sNMPv3TrapResponse
	err := c.get(ctx, c.endpoints().SNMPv3Traps(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Traps, nil
}
