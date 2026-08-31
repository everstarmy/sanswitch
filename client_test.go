package sanswitch

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Authorization", "Bearer test-token-12345")
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(loginXML))
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)

	resp, err := c.Login(t.Context(), testCredentials)
	if err != nil {
		t.Fatalf("Login() error: %v", err)
	}
	if !c.IsLoggedIn() {
		t.Error("expected IsLoggedIn() = true after Login()")
	}
	if resp.UserName != "admin" {
		t.Errorf("expected username 'admin', got %q", resp.UserName)
	}
	if resp.FirmwareVersion != "v9.2.0a" {
		t.Errorf("expected firmware 'v9.2.0a', got %q", resp.FirmwareVersion)
	}
	if resp.Model != "G620" {
		t.Errorf("expected model 'G620', got %q", resp.Model)
	}
}

func TestLoginConfiguresFOS91SaveEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Authorization", "Bearer test-token-12345")
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(strings.Replace(loginXML, "v9.2.0a", "v9.1.1", 1)))
	})
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-action/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)

	if _, err := c.Login(t.Context(), testCredentials); err != nil {
		t.Fatalf("Login() error: %v", err)
	}
	if err := c.SaveZoneConfig(t.Context(), "abc"); err != nil {
		t.Fatalf("SaveZoneConfig() error: %v", err)
	}
}

func TestLoginWithoutBodyMarksLegacyVersionAndBlocksWrites(t *testing.T) {
	var postCalled bool
	var patchCalled bool
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Authorization", "Bearer test-token-12345")
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/rest/running/test/create", func(w http.ResponseWriter, r *http.Request) {
		postCalled = true
		w.WriteHeader(http.StatusCreated)
	})
	mux.HandleFunc("/rest/running/test/update", func(w http.ResponseWriter, r *http.Request) {
		patchCalled = true
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)

	resp, err := c.Login(t.Context(), testCredentials)
	if err != nil {
		t.Fatalf("Login() error: %v", err)
	}
	if resp.FirmwareVersion != "" {
		t.Fatalf("expected empty firmware version, got %q", resp.FirmwareVersion)
	}
	if err := c.post(t.Context(), "/test/create", nil); !errors.Is(err, ErrUnknownFOSVersion) {
		t.Fatalf("expected ErrUnknownFOSVersion for POST, got %v", err)
	}
	if err := c.patch(t.Context(), "/test/update", nil); !errors.Is(err, ErrUnknownFOSVersion) {
		t.Fatalf("expected ErrUnknownFOSVersion for PATCH, got %v", err)
	}
	if postCalled || patchCalled {
		t.Fatalf("expected write requests to be blocked before HTTP call; post=%v patch=%v", postCalled, patchCalled)
	}
}

func TestLoginNoContentMarksLegacyVersion(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Authorization", "Bearer test-token-12345")
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)

	resp, err := c.Login(t.Context(), testCredentials)
	if err != nil {
		t.Fatalf("Login() error: %v", err)
	}
	if resp.FirmwareVersion != "" {
		t.Fatalf("expected empty firmware version, got %q", resp.FirmwareVersion)
	}
	if err := c.patch(t.Context(), "/test/update", nil); !errors.Is(err, ErrUnknownFOSVersion) {
		t.Fatalf("expected ErrUnknownFOSVersion for PATCH, got %v", err)
	}
}

func TestFOSVersionBelow91BlocksWrites(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/test/create", func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("POST should not be sent for FOS < 9.1")
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost", WithFOSVersion("v8.2.3"))
	pointClientAt(t, c, ts.URL)

	if err := c.post(t.Context(), "/test/create", nil); !errors.Is(err, ErrUnsupportedOperation) {
		t.Fatalf("expected ErrUnsupportedOperation, got %v", err)
	}
}

func TestLoginFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)

	_, err := c.Login(t.Context(), testCredentials)
	if err == nil {
		t.Fatal("expected error for invalid credentials")
	}
}

func TestLogout(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Authorization", "Bearer tok")
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(loginXML))
	})
	mux.HandleFunc("/rest/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer tok" {
			t.Errorf("expected logout token, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)

	if _, err := c.Login(t.Context(), testCredentials); err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if err := c.Logout(t.Context()); err != nil {
		t.Fatalf("Logout() error: %v", err)
	}
	if c.IsLoggedIn() {
		t.Error("expected IsLoggedIn() = false after Logout()")
	}
}

func TestGetXMLParsing(t *testing.T) {
	type testItem struct {
		XMLName xml.Name `xml:"item"`
		Name    string   `xml:"name"`
		Value   int      `xml:"value"`
	}
	type testResponse struct {
		XMLName xml.Name   `xml:"Response"`
		Items   []testItem `xml:"item"`
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/test/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
  <item><name>alpha</name><value>1</value></item>
  <item><name>beta</name><value>2</value></item>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	var resp testResponse
	err := c.get(t.Context(), "/test/items", &resp)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}
	if resp.Items[0].Name != "alpha" || resp.Items[0].Value != 1 {
		t.Errorf("unexpected first item: %+v", resp.Items[0])
	}
	if resp.Items[1].Name != "beta" || resp.Items[1].Value != 2 {
		t.Errorf("unexpected second item: %+v", resp.Items[1])
	}
}

func TestGetUnauthorized(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	var resp struct{}
	err := c.get(t.Context(), "/test", &resp)
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestPostXML(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/test/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/yang-data+xml" {
			t.Errorf("expected Content-Type application/yang-data+xml, got %s", ct)
		}
		w.WriteHeader(http.StatusCreated)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	payload := struct {
		XMLName xml.Name `xml:"item"`
		Name    string   `xml:"name"`
	}{Name: "test-item"}

	err := c.post(t.Context(), "/test/create", payload)
	if err != nil {
		t.Fatalf("Post() error: %v", err)
	}
}

func TestPatchAndDelete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/test/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/rest/running/test/delete/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	if err := c.patch(t.Context(), "/test/update", struct {
		XMLName xml.Name `xml:"item"`
		Name    string   `xml:"name"`
	}{Name: "updated"}); err != nil {
		t.Fatalf("Patch() error: %v", err)
	}

	if err := c.delete(t.Context(), "/test/delete/me"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}

func TestAPIErrorParsing(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/test/err", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<errors>
  <error>
    <error-code>invalid-input</error-code>
    <error-message>Zone name is required</error-message>
  </error>
</errors>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	var resp struct{}
	err := c.get(t.Context(), "/test/err", &resp)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
	if apiErr.ErrorCode != "invalid-input" {
		t.Errorf("expected error code 'invalid-input', got %q", apiErr.ErrorCode)
	}
	if apiErr.Message != "Zone name is required" {
		t.Errorf("expected message 'Zone name is required', got %q", apiErr.Message)
	}
}

func TestBuildURLAddsVFIDToExistingQuery(t *testing.T) {
	c := mustNewClient(t, "switch.example")
	c.SetVFID(128)

	got := c.buildURL("/brocade-test/resource?depth=2")
	want := "https://switch.example/rest/running/brocade-test/resource?depth=2&vf-id=128"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSetVerboseTogglesDebugLogging(t *testing.T) {
	type testResponse struct {
		OK string `xml:"ok"`
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/log-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><ok>true</ok></Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	var quiet bytes.Buffer
	c.logger = slog.New(slog.NewTextHandler(&quiet, &slog.HandlerOptions{Level: slog.LevelDebug}))
	c.SetVerbose(false)
	var quietResp testResponse
	if err := c.get(t.Context(), "/log-test", &quietResp); err != nil {
		t.Fatalf("Get() with verbose false error: %v", err)
	}
	if quiet.Len() != 0 {
		t.Fatalf("expected no debug log with verbose false, got %q", quiet.String())
	}

	var verbose bytes.Buffer
	c.SetLogOutput(&verbose)
	c.SetVerbose(true)
	var verboseResp testResponse
	if err := c.get(t.Context(), "/log-test", &verboseResp); err != nil {
		t.Fatalf("Get() with verbose true error: %v", err)
	}
	if !strings.Contains(verbose.String(), "GET response") {
		t.Fatalf("expected debug log with verbose true, got %q", verbose.String())
	}
}

func TestContextCancellation(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/slow", func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost",
		WithRetry(0), // 禁用重试以精确测试 context 取消
	)
	pointClientAt(t, c, ts.URL)

	ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
	defer cancel()

	var resp struct{}
	err := c.get(ctx, "/slow", &resp)
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
}

func TestResponseBodyLimit(t *testing.T) {
	attempts := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/large", func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(strings.Repeat("x", 64)))
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost",
		WithRetry(3),
		WithMaxResponseBodyBytes(32),
	)
	pointClientAt(t, c, ts.URL)

	var resp struct{}
	err := c.get(t.Context(), "/large", &resp)
	if !errors.Is(err, ErrResponseBodyTooLarge) {
		t.Fatalf("Get() error = %v; want ErrResponseBodyTooLarge", err)
	}
	if attempts != 1 {
		t.Fatalf("large response was retried %d times; want 1 attempt", attempts)
	}
}

func TestClientCloseIsIdempotent(t *testing.T) {
	c := mustNewClient(t, "switch.example")
	if err := c.Close(); err != nil {
		t.Fatalf("first Close() error: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("second Close() error: %v", err)
	}
}

func TestRetryOnServerError(t *testing.T) {
	attempts := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/flaky", func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><ok>true</ok></Response>`))
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost",
		WithRetry(3),
		WithRetryWait(50*time.Millisecond),
		WithRetryMaxWait(200*time.Millisecond),
	)
	pointClientAt(t, c, ts.URL)

	var resp struct {
		OK string `xml:"ok"`
	}
	err := c.get(t.Context(), "/flaky", &resp)
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if resp.OK != "true" {
		t.Errorf("expected ok=true, got %q", resp.OK)
	}
	if attempts < 3 {
		t.Errorf("expected at least 3 attempts, got %d", attempts)
	}
}

func TestPostDoesNotRetryOnServerError(t *testing.T) {
	attempts := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/mutate", func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusInternalServerError)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost",
		WithRetry(3),
		WithRetryWait(10*time.Millisecond),
		WithRetryMaxWait(20*time.Millisecond),
		WithFOSVersion("v9.2.0"),
	)
	pointClientAt(t, c, ts.URL)

	err := c.post(t.Context(), "/mutate", nil)
	if err == nil {
		t.Fatal("expected server error, got nil")
	}
	if attempts != 1 {
		t.Fatalf("expected POST to be attempted once, got %d attempts", attempts)
	}
}

func TestWriteRequiresKnownFOSVersion(t *testing.T) {
	called := false
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/test/create", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)

	err := c.post(t.Context(), "/test/create", nil)
	if !errors.Is(err, ErrUnknownFOSVersion) {
		t.Fatalf("expected ErrUnknownFOSVersion, got %v", err)
	}
	if called {
		t.Fatal("expected unknown FOS version to block the write before HTTP")
	}
}

func TestUnknownFOSVersionWritesCanBeExplicitlyEnabled(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/test/create", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost", WithAllowUnknownFOSVersionWrites())
	pointClientAt(t, c, ts.URL)

	if err := c.post(t.Context(), "/test/create", nil); err != nil {
		t.Fatalf("expected explicitly enabled write to succeed, got %v", err)
	}
}

func TestLoginRequiresAuthorizationHeader(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(loginXML))
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)

	_, err := c.Login(t.Context(), testCredentials)
	if !errors.Is(err, ErrInvalidResponse) {
		t.Fatalf("expected ErrInvalidResponse, got %v", err)
	}
	if c.IsLoggedIn() {
		t.Fatal("client must remain logged out when the response has no token")
	}
}

func TestFailedReloginClearsPreviousToken(t *testing.T) {
	attempts := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Authorization", "Bearer first-token")
			w.Header().Set("Content-Type", "application/yang-data+xml")
			w.Write([]byte(loginXML))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)
	if _, err := c.Login(t.Context(), testCredentials); err != nil {
		t.Fatalf("first Login() error: %v", err)
	}
	if !c.IsLoggedIn() {
		t.Fatal("expected client to be logged in after first login")
	}

	if _, err := c.Login(t.Context(), testCredentials); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if c.IsLoggedIn() {
		t.Fatal("failed re-login must not leave the old token installed")
	}
}

func TestUnauthorizedResponseClearsToken(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Authorization", "Bearer token")
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(loginXML))
	})
	mux.HandleFunc("/rest/running/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	ts := newMockFOS(t, mux)
	c := mustNewClient(t, "localhost")
	pointClientAt(t, c, ts.URL)
	if _, err := c.Login(t.Context(), testCredentials); err != nil {
		t.Fatalf("Login() error: %v", err)
	}

	var result struct{}
	if err := c.get(t.Context(), "/test", &result); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if c.IsLoggedIn() {
		t.Fatal("401 response must clear the local token")
	}
}

func TestDebugLogsDoNotContainResponseBody(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/log-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><secret>do-not-log-this</secret></Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)
	var logs bytes.Buffer
	c.SetLogOutput(&logs)
	c.SetVerbose(true)

	var result struct {
		Secret string `xml:"secret"`
	}
	if err := c.get(t.Context(), "/log-test", &result); err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if strings.Contains(logs.String(), "do-not-log-this") {
		t.Fatalf("debug log leaked response body: %q", logs.String())
	}
	if !strings.Contains(logs.String(), "GET response") {
		t.Fatalf("expected request metadata in debug log, got %q", logs.String())
	}
}

func TestClientTLSConfiguration(t *testing.T) {
	secure := mustNewClient(t, "switch.example")
	secureTransport, ok := secure.client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", secure.client.Transport)
	}
	if secureTransport.TLSClientConfig != nil && secureTransport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("TLS certificate verification must be enabled by default")
	}

	insecure := mustNewClient(t, "switch.example", WithInsecureSkipVerify())
	insecureTransport, ok := insecure.client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", insecure.client.Transport)
	}
	if insecureTransport.TLSClientConfig == nil || !insecureTransport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("WithInsecureSkipVerify() must explicitly disable certificate verification")
	}
}

func TestParseRetryAfter(t *testing.T) {
	wait, ok := parseRetryAfter("2")
	if !ok || wait != 2*time.Second {
		t.Fatalf("parseRetryAfter(2) = %v, %v; want 2s, true", wait, ok)
	}
	if _, ok := parseRetryAfter("not-a-duration"); ok {
		t.Fatal("expected invalid Retry-After to be rejected")
	}
}
