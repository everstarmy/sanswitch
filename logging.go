package san

import (
	"encoding/xml"
	"fmt"
)

// ==================== Error Log (RAS Log) ====================

// ErrorLogResponse 是 GET /brocade-logging/error-log 的 XML 响应包装
type ErrorLogResponse struct {
	XMLName   xml.Name       `xml:"Response"`
	ErrorLogs []ErrorLogInfo `xml:"error-log"`
}

// ErrorLogInfo 描述一条 RAS 错误日志记录，包含序列号、时间戳、消息 ID、
// 严重级别、Fabric ID、插槽 ID、消息文本等字段。
// 对应 YANG 模型: brocade-logging/error-log
type ErrorLogInfo struct {
	XMLName                xml.Name `xml:"error-log" json:"-"`
	SequenceNumber         uint32   `xml:"sequence-number" json:"sequence_number"`
	TimeStamp              string   `xml:"time-stamp" json:"time_stamp"`
	MessageID              string   `xml:"message-id" json:"message_id"`
	FabricID               uint32   `xml:"fabric-id" json:"fabric_id"`
	SlotID                 uint32   `xml:"slot-id" json:"slot_id"`
	SeverityLevel          string   `xml:"severity-level" json:"severity_level"`
	FFDCGeneratedEvent     bool     `xml:"ffdc-generated-event" json:"ffdc_generated_event"`
	SwitchUserFriendlyName string   `xml:"switch-user-friendly-name" json:"switch_user_friendly_name"`
	MessageText            string   `xml:"message-text" json:"message_text"`
	EventInfo              string   `xml:"event-info" json:"event_info"`
}

// GetErrorLogs 获取交换机上的所有 RAS 错误日志记录。
// 对应 API: GET /brocade-logging/error-log
func (c *Client) GetErrorLogs() ([]ErrorLogInfo, error) {
	if err := c.ensureLoggingSupported("error-log"); err != nil {
		return nil, err
	}

	var resp ErrorLogResponse
	err := c.Get(c.endpoints().ErrorLogs(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.ErrorLogs, nil
}

// ==================== Audit Log ====================

// AuditLogResponse 是 GET /brocade-logging/audit-log 的 XML 响应包装
type AuditLogResponse struct {
	XMLName   xml.Name       `xml:"Response"`
	AuditLogs []AuditLogInfo `xml:"audit-log"`
}

// AuditLogInfo 描述一条审计日志记录，包含序列号、时间戳、消息 ID、
// 严重级别、事件类别、操作用户、IP 地址、角色、接口等字段。
// 对应 YANG 模型: brocade-logging/audit-log
type AuditLogInfo struct {
	XMLName                xml.Name `xml:"audit-log" json:"-"`
	SequenceNumber         uint32   `xml:"sequence-number" json:"sequence_number"`
	TimeStamp              string   `xml:"time-stamp" json:"time_stamp"`
	MessageID              string   `xml:"message-id" json:"message_id"`
	SwitchUserFriendlyName string   `xml:"switch-user-friendly-name" json:"switch_user_friendly_name"`
	MessageText            string   `xml:"message-text" json:"message_text"`
	SeverityLevel          string   `xml:"severity-level" json:"severity_level"`
	EventClass             string   `xml:"event-class" json:"event_class"`
	IPAddress              string   `xml:"ip-address" json:"ip_address"`
	UserName               string   `xml:"user-name" json:"user_name"`
	Role                   string   `xml:"role" json:"role"`
	Interface              string   `xml:"interface" json:"interface"`
	ApplicationName        string   `xml:"application-name" json:"application_name"`
	FabricID               uint32   `xml:"fabric-id" json:"fabric_id"`
	EventInfo              string   `xml:"event-info" json:"event_info"`
	ApplicationUserName    string   `xml:"application-user-name" json:"application_user_name"`
}

// GetAuditLogs 获取交换机上的所有审计日志记录。
// 对应 API: GET /brocade-logging/audit-log
func (c *Client) GetAuditLogs() ([]AuditLogInfo, error) {
	if err := c.ensureLoggingSupported("audit-log"); err != nil {
		return nil, err
	}

	var resp AuditLogResponse
	err := c.Get(c.endpoints().AuditLogs(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.AuditLogs, nil
}

func (c *Client) ensureLoggingSupported(operation string) error {
	if c.endpoints().allowLogging() {
		return nil
	}
	return fmt.Errorf("%w: FOS %s does not support %s endpoint", ErrUnsupportedOperation, c.endpoints().version, operation)
}
