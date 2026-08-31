package sanswitch

import (
	"context"
	"encoding/xml"
	"errors"
	"strings"
)

// definedConfigAPI 表示 Zone 定义配置中的 cfg（用于 XML 请求/响应序列化）
type definedConfigAPI struct {
	XMLName     xml.Name `xml:"cfg"`
	Name        string   `xml:"cfg-name"`
	MemberZones []string `xml:"member-zone>zone-name"` // cfg 包含的 Zone 名称列表
}

// definedConfigResponse 是 GET /brocade-zone/defined-configuration/cfg 的 XML 响应包装
type definedConfigResponse struct {
	XMLName xml.Name           `xml:"Response"`
	Configs []definedConfigAPI `xml:"cfg"`
}

// effectiveConfigAPI 表示已生效的 Zone 配置信息（用于 XML 响应序列化），
// 包含配置名称、校验和、默认 Zone 访问策略、数据库容量和事务状态等字段。
type effectiveConfigAPI struct {
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

// effectiveConfigResponse 是 GET /brocade-zone/effective-configuration 的 XML 响应包装
type effectiveConfigResponse struct {
	XMLName xml.Name           `xml:"Response"`
	Config  effectiveConfigAPI `xml:"effective-configuration"`
}

// DefinedConfigs 获取 Zone 定义配置中的所有 cfg 列表。
func (c *Client) DefinedConfigs(ctx context.Context) ([]ZoneConfig, error) {
	var resp definedConfigResponse
	err := c.get(ctx, c.endpoints().DefinedConfigs(), &resp)
	if err != nil {
		return nil, err
	}

	var configs []ZoneConfig
	for _, cfg := range resp.Configs {
		configs = append(configs, ZoneConfig{
			Name:        cfg.Name,
			Type:        "defined",
			MemberZones: cfg.MemberZones,
		})
	}

	return configs, nil
}

// EffectiveConfig 获取当前已生效的 Zone 配置信息。
func (c *Client) EffectiveConfig(ctx context.Context) (*ZoneConfig, error) {
	var resp effectiveConfigResponse
	err := c.get(ctx, c.endpoints().EffectiveConfig(), &resp)
	if err != nil {
		return nil, err
	}

	return &ZoneConfig{
		Name:              resp.Config.ConfigName,
		Type:              "effective",
		Checksum:          resp.Config.Checksum,
		DefaultZoneAccess: resp.Config.DefaultZoneAccess,
	}, nil
}

// ZoneDatabase returns Zone database capacity and transaction information.
func (c *Client) ZoneDatabase(ctx context.Context) (*ZoneDatabase, error) {
	var resp effectiveConfigResponse
	err := c.get(ctx, c.endpoints().EffectiveConfig(), &resp)
	if err != nil {
		return nil, err
	}

	return &ZoneDatabase{
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

// ZoneChecksum 获取当前 Zone 配置的校验和。
func (c *Client) ZoneChecksum(ctx context.Context) (string, error) {
	var resp effectiveConfigResponse
	err := c.get(ctx, c.endpoints().ZoneChecksum(), &resp)
	if err != nil {
		return "", err
	}
	return resp.Config.Checksum, nil
}

// UpdateDefinedConfig 更新指定 cfg 的成员 Zone 列表。
func (c *Client) UpdateDefinedConfig(ctx context.Context, name string, memberZones []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("config name required")
	}
	payload := definedConfigAPI{
		Name:        name,
		MemberZones: memberZones,
	}
	return c.patch(ctx, c.endpoints().DefinedConfigs(), payload)
}

// PatcheffectiveConfigAPI 用于 Save/Activate Zone 配置操作的请求体
type patchEffectiveConfigAPI struct {
	XMLName  xml.Name `xml:"checksum"`
	Checksum string   `xml:",chardata"` // 当前配置校验和，用于防止并发修改冲突
}

// SaveZoneConfig 保存当前 Zone 定义配置。
func (c *Client) SaveZoneConfig(ctx context.Context, checksum string) error {
	if strings.TrimSpace(checksum) == "" {
		return errors.New("zone checksum required")
	}
	payload := patchEffectiveConfigAPI{
		Checksum: checksum,
	}
	return c.patch(ctx, c.endpoints().ZoneSaveConfig(), payload)
}

// ActivateZoneConfig 激活指定的 Zone 配置。
func (c *Client) ActivateZoneConfig(ctx context.Context, name string, checksum string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(checksum) == "" {
		return errors.New("config name and zone checksum required")
	}
	payload := patchEffectiveConfigAPI{
		Checksum: checksum,
	}
	return c.patch(ctx, c.endpoints().ZoneActivateConfig(name), payload)
}
