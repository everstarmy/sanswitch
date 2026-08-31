package sanswitch

// Switch 表示单台交换机的摘要信息，由 FabricSwitch 中的原始数据提取。
type Switch struct {
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

// ZoneMember contains the ordinary and principal members of a Zone.
type ZoneMember struct {
	MemberEntries    []string `json:"member"`
	PrincipalEntries []string `json:"principal"`
}

// Zone 表示一个 Zone（已定义或已生效），包含名称、成员列表和类型
type Zone struct {
	Name    string     `json:"name"`
	Members ZoneMember `json:"members"`

	Description string `json:"description"`
	Type        string `json:"type"`
	TypeString  string `json:"type_string"`
}

// Alias 表示一个 Zone Alias，包含名称和成员列表
type Alias struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

// ZoneConfig 表示一个 Zone 配置（cfg），包含配置名称、类型、成员 Zone 列表和校验和
type ZoneConfig struct {
	Name              string   `json:"name"`
	Type              string   `json:"type"`
	MemberZones       []string `json:"member_zones"`
	Checksum          string   `json:"checksum"`
	DefaultZoneAccess string   `json:"default_zone_access"`
}

// ZoneDatabase 表示 Zone 数据库的容量和事务状态信息
type ZoneDatabase struct {
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
