package sanswitch

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

type retryableNetError struct{}

func (retryableNetError) Error() string   { return "temporary network failure" }
func (retryableNetError) Timeout() bool   { return false }
func (retryableNetError) Temporary() bool { return true }

func TestNewValidatesConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		options  []ClientOption
	}{
		{name: "empty endpoint"},
		{name: "unsupported scheme", endpoint: "ftp://switch.example"},
		{name: "non-positive timeout", endpoint: "switch.example", options: []ClientOption{WithTimeout(0)}},
		{name: "negative retry count", endpoint: "switch.example", options: []ClientOption{WithRetry(-1)}},
		{name: "negative retry wait", endpoint: "switch.example", options: []ClientOption{WithRetryWait(-time.Second)}},
		{name: "inverted retry waits", endpoint: "switch.example", options: []ClientOption{WithRetryWait(time.Second), WithRetryMaxWait(time.Millisecond)}},
		{name: "non-positive body limit", endpoint: "switch.example", options: []ClientOption{WithMaxResponseBodyBytes(0)}},
		{name: "invalid FOS version", endpoint: "switch.example", options: []ClientOption{WithFOSVersion("nine")}},
		{name: "transport and TLS", endpoint: "switch.example", options: []ClientOption{WithTransport(roundTripFunc(nil)), WithTLSConfig(&tls.Config{})}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := New(test.endpoint, test.options...)
			if err == nil {
				t.Fatal("New() returned nil error")
			}
			if test.name == "invalid FOS version" {
				if !errors.Is(err, ErrInvalidVersion) {
					t.Fatalf("New() error = %v, want ErrInvalidVersion", err)
				}
				return
			}
			if !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("New() error = %v, want ErrInvalidConfig", err)
			}
		})
	}
}

func TestOpenUsesShortLivedCredentialsAndLearnsCapabilities(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /rest/login", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "admin" || password != "secret" {
			t.Fatalf("BasicAuth() = %q, %q, %v", username, password, ok)
		}
		w.Header().Set("Authorization", "Bearer token")
		_, _ = w.Write([]byte(loginXML))
	})
	server := newMockFOS(t, mux)

	client, err := Open(t.Context(), server.URL, Credentials{Username: "admin", Password: "secret"})
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}
	version, known := client.Version()
	if !known || version != (Version{Major: 9, Minor: 2}) {
		t.Fatalf("Version() = %v, %v", version, known)
	}
	capabilities := client.Capabilities()
	if !capabilities.ZoneWrite || !capabilities.FirmwareHistory || !capabilities.ModernZoneEndpoint {
		t.Fatalf("Capabilities() = %+v", capabilities)
	}
}

func TestClientSessionStateIsRaceSafe(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /rest/login", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Authorization", "Bearer token")
		_, _ = w.Write([]byte(loginXML))
	})
	mux.HandleFunc("POST /rest/logout", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	server := newMockFOS(t, mux)
	client := mustNewClient(t, server.URL)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var failures []error
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			if i%2 == 0 {
				_, err = client.Login(t.Context(), testCredentials)
			} else {
				err = client.Logout(t.Context())
			}
			if err != nil {
				mu.Lock()
				failures = append(failures, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if len(failures) != 0 {
		t.Fatalf("concurrent session failures: %v", failures)
	}
}

func TestRetryStatusClassification(t *testing.T) {
	retryable := []int{408, 425, 429, 500, 502, 503, 504}
	for _, status := range retryable {
		if !isRetryableStatus(status) {
			t.Errorf("status %d should be retryable", status)
		}
	}
	for _, status := range []int{400, 401, 404, 501, 505} {
		if isRetryableStatus(status) {
			t.Errorf("status %d should not be retryable", status)
		}
	}
}

func TestRetryRequestClassification(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		if !isRetryableMethod(method) {
			t.Errorf("method %s should be retryable", method)
		}
	}
	for _, method := range []string{http.MethodPost, http.MethodPatch, http.MethodDelete} {
		if isRetryableMethod(method) {
			t.Errorf("method %s should not be retryable", method)
		}
	}

	if !isRetryableError(retryableNetError{}) {
		t.Error("network failure should be retryable")
	}
	if isRetryableError(errors.New("permanent failure")) {
		t.Error("plain error should not be retryable")
	}
	if isRetryableError(context.Canceled) {
		t.Error("canceled context should not be retryable")
	}
	certificateError := &tls.CertificateVerificationError{Err: x509.UnknownAuthorityError{}}
	if isRetryableError(certificateError) {
		t.Error("certificate verification error should not be retryable")
	}
}

func TestNilContextIsRejected(t *testing.T) {
	client := mustNewClient(t, "switch.example")
	_, err := client.Login(nil, testCredentials)
	if err == nil {
		t.Fatal("Login(nil) returned nil error")
	}
	if errors.Is(err, context.Canceled) {
		t.Fatalf("Login(nil) error = %v", err)
	}
}
