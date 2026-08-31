# sanswitch

`sanswitch` 是一个仅使用 Go 标准库的 Brocade Fabric OS REST API 客户端，支持交换机、端口、FRU、Zone、SFP、MAPS、Trunk、SNMP、NTP 和日志采集。

## 设计原则

- 运行时零第三方依赖，无 `go.sum` 和 vendor
- Go 1.24+
- 所有网络操作都要求显式传入 `context.Context`
- 默认 HTTPS 并验证服务器证书
- 严格限制响应体大小
- 只对只读、可重试请求执行带抖动的指数退避
- 未知 Fabric OS 版本默认禁止写入
- Zone 写入使用显式事务和 checksum 乐观并发控制
- XML/YANG 传输结构不进入公共 API

## 安装

```bash
go get github.com/everstarmy/sanswitch@latest
```

## 创建并登录

`Open` 会创建客户端并立即登录。凭据只用于本次登录，不会保存在 `Client` 中。

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

client, err := sanswitch.Open(ctx, "switch.example.com", sanswitch.Credentials{
	Username: os.Getenv("SAN_USERNAME"),
	Password: os.Getenv("SAN_PASSWORD"),
})
if err != nil {
	log.Fatal(err)
}
defer client.Close()
defer func() {
	logoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = client.Logout(logoutCtx)
}()

ports, err := client.Ports(ctx)
if err != nil {
	log.Fatal(err)
}
```

如果希望分别控制构造和登录：

```go
client, err := sanswitch.New("switch.example.com")
if err != nil {
	return err
}
session, err := client.Login(ctx, credentials)
```

## 配置

```go
client, err := sanswitch.New("switch.example.com",
	sanswitch.WithTimeout(45*time.Second),
	sanswitch.WithRetry(4),
	sanswitch.WithRetryWait(500*time.Millisecond),
	sanswitch.WithRetryMaxWait(15*time.Second),
	sanswitch.WithMaxResponseBodyBytes(8<<20),
	sanswitch.WithLogger(logger),
	sanswitch.WithTLSConfig(tlsConfig),
)
```

| 选项 | 默认值 | 说明 |
| --- | --- | --- |
| `WithTimeout` | `30s` | HTTP 请求总超时 |
| `WithRetry` | `3` | 只读请求最大重试次数 |
| `WithRetryWait` | `1s` | 初始退避时间 |
| `WithRetryMaxWait` | `30s` | 最大退避时间 |
| `WithMaxResponseBodyBytes` | `16 MiB` | 响应体硬限制 |
| `WithLogger` | `slog.Default()` | 自定义结构化日志 |
| `WithTLSConfig` | 系统 TLS 配置 | 自定义 CA、客户端证书等 |
| `WithTransport` | 克隆的默认 Transport | 注入代理、追踪或测试 Transport |
| `WithFOSVersion` | 登录自动识别 | 显式指定 FOS 版本 |
| `WithHTTP` | 关闭 | 使用明文 HTTP，仅适合隔离测试环境 |
| `WithInsecureSkipVerify` | 关闭 | 禁用证书验证，仅适合受控测试 |
| `WithAllowUnknownFOSVersionWrites` | 关闭 | 显式允许未知版本写入 |

非法 endpoint、超时、重试参数、响应限制和 FOS 版本会让 `New` 返回 `ErrInvalidConfig` 或 `ErrInvalidVersion`，不会被静默修正。

## 常用读取 API

所有方法的第一个参数都是 `context.Context`。

| 分类 | 方法 |
| --- | --- |
| 交换机 | `Switch`、`FabricSwitches`、`Ports`、`Hardware`、`LogicalSwitches` |
| Zone | `DefinedZones`、`DefinedZone`、`EffectiveZones`、`DefinedAliases`、`DefinedConfigs`、`EffectiveConfig`、`ZoneDatabase` |
| FRU | `Blades`、`Fans`、`PowerSupplies`、`HistoryLogs`、`Sensors` |
| Fabric | `FibreChannelStatistics`、`FibreChannelNameServers`、`FDMIHBAs`、`FDMIPorts` |
| 监控 | `SwitchStatus`、`SystemResources`、`ErrorLogs`、`AuditLogs` |
| 其他 | `MediaRDPs`、`Trunks`、`TrunkPerformances`、`TrunkAreas`、`FirmwareHistory`、`SNMPSystem`、`TimeZone`、`ClockServer` |

调用方需要 mock 时，应在消费包声明最小接口：

```go
type portReader interface {
	Ports(context.Context) ([]sanswitch.Port, error)
}
```

## Zone 事务

组合写操作不再隐藏多个远端步骤。调用方显式管理事务：

```go
tx, err := client.BeginZoneTransaction(ctx)
if err != nil {
	return err
}

if err := tx.CreateZone(ctx, "zone_app_storage", []string{
	"10:00:00:00:00:00:00:01",
	"20:00:00:00:00:00:00:01",
}, nil); err != nil {
	_ = tx.Abort(ctx)
	return err
}
if err := tx.AddZoneToConfig(ctx, "cfg_prod", "zone_app_storage"); err != nil {
	_ = tx.Abort(ctx)
	return err
}
if err := tx.Commit(ctx, "cfg_prod"); err != nil {
	_ = tx.Abort(ctx)
	return err
}
```

`Commit` 使用事务开始时取得的 checksum。其他客户端或进程抢先修改配置时，Fabric OS 会拒绝保存。已经尝试过远端修改的错误会包装为 `*PartialMutationError`，可通过 `errors.As` 判断。

## Fabric OS 能力

```go
version, known := client.Version()
capabilities := client.Capabilities()
```

- FOS 9.0+：FRU history、sensor、logging
- FOS 9.1+：Zone 写入、firmware history
- FOS 9.2+：新版 Zone save/abort endpoint
- 版本未知：允许只读探测，默认禁止写入并返回 `ErrUnknownFOSVersion`

## 错误判断

```go
switch {
case errors.Is(err, sanswitch.ErrUnauthorized):
case errors.Is(err, sanswitch.ErrNotFound):
case errors.Is(err, sanswitch.ErrTimeout):
case errors.Is(err, sanswitch.ErrUnknownFOSVersion):
case errors.Is(err, sanswitch.ErrUnsupportedOperation):
}

var apiErr *sanswitch.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.StatusCode, apiErr.ErrorCode, apiErr.Message)
}
```

## 验证

```bash
make check
```

该命令执行格式检查、module 纯净度检查、普通测试、随机顺序测试、race、覆盖率阈值和 `go vet`。
