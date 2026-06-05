package san

import (
	"encoding/xml"
)

// DefinedConfigAPI 表示 Zone 定义配置中的 cfg（用于 XML 请求/响应序列化）
type DefinedConfigAPI struct {
	XMLName     xml.Name `xml:"cfg"`
	Name        string   `xml:"cfg-name"`
	MemberZones []string `xml:"member-zone>zone-name"` // cfg 包含的 Zone 名称列表
}

// DefinedConfigResponse 是 GET /brocade-zone/defined-configuration/cfg 的 XML 响应包装
type DefinedConfigResponse struct {
	XMLName xml.Name           `xml:"Response"`
	Configs []DefinedConfigAPI `xml:"cfg"`
}

// EffectiveConfigAPI 表示已生效的 Zone 配置信息（用于 XML 响应序列化），
// 包含配置名称、校验和、默认 Zone 访问策略、数据库容量和事务状态等字段。
type EffectiveConfigAPI struct {
	XMLName                xml.Name `xml:"effective-configuration"`
	ConfigName             string   `xml:"cfg-name"`
	Checksum               string   `xml:"checksum"`
	DefaultZoneAccess      string   `xml:"default-zone-access-v2"`
	DBMax                  uint32   `xml:"db-max"`
	DBAvail                uint32   `xml:"db-avail"`
	DBCommitted            uint32   `xml:"db-committed"`
	DBTransaction          uint32   `xml:"db-transaction"`
	TransactionToken       uint32   `xml:"transaction-token"`
	DBChassisWideCommitted uint32   `xml:"db-chassis-wide-committed"`
	DBChassisWideMax       uint32   `xml:"db-chassis-wide-max"`
	DBFabricWideMax        uint32   `xml:"db-fabric-wide-max"`
	DomainWithLowestDBMax  uint32   `xml:"domain-with-lowest-db-max"`
}

// EffectiveConfigResponse 是 GET /brocade-zone/effective-configuration 的 XML 响应包装
type EffectiveConfigResponse struct {
	XMLName xml.Name           `xml:"Response"`
	Config  EffectiveConfigAPI `xml:"effective-configuration"`
}

// GetDefinedConfigs 获取 Zone 定义配置中的所有 cfg 列表。
// 对应 API: GET /brocade-zone/defined-configuration/cfg
func (c *Client) GetDefinedConfigs() ([]ConfigInfo, error) {
	var resp DefinedConfigResponse
	err := c.Get(c.endpoints().DefinedConfigs(), &resp)
	if err != nil {
		return nil, err
	}

	var configs []ConfigInfo
	for _, cfg := range resp.Configs {
		configs = append(configs, ConfigInfo{
			Name:        cfg.Name,
			Type:        "defined",
			MemberZones: cfg.MemberZones,
		})
	}

	return configs, nil
}

// GetEffectiveConfig 获取当前已生效的 Zone 配置信息，包含配置名称、校验和、默认 Zone 访问策略。
// 对应 API: GET /brocade-zone/effective-configuration
func (c *Client) GetEffectiveConfig() (*ConfigInfo, error) {
	var resp EffectiveConfigResponse
	err := c.Get(c.endpoints().EffectiveConfig(), &resp)
	if err != nil {
		return nil, err
	}

	return &ConfigInfo{
		Name:              resp.Config.ConfigName,
		Type:              "effective",
		Checksum:          resp.Config.Checksum,
		DefaultZoneAccess: resp.Config.DefaultZoneAccess,
	}, nil
}

// GetZoneDatabaseInfo 获取 Zone 数据库的容量和事务状态信息。
// 对应 API: GET /brocade-zone/effective-configuration
func (c *Client) GetZoneDatabaseInfo() (*ZoneDatabaseInfo, error) {
	var resp EffectiveConfigResponse
	err := c.Get(c.endpoints().EffectiveConfig(), &resp)
	if err != nil {
		return nil, err
	}

	return &ZoneDatabaseInfo{
		DBMax:                  resp.Config.DBMax,
		DBAvail:                resp.Config.DBAvail,
		DBCommitted:            resp.Config.DBCommitted,
		DBTransaction:          resp.Config.DBTransaction,
		TransactionToken:       resp.Config.TransactionToken,
		DBChassisWideCommitted: resp.Config.DBChassisWideCommitted,
		DBChassisWideMax:       resp.Config.DBChassisWideMax,
		DBFabricWideMax:        resp.Config.DBFabricWideMax,
		DomainWithLowestDBMax:  resp.Config.DomainWithLowestDBMax,
	}, nil
}

// GetZoneChecksum 获取当前 Zone 配置的校验和，用于后续的 Save/Activate 操作。
// 对应 API: GET /brocade-zone/effective-configuration/checksum
func (c *Client) GetZoneChecksum() (string, error) {
	var resp EffectiveConfigResponse
	err := c.Get(c.endpoints().ZoneChecksum(), &resp)
	if err != nil {
		return "", err
	}
	return resp.Config.Checksum, nil
}

// UpdateDefinedConfig 更新 Zone 定义配置中指定 cfg 的成员 Zone 列表。
// 对应 API: PATCH /brocade-zone/defined-configuration/cfg
func (c *Client) UpdateDefinedConfig(name string, memberZones []string) error {
	payload := DefinedConfigAPI{
		Name:        name,
		MemberZones: memberZones,
	}
	return c.Patch(c.endpoints().DefinedConfigs(), payload)
}

// PatchEffectiveConfigAPI 用于 Save/Activate Zone 配置操作的请求体
type PatchEffectiveConfigAPI struct {
	XMLName  xml.Name `xml:"checksum"`
	Checksum string   `xml:",chardata"` // 当前配置校验和，用于防止并发修改冲突
}

// SaveZoneConfig 保存当前 Zone 定义配置到持久存储。
// checksum 为当前有效配置的校验和，用于防止并发修改冲突。
// 对应 API: PATCH /brocade-zone/effective-configuration/cfg-action-v2/save
func (c *Client) SaveZoneConfig(checksum string) error {
	payload := PatchEffectiveConfigAPI{
		Checksum: checksum,
	}
	return c.Patch(c.endpoints().ZoneSaveConfig(), payload)
}

// ActivateZoneConfig 激活指定的 Zone 配置（cfg），使其成为当前生效配置。
// checksum 为当前有效配置的校验和，用于防止并发修改冲突。
// 对应 API: PATCH /brocade-zone/effective-configuration/cfg-name/{name}
func (c *Client) ActivateZoneConfig(name string, checksum string) error {
	payload := PatchEffectiveConfigAPI{
		Checksum: checksum,
	}
	return c.Patch(c.endpoints().ZoneActivateConfig(name), payload)
}
