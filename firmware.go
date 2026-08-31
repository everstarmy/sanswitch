package sanswitch

import (
	"context"
	"encoding/xml"
	"fmt"
)

// ==================== Firmware History ====================

// firmwareHistoryResponse 对应 GET /brocade-firmware/firmware-history
type firmwareHistoryResponse struct {
	XMLName         xml.Name          `xml:"Response"`
	FirmwareHistory []FirmwareHistory `xml:"firmware-history"`
}

// FirmwareHistory 描述一条固件安装历史记录
type FirmwareHistory struct {
	XMLName         xml.Name `xml:"firmware-history" json:"-"`
	SequenceNumber  uint16   `xml:"sequence-number" json:"sequence_number"`
	TimeStamp       string   `xml:"time-stamp" json:"time_stamp"`
	SwitchName      string   `xml:"switch-name" json:"switch_name"`
	SlotNumber      uint16   `xml:"slot-number" json:"slot_number"`
	ProcessID       uint32   `xml:"process-id" json:"process_id"`
	FirmwareVersion string   `xml:"firmware-version" json:"firmware_version"`
}

// ==================== Client Methods ====================

// FirmwareHistory 获取固件版本安装历史并允许取消请求。
func (c *Client) FirmwareHistory(ctx context.Context) ([]FirmwareHistory, error) {
	if err := c.ensureFirmwareHistorySupported(); err != nil {
		return nil, err
	}
	var resp firmwareHistoryResponse
	err := c.get(ctx, c.endpoints().FirmwareHistory(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.FirmwareHistory, nil
}

func (c *Client) ensureFirmwareHistorySupported() error {
	if c.endpoints().allowFirmwareHistory() {
		return nil
	}
	return fmt.Errorf("%w: FOS %s does not support firmware-history endpoint", ErrUnsupportedOperation, c.endpoints().version)
}
