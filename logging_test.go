package san

import (
	"errors"
	"net/http"
	"testing"
)

func TestLoggingEndpointsRequireFOS90(t *testing.T) {
	var called bool
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "password", WithFOSVersion("v8.2.3"))
	c.baseURL = ts.URL + "/rest/running"

	if _, err := c.GetErrorLogs(); !errors.Is(err, ErrUnsupportedOperation) {
		t.Fatalf("expected ErrUnsupportedOperation for GetErrorLogs, got %v", err)
	}
	if _, err := c.GetAuditLogs(); !errors.Is(err, ErrUnsupportedOperation) {
		t.Fatalf("expected ErrUnsupportedOperation for GetAuditLogs, got %v", err)
	}
	if called {
		t.Fatal("expected unsupported logging endpoints to return before HTTP call")
	}
}

func TestLoggingEndpointsAllowFOS90(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-logging/error-log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><error-log><sequence-number>1</sequence-number></error-log></Response>`))
	})
	mux.HandleFunc("/rest/running/brocade-logging/audit-log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><audit-log><sequence-number>1</sequence-number></audit-log></Response>`))
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "password", WithFOSVersion("v9.0.0"))
	c.baseURL = ts.URL + "/rest/running"

	errorLogs, err := c.GetErrorLogs()
	if err != nil {
		t.Fatalf("GetErrorLogs() error: %v", err)
	}
	if len(errorLogs) != 1 {
		t.Fatalf("expected 1 error log, got %d", len(errorLogs))
	}
	auditLogs, err := c.GetAuditLogs()
	if err != nil {
		t.Fatalf("GetAuditLogs() error: %v", err)
	}
	if len(auditLogs) != 1 {
		t.Fatalf("expected 1 audit log, got %d", len(auditLogs))
	}
}
