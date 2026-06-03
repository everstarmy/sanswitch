package san

import (
	"bytes"
	"context"
	"encoding/xml"
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
	c := NewClient("localhost", "admin", "password")
	c.baseURL = ts.URL + "/rest/running"

	resp, err := c.Login()
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

func TestLoginFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "wrong")
	c.baseURL = ts.URL + "/rest/running"

	_, err := c.Login()
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
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "password")
	c.baseURL = ts.URL + "/rest/running"

	if _, err := c.Login(); err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if err := c.Logout(); err != nil {
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
	err := c.Get("/test/items", &resp)
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
	err := c.Get("/test", &resp)
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

	err := c.Post("/test/create", payload)
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

	if err := c.Patch("/test/update", struct {
		XMLName xml.Name `xml:"item"`
		Name    string   `xml:"name"`
	}{Name: "updated"}); err != nil {
		t.Fatalf("Patch() error: %v", err)
	}

	if err := c.Delete("/test/delete/me"); err != nil {
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
	err := c.Get("/test/err", &resp)
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
	c := NewClient("switch.example", "admin", "password")
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
	if err := c.Get("/log-test", &quietResp); err != nil {
		t.Fatalf("Get() with verbose false error: %v", err)
	}
	if quiet.Len() != 0 {
		t.Fatalf("expected no debug log with verbose false, got %q", quiet.String())
	}

	var verbose bytes.Buffer
	c.SetLogOutput(&verbose)
	c.SetVerbose(true)
	var verboseResp testResponse
	if err := c.Get("/log-test", &verboseResp); err != nil {
		t.Fatalf("Get() with verbose true error: %v", err)
	}
	if !strings.Contains(verbose.String(), "GET response") {
		t.Fatalf("expected debug log with verbose true, got %q", verbose.String())
	}
}

func TestGetWithContextCancel(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "password",
		WithRetry(0), // 禁用重试以精确测试 context 取消
	)
	c.baseURL = ts.URL + "/rest/running"

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var resp struct{}
	err := c.GetWithContext(ctx, "/slow", &resp)
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
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
	c := NewClient("localhost", "admin", "password",
		WithRetry(3),
		WithRetryWait(50*time.Millisecond),
		WithRetryMaxWait(200*time.Millisecond),
	)
	c.baseURL = ts.URL + "/rest/running"

	var resp struct {
		OK string `xml:"ok"`
	}
	err := c.Get("/flaky", &resp)
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
	c := NewClient("localhost", "admin", "password",
		WithRetry(3),
		WithRetryWait(10*time.Millisecond),
		WithRetryMaxWait(20*time.Millisecond),
	)
	c.baseURL = ts.URL + "/rest/running"

	err := c.Post("/mutate", nil)
	if err == nil {
		t.Fatal("expected server error, got nil")
	}
	if attempts != 1 {
		t.Fatalf("expected POST to be attempted once, got %d attempts", attempts)
	}
}
