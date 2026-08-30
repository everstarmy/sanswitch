package san

import "context"

// Session describes the lifecycle operations supported by SANSwitch.
// All network operations accept a context so callers can bound their lifetime.
type Session interface {
	LoginWithContext(context.Context) (*LoginResponse, error)
	LogoutWithContext(context.Context) error
	IsLoggedIn() bool
	Close() error
}

// SwitchReader describes read-only switch, port, chassis, and logical-switch
// operations.
type SwitchReader interface {
	GetSwitchInfoWithContext(context.Context) (*SwitchInfo, error)
	GetFabricSwitchesWithContext(context.Context) ([]FabricSwitch, error)
	GetPortsWithContext(context.Context) ([]PortInfo, error)
	GetHardwareInfoWithContext(context.Context) (*HardwareInfo, error)
	GetLogicalSwitchesWithContext(context.Context) ([]LogicalSwitchInfo, error)
}

// ZoneReader describes read-only Zone database operations.
type ZoneReader interface {
	GetDefinedZonesWithContext(context.Context) ([]ZoneInfo, error)
	GetDefinedZoneWithContext(context.Context, string) (*ZoneInfo, error)
	GetEffectiveZonesWithContext(context.Context) ([]ZoneInfo, error)
	GetDefinedAliasesWithContext(context.Context) ([]AliasInfo, error)
	GetDefinedConfigsWithContext(context.Context) ([]ConfigInfo, error)
	GetEffectiveConfigWithContext(context.Context) (*ConfigInfo, error)
	GetZoneDatabaseInfoWithContext(context.Context) (*ZoneDatabaseInfo, error)
	GetZoneChecksumWithContext(context.Context) (string, error)
	GetZoneTransactionStatusWithContext(context.Context) (*ZoneTransactionStatus, error)
}

// ZoneWriter describes mutating Zone database operations.
type ZoneWriter interface {
	CreateAliasWithContext(context.Context, string, []string) error
	UpdateAliasWithContext(context.Context, string, []string) error
	RenameAliasWithContext(context.Context, string, string) error
	DeleteAliasWithContext(context.Context, string) error
	CreateZoneWithContext(context.Context, string, []string, []string) error
	UpdateZoneWithContext(context.Context, string, []string, []string) error
	RenameZoneWithContext(context.Context, string, string) error
	DeleteZoneWithContext(context.Context, string) error
	CreateZoneAndActivateWithContext(context.Context, string, string, []string, []string) error
	ReplaceZoneAndActivateWithContext(context.Context, string, string, []string, []string) error
	DeleteZoneAndActivateWithContext(context.Context, string, string) error
	UpdateDefinedConfigWithContext(context.Context, string, []string) error
	SaveZoneConfigWithContext(context.Context, string) error
	ActivateZoneConfigWithContext(context.Context, string, string) error
	AbortZoneTransactionWithContext(context.Context) error
}

// InventoryReader describes read-only FRU, media, statistics, FDMI, trunk,
// and firmware inventory operations.
type InventoryReader interface {
	GetBladesWithContext(context.Context) ([]BladeInfo, error)
	GetFansWithContext(context.Context) ([]FanInfo, error)
	GetPowerSuppliesWithContext(context.Context) ([]PowerSupplyInfo, error)
	GetHistoryLogsWithContext(context.Context) ([]HistoryLogInfo, error)
	GetSensorsWithContext(context.Context) ([]SensorInfo, error)
	GetMediaRDPsWithContext(context.Context) ([]MediaRDPInfo, error)
	GetFibreChannelStatisticsWithContext(context.Context) ([]FibreChannelStatisticsInfo, error)
	GetFibreChannelNameServersWithContext(context.Context) ([]FibreChannelNameServerInfo, error)
	GetFDMIHBAsWithContext(context.Context) ([]FDMIHBAInfo, error)
	GetFDMIPortsWithContext(context.Context) ([]FDMIPortInfo, error)
	GetTrunksWithContext(context.Context) ([]TrunkInfo, error)
	GetTrunkPerformancesWithContext(context.Context) ([]TrunkPerformanceInfo, error)
	GetTrunkAreasWithContext(context.Context) ([]TrunkAreaInfo, error)
	GetFirmwareHistoryWithContext(context.Context) ([]FirmwareHistoryInfo, error)
}

// MonitoringReader describes read-only MAPS, SNMP, time, and logging
// operations.
type MonitoringReader interface {
	GetSwitchStatusPolicyReportWithContext(context.Context) (*SwitchStatusPolicyReportInfo, error)
	GetSystemResourcesWithContext(context.Context) (*SystemResourcesInfo, error)
	GetSNMPSystemWithContext(context.Context) (*SNMPSystemInfo, error)
	GetSNMPv1AccountsWithContext(context.Context) ([]SNMPv1AccountInfo, error)
	GetSNMPv1TrapsWithContext(context.Context) ([]SNMPv1TrapInfo, error)
	GetSNMPv3AccountsWithContext(context.Context) ([]SNMPv3AccountInfo, error)
	GetSNMPv3TrapsWithContext(context.Context) ([]SNMPv3TrapInfo, error)
	GetTimeZoneWithContext(context.Context) (*TimeZoneInfo, error)
	GetClockServerWithContext(context.Context) (*ClockServerInfo, error)
	GetErrorLogsWithContext(context.Context) ([]ErrorLogInfo, error)
	GetAuditLogsWithContext(context.Context) ([]AuditLogInfo, error)
}
