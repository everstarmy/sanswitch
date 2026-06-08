package san

import "encoding/xml"

// ZoneTransactionStatus 表示当前 Zone 事务的状态
type ZoneTransactionStatus struct {
	XMLName          xml.Name `xml:"effective-configuration" json:"-"`
	TransactionToken uint32   `xml:"transaction-token" json:"transaction_token"`
}

// ZoneTransactionStatusResponse 是 GET transaction-token 接口的响应包装
type ZoneTransactionStatusResponse struct {
	XMLName xml.Name              `xml:"Response"`
	Status  ZoneTransactionStatus `xml:"effective-configuration"`
}

// AbortZoneTransaction 中止当前挂起的 Zone 事务（释放事务锁）
// FOS 9.1 及以下使用 PATCH /brocade-zone/effective-configuration/cfg-action/4，
// FOS 9.2+ 使用 PATCH /brocade-zone/effective-configuration/cfg-action-v2/transaction-abort。
func (c *Client) AbortZoneTransaction() error {
	return c.patchWithoutVersionGate(c.endpoints().ZoneAbortTransaction(), nil)
}

// GetZoneTransactionStatus 查询当前 Zone 事务状态
// transaction-token 为 0 表示无挂起事务
// 对应 API: GET /brocade-zone/effective-configuration/transaction-token
func (c *Client) GetZoneTransactionStatus() (*ZoneTransactionStatus, error) {
	var resp ZoneTransactionStatusResponse
	err := c.Get(c.endpoints().ZoneTransactionStatus(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Status, nil
}
