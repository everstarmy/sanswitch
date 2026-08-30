# sanswitch

Brocade FOS REST API Go 客户端库，用于采集和管理 Brocade SAN 交换机（Fibre Channel Fabric）的硬件信息、端口状态、Zone 配置等。

## 特性

- 覆盖 Brocade FOS 9.x REST API 的 **35+ 个采集端点**，涵盖 Fabric、端口、FRU、Zone、SFP、MAPS、Trunk、SNMP、NTP 等模块
- `NewSANSwitch` 创建即自动登录，失败返回 `nil, err`
- 默认 HTTPS 并校验证书；自签名证书可通过 `WithTLSConfig()` 配置，测试环境可显式使用 `WithInsecureSkipVerify()`
- 按会话、交换机、Zone、清单和监控拆分的细粒度能力接口，便于依赖注入
- `context.Context` 请求级取消与超时控制
- `log/slog` 结构化日志
- 自动重试 + 指数退避（网络错误、429、5xx）
- 默认限制响应体为 16 MiB，避免异常响应造成无界内存增长
- 仅依赖 Go 标准库，无第三方运行时依赖
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
	"os"

	"github.com/everstarmy/sanswitch"
)

func main() {
	host := os.Getenv("SAN_HOST")
	username := os.Getenv("SAN_USERNAME")
	password := os.Getenv("SAN_PASSWORD")
	if host == "" || username == "" || password == "" {
		log.Fatal("SAN_HOST、SAN_USERNAME、SAN_PASSWORD 必须设置")
	}

	// 创建并自动登录（默认 HTTPS 并校验证书）
	sw, err := san.NewSANSwitch(host, username, password)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	defer sw.Close()
	defer sw.Logout()

	// 获取 Fabric 中所有交换机
	switches, err := sw.GetFabricSwitches()
	if err != nil {
		log.Fatalf("获取交换机列表失败: %v", err)
	}
	for _, s := range switches {
		fmt.Printf("%s (Domain %d, IP %s)\n", s.Name, s.DomainID, s.IPAddress)
	}

	// 获取端口列表
	ports, err := sw.GetPorts()
	if err != nil {
		log.Fatalf("获取端口列表失败: %v", err)
	}
	for _, p := range ports {
		fmt.Printf("端口 %s: %s, 速率 %s\n", p.Name, p.OperationalStatusString, p.Speed)
	}
}
```

更多示例请参考 [example/main.go](example/main.go)。

运行示例前设置凭据，避免将密码写入源码或命令历史：

```bash
export SAN_HOST=192.168.1.100
export SAN_USERNAME=admin
export SAN_PASSWORD='change-me'
go run ./example
```

## 客户端配置

通过 `ClientOption` 函数选项自定义客户端行为：

```go
// HTTPS + 自动登录（默认）
sw, err := san.NewSANSwitch(host, username, password,
	san.WithTimeout(60*time.Second),       // 请求超时（默认 30s）
	san.WithRetry(5),                      // 重试次数（默认 3）
	san.WithRetryWait(2*time.Second),      // 重试初始等待（默认 1s）
	san.WithRetryMaxWait(60*time.Second),  // 重试最大等待（默认 30s）
	san.WithFOSVersion("v9.1.1"),          // 可选：直接使用 Client 时指定 FOS 版本
	san.WithLogger(myLogger),              // 注入自定义 slog.Logger
	san.WithTLSConfig(tlsConfig),          // 自定义 CA 或证书校验策略
)

// HTTP 模式（仅用于开发/测试环境）
sw, err := san.NewSANSwitch(host, username, password,
	san.WithHTTP(),
)
```

| 选项 | 默认值 | 说明 |
|------|--------|------|
| `WithTimeout(d)` | `30s` | 请求超时时间 |
| `WithRetry(n)` | `3` | 最大重试次数 |
| `WithRetryWait(d)` | `1s` | 重试初始等待（指数退避起点） |
| `WithRetryMaxWait(d)` | `30s` | 重试最大等待上限 |
| `WithMaxResponseBodyBytes(n)` | `16 MiB` | 单个响应体的最大字节数，超限返回 `ErrResponseBodyTooLarge` |
| `WithFOSVersion(v)` | 自动登录识别 / 未知时只读 | 指定 FOS 版本以兼容版本差异 endpoint |
| `WithLogger(l)` | `slog.Default()` | 自定义结构化日志 |
| `WithTLSConfig(c)` | 系统证书池 | 自定义 TLS/CA 配置 |
| `WithInsecureSkipVerify()` | 关闭 | 显式关闭 HTTPS 证书校验，仅建议测试使用 |
| `WithAllowUnknownFOSVersionWrites()` | 关闭 | 显式允许未知 FOS 版本执行写操作 |
| `WithHTTP()` | HTTPS | 使用 HTTP 替代 HTTPS |

版本识别规则：登录响应包含 `firmware-version` 时会自动记录版本，并按 `vX.Y` 选择兼容 endpoint；若登录无响应体，则按低于 9.1 的旧版 FOS 处理。未知版本默认禁止写操作，可通过 `WithFOSVersion()` 或显式的 `WithAllowUnknownFOSVersionWrites()` 配置。低于 9.1 的版本不允许执行写操作。`GetHistoryLogs`、`GetSensors`、`GetErrorLogs` 和 `GetAuditLogs` 需要 FOS 9.0+，`GetFirmwareHistory` 需要 FOS 9.1+。

### 调试日志

`SetVerbose(true)` 会启用 `slog.LevelDebug`，输出 HTTP 请求/响应元数据；为避免泄露 Token、SNMP 凭据等敏感信息，不输出请求和响应 Body：

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
sw, err := san.NewSANSwitch(host, username, password, san.WithLogger(logger))
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
| `GetDefinedZone(name)` | `/brocade-zone/defined-configuration/zone/zone-name/{name}` | 获取指定已定义 Zone |
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
| `SaveZoneConfig(checksum)` | 9.2+: `PATCH /brocade-zone/effective-configuration/cfg-action-v2/save`; 9.1: `PATCH /brocade-zone/effective-configuration/cfg-action/1` | 保存 Zone 配置 |
| `ActivateZoneConfig(name, checksum)` | `PATCH /brocade-zone/effective-configuration/cfg-name/{name}` | 激活 Zone 配置 |
| `CreateZoneAndActivate(cfg, zone, members, principalMembers)` | 组合流程 | 创建 Zone、加入 cfg、保存并激活 |
| `ReplaceZoneAndActivate(cfg, zone, members, principalMembers)` | 组合流程 | 全量替换 Zone 成员、保存并激活 |
| `DeleteZoneAndActivate(cfg, zone)` | 组合流程 | 删除 Zone、保存并激活 |
| `AbortZoneTransaction()` | 9.2+: `PATCH /brocade-zone/effective-configuration/cfg-action-v2/transaction-abort`; <=9.1: `PATCH /brocade-zone/effective-configuration/cfg-action/4` | 中止 Zone 事务 |
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

创建时如果传入 `principalMembers`，请求体会自动使用 `zone-type-string=user-created-peer-zone`：

```go
err := sw.CreateZone(
	"peer_zone_app01",
	[]string{
		"10:10:10:27:f8:f0:2a:e8",
		"10:10:10:27:f8:f0:3a:70",
		"10:10:10:27:f8:f0:38:65",
	},
	[]string{"10:10:10:27:f8:8f:44:cd"},
)
```

注意：`CreateZone` 会先检查同名 Zone 是否已存在，`UpdateZone` 和 `DeleteZone` 会先检查目标 Zone 是否存在。`UpdateZone` 和 `ReplaceZoneAndActivate` 会按 Brocade REST API 语义全量覆盖 `member-entry` leaf-list，请传入所有需要保留的成员。更新已有 Zone 时会沿用当前 defined zone 的 `zone-type-string`，不会在 `zone` 和 `user-created-peer-zone` 之间转换。

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
| `GetFDMIPorts()` | `/brocade-fdmi/port` | 获取 FDMI 端口信息 |

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
| Zone 管理 | 9 | 15 | 24 |
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
| **合计** | **39** | **15** | **54** |

## Context 支持

所有网络相关的 Client 和 SANSwitch 方法都提供 `WithContext` 变体，支持请求级取消与超时控制。无 Context 方法仅作为使用 `context.Background()` 的便捷封装：

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

sw, err := san.NewSANSwitchWithContext(ctx, host, username, password)
if err != nil {
	log.Fatal(err)
}
switches, err := sw.GetFabricSwitchesWithContext(ctx)
if err != nil {
	log.Fatal(err)
}
```

底层方法包括：`GetWithContext`、`PostWithContext`、`PatchWithContext`、`DeleteWithContext`、`LoginWithContext`、`LogoutWithContext`；所有高层采集、Zone 写操作和事务操作也提供对应方法。

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

预定义错误：`ErrNotFound`、`ErrUnauthorized`、`ErrConnectionFailed`、`ErrInvalidResponse`、`ErrTimeout`、`ErrUnsupportedOperation`、`ErrResponseBodyTooLarge`。

Zone 多步操作若在远端变更后失败，会返回 `PartialMutationError`；客户端会尝试自动中止事务，但调用方仍应检查交换机最终状态。

## 接口抽象

库不再暴露一个包含所有能力的巨型接口。`Session`、`SwitchReader`、`ZoneReader`、`ZoneWriter`、`InventoryReader` 和 `MonitoringReader` 均以 Context 方法为主，调用方可以按依赖范围选择：

```go
type PortReader interface {
	GetPortsWithContext(context.Context) ([]san.PortInfo, error)
}

func collectPorts(ctx context.Context, reader PortReader) ([]san.PortInfo, error) {
	return reader.GetPortsWithContext(ctx)
}
```

如果业务只需要一两个方法，建议像上面一样在业务包内定义更窄的接口；`SANSwitch` 已通过编译期断言实现这些能力接口。

## 技术栈

| 组件 | 选型 |
|------|------|
| 语言 | Go 1.21+ |
| HTTP 客户端 | `net/http`（Go 标准库） |
| 数据格式 | XML（Brocade YANG Data Model `application/yang-data+xml`） |
| 日志 | `log/slog`（Go 标准库） |
| 测试 | `net/http/httptest`（标准库） |

## 项目结构

```
├── client.go           # HTTP 客户端、认证、重试、Context 支持
├── doc.go              # 包文档
├── san.go              # SANSwitch facade
├── interfaces.go       # 细粒度能力接口
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
├── Makefile            # 本地质量检查入口
├── *_test.go           # 单元测试
├── example/
│   └── main.go         # 使用示例
└── docs/
    └── brocade-rest-api-doc/  # Brocade FOS REST API 参考文档
```

## 参考文档

- [Fabric OS REST API Overview](https://techdocs.broadcom.com/us/en/fibre-channel-networking/fabric-os/fabric-os-rest-api/9-2-x.html)
- Brocade FOS 9.2.x REST API Reference (fos-92x-restapi.pdf)

## 模块与发布

- 模块路径：`github.com/everstarmy/sanswitch`
- 支持 Go 1.21 及以上；库本身不锁定 `toolchain`，由调用方选择工具链补丁版本
- 运行时仅使用 Go 标准库，`go list -m all` 不应出现第三方模块
- 每次发布前执行 `go mod tidy -diff`、`go mod verify`、`go test -race ./...` 和 `go vet ./...`
- 本地可直接运行 `make check` 执行格式、模块、测试、竞态和 vet 检查
- 版本变更记录维护在 [CHANGELOG.md](CHANGELOG.md)，Git tag 遵循 SemVer
