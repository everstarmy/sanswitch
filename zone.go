package san

import (
	"encoding/xml"
	"errors"
	"net/url"
	"strings"
)

// DefinedZoneAPI 表示 Zone 定义配置中的 Zone（用于 XML 请求/响应序列化）。
// MemberEntryNames 为普通成员（entry-name），PrincipalEntryNames 为 Principal 成员（principal-entry-name）。
type DefinedZoneAPI struct {
	XMLName             xml.Name `xml:"zone"`
	Name                string   `xml:"zone-name"`
	ZoneType            string   `xml:"zone-type"`
	ZoneTypeString      string   `xml:"zone-type-string"`
	MemberEntryNames    []string `xml:"member-entry>entry-name"`           // 普通成员列表
	PrincipalEntryNames []string `xml:"member-entry>principal-entry-name"` // Principal 成员列表
}

// DefinedZoneResponse 是 GET /brocade-zone/defined-configuration/zone 的 XML 响应包装
type DefinedZoneResponse struct {
	XMLName xml.Name         `xml:"Response"`
	Zones   []DefinedZoneAPI `xml:"zone"`
}

// EffectiveZoneAPI 表示已生效配置中的 Zone（用于 XML 响应序列化）。
// 与 DefinedZoneAPI 类似，但对应 YANG 模型中的 enabled-zone 节点。
type EffectiveZoneAPI struct {
	XMLName             xml.Name `xml:"enabled-zone"`
	Name                string   `xml:"zone-name"`
	ZoneType            string   `xml:"zone-type"`
	ZoneTypeString      string   `xml:"zone-type-string"`
	MemberEntryNames    []string `xml:"member-entry>entry-name"`           // 普通成员列表
	PrincipalEntryNames []string `xml:"member-entry>principal-entry-name"` // Principal 成员列表
}

// EffectiveZoneResponse 是 GET /brocade-zone/effective-configuration/enabled-zone 的 XML 响应包装
type EffectiveZoneResponse struct {
	XMLName xml.Name           `xml:"Response"`
	Zones   []EffectiveZoneAPI `xml:"enabled-zone"`
}

// GetDefinedZones 获取 Zone 定义配置中的所有 Zone 列表。
// 返回的 ZoneInfo 中 Members 字段包含普通成员和 Principal 成员的合并列表。
// 对应 API: GET /brocade-zone/defined-configuration/zone
func (c *Client) GetDefinedZones() ([]ZoneInfo, error) {
	var resp DefinedZoneResponse
	err := c.Get("/brocade-zone/defined-configuration/zone", &resp)
	if err != nil {
		return nil, err
	}

	var zones []ZoneInfo
	for _, z := range resp.Zones {
		members := append(z.MemberEntryNames, z.PrincipalEntryNames...)
		zones = append(zones, ZoneInfo{
			Name:        z.Name,
			Members:     members,
			Description: z.ZoneTypeString,
			Type:        z.ZoneType,
		})
	}

	return zones, nil
}

// GetEffectiveZones 获取已生效配置中的所有 Zone 列表。
// 返回的 ZoneInfo 中 Members 字段包含普通成员和 Principal 成员的合并列表。
// 对应 API: GET /brocade-zone/effective-configuration/enabled-zone
func (c *Client) GetEffectiveZones() ([]ZoneInfo, error) {
	var resp EffectiveZoneResponse
	err := c.Get("/brocade-zone/effective-configuration/enabled-zone", &resp)
	if err != nil {
		return nil, err
	}

	var zones []ZoneInfo
	for _, z := range resp.Zones {
		members := append(z.MemberEntryNames, z.PrincipalEntryNames...)
		zones = append(zones, ZoneInfo{
			Name:        z.Name,
			Members:     members,
			Description: z.ZoneTypeString,
			Type:        z.ZoneType,
		})
	}

	return zones, nil
}

// CreateZone 在 Zone 定义配置中创建一个新的 Zone。
// members 为普通成员（entry-name），principalMembers 为 Principal 成员（principal-entry-name）。
// 对应 API: POST /brocade-zone/defined-configuration/zone
func (c *Client) CreateZone(name string, members []string, principalMembers []string) error {
	payload := DefinedZoneAPI{
		Name:                name,
		ZoneTypeString:      "zone",
		MemberEntryNames:    members,
		PrincipalEntryNames: principalMembers,
	}
	return c.Post("/brocade-zone/defined-configuration/zone", payload)
}

// CreateZoneAndActivate 执行完整的 Zone 创建并激活流程：
// 1. 获取当前 checksum
// 2. 创建 Zone
// 3. 将 Zone 添加到指定 cfg 配置中
// 4. 保存配置并激活
// 若 Zone 已在 cfg 中则不会重复添加。
func (c *Client) CreateZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error {
	if err := validateZoneActivationInput(cfgName, zoneName, members); err != nil {
		return err
	}

	checksum, err := c.GetZoneChecksum()
	if err != nil {
		return err
	}
	if err := c.CreateZone(zoneName, members, principalMembers); err != nil {
		return err
	}
	configs, err := c.GetDefinedConfigs()
	if err != nil {
		return err
	}
	memberZones, err := memberZonesForConfig(configs, cfgName)
	if err != nil {
		return err
	}
	if !containsString(memberZones, zoneName) {
		memberZones = append(memberZones, zoneName)
	}
	if err := c.UpdateDefinedConfig(cfgName, memberZones); err != nil {
		return err
	}
	return c.saveAndActivateZoneConfig(cfgName, checksum)
}

// UpdateZone 更新 Zone 定义配置中已有 Zone 的成员列表。
// members 为普通成员，principalMembers 为 Principal 成员。
// 对应 API: PATCH /brocade-zone/defined-configuration/zone
func (c *Client) UpdateZone(name string, members []string, principalMembers []string) error {
	payload := DefinedZoneAPI{
		Name:                name,
		ZoneTypeString:      "zone",
		MemberEntryNames:    members,
		PrincipalEntryNames: principalMembers,
	}
	return c.Patch("/brocade-zone/defined-configuration/zone", payload)
}

// RenameZone 重命名 Zone 定义配置中的一个 Zone。
// 对应 API: PATCH /brocade-zone/defined-configuration/zone/zone-name/{oldName}
func (c *Client) RenameZone(oldName, newName string) error {
	payload := DefinedZoneAPI{
		Name: newName,
	}
	endpoint := "/brocade-zone/defined-configuration/zone/zone-name/" + url.PathEscape(oldName)
	return c.Patch(endpoint, payload)
}

// ReplaceZoneAndActivate 执行完整的 Zone 替换并激活流程：
// 1. 获取当前 checksum
// 2. 更新 Zone 成员列表（覆盖原有成员）
// 3. 保存配置并激活
func (c *Client) ReplaceZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error {
	if err := validateZoneActivationInput(cfgName, zoneName, members); err != nil {
		return err
	}

	checksum, err := c.GetZoneChecksum()
	if err != nil {
		return err
	}
	if err := c.UpdateZone(zoneName, members, principalMembers); err != nil {
		return err
	}
	return c.saveAndActivateZoneConfig(cfgName, checksum)
}

// DeleteZone 从 Zone 定义配置中删除一个 Zone。
// 对应 API: DELETE /brocade-zone/defined-configuration/zone/zone-name/{name}
func (c *Client) DeleteZone(name string) error {
	endpoint := "/brocade-zone/defined-configuration/zone/zone-name/" + url.PathEscape(name)
	return c.Delete(endpoint)
}

// DeleteZoneAndActivate 执行完整的 Zone 删除并激活流程：
// 1. 获取当前 checksum
// 2. 删除 Zone
// 3. 保存配置并激活
func (c *Client) DeleteZoneAndActivate(cfgName, zoneName string) error {
	if strings.TrimSpace(cfgName) == "" {
		return errors.New("cfg name required")
	}
	if strings.TrimSpace(zoneName) == "" {
		return errors.New("zone name required")
	}

	checksum, err := c.GetZoneChecksum()
	if err != nil {
		return err
	}
	if err := c.DeleteZone(zoneName); err != nil {
		return err
	}
	return c.saveAndActivateZoneConfig(cfgName, checksum)
}

// saveAndActivateZoneConfig 保存 Zone 配置并使用新 checksum 激活指定的 cfg
func (c *Client) saveAndActivateZoneConfig(cfgName, checksum string) error {
	if err := c.SaveZoneConfig(checksum); err != nil {
		return err
	}
	newChecksum, err := c.GetZoneChecksum()
	if err != nil {
		return err
	}
	return c.ActivateZoneConfig(cfgName, newChecksum)
}

// validateZoneActivationInput 校验 Zone 激活操作所需的输入参数
func validateZoneActivationInput(cfgName, zoneName string, members []string) error {
	if strings.TrimSpace(cfgName) == "" {
		return errors.New("cfg name required")
	}
	if strings.TrimSpace(zoneName) == "" {
		return errors.New("zone name required")
	}
	if len(members) == 0 {
		return errors.New("zone members required")
	}
	return nil
}

// memberZonesForConfig 从已定义的 cfg 列表中查找指定配置名称的成员 Zone 列表
func memberZonesForConfig(configs []ConfigInfo, cfgName string) ([]string, error) {
	for _, cfg := range configs {
		if cfg.Name == cfgName {
			return append([]string(nil), cfg.MemberZones...), nil
		}
	}
	return nil, ErrNotFound
}

// containsString 判断字符串切片中是否包含目标字符串
func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
