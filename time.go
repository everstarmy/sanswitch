package sanswitch

import (
	"context"
	"encoding/xml"
)

// ==================== Time Zone ====================

// timeZoneResponse 对应 GET /brocade-time/time-zone
type timeZoneResponse struct {
	XMLName  xml.Name `xml:"Response"`
	TimeZone TimeZone `xml:"time-zone"`
}

// TimeZone 描述交换机时区配置
type TimeZone struct {
	XMLName          xml.Name `xml:"time-zone" json:"-"`
	Name             string   `xml:"name" json:"name"`
	GMTOffsetHours   int16    `xml:"gmt-offset-hours" json:"gmt_offset_hours"`
	GMTOffsetMinutes int16    `xml:"gmt-offset-minutes" json:"gmt_offset_minutes"`
}

// ==================== Clock Server (NTP) ====================

// clockServerResponse 对应 GET /brocade-time/clock-server
type clockServerResponse struct {
	XMLName     xml.Name    `xml:"Response"`
	ClockServer ClockServer `xml:"clock-server"`
}

// ClockServer 描述交换机 NTP 时钟服务器配置
type ClockServer struct {
	XMLName            xml.Name `xml:"clock-server" json:"-"`
	NTPServerAddresses []string `xml:"ntp-server-address>server-address" json:"ntp_server_addresses"`
	ActiveServer       string   `xml:"active-server" json:"active_server"`
	TSAuthSpec         string   `xml:"ts-auth-spec" json:"ts_auth_spec"`
	TSLegacyMode       bool     `xml:"ts-legacy-mode" json:"ts_legacy_mode"`
}

// ==================== Client Methods ====================

// TimeZone 获取交换机时区配置并允许取消请求。
func (c *Client) TimeZone(ctx context.Context) (*TimeZone, error) {
	var resp timeZoneResponse
	err := c.get(ctx, c.endpoints().TimeZone(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.TimeZone, nil
}

// ClockServer 获取交换机 NTP 时钟服务器配置并允许取消请求。
func (c *Client) ClockServer(ctx context.Context) (*ClockServer, error) {
	var resp clockServerResponse
	err := c.get(ctx, c.endpoints().ClockServer(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.ClockServer, nil
}
