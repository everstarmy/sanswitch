package main

import (
	"fmt"
	"log"
	"time"

	san "github.com/everstarmy/sanswitch"
)

func main() {
	host := "192.168.1.100"
	username := "admin"
	password := "password"

	// 创建并自动登录（默认 HTTPS + 跳过证书验证）
	switchClient, err := san.NewSANSwitch(host, username, password,
		san.WithTimeout(60*time.Second),
	)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	switchClient.SetVerbose(true)

	fmt.Println("=== 登录成功 ===")

	defer func() {
		fmt.Println("\n=== 登出交换机 ===")
		if err := switchClient.Logout(); err != nil {
			log.Printf("登出失败: %v", err)
		} else {
			fmt.Println("登出成功")
		}
	}()

	fmt.Println("\n=== 获取Fabric中的所有交换机 ===")
	switches, err := switchClient.GetFabricSwitches()
	if err != nil {
		log.Fatalf("获取交换机列表失败: %v", err)
	}
	for _, sw := range switches {
		fmt.Printf("交换机: %s (WWNN: %s, Domain: %d, IP: %s)\n",
			sw.SwitchUserFriendlyName, sw.Name, sw.DomainID, sw.IPAddress)
	}

	fmt.Println("\n=== 获取交换机基本信息 ===")
	switchInfo, err := switchClient.GetSwitchInfo()
	if err != nil {
		log.Fatalf("获取交换机信息失败: %v", err)
	}
	fmt.Printf("交换机名称: %s\n", switchInfo.Name)
	fmt.Printf("WWN: %s\n", switchInfo.WWN)
	fmt.Printf("Chassis WWN: %s\n", switchInfo.ChassisWWN)
	fmt.Printf("域ID: %d\n", switchInfo.DomainID)
	fmt.Printf("固件版本: %s\n", switchInfo.FirmwareVersion)
	fmt.Printf("型号: %s\n", switchInfo.ModelName)
	fmt.Printf("序列号: %s\n", switchInfo.SerialNumber)
	fmt.Printf("FCID: %s (%s)\n", switchInfo.Fcid, switchInfo.FcidHex)
	fmt.Printf("IP地址: %s\n", switchInfo.IPAddress)

	fmt.Println("\n=== 获取端口信息 ===")
	ports, err := switchClient.GetPorts()
	if err != nil {
		log.Fatalf("获取端口信息失败: %v", err)
	}
	for _, port := range ports {
		fmt.Printf("端口 %s: 状态=%s, 速率=%s, WWN=%s, FCID=%s, 端口类型=%s\n",
			port.Name, port.OperationalStatusString, port.Speed, port.WWN, port.FCID, port.PortType)
	}

	fmt.Println("\n=== 获取硬件信息 ===")
	hardware, err := switchClient.GetHardwareInfo()
	if err != nil {
		log.Fatalf("获取硬件信息失败: %v", err)
	}
	fmt.Printf("机箱类型: %s\n", hardware.ChassisType)
	fmt.Printf("机箱序列号: %s\n", hardware.ChassisSerial)
	fmt.Printf("插槽数: %d\n", hardware.NumberOfSlots)
	fmt.Printf("端口数: %d\n", hardware.NumberOfPorts)
	fmt.Printf("CPU型号: %s\n", hardware.CPUModel)

	fmt.Println("\n=== 获取Switch Status Policy Report (MAPS) ===")
	switchStatus, err := switchClient.GetSwitchStatusPolicyReport()
	if err != nil {
		log.Fatalf("获取Switch Status Policy Report失败: %v", err)
	}
	fmt.Printf("交换机健康状态: %s\n", switchStatus.SwitchStatus)
	fmt.Printf("电源健康状态: %s\n", switchStatus.PowerSupplyHealth)
	fmt.Printf("风扇健康状态: %s\n", switchStatus.FanHealth)
	fmt.Printf("温度传感器健康状态: %s\n", switchStatus.TemperatureSensorHealth)
	fmt.Printf("HA健康状态: %s\n", switchStatus.HAHealth)
	fmt.Printf("Flash健康状态: %s\n", switchStatus.FlashHealth)
	fmt.Printf("边缘端口健康状态: %s\n", switchStatus.MarginalPortHealth)
	fmt.Printf("故障端口健康状态: %s\n", switchStatus.FaultyPortHealth)

	fmt.Println("\n=== 获取System Resources (MAPS) ===")
	sysResources, err := switchClient.GetSystemResources()
	if err != nil {
		log.Fatalf("获取System Resources失败: %v", err)
	}
	fmt.Printf("CPU使用率: %d%%\n", sysResources.CPUUsage)
	fmt.Printf("内存使用: %d\n", sysResources.MemoryUsage)
	fmt.Printf("总内存: %d\n", sysResources.TotalMemory)
	fmt.Printf("Flash使用: %d\n", sysResources.FlashUsage)
	fmt.Printf("可用内核内存: %d\n", sysResources.FreeKernelMemory)

	fmt.Println("\n=== 获取SFP信息 ===")
	sfps, err := switchClient.GetMediaRDPs()
	if err != nil {
		log.Fatalf("获取SFP信息失败: %v", err)
	}
	fmt.Printf("SFP数量: %d\n", len(sfps))
	for _, sfp := range sfps {
		fmt.Printf("端口 %s: 厂商=%s, 型号=%s, 序列号=%s, 温度=%f, Tx=%f, Rx=%f\n",
			sfp.Name, sfp.VendorName, sfp.PartNumber, sfp.SerialNumber,
			sfp.Temperature, sfp.TxPower, sfp.RxPower)
	}

	fmt.Println("\n=== 获取Blade信息 ===")
	blades, err := switchClient.GetBlades()
	if err != nil {
		log.Fatalf("获取Blade信息失败: %v", err)
	}
	fmt.Printf("Blade数量: %d\n", len(blades))
	for _, blade := range blades {
		fmt.Printf("Blade插槽%d (类型: %s): 状态=%s, 序列号=%s, 固件=%s\n",
			blade.SlotNumber, blade.BladeTypeString, blade.BladeStateString, blade.SerialNumber, blade.FirmwareVersion)
	}

	fmt.Println("\n=== 获取风扇信息 ===")
	fans, err := switchClient.GetFans()
	if err != nil {
		log.Fatalf("获取风扇信息失败: %v", err)
	}
	fmt.Printf("风扇数量: %d\n", len(fans))
	for _, fan := range fans {
		fmt.Printf("风扇单元%d: 状态=%s, 转速=%d RPM\n",
			fan.UnitNumber, fan.OperationalState, fan.SpeedRPM)
	}

	fmt.Println("\n=== 获取电源信息 ===")
	psus, err := switchClient.GetPowerSupplies()
	if err != nil {
		log.Fatalf("获取电源信息失败: %v", err)
	}
	fmt.Printf("电源数量: %d\n", len(psus))
	for _, psu := range psus {
		fmt.Printf("电源单元%d: 状态=%s, 输入功率=%dW\n",
			psu.UnitNumber, psu.OperationalState, psu.InputPower)
	}

	fmt.Println("\n=== 获取传感器信息 ===")
	sensors, err := switchClient.GetSensors()
	if err != nil {
		log.Fatalf("获取传感器信息失败: %v", err)
	}
	fmt.Printf("传感器数量: %d\n", len(sensors))
	for _, sensor := range sensors {
		fmt.Printf("传感器 %s (类型: %s): 值=%d, 状态=%s\n",
			sensor.Index, sensor.Category, sensor.Temperature, sensor.State)
	}

	fmt.Println("\n=== 获取FRU历史日志 ===")
	logs, err := switchClient.GetHistoryLogs()
	if err != nil {
		log.Fatalf("获取FRU历史日志失败: %v", err)
	}
	fmt.Printf("历史日志数量: %d\n", len(logs))
	for _, log := range logs {
		fmt.Printf("[%s] %s slot=%d state=%s sn=%s\n",
			log.TimeStamp, log.FRUType, log.Position, log.State, log.SerialNumber)
	}

	fmt.Println("\n=== 获取Defined Zones (定义的Zone) ===")
	definedZones, err := switchClient.GetDefinedZones()
	if err != nil {
		log.Fatalf("获取Defined Zones失败: %v", err)
	}
	fmt.Printf("定义的Zone数量: %d\n", len(definedZones))
	for _, zone := range definedZones {
		fmt.Printf("Zone %s (类型: %s): 成员=%v\n", zone.Name, zone.Type, zone.Members)
	}

	fmt.Println("\n=== 获取Effective Zones (生效的Zone) ===")
	effectiveZones, err := switchClient.GetEffectiveZones()
	if err != nil {
		log.Fatalf("获取Effective Zones失败: %v", err)
	}
	fmt.Printf("生效的Zone数量: %d\n", len(effectiveZones))
	for _, zone := range effectiveZones {
		fmt.Printf("Zone %s (类型: %s): 成员=%v\n", zone.Name, zone.Type, zone.Members)
	}

	fmt.Println("\n=== 获取Defined Aliases (定义的Alias) ===")
	aliases, err := switchClient.GetDefinedAliases()
	if err != nil {
		log.Fatalf("获取Alias信息失败: %v", err)
	}
	for _, alias := range aliases {
		fmt.Printf("Alias %s: 成员=%v\n", alias.Name, alias.Members)
	}

	fmt.Println("\n=== 获取Defined Configs (定义的配置) ===")
	definedConfigs, err := switchClient.GetDefinedConfigs()
	if err != nil {
		log.Fatalf("获取Defined Configs失败: %v", err)
	}
	for _, cfg := range definedConfigs {
		fmt.Printf("配置 %s: 包含Zone=%v\n", cfg.Name, cfg.MemberZones)
	}

	fmt.Println("\n=== 获取Effective Config (生效的配置) ===")
	effectiveConfig, err := switchClient.GetEffectiveConfig()
	if err != nil {
		log.Fatalf("获取Effective Config失败: %v", err)
	}
	fmt.Printf("活动配置名称: %s\n", effectiveConfig.Name)
	fmt.Printf("校验和: %s\n", effectiveConfig.Checksum)
	fmt.Printf("默认Zone访问: %s\n", effectiveConfig.DefaultZoneAccess)

	fmt.Println("\n=== 获取Zone数据库信息 ===")
	dbInfo, err := switchClient.GetZoneDatabaseInfo()
	if err != nil {
		log.Fatalf("获取Zone数据库信息失败: %v", err)
	}
	fmt.Printf("数据库最大容量: %d\n", dbInfo.DBMax)
	fmt.Printf("数据库可用空间: %d\n", dbInfo.DBAvail)
	fmt.Printf("已提交的数据库大小: %d\n", dbInfo.DBCommitted)
	fmt.Printf("事务Token: %d\n", dbInfo.TransactionToken)

	fmt.Println("\n=== 获取逻辑交换机信息 ===")
	logicalSwitches, err := switchClient.GetLogicalSwitches()
	if err != nil {
		log.Fatalf("获取逻辑交换机信息失败: %v", err)
	}
	for _, ls := range logicalSwitches {
		fmt.Printf("逻辑交换机: FabricID=%d, WWN=%s, BaseSwitch=%v, DefaultSwitch=%v\n",
			ls.FabricID, ls.SwitchWWN, ls.BaseSwitchEnabled, ls.DefaultSwitch)
		fmt.Printf("  端口成员: %v\n", ls.PortMembers)
		fmt.Printf("  GE端口成员: %v\n", ls.GePortMembers)
	}

	fmt.Println("\n=== 设置Virtual Fabric ID (示例: VF 128) ===")
	switchClient.SetVFID(128)

	fmt.Println("\n=== 获取VF 128的逻辑交换机信息 ===")
	logicalSwitches, err = switchClient.GetLogicalSwitches()
	if err != nil {
		log.Fatalf("获取VF 128逻辑交换机信息失败: %v", err)
	}
	for _, ls := range logicalSwitches {
		fmt.Printf("VF %d 逻辑交换机: FabricID=%d, WWN=%s\n",
			128, ls.FabricID, ls.SwitchWWN)
	}
}
