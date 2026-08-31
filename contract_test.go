package sanswitch

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFOS92ReadContracts(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rest/running/brocade-chassis/chassis", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		_, _ = w.Write([]byte("<Response><chassis><serial-number>contract</serial-number></chassis></Response>"))
	})
	mux.HandleFunc("GET /rest/running/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		_, _ = w.Write([]byte("<Response></Response>"))
	})
	server := newMockFOS(t, mux)
	client := newTestClient(t, server)

	tests := []struct {
		name string
		read func(context.Context) error
	}{
		{name: "logical switches", read: func(ctx context.Context) error { _, err := client.LogicalSwitches(ctx); return err }},
		{name: "defined aliases", read: func(ctx context.Context) error { _, err := client.DefinedAliases(ctx); return err }},
		{name: "effective config", read: func(ctx context.Context) error { _, err := client.EffectiveConfig(ctx); return err }},
		{name: "zone database", read: func(ctx context.Context) error { _, err := client.ZoneDatabase(ctx); return err }},
		{name: "statistics", read: func(ctx context.Context) error { _, err := client.FibreChannelStatistics(ctx); return err }},
		{name: "blades", read: func(ctx context.Context) error { _, err := client.Blades(ctx); return err }},
		{name: "fans", read: func(ctx context.Context) error { _, err := client.Fans(ctx); return err }},
		{name: "power supplies", read: func(ctx context.Context) error { _, err := client.PowerSupplies(ctx); return err }},
		{name: "hardware", read: func(ctx context.Context) error { _, err := client.Hardware(ctx); return err }},
		{name: "switch status", read: func(ctx context.Context) error { _, err := client.SwitchStatus(ctx); return err }},
		{name: "system resources", read: func(ctx context.Context) error { _, err := client.SystemResources(ctx); return err }},
		{name: "media RDPs", read: func(ctx context.Context) error { _, err := client.MediaRDPs(ctx); return err }},
		{name: "name server", read: func(ctx context.Context) error { _, err := client.FibreChannelNameServers(ctx); return err }},
		{name: "FDMI HBAs", read: func(ctx context.Context) error { _, err := client.FDMIHBAs(ctx); return err }},
		{name: "FDMI ports", read: func(ctx context.Context) error { _, err := client.FDMIPorts(ctx); return err }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if err := test.read(t.Context()); err != nil {
				t.Fatalf("read contract: %v", err)
			}
		})
	}
}

func TestAliasWriteContracts(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	for _, method := range []string{http.MethodPost, http.MethodPatch} {
		mux.HandleFunc(method+" /rest/running/brocade-zone/defined-configuration/alias", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read alias payload: %v", err)
			}
			payload := string(body)
			if !strings.Contains(payload, "<alias-name>host-a</alias-name>") ||
				!strings.Contains(payload, "<alias-entry-name>10:00:00:00:00:00:00:01</alias-entry-name>") {
				t.Fatalf("alias payload = %s", payload)
			}
			w.WriteHeader(http.StatusNoContent)
		})
	}
	server := newMockFOS(t, mux)
	client := newTestClient(t, server)
	members := []string{"10:00:00:00:00:00:00:01"}
	if err := client.CreateAlias(t.Context(), "host-a", members); err != nil {
		t.Fatalf("CreateAlias() error: %v", err)
	}
	if err := client.UpdateAlias(t.Context(), "host-a", members); err != nil {
		t.Fatalf("UpdateAlias() error: %v", err)
	}
}
