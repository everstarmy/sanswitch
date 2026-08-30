package san

import (
	"errors"
	"fmt"
)

// 预定义错误，用于在 HTTP 响应和内部逻辑中统一判断
var (
	// ErrNotFound 请求的资源不存在（如交换机信息、Zone 配置为空）
	ErrNotFound = errors.New("resource not found")
	// ErrUnauthorized HTTP 401 认证失败，Token 无效或已过期
	ErrUnauthorized = errors.New("unauthorized access")
	// ErrConnectionFailed 无法与交换机建立 TCP 连接
	ErrConnectionFailed = errors.New("connection failed")
	// ErrInvalidResponse 服务端返回的响应格式异常，无法解析
	ErrInvalidResponse = errors.New("invalid response from server")
	// ErrTimeout 请求超时
	ErrTimeout = errors.New("request timeout")
	// ErrUnsupportedOperation 当前 FOS 版本不支持该操作
	ErrUnsupportedOperation = errors.New("operation unsupported by FOS version")
	// ErrResponseBodyTooLarge 服务端响应超过客户端允许的大小上限
	ErrResponseBodyTooLarge = errors.New("response body too large")
)

// APIError 表示 FOS REST API 返回的结构化错误
type APIError struct {
	XMLName    struct{} `xml:"errors" json:"-"`
	StatusCode int      `xml:"-" json:"status_code"`
	Message    string   `xml:"error>error-message" json:"message"`
	ErrorCode  string   `xml:"error>error-code" json:"error_code"`
}

// PartialMutationError indicates that a multi-step Zone operation failed after
// at least one remote mutation had been attempted. The underlying error is
// preserved for errors.Is and errors.As.
type PartialMutationError struct {
	Err error
}

func (e *PartialMutationError) Error() string {
	return fmt.Sprintf("zone operation may have left a remote mutation: %v", e.Err)
}

func (e *PartialMutationError) Unwrap() error {
	return e.Err
}

func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("fos api error (status %d, code %s): %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("fos api error (status %d): %s", e.StatusCode, e.Message)
}

// Unwrap maps common HTTP statuses to package-level sentinel errors.
func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case 404:
		return ErrNotFound
	case 401:
		return ErrUnauthorized
	default:
		return nil
	}
}
