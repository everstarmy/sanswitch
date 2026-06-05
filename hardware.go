package san

import "encoding/xml"

// ChassisResponse 是 GET /brocade-chassis/chassis 的 XML 响应包装
type ChassisResponse struct {
	XMLName xml.Name       `xml:"Response"`
	Chassis []HardwareInfo `xml:"chassis"`
}

// HardwareInfo 描述交换机机箱硬件信息，包含机箱类型、序列号、CPU、内存、Flash、
// 电源数量、风扇数量和温度等字段。
// 对应 YANG 模型: brocade-chassis/chassis
type HardwareInfo struct {
	XMLName          xml.Name `xml:"chassis" json:"-"`
	ChassisType      string   `xml:"chassis-type" json:"chassis_type"`             // 机箱类型（如 "SAN Director"）
	ChassisSerial    string   `xml:"serial-number" json:"chassis_serial"`          // 机箱序列号
	NumberOfSlots    int      `xml:"number-of-slots" json:"number_of_slots"`       // 插槽总数
	NumberOfPorts    int      `xml:"number-of-ports" json:"number_of_ports"`       // 端口总数
	CPUModel         string   `xml:"cpu-model" json:"cpu_model"`                   // CPU 型号
	MemorySize       string   `xml:"memory-size" json:"memory_size"`               // 内存大小
	FlashSize        string   `xml:"flash-size" json:"flash_size"`                 // Flash 存储大小
	PowerSupplyCount int      `xml:"power-supply-count" json:"power_supply_count"` // 电源数量
	FanCount         int      `xml:"fan-count" json:"fan_count"`                   // 风扇数量
	Temperature      string   `xml:"temperature" json:"temperature"`               // 温度（字符串格式）
}

// GetHardwareInfo 获取交换机机箱硬件信息（取列表中的第一条记录）。
// 若响应为空则返回 ErrNotFound。
// 对应 API: GET /brocade-chassis/chassis
func (c *Client) GetHardwareInfo() (*HardwareInfo, error) {
	var resp ChassisResponse
	err := c.Get(c.endpoints().Chassis(), &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Chassis) == 0 {
		return nil, ErrNotFound
	}

	return &resp.Chassis[0], nil
}
