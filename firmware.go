package san

import "encoding/xml"

// ==================== Firmware History ====================

// FirmwareHistoryResponse 对应 GET /brocade-firmware/firmware-history
type FirmwareHistoryResponse struct {
	XMLName         xml.Name              `xml:"Response"`
	FirmwareHistory []FirmwareHistoryInfo `xml:"firmware-history"`
}

// FirmwareHistoryInfo 描述一条固件安装历史记录
type FirmwareHistoryInfo struct {
	XMLName         xml.Name `xml:"firmware-history" json:"-"`
	SequenceNumber  uint16   `xml:"sequence-number" json:"sequence_number"`
	TimeStamp       string   `xml:"time-stamp" json:"time_stamp"`
	SwitchName      string   `xml:"switch-name" json:"switch_name"`
	SlotNumber      uint16   `xml:"slot-number" json:"slot_number"`
	ProcessID       uint32   `xml:"process-id" json:"process_id"`
	FirmwareVersion string   `xml:"firmware-version" json:"firmware_version"`
}

// ==================== Client Methods ====================

// GetFirmwareHistory 获取固件版本安装历史
func (c *Client) GetFirmwareHistory() ([]FirmwareHistoryInfo, error) {
	var resp FirmwareHistoryResponse
	err := c.Get("/brocade-firmware/firmware-history", &resp)
	if err != nil {
		return nil, err
	}
	return resp.FirmwareHistory, nil
}
