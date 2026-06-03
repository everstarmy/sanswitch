package san

import "encoding/xml"

// ==================== Blade ====================

// BladeResponse 是 GET /brocade-fru/blade 的 XML 响应包装
type BladeResponse struct {
	XMLName xml.Name    `xml:"Response"`
	Blades  []BladeInfo `xml:"blade"`
}

// BladeInfo 描述一个 FRU 板卡的详细信息，包含插槽号、类型、状态、固件版本、
// 端口数量、功耗、温度、网络地址等字段。
// 对应 YANG 模型: brocade-fru/blade
type BladeInfo struct {
	XMLName                  xml.Name `xml:"blade" json:"-"`
	SlotNumber               uint32   `xml:"slot-number" json:"slot_number"`
	BladeType                string   `xml:"blade-type" json:"blade_type"`
	BladeTypeString          string   `xml:"blade-type-string" json:"blade_type_string"`
	BladeID                  uint16   `xml:"blade-id" json:"blade_id"`
	BladeState               string   `xml:"blade-state" json:"blade_state"`
	BladeStateString         string   `xml:"blade-state-string" json:"blade_state_string"`
	ModelName                string   `xml:"model-name" json:"model_name"`
	PartNumber               string   `xml:"part-number" json:"part_number"`
	SerialNumber             string   `xml:"serial-number" json:"serial_number"`
	VendorSerialNumber       string   `xml:"vendor-serial-number" json:"vendor_serial_number"`
	VendorPartNumber         string   `xml:"vendor-part-number" json:"vendor_part_number"`
	VendorRevisionNumber     string   `xml:"vendor-revision-number" json:"vendor_revision_number"`
	VendorID                 string   `xml:"vendor-id" json:"vendor_id"`
	Manufacturer             string   `xml:"manufacturer" json:"manufacturer"`
	ManufactureDate          string   `xml:"manufacture-date" json:"manufacture_date"`
	LastHeaderUpdateDate     string   `xml:"last-header-update-date" json:"last_header_update_date"`
	HeaderVersion            uint16   `xml:"header-version" json:"header_version"`
	FirmwareVersion          string   `xml:"firmware-version" json:"firmware_version"`
	PrimaryFirmwareVersion   string   `xml:"primary-firmware-version" json:"primary_firmware_version"`
	SecondaryFirmwareVersion string   `xml:"secondary-firmware-version" json:"secondary_firmware_version"`
	SoftwareVersion          string   `xml:"software-version" json:"software_version"`
	BootVersion              string   `xml:"boot-version" json:"boot_version"`
	HardwareVersion          string   `xml:"hardware-version" json:"hardware_version"`
	FCPortCount              uint16   `xml:"fc-port-count" json:"fc_port_count"`
	GEPortCount              uint16   `xml:"ge-port-count" json:"ge_port_count"`
	PowerUsage               uint32   `xml:"power-usage" json:"power_usage"`
	PowerUsageWatts          float32  `xml:"power-usage-watts" json:"power_usage_watts"`
	PowerCapability          uint32   `xml:"power-capability" json:"power_capability"`
	PowerCapabilityWatts     float32  `xml:"power-capability-watts" json:"power_capability_watts"`
	PowerConsumption         uint32   `xml:"power-consumption" json:"power_consumption"`
	Temperature              uint32   `xml:"temperature" json:"temperature"`
	TemperatureCelsius       int16    `xml:"temperature-celsius" json:"temperature_celsius"`
	TimeAlive                uint32   `xml:"time-alive" json:"time_alive"`
	TimeAwake                uint32   `xml:"time-awake" json:"time_awake"`
	OperationalStatus        uint32   `xml:"operational-status" json:"operational_status"`
	OperationalStatusString  string   `xml:"operational-status-string" json:"operational_status_string"`
	DiagnosticsPostStatus    string   `xml:"diagnostics-post-status" json:"diagnostics_post_status"`
	StatusLED                string   `xml:"status-led" json:"status_led"`
	ExtensionEnabled         bool     `xml:"extension-enabled" json:"extension_enabled"`
	ExtensionAppMode         string   `xml:"extension-app-mode" json:"extension_app_mode"`
	ExtensionVEMode          string   `xml:"extension-ve-mode" json:"extension_ve_mode"`
	ExtensionGEMode          string   `xml:"extension-ge-mode" json:"extension_ge_mode"`
	Presence                 bool     `xml:"presence" json:"presence"`
	HotSwappable             bool     `xml:"hot-swappable" json:"hot_swappable"`
	HotPluggable             bool     `xml:"hot-pluggable" json:"hot_pluggable"`
	PowerOn                  bool     `xml:"power-on" json:"power_on"`
	PowerOff                 bool     `xml:"power-off" json:"power_off"`
	PostStatus               string   `xml:"post-status" json:"post_status"`
	LinkLocalAddress         string   `xml:"link-local-address" json:"link_local_address"`
	StaticIPAddresses        []string `xml:"static-ip-addresses>ip-address" json:"static_ip_addresses"`
	StatefulIPAddresses      []string `xml:"stateful-ip-addresses>ip-address" json:"stateful_ip_addresses"`
	StatelessIPv6Addresses   []string `xml:"stateless-ipv6-addresses>ip-address" json:"stateless_ipv6_addresses"`
}

// ==================== Fan ====================

// FanResponse 是 GET /brocade-fru/fan 的 XML 响应包装
type FanResponse struct {
	XMLName xml.Name  `xml:"Response"`
	Fans    []FanInfo `xml:"fan"`
}

// FanInfo 描述一个风扇单元的详细信息，包含类型、状态、风向、转速、
// 功耗、序列号等字段。
// 对应 YANG 模型: brocade-fru/fan
type FanInfo struct {
	XMLName                 xml.Name `xml:"fan" json:"-"`
	UnitNumber              uint32   `xml:"unit-number" json:"unit_number"`
	FanType                 string   `xml:"fan-type" json:"fan_type"`
	FanTypeString           string   `xml:"fan-type-string" json:"fan_type_string"`
	FanState                string   `xml:"fan-state" json:"fan_state"`
	FanStateString          string   `xml:"fan-state-string" json:"fan_state_string"`
	FanDirection            string   `xml:"fan-direction" json:"fan_direction"`
	FanDirectionString      string   `xml:"fan-direction-string" json:"fan_direction_string"`
	AirflowDirection        string   `xml:"airflow-direction" json:"airflow_direction"`
	OperationalState        string   `xml:"operational-state" json:"operational_state"`
	PartNumber              string   `xml:"part-number" json:"part_number"`
	SerialNumber            string   `xml:"serial-number" json:"serial_number"`
	VendorSerialNumber      string   `xml:"vendor-serial-number" json:"vendor_serial_number"`
	VendorPartNumber        string   `xml:"vendor-part-number" json:"vendor_part_number"`
	VendorRevisionNumber    string   `xml:"vendor-revision-number" json:"vendor_revision_number"`
	VendorID                string   `xml:"vendor-id" json:"vendor_id"`
	Manufacturer            string   `xml:"manufacturer" json:"manufacturer"`
	ManufactureDate         string   `xml:"manufacture-date" json:"manufacture_date"`
	LastHeaderUpdateDate    string   `xml:"last-header-update-date" json:"last_header_update_date"`
	HeaderVersion           uint16   `xml:"header-version" json:"header_version"`
	Speed                   uint32   `xml:"speed" json:"speed"`
	SpeedRPM                uint32   `xml:"speed-rpm" json:"speed_rpm"`
	AirFlow                 uint32   `xml:"air-flow" json:"air_flow"`
	PowerConsumption        uint32   `xml:"power-consumption" json:"power_consumption"`
	TimeAlive               uint32   `xml:"time-alive" json:"time_alive"`
	TimeAwake               uint32   `xml:"time-awake" json:"time_awake"`
	OperationalStatus       uint32   `xml:"operational-status" json:"operational_status"`
	OperationalStatusString string   `xml:"operational-status-string" json:"operational_status_string"`
	StatusLED               string   `xml:"status-led" json:"status_led"`
	Presence                bool     `xml:"presence" json:"presence"`
	HotSwappable            bool     `xml:"hot-swappable" json:"hot_swappable"`
}

// ==================== Power Supply ====================

// PowerSupplyResponse 是 GET /brocade-fru/power-supply 的 XML 响应包装
type PowerSupplyResponse struct {
	XMLName       xml.Name          `xml:"Response"`
	PowerSupplies []PowerSupplyInfo `xml:"power-supply"`
}

// PowerSupplyInfo 描述一个电源单元的详细信息，包含类型、状态、输入/输出电压电流、
// 功耗、温度、风扇转速、序列号等字段。
// 对应 YANG 模型: brocade-fru/power-supply
type PowerSupplyInfo struct {
	XMLName                    xml.Name `xml:"power-supply" json:"-"`
	UnitNumber                 uint32   `xml:"unit-number" json:"unit_number"`
	PowerSupplyType            string   `xml:"power-supply-type" json:"power_supply_type"`
	PowerSupplyTypeString      string   `xml:"power-supply-type-string" json:"power_supply_type_string"`
	PowerSupplyState           string   `xml:"power-supply-state" json:"power_supply_state"`
	PowerSupplyStateString     string   `xml:"power-supply-state-string" json:"power_supply_state_string"`
	OperationalState           string   `xml:"operational-state" json:"operational_state"`
	PartNumber                 string   `xml:"part-number" json:"part_number"`
	SerialNumber               string   `xml:"serial-number" json:"serial_number"`
	VendorSerialNumber         string   `xml:"vendor-serial-number" json:"vendor_serial_number"`
	VendorPartNumber           string   `xml:"vendor-part-number" json:"vendor_part_number"`
	VendorRevisionNumber       string   `xml:"vendor-revision-number" json:"vendor_revision_number"`
	VendorID                   string   `xml:"vendor-id" json:"vendor_id"`
	Manufacturer               string   `xml:"manufacturer" json:"manufacturer"`
	ManufactureDate            string   `xml:"manufacture-date" json:"manufacture_date"`
	LastHeaderUpdateDate       string   `xml:"last-header-update-date" json:"last_header_update_date"`
	HeaderVersion              uint16   `xml:"header-version" json:"header_version"`
	PowerProduction            uint32   `xml:"power-production" json:"power_production"`
	PowerSource                string   `xml:"power-source" json:"power_source"`
	AirflowDirection           string   `xml:"airflow-direction" json:"airflow_direction"`
	TemperatureSensorSupported bool     `xml:"temperature-sensor-supported" json:"temperature_sensor_supported"`
	InputVoltage               float32  `xml:"input-voltage" json:"input_voltage"`
	InputVoltageMillivolts     float64  `xml:"input-voltage-millivolts" json:"input_voltage_millivolts"`
	InputCurrent               uint32   `xml:"input-current" json:"input_current"`
	InputCurrentMilliamps      uint32   `xml:"input-current-milliamps" json:"input_current_milliamps"`
	InputPower                 uint32   `xml:"input-power" json:"input_power"`
	InputPowerMilliwatts       uint32   `xml:"input-power-milliwatts" json:"input_power_milliwatts"`
	OutputVoltage              uint32   `xml:"output-voltage" json:"output_voltage"`
	OutputVoltageMillivolts    uint32   `xml:"output-voltage-millivolts" json:"output_voltage_millivolts"`
	OutputCurrent              uint32   `xml:"output-current" json:"output_current"`
	OutputCurrentMilliamps     uint32   `xml:"output-current-milliamps" json:"output_current_milliamps"`
	OutputPower                uint32   `xml:"output-power" json:"output_power"`
	OutputPowerMilliwatts      uint32   `xml:"output-power-milliwatts" json:"output_power_milliwatts"`
	PowerUsage                 int32    `xml:"power-usage" json:"power_usage"`
	Temperature                float64  `xml:"temperature" json:"temperature"`
	TemperatureCelsius         int16    `xml:"temperature-celsius" json:"temperature_celsius"`
	FanSpeed                   uint32   `xml:"fan-speed" json:"fan_speed"`
	FanSpeedRPM                uint32   `xml:"fan-speed-rpm" json:"fan_speed_rpm"`
	TimeAlive                  uint32   `xml:"time-alive" json:"time_alive"`
	TimeAwake                  uint32   `xml:"time-awake" json:"time_awake"`
	OperationalStatus          uint32   `xml:"operational-status" json:"operational_status"`
	OperationalStatusString    string   `xml:"operational-status-string" json:"operational_status_string"`
	StatusLED                  string   `xml:"status-led" json:"status_led"`
	Presence                   bool     `xml:"presence" json:"presence"`
	HotSwappable               bool     `xml:"hot-swappable" json:"hot_swappable"`
}

// ==================== History Log ====================

// HistoryLogResponse 是 GET /brocade-fru/history-log 的 XML 响应包装
type HistoryLogResponse struct {
	XMLName     xml.Name         `xml:"Response"`
	HistoryLogs []HistoryLogInfo `xml:"history-log"`
}

// HistoryLogInfo 描述一条 FRU 历史日志记录，包含 FRU 类型、位置、状态、
// 时间戳、序列号和部件号。
// 对应 YANG 模型: brocade-fru/history-log
type HistoryLogInfo struct {
	XMLName      xml.Name `xml:"history-log" json:"-"`
	FRUType      string   `xml:"fru-type" json:"fru_type"`
	Position     uint16   `xml:"position" json:"position"`
	State        string   `xml:"state" json:"state"`
	TimeStamp    string   `xml:"time-stamp" json:"time_stamp"`
	SerialNumber string   `xml:"serial-number" json:"serial_number"`
	PartNumber   string   `xml:"part-number" json:"part_number"`
}

// ==================== Sensor ====================

// SensorResponse 是 GET /brocade-fru/sensor 的 XML 响应包装
type SensorResponse struct {
	XMLName xml.Name     `xml:"Response"`
	Sensors []SensorInfo `xml:"sensor"`
}

// SensorInfo 描述一个传感器（如温度传感器）的详细信息，包含 ID、插槽号、
// 状态、类别、温度值及告警阈值等字段。
// 对应 YANG 模型: brocade-fru/sensor
type SensorInfo struct {
	XMLName                 xml.Name `xml:"sensor" json:"-"`
	ID                      uint16   `xml:"id" json:"id"`
	SlotNumber              uint16   `xml:"slot-number" json:"slot_number"`
	Index                   uint16   `xml:"index" json:"index"`
	State                   string   `xml:"state" json:"state"`
	Category                string   `xml:"category" json:"category"`
	Temperature             uint16   `xml:"temperature" json:"temperature"`
	TemperatureInFahrenheit int32    `xml:"temperature-in-fahrenheit" json:"temperature_in_fahrenheit"`
	// Legacy fields (kept for backward compatibility with older FOS versions)
	Name                    string `xml:"name" json:"name"`
	Type                    string `xml:"sensor-type" json:"type"`
	TypeString              string `xml:"sensor-type-string" json:"type_string"`
	Value                   int32  `xml:"sensor-value" json:"value"`
	Status                  string `xml:"sensor-status" json:"status"`
	StatusString            string `xml:"sensor-status-string" json:"status_string"`
	Unit                    string `xml:"sensor-unit" json:"unit"`
	UnitString              string `xml:"sensor-unit-string" json:"unit_string"`
	LowThreshold            int32  `xml:"low-threshold" json:"low_threshold"`
	HighThreshold           int32  `xml:"high-threshold" json:"high_threshold"`
	LowWarningThreshold     int32  `xml:"low-warning-threshold" json:"low_warning_threshold"`
	HighWarningThreshold    int32  `xml:"high-warning-threshold" json:"high_warning_threshold"`
	OperationalStatus       uint32 `xml:"operational-status" json:"operational_status"`
	OperationalStatusString string `xml:"operational-status-string" json:"operational_status_string"`
}

// ==================== Client Methods ====================

// GetBlades 获取交换机上所有 FRU 板卡的详细信息。
// 对应 API: GET /brocade-fru/blade
func (c *Client) GetBlades() ([]BladeInfo, error) {
	var resp BladeResponse
	err := c.Get("/brocade-fru/blade", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Blades, nil
}

// GetFans 获取交换机上所有风扇单元的详细信息。
// 对应 API: GET /brocade-fru/fan
func (c *Client) GetFans() ([]FanInfo, error) {
	var resp FanResponse
	err := c.Get("/brocade-fru/fan", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Fans, nil
}

// GetPowerSupplies 获取交换机上所有电源单元的详细信息。
// 对应 API: GET /brocade-fru/power-supply
func (c *Client) GetPowerSupplies() ([]PowerSupplyInfo, error) {
	var resp PowerSupplyResponse
	err := c.Get("/brocade-fru/power-supply", &resp)
	if err != nil {
		return nil, err
	}
	return resp.PowerSupplies, nil
}

// GetHistoryLogs 获取 FRU 组件的历史日志记录。
// 对应 API: GET /brocade-fru/history-log
func (c *Client) GetHistoryLogs() ([]HistoryLogInfo, error) {
	var resp HistoryLogResponse
	err := c.Get("/brocade-fru/history-log", &resp)
	if err != nil {
		return nil, err
	}
	return resp.HistoryLogs, nil
}

// GetSensors 获取交换机上所有传感器的详细信息。
// 对应 API: GET /brocade-fru/sensor
func (c *Client) GetSensors() ([]SensorInfo, error) {
	var resp SensorResponse
	err := c.Get("/brocade-fru/sensor", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Sensors, nil
}
