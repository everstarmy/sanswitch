package san

import "io"

// SwitchAPI 定义 SAN 交换机的核心操作接口，支持 Mock 测试和替换底层实现
type SwitchAPI interface {
	// 连接管理
	Login() (*LoginResponse, error)
	Logout() error
	IsLoggedIn() bool

	// 交换机与端口
	GetSwitchInfo() (*SwitchInfo, error)
	GetFabricSwitches() ([]FabricSwitch, error)
	GetPorts() ([]PortInfo, error)
	GetHardwareInfo() (*HardwareInfo, error)

	// FRU 组件
	GetBlades() ([]BladeInfo, error)
	GetFans() ([]FanInfo, error)
	GetPowerSupplies() ([]PowerSupplyInfo, error)
	GetHistoryLogs() ([]HistoryLogInfo, error)
	GetSensors() ([]SensorInfo, error)

	// Zone 管理
	GetDefinedZones() ([]ZoneInfo, error)
	GetDefinedZone(name string) (*ZoneInfo, error)
	GetEffectiveZones() ([]ZoneInfo, error)
	GetDefinedAliases() ([]AliasInfo, error)
	GetDefinedConfigs() ([]ConfigInfo, error)
	GetEffectiveConfig() (*ConfigInfo, error)
	GetZoneDatabaseInfo() (*ZoneDatabaseInfo, error)
	CreateAlias(name string, members []string) error
	UpdateAlias(name string, members []string) error
	RenameAlias(oldName, newName string) error
	DeleteAlias(name string) error
	CreateZone(name string, members []string, principalMembers []string) error
	UpdateZone(name string, members []string, principalMembers []string) error
	RenameZone(oldName, newName string) error
	DeleteZone(name string) error
	CreateZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error
	ReplaceZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error
	DeleteZoneAndActivate(cfgName, zoneName string) error
	GetZoneChecksum() (string, error)
	UpdateDefinedConfig(name string, memberZones []string) error
	SaveZoneConfig(checksum string) error
	ActivateZoneConfig(name string, checksum string) error
	AbortZoneTransaction() error
	GetZoneTransactionStatus() (*ZoneTransactionStatus, error)

	// 逻辑交换机
	GetLogicalSwitches() ([]LogicalSwitchInfo, error)

	// MAPS 监控
	GetSwitchStatusPolicyReport() (*SwitchStatusPolicyReportInfo, error)
	GetSystemResources() (*SystemResourcesInfo, error)

	// Media (SFP)
	GetMediaRDPs() ([]MediaRDPInfo, error)

	// Statistics
	GetFibreChannelStatistics() ([]FibreChannelStatisticsInfo, error)

	// Name Server
	GetFibreChannelNameServers() ([]FibreChannelNameServerInfo, error)

	// FDMI
	GetFDMIHBAs() ([]FDMIHBAInfo, error)
	GetFDMIports() ([]FDMIportInfo, error)

	// Trunk (ISL 链路聚合)
	GetTrunks() ([]TrunkInfo, error)
	GetTrunkPerformances() ([]TrunkPerformanceInfo, error)
	GetTrunkAreas() ([]TrunkAreaInfo, error)

	// Firmware
	GetFirmwareHistory() ([]FirmwareHistoryInfo, error)

	// SNMP
	GetSNMPSystem() (*SNMPSystemInfo, error)
	GetSNMPv1Accounts() ([]SNMPv1AccountInfo, error)
	GetSNMPv1Traps() ([]SNMPv1TrapInfo, error)
	GetSNMPv3Accounts() ([]SNMPv3AccountInfo, error)
	GetSNMPv3Traps() ([]SNMPv3TrapInfo, error)

	// Time / NTP
	GetTimeZone() (*TimeZoneInfo, error)
	GetClockServer() (*ClockServerInfo, error)

	// Logging
	GetErrorLogs() ([]ErrorLogInfo, error)
	GetAuditLogs() ([]AuditLogInfo, error)
}

// SANSwitch 是 SwitchAPI 接口的默认实现，封装了底层 Client 并提供统一的 facade 方法。
// 通过 NewSANSwitch 创建实例后自动完成登录，可直接调用各业务方法。
type SANSwitch struct {
	client *Client // 底层 HTTP 客户端，负责认证、请求发送和响应解析
}

// NewSANSwitch 创建并登录 SAN 交换机客户端。
// 默认使用 HTTPS（跳过证书验证），可通过 WithHTTP() 切换到 HTTP。
// 登录失败时返回 nil 和 error。
func NewSANSwitch(host, username, password string, opts ...ClientOption) (*SANSwitch, error) {
	c := NewClient(host, username, password, opts...)
	sw := &SANSwitch{client: c}
	if _, err := c.Login(); err != nil {
		return nil, err
	}
	return sw, nil
}

// SetVerbose 开启或关闭调试级别日志输出
func (s *SANSwitch) SetVerbose(verbose bool) {
	s.client.SetVerbose(verbose)
}

// SetLogOutput 设置日志输出目标（nil 则恢复为 os.Stderr）
func (s *SANSwitch) SetLogOutput(w io.Writer) {
	s.client.SetLogOutput(w)
}

// SetVFID 设置虚拟 Fabric ID，用于 Virtual Fabric 场景下的请求路由
func (s *SANSwitch) SetVFID(vfID int) {
	s.client.SetVFID(vfID)
}

// Login 手动登录交换机，获取认证 Token。
// 通常由 NewSANSwitch 自动调用，仅在 Logout 后需要重新登录时使用。
func (s *SANSwitch) Login() (*LoginResponse, error) {
	return s.client.Login()
}

// Logout 注销当前会话，清除认证 Token
func (s *SANSwitch) Logout() error {
	return s.client.Logout()
}

// IsLoggedIn 返回当前是否已持有有效的认证 Token
func (s *SANSwitch) IsLoggedIn() bool {
	return s.client.IsLoggedIn()
}

// GetSwitchInfo 获取当前登录交换机的摘要信息
func (s *SANSwitch) GetSwitchInfo() (*SwitchInfo, error) {
	return s.client.GetSwitchInfo()
}

// GetFabricSwitches 获取 Fabric 中所有交换机的详细信息
func (s *SANSwitch) GetFabricSwitches() ([]FabricSwitch, error) {
	return s.client.GetFabricSwitches()
}

// GetPorts 获取交换机上所有 FC 端口的详细信息
func (s *SANSwitch) GetPorts() ([]PortInfo, error) {
	return s.client.GetPorts()
}

// GetHardwareInfo 获取交换机机箱硬件信息
func (s *SANSwitch) GetHardwareInfo() (*HardwareInfo, error) {
	return s.client.GetHardwareInfo()
}

// ==================== FRU 相关方法 ====================

// GetBlades 获取交换机上所有 FRU 板卡的详细信息
func (s *SANSwitch) GetBlades() ([]BladeInfo, error) {
	return s.client.GetBlades()
}

// GetFans 获取交换机上所有风扇单元的详细信息
func (s *SANSwitch) GetFans() ([]FanInfo, error) {
	return s.client.GetFans()
}

// GetPowerSupplies 获取交换机上所有电源单元的详细信息
func (s *SANSwitch) GetPowerSupplies() ([]PowerSupplyInfo, error) {
	return s.client.GetPowerSupplies()
}

// GetHistoryLogs 获取 FRU 组件的历史日志记录
func (s *SANSwitch) GetHistoryLogs() ([]HistoryLogInfo, error) {
	return s.client.GetHistoryLogs()
}

// GetSensors 获取交换机上所有传感器的详细信息
func (s *SANSwitch) GetSensors() ([]SensorInfo, error) {
	return s.client.GetSensors()
}

// ==================== Zone 管理相关方法 ====================

// GetDefinedZones 获取 Zone 定义配置中的所有 Zone 列表
func (s *SANSwitch) GetDefinedZones() ([]ZoneInfo, error) {
	return s.client.GetDefinedZones()
}

// GetDefinedZone 获取 Zone 定义配置中的单个 Zone
func (s *SANSwitch) GetDefinedZone(name string) (*ZoneInfo, error) {
	return s.client.GetDefinedZone(name)
}

// GetEffectiveZones 获取已生效配置中的所有 Zone 列表
func (s *SANSwitch) GetEffectiveZones() ([]ZoneInfo, error) {
	return s.client.GetEffectiveZones()
}

// GetDefinedAliases 获取 Zone 定义配置中的所有 Alias 列表
func (s *SANSwitch) GetDefinedAliases() ([]AliasInfo, error) {
	return s.client.GetDefinedAliases()
}

// GetDefinedConfigs 获取 Zone 定义配置中的所有 cfg 列表
func (s *SANSwitch) GetDefinedConfigs() ([]ConfigInfo, error) {
	return s.client.GetDefinedConfigs()
}

// GetEffectiveConfig 获取当前已生效的 Zone 配置信息
func (s *SANSwitch) GetEffectiveConfig() (*ConfigInfo, error) {
	return s.client.GetEffectiveConfig()
}

// GetZoneDatabaseInfo 获取 Zone 数据库的容量和事务状态信息
func (s *SANSwitch) GetZoneDatabaseInfo() (*ZoneDatabaseInfo, error) {
	return s.client.GetZoneDatabaseInfo()
}

// GetLogicalSwitches 获取交换机上所有逻辑交换机的配置信息
func (s *SANSwitch) GetLogicalSwitches() ([]LogicalSwitchInfo, error) {
	return s.client.GetLogicalSwitches()
}

// ==================== MAPS 监控相关方法 ====================

// GetSwitchStatusPolicyReport 获取交换机各组件的健康状态策略报告
func (s *SANSwitch) GetSwitchStatusPolicyReport() (*SwitchStatusPolicyReportInfo, error) {
	return s.client.GetSwitchStatusPolicyReport()
}

// GetSystemResources 获取交换机系统资源使用情况
func (s *SANSwitch) GetSystemResources() (*SystemResourcesInfo, error) {
	return s.client.GetSystemResources()
}

// ==================== Media (SFP) 相关方法 ====================

// GetMediaRDPs 获取所有 SFP/XFP 光模块的原始诊断参数信息
func (s *SANSwitch) GetMediaRDPs() ([]MediaRDPInfo, error) {
	return s.client.GetMediaRDPs()
}

// ==================== Statistics 相关方法 ====================

// GetFibreChannelStatistics 获取所有 FC 端口的性能统计计数器
func (s *SANSwitch) GetFibreChannelStatistics() ([]FibreChannelStatisticsInfo, error) {
	return s.client.GetFibreChannelStatistics()
}

// ==================== Name Server 相关方法 ====================

// GetFibreChannelNameServers 获取 Fabric Name Server 中注册的所有设备条目
func (s *SANSwitch) GetFibreChannelNameServers() ([]FibreChannelNameServerInfo, error) {
	return s.client.GetFibreChannelNameServers()
}

// ==================== FDMI 相关方法 ====================

// GetFDMIHBAs 获取 FDMI 中注册的所有 HBA 适配器信息
func (s *SANSwitch) GetFDMIHBAs() ([]FDMIHBAInfo, error) {
	return s.client.GetFDMIHBAs()
}

// GetFDMIports 获取 FDMI 中注册的所有 FC 端口信息
func (s *SANSwitch) GetFDMIports() ([]FDMIportInfo, error) {
	return s.client.GetFDMIports()
}

// ==================== Logging 相关方法 ====================

// GetErrorLogs 获取交换机上的所有 RAS 错误日志记录
func (s *SANSwitch) GetErrorLogs() ([]ErrorLogInfo, error) {
	return s.client.GetErrorLogs()
}

// GetAuditLogs 获取交换机上的所有审计日志记录
func (s *SANSwitch) GetAuditLogs() ([]AuditLogInfo, error) {
	return s.client.GetAuditLogs()
}

// ==================== Alias 管理方法 ====================

// CreateAlias 在 Zone 定义配置中创建一个新的 Alias
func (s *SANSwitch) CreateAlias(name string, members []string) error {
	return s.client.CreateAlias(name, members)
}

// UpdateAlias 更新 Zone 定义配置中已有 Alias 的成员列表
func (s *SANSwitch) UpdateAlias(name string, members []string) error {
	return s.client.UpdateAlias(name, members)
}

// RenameAlias 重命名 Zone 定义配置中的一个 Alias
func (s *SANSwitch) RenameAlias(oldName, newName string) error {
	return s.client.RenameAlias(oldName, newName)
}

// DeleteAlias 从 Zone 定义配置中删除一个 Alias
func (s *SANSwitch) DeleteAlias(name string) error {
	return s.client.DeleteAlias(name)
}

// ==================== Zone 操作方法 ====================

// CreateZone 在 Zone 定义配置中创建一个新的 Zone
func (s *SANSwitch) CreateZone(name string, members []string, principalMembers []string) error {
	return s.client.CreateZone(name, members, principalMembers)
}

// UpdateZone 更新 Zone 定义配置中已有 Zone 的成员列表
func (s *SANSwitch) UpdateZone(name string, members []string, principalMembers []string) error {
	return s.client.UpdateZone(name, members, principalMembers)
}

// RenameZone 重命名 Zone 定义配置中的一个 Zone
func (s *SANSwitch) RenameZone(oldName, newName string) error {
	return s.client.RenameZone(oldName, newName)
}

// DeleteZone 从 Zone 定义配置中删除一个 Zone
func (s *SANSwitch) DeleteZone(name string) error {
	return s.client.DeleteZone(name)
}

// CreateZoneAndActivate 执行完整的 Zone 创建并激活流程
func (s *SANSwitch) CreateZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error {
	return s.client.CreateZoneAndActivate(cfgName, zoneName, members, principalMembers)
}

// ReplaceZoneAndActivate 执行完整的 Zone 替换并激活流程
func (s *SANSwitch) ReplaceZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error {
	return s.client.ReplaceZoneAndActivate(cfgName, zoneName, members, principalMembers)
}

// DeleteZoneAndActivate 执行完整的 Zone 删除并激活流程
func (s *SANSwitch) DeleteZoneAndActivate(cfgName, zoneName string) error {
	return s.client.DeleteZoneAndActivate(cfgName, zoneName)
}

// GetZoneChecksum 获取当前 Zone 配置的校验和
func (s *SANSwitch) GetZoneChecksum() (string, error) {
	return s.client.GetZoneChecksum()
}

// UpdateDefinedConfig 更新 Zone 定义配置中指定 cfg 的成员 Zone 列表
func (s *SANSwitch) UpdateDefinedConfig(name string, memberZones []string) error {
	return s.client.UpdateDefinedConfig(name, memberZones)
}

// SaveZoneConfig 保存当前 Zone 定义配置到持久存储
func (s *SANSwitch) SaveZoneConfig(checksum string) error {
	return s.client.SaveZoneConfig(checksum)
}

// ActivateZoneConfig 激活指定的 Zone 配置（cfg）
func (s *SANSwitch) ActivateZoneConfig(name string, checksum string) error {
	return s.client.ActivateZoneConfig(name, checksum)
}

// AbortZoneTransaction 中止当前未完成的 Zone 事务
func (s *SANSwitch) AbortZoneTransaction() error {
	return s.client.AbortZoneTransaction()
}

// GetZoneTransactionStatus 获取当前 Zone 事务的状态信息
func (s *SANSwitch) GetZoneTransactionStatus() (*ZoneTransactionStatus, error) {
	return s.client.GetZoneTransactionStatus()
}

// ==================== Trunk (ISL 链路聚合) 相关方法 ====================

// GetTrunks 获取所有 ISL Trunk 链路的详细信息
func (s *SANSwitch) GetTrunks() ([]TrunkInfo, error) {
	return s.client.GetTrunks()
}

// GetTrunkPerformances 获取所有 Trunk 链路的性能统计信息
func (s *SANSwitch) GetTrunkPerformances() ([]TrunkPerformanceInfo, error) {
	return s.client.GetTrunkPerformances()
}

// GetTrunkAreas 获取所有 Trunk 的 Area 信息
func (s *SANSwitch) GetTrunkAreas() ([]TrunkAreaInfo, error) {
	return s.client.GetTrunkAreas()
}

// ==================== Firmware 相关方法 ====================

// GetFirmwareHistory 获取固件升级历史记录
func (s *SANSwitch) GetFirmwareHistory() ([]FirmwareHistoryInfo, error) {
	return s.client.GetFirmwareHistory()
}

// ==================== SNMP 相关方法 ====================

// GetSNMPSystem 获取 SNMP 系统全局配置信息
func (s *SANSwitch) GetSNMPSystem() (*SNMPSystemInfo, error) {
	return s.client.GetSNMPSystem()
}

// GetSNMPv1Accounts 获取所有 SNMPv1 社区账户信息
func (s *SANSwitch) GetSNMPv1Accounts() ([]SNMPv1AccountInfo, error) {
	return s.client.GetSNMPv1Accounts()
}

// GetSNMPv1Traps 获取所有 SNMPv1 Trap 接收器配置
func (s *SANSwitch) GetSNMPv1Traps() ([]SNMPv1TrapInfo, error) {
	return s.client.GetSNMPv1Traps()
}

// GetSNMPv3Accounts 获取所有 SNMPv3 用户账户信息
func (s *SANSwitch) GetSNMPv3Accounts() ([]SNMPv3AccountInfo, error) {
	return s.client.GetSNMPv3Accounts()
}

// GetSNMPv3Traps 获取所有 SNMPv3 Trap 接收器配置
func (s *SANSwitch) GetSNMPv3Traps() ([]SNMPv3TrapInfo, error) {
	return s.client.GetSNMPv3Traps()
}

// ==================== Time / NTP 相关方法 ====================

// GetTimeZone 获取交换机当前时区配置
func (s *SANSwitch) GetTimeZone() (*TimeZoneInfo, error) {
	return s.client.GetTimeZone()
}

// GetClockServer 获取交换机的 NTP 时钟服务器配置
func (s *SANSwitch) GetClockServer() (*ClockServerInfo, error) {
	return s.client.GetClockServer()
}

// 编译期断言：确保 *SANSwitch 实现 SwitchAPI 接口
var _ SwitchAPI = (*SANSwitch)(nil)
