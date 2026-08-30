package san

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 默认配置常量
const (
	DefaultTimeout          = 30 * time.Second
	DefaultRetryCount       = 3
	DefaultRetryWaitTime    = 1 * time.Second
	DefaultRetryMaxWaitTime = 30 * time.Second
	// DefaultMaxResponseBodyBytes limits the amount of response data retained
	// for XML decoding and error reporting.
	DefaultMaxResponseBodyBytes int64 = 16 << 20
	zoneCleanupTimeout                = 5 * time.Second
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

// WithMaxResponseBodyBytes limits the maximum response body size accepted by
// the client. A non-positive value leaves the default limit unchanged.
func WithMaxResponseBodyBytes(maxBytes int64) ClientOption {
	return func(c *Client) {
		if maxBytes > 0 {
			c.maxResponseBodyBytes = maxBytes
		}
	}
}

// WithLogger 注入自定义的结构化日志 logger
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
		c.customLogger = logger != nil
	}
}

// WithTLSConfig 设置 HTTPS 请求使用的 TLS 配置。
// 配置会在创建 Client 时复制，调用方之后可以安全地复用或修改原配置。
func WithTLSConfig(config *tls.Config) ClientOption {
	return func(c *Client) {
		if config == nil {
			c.tlsConfig = nil
			return
		}
		c.tlsConfig = config.Clone()
	}
}

// WithInsecureSkipVerify 显式关闭 HTTPS 服务器证书校验。
// 仅建议用于测试或明确使用自签名证书且有其他网络隔离措施的环境。
func WithInsecureSkipVerify() ClientOption {
	return func(c *Client) {
		// This is intentionally opt-in; the secure default leaves verification enabled.
		c.tlsConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402 -- explicitly requested by the caller
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
		c.fosVersion = strings.TrimSpace(version)
		c.fosVersionConfigured = c.fosVersion != ""
	}
}

// WithAllowUnknownFOSVersionWrites explicitly permits write requests when the
// FOS version is unknown. The safer default is to block writes until the
// version is supplied explicitly or learned during login.
func WithAllowUnknownFOSVersionWrites() ClientOption {
	return func(c *Client) {
		c.allowUnknownFOSVersionWrites = true
	}
}

// Client 是 Brocade FOS REST API 的 HTTP 客户端，封装了认证、请求发送、重试、
// 日志、Virtual Fabric 路由等核心功能。
// 通过 NewClient 创建实例，或使用 NewSANSwitch 自动完成登录。
type Client struct {
	host                         string
	username                     string
	password                     string
	authToken                    string
	client                       *http.Client
	logger                       *slog.Logger
	baseURL                      string // 测试用：覆盖默认 URL 前缀
	tlsConfig                    *tls.Config
	vfID                         int
	useHTTP                      bool
	logOutput                    io.Writer
	verbose                      bool
	customLogger                 bool
	timeout                      time.Duration
	retryCount                   int
	retryWaitTime                time.Duration
	retryMaxWaitTime             time.Duration
	maxResponseBodyBytes         int64
	fosVersion                   string
	fosVersionConfigured         bool
	allowUnknownFOSVersionWrites bool
	stateMu                      sync.RWMutex
	loggerMu                     sync.RWMutex
	zoneWriteMu                  sync.Mutex
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

type apiResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// NewClient 创建一个新的 FOS REST API 客户端实例（不自动登录）。
// 默认使用 HTTPS（校验证书）、30 秒超时、3 次指数退避重试。
// 可通过 ClientOption 函数选项自定义超时、重试、日志和协议等配置。
func NewClient(host, username, password string, opts ...ClientOption) *Client {
	c := &Client{
		host:                 host,
		username:             username,
		password:             password,
		logOutput:            os.Stderr,
		timeout:              DefaultTimeout,
		retryCount:           DefaultRetryCount,
		retryWaitTime:        DefaultRetryWaitTime,
		retryMaxWaitTime:     DefaultRetryMaxWaitTime,
		maxResponseBodyBytes: DefaultMaxResponseBodyBytes,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	if c.timeout <= 0 {
		c.timeout = DefaultTimeout
	}
	if c.retryCount < 0 {
		c.retryCount = 0
	}
	if c.retryWaitTime < 0 {
		c.retryWaitTime = 0
	}
	if c.retryMaxWaitTime < c.retryWaitTime {
		c.retryMaxWaitTime = c.retryWaitTime
	}
	if c.maxResponseBodyBytes <= 0 {
		c.maxResponseBodyBytes = DefaultMaxResponseBodyBytes
	}

	// 若未注入自定义 logger，则使用默认 logger。
	if c.logger == nil {
		c.logger = slog.Default()
	}

	c.client = &http.Client{
		Transport: cloneHTTPTransport(c.tlsConfig),
		Timeout:   c.timeout,
	}
	return c
}

// Timeout 返回当前配置的请求超时时间
func (c *Client) Timeout() time.Duration {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.timeout
}

// RetryCount 返回当前配置的最大重试次数
func (c *Client) RetryCount() int {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.retryCount
}

// Close releases idle HTTP connections held by the client.
// It is safe to call Close more than once.
func (c *Client) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	c.client.CloseIdleConnections()
	return nil
}

// SetVerbose 开启或关闭调试级别日志输出
func (c *Client) SetVerbose(verbose bool) {
	c.loggerMu.Lock()
	defer c.loggerMu.Unlock()
	c.verbose = verbose
	if c.customLogger {
		return
	}
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
	c.loggerMu.Lock()
	defer c.loggerMu.Unlock()
	if w == nil {
		c.logOutput = os.Stderr
	} else {
		c.logOutput = w
	}
	if c.verbose && !c.customLogger {
		c.logger = slog.New(slog.NewTextHandler(c.logOutput, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
}

// SetVFID 设置虚拟 Fabric ID（vf-id 查询参数），用于 Virtual Fabric 场景下的请求路由。
// vfID <= 0 时不会在请求中附加 vf-id 参数。
func (c *Client) SetVFID(vfID int) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.vfID = vfID
}

// Login 使用 Basic Auth 登录交换机并获取认证 Token。
// 成功后将 Token 设置到后续请求的 Authorization 头中。
// 对应 API: POST /rest/login
func (c *Client) Login() (*LoginResponse, error) {
	return c.LoginWithContext(context.Background())
}

// LoginWithContext 使用 Basic Auth 登录交换机并获取认证 Token。
// ctx 取消时会中止当前登录请求。
func (c *Client) LoginWithContext(ctx context.Context) (*LoginResponse, error) {
	ctx = nonNilContext(ctx)
	c.clearAuth()
	url := c.restBase() + c.endpoints().Login()
	headers := c.baseHeaders()
	headers.Set("Authorization", "Basic "+base64Encode(c.username+":"+c.password))

	resp, err := c.do(ctx, http.MethodPost, url, nil, headers, false)

	if err != nil {
		return nil, fmt.Errorf("login: %w", wrapRequestError(err))
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("login: %w", c.responseError(resp))
	}

	authHeader := resp.Header.Get("Authorization")
	if authHeader == "" {
		authHeader = resp.Header.Get("authorization")
	}
	if authHeader == "" {
		return nil, fmt.Errorf("%w: login response missing authorization header", ErrInvalidResponse)
	}

	if len(resp.Body) == 0 {
		c.setAuthenticated(authHeader, "")
		return &LoginResponse{}, nil
	}

	var result LoginResponse
	if err := xml.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("login response: %w", invalidResponseError(err))
	}
	c.setAuthenticated(authHeader, result.FirmwareVersion)

	return &result, nil
}

// Logout 注销当前会话，清除认证 Token。
// 对应 API: POST /rest/logout
func (c *Client) Logout() error {
	return c.LogoutWithContext(context.Background())
}

// LogoutWithContext 注销当前会话并清除认证 Token。
// ctx 取消时会中止当前注销请求。
func (c *Client) LogoutWithContext(ctx context.Context) error {
	ctx = nonNilContext(ctx)
	url := c.restBase() + c.endpoints().Logout()

	resp, err := c.do(ctx, http.MethodPost, url, nil, c.authHeaders(), false)

	if err != nil {
		c.clearAuth()
		return fmt.Errorf("logout: %w", wrapRequestError(err))
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		c.clearAuth()
		return fmt.Errorf("logout: %w", c.responseError(resp))
	}

	c.clearAuth()

	return nil
}

func (c *Client) buildURL(endpoint string) string {
	var rawURL string
	if c.baseURL != "" {
		rawURL = c.baseURL + endpoint
	} else {
		rawURL = fmt.Sprintf("%s://%s/rest/running%s", c.scheme(), c.host, endpoint)
	}
	c.stateMu.RLock()
	vfID := c.vfID
	c.stateMu.RUnlock()
	if vfID > 0 {
		parsed, err := url.Parse(rawURL)
		if err == nil {
			query := parsed.Query()
			query.Set("vf-id", fmt.Sprintf("%d", vfID))
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
		baseURL := strings.TrimRight(c.baseURL, "/")
		if strings.HasSuffix(baseURL, "/rest/running") {
			return strings.TrimSuffix(baseURL, "/running")
		}
		return baseURL
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

// Get 执行 GET 请求并将 XML 响应解析到 result 中。
func (c *Client) Get(endpoint string, result any) error {
	return c.GetWithContext(context.Background(), endpoint, result)
}

// Post 执行 POST 请求，将 payload 序列化为 XML 发送到指定端点。
func (c *Client) Post(endpoint string, payload any) error {
	return c.PostWithContext(context.Background(), endpoint, payload)
}

// Patch 执行 PATCH 请求，将 payload 序列化为 XML 发送到指定端点。
func (c *Client) Patch(endpoint string, payload any) error {
	return c.PatchWithContext(context.Background(), endpoint, payload)
}

// Delete 执行 DELETE 请求，删除指定端点的资源。
func (c *Client) Delete(endpoint string) error {
	return c.DeleteWithContext(context.Background(), endpoint)
}

// IsLoggedIn 返回当前是否已持有认证 Token。
func (c *Client) IsLoggedIn() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
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

// cloneHTTPTransport returns an isolated transport so TLS settings do not
// mutate the process-wide http.DefaultTransport.
func cloneHTTPTransport(tlsConfig *tls.Config) http.RoundTripper {
	if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
		transport := defaultTransport.Clone()
		if tlsConfig != nil {
			transport.TLSClientConfig = tlsConfig.Clone()
		}
		return transport
	}
	if tlsConfig != nil {
		return &http.Transport{TLSClientConfig: tlsConfig.Clone()}
	}
	return http.DefaultTransport
}

// retryAfter returns a bounded server-provided or exponential backoff delay.
func (c *Client) retryAfter(resp *apiResponse, attempt int) time.Duration {
	if resp != nil {
		if wait, ok := parseRetryAfter(resp.Header.Get("Retry-After")); ok {
			return min(wait, c.retryMaxWaitTime)
		}
	}

	wait := c.retryWaitTime
	for i := 1; i < attempt && wait < c.retryMaxWaitTime; i++ {
		if wait > c.retryMaxWaitTime/2 {
			wait = c.retryMaxWaitTime
			break
		}
		wait *= 2
	}
	return min(wait, c.retryMaxWaitTime)
}

func parseRetryAfter(value string) (time.Duration, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil && seconds >= 0 {
		maxSeconds := int64((time.Duration(1<<63 - 1)) / time.Second)
		if seconds > maxSeconds {
			return time.Duration(1<<63 - 1), true
		}
		return time.Duration(seconds) * time.Second, true
	}
	when, err := http.ParseTime(value)
	if err != nil {
		return 0, false
	}
	wait := time.Until(when)
	if wait < 0 {
		return 0, true
	}
	return wait, true
}

// parseAPIError 尝试将非 200 响应解析为 FOS 结构化错误；若解析失败则回退到通用错误
func (c *Client) parseAPIError(resp *apiResponse) error {
	apiErr := &APIError{StatusCode: resp.StatusCode}
	body := resp.Body
	if err := xml.Unmarshal(body, apiErr); err != nil || apiErr.Message == "" {
		return fmt.Errorf("api request failed with status %d: %s", resp.StatusCode, responseSnippet(body))
	}
	return apiErr
}

func (c *Client) responseError(resp *apiResponse) error {
	if resp.StatusCode == http.StatusUnauthorized {
		c.clearAuth()
		return ErrUnauthorized
	}
	err := c.parseAPIError(resp)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	}
	return err
}

func (c *Client) baseHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/yang-data+xml")
	headers.Set("Accept", "application/yang-data+xml")
	return headers
}

func (c *Client) authHeaders() http.Header {
	headers := c.baseHeaders()
	c.stateMu.RLock()
	authToken := c.authToken
	c.stateMu.RUnlock()
	if authToken != "" {
		headers.Set("Authorization", authToken)
	}
	return headers
}

func (c *Client) newRequest(ctx context.Context, method, url string, body []byte, headers http.Header) (*http.Request, error) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(nonNilContext(ctx), method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, nil
}

func (c *Client) do(ctx context.Context, method, url string, body []byte, headers http.Header, retry bool) (*apiResponse, error) {
	ctx = nonNilContext(ctx)
	retryCount := c.RetryCount()
	if !retry || !isRetryableMethod(method) {
		retryCount = 0
	}
	attempts := retryCount + 1

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		req, err := c.newRequest(ctx, method, url, body, headers)
		if err != nil {
			return nil, err
		}
		resp, err := c.client.Do(req)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			lastErr = err
			if attempt >= attempts {
				break
			}
			if err := waitForRetry(ctx, c.retryAfter(nil, attempt)); err != nil {
				return nil, err
			}
			continue
		}

		responseBody, readErr := readResponseBody(resp.Body, c.maxResponseBodyBytes)
		closeErr := resp.Body.Close()
		if readErr != nil {
			if errors.Is(readErr, ErrResponseBodyTooLarge) {
				return nil, readErr
			}
			lastErr = readErr
			if attempt >= attempts {
				break
			}
			if err := waitForRetry(ctx, c.retryAfter(nil, attempt)); err != nil {
				return nil, err
			}
			continue
		}
		if closeErr != nil {
			lastErr = closeErr
			if attempt >= attempts {
				break
			}
			if err := waitForRetry(ctx, c.retryAfter(nil, attempt)); err != nil {
				return nil, err
			}
			continue
		}

		apiResp := &apiResponse{
			StatusCode: resp.StatusCode,
			Header:     resp.Header.Clone(),
			Body:       responseBody,
		}
		if isRetryableStatus(apiResp.StatusCode) && attempt < attempts {
			if err := waitForRetry(ctx, c.retryAfter(apiResp, attempt)); err != nil {
				return nil, err
			}
			continue
		}
		return apiResp, nil
	}

	return nil, lastErr
}

func isRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}

func waitForRetry(ctx context.Context, wait time.Duration) error {
	if wait <= 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) setAuthenticated(authToken, firmwareVersion string) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.authToken = authToken
	if firmwareVersion != "" {
		c.fosVersion = firmwareVersion
	} else if !c.fosVersionConfigured {
		c.fosVersion = legacyFOSVersion
	}
}

func (c *Client) clearAuth() {
	c.stateMu.Lock()
	c.authToken = ""
	c.stateMu.Unlock()
}

func (c *Client) debug(message string, args ...any) {
	c.loggerMu.RLock()
	logger := c.logger
	c.loggerMu.RUnlock()
	if logger != nil {
		logger.Debug(message, args...)
	}
}

func nonNilContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func wrapRequestError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return err
	}
	if errors.Is(err, ErrResponseBodyTooLarge) {
		return err
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %w", ErrTimeout, err)
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Errorf("%w: %w", ErrTimeout, err)
	}
	return fmt.Errorf("%w: %w", ErrConnectionFailed, err)
}

func invalidResponseError(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidResponse, err)
}

func responseSnippet(body []byte) string {
	const maxErrorBodyBytes = 4 << 10
	if len(body) <= maxErrorBodyBytes {
		return string(body)
	}
	return string(body[:maxErrorBodyBytes]) + "..."
}

func readResponseBody(body io.Reader, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 || maxBytes == math.MaxInt64 {
		return io.ReadAll(body)
	}

	responseBody, err := io.ReadAll(io.LimitReader(body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(responseBody)) > maxBytes {
		return nil, fmt.Errorf("%w: limit %d bytes", ErrResponseBodyTooLarge, maxBytes)
	}
	return responseBody, nil
}

// ---------- WithContext 变体：支持请求级 context 取消与超时控制 ----------

// GetWithContext 执行带 context 的 GET 请求
func (c *Client) GetWithContext(ctx context.Context, endpoint string, result any) error {
	ctx = nonNilContext(ctx)
	url := c.buildURL(endpoint)

	resp, err := c.do(ctx, http.MethodGet, url, nil, c.authHeaders(), true)
	if err != nil {
		return wrapRequestError(err)
	}

	if resp.StatusCode != http.StatusOK {
		return c.responseError(resp)
	}

	c.debug("GET response", "url", url, "status", resp.StatusCode, "response_bytes", len(resp.Body))

	if err := xml.Unmarshal(resp.Body, result); err != nil {
		return invalidResponseError(err)
	}
	return nil
}

// PostWithContext 执行带 context 的 POST 请求
func (c *Client) PostWithContext(ctx context.Context, endpoint string, payload any) error {
	return c.mutateWithContext(nonNilContext(ctx), http.MethodPost, endpoint, payload, true)
}

// PatchWithContext 执行带 context 的 PATCH 请求
func (c *Client) PatchWithContext(ctx context.Context, endpoint string, payload any) error {
	return c.mutateWithContext(nonNilContext(ctx), http.MethodPatch, endpoint, payload, true)
}

// DeleteWithContext 执行带 context 的 DELETE 请求
func (c *Client) DeleteWithContext(ctx context.Context, endpoint string) error {
	return c.mutateWithContext(nonNilContext(ctx), http.MethodDelete, endpoint, nil, true)
}

func (c *Client) patchWithoutVersionGate(endpoint string, payload any) error {
	return c.patchWithoutVersionGateWithContext(context.Background(), endpoint, payload)
}

func (c *Client) patchWithoutVersionGateWithContext(ctx context.Context, endpoint string, payload any) error {
	return c.mutateWithContext(nonNilContext(ctx), http.MethodPatch, endpoint, payload, false)
}

func (c *Client) mutateWithContext(ctx context.Context, method, endpoint string, payload any, enforceVersionGate bool) error {
	if enforceVersionGate {
		if err := c.ensureWriteSupported(); err != nil {
			return err
		}
	}

	url := c.buildURL(endpoint)
	reqBody, err := marshalPayload(payload)
	if err != nil {
		return err
	}

	c.debug(method+" request", "url", url, "payload_bytes", len(reqBody))

	if method != http.MethodPost && method != http.MethodPatch && method != http.MethodDelete {
		return fmt.Errorf("unsupported mutation method %q", method)
	}
	resp, err := c.do(ctx, method, url, reqBody, c.authHeaders(), false)
	if err != nil {
		return wrapRequestError(err)
	}

	if !mutationStatusOK(method, resp.StatusCode) {
		return c.responseError(resp)
	}

	c.debug(method+" response", "url", url, "status", resp.StatusCode, "response_bytes", len(resp.Body))
	return nil
}

func marshalPayload(payload any) ([]byte, error) {
	if payload == nil {
		return nil, nil
	}
	body, err := xml.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request payload: %w", err)
	}
	return body, nil
}

func mutationStatusOK(method string, status int) bool {
	if method == http.MethodDelete {
		return status == http.StatusOK || status == http.StatusNoContent
	}
	return status == http.StatusOK || status == http.StatusCreated || status == http.StatusNoContent
}

// base64Encode 对字符串进行 Base64 编码，用于 Basic Auth 认证
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (c *Client) ensureWriteSupported() error {
	if c.endpoints().allowWrite() {
		return nil
	}
	return fmt.Errorf("%w: FOS %s does not support write operations", ErrUnsupportedOperation, c.endpoints().version)
}
