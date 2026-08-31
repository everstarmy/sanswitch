package sanswitch

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand/v2"
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

// Credentials contains the short-lived credentials used for Login.
type Credentials struct {
	Username string
	Password string
}

type clientOptions struct {
	timeout                      time.Duration
	retryCount                   int
	retryWaitTime                time.Duration
	retryMaxWaitTime             time.Duration
	maxResponseBodyBytes         int64
	logger                       *slog.Logger
	tlsConfig                    *tls.Config
	transport                    http.RoundTripper
	useHTTP                      bool
	fosVersion                   string
	allowUnknownFOSVersionWrites bool
}

// ClientOption configures a Client before it is constructed.
type ClientOption func(*clientOptions)

// WithTimeout 设置 HTTP 请求的超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithRetry 设置请求的重试次数（0 表示不重试）
func WithRetry(count int) ClientOption {
	return func(o *clientOptions) {
		o.retryCount = count
	}
}

// WithRetryWait 设置重试的初始等待时间（用于指数退避计算起点）
func WithRetryWait(wait time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.retryWaitTime = wait
	}
}

// WithRetryMaxWait 设置重试的最大等待时间上限
func WithRetryMaxWait(maxWait time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.retryMaxWaitTime = maxWait
	}
}

// WithMaxResponseBodyBytes limits the maximum response body size accepted by
// the client. New rejects non-positive values.
func WithMaxResponseBodyBytes(maxBytes int64) ClientOption {
	return func(o *clientOptions) {
		o.maxResponseBodyBytes = maxBytes
	}
}

// WithLogger 注入自定义的结构化日志 logger
func WithLogger(logger *slog.Logger) ClientOption {
	return func(o *clientOptions) {
		o.logger = logger
	}
}

// WithTLSConfig 设置 HTTPS 请求使用的 TLS 配置。
// 配置会在创建 Client 时复制，调用方之后可以安全地复用或修改原配置。
func WithTLSConfig(config *tls.Config) ClientOption {
	return func(o *clientOptions) {
		if config == nil {
			o.tlsConfig = nil
			return
		}
		o.tlsConfig = config.Clone()
	}
}

// WithTransport injects the HTTP transport used by the client. It is useful
// for proxies, observability, and deterministic tests. It cannot be combined
// with WithTLSConfig or WithInsecureSkipVerify.
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(o *clientOptions) {
		o.transport = transport
	}
}

// WithInsecureSkipVerify 显式关闭 HTTPS 服务器证书校验。
// 仅建议用于测试或明确使用自签名证书且有其他网络隔离措施的环境。
func WithInsecureSkipVerify() ClientOption {
	return func(o *clientOptions) {
		// This is intentionally opt-in; the secure default leaves verification enabled.
		o.tlsConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402 -- explicitly requested by the caller
	}
}

// WithHTTP 使用 HTTP 替代 HTTPS（不推荐用于生产环境）
func WithHTTP() ClientOption {
	return func(o *clientOptions) {
		o.useHTTP = true
	}
}

// WithFOSVersion 设置 Fabric OS 版本，用于兼容不同版本的 REST endpoint。
// Open records the version reported at login; New callers may use this option
// when the version is known out of band.
func WithFOSVersion(version string) ClientOption {
	return func(o *clientOptions) {
		o.fosVersion = strings.TrimSpace(version)
	}
}

// WithAllowUnknownFOSVersionWrites explicitly permits write requests when the
// FOS version is unknown. The safer default is to block writes until the
// version is supplied explicitly or learned during login.
func WithAllowUnknownFOSVersionWrites() ClientOption {
	return func(o *clientOptions) {
		o.allowUnknownFOSVersionWrites = true
	}
}

// Client is a concurrency-safe Brocade Fabric OS REST API client.
type Client struct {
	endpoint                     *url.URL
	switchAddress                string
	authToken                    string
	client                       *http.Client
	logger                       *slog.Logger
	vfID                         int
	logOutput                    io.Writer
	verbose                      bool
	customLogger                 bool
	timeout                      time.Duration
	retryCount                   int
	retryWaitTime                time.Duration
	retryMaxWaitTime             time.Duration
	maxResponseBodyBytes         int64
	fosVersion                   Version
	fosVersionConfigured         bool
	allowUnknownFOSVersionWrites bool
	waitForRetry                 func(context.Context, time.Duration) error
	stateMu                      sync.RWMutex
	loggerMu                     sync.RWMutex
}

// SessionInfo 是 POST /login 的 XML 响应，包含登录后的用户信息和交换机参数
type SessionInfo struct {
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

// New constructs a client without performing network I/O.
func New(endpoint string, opts ...ClientOption) (*Client, error) {
	options := clientOptions{
		timeout:              DefaultTimeout,
		retryCount:           DefaultRetryCount,
		retryWaitTime:        DefaultRetryWaitTime,
		retryMaxWaitTime:     DefaultRetryMaxWaitTime,
		maxResponseBodyBytes: DefaultMaxResponseBodyBytes,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	if options.timeout <= 0 {
		return nil, fmt.Errorf("%w: timeout must be positive", ErrInvalidConfig)
	}
	if options.retryCount < 0 {
		return nil, fmt.Errorf("%w: retry count must not be negative", ErrInvalidConfig)
	}
	if options.retryWaitTime < 0 {
		return nil, fmt.Errorf("%w: retry wait must not be negative", ErrInvalidConfig)
	}
	if options.retryMaxWaitTime < options.retryWaitTime {
		return nil, fmt.Errorf("%w: maximum retry wait is less than initial retry wait", ErrInvalidConfig)
	}
	if options.maxResponseBodyBytes <= 0 {
		return nil, fmt.Errorf("%w: response body limit must be positive", ErrInvalidConfig)
	}
	if options.transport != nil && options.tlsConfig != nil {
		return nil, fmt.Errorf("%w: custom transport and TLS config are mutually exclusive", ErrInvalidConfig)
	}

	parsedEndpoint, err := parseEndpoint(endpoint, options.useHTTP)
	if err != nil {
		return nil, err
	}

	var version Version
	if options.fosVersion != "" {
		version, err = ParseVersion(options.fosVersion)
		if err != nil {
			return nil, err
		}
	}

	transport := options.transport
	if transport == nil {
		transport = cloneHTTPTransport(options.tlsConfig)
	}
	logger := options.logger
	if logger == nil {
		logger = slog.Default()
	}
	c := &Client{
		endpoint:                     parsedEndpoint,
		switchAddress:                parsedEndpoint.Hostname(),
		client:                       &http.Client{Transport: transport, Timeout: options.timeout},
		logger:                       logger,
		logOutput:                    os.Stderr,
		customLogger:                 options.logger != nil,
		timeout:                      options.timeout,
		retryCount:                   options.retryCount,
		retryWaitTime:                options.retryWaitTime,
		retryMaxWaitTime:             options.retryMaxWaitTime,
		maxResponseBodyBytes:         options.maxResponseBodyBytes,
		fosVersion:                   version,
		fosVersionConfigured:         version.Valid(),
		allowUnknownFOSVersionWrites: options.allowUnknownFOSVersionWrites,
		waitForRetry:                 waitForRetry,
	}
	return c, nil
}

// Open constructs a client and authenticates it with short-lived credentials.
func Open(ctx context.Context, endpoint string, credentials Credentials, opts ...ClientOption) (*Client, error) {
	client, err := New(endpoint, opts...)
	if err != nil {
		return nil, err
	}
	if _, err := client.Login(ctx, credentials); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
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
// ctx 取消时会中止当前登录请求。
func (c *Client) Login(ctx context.Context, credentials Credentials) (*SessionInfo, error) {
	if strings.TrimSpace(credentials.Username) == "" || credentials.Password == "" {
		return nil, fmt.Errorf("%w: username and password are required", ErrInvalidConfig)
	}
	c.clearAuth()
	url := c.restBase() + c.endpoints().Login()
	headers := c.baseHeaders()
	headers.Set("Authorization", "Basic "+base64Encode(credentials.Username+":"+credentials.Password))

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
		c.setAuthenticated(authHeader, Version{})
		return &SessionInfo{}, nil
	}

	var result SessionInfo
	if err := xml.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("login response: %w", invalidResponseError(err))
	}
	var version Version
	if result.FirmwareVersion != "" {
		version, err = ParseVersion(result.FirmwareVersion)
		if err != nil {
			return nil, fmt.Errorf("login response firmware version: %w", err)
		}
	}
	c.setAuthenticated(authHeader, version)

	return &result, nil
}

// Logout 注销当前会话并清除认证 Token。
// ctx 取消时会中止当前注销请求。
func (c *Client) Logout(ctx context.Context) error {
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
	base := *c.endpoint
	relative, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}
	prefixPath := strings.TrimRight(base.Path, "/") + "/rest/running"
	prefixRawPath := strings.TrimRight(base.EscapedPath(), "/") + "/rest/running"
	base.Path = prefixPath + relative.Path
	if relative.RawPath != "" || base.RawPath != "" {
		base.RawPath = prefixRawPath + relative.EscapedPath()
	}
	base.RawQuery = relative.RawQuery
	c.stateMu.RLock()
	vfID := c.vfID
	c.stateMu.RUnlock()
	if vfID > 0 {
		query := base.Query()
		query.Set("vf-id", strconv.Itoa(vfID))
		base.RawQuery = query.Encode()
	}
	return base.String()
}

func (c *Client) restBase() string {
	base := *c.endpoint
	base.Path = strings.TrimRight(base.Path, "/") + "/rest"
	return base.String()
}

func parseEndpoint(raw string, useHTTP bool) (*url.URL, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, fmt.Errorf("%w: endpoint is required", ErrInvalidConfig)
	}
	if !strings.Contains(value, "://") {
		scheme := "https"
		if useHTTP {
			scheme = "http"
		}
		value = scheme + "://" + value
	}
	endpoint, err := url.Parse(value)
	if err != nil || endpoint.Host == "" {
		return nil, fmt.Errorf("%w: invalid endpoint %q", ErrInvalidConfig, raw)
	}
	if endpoint.Scheme != "https" && endpoint.Scheme != "http" {
		return nil, fmt.Errorf("%w: unsupported endpoint scheme %q", ErrInvalidConfig, endpoint.Scheme)
	}
	if useHTTP {
		endpoint.Scheme = "http"
	}
	if endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return nil, fmt.Errorf("%w: endpoint must not contain user info, query, or fragment", ErrInvalidConfig)
	}
	endpoint.Path = strings.TrimRight(endpoint.Path, "/")
	return endpoint, nil
}

// IsLoggedIn 返回当前是否已持有认证 Token。
func (c *Client) IsLoggedIn() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.authToken != ""
}

// Version returns the Fabric OS version learned during Login or configured at
// construction. The boolean is false when the version is unknown.
func (c *Client) Version() (Version, bool) {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.fosVersion, c.fosVersion.Valid()
}

// Capabilities returns the operations supported by the known Fabric OS
// version. Unknown versions remain read-only unless explicitly overridden.
func (c *Client) Capabilities() Capabilities {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return capabilitiesFor(c.fosVersion, c.allowUnknownFOSVersionWrites)
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
	wait = min(wait, c.retryMaxWaitTime)
	if wait <= 1 {
		return wait
	}
	// Full jitter prevents clients retrying a recovering switch in lockstep.
	return time.Duration(rand.Int64N(int64(wait) + 1))
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
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
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
			if attempt >= attempts || !isRetryableError(err) {
				break
			}
			if err := c.waitForRetry(ctx, c.retryAfter(nil, attempt)); err != nil {
				return nil, err
			}
			continue
		}

		responseBody, readErr := readResponseBody(resp.Body, c.maxResponseBodyBytes)
		_ = resp.Body.Close()
		if readErr != nil {
			if errors.Is(readErr, ErrResponseBodyTooLarge) {
				return nil, readErr
			}
			lastErr = readErr
			if attempt >= attempts {
				break
			}
			if err := c.waitForRetry(ctx, c.retryAfter(nil, attempt)); err != nil {
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
			if err := c.waitForRetry(ctx, c.retryAfter(apiResp, attempt)); err != nil {
				return nil, err
			}
			continue
		}
		return apiResp, nil
	}

	return nil, lastErr
}

func isRetryableStatus(status int) bool {
	switch status {
	case http.StatusRequestTimeout,
		http.StatusTooEarly,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func isRetryableError(err error) bool {
	var certificateError *tls.CertificateVerificationError
	if errors.As(err, &certificateError) {
		return false
	}
	var unknownAuthority x509.UnknownAuthorityError
	if errors.As(err, &unknownAuthority) {
		return false
	}
	var hostnameError x509.HostnameError
	if errors.As(err, &hostnameError) {
		return false
	}
	var netErr net.Error
	return errors.As(err, &netErr) && (netErr.Timeout() || !errors.Is(err, context.Canceled))
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

func (c *Client) setAuthenticated(authToken string, version Version) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.authToken = authToken
	if version.Valid() {
		c.fosVersion = version
	} else if !c.fosVersionConfigured {
		c.fosVersion = Version{}
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

func (c *Client) get(ctx context.Context, endpoint string, result any) error {
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

func (c *Client) post(ctx context.Context, endpoint string, payload any) error {
	return c.mutate(ctx, http.MethodPost, endpoint, payload, true)
}

func (c *Client) patch(ctx context.Context, endpoint string, payload any) error {
	return c.mutate(ctx, http.MethodPatch, endpoint, payload, true)
}

func (c *Client) delete(ctx context.Context, endpoint string) error {
	return c.mutate(ctx, http.MethodDelete, endpoint, nil, true)
}

func (c *Client) patchWithoutVersionGate(ctx context.Context, endpoint string, payload any) error {
	return c.mutate(ctx, http.MethodPatch, endpoint, payload, false)
}

func (c *Client) mutate(ctx context.Context, method, endpoint string, payload any, enforceVersionGate bool) error {
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
	endpointSet := c.endpoints()
	if endpointSet.allowWrite() {
		return nil
	}
	if !endpointSet.version.Valid() {
		return fmt.Errorf("%w: write operations require a known FOS version", ErrUnknownFOSVersion)
	}
	return fmt.Errorf("%w: FOS %s does not support write operations", ErrUnsupportedOperation, endpointSet.version)
}
