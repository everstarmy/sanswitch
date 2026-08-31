package sanswitch

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestZoneTransactionCommit(t *testing.T) {
	var calls []string
	checksumReads := 0
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rest/running/brocade-zone/effective-configuration/checksum", func(w http.ResponseWriter, _ *http.Request) {
		checksumReads++
		checksum := "before"
		if checksumReads > 1 {
			checksum = "after"
		}
		_, _ = io.WriteString(w, checksumXML(checksum))
	})
	mux.HandleFunc("GET /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone-new", func(w http.ResponseWriter, _ *http.Request) {
		writeNotFound(w)
	})
	mux.HandleFunc("POST /rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, _ *http.Request) {
		calls = append(calls, "create")
		w.WriteHeader(http.StatusCreated)
	})
	mux.HandleFunc("GET /rest/running/brocade-zone/defined-configuration/cfg", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, definedConfigXML("cfg1", "zone-old"))
	})
	mux.HandleFunc("PATCH /rest/running/brocade-zone/defined-configuration/cfg", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read config body: %v", err)
		}
		if !strings.Contains(string(body), "zone-new") {
			t.Fatalf("config update omitted zone-new: %s", body)
		}
		calls = append(calls, "add-to-config")
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("PATCH /rest/running/brocade-zone/effective-configuration/cfg-action-v2/save", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "<checksum>before</checksum>" {
			t.Fatalf("save checksum = %s, want before", body)
		}
		calls = append(calls, "save")
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("PATCH /rest/running/brocade-zone/effective-configuration/cfg-name/cfg1", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "<checksum>after</checksum>" {
			t.Fatalf("activate checksum = %s, want after", body)
		}
		calls = append(calls, "activate")
		w.WriteHeader(http.StatusNoContent)
	})

	server := newMockFOS(t, mux)
	client := newTestClient(t, server)
	tx, err := client.BeginZoneTransaction(t.Context())
	if err != nil {
		t.Fatalf("BeginZoneTransaction() error: %v", err)
	}
	if err := tx.CreateZone(t.Context(), "zone-new", []string{"10:00:00:00:00:00:00:01"}, nil); err != nil {
		t.Fatalf("CreateZone() error: %v", err)
	}
	if err := tx.AddZoneToConfig(t.Context(), "cfg1", "zone-new"); err != nil {
		t.Fatalf("AddZoneToConfig() error: %v", err)
	}
	if err := tx.Commit(t.Context(), "cfg1"); err != nil {
		t.Fatalf("Commit() error: %v", err)
	}

	want := []string{"create", "add-to-config", "save", "activate"}
	if strings.Join(calls, ",") != strings.Join(want, ",") {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
	if err := tx.Abort(t.Context()); !errors.Is(err, ErrTransactionClosed) {
		t.Fatalf("Abort() after commit = %v, want ErrTransactionClosed", err)
	}
}

func TestZoneTransactionPartialFailureCanAbort(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rest/running/brocade-zone/effective-configuration/checksum", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, checksumXML("before"))
	})
	mux.HandleFunc("GET /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone-new", func(w http.ResponseWriter, _ *http.Request) {
		writeNotFound(w)
	})
	mux.HandleFunc("POST /rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	mux.HandleFunc("GET /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone-missing", func(w http.ResponseWriter, _ *http.Request) {
		writeNotFound(w)
	})
	mux.HandleFunc("PATCH /rest/running/brocade-zone/effective-configuration/cfg-action-v2/transaction-abort", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	server := newMockFOS(t, mux)
	client := newTestClient(t, server)
	tx, err := client.BeginZoneTransaction(t.Context())
	if err != nil {
		t.Fatalf("BeginZoneTransaction() error: %v", err)
	}
	if err := tx.CreateZone(t.Context(), "zone-new", []string{"member"}, nil); err != nil {
		t.Fatalf("CreateZone() error: %v", err)
	}
	err = tx.ReplaceZone(t.Context(), "zone-missing", []string{"member"}, nil)
	var partial *PartialMutationError
	if !errors.As(err, &partial) {
		t.Fatalf("ReplaceZone() error = %v, want PartialMutationError", err)
	}
	if err := tx.Abort(t.Context()); err != nil {
		t.Fatalf("Abort() error: %v", err)
	}
}

func TestZoneTransactionRequiresMutationBeforeCommit(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rest/running/brocade-zone/effective-configuration/checksum", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, checksumXML("before"))
	})
	server := newMockFOS(t, mux)
	tx, err := newTestClient(t, server).BeginZoneTransaction(t.Context())
	if err != nil {
		t.Fatalf("BeginZoneTransaction() error: %v", err)
	}
	if err := tx.Commit(t.Context(), "cfg1"); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("Commit() error = %v, want ErrInvalidConfig", err)
	}
}

func TestZoneTransactionRemoveAndDelete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rest/running/brocade-zone/effective-configuration/checksum", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, checksumXML("before"))
	})
	mux.HandleFunc("GET /rest/running/brocade-zone/defined-configuration/cfg", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, definedConfigXML("cfg1", "zone_A"))
	})
	mux.HandleFunc("PATCH /rest/running/brocade-zone/defined-configuration/cfg", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "zone_A") {
			t.Fatalf("removed zone remained in config: %s", body)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_A", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, definedZoneXML)
	})
	mux.HandleFunc("DELETE /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_A", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("PATCH /rest/running/brocade-zone/effective-configuration/cfg-action-v2/transaction-abort", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	server := newMockFOS(t, mux)
	tx, err := newTestClient(t, server).BeginZoneTransaction(t.Context())
	if err != nil {
		t.Fatalf("BeginZoneTransaction() error: %v", err)
	}
	if err := tx.RemoveZoneFromConfig(t.Context(), "cfg1", "zone_A"); err != nil {
		t.Fatalf("RemoveZoneFromConfig() error: %v", err)
	}
	if err := tx.DeleteZone(t.Context(), "zone_A"); err != nil {
		t.Fatalf("DeleteZone() error: %v", err)
	}
	if err := tx.Abort(t.Context()); err != nil {
		t.Fatalf("Abort() error: %v", err)
	}
}
