package sanswitch

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

const (
	// ZoneTypeZone 表示普通 user-created zone 的 zone-type-string。
	ZoneTypeZone = "zone"
	// ZoneTypeUserCreatedPeerZone 表示 peer zone 的 zone-type-string。
	ZoneTypeUserCreatedPeerZone = "user-created-peer-zone"
)

// definedZoneAPI 表示 Zone 定义配置中的 Zone（用于 XML 请求/响应序列化）。
// MemberEntryNames 为普通成员（entry-name），PrincipalEntryNames 为 Principal 成员（principal-entry-name）。
type definedZoneAPI struct {
	XMLName             xml.Name `xml:"zone"`
	Name                string   `xml:"zone-name"`
	ZoneType            string   `xml:"zone-type,omitempty"`
	ZoneTypeString      string   `xml:"zone-type-string"`
	MemberEntryNames    []string `xml:"member-entry>entry-name"`           // 普通成员列表
	PrincipalEntryNames []string `xml:"member-entry>principal-entry-name"` // Principal 成员列表
}

// definedZoneResponse 是 GET /brocade-zone/defined-configuration/zone 的 XML 响应包装
type definedZoneResponse struct {
	XMLName xml.Name         `xml:"Response"`
	Zones   []definedZoneAPI `xml:"zone"`
}

// effectiveZoneAPI 表示已生效配置中的 Zone（用于 XML 响应序列化）。
// 与 definedZoneAPI 类似，但对应 YANG 模型中的 enabled-zone 节点。
type effectiveZoneAPI struct {
	XMLName             xml.Name `xml:"enabled-zone"`
	Name                string   `xml:"zone-name"`
	ZoneType            string   `xml:"zone-type"`
	ZoneTypeString      string   `xml:"zone-type-string"`
	MemberEntryNames    []string `xml:"member-entry>entry-name"`           // 普通成员列表
	PrincipalEntryNames []string `xml:"member-entry>principal-entry-name"` // Principal 成员列表
}

// effectiveZoneResponse 是 GET /brocade-zone/effective-configuration/enabled-zone 的 XML 响应包装
type effectiveZoneResponse struct {
	XMLName xml.Name           `xml:"Response"`
	Zones   []effectiveZoneAPI `xml:"effective-configuration>enabled-zone"`
}

// DefinedZones 获取 Zone 定义配置中的所有 Zone 列表。
func (c *Client) DefinedZones(ctx context.Context) ([]Zone, error) {
	var resp definedZoneResponse
	err := c.get(ctx, c.endpoints().DefinedZones(), &resp)
	if err != nil {
		return nil, err
	}

	zones := make([]Zone, 0, len(resp.Zones))
	for _, z := range resp.Zones {
		zones = append(zones, Zone{
			Name:        z.Name,
			Members:     ZoneMember{MemberEntries: z.MemberEntryNames, PrincipalEntries: z.PrincipalEntryNames},
			Description: z.ZoneTypeString,
			Type:        z.ZoneType,
			TypeString:  z.ZoneTypeString,
		})
	}

	return zones, nil
}

// DefinedZone 获取 Zone 定义配置中的单个 Zone。
func (c *Client) DefinedZone(ctx context.Context, name string) (*Zone, error) {
	var resp definedZoneResponse
	if err := c.get(ctx, c.endpoints().DefinedZone(name), &resp); err != nil {
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

// EffectiveZones 获取已生效配置中的所有 Zone 列表。
func (c *Client) EffectiveZones(ctx context.Context) ([]Zone, error) {
	var resp effectiveZoneResponse
	err := c.get(ctx, c.endpoints().EffectiveZones(), &resp)
	if err != nil {
		return nil, err
	}

	zones := make([]Zone, 0, len(resp.Zones))
	for _, z := range resp.Zones {
		zones = append(zones, Zone{
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
func (c *Client) CreateZone(ctx context.Context, name string, members []string, principalMembers []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("zone name required")
	}
	if len(members) == 0 && len(principalMembers) == 0 {
		return errors.New("zone members required")
	}
	if err := c.ensureZoneAbsent(ctx, name); err != nil {
		return err
	}
	return c.post(ctx, c.endpoints().DefinedZones(), definedZonePayload(name, members, principalMembers))
}

// UpdateZone 更新 Zone 定义配置中已有 Zone 的成员列表。
func (c *Client) UpdateZone(ctx context.Context, name string, members []string, principalMembers []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("zone name required")
	}
	if len(members) == 0 && len(principalMembers) == 0 {
		return errors.New("zone members required")
	}
	zoneTypeString, err := c.zoneTypeStringForUpdate(ctx, name, principalMembers)
	if err != nil {
		return err
	}
	return c.patchZone(ctx, zoneUpdatePayload(name, zoneTypeString, members, principalMembers))
}

// RenameZone 重命名 Zone 定义配置中的一个 Zone。
func (c *Client) RenameZone(ctx context.Context, oldName, newName string) error {
	if strings.TrimSpace(oldName) == "" || strings.TrimSpace(newName) == "" {
		return errors.New("old and new zone names required")
	}
	payload := definedZoneAPI{
		Name: newName,
	}
	return c.patch(ctx, c.endpoints().DefinedZone(oldName), payload)
}

// DeleteZone 从 Zone 定义配置中删除一个 Zone。
func (c *Client) DeleteZone(ctx context.Context, name string) error {
	if _, err := c.DefinedZone(ctx, name); err != nil {
		return err
	}
	return c.delete(ctx, c.endpoints().DefinedZone(name))
}

func zoneTypeStringForCreate(principalMembers []string) string {
	if len(principalMembers) > 0 {
		return ZoneTypeUserCreatedPeerZone
	}
	return ZoneTypeZone
}

func (c *Client) ensureZoneAbsent(ctx context.Context, name string) error {
	_, err := c.DefinedZone(ctx, name)
	if err == nil {
		return errors.New("zone already exists")
	}
	if isNotFoundError(err) {
		return nil
	}
	return err
}

func (c *Client) zoneTypeStringForUpdate(ctx context.Context, name string, principalMembers []string) (string, error) {
	zone, err := c.DefinedZone(ctx, name)
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

func definedZonePayload(name string, members, principalMembers []string) definedZoneAPI {
	return definedZoneAPI{
		Name:                name,
		ZoneTypeString:      zoneTypeStringForCreate(principalMembers),
		MemberEntryNames:    members,
		PrincipalEntryNames: principalMembers,
	}
}

func zoneUpdatePayload(name, zoneTypeString string, members, principalMembers []string) definedZoneAPI {
	return definedZoneAPI{
		Name:                name,
		ZoneTypeString:      zoneTypeString,
		MemberEntryNames:    members,
		PrincipalEntryNames: principalMembers,
	}
}

func (c *Client) patchZone(ctx context.Context, payload definedZoneAPI) error {
	return c.patch(ctx, c.endpoints().DefinedZones(), payload)
}

func zoneInfoFromDefinedZone(z definedZoneAPI) Zone {
	return Zone{
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
func memberZonesForConfig(configs []ZoneConfig, cfgName string) ([]string, error) {
	for _, cfg := range configs {
		if cfg.Name == cfgName {
			return slices.Clone(cfg.MemberZones), nil
		}
	}
	return nil, ErrNotFound
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
