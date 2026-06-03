# sanswitch

Brocade FOS REST API Go 客户端库，用于采集和管理 Brocade SAN 交换机（Fibre Channel Fabric）的硬件信息、端口状态、Zone 配置等。

## 特性

- 覆盖 Brocade FOS 9.x REST API 的 **35+ 个采集端点**，涵盖 Fabric、端口、FRU、Zone、SFP、MAPS、Trunk、SNMP、NTP 等模块
- `NewSANSwitch` 创建即自动登录，失败返回 `nil, err`
- 默认 HTTPS + 跳过证书验证，可通过 `WithHTTP()` 切换为 HTTP
- `SwitchAPI` 接口抽象，支持 Mock 测试与依赖注入
- `context.Context` 请求级取消与超时控制
- `log/slog` 结构化日志
- 自动重试 + 指数退避（网络错误、429、5xx）
- FOS 结构化错误解析（`APIError`）
- Virtual Fabric（VFID）支持
- 函数选项模式（`ClientOption`）灵活配置

## 安装

```bash
go get github.com/everstarmy/sanswitch
```

## 快速开始

```go
package main

import (
	"fmt"
	"log"

	"github.com/everstarmy/sanswitch"
)

func main() {
	// 创建并自动登录（默认 HTTPS + 跳过证书验证）
	sw, err := san.NewSANSwitch("192.168.1.100", "admin", "password")
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	defer sw.Logout()

	// 获取 Fabric 中所有交换机
	switches, _ := sw.GetFabricSwitches()
	for _, s := range switches {
		fmt.Printf("%s (Domain %d, IP %s)\n", s.Name, s.DomainID, s.IPAddress)
	}

	// 获取端口列表
	ports, _ := sw.GetPorts()
	for _, p := range ports {
		fmt.Printf("端口 %s: %s, 速率 %s\n", p.Name, p.OperationalStatusString, p.Speed)
	}
}
```

更多示例请参考 [example/main.go](example/main.go)。

## 客户端配置

通过 `ClientOption` 函数选项自定义客户端行为：

```go
// HTTPS + 自动登录（默认）
sw, err := san.NewSANSwitch("192.168.1.100", "admin", "password",
	san.WithTimeout(60*time.Second),       // 请求超时（默认 30s）
	san.WithRetry(5),                      // 重试次数（默认 3）
	san.WithRetryWait(2*time.Second),      // 重试初始等待（默认 1s）
	san.WithRetryMaxWait(60*time.Second),  // 重试最大等待（默认 30s）
	san.WithLogger(myLogger),              // 注入自定义 slog.Logger
)

// HTTP 模式（仅用于开发/测试环境）
sw, err := san.NewSANSwitch("192.168.1.100", "admin", "password",
	san.WithHTTP(),
)
```

| 选项 | 默认值 | 说明 |
|------|--------|------|
| `WithTimeout(d)` | `30s` | 请求超时时间 |
| `WithRetry(n)` | `3` | 最大重试次数 |
| `WithRetryWait(d)` | `1s` | 重试初始等待（指数退避起点） |
| `WithRetryMaxWait(d)` | `30s` | 重试最大等待上限 |
| `WithLogger(l)` | `slog.Default()` | 自定义结构化日志 |
| `WithHTTP()` | HTTPS | 使用 HTTP 替代 HTTPS |
```

### 调试日志

`SetVerbose(true)` 会启用 `slog.LevelDebug`，输出 HTTP 请求/响应调试日志：

```go
sw.SetVerbose(true)
```

如需指定日志输出位置，请先设置输出，再开启 verbose：

```go
sw.SetLogOutput(os.Stdout)
sw.SetVerbose(true)
```

也可以通过 `WithLogger` 注入完整的自定义 `slog.Logger`：

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))
sw, err := san.NewSANSwitch("192.168.1.100", "admin", "password", san.WithLogger(logger))
```

### Virtual Fabric 支持

```go
sw.SetVFID(128) // 后续请求将自动附加 ?vf-id=128 参数
sw.SetVFID(0)   // 设回 0 可取消 VFID
```

## API 接口总览

### 1. 基本信息

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetSwitchInfo()` | — | 获取本交换机摘要信息 |
| `GetFabricSwitches()` | `/brocade-fabric/fabric-switch` | 获取 Fabric 中所有交换机 |
| `GetHardwareInfo()` | — | 获取硬件信息（机箱、CPU、端口数等） |

### 2. 端口与统计

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetPorts()` | `/brocade-interface/fibrechannel` | 获取所有 FC 端口信息 |
| `GetFibreChannelStatistics()` | `/brocade-interface/fibrechannel-statistics` | 获取端口性能统计计数器 |

### 3. 逻辑交换机

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetLogicalSwitches()` | `/brocade-fibrechannel-logical-switch/fibrechannel-logical-switch` | 获取逻辑交换机及端口成员 |

### 4. Zone 管理

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetDefinedZones()` | `/brocade-zone/defined-configuration/zone` | 获取已定义的 Zone |
| `GetEffectiveZones()` | `/brocade-zone/effective-configuration/enabled-zone` | 获取生效的 Zone |
| `GetDefinedAliases()` | `/brocade-zone/defined-configuration/alias` | 获取已定义的 Alias |
| `GetDefinedConfigs()` | `/brocade-zone/defined-configuration/cfg` | 获取已定义的 Zone 配置 |
| `GetEffectiveConfig()` | `/brocade-zone/effective-configuration` | 获取生效的 Zone 配置 |
| `GetZoneDatabaseInfo()` | `/brocade-zone/effective-configuration` | 获取 Zone 数据库容量信息 |
| `GetZoneChecksum()` | `/brocade-zone/effective-configuration/checksum` | 获取 Zone DB checksum |
| `CreateZone(name, members, principalMembers)` | `POST /brocade-zone/defined-configuration/zone` | 创建 defined Zone，不自动生效 |
| `UpdateZone(name, members, principalMembers)` | `PATCH /brocade-zone/defined-configuration/zone` | 全量替换 defined Zone 成员，不自动生效 |
| `RenameZone(oldName, newName)` | `PATCH /brocade-zone/defined-configuration/zone/zone-name/{oldName}` | 重命名 defined Zone，不自动生效 |
| `DeleteZone(name)` | `DELETE /brocade-zone/defined-configuration/zone/zone-name/{name}` | 删除 defined Zone，不自动生效 |
| `CreateAlias(name, members)` | `POST /brocade-zone/defined-configuration/alias` | 创建 Alias |
| `UpdateAlias(name, members)` | `PATCH /brocade-zone/defined-configuration/alias` | 更新 Alias 成员 |
| `RenameAlias(oldName, newName)` | `PATCH /brocade-zone/defined-configuration/alias/alias-name/{oldName}` | 重命名 Alias，不自动生效 |
| `DeleteAlias(name)` | `DELETE /brocade-zone/defined-configuration/alias/alias-name/{name}` | 删除 Alias |
| `UpdateDefinedConfig(name, memberZones)` | `PATCH /brocade-zone/defined-configuration/cfg` | 全量替换 cfg 成员 Zone |
| `SaveZoneConfig(checksum)` | `PATCH /brocade-zone/effective-configuration/cfg-action-v2/save` | 保存 Zone 配置 |
| `ActivateZoneConfig(name, checksum)` | `PATCH /brocade-zone/effective-configuration/cfg-name/{name}` | 激活 Zone 配置 |
| `CreateZoneAndActivate(cfg, zone, members, principalMembers)` | 组合流程 | 创建 Zone、加入 cfg、保存并激活 |
| `ReplaceZoneAndActivate(cfg, zone, members, principalMembers)` | 组合流程 | 全量替换 Zone 成员、保存并激活 |
| `DeleteZoneAndActivate(cfg, zone)` | 组合流程 | 删除 Zone、保存并激活 |
| `AbortZoneTransaction()` | `PATCH /brocade-zone/effective-configuration/cfg-action-v2/transaction-abort` | 中止 Zone 事务 |
| `GetZoneTransactionStatus()` | `/brocade-zone/effective-configuration/transaction-token` | 查询 Zone 事务状态 |

低层 Zone 写方法只修改 defined configuration。若需要让变更在 Fabric 中生效，推荐使用组合流程：

```go
err := sw.CreateZoneAndActivate(
	"cfg1",
	"zone_app01_storage01",
	[]string{"10:00:00:00:00:00:00:01", "20:00:00:00:00:00:00:01"},
	nil,
)
```

注意：`UpdateZone` 和 `ReplaceZoneAndActivate` 会按 Brocade REST API 语义全量覆盖 `member-entry` leaf-list，请传入所有需要保留的成员。

### 5. FRU 组件

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetBlades()` | `/brocade-fru/blade` | 获取刀片（CP/SW/Core）信息 |
| `GetFans()` | `/brocade-fru/fan` | 获取风扇信息 |
| `GetPowerSupplies()` | `/brocade-fru/power-supply` | 获取电源信息 |
| `GetSensors()` | `/brocade-fru/sensor` | 获取传感器信息 |
| `GetHistoryLogs()` | `/brocade-fru/history-log` | 获取 FRU 历史日志 |

### 6. MAPS 监控

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetSwitchStatusPolicyReport()` | `/brocade-maps/switch-status-policy-report` | 获取交换机健康状态策略报告 |
| `GetSystemResources()` | `/brocade-maps/system-resources` | 获取 CPU / 内存 / Flash 使用率 |

### 7. SFP / Media

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetMediaRDPs()` | `/brocade-media/media-rdp` | 获取 SFP 光模块详细信息（温度、光功率、厂商等） |

### 8. 名称服务器

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetFibreChannelNameServers()` | `/brocade-name-server/fibrechannel-name-server` | 获取 FC 名称服务器注册信息 |

### 9. FDMI

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetFDMIHBAs()` | `/brocade-fdmi/hba` | 获取 HBA 卡信息 |
| `GetFDMIports()` | `/brocade-fdmi/port` | 获取 FDMI 端口信息 |

### 10. Trunk（ISL 链路聚合）

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetTrunks()` | `/brocade-fibrechannel-trunk/trunk` | 获取 ISL Trunk 信息 |
| `GetTrunkPerformances()` | `/brocade-fibrechannel-trunk/performance` | 获取 Trunk 性能统计 |
| `GetTrunkAreas()` | `/brocade-fibrechannel-trunk/trunk-area` | 获取 Trunk Area 信息 |

### 11. Firmware

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetFirmwareHistory()` | `/brocade-firmware/firmware-history` | 获取固件升级历史记录 |

### 12. SNMP

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetSNMPSystem()` | `/brocade-snmp/system` | 获取 SNMP 系统配置 |
| `GetSNMPv1Accounts()` | `/brocade-snmp/v1-account` | 获取 SNMPv1 社区账户 |
| `GetSNMPv1Traps()` | `/brocade-snmp/v1-trap` | 获取 SNMPv1 Trap 目标 |
| `GetSNMPv3Accounts()` | `/brocade-snmp/v3-account` | 获取 SNMPv3 用户账户 |
| `GetSNMPv3Traps()` | `/brocade-snmp/v3-trap` | 获取 SNMPv3 Trap 目标 |

### 13. Time / NTP

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetTimeZone()` | `/brocade-time/time-zone` | 获取时区配置 |
| `GetClockServer()` | `/brocade-time/clock-server` | 获取 NTP 时钟服务器 |

### 14. 日志

| 方法 | REST 端点 | 说明 |
|------|----------|------|
| `GetErrorLogs()` | `/brocade-logging/error-log` | 获取错误日志（RAS Log） |
| `GetAuditLogs()` | `/brocade-logging/audit-log` | 获取审计日志 |

## 接口统计

| 类别 | 采集 (GET) | 写入 (POST/PATCH/DELETE) | 合计 |
|------|-----------|------------------------|------|
| 基本信息 | 3 | — | 3 |
| 端口与统计 | 2 | — | 2 |
| 逻辑交换机 | 1 | — | 1 |
| Zone 管理 | 7 | 15 | 22 |
| FRU 组件 | 5 | — | 5 |
| MAPS 监控 | 2 | — | 2 |
| SFP / Media | 1 | — | 1 |
| 名称服务器 | 1 | — | 1 |
| FDMI | 2 | — | 2 |
| Trunk | 3 | — | 3 |
| Firmware | 1 | — | 1 |
| SNMP | 5 | — | 5 |
| Time / NTP | 2 | — | 2 |
| 日志 | 2 | — | 2 |
| **合计** | **37** | **15** | **52** |

## Context 支持

所有 HTTP 方法均提供 `WithContext` 变体，支持请求级取消与超时控制：

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

var resp LoginResponse
err := client.GetWithContext(ctx, "/login", &resp)
```

可用的 Context 方法：`GetWithContext`、`PostWithContext`、`PatchWithContext`、`DeleteWithContext`。

## 错误处理

```go
import "errors"

switches, err := sw.GetFabricSwitches()
if errors.Is(err, san.ErrUnauthorized) {
	// 认证失败，需要重新登录
}

var apiErr *san.APIError
if errors.As(err, &apiErr) {
	fmt.Printf("FOS 错误码: %s, 消息: %s\n", apiErr.ErrorCode, apiErr.Message)
}
```

预定义错误：`ErrNotFound`、`ErrUnauthorized`、`ErrConnectionFailed`、`ErrInvalidResponse`、`ErrTimeout`。

## 接口抽象

`SwitchAPI` 接口定义了所有操作，支持 Mock 测试：

```go
// 编译期断言
var _ san.SwitchAPI = (*san.SANSwitch)(nil)

// Mock 示例
type MockSwitch struct { san.SwitchAPI }
func (m *MockSwitch) GetPorts() ([]san.PortInfo, error) {
	return []san.PortInfo{{Name: "0/0", Speed: "32000000000"}}, nil
}
```

## 技术栈

| 组件 | 选型 |
|------|------|
| 语言 | Go 1.21+ |
| HTTP 客户端 | [go-resty/resty/v2](https://github.com/go-resty/resty) |
| 数据格式 | XML（Brocade YANG Data Model `application/yang-data+xml`） |
| 日志 | `log/slog`（Go 标准库） |
| 测试 | `net/http/httptest`（标准库） |

## 项目结构

```
├── client.go           # HTTP 客户端、认证、重试、Context 支持
├── san.go              # SwitchAPI 接口 + SANSwitch facade
├── types.go            # 公共类型定义
├── errors.go           # 错误类型与 APIError
├── switch.go           # Fabric Switch 采集
├── port.go             # FC 端口采集
├── statistics.go       # 端口性能统计
├── hardware.go         # 硬件信息
├── logical_switch.go   # 逻辑交换机
├── zone.go             # Zone 查询与写操作
├── alias.go            # Alias 管理
├── config.go           # Zone 配置管理
├── fru.go              # FRU 组件（blade/fan/psu/sensor/history-log）
├── maps.go             # MAPS 监控
├── media.go            # SFP 光模块
├── nameserver.go       # FC 名称服务器
├── fdmi.go             # FDMI HBA/Port
├── trunk.go            # ISL Trunk 链路聚合
├── firmware.go         # 固件升级历史
├── snmp.go             # SNMP 配置
├── time.go             # NTP 时钟与时间区
├── logging.go          # 日志（error-log / audit-log）
├── *_test.go           # 单元测试
├── example/
│   └── main.go         # 使用示例
└── docs/
    └── brocade-rest-api-doc/  # Brocade FOS REST API 参考文档
```

## 参考文档

- [Fabric OS REST API Overview](https://techdocs.broadcom.com/us/en/fibre-channel-networking/fabric-os/fabric-os-rest-api/9-2-x.html)
- Brocade FOS 9.2.x REST API Reference (fos-92x-restapi.pdf)
