package san

import (
	"encoding/xml"
	"errors"
	"net/http"
	"strings"
)

const (
	// ZoneTypeZone 表示普通 user-created zone 的 zone-type-string。
	ZoneTypeZone = "zone"
	// ZoneTypeUserCreatedPeerZone 表示 peer zone 的 zone-type-string。
	ZoneTypeUserCreatedPeerZone = "user-created-peer-zone"
)

// DefinedZoneAPI 表示 Zone 定义配置中的 Zone（用于 XML 请求/响应序列化）。
// MemberEntryNames 为普通成员（entry-name），PrincipalEntryNames 为 Principal 成员（principal-entry-name）。
type DefinedZoneAPI struct {
	XMLName             xml.Name `xml:"zone"`
	Name                string   `xml:"zone-name"`
	ZoneType            string   `xml:"zone-type,omitempty"`
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
	Zones   []EffectiveZoneAPI `xml:"effective-configuration>enabled-zone"`
}

// GetDefinedZones 获取 Zone 定义配置中的所有 Zone 列表。
// 返回的 ZoneInfo 中 Members 字段包含普通成员和 Principal 成员的合并列表。
// 对应 API: GET /brocade-zone/defined-configuration/zone
func (c *Client) GetDefinedZones() ([]ZoneInfo, error) {
	var resp DefinedZoneResponse
	err := c.Get(c.endpoints().DefinedZones(), &resp)
	if err != nil {
		return nil, err
	}

	var zones []ZoneInfo
	for _, z := range resp.Zones {
		zones = append(zones, ZoneInfo{
			Name:        z.Name,
			Members:     ZoneMember{MemberEntries: z.MemberEntryNames, PrincipalEntries: z.PrincipalEntryNames},
			Description: z.ZoneTypeString,
			Type:        z.ZoneType,
			TypeString:  z.ZoneTypeString,
		})
	}

	return zones, nil
}

// GetDefinedZone 获取 Zone 定义配置中的单个 Zone。
// 对应 API: GET /brocade-zone/defined-configuration/zone/zone-name/{name}
func (c *Client) GetDefinedZone(name string) (*ZoneInfo, error) {
	var resp DefinedZoneResponse
	if err := c.Get(c.endpoints().DefinedZone(name), &resp); err != nil {
		if isNotFoundError(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if len(resp.Zones) == 0 {
		return nil, ErrNotFound
	}
	zone := zoneInfoFromDefinedZone(resp.Zones[0])
	return &zone, nil
}

// GetEffectiveZones 获取已生效配置中的所有 Zone 列表。
// 返回的 ZoneInfo 中 Members 字段包含普通成员和 Principal 成员的合并列表。
// 对应 API: GET /brocade-zone/effective-configuration/enabled-zone
func (c *Client) GetEffectiveZones() ([]ZoneInfo, error) {
	var resp EffectiveZoneResponse
	err := c.Get(c.endpoints().EffectiveZones(), &resp)
	if err != nil {
		return nil, err
	}

	var zones []ZoneInfo
	for _, z := range resp.Zones {
		zones = append(zones, ZoneInfo{
			Name:        z.Name,
			Members:     ZoneMember{MemberEntries: z.MemberEntryNames, PrincipalEntries: z.PrincipalEntryNames},
			Description: z.ZoneTypeString,
			Type:        z.ZoneType,
			TypeString:  z.ZoneTypeString,
		})
	}

	return zones, nil
}

// CreateZone 在 Zone 定义配置中创建一个新的 Zone。
// members 为普通成员（entry-name），principalMembers 为 Principal 成员（principal-entry-name）。
// 对应 API: POST /brocade-zone/defined-configuration/zone
func (c *Client) CreateZone(name string, members []string, principalMembers []string) error {
	if err := c.ensureZoneAbsent(name); err != nil {
		return err
	}
	payload := DefinedZoneAPI{
		Name:                name,
		ZoneTypeString:      zoneTypeStringForCreate(principalMembers),
		MemberEntryNames:    members,
		PrincipalEntryNames: principalMembers,
	}
	return c.Post(c.endpoints().DefinedZones(), payload)
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
	zoneTypeString, err := c.zoneTypeStringForUpdate(name, principalMembers)
	if err != nil {
		return err
	}
	payload := DefinedZoneAPI{
		Name:                name,
		ZoneTypeString:      zoneTypeString,
		MemberEntryNames:    members,
		PrincipalEntryNames: principalMembers,
	}
	return c.Patch(c.endpoints().DefinedZones(), payload)
}

// RenameZone 重命名 Zone 定义配置中的一个 Zone。
// 对应 API: PATCH /brocade-zone/defined-configuration/zone/zone-name/{oldName}
func (c *Client) RenameZone(oldName, newName string) error {
	payload := DefinedZoneAPI{
		Name: newName,
	}
	return c.Patch(c.endpoints().DefinedZone(oldName), payload)
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
	if _, err := c.GetDefinedZone(name); err != nil {
		return err
	}
	return c.Delete(c.endpoints().DefinedZone(name))
}

// DeleteZoneAndActivate 执行完整的 Zone 删除并激活流程：
// 1. 获取当前 checksum
// 2. 从所有包含该 Zone 的 cfg 中移除 Zone
// 3. 删除 Zone
// 4. 保存配置并激活
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
	if _, err := c.GetDefinedZone(zoneName); err != nil {
		return err
	}
	configs, err := c.GetDefinedConfigs()
	if err != nil {
		return err
	}
	for _, cfg := range configs {
		memberZones, removed := removeString(cfg.MemberZones, zoneName)
		if !removed {
			continue
		}
		if err := c.UpdateDefinedConfig(cfg.Name, memberZones); err != nil {
			return err
		}
	}
	if err := c.Delete(c.endpoints().DefinedZone(zoneName)); err != nil {
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

func zoneTypeStringForCreate(principalMembers []string) string {
	if len(principalMembers) > 0 {
		return ZoneTypeUserCreatedPeerZone
	}
	return ZoneTypeZone
}

func (c *Client) ensureZoneAbsent(name string) error {
	_, err := c.GetDefinedZone(name)
	if err == nil {
		return errors.New("zone already exists")
	}
	if isNotFoundError(err) {
		return nil
	}
	return err
}

func (c *Client) zoneTypeStringForUpdate(name string, principalMembers []string) (string, error) {
	zone, err := c.GetDefinedZone(name)
	if err != nil {
		return "", err
	}
	if zone.TypeString == ZoneTypeZone && len(principalMembers) > 0 {
		return "", errors.New("cannot convert zone to peer zone")
	}
	return zone.TypeString, nil
}

func zoneInfoFromDefinedZone(z DefinedZoneAPI) ZoneInfo {
	return ZoneInfo{
		Name:        z.Name,
		Members:     ZoneMember{MemberEntries: z.MemberEntryNames, PrincipalEntries: z.PrincipalEntryNames},
		Description: z.ZoneTypeString,
		Type:        z.ZoneType,
		TypeString:  z.ZoneTypeString,
	}
}

func isNotFoundError(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
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

func removeString(values []string, target string) ([]string, bool) {
	filtered := make([]string, 0, len(values))
	removed := false
	for _, value := range values {
		if value == target {
			removed = true
			continue
		}
		filtered = append(filtered, value)
	}
	return filtered, removed
}
