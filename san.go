package san

import (
	"context"
	"io"
)

// SANSwitch 是已登录 SAN 交换机客户端，封装了底层 Client 并提供统一的
// 业务 facade 方法。通过 NewSANSwitch 创建实例后自动完成登录。
type SANSwitch struct {
	client *Client // 底层 HTTP 客户端，负责认证、请求发送和响应解析
}

// NewSANSwitch 创建并登录 SAN 交换机客户端。
// 默认使用 HTTPS 并校验证书，可通过 WithInsecureSkipVerify 或 WithHTTP 显式降低安全性。
// 登录失败时返回 nil 和 error。
func NewSANSwitch(host, username, password string, opts ...ClientOption) (*SANSwitch, error) {
	return NewSANSwitchWithContext(context.Background(), host, username, password, opts...)
}

// NewSANSwitchWithContext 创建并登录 SAN 交换机客户端，并允许取消登录请求。
func NewSANSwitchWithContext(ctx context.Context, host, username, password string, opts ...ClientOption) (*SANSwitch, error) {
	c := NewClient(host, username, password, opts...)
	sw := &SANSwitch{client: c}
	if _, err := c.LoginWithContext(ctx); err != nil {
		_ = c.Close()
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

// LoginWithContext 手动登录交换机并允许取消请求。
func (s *SANSwitch) LoginWithContext(ctx context.Context) (*LoginResponse, error) {
	return s.client.LoginWithContext(ctx)
}

// Logout 注销当前会话，清除认证 Token
func (s *SANSwitch) Logout() error {
	return s.client.Logout()
}

// LogoutWithContext 注销当前会话并允许取消请求。
func (s *SANSwitch) LogoutWithContext(ctx context.Context) error {
	return s.client.LogoutWithContext(ctx)
}

// IsLoggedIn 返回当前是否已持有有效的认证 Token
func (s *SANSwitch) IsLoggedIn() bool {
	return s.client.IsLoggedIn()
}

// Close releases idle HTTP connections held by the switch client.
func (s *SANSwitch) Close() error {
	if s == nil || s.client == nil {
		return nil
	}
	return s.client.Close()
}

// GetSwitchInfo 获取当前登录交换机的摘要信息
func (s *SANSwitch) GetSwitchInfo() (*SwitchInfo, error) {
	return s.client.GetSwitchInfo()
}

// GetSwitchInfoWithContext 获取当前登录交换机的摘要信息并允许取消请求。
func (s *SANSwitch) GetSwitchInfoWithContext(ctx context.Context) (*SwitchInfo, error) {
	return s.client.GetSwitchInfoWithContext(ctx)
}

// GetFabricSwitches 获取 Fabric 中所有交换机的详细信息
func (s *SANSwitch) GetFabricSwitches() ([]FabricSwitch, error) {
	return s.client.GetFabricSwitches()
}

// GetFabricSwitchesWithContext 获取 Fabric 中所有交换机的信息并允许取消请求。
func (s *SANSwitch) GetFabricSwitchesWithContext(ctx context.Context) ([]FabricSwitch, error) {
	return s.client.GetFabricSwitchesWithContext(ctx)
}

// GetPorts 获取交换机上所有 FC 端口的详细信息
func (s *SANSwitch) GetPorts() ([]PortInfo, error) {
	return s.client.GetPorts()
}

// GetPortsWithContext 获取交换机上所有 FC 端口并允许取消请求。
func (s *SANSwitch) GetPortsWithContext(ctx context.Context) ([]PortInfo, error) {
	return s.client.GetPortsWithContext(ctx)
}

// GetHardwareInfo 获取交换机机箱硬件信息
func (s *SANSwitch) GetHardwareInfo() (*HardwareInfo, error) {
	return s.client.GetHardwareInfo()
}

// GetHardwareInfoWithContext 获取交换机机箱硬件信息并允许取消请求。
func (s *SANSwitch) GetHardwareInfoWithContext(ctx context.Context) (*HardwareInfo, error) {
	return s.client.GetHardwareInfoWithContext(ctx)
}

// ==================== FRU 相关方法 ====================

// GetBlades 获取交换机上所有 FRU 板卡的详细信息
func (s *SANSwitch) GetBlades() ([]BladeInfo, error) {
	return s.client.GetBlades()
}

// GetBladesWithContext 获取交换机上所有 FRU 板卡并允许取消请求。
func (s *SANSwitch) GetBladesWithContext(ctx context.Context) ([]BladeInfo, error) {
	return s.client.GetBladesWithContext(ctx)
}

// GetFans 获取交换机上所有风扇单元的详细信息
func (s *SANSwitch) GetFans() ([]FanInfo, error) {
	return s.client.GetFans()
}

// GetFansWithContext 获取交换机上所有风扇并允许取消请求。
func (s *SANSwitch) GetFansWithContext(ctx context.Context) ([]FanInfo, error) {
	return s.client.GetFansWithContext(ctx)
}

// GetPowerSupplies 获取交换机上所有电源单元的详细信息
func (s *SANSwitch) GetPowerSupplies() ([]PowerSupplyInfo, error) {
	return s.client.GetPowerSupplies()
}

// GetPowerSuppliesWithContext 获取交换机上所有电源并允许取消请求。
func (s *SANSwitch) GetPowerSuppliesWithContext(ctx context.Context) ([]PowerSupplyInfo, error) {
	return s.client.GetPowerSuppliesWithContext(ctx)
}

// GetHistoryLogs 获取 FRU 组件的历史日志记录
func (s *SANSwitch) GetHistoryLogs() ([]HistoryLogInfo, error) {
	return s.client.GetHistoryLogs()
}

// GetHistoryLogsWithContext 获取 FRU 历史日志并允许取消请求。
func (s *SANSwitch) GetHistoryLogsWithContext(ctx context.Context) ([]HistoryLogInfo, error) {
	return s.client.GetHistoryLogsWithContext(ctx)
}

// GetSensors 获取交换机上所有传感器的详细信息
func (s *SANSwitch) GetSensors() ([]SensorInfo, error) {
	return s.client.GetSensors()
}

// GetSensorsWithContext 获取交换机上所有传感器并允许取消请求。
func (s *SANSwitch) GetSensorsWithContext(ctx context.Context) ([]SensorInfo, error) {
	return s.client.GetSensorsWithContext(ctx)
}

// ==================== Zone 管理相关方法 ====================

// GetDefinedZones 获取 Zone 定义配置中的所有 Zone 列表
func (s *SANSwitch) GetDefinedZones() ([]ZoneInfo, error) {
	return s.client.GetDefinedZones()
}

// GetDefinedZonesWithContext 获取 Zone 定义配置并允许取消请求。
func (s *SANSwitch) GetDefinedZonesWithContext(ctx context.Context) ([]ZoneInfo, error) {
	return s.client.GetDefinedZonesWithContext(ctx)
}

// GetDefinedZone 获取 Zone 定义配置中的单个 Zone
func (s *SANSwitch) GetDefinedZone(name string) (*ZoneInfo, error) {
	return s.client.GetDefinedZone(name)
}

// GetDefinedZoneWithContext 获取指定 Zone 定义并允许取消请求。
func (s *SANSwitch) GetDefinedZoneWithContext(ctx context.Context, name string) (*ZoneInfo, error) {
	return s.client.GetDefinedZoneWithContext(ctx, name)
}

// GetEffectiveZones 获取已生效配置中的所有 Zone 列表
func (s *SANSwitch) GetEffectiveZones() ([]ZoneInfo, error) {
	return s.client.GetEffectiveZones()
}

// GetEffectiveZonesWithContext 获取已生效 Zone 配置并允许取消请求。
func (s *SANSwitch) GetEffectiveZonesWithContext(ctx context.Context) ([]ZoneInfo, error) {
	return s.client.GetEffectiveZonesWithContext(ctx)
}

// GetDefinedAliases 获取 Zone 定义配置中的所有 Alias 列表
func (s *SANSwitch) GetDefinedAliases() ([]AliasInfo, error) {
	return s.client.GetDefinedAliases()
}

// GetDefinedAliasesWithContext 获取 Zone Alias 列表并允许取消请求。
func (s *SANSwitch) GetDefinedAliasesWithContext(ctx context.Context) ([]AliasInfo, error) {
	return s.client.GetDefinedAliasesWithContext(ctx)
}

// GetDefinedConfigs 获取 Zone 定义配置中的所有 cfg 列表
func (s *SANSwitch) GetDefinedConfigs() ([]ConfigInfo, error) {
	return s.client.GetDefinedConfigs()
}

// GetDefinedConfigsWithContext 获取 Zone cfg 列表并允许取消请求。
func (s *SANSwitch) GetDefinedConfigsWithContext(ctx context.Context) ([]ConfigInfo, error) {
	return s.client.GetDefinedConfigsWithContext(ctx)
}

// GetEffectiveConfig 获取当前已生效的 Zone 配置信息
func (s *SANSwitch) GetEffectiveConfig() (*ConfigInfo, error) {
	return s.client.GetEffectiveConfig()
}

// GetEffectiveConfigWithContext 获取当前生效 cfg 并允许取消请求。
func (s *SANSwitch) GetEffectiveConfigWithContext(ctx context.Context) (*ConfigInfo, error) {
	return s.client.GetEffectiveConfigWithContext(ctx)
}

// GetZoneDatabaseInfo 获取 Zone 数据库的容量和事务状态信息
func (s *SANSwitch) GetZoneDatabaseInfo() (*ZoneDatabaseInfo, error) {
	return s.client.GetZoneDatabaseInfo()
}

// GetZoneDatabaseInfoWithContext 获取 Zone 数据库信息并允许取消请求。
func (s *SANSwitch) GetZoneDatabaseInfoWithContext(ctx context.Context) (*ZoneDatabaseInfo, error) {
	return s.client.GetZoneDatabaseInfoWithContext(ctx)
}

// GetLogicalSwitches 获取交换机上所有逻辑交换机的配置信息
func (s *SANSwitch) GetLogicalSwitches() ([]LogicalSwitchInfo, error) {
	return s.client.GetLogicalSwitches()
}

// GetLogicalSwitchesWithContext 获取逻辑交换机并允许取消请求。
func (s *SANSwitch) GetLogicalSwitchesWithContext(ctx context.Context) ([]LogicalSwitchInfo, error) {
	return s.client.GetLogicalSwitchesWithContext(ctx)
}

// ==================== MAPS 监控相关方法 ====================

// GetSwitchStatusPolicyReport 获取交换机各组件的健康状态策略报告
func (s *SANSwitch) GetSwitchStatusPolicyReport() (*SwitchStatusPolicyReportInfo, error) {
	return s.client.GetSwitchStatusPolicyReport()
}

// GetSwitchStatusPolicyReportWithContext 获取健康状态策略报告并允许取消请求。
func (s *SANSwitch) GetSwitchStatusPolicyReportWithContext(ctx context.Context) (*SwitchStatusPolicyReportInfo, error) {
	return s.client.GetSwitchStatusPolicyReportWithContext(ctx)
}

// GetSystemResources 获取交换机系统资源使用情况
func (s *SANSwitch) GetSystemResources() (*SystemResourcesInfo, error) {
	return s.client.GetSystemResources()
}

// GetSystemResourcesWithContext 获取系统资源使用情况并允许取消请求。
func (s *SANSwitch) GetSystemResourcesWithContext(ctx context.Context) (*SystemResourcesInfo, error) {
	return s.client.GetSystemResourcesWithContext(ctx)
}

// ==================== Media (SFP) 相关方法 ====================

// GetMediaRDPs 获取所有 SFP/XFP 光模块的原始诊断参数信息
func (s *SANSwitch) GetMediaRDPs() ([]MediaRDPInfo, error) {
	return s.client.GetMediaRDPs()
}

// GetMediaRDPsWithContext 获取 SFP/XFP 诊断信息并允许取消请求。
func (s *SANSwitch) GetMediaRDPsWithContext(ctx context.Context) ([]MediaRDPInfo, error) {
	return s.client.GetMediaRDPsWithContext(ctx)
}

// ==================== Statistics 相关方法 ====================

// GetFibreChannelStatistics 获取所有 FC 端口的性能统计计数器
func (s *SANSwitch) GetFibreChannelStatistics() ([]FibreChannelStatisticsInfo, error) {
	return s.client.GetFibreChannelStatistics()
}

// GetFibreChannelStatisticsWithContext 获取 FC 性能统计并允许取消请求。
func (s *SANSwitch) GetFibreChannelStatisticsWithContext(ctx context.Context) ([]FibreChannelStatisticsInfo, error) {
	return s.client.GetFibreChannelStatisticsWithContext(ctx)
}

// ==================== Name Server 相关方法 ====================

// GetFibreChannelNameServers 获取 Fabric Name Server 中注册的所有设备条目
func (s *SANSwitch) GetFibreChannelNameServers() ([]FibreChannelNameServerInfo, error) {
	return s.client.GetFibreChannelNameServers()
}

// GetFibreChannelNameServersWithContext 获取 Name Server 条目并允许取消请求。
func (s *SANSwitch) GetFibreChannelNameServersWithContext(ctx context.Context) ([]FibreChannelNameServerInfo, error) {
	return s.client.GetFibreChannelNameServersWithContext(ctx)
}

// ==================== FDMI 相关方法 ====================

// GetFDMIHBAs 获取 FDMI 中注册的所有 HBA 适配器信息
func (s *SANSwitch) GetFDMIHBAs() ([]FDMIHBAInfo, error) {
	return s.client.GetFDMIHBAs()
}

// GetFDMIHBAsWithContext 获取 FDMI HBA 信息并允许取消请求。
func (s *SANSwitch) GetFDMIHBAsWithContext(ctx context.Context) ([]FDMIHBAInfo, error) {
	return s.client.GetFDMIHBAsWithContext(ctx)
}

// GetFDMIPorts 获取 FDMI 中注册的所有 FC 端口信息。
func (s *SANSwitch) GetFDMIPorts() ([]FDMIPortInfo, error) {
	return s.client.GetFDMIPorts()
}

// GetFDMIPortsWithContext 获取 FDMI 端口信息并允许取消请求。
func (s *SANSwitch) GetFDMIPortsWithContext(ctx context.Context) ([]FDMIPortInfo, error) {
	return s.client.GetFDMIPortsWithContext(ctx)
}

// ==================== Logging 相关方法 ====================

// GetErrorLogs 获取交换机上的所有 RAS 错误日志记录
func (s *SANSwitch) GetErrorLogs() ([]ErrorLogInfo, error) {
	return s.client.GetErrorLogs()
}

// GetErrorLogsWithContext 获取 RAS 错误日志并允许取消请求。
func (s *SANSwitch) GetErrorLogsWithContext(ctx context.Context) ([]ErrorLogInfo, error) {
	return s.client.GetErrorLogsWithContext(ctx)
}

// GetAuditLogs 获取交换机上的所有审计日志记录
func (s *SANSwitch) GetAuditLogs() ([]AuditLogInfo, error) {
	return s.client.GetAuditLogs()
}

// GetAuditLogsWithContext 获取审计日志并允许取消请求。
func (s *SANSwitch) GetAuditLogsWithContext(ctx context.Context) ([]AuditLogInfo, error) {
	return s.client.GetAuditLogsWithContext(ctx)
}

// ==================== Alias 管理方法 ====================

// CreateAlias 在 Zone 定义配置中创建一个新的 Alias
func (s *SANSwitch) CreateAlias(name string, members []string) error {
	return s.client.CreateAlias(name, members)
}

// CreateAliasWithContext 创建 Alias 并允许取消请求。
func (s *SANSwitch) CreateAliasWithContext(ctx context.Context, name string, members []string) error {
	return s.client.CreateAliasWithContext(ctx, name, members)
}

// UpdateAlias 更新 Zone 定义配置中已有 Alias 的成员列表
func (s *SANSwitch) UpdateAlias(name string, members []string) error {
	return s.client.UpdateAlias(name, members)
}

// UpdateAliasWithContext 更新 Alias 并允许取消请求。
func (s *SANSwitch) UpdateAliasWithContext(ctx context.Context, name string, members []string) error {
	return s.client.UpdateAliasWithContext(ctx, name, members)
}

// RenameAlias 重命名 Zone 定义配置中的一个 Alias
func (s *SANSwitch) RenameAlias(oldName, newName string) error {
	return s.client.RenameAlias(oldName, newName)
}

// RenameAliasWithContext 重命名 Alias 并允许取消请求。
func (s *SANSwitch) RenameAliasWithContext(ctx context.Context, oldName, newName string) error {
	return s.client.RenameAliasWithContext(ctx, oldName, newName)
}

// DeleteAlias 从 Zone 定义配置中删除一个 Alias
func (s *SANSwitch) DeleteAlias(name string) error {
	return s.client.DeleteAlias(name)
}

// DeleteAliasWithContext 删除 Alias 并允许取消请求。
func (s *SANSwitch) DeleteAliasWithContext(ctx context.Context, name string) error {
	return s.client.DeleteAliasWithContext(ctx, name)
}

// ==================== Zone 操作方法 ====================

// CreateZone 在 Zone 定义配置中创建一个新的 Zone
func (s *SANSwitch) CreateZone(name string, members []string, principalMembers []string) error {
	return s.client.CreateZone(name, members, principalMembers)
}

// CreateZoneWithContext 创建 Zone 并允许取消请求。
func (s *SANSwitch) CreateZoneWithContext(ctx context.Context, name string, members []string, principalMembers []string) error {
	return s.client.CreateZoneWithContext(ctx, name, members, principalMembers)
}

// UpdateZone 更新 Zone 定义配置中已有 Zone 的成员列表
func (s *SANSwitch) UpdateZone(name string, members []string, principalMembers []string) error {
	return s.client.UpdateZone(name, members, principalMembers)
}

// UpdateZoneWithContext 更新 Zone 并允许取消请求。
func (s *SANSwitch) UpdateZoneWithContext(ctx context.Context, name string, members []string, principalMembers []string) error {
	return s.client.UpdateZoneWithContext(ctx, name, members, principalMembers)
}

// RenameZone 重命名 Zone 定义配置中的一个 Zone
func (s *SANSwitch) RenameZone(oldName, newName string) error {
	return s.client.RenameZone(oldName, newName)
}

// RenameZoneWithContext 重命名 Zone 并允许取消请求。
func (s *SANSwitch) RenameZoneWithContext(ctx context.Context, oldName, newName string) error {
	return s.client.RenameZoneWithContext(ctx, oldName, newName)
}

// DeleteZone 从 Zone 定义配置中删除一个 Zone
func (s *SANSwitch) DeleteZone(name string) error {
	return s.client.DeleteZone(name)
}

// DeleteZoneWithContext 删除 Zone 并允许取消请求。
func (s *SANSwitch) DeleteZoneWithContext(ctx context.Context, name string) error {
	return s.client.DeleteZoneWithContext(ctx, name)
}

// CreateZoneAndActivate 执行完整的 Zone 创建并激活流程
func (s *SANSwitch) CreateZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error {
	return s.client.CreateZoneAndActivate(cfgName, zoneName, members, principalMembers)
}

// CreateZoneAndActivateWithContext 创建并激活 Zone，并允许取消请求。
func (s *SANSwitch) CreateZoneAndActivateWithContext(ctx context.Context, cfgName, zoneName string, members []string, principalMembers []string) error {
	return s.client.CreateZoneAndActivateWithContext(ctx, cfgName, zoneName, members, principalMembers)
}

// ReplaceZoneAndActivate 执行完整的 Zone 替换并激活流程
func (s *SANSwitch) ReplaceZoneAndActivate(cfgName, zoneName string, members []string, principalMembers []string) error {
	return s.client.ReplaceZoneAndActivate(cfgName, zoneName, members, principalMembers)
}

// ReplaceZoneAndActivateWithContext 替换并激活 Zone，并允许取消请求。
func (s *SANSwitch) ReplaceZoneAndActivateWithContext(ctx context.Context, cfgName, zoneName string, members []string, principalMembers []string) error {
	return s.client.ReplaceZoneAndActivateWithContext(ctx, cfgName, zoneName, members, principalMembers)
}

// DeleteZoneAndActivate 执行完整的 Zone 删除并激活流程
func (s *SANSwitch) DeleteZoneAndActivate(cfgName, zoneName string) error {
	return s.client.DeleteZoneAndActivate(cfgName, zoneName)
}

// DeleteZoneAndActivateWithContext 删除并激活 Zone，并允许取消请求。
func (s *SANSwitch) DeleteZoneAndActivateWithContext(ctx context.Context, cfgName, zoneName string) error {
	return s.client.DeleteZoneAndActivateWithContext(ctx, cfgName, zoneName)
}

// GetZoneChecksum 获取当前 Zone 配置的校验和
func (s *SANSwitch) GetZoneChecksum() (string, error) {
	return s.client.GetZoneChecksum()
}

// GetZoneChecksumWithContext 获取 Zone checksum 并允许取消请求。
func (s *SANSwitch) GetZoneChecksumWithContext(ctx context.Context) (string, error) {
	return s.client.GetZoneChecksumWithContext(ctx)
}

// UpdateDefinedConfig 更新 Zone 定义配置中指定 cfg 的成员 Zone 列表
func (s *SANSwitch) UpdateDefinedConfig(name string, memberZones []string) error {
	return s.client.UpdateDefinedConfig(name, memberZones)
}

// UpdateDefinedConfigWithContext 更新 Zone cfg 并允许取消请求。
func (s *SANSwitch) UpdateDefinedConfigWithContext(ctx context.Context, name string, memberZones []string) error {
	return s.client.UpdateDefinedConfigWithContext(ctx, name, memberZones)
}

// SaveZoneConfig 保存当前 Zone 定义配置到持久存储
func (s *SANSwitch) SaveZoneConfig(checksum string) error {
	return s.client.SaveZoneConfig(checksum)
}

// SaveZoneConfigWithContext 保存 Zone cfg 并允许取消请求。
func (s *SANSwitch) SaveZoneConfigWithContext(ctx context.Context, checksum string) error {
	return s.client.SaveZoneConfigWithContext(ctx, checksum)
}

// ActivateZoneConfig 激活指定的 Zone 配置（cfg）
func (s *SANSwitch) ActivateZoneConfig(name string, checksum string) error {
	return s.client.ActivateZoneConfig(name, checksum)
}

// ActivateZoneConfigWithContext 激活 Zone cfg 并允许取消请求。
func (s *SANSwitch) ActivateZoneConfigWithContext(ctx context.Context, name string, checksum string) error {
	return s.client.ActivateZoneConfigWithContext(ctx, name, checksum)
}

// AbortZoneTransaction 中止当前未完成的 Zone 事务
func (s *SANSwitch) AbortZoneTransaction() error {
	return s.client.AbortZoneTransaction()
}

// AbortZoneTransactionWithContext 中止 Zone 事务并允许取消请求。
func (s *SANSwitch) AbortZoneTransactionWithContext(ctx context.Context) error {
	return s.client.AbortZoneTransactionWithContext(ctx)
}

// GetZoneTransactionStatus 获取当前 Zone 事务的状态信息
func (s *SANSwitch) GetZoneTransactionStatus() (*ZoneTransactionStatus, error) {
	return s.client.GetZoneTransactionStatus()
}

// GetZoneTransactionStatusWithContext 查询 Zone 事务状态并允许取消请求。
func (s *SANSwitch) GetZoneTransactionStatusWithContext(ctx context.Context) (*ZoneTransactionStatus, error) {
	return s.client.GetZoneTransactionStatusWithContext(ctx)
}

// ==================== Trunk (ISL 链路聚合) 相关方法 ====================

// GetTrunks 获取所有 ISL Trunk 链路的详细信息
func (s *SANSwitch) GetTrunks() ([]TrunkInfo, error) {
	return s.client.GetTrunks()
}

// GetTrunksWithContext 获取 Trunk 成员并允许取消请求。
func (s *SANSwitch) GetTrunksWithContext(ctx context.Context) ([]TrunkInfo, error) {
	return s.client.GetTrunksWithContext(ctx)
}

// GetTrunkPerformances 获取所有 Trunk 链路的性能统计信息
func (s *SANSwitch) GetTrunkPerformances() ([]TrunkPerformanceInfo, error) {
	return s.client.GetTrunkPerformances()
}

// GetTrunkPerformancesWithContext 获取 Trunk 性能统计并允许取消请求。
func (s *SANSwitch) GetTrunkPerformancesWithContext(ctx context.Context) ([]TrunkPerformanceInfo, error) {
	return s.client.GetTrunkPerformancesWithContext(ctx)
}

// GetTrunkAreas 获取所有 Trunk 的 Area 信息
func (s *SANSwitch) GetTrunkAreas() ([]TrunkAreaInfo, error) {
	return s.client.GetTrunkAreas()
}

// GetTrunkAreasWithContext 获取 Trunk Area 信息并允许取消请求。
func (s *SANSwitch) GetTrunkAreasWithContext(ctx context.Context) ([]TrunkAreaInfo, error) {
	return s.client.GetTrunkAreasWithContext(ctx)
}

// ==================== Firmware 相关方法 ====================

// GetFirmwareHistory 获取固件升级历史记录
func (s *SANSwitch) GetFirmwareHistory() ([]FirmwareHistoryInfo, error) {
	return s.client.GetFirmwareHistory()
}

// GetFirmwareHistoryWithContext 获取固件升级历史并允许取消请求。
func (s *SANSwitch) GetFirmwareHistoryWithContext(ctx context.Context) ([]FirmwareHistoryInfo, error) {
	return s.client.GetFirmwareHistoryWithContext(ctx)
}

// ==================== SNMP 相关方法 ====================

// GetSNMPSystem 获取 SNMP 系统全局配置信息
func (s *SANSwitch) GetSNMPSystem() (*SNMPSystemInfo, error) {
	return s.client.GetSNMPSystem()
}

// GetSNMPSystemWithContext 获取 SNMP 系统配置并允许取消请求。
func (s *SANSwitch) GetSNMPSystemWithContext(ctx context.Context) (*SNMPSystemInfo, error) {
	return s.client.GetSNMPSystemWithContext(ctx)
}

// GetSNMPv1Accounts 获取所有 SNMPv1 社区账户信息
func (s *SANSwitch) GetSNMPv1Accounts() ([]SNMPv1AccountInfo, error) {
	return s.client.GetSNMPv1Accounts()
}

// GetSNMPv1AccountsWithContext 获取 SNMPv1 账户并允许取消请求。
func (s *SANSwitch) GetSNMPv1AccountsWithContext(ctx context.Context) ([]SNMPv1AccountInfo, error) {
	return s.client.GetSNMPv1AccountsWithContext(ctx)
}

// GetSNMPv1Traps 获取所有 SNMPv1 Trap 接收器配置
func (s *SANSwitch) GetSNMPv1Traps() ([]SNMPv1TrapInfo, error) {
	return s.client.GetSNMPv1Traps()
}

// GetSNMPv1TrapsWithContext 获取 SNMPv1 Trap 配置并允许取消请求。
func (s *SANSwitch) GetSNMPv1TrapsWithContext(ctx context.Context) ([]SNMPv1TrapInfo, error) {
	return s.client.GetSNMPv1TrapsWithContext(ctx)
}

// GetSNMPv3Accounts 获取所有 SNMPv3 用户账户信息
func (s *SANSwitch) GetSNMPv3Accounts() ([]SNMPv3AccountInfo, error) {
	return s.client.GetSNMPv3Accounts()
}

// GetSNMPv3AccountsWithContext 获取 SNMPv3 账户并允许取消请求。
func (s *SANSwitch) GetSNMPv3AccountsWithContext(ctx context.Context) ([]SNMPv3AccountInfo, error) {
	return s.client.GetSNMPv3AccountsWithContext(ctx)
}

// GetSNMPv3Traps 获取所有 SNMPv3 Trap 接收器配置
func (s *SANSwitch) GetSNMPv3Traps() ([]SNMPv3TrapInfo, error) {
	return s.client.GetSNMPv3Traps()
}

// GetSNMPv3TrapsWithContext 获取 SNMPv3 Trap 配置并允许取消请求。
func (s *SANSwitch) GetSNMPv3TrapsWithContext(ctx context.Context) ([]SNMPv3TrapInfo, error) {
	return s.client.GetSNMPv3TrapsWithContext(ctx)
}

// ==================== Time / NTP 相关方法 ====================

// GetTimeZone 获取交换机当前时区配置
func (s *SANSwitch) GetTimeZone() (*TimeZoneInfo, error) {
	return s.client.GetTimeZone()
}

// GetTimeZoneWithContext 获取时区配置并允许取消请求。
func (s *SANSwitch) GetTimeZoneWithContext(ctx context.Context) (*TimeZoneInfo, error) {
	return s.client.GetTimeZoneWithContext(ctx)
}

// GetClockServer 获取交换机的 NTP 时钟服务器配置
func (s *SANSwitch) GetClockServer() (*ClockServerInfo, error) {
	return s.client.GetClockServer()
}

// GetClockServerWithContext 获取 NTP 时钟服务器并允许取消请求。
func (s *SANSwitch) GetClockServerWithContext(ctx context.Context) (*ClockServerInfo, error) {
	return s.client.GetClockServerWithContext(ctx)
}

var (
	_ Session          = (*SANSwitch)(nil)
	_ SwitchReader     = (*SANSwitch)(nil)
	_ ZoneReader       = (*SANSwitch)(nil)
	_ ZoneWriter       = (*SANSwitch)(nil)
	_ InventoryReader  = (*SANSwitch)(nil)
	_ MonitoringReader = (*SANSwitch)(nil)
	_ Session          = (*Client)(nil)
	_ SwitchReader     = (*Client)(nil)
	_ ZoneReader       = (*Client)(nil)
	_ ZoneWriter       = (*Client)(nil)
	_ InventoryReader  = (*Client)(nil)
	_ MonitoringReader = (*Client)(nil)
)
