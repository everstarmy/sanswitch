package sanswitch

import (
	"context"
	"encoding/xml"
	"errors"
	"strings"
)

// definedAliasAPI 表示 Zone 定义配置中的 Alias（用于 XML 请求/响应序列化）
type definedAliasAPI struct {
	XMLName          xml.Name `xml:"alias"`
	Name             string   `xml:"alias-name"`
	MemberEntryNames []string `xml:"member-entry>alias-entry-name"` // Alias 成员列表（WWN 或别名）
}

// definedAliasResponse 是 GET /brocade-zone/defined-configuration/alias 的 XML 响应包装
type definedAliasResponse struct {
	XMLName xml.Name          `xml:"Response"`
	Aliases []definedAliasAPI `xml:"alias"`
}

// DefinedAliases 获取 Zone 定义配置中的所有 Alias 列表并允许取消请求。
func (c *Client) DefinedAliases(ctx context.Context) ([]Alias, error) {
	var resp definedAliasResponse
	err := c.get(ctx, c.endpoints().DefinedAliases(), &resp)
	if err != nil {
		return nil, err
	}

	var aliases []Alias
	for _, a := range resp.Aliases {
		aliases = append(aliases, Alias{
			Name:    a.Name,
			Members: a.MemberEntryNames,
		})
	}

	return aliases, nil
}

// CreateAlias 在 Zone 定义配置中创建 Alias 并允许取消请求。
func (c *Client) CreateAlias(ctx context.Context, name string, members []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("alias name required")
	}
	if len(members) == 0 {
		return errors.New("alias members required")
	}
	payload := definedAliasAPI{
		Name:             name,
		MemberEntryNames: members,
	}
	return c.post(ctx, c.endpoints().DefinedAliases(), payload)
}

// UpdateAlias 更新 Alias 成员列表并允许取消请求。
func (c *Client) UpdateAlias(ctx context.Context, name string, members []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("alias name required")
	}
	if len(members) == 0 {
		return errors.New("alias members required")
	}
	payload := definedAliasAPI{
		Name:             name,
		MemberEntryNames: members,
	}
	return c.patch(ctx, c.endpoints().DefinedAliases(), payload)
}

// RenameAlias 重命名 Alias 并允许取消请求。
func (c *Client) RenameAlias(ctx context.Context, oldName, newName string) error {
	if strings.TrimSpace(oldName) == "" || strings.TrimSpace(newName) == "" {
		return errors.New("old and new alias names required")
	}
	payload := definedAliasAPI{
		Name: newName,
	}
	return c.patch(ctx, c.endpoints().DefinedAlias(oldName), payload)
}

// DeleteAlias 删除 Alias 并允许取消请求。
func (c *Client) DeleteAlias(ctx context.Context, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("alias name required")
	}
	return c.delete(ctx, c.endpoints().DefinedAlias(name))
}
