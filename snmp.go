package san

import "encoding/xml"

// ==================== SNMP System ====================

// SNMPSystemResponse 对应 GET /brocade-snmp/system
type SNMPSystemResponse struct {
	XMLName xml.Name       `xml:"Response"`
	System  SNMPSystemInfo `xml:"system"`
}

// SNMPSystemInfo 描述 SNMP 系统级配置
type SNMPSystemInfo struct {
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

// SNMPv1AccountResponse 对应 GET /brocade-snmp/v1-account
type SNMPv1AccountResponse struct {
	XMLName  xml.Name            `xml:"Response"`
	Accounts []SNMPv1AccountInfo `xml:"v1-account"`
}

// SNMPv1AccountInfo 描述一个 SNMPv1 社区字符串账户
type SNMPv1AccountInfo struct {
	XMLName        xml.Name `xml:"v1-account" json:"-"`
	Index          uint16   `xml:"index" json:"index"`
	CommunityName  string   `xml:"community-name" json:"community_name"`
	CommunityGroup string   `xml:"community-group" json:"community_group"`
}

// ==================== SNMPv1 Trap ====================

// SNMPv1TrapResponse 对应 GET /brocade-snmp/v1-trap
type SNMPv1TrapResponse struct {
	XMLName xml.Name         `xml:"Response"`
	Traps   []SNMPv1TrapInfo `xml:"v1-trap"`
}

// SNMPv1TrapInfo 描述一个 SNMPv1 Trap 接收者配置
type SNMPv1TrapInfo struct {
	XMLName           xml.Name `xml:"v1-trap" json:"-"`
	Index             uint16   `xml:"index" json:"index"`
	Host              string   `xml:"host" json:"host"`
	TrapSeverityLevel string   `xml:"trap-severity-level" json:"trap_severity_level"`
	PortNumber        uint16   `xml:"port-number" json:"port_number"`
}

// ==================== SNMPv3 Account ====================

// SNMPv3AccountResponse 对应 GET /brocade-snmp/v3-account
type SNMPv3AccountResponse struct {
	XMLName  xml.Name            `xml:"Response"`
	Accounts []SNMPv3AccountInfo `xml:"v3-account"`
}

// SNMPv3AccountInfo 描述一个 SNMPv3 用户账户
type SNMPv3AccountInfo struct {
	XMLName                xml.Name `xml:"v3-account" json:"-"`
	Index                  uint16   `xml:"index" json:"index"`
	UserName               string   `xml:"user-name" json:"user_name"`
	UserGroup              string   `xml:"user-group" json:"user_group"`
	AuthenticationProtocol string   `xml:"authentication-protocol" json:"authentication_protocol"`
	PrivacyProtocol        string   `xml:"privacy-protocol" json:"privacy_protocol"`
	AuthenticationPassword string   `xml:"authentication-password" json:"authentication_password,omitempty"`
	PrivacyPassword        string   `xml:"privacy-password" json:"privacy_password,omitempty"`
	ManagerEngineID        string   `xml:"manager-engine-id" json:"manager_engine_id"`
}

// ==================== SNMPv3 Trap ====================

// SNMPv3TrapResponse 对应 GET /brocade-snmp/v3-trap
type SNMPv3TrapResponse struct {
	XMLName xml.Name         `xml:"Response"`
	Traps   []SNMPv3TrapInfo `xml:"v3-trap"`
}

// SNMPv3TrapInfo 描述一个 SNMPv3 Trap 接收者配置
type SNMPv3TrapInfo struct {
	XMLName           xml.Name `xml:"v3-trap" json:"-"`
	TrapIndex         uint16   `xml:"trap-index" json:"trap_index"`
	USMIndex          uint16   `xml:"usm-index" json:"usm_index"`
	Host              string   `xml:"host" json:"host"`
	TrapSeverityLevel string   `xml:"trap-severity-level" json:"trap_severity_level"`
	PortNumber        uint16   `xml:"port-number" json:"port_number"`
	InformsEnabled    bool     `xml:"informs-enabled" json:"informs_enabled"`
}

// ==================== Client Methods ====================

// GetSNMPSystem 获取 SNMP 系统级配置
func (c *Client) GetSNMPSystem() (*SNMPSystemInfo, error) {
	var resp SNMPSystemResponse
	err := c.Get(c.endpoints().SNMPSystem(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.System, nil
}

// GetSNMPv1Accounts 获取所有 SNMPv1 社区字符串账户
func (c *Client) GetSNMPv1Accounts() ([]SNMPv1AccountInfo, error) {
	var resp SNMPv1AccountResponse
	err := c.Get(c.endpoints().SNMPv1Accounts(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}

// GetSNMPv1Traps 获取所有 SNMPv1 Trap 接收者配置
func (c *Client) GetSNMPv1Traps() ([]SNMPv1TrapInfo, error) {
	var resp SNMPv1TrapResponse
	err := c.Get(c.endpoints().SNMPv1Traps(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Traps, nil
}

// GetSNMPv3Accounts 获取所有 SNMPv3 用户账户
func (c *Client) GetSNMPv3Accounts() ([]SNMPv3AccountInfo, error) {
	var resp SNMPv3AccountResponse
	err := c.Get(c.endpoints().SNMPv3Accounts(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}

// GetSNMPv3Traps 获取所有 SNMPv3 Trap 接收者配置
func (c *Client) GetSNMPv3Traps() ([]SNMPv3TrapInfo, error) {
	var resp SNMPv3TrapResponse
	err := c.Get(c.endpoints().SNMPv3Traps(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Traps, nil
}
