package san

import "encoding/xml"

// ==================== Time Zone ====================

// TimeZoneResponse 对应 GET /brocade-time/time-zone
type TimeZoneResponse struct {
	XMLName  xml.Name     `xml:"Response"`
	TimeZone TimeZoneInfo `xml:"time-zone"`
}

// TimeZoneInfo 描述交换机时区配置
type TimeZoneInfo struct {
	XMLName          xml.Name `xml:"time-zone" json:"-"`
	Name             string   `xml:"name" json:"name"`
	GMTOffsetHours   int16    `xml:"gmt-offset-hours" json:"gmt_offset_hours"`
	GMTOffsetMinutes int16    `xml:"gmt-offset-minutes" json:"gmt_offset_minutes"`
}

// ==================== Clock Server (NTP) ====================

// ClockServerResponse 对应 GET /brocade-time/clock-server
type ClockServerResponse struct {
	XMLName     xml.Name        `xml:"Response"`
	ClockServer ClockServerInfo `xml:"clock-server"`
}

// ClockServerInfo 描述交换机 NTP 时钟服务器配置
type ClockServerInfo struct {
	XMLName            xml.Name `xml:"clock-server" json:"-"`
	NTPServerAddresses []string `xml:"ntp-server-address>server-address" json:"ntp_server_addresses"`
	ActiveServer       string   `xml:"active-server" json:"active_server"`
	TSAuthSpec         string   `xml:"ts-auth-spec" json:"ts_auth_spec"`
	TSLegacyMode       bool     `xml:"ts-legacy-mode" json:"ts_legacy_mode"`
}

// ==================== Client Methods ====================

// GetTimeZone 获取交换机时区配置
func (c *Client) GetTimeZone() (*TimeZoneInfo, error) {
	var resp TimeZoneResponse
	err := c.Get("/brocade-time/time-zone", &resp)
	if err != nil {
		return nil, err
	}
	return &resp.TimeZone, nil
}

// GetClockServer 获取交换机 NTP 时钟服务器配置
func (c *Client) GetClockServer() (*ClockServerInfo, error) {
	var resp ClockServerResponse
	err := c.Get("/brocade-time/clock-server", &resp)
	if err != nil {
		return nil, err
	}
	return &resp.ClockServer, nil
}
