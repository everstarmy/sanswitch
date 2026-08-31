package sanswitch

import (
	"context"
	"encoding/xml"
)

// ZoneTransactionStatus 表示当前 Zone 事务的状态
type ZoneTransactionStatus struct {
	XMLName          xml.Name `xml:"effective-configuration" json:"-"`
	TransactionToken uint32   `xml:"transaction-token" json:"transaction_token"`
}

// zoneTransactionStatusResponse 是 GET transaction-token 接口的响应包装
type zoneTransactionStatusResponse struct {
	XMLName xml.Name              `xml:"Response"`
	Status  ZoneTransactionStatus `xml:"effective-configuration"`
}

// AbortZoneTransaction 中止当前挂起的 Zone 事务。
func (c *Client) AbortZoneTransaction(ctx context.Context) error {
	return c.patchWithoutVersionGate(ctx, c.endpoints().ZoneAbortTransaction(), nil)
}

// ZoneTransactionStatus 查询当前 Zone 事务状态。
func (c *Client) ZoneTransactionStatus(ctx context.Context) (*ZoneTransactionStatus, error) {
	var resp zoneTransactionStatusResponse
	err := c.get(ctx, c.endpoints().ZoneTransactionStatus(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Status, nil
}
