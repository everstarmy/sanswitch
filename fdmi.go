package san

import "encoding/xml"

// ==================== FDMI HBA ====================

// FDMIHBAResponse 是 GET /brocade-fdmi/hba 的 XML 响应包装
type FDMIHBAResponse struct {
	XMLName xml.Name      `xml:"Response"`
	HBAs    []FDMIHBAInfo `xml:"hba"`
}

// FDMIHBAInfo 描述 FDMI（Fabric Device Management Interface）中注册的一个 HBA 适配器，
// 包含 HBA ID、制造商、型号、序列号、固件版本、驱动版本、操作系统等字段。
// 对应 YANG 模型: brocade-fdmi/hba
type FDMIHBAInfo struct {
	XMLName            xml.Name `xml:"hba" json:"-"`
	HBAID              string   `xml:"hba-id" json:"hba_id"`
	DomainID           string   `xml:"domain-id" json:"domain_id"`
	Manufacturer       string   `xml:"manufacturer" json:"manufacturer"`
	SerialNumber       string   `xml:"serial-number" json:"serial_number"`
	Model              string   `xml:"model" json:"model"`
	ModelDescription   string   `xml:"model-description" json:"model_description"`
	NodeName           string   `xml:"node-name" json:"node_name"`
	NodeSymbolicName   string   `xml:"node-symbolic-name" json:"node_symbolic_name"`
	HardwareVersion    string   `xml:"hardware-version" json:"hardware_version"`
	DriverVersion      string   `xml:"driver-version" json:"driver_version"`
	OptionROMVersion   string   `xml:"option-rom-version" json:"option_rom_version"`
	FirmwareVersion    string   `xml:"firmware-version" json:"firmware_version"`
	OSNameAndVersion   string   `xml:"os-name-and-version" json:"os_name_and_version"`
	MaxCTPayload       uint32   `xml:"max-ct-payload" json:"max_ct_payload"`
	VendorID           string   `xml:"vendor-id" json:"vendor_id"`
	VendorSpecificInfo string   `xml:"vendor-specific-info" json:"vendor_specific_info"`
	NumberOfPorts      uint32   `xml:"number-of-ports" json:"number_of_ports"`
	FabricName         string   `xml:"fabric-name" json:"fabric_name"`
	BootBIOSVersion    string   `xml:"boot-bios-version" json:"boot_bios_version"`
	BootBIOSEnabledV2  bool     `xml:"boot-bios-enabled-v2" json:"boot_bios_enabled_v2"`
	HBAPortList        []string `xml:"hba-port-list>wwn" json:"hba_port_list"`
}

// GetFDMIHBAs 获取 FDMI 中注册的所有 HBA 适配器信息。
// 对应 API: GET /brocade-fdmi/hba
func (c *Client) GetFDMIHBAs() ([]FDMIHBAInfo, error) {
	var resp FDMIHBAResponse
	err := c.Get("/brocade-fdmi/hba", &resp)
	if err != nil {
		return nil, err
	}
	return resp.HBAs, nil
}

// ==================== FDMI Port ====================

// FDMIResponse 是 GET /brocade-fdmi/port 的 XML 响应包装
type FDMIResponse struct {
	XMLName xml.Name       `xml:"Response"`
	Ports   []FDMIportInfo `xml:"port"`
}

// FDMIportInfo 描述 FDMI 中注册的一个 FC 端口，包含端口 WWN、所属 HBA、
// 端口类型、支持的协议速率、当前速率、最大帧大小、主机名、VSA 等字段。
// 对应 YANG 模型: brocade-fdmi/port
type FDMIportInfo struct {
	XMLName                    xml.Name `xml:"port" json:"-"`
	PortName                   string   `xml:"port-name" json:"port_name"`
	HBAID                      string   `xml:"hba-id" json:"hba_id"`
	DomainID                   string   `xml:"domain-id" json:"domain_id"`
	PortSymbolicName           string   `xml:"port-symbolic-name" json:"port_symbolic_name"`
	PortID                     string   `xml:"port-id" json:"port_id"`
	PortType                   string   `xml:"port-type" json:"port_type"`
	SupportedClassOfService    string   `xml:"supported-class-of-service" json:"supported_class_of_service"`
	SupportedFC4Type           string   `xml:"supported-fc4-type" json:"supported_fc4_type"`
	ActiveFC4Type              string   `xml:"active-fc4-type" json:"active_fc4_type"`
	SupportedProtocolSpeeds    []string `xml:"supported-protocol-speeds>supported-protocol-speed" json:"supported_protocol_speeds"`
	CurrentProtocolSpeed       string   `xml:"current-protocol-speed" json:"current_protocol_speed"`
	MaximumFrameSize           uint32   `xml:"maximum-frame-size" json:"maximum_frame_size"`
	OSDeviceName               string   `xml:"os-device-name" json:"os_device_name"`
	HostName                   string   `xml:"host-name" json:"host_name"`
	NodeName                   string   `xml:"node-name" json:"node_name"`
	FabricName                 string   `xml:"fabric-name" json:"fabric_name"`
	PortState                  string   `xml:"port-state" json:"port_state"`
	NumberOfDiscoveredPorts    uint32   `xml:"number-of-discovered-ports" json:"number_of_discovered_ports"`
	VSAServiceCategory         string   `xml:"vsa-service-category" json:"vsa_service_category"`
	VSAGUID                    string   `xml:"vsa-guid" json:"vsa_guid"`
	VSAVersion                 string   `xml:"vsa-version" json:"vsa_version"`
	VSAProductName             string   `xml:"vsa-product-name" json:"vsa_product_name"`
	VSAPortInfo                string   `xml:"vsa-port-info" json:"vsa_port_info"`
	VSAQoSSupported            string   `xml:"vsa-qos-supported" json:"vsa_qos_supported"`
	VSASecurity                string   `xml:"vsa-security" json:"vsa_security"`
	VSAStorageArrayFamily      string   `xml:"vsa-storage-array-family" json:"vsa_storage_array_family"`
	VSAStorageArrayName        string   `xml:"vsa-storage-array-name" json:"vsa_storage_array_name"`
	VSAStorageArraySystemModel string   `xml:"vsa-storage-array-system-model" json:"vsa_storage_array_system_model"`
	VSAStorageArrayOS          string   `xml:"vsa-storage-array-os" json:"vsa_storage_array_os"`
	VSAStorageArrayNodeCount   uint32   `xml:"vsa-storage-array-node-count" json:"vsa_storage_array_node_count"`
	VSAStorageArrayNodes       []string `xml:"vsa-storage-array-nodes>nodes" json:"vsa_storage_array_nodes"`
	VSAConnectedPorts          []string `xml:"vsa-connected-ports>wwns" json:"vsa_connected_ports"`
	VSAEndToEndVersion         string   `xml:"vsa-end-to-end-version" json:"vsa_end_to_end_version"`
}

// GetFDMIports 获取 FDMI 中注册的所有 FC 端口信息。
// 对应 API: GET /brocade-fdmi/port
func (c *Client) GetFDMIports() ([]FDMIportInfo, error) {
	var resp FDMIResponse
	err := c.Get("/brocade-fdmi/port", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Ports, nil
}
