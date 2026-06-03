package san

// SwitchInfo 表示单台交换机的摘要信息，由 GetSwitchInfo 从 FabricSwitch 中提取
type SwitchInfo struct {
	Name            string `json:"name"`
	WWN             string `json:"wwn"`
	ChassisWWN      string `json:"chassis_wwn"`
	DomainID        int    `json:"domain_id"`
	FirmwareVersion string `json:"firmware_version"`
	ModelName       string `json:"model_name"`
	SerialNumber    string `json:"serial_number"`
	IPAddress       string `json:"ip_address"`
	IPv6Address     string `json:"ipv6_address"`
	Fcid            string `json:"fcid"`
	FcidHex         string `json:"fcid_hex"`
	Principal       bool   `json:"principal"`
}

type ZoneMember struct {
	MemberEntries    []string `json:"member"`
	PrincipalEntries []string `json:"principal"`
}

// ZoneInfo 表示一个 Zone（已定义或已生效），包含名称、成员列表和类型
type ZoneInfo struct {
	Name    string     `json:"name"`
	Members ZoneMember `json:"members"`

	Description string `json:"description"`
	Type        string `json:"type"`
	TypeString  string `json:"type_string"`
}

// AliasInfo 表示一个 Zone Alias，包含名称和成员列表
type AliasInfo struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

// ConfigInfo 表示一个 Zone 配置（cfg），包含配置名称、类型、成员 Zone 列表和校验和
type ConfigInfo struct {
	Name              string   `json:"name"`
	Type              string   `json:"type"`
	MemberZones       []string `json:"member_zones"`
	Checksum          string   `json:"checksum"`
	DefaultZoneAccess string   `json:"default_zone_access"`
}

// ZoneDatabaseInfo 表示 Zone 数据库的容量和事务状态信息
type ZoneDatabaseInfo struct {
	DBMax                  uint32 `json:"db_max"`
	DBAvail                uint32 `json:"db_avail"`
	DBCommitted            uint32 `json:"db_committed"`
	DBTransaction          uint32 `json:"db_transaction"`
	TransactionToken       uint32 `json:"transaction_token"`
	DBChassisWideCommitted uint32 `json:"db_chassis_wide_committed"`
	DBChassisWideMax       uint32 `json:"db_chassis_wide_max"`
	DBFabricWideMax        uint32 `json:"db_fabric_wide_max"`
	DomainWithLowestDBMax  uint32 `json:"domain_with_lowest_db_max"`
}
