package san

import "encoding/xml"

// SSPStateType 表示 MAPS（Monitoring and Alerting Policy Suite）中各组件的健康状态等级。
// 可能的值: ok, warning, critical, unknown
type SSPStateType string

const (
	SSPStateOK       SSPStateType = "ok"       // 正常
	SSPStateWarning  SSPStateType = "warning"  // 警告
	SSPStateCritical SSPStateType = "critical" // 严重
	SSPStateUnknown  SSPStateType = "unknown"  // 未知
)

// ==================== Switch Status Policy Report ====================

// SwitchStatusPolicyReportResponse 是 GET /brocade-maps/switch-status-policy-report 的 XML 响应包装
type SwitchStatusPolicyReportResponse struct {
	XMLName                  xml.Name                 `xml:"Response"`
	SwitchStatusPolicyReport SwitchStatusPolicyReport `xml:"switch-status-policy-report"`
}

// SwitchStatusPolicyReport 描述交换机各组件的健康状态策略报告（XML 原始结构），
// 包含交换机整体健康度、电源、风扇、温度传感器、端口、SFP 等组件状态。
// 对应 YANG 模型: brocade-maps/switch-status-policy-report
type SwitchStatusPolicyReport struct {
	XMLName                     xml.Name     `xml:"switch-status-policy-report"`
	SwitchHealth                SSPStateType `xml:"switch-health"`
	PowerSupplyHealth           SSPStateType `xml:"power-supply-health"`
	FanHealth                   SSPStateType `xml:"fan-health"`
	WWNHealth                   SSPStateType `xml:"wwn-health"`
	TemperatureSensorHealth     SSPStateType `xml:"temperature-sensor-health"`
	HAHealth                    SSPStateType `xml:"ha-health"`
	ControlProcessorHealth      SSPStateType `xml:"control-processor-health"`
	CoreBladeHealth             SSPStateType `xml:"core-blade-health"`
	BladeHealth                 SSPStateType `xml:"blade-health"`
	FlashHealth                 SSPStateType `xml:"flash-health"`
	MarginalPortHealth          SSPStateType `xml:"marginal-port-health"`
	FaultyPortHealth            SSPStateType `xml:"faulty-port-health"`
	MissingSFPHealth            SSPStateType `xml:"missing-sfp-health"`
	ErrorPortHealth             SSPStateType `xml:"error-port-health"`
	ExpiredCertificateHealth    SSPStateType `xml:"expired-certificate-health"`
	AirflowMismatchHealth       SSPStateType `xml:"airflow-mismatch-health"`
	MarginalSFPHealth           SSPStateType `xml:"marginal-sfp-health"`
	TrustedFOSCertificateHealth SSPStateType `xml:"trusted-fos-cert-health"`
}

// SwitchStatusPolicyReportInfo 是交换机健康状态策略报告的 JSON 友好表示，
// 由 GetSwitchStatusPolicyReport 从 XML 响应中转换而来
type SwitchStatusPolicyReportInfo struct {
	SwitchStatus                SSPStateType `json:"switch_status"`
	PowerSupplyHealth           SSPStateType `json:"power_supply_health"`
	FanHealth                   SSPStateType `json:"fan_health"`
	WWNHealth                   SSPStateType `json:"wwn_health"`
	TemperatureSensorHealth     SSPStateType `json:"temperature_sensor_health"`
	HAHealth                    SSPStateType `json:"ha_health"`
	ControlProcessorHealth      SSPStateType `json:"control_processor_health"`
	CoreBladeHealth             SSPStateType `json:"core_blade_health"`
	BladeHealth                 SSPStateType `json:"blade_health"`
	FlashHealth                 SSPStateType `json:"flash_health"`
	MarginalPortHealth          SSPStateType `json:"marginal_port_health"`
	FaultyPortHealth            SSPStateType `json:"faulty_port_health"`
	MissingSFPHealth            SSPStateType `json:"missing_sfp_health"`
	ErrorPortHealth             SSPStateType `json:"error_port_health"`
	ExpiredCertificateHealth    SSPStateType `json:"expired_certificate_health"`
	AirflowMismatchHealth       SSPStateType `json:"airflow_mismatch_health"`
	MarginalSFPHealth           SSPStateType `json:"marginal_sfp_health"`
	TrustedFOSCertificateHealth SSPStateType `json:"trusted_fos_certificate_health"`
}

// ==================== System Resources ====================

// SystemResourcesResponse 是 GET /brocade-maps/system-resources 的 XML 响应包装
type SystemResourcesResponse struct {
	XMLName         xml.Name        `xml:"Response"`
	SystemResources SystemResources `xml:"system-resources"`
}

// SystemResources 描述交换机系统资源使用情况（XML 原始结构），
// 包含 CPU、内存、Flash 使用率和可用内核内存。
// 对应 YANG 模型: brocade-maps/system-resources
type SystemResources struct {
	XMLName          xml.Name `xml:"system-resources"`
	CPUUsage         uint32   `xml:"cpu-usage"`
	MemoryUsage      uint32   `xml:"memory-usage"`
	TotalMemory      uint32   `xml:"total-memory"`
	FlashUsage       uint32   `xml:"flash-usage"`
	FreeKernelMemory uint32   `xml:"free-kernel-memory"`
}

// SystemResourcesInfo 是交换机系统资源使用情况的 JSON 友好表示，
// 由 GetSystemResources 从 XML 响应中转换而来
type SystemResourcesInfo struct {
	CPUUsage         uint32 `json:"cpu_usage"`
	MemoryUsage      uint32 `json:"memory_usage"`
	TotalMemory      uint32 `json:"total_memory"`
	FlashUsage       uint32 `json:"flash_usage"`
	FreeKernelMemory uint32 `json:"free_kernel_memory"`
}

// ==================== Client Methods ====================

// GetSwitchStatusPolicyReport 获取交换机各组件的健康状态策略报告，
// 包括交换机整体健康度、电源、风扇、温度传感器、端口、SFP 等状态。
// 对应 API: GET /brocade-maps/switch-status-policy-report
func (c *Client) GetSwitchStatusPolicyReport() (*SwitchStatusPolicyReportInfo, error) {
	var resp SwitchStatusPolicyReportResponse
	err := c.Get(c.endpoints().SwitchStatusPolicyReport(), &resp)
	if err != nil {
		return nil, err
	}

	return &SwitchStatusPolicyReportInfo{
		SwitchStatus:                resp.SwitchStatusPolicyReport.SwitchHealth,
		PowerSupplyHealth:           resp.SwitchStatusPolicyReport.PowerSupplyHealth,
		FanHealth:                   resp.SwitchStatusPolicyReport.FanHealth,
		WWNHealth:                   resp.SwitchStatusPolicyReport.WWNHealth,
		TemperatureSensorHealth:     resp.SwitchStatusPolicyReport.TemperatureSensorHealth,
		HAHealth:                    resp.SwitchStatusPolicyReport.HAHealth,
		ControlProcessorHealth:      resp.SwitchStatusPolicyReport.ControlProcessorHealth,
		CoreBladeHealth:             resp.SwitchStatusPolicyReport.CoreBladeHealth,
		BladeHealth:                 resp.SwitchStatusPolicyReport.BladeHealth,
		FlashHealth:                 resp.SwitchStatusPolicyReport.FlashHealth,
		MarginalPortHealth:          resp.SwitchStatusPolicyReport.MarginalPortHealth,
		FaultyPortHealth:            resp.SwitchStatusPolicyReport.FaultyPortHealth,
		MissingSFPHealth:            resp.SwitchStatusPolicyReport.MissingSFPHealth,
		ErrorPortHealth:             resp.SwitchStatusPolicyReport.ErrorPortHealth,
		ExpiredCertificateHealth:    resp.SwitchStatusPolicyReport.ExpiredCertificateHealth,
		AirflowMismatchHealth:       resp.SwitchStatusPolicyReport.AirflowMismatchHealth,
		MarginalSFPHealth:           resp.SwitchStatusPolicyReport.MarginalSFPHealth,
		TrustedFOSCertificateHealth: resp.SwitchStatusPolicyReport.TrustedFOSCertificateHealth,
	}, nil
}

// GetSystemResources 获取交换机系统资源使用情况，包括 CPU、内存、Flash 使用率等。
// 对应 API: GET /brocade-maps/system-resources
func (c *Client) GetSystemResources() (*SystemResourcesInfo, error) {
	var resp SystemResourcesResponse
	err := c.Get(c.endpoints().SystemResources(), &resp)
	if err != nil {
		return nil, err
	}

	return &SystemResourcesInfo{
		CPUUsage:         resp.SystemResources.CPUUsage,
		MemoryUsage:      resp.SystemResources.MemoryUsage,
		TotalMemory:      resp.SystemResources.TotalMemory,
		FlashUsage:       resp.SystemResources.FlashUsage,
		FreeKernelMemory: resp.SystemResources.FreeKernelMemory,
	}, nil
}
