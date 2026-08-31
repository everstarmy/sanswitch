package sanswitch

import (
	"context"
	"encoding/xml"
	"fmt"
)

// ==================== Error Log (RAS Log) ====================

// errorLogResponse 是 GET /brocade-logging/error-log 的 XML 响应包装
type errorLogResponse struct {
	XMLName   xml.Name   `xml:"Response"`
	ErrorLogs []ErrorLog `xml:"error-log"`
}

// ErrorLog 描述一条 RAS 错误日志记录，包含序列号、时间戳、消息 ID、
// 严重级别、Fabric ID、插槽 ID、消息文本等字段。
// 对应 YANG 模型: brocade-logging/error-log
type ErrorLog struct {
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

// ErrorLogs 获取 RAS 错误日志并允许取消请求。
func (c *Client) ErrorLogs(ctx context.Context) ([]ErrorLog, error) {
	if err := c.ensureLoggingSupported("error-log"); err != nil {
		return nil, err
	}

	var resp errorLogResponse
	err := c.get(ctx, c.endpoints().ErrorLogs(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.ErrorLogs, nil
}

// ==================== Audit Log ====================

// auditLogResponse 是 GET /brocade-logging/audit-log 的 XML 响应包装
type auditLogResponse struct {
	XMLName   xml.Name   `xml:"Response"`
	AuditLogs []AuditLog `xml:"audit-log"`
}

// AuditLog 描述一条审计日志记录，包含序列号、时间戳、消息 ID、
// 严重级别、事件类别、操作用户、IP 地址、角色、接口等字段。
// 对应 YANG 模型: brocade-logging/audit-log
type AuditLog struct {
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

// AuditLogs 获取审计日志并允许取消请求。
func (c *Client) AuditLogs(ctx context.Context) ([]AuditLog, error) {
	if err := c.ensureLoggingSupported("audit-log"); err != nil {
		return nil, err
	}

	var resp auditLogResponse
	err := c.get(ctx, c.endpoints().AuditLogs(), &resp)
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
