package san

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
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
	return c.GetDefinedZonesWithContext(context.Background())
}

// GetDefinedZonesWithContext 获取 Zone 定义配置中的所有 Zone 列表。
func (c *Client) GetDefinedZonesWithContext(ctx context.Context) ([]ZoneInfo, error) {
	var resp DefinedZoneResponse
	err := c.GetWithContext(ctx, c.endpoints().DefinedZones(), &resp)
	if err != nil {
		return nil, err
	}

	zones := make([]ZoneInfo, 0, len(resp.Zones))
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
	return c.GetDefinedZoneWithContext(context.Background(), name)
}

// GetDefinedZoneWithContext 获取 Zone 定义配置中的单个 Zone。
func (c *Client) GetDefinedZoneWithContext(ctx context.Context, name string) (*ZoneInfo, error) {
	var resp DefinedZoneResponse
	if err := c.GetWithContext(ctx, c.endpoints().DefinedZone(name), &resp); err != nil {
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
	return c.GetEffectiveZonesWithContext(context.Background())
}

// GetEffectiveZonesWithContext 获取已生效配置中的所有 Zone 列表。
func (c *Client) GetEffectiveZonesWithContext(ctx context.Context) ([]ZoneInfo, error) {
	var resp EffectiveZoneResponse
	err := c.GetWithContext(ctx, c.endpoints().EffectiveZones(), &resp)
	if err != nil {
		return nil, err
	}

	zones := make([]ZoneInfo, 0, len(resp.Zones))
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
	return c.CreateZoneWithContext(context.Background(), name, members, principalMembers)
}

// CreateZoneWithContext 在 Zone 定义配置中创建一个新的 Zone。
func (c *Client) CreateZoneWithContext(ctx context.Context, name string, members []string, principalMembers []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("zone name required")
	}
	if err := c.ensureZoneAbsentWithContext(ctx, name); err != nil {
		return err
	}
	return c.PostWithContext(ctx, c.endpoints().DefinedZones(), definedZonePayload(name, members, principalMembers))
}

// CreateZoneAndActivate 执行完整的 Zone 创建并激活流程：
// 1. 获取当前 checksum
// 2. 创建 Zone
// 3. 将 Zone 添加到指定 cfg 配置中
// 4. 保存配置并激活
// 创建前要求 Zone 不存在；若 cfg 已包含该名称则不会重复添加。
func (c *Client) CreateZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error {
	return c.CreateZoneAndActivateWithContext(context.Background(), cfgName, zoneName, members, principalMembers)
}

// CreateZoneAndActivateWithContext 执行带 context 的 Zone 创建并激活流程。
func (c *Client) CreateZoneAndActivateWithContext(ctx context.Context, cfgName, zoneName string, members []string, principalMembers []string) (err error) {
	ctx = nonNilContext(ctx)
	if err := validateZoneActivationInput(cfgName, zoneName, members); err != nil {
		return err
	}

	c.zoneWriteMu.Lock()
	defer c.zoneWriteMu.Unlock()

	configs, err := c.GetDefinedConfigsWithContext(ctx)
	if err != nil {
		return err
	}
	memberZones, err := memberZonesForConfig(configs, cfgName)
	if err != nil {
		return err
	}
	if err := c.ensureZoneAbsentWithContext(ctx, zoneName); err != nil {
		return err
	}
	checksum, err := c.GetZoneChecksumWithContext(ctx)
	if err != nil {
		return err
	}
	if err = c.ensureWriteSupported(); err != nil {
		return err
	}

	mutationAttempted := false
	defer func() {
		if err == nil || !mutationAttempted {
			return
		}
		if abortErr := c.abortZoneTransactionAfterFailure(); abortErr != nil {
			err = errors.Join(err, fmt.Errorf("abort zone transaction: %w", abortErr))
		}
		err = &PartialMutationError{Err: err}
	}()

	mutationAttempted = true
	if err = c.PostWithContext(ctx, c.endpoints().DefinedZones(), definedZonePayload(zoneName, members, principalMembers)); err != nil {
		return err
	}
	if !containsString(memberZones, zoneName) {
		memberZones = append(memberZones, zoneName)
	}
	if err = c.UpdateDefinedConfigWithContext(ctx, cfgName, memberZones); err != nil {
		return err
	}
	return c.saveAndActivateZoneConfigWithContext(ctx, cfgName, checksum)
}

// UpdateZone 更新 Zone 定义配置中已有 Zone 的成员列表。
// members 为普通成员，principalMembers 为 Principal 成员。
// 对应 API: PATCH /brocade-zone/defined-configuration/zone
func (c *Client) UpdateZone(name string, members []string, principalMembers []string) error {
	return c.UpdateZoneWithContext(context.Background(), name, members, principalMembers)
}

// UpdateZoneWithContext 更新 Zone 定义配置中已有 Zone 的成员列表。
func (c *Client) UpdateZoneWithContext(ctx context.Context, name string, members []string, principalMembers []string) error {
	zoneTypeString, err := c.zoneTypeStringForUpdateWithContext(ctx, name, principalMembers)
	if err != nil {
		return err
	}
	return c.patchZoneWithContext(ctx, zoneUpdatePayload(name, zoneTypeString, members, principalMembers))
}

// RenameZone 重命名 Zone 定义配置中的一个 Zone。
// 对应 API: PATCH /brocade-zone/defined-configuration/zone/zone-name/{oldName}
func (c *Client) RenameZone(oldName, newName string) error {
	return c.RenameZoneWithContext(context.Background(), oldName, newName)
}

// RenameZoneWithContext 重命名 Zone 定义配置中的一个 Zone。
func (c *Client) RenameZoneWithContext(ctx context.Context, oldName, newName string) error {
	payload := DefinedZoneAPI{
		Name: newName,
	}
	return c.PatchWithContext(ctx, c.endpoints().DefinedZone(oldName), payload)
}

// ReplaceZoneAndActivate 执行完整的 Zone 替换并激活流程：
// 1. 获取当前 checksum
// 2. 更新 Zone 成员列表（覆盖原有成员）
// 3. 保存配置并激活
func (c *Client) ReplaceZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error {
	return c.ReplaceZoneAndActivateWithContext(context.Background(), cfgName, zoneName, members, principalMembers)
}

// ReplaceZoneAndActivateWithContext 执行带 context 的 Zone 替换并激活流程。
func (c *Client) ReplaceZoneAndActivateWithContext(ctx context.Context, cfgName, zoneName string, members []string, principalMembers []string) (err error) {
	ctx = nonNilContext(ctx)
	if err := validateZoneActivationInput(cfgName, zoneName, members); err != nil {
		return err
	}

	c.zoneWriteMu.Lock()
	defer c.zoneWriteMu.Unlock()

	configs, err := c.GetDefinedConfigsWithContext(ctx)
	if err != nil {
		return err
	}
	if _, err := memberZonesForConfig(configs, cfgName); err != nil {
		return err
	}
	checksum, err := c.GetZoneChecksumWithContext(ctx)
	if err != nil {
		return err
	}
	zoneTypeString, err := c.zoneTypeStringForUpdateWithContext(ctx, zoneName, principalMembers)
	if err != nil {
		return err
	}
	if err = c.ensureWriteSupported(); err != nil {
		return err
	}

	mutationAttempted := false
	defer func() {
		if err == nil || !mutationAttempted {
			return
		}
		if abortErr := c.abortZoneTransactionAfterFailure(); abortErr != nil {
			err = errors.Join(err, fmt.Errorf("abort zone transaction: %w", abortErr))
		}
		err = &PartialMutationError{Err: err}
	}()

	mutationAttempted = true
	if err = c.patchZoneWithContext(ctx, zoneUpdatePayload(zoneName, zoneTypeString, members, principalMembers)); err != nil {
		return err
	}
	return c.saveAndActivateZoneConfigWithContext(ctx, cfgName, checksum)
}

// DeleteZone 从 Zone 定义配置中删除一个 Zone。
// 对应 API: DELETE /brocade-zone/defined-configuration/zone/zone-name/{name}
func (c *Client) DeleteZone(name string) error {
	return c.DeleteZoneWithContext(context.Background(), name)
}

// DeleteZoneWithContext 从 Zone 定义配置中删除一个 Zone。
func (c *Client) DeleteZoneWithContext(ctx context.Context, name string) error {
	if _, err := c.GetDefinedZoneWithContext(ctx, name); err != nil {
		return err
	}
	return c.DeleteWithContext(ctx, c.endpoints().DefinedZone(name))
}

// DeleteZoneAndActivate 执行完整的 Zone 删除并激活流程：
// 1. 获取当前 checksum
// 2. 从所有包含该 Zone 的 cfg 中移除 Zone
// 3. 删除 Zone
// 4. 保存配置并激活
func (c *Client) DeleteZoneAndActivate(cfgName, zoneName string) error {
	return c.DeleteZoneAndActivateWithContext(context.Background(), cfgName, zoneName)
}

// DeleteZoneAndActivateWithContext 执行带 context 的 Zone 删除并激活流程。
func (c *Client) DeleteZoneAndActivateWithContext(ctx context.Context, cfgName, zoneName string) (err error) {
	ctx = nonNilContext(ctx)
	if strings.TrimSpace(cfgName) == "" {
		return errors.New("cfg name required")
	}
	if strings.TrimSpace(zoneName) == "" {
		return errors.New("zone name required")
	}

	c.zoneWriteMu.Lock()
	defer c.zoneWriteMu.Unlock()

	configs, err := c.GetDefinedConfigsWithContext(ctx)
	if err != nil {
		return err
	}
	if _, err := memberZonesForConfig(configs, cfgName); err != nil {
		return err
	}
	if _, err := c.GetDefinedZoneWithContext(ctx, zoneName); err != nil {
		return err
	}
	checksum, err := c.GetZoneChecksumWithContext(ctx)
	if err != nil {
		return err
	}
	if err = c.ensureWriteSupported(); err != nil {
		return err
	}

	mutationAttempted := false
	defer func() {
		if err == nil || !mutationAttempted {
			return
		}
		if abortErr := c.abortZoneTransactionAfterFailure(); abortErr != nil {
			err = errors.Join(err, fmt.Errorf("abort zone transaction: %w", abortErr))
		}
		err = &PartialMutationError{Err: err}
	}()

	mutationAttempted = true
	for _, cfg := range configs {
		memberZones, removed := removeString(cfg.MemberZones, zoneName)
		if !removed {
			continue
		}
		if err = c.UpdateDefinedConfigWithContext(ctx, cfg.Name, memberZones); err != nil {
			return err
		}
	}
	if err = c.DeleteWithContext(ctx, c.endpoints().DefinedZone(zoneName)); err != nil {
		return err
	}
	return c.saveAndActivateZoneConfigWithContext(ctx, cfgName, checksum)
}

// saveAndActivateZoneConfig 保存 Zone 配置并使用新 checksum 激活指定的 cfg
func (c *Client) saveAndActivateZoneConfig(cfgName, checksum string) error {
	return c.saveAndActivateZoneConfigWithContext(context.Background(), cfgName, checksum)
}

func (c *Client) saveAndActivateZoneConfigWithContext(ctx context.Context, cfgName, checksum string) error {
	if err := c.SaveZoneConfigWithContext(ctx, checksum); err != nil {
		return err
	}
	newChecksum, err := c.GetZoneChecksumWithContext(ctx)
	if err != nil {
		return err
	}
	return c.ActivateZoneConfigWithContext(ctx, cfgName, newChecksum)
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
	return c.ensureZoneAbsentWithContext(context.Background(), name)
}

func (c *Client) ensureZoneAbsentWithContext(ctx context.Context, name string) error {
	_, err := c.GetDefinedZoneWithContext(ctx, name)
	if err == nil {
		return errors.New("zone already exists")
	}
	if isNotFoundError(err) {
		return nil
	}
	return err
}

func (c *Client) zoneTypeStringForUpdate(name string, principalMembers []string) (string, error) {
	return c.zoneTypeStringForUpdateWithContext(context.Background(), name, principalMembers)
}

func (c *Client) zoneTypeStringForUpdateWithContext(ctx context.Context, name string, principalMembers []string) (string, error) {
	zone, err := c.GetDefinedZoneWithContext(ctx, name)
	if err != nil {
		return "", err
	}
	if zone.TypeString == ZoneTypeZone && len(principalMembers) > 0 {
		return "", errors.New("cannot convert zone to peer zone")
	}
	if zone.TypeString == ZoneTypeUserCreatedPeerZone && len(principalMembers) == 0 {
		return "", errors.New("peer zone requires principal members")
	}
	if zone.TypeString != ZoneTypeZone && zone.TypeString != ZoneTypeUserCreatedPeerZone {
		return "", fmt.Errorf("%w: unsupported zone type %q", ErrInvalidResponse, zone.TypeString)
	}
	return zone.TypeString, nil
}

func definedZonePayload(name string, members, principalMembers []string) DefinedZoneAPI {
	return DefinedZoneAPI{
		Name:                name,
		ZoneTypeString:      zoneTypeStringForCreate(principalMembers),
		MemberEntryNames:    members,
		PrincipalEntryNames: principalMembers,
	}
}

func zoneUpdatePayload(name, zoneTypeString string, members, principalMembers []string) DefinedZoneAPI {
	return DefinedZoneAPI{
		Name:                name,
		ZoneTypeString:      zoneTypeString,
		MemberEntryNames:    members,
		PrincipalEntryNames: principalMembers,
	}
}

func (c *Client) patchZoneWithContext(ctx context.Context, payload DefinedZoneAPI) error {
	return c.PatchWithContext(ctx, c.endpoints().DefinedZones(), payload)
}

func (c *Client) abortZoneTransactionAfterFailure() error {
	timeout := min(c.Timeout(), zoneCleanupTimeout)
	if timeout <= 0 {
		timeout = zoneCleanupTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.AbortZoneTransactionWithContext(ctx)
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
