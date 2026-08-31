package sanswitch

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// newMockFOS 创建一个模拟 FOS REST API 的 HTTP 测试服务器。
// mux 中注册的路由将处理 /rest/... 路径。
func newMockFOS(t *testing.T, mux *http.ServeMux) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// newTestClient 创建一个指向测试服务器的 Client，跳过 Login。
func newTestClient(t *testing.T, ts *httptest.Server) *Client {
	t.Helper()
	return mustNewClient(t, ts.URL, WithFOSVersion("v9.2.0"))
}

func mustNewClient(t *testing.T, endpoint string, opts ...ClientOption) *Client {
	t.Helper()
	client, err := New(endpoint, opts...)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	return client
}

func pointClientAt(t *testing.T, client *Client, endpoint string) {
	t.Helper()
	parsed, err := parseEndpoint(endpoint, false)
	if err != nil {
		t.Fatalf("parse endpoint: %v", err)
	}
	client.endpoint = parsed
}

var testCredentials = Credentials{Username: "admin", Password: "password"}

// loginXML 是标准的 Login 响应 XML
const loginXML = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <switch-parameters>
    <user-name>admin</user-name>
    <chassis-access-role>admin</chassis-access-role>
    <home-virtual-fabric>1</home-virtual-fabric>
    <firmware-version>v9.2.0a</firmware-version>
    <model>G620</model>
  </switch-parameters>
</Response>`
