package san

import (
	"errors"
	"net/http"
	"testing"
)

func TestFRUHistoryLogAndSensorRequireFOS90(t *testing.T) {
	var called bool
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNotFound)
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "password", WithFOSVersion("v8.2.3"))
	c.baseURL = ts.URL + "/rest/running"

	if _, err := c.GetHistoryLogs(); !errors.Is(err, ErrUnsupportedOperation) {
		t.Fatalf("expected ErrUnsupportedOperation for GetHistoryLogs, got %v", err)
	}
	if _, err := c.GetSensors(); !errors.Is(err, ErrUnsupportedOperation) {
		t.Fatalf("expected ErrUnsupportedOperation for GetSensors, got %v", err)
	}
	if called {
		t.Fatal("expected unsupported FRU endpoints to return before HTTP call")
	}
}

func TestFRUHistoryLogAndSensorAllowFOS90(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-fru/history-log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><history-log><fru-type>fan</fru-type></history-log></Response>`))
	})
	mux.HandleFunc("/rest/running/brocade-fru/sensor", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><sensor><id>1</id></sensor></Response>`))
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "password", WithFOSVersion("v9.0.0"))
	c.baseURL = ts.URL + "/rest/running"

	historyLogs, err := c.GetHistoryLogs()
	if err != nil {
		t.Fatalf("GetHistoryLogs() error: %v", err)
	}
	if len(historyLogs) != 1 {
		t.Fatalf("expected 1 history log, got %d", len(historyLogs))
	}
	sensors, err := c.GetSensors()
	if err != nil {
		t.Fatalf("GetSensors() error: %v", err)
	}
	if len(sensors) != 1 {
		t.Fatalf("expected 1 sensor, got %d", len(sensors))
	}
}
