package san

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// 默认配置常量
const (
	DefaultTimeout          = 30 * time.Second
	DefaultRetryCount       = 3
	DefaultRetryWaitTime    = 1 * time.Second
	DefaultRetryMaxWaitTime = 30 * time.Second
)

// ClientOption 是一个函数选项类型，用于配置 Client 的参数
type ClientOption func(*Client)

// WithTimeout 设置 HTTP 请求的超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithRetry 设置请求的重试次数（0 表示不重试）
func WithRetry(count int) ClientOption {
	return func(c *Client) {
		c.retryCount = count
	}
}

// WithRetryWait 设置重试的初始等待时间（用于指数退避计算起点）
func WithRetryWait(wait time.Duration) ClientOption {
	return func(c *Client) {
		c.retryWaitTime = wait
	}
}

// WithRetryMaxWait 设置重试的最大等待时间上限
func WithRetryMaxWait(maxWait time.Duration) ClientOption {
	return func(c *Client) {
		c.retryMaxWaitTime = maxWait
	}
}

// WithLogger 注入自定义的结构化日志 logger
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithHTTP 使用 HTTP 替代 HTTPS（不推荐用于生产环境）
func WithHTTP() ClientOption {
	return func(c *Client) {
		c.useHTTP = true
	}
}

// WithFOSVersion 设置 Fabric OS 版本，用于兼容不同版本的 REST endpoint。
// NewSANSwitch 会在登录后自动记录版本；直接使用 NewClient 时可用该选项显式指定。
func WithFOSVersion(version string) ClientOption {
	return func(c *Client) {
		c.fosVersion = version
	}
}

// Client 是 Brocade FOS REST API 的 HTTP 客户端，封装了认证、请求发送、重试、
// 日志、Virtual Fabric 路由等核心功能。
// 通过 NewClient 创建实例，或使用 NewSANSwitch 自动完成登录。
type Client struct {
	host             string
	username         string
	password         string
	authToken        string
	client           *resty.Client
	logger           *slog.Logger
	baseURL          string // 测试用：覆盖默认 URL 前缀
	vfID             int
	useHTTP          bool
	logOutput        io.Writer
	timeout          time.Duration
	retryCount       int
	retryWaitTime    time.Duration
	retryMaxWaitTime time.Duration
	fosVersion       string
}

// LoginResponse 是 POST /login 的 XML 响应，包含登录后的用户信息和交换机参数
type LoginResponse struct {
	XMLName           xml.Name `xml:"Response"`
	UserName          string   `xml:"switch-parameters>user-name"`
	ChassisAccessRole string   `xml:"switch-parameters>chassis-access-role"`
	HomeVirtualFabric int      `xml:"switch-parameters>home-virtual-fabric"`
	FirmwareVersion   string   `xml:"switch-parameters>firmware-version"`
	Model             string   `xml:"switch-parameters>model"`
}

// NewClient 创建一个新的 FOS REST API 客户端实例（不自动登录）。
// 默认使用 HTTPS（跳过证书验证）、30 秒超时、3 次指数退避重试。
// 可通过 ClientOption 函数选项自定义超时、重试、日志和协议等配置。
func NewClient(host, username, password string, opts ...ClientOption) *Client {
	c := &Client{
		host:             host,
		username:         username,
		password:         password,
		logOutput:        os.Stderr,
		timeout:          DefaultTimeout,
		retryCount:       DefaultRetryCount,
		retryWaitTime:    DefaultRetryWaitTime,
		retryMaxWaitTime: DefaultRetryMaxWaitTime,
	}

	for _, opt := range opts {
		opt(c)
	}

	// 若未注入自定义 logger，则使用默认的空 logger
	if c.logger == nil {
		c.logger = slog.Default()
	}

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetHeader("Content-Type", "application/yang-data+xml")
	client.SetHeader("Accept", "application/yang-data+xml")
	client.SetTimeout(c.timeout)

	// 配置重试策略与指数退避
	client.SetRetryCount(c.retryCount)
	client.SetRetryWaitTime(c.retryWaitTime)
	client.SetRetryMaxWaitTime(c.retryMaxWaitTime)
	// 指数退避: 重试间隔时间 = baseWait * 2^(attempt-1)
	client.SetRetryAfter(func(_ *resty.Client, resp *resty.Response) (time.Duration, error) {
		attempt := resp.Request.Attempt
		backoff := float64(c.retryWaitTime) * math.Pow(2, float64(attempt-1))
		wait := time.Duration(backoff)
		if wait > c.retryMaxWaitTime {
			wait = c.retryMaxWaitTime
		}
		return wait, nil
	})
	// 在网络错误、429 限流或 5xx 服务端错误时触发重试
	client.AddRetryCondition(func(r *resty.Response, err error) bool {
		if r == nil || r.Request == nil || !isRetryableMethod(r.Request.Method) {
			return false
		}
		if err != nil {
			return true
		}
		return r.StatusCode() == http.StatusTooManyRequests ||
			r.StatusCode() >= http.StatusInternalServerError
	})

	c.client = client
	return c
}

// Timeout 返回当前配置的请求超时时间
func (c *Client) Timeout() time.Duration {
	return c.timeout
}

// RetryCount 返回当前配置的最大重试次数
func (c *Client) RetryCount() int {
	return c.retryCount
}

// SetVerbose 开启或关闭调试级别日志输出
func (c *Client) SetVerbose(verbose bool) {
	if verbose {
		c.logger = slog.New(slog.NewTextHandler(c.logOutput, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		// 使用 discard handler 抑制日志输出
		c.logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
}

// SetLogOutput 设置日志输出目标（nil 则恢复为 os.Stderr）
func (c *Client) SetLogOutput(w io.Writer) {
	if w == nil {
		c.logOutput = os.Stderr
		return
	}
	c.logOutput = w
}

// SetVFID 设置虚拟 Fabric ID（vf-id 查询参数），用于 Virtual Fabric 场景下的请求路由。
// vfID <= 0 时不会在请求中附加 vf-id 参数。
func (c *Client) SetVFID(vfID int) {
	c.vfID = vfID
}

// Login 使用 Basic Auth 登录交换机并获取认证 Token。
// 成功后将 Token 设置到后续请求的 Authorization 头中。
// 对应 API: POST /rest/login
func (c *Client) Login() (*LoginResponse, error) {
	url := c.restBase() + c.endpoints().Login()

	resp, err := c.client.R().
		SetHeader("Authorization", "Basic "+base64Encode(c.username+":"+c.password)).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	authHeader := resp.Header().Get("Authorization")
	if authHeader == "" {
		authHeader = resp.Header().Get("authorization")
	}

	if authHeader != "" {
		c.authToken = authHeader
		c.client.SetHeader("Authorization", c.authToken)
	}

	if len(resp.Body()) == 0 {
		c.fosVersion = legacyFOSVersion
		return &LoginResponse{}, nil
	}

	var result LoginResponse
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("login response: %w", err)
	}
	c.fosVersion = result.FirmwareVersion

	return &result, nil
}

// Logout 注销当前会话，清除认证 Token。
// 对应 API: POST /rest/logout
func (c *Client) Logout() error {
	url := c.restBase() + c.endpoints().Logout()

	resp, err := c.client.R().Post(url)

	if err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("logout failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	c.authToken = ""
	c.client.Header.Del("Authorization")

	return nil
}

func (c *Client) buildURL(endpoint string) string {
	var rawURL string
	if c.baseURL != "" {
		rawURL = c.baseURL + endpoint
	} else {
		rawURL = fmt.Sprintf("%s://%s/rest/running%s", c.scheme(), c.host, endpoint)
	}
	if c.vfID > 0 {
		parsed, err := url.Parse(rawURL)
		if err == nil {
			query := parsed.Query()
			query.Set("vf-id", fmt.Sprintf("%d", c.vfID))
			parsed.RawQuery = query.Encode()
			rawURL = parsed.String()
		}
	}
	return rawURL
}

// restBase 返回 REST API 的基础 URL 前缀（不含 /running），用于 login/logout 等端点
func (c *Client) restBase() string {
	if c.baseURL != "" {
		// baseURL 格式: http://host/rest/running  →  截取到 /rest
		return c.baseURL[:len(c.baseURL)-len("/running")]
	}
	return fmt.Sprintf("%s://%s/rest", c.scheme(), c.host)
}

// scheme 根据配置返回 URL 协议
func (c *Client) scheme() string {
	if c.useHTTP {
		return "http"
	}
	return "https"
}

// Get 执行 GET 请求并将 XML 响应解析到 result 中
func (c *Client) Get(endpoint string, result interface{}) error {
	url := c.buildURL(endpoint)

	resp, err := c.client.R().Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return c.parseAPIError(resp)
	}

	c.logger.Debug("GET response", "url", url, "body", string(resp.Body()))

	err = xml.Unmarshal(resp.Body(), result)
	if err != nil {
		return fmt.Errorf("failed to parse XML response: %w", err)
	}

	return nil
}

// Post 执行 POST 请求，将 payload 序列化为 XML 发送到指定端点
func (c *Client) Post(endpoint string, payload interface{}) error {
	if err := c.ensureWriteSupported(); err != nil {
		return err
	}
	url := c.buildURL(endpoint)

	var reqBody []byte
	var err error
	if payload != nil {
		reqBody, err = xml.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal request payload: %w", err)
		}
	}

	req := c.client.R()
	if payload != nil {
		req.SetBody(reqBody)
	}

	c.logger.Debug("POST request", "url", url, "payload", string(reqBody))

	resp, err := req.Post(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated && resp.StatusCode() != http.StatusNoContent {
		if resp.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return c.parseAPIError(resp)
	}

	c.logger.Debug("POST response", "url", url, "body", string(resp.Body()))

	return nil
}

// Patch 执行 PATCH 请求，将 payload 序列化为 XML 发送到指定端点
func (c *Client) Patch(endpoint string, payload interface{}) error {
	if err := c.ensureWriteSupported(); err != nil {
		return err
	}
	url := c.buildURL(endpoint)

	var reqBody []byte
	var err error
	if payload != nil {
		reqBody, err = xml.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal request payload: %w", err)
		}
	}

	req := c.client.R()
	if payload != nil {
		req.SetBody(reqBody)
	}

	c.logger.Debug("PATCH request", "url", url, "payload", string(reqBody))

	resp, err := req.Patch(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated && resp.StatusCode() != http.StatusNoContent {
		if resp.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return c.parseAPIError(resp)
	}

	c.logger.Debug("PATCH response", "url", url, "body", string(resp.Body()))

	return nil
}

// Delete 执行 DELETE 请求，删除指定端点的资源
func (c *Client) Delete(endpoint string) error {
	url := c.buildURL(endpoint)

	c.logger.Debug("DELETE request", "url", url)

	resp, err := c.client.R().Delete(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		if resp.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return c.parseAPIError(resp)
	}

	c.logger.Debug("DELETE response", "url", url, "body", string(resp.Body()))

	return nil
}

// IsLoggedIn 返回当前是否已持有有效的认证 Token
func (c *Client) IsLoggedIn() bool {
	return c.authToken != ""
}

// isRetryableMethod 判断 HTTP 方法是否支持自动重试（仅 GET/HEAD/OPTIONS 等幂等方法）
func isRetryableMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

// parseAPIError 尝试将非 200 响应解析为 FOS 结构化错误；若解析失败则回退到通用错误
func (c *Client) parseAPIError(resp *resty.Response) error {
	apiErr := &APIError{StatusCode: resp.StatusCode()}
	if err := xml.Unmarshal(resp.Body(), apiErr); err != nil || apiErr.Message == "" {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}
	return apiErr
}

// ---------- WithContext 变体：支持请求级 context 取消与超时控制 ----------

// GetWithContext 执行带 context 的 GET 请求
func (c *Client) GetWithContext(ctx context.Context, endpoint string, result interface{}) error {
	url := c.buildURL(endpoint)

	resp, err := c.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return c.parseAPIError(resp)
	}

	c.logger.Debug("GET response", "url", url, "body", string(resp.Body()))

	if err := xml.Unmarshal(resp.Body(), result); err != nil {
		return fmt.Errorf("failed to parse XML response: %w", err)
	}
	return nil
}

// PostWithContext 执行带 context 的 POST 请求
func (c *Client) PostWithContext(ctx context.Context, endpoint string, payload interface{}) error {
	if err := c.ensureWriteSupported(); err != nil {
		return err
	}
	url := c.buildURL(endpoint)

	var reqBody []byte
	var err error
	if payload != nil {
		reqBody, err = xml.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal request payload: %w", err)
		}
	}

	req := c.client.R().SetContext(ctx)
	if payload != nil {
		req.SetBody(reqBody)
	}

	c.logger.Debug("POST request", "url", url, "payload", string(reqBody))

	resp, err := req.Post(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated && resp.StatusCode() != http.StatusNoContent {
		if resp.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return c.parseAPIError(resp)
	}

	c.logger.Debug("POST response", "url", url, "body", string(resp.Body()))
	return nil
}

// PatchWithContext 执行带 context 的 PATCH 请求
func (c *Client) PatchWithContext(ctx context.Context, endpoint string, payload interface{}) error {
	if err := c.ensureWriteSupported(); err != nil {
		return err
	}
	url := c.buildURL(endpoint)

	var reqBody []byte
	var err error
	if payload != nil {
		reqBody, err = xml.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal request payload: %w", err)
		}
	}

	req := c.client.R().SetContext(ctx)
	if payload != nil {
		req.SetBody(reqBody)
	}

	c.logger.Debug("PATCH request", "url", url, "payload", string(reqBody))

	resp, err := req.Patch(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated && resp.StatusCode() != http.StatusNoContent {
		if resp.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return c.parseAPIError(resp)
	}

	c.logger.Debug("PATCH response", "url", url, "body", string(resp.Body()))
	return nil
}

// DeleteWithContext 执行带 context 的 DELETE 请求
func (c *Client) DeleteWithContext(ctx context.Context, endpoint string) error {
	url := c.buildURL(endpoint)

	c.logger.Debug("DELETE request", "url", url)

	resp, err := c.client.R().SetContext(ctx).Delete(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		if resp.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return c.parseAPIError(resp)
	}

	c.logger.Debug("DELETE response", "url", url, "body", string(resp.Body()))
	return nil
}

// base64Encode 对字符串进行 Base64 编码，用于 Basic Auth 认证
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (c *Client) ensureWriteSupported() error {
	if c.endpoints().allowWrite() {
		return nil
	}
	return fmt.Errorf("%w: FOS %s does not support POST/PATCH operations", ErrUnsupportedOperation, c.endpoints().version)
}
