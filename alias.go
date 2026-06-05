package san

import (
	"encoding/xml"
)

// DefinedAliasAPI 表示 Zone 定义配置中的 Alias（用于 XML 请求/响应序列化）
type DefinedAliasAPI struct {
	XMLName          xml.Name `xml:"alias"`
	Name             string   `xml:"alias-name"`
	MemberEntryNames []string `xml:"member-entry>alias-entry-name"` // Alias 成员列表（WWN 或别名）
}

// DefinedAliasResponse 是 GET /brocade-zone/defined-configuration/alias 的 XML 响应包装
type DefinedAliasResponse struct {
	XMLName xml.Name          `xml:"Response"`
	Aliases []DefinedAliasAPI `xml:"alias"`
}

// GetDefinedAliases 获取 Zone 定义配置中的所有 Alias 列表。
// 对应 API: GET /brocade-zone/defined-configuration/alias
func (c *Client) GetDefinedAliases() ([]AliasInfo, error) {
	var resp DefinedAliasResponse
	err := c.Get(c.endpoints().DefinedAliases(), &resp)
	if err != nil {
		return nil, err
	}

	var aliases []AliasInfo
	for _, a := range resp.Aliases {
		aliases = append(aliases, AliasInfo{
			Name:    a.Name,
			Members: a.MemberEntryNames,
		})
	}

	return aliases, nil
}

// CreateAlias 在 Zone 定义配置中创建一个新的 Alias。
// 对应 API: POST /brocade-zone/defined-configuration/alias
func (c *Client) CreateAlias(name string, members []string) error {
	payload := DefinedAliasAPI{
		Name:             name,
		MemberEntryNames: members,
	}
	return c.Post(c.endpoints().DefinedAliases(), payload)
}

// UpdateAlias 更新 Zone 定义配置中已有 Alias 的成员列表。
// 对应 API: PATCH /brocade-zone/defined-configuration/alias
func (c *Client) UpdateAlias(name string, members []string) error {
	payload := DefinedAliasAPI{
		Name:             name,
		MemberEntryNames: members,
	}
	return c.Patch(c.endpoints().DefinedAliases(), payload)
}

// RenameAlias 重命名 Zone 定义配置中的一个 Alias。
// 对应 API: PATCH /brocade-zone/defined-configuration/alias/alias-name/{oldName}
func (c *Client) RenameAlias(oldName, newName string) error {
	payload := DefinedAliasAPI{
		Name: newName,
	}
	return c.Patch(c.endpoints().DefinedAlias(oldName), payload)
}

// DeleteAlias 从 Zone 定义配置中删除一个 Alias。
// 对应 API: DELETE /brocade-zone/defined-configuration/alias/alias-name/{name}
func (c *Client) DeleteAlias(name string) error {
	return c.Delete(c.endpoints().DefinedAlias(name))
}
