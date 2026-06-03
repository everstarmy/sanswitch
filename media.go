package san

import "encoding/xml"

// MediaRDPResponse 是 GET /brocade-media/media-rdp 的 XML 响应包装
type MediaRDPResponse struct {
	XMLName   xml.Name       `xml:"Response"`
	MediaRDPs []MediaRDPInfo `xml:"media-rdp"`
}

// MediaRDPInfo 描述一个 SFP/XFP 光模块的完整诊断信息（RDP = Raw Diagnostic Parameters），
// 包含模块标识、速率能力、光纤类型、光功率、温度、电流、电压、
// 远端模块信息及各类告警标志等 100+ 个字段。
// 对应 YANG 模型: brocade-media/media-rdp
type MediaRDPInfo struct {
	XMLName                        xml.Name `xml:"media-rdp" json:"-"`
	Name                           string   `xml:"name" json:"name"`
	Identifier                     string   `xml:"identifier" json:"identifier"`
	Connector                      string   `xml:"connector" json:"connector"`
	NominalBaudRate                int32    `xml:"nominal-baud-rate" json:"nominal_baud_rate"`
	BaudRateMaximum                int32    `xml:"baud-rate-maximum" json:"baud_rate_maximum"`
	BaudRateMinimum                int32    `xml:"baud-rate-minimum" json:"baud_rate_minimum"`
	OpticalCableTypes              []string `xml:"optical-cable-types>cable-type" json:"optical_cable_types"`
	TransmissionType               string   `xml:"transmission-type" json:"transmission_type"`
	Fiber9umLinkLength             uint32   `xml:"fiber-9um-link-length" json:"fiber_9um_link_length"`
	Fiber9mLinkLength              uint32   `xml:"fiber-9m-link-length" json:"fiber_9m_link_length"`
	MultimodeFiberOM250umLength    uint32   `xml:"multimode-fiber-om2-50um-length" json:"multimode_fiber_om2_50um_length"`
	MultimodeFiberOM350umLength    uint32   `xml:"multimode-fiber-om3-50um-length" json:"multimode_fiber_om3_50um_length"`
	MultimodeFiberOM450umLength    uint32   `xml:"multimode-fiber-om4-50um-length" json:"multimode_fiber_om4_50um_length"`
	MultimodeFiber625umLength      uint32   `xml:"multimode-fiber-625um-length" json:"multimode_fiber_625um_length"`
	CopperLinkLength               uint32   `xml:"copper-link-length" json:"copper_link_length"`
	EnhancedOptions                string   `xml:"enhanced-options" json:"enhanced_options"`
	DigitalDiagnosticsType         string   `xml:"digital-diagnostics-type" json:"digital_diagnostics_type"`
	MediaOptionalSignals           []string `xml:"media-optional-signals>signal" json:"media_optional_signals"`
	MediaOptionalSignalValue       string   `xml:"media-optional-signal-value" json:"media_optional_signal_value"`
	StatusControl                  string   `xml:"status-control" json:"status_control"`
	ElectricLoopbackEnabled        bool     `xml:"electric-loopback-enabled" json:"electric_loopback_enabled"`
	OpticalLoopbackEnabled         bool     `xml:"optical-loopback-enabled" json:"optical_loopback_enabled"`
	MediaSpeedCapability           []string `xml:"media-speed-capability>speed" json:"media_speed_capability"`
	MediaDistance                  []string `xml:"media-distance>distance" json:"media_distance"`
	Encoding                       string   `xml:"encoding" json:"encoding"`
	VendorOUI                      string   `xml:"vendor-oui" json:"vendor_oui"`
	PartNumber                     string   `xml:"part-number" json:"part_number"`
	SerialNumber                   string   `xml:"serial-number" json:"serial_number"`
	VendorName                     string   `xml:"vendor-name" json:"vendor_name"`
	VendorRevision                 string   `xml:"vendor-revision" json:"vendor_revision"`
	DateCode                       string   `xml:"date-code" json:"date_code"`
	Temperature                    float64  `xml:"temperature" json:"temperature"`
	RxPower                        float64  `xml:"rx-power" json:"rx_power"`
	TxPower                        float64  `xml:"tx-power" json:"tx_power"`
	Current                        float64  `xml:"current" json:"current"`
	Voltage                        float64  `xml:"voltage" json:"voltage"`
	Wavelength                     uint32   `xml:"wavelength" json:"wavelength"`
	PowerOnTime                    int32    `xml:"power-on-time" json:"power_on_time"`
	AlarmFlagsHighAlarm            string   `xml:"alarm-flags>high-alarm" json:"alarm_flags_high_alarm"`
	AlarmFlagsLowAlarm             string   `xml:"alarm-flags>low-alarm" json:"alarm_flags_low_alarm"`
	WarningFlagsHighWarning        string   `xml:"warning-flags>high-warning" json:"warning_flags_high_warning"`
	WarningFlagsLowWarning         string   `xml:"warning-flags>low-warning" json:"warning_flags_low_warning"`
	StateTransitions               uint32   `xml:"state-transitions" json:"state_transitions"`
	RoundTripLinkLatency           int32    `xml:"round-trip-link-latency" json:"round_trip_link_latency"`
	LastPollTime                   string   `xml:"last-poll-time" json:"last_poll_time"`
	PeerDataAvailable              bool     `xml:"peer-data-available" json:"peer_data_available"`
	RemoteIdentifier               string   `xml:"remote-identifier" json:"remote_identifier"`
	RemoteLaserType                string   `xml:"remote-laser-type" json:"remote_laser_type"`
	RemoteOpticalType              string   `xml:"remote-optical-type" json:"remote_optical_type"`
	RemoteMediaSpeedCapability     []string `xml:"remote-media-speed-capability>speed" json:"remote_media_speed_capability"`
	RemoteOpticalProductDataPN     string   `xml:"remote-optical-product-data>part-number" json:"remote_optical_product_data_part_number"`
	RemoteOpticalProductDataSN     string   `xml:"remote-optical-product-data>serial-number" json:"remote_optical_product_data_serial_number"`
	RemoteOpticalProductDataVN     string   `xml:"remote-optical-product-data>vendor-name" json:"remote_optical_product_data_vendor_name"`
	RemoteOpticalProductDataVR     string   `xml:"remote-optical-product-data>vendor-revision" json:"remote_optical_product_data_vendor_revision"`
	RemoteOpticalProductDataDC     string   `xml:"remote-optical-product-data>date-code" json:"remote_optical_product_data_date_code"`
	RemoteMediaTemperature         float64  `xml:"remote-media-temperature" json:"remote_media_temperature"`
	RemoteMediaRxPower             float64  `xml:"remote-media-rx-power" json:"remote_media_rx_power"`
	RemoteMediaTxPower             float64  `xml:"remote-media-tx-power" json:"remote_media_tx_power"`
	RemoteMediaCurrent             float64  `xml:"remote-media-current" json:"remote_media_current"`
	RemoteMediaVoltage             float64  `xml:"remote-media-voltage" json:"remote_media_voltage"`
	RemoteMediaVoltageAlertHA      string   `xml:"remote-media-voltage-alert>high-alarm" json:"remote_media_voltage_alert_high_alarm"`
	RemoteMediaVoltageAlertLA      string   `xml:"remote-media-voltage-alert>low-alarm" json:"remote_media_voltage_alert_low_alarm"`
	RemoteMediaVoltageAlertHW      string   `xml:"remote-media-voltage-alert>high-warning" json:"remote_media_voltage_alert_high_warning"`
	RemoteMediaVoltageAlertLW      string   `xml:"remote-media-voltage-alert>low-warning" json:"remote_media_voltage_alert_low_warning"`
	RemoteMediaSignalLossUS        string   `xml:"remote-media-signal-loss>up-stream" json:"remote_media_signal_loss_up_stream"`
	RemoteMediaSignalLossDS        string   `xml:"remote-media-signal-loss>down-stream" json:"remote_media_signal_loss_down_stream"`
	RemoteMediaSignalLossV2US      float64  `xml:"remote-media-signal-loss-v2>up-stream-v2" json:"remote_media_signal_loss_v2_up_stream"`
	RemoteMediaSignalLossV2DS      float64  `xml:"remote-media-signal-loss-v2>down-stream-v2" json:"remote_media_signal_loss_v2_down_stream"`
	RemoteMediaTemperatureAlertHA  string   `xml:"remote-media-temperature-alert>high-alarm" json:"remote_media_temperature_alert_high_alarm"`
	RemoteMediaTemperatureAlertLA  string   `xml:"remote-media-temperature-alert>low-alarm" json:"remote_media_temperature_alert_low_alarm"`
	RemoteMediaTemperatureAlertHW  string   `xml:"remote-media-temperature-alert>high-warning" json:"remote_media_temperature_alert_high_warning"`
	RemoteMediaTemperatureAlertLW  string   `xml:"remote-media-temperature-alert>low-warning" json:"remote_media_temperature_alert_low_warning"`
	RemoteMediaTXBiasAlertHA       string   `xml:"remote-media-tx-bias-alert>high-alarm" json:"remote_media_tx_bias_alert_high_alarm"`
	RemoteMediaTXBiasAlertLA       string   `xml:"remote-media-tx-bias-alert>low-alarm" json:"remote_media_tx_bias_alert_low_alarm"`
	RemoteMediaTXBiasAlertHW       string   `xml:"remote-media-tx-bias-alert>high-warning" json:"remote_media_tx_bias_alert_high_warning"`
	RemoteMediaTXBiasAlertLW       string   `xml:"remote-media-tx-bias-alert>low-warning" json:"remote_media_tx_bias_alert_low_warning"`
	RemoteMediaTXPowerAlertHA      string   `xml:"remote-media-tx-power-alert>high-alarm" json:"remote_media_tx_power_alert_high_alarm"`
	RemoteMediaTXPowerAlertLA      string   `xml:"remote-media-tx-power-alert>low-alarm" json:"remote_media_tx_power_alert_low_alarm"`
	RemoteMediaTXPowerAlertHW      string   `xml:"remote-media-tx-power-alert>high-warning" json:"remote_media_tx_power_alert_high_warning"`
	RemoteMediaTXPowerAlertLW      string   `xml:"remote-media-tx-power-alert>low-warning" json:"remote_media_tx_power_alert_low_warning"`
	RemoteMediaRXPowerAlertHA      string   `xml:"remote-media-rx-power-alert>high-alarm" json:"remote_media_rx_power_alert_high_alarm"`
	RemoteMediaRXPowerAlertLA      string   `xml:"remote-media-rx-power-alert>low-alarm" json:"remote_media_rx_power_alert_low_alarm"`
	RemoteMediaRXPowerAlertHW      string   `xml:"remote-media-rx-power-alert>high-warning" json:"remote_media_rx_power_alert_high_warning"`
	RemoteMediaRXPowerAlertLW      string   `xml:"remote-media-rx-power-alert>low-warning" json:"remote_media_rx_power_alert_low_warning"`
	RemoteMediaVoltageAlarmTypeHA  string   `xml:"remote-media-voltage-alarm-type>high-alarm" json:"remote_media_voltage_alarm_type_high_alarm"`
	RemoteMediaVoltageAlarmTypeLA  string   `xml:"remote-media-voltage-alarm-type>low-alarm" json:"remote_media_voltage_alarm_type_low_alarm"`
	RemoteMediaVoltageAlarmTypeHW  string   `xml:"remote-media-voltage-alarm-type>high-warning" json:"remote_media_voltage_alarm_type_high_warning"`
	RemoteMediaVoltageAlarmTypeLW  string   `xml:"remote-media-voltage-alarm-type>low-warning" json:"remote_media_voltage_alarm_type_low_warning"`
	RemoteMediaTemperatureAlarmTHA string   `xml:"remote-media-temperature-alarm-type>high-alarm" json:"remote_media_temperature_alarm_type_high_alarm"`
	RemoteMediaTemperatureAlarmTLA string   `xml:"remote-media-temperature-alarm-type>low-alarm" json:"remote_media_temperature_alarm_type_low_alarm"`
	RemoteMediaTemperatureAlarmTHW string   `xml:"remote-media-temperature-alarm-type>high-warning" json:"remote_media_temperature_alarm_type_high_warning"`
	RemoteMediaTemperatureAlarmTLW string   `xml:"remote-media-temperature-alarm-type>low-warning" json:"remote_media_temperature_alarm_type_low_warning"`
	RemoteMediaTXBiasAlarmTypeHA   string   `xml:"remote-media-tx-bias-alarm-type>high-alarm" json:"remote_media_tx_bias_alarm_type_high_alarm"`
	RemoteMediaTXBiasAlarmTypeLA   string   `xml:"remote-media-tx-bias-alarm-type>low-alarm" json:"remote_media_tx_bias_alarm_type_low_alarm"`
	RemoteMediaTXBiasAlarmTypeHW   string   `xml:"remote-media-tx-bias-alarm-type>high-warning" json:"remote_media_tx_bias_alarm_type_high_warning"`
	RemoteMediaTXBiasAlarmTypeLW   string   `xml:"remote-media-tx-bias-alarm-type>low-warning" json:"remote_media_tx_bias_alarm_type_low_warning"`
	RemoteMediaTXPowerAlarmTypeHA  string   `xml:"remote-media-tx-power-alarm-type>high-alarm" json:"remote_media_tx_power_alarm_type_high_alarm"`
	RemoteMediaTXPowerAlarmTypeLA  string   `xml:"remote-media-tx-power-alarm-type>low-alarm" json:"remote_media_tx_power_alarm_type_low_alarm"`
	RemoteMediaTXPowerAlarmTypeHW  string   `xml:"remote-media-tx-power-alarm-type>high-warning" json:"remote_media_tx_power_alarm_type_high_warning"`
	RemoteMediaTXPowerAlarmTypeLW  string   `xml:"remote-media-tx-power-alarm-type>low-warning" json:"remote_media_tx_power_alarm_type_low_warning"`
	RemoteMediaRXPowerAlarmTypeHA  string   `xml:"remote-media-rx-power-alarm-type>high-alarm" json:"remote_media_rx_power_alarm_type_high_alarm"`
	RemoteMediaRXPowerAlarmTypeLA  string   `xml:"remote-media-rx-power-alarm-type>low-alarm" json:"remote_media_rx_power_alarm_type_low_alarm"`
	RemoteMediaRXPowerAlarmTypeHW  string   `xml:"remote-media-rx-power-alarm-type>high-warning" json:"remote_media_rx_power_alarm_type_high_warning"`
	RemoteMediaRXPowerAlarmTypeLW  string   `xml:"remote-media-rx-power-alarm-type>low-warning" json:"remote_media_rx_power_alarm_type_low_warning"`
}

// GetMediaRDPs 获取交换机上所有 SFP/XFP 光模块的原始诊断参数信息。
// 对应 API: GET /brocade-media/media-rdp
func (c *Client) GetMediaRDPs() ([]MediaRDPInfo, error) {
	var resp MediaRDPResponse
	err := c.Get("/brocade-media/media-rdp", &resp)
	if err != nil {
		return nil, err
	}
	return resp.MediaRDPs, nil
}
