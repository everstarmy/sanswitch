package san

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"testing"
)

const definedZoneXML = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <zone>
    <zone-name>zone_A</zone-name>
    <zone-type>0</zone-type>
    <zone-type-string>zone</zone-type-string>
    <member-entry>
      <entry-name>10:00:00:00:c9:f8:04:35</entry-name>
      <entry-name>20:00:00:00:c9:f8:04:35</entry-name>
    </member-entry>
  </zone>
  <zone>
    <zone-name>zone_B</zone-name>
    <zone-type>0</zone-type>
    <zone-type-string>zone</zone-type-string>
    <member-entry>
      <entry-name>10:00:00:00:c9:f8:04:36</entry-name>
      <entry-name>20:00:00:00:c9:f8:04:36</entry-name>
      <principal-entry-name>10:00:00:00:c9:f8:04:99</principal-entry-name>
    </member-entry>
  </zone>
</Response>`

const effectiveZoneXML = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <enabled-zone>
    <zone-name>zone_A</zone-name>
    <zone-type>0</zone-type>
    <zone-type-string>zone</zone-type-string>
    <member-entry>
      <entry-name>10:00:00:00:c9:f8:04:35</entry-name>
      <entry-name>20:00:00:00:c9:f8:04:35</entry-name>
    </member-entry>
  </enabled-zone>
</Response>`

const transactionTokenXML = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <effective-configuration>
    <transaction-token>12345</transaction-token>
  </effective-configuration>
</Response>`

func checksumXML(checksum string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <effective-configuration>
    <checksum>` + checksum + `</checksum>
  </effective-configuration>
</Response>`
}

func definedConfigXML(name string, zones ...string) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?><Response><cfg><cfg-name>`)
	b.WriteString(name)
	b.WriteString(`</cfg-name><member-zone>`)
	for _, zone := range zones {
		b.WriteString(`<zone-name>`)
		b.WriteString(zone)
		b.WriteString(`</zone-name>`)
	}
	b.WriteString(`</member-zone></cfg></Response>`)
	return b.String()
}

func TestGetDefinedZones(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(definedZoneXML))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	zones, err := c.GetDefinedZones()
	if err != nil {
		t.Fatalf("GetDefinedZones() error: %v", err)
	}
	if len(zones) != 2 {
		t.Fatalf("expected 2 zones, got %d", len(zones))
	}

	z := zones[0]
	if z.Name != "zone_A" {
		t.Errorf("unexpected zone name: %q", z.Name)
	}
	if len(z.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(z.Members))
	}

	z2 := zones[1]
	if z2.Name != "zone_B" {
		t.Errorf("unexpected zone name: %q", z2.Name)
	}
	// zone_B has 2 entry-name + 1 principal-entry-name = 3 members
	if len(z2.Members) != 3 {
		t.Errorf("expected 3 members for zone_B, got %d", len(z2.Members))
	}
}

func TestGetEffectiveZones(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/enabled-zone", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(effectiveZoneXML))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	zones, err := c.GetEffectiveZones()
	if err != nil {
		t.Fatalf("GetEffectiveZones() error: %v", err)
	}
	if len(zones) != 1 {
		t.Fatalf("expected 1 effective zone, got %d", len(zones))
	}
	if zones[0].Name != "zone_A" {
		t.Errorf("unexpected zone name: %q", zones[0].Name)
	}
}

func TestCreateZone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.CreateZone("new_zone", []string{"10:00:00:00:c9:f8:04:35", "20:00:00:00:c9:f8:04:35"}, []string{})
	if err != nil {
		t.Fatalf("CreateZone() error: %v", err)
	}
}

func TestCreateZoneAndActivateWorkflow(t *testing.T) {
	var calls []string
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/checksum", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/yang-data+xml")
		if len(calls) == 1 {
			w.Write([]byte(checksumXML("old-checksum")))
			return
		}
		w.Write([]byte(checksumXML("new-checksum")))
	})
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read zone body: %v", err)
		}
		var payload DefinedZoneAPI
		if err := xml.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal zone body: %v", err)
		}
		if payload.Name != "zone_new" || len(payload.MemberEntryNames) != 2 {
			t.Fatalf("unexpected zone payload: %+v", payload)
		}
		w.WriteHeader(http.StatusCreated)
	})
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/cfg", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/yang-data+xml")
			w.Write([]byte(definedConfigXML("cfg1", "zone_existing")))
		case http.MethodPatch:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read cfg body: %v", err)
			}
			var payload DefinedConfigAPI
			if err := xml.Unmarshal(body, &payload); err != nil {
				t.Fatalf("unmarshal cfg body: %v", err)
			}
			if payload.Name != "cfg1" {
				t.Fatalf("expected cfg1, got %q", payload.Name)
			}
			if got := strings.Join(payload.MemberZones, ","); got != "zone_existing,zone_new" {
				t.Fatalf("expected cfg zones zone_existing,zone_new; got %s", got)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-action-v2/save", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read save body: %v", err)
		}
		var payload PatchEffectiveConfigAPI
		if err := xml.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal save body: %v", err)
		}
		if payload.Checksum != "old-checksum" {
			t.Fatalf("expected old checksum, got %q", payload.Checksum)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-name/cfg1", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read activate body: %v", err)
		}
		var payload PatchEffectiveConfigAPI
		if err := xml.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal activate body: %v", err)
		}
		if payload.Checksum != "new-checksum" {
			t.Fatalf("expected new checksum, got %q", payload.Checksum)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.CreateZoneAndActivate("cfg1", "zone_new", []string{"10:00:00:00:00:00:00:01", "10:00:00:00:00:00:00:02"}, nil)
	if err != nil {
		t.Fatalf("CreateZoneAndActivate() error: %v", err)
	}

	want := []string{
		"GET /rest/running/brocade-zone/effective-configuration/checksum",
		"POST /rest/running/brocade-zone/defined-configuration/zone",
		"GET /rest/running/brocade-zone/defined-configuration/cfg",
		"PATCH /rest/running/brocade-zone/defined-configuration/cfg",
		"PATCH /rest/running/brocade-zone/effective-configuration/cfg-action-v2/save",
		"GET /rest/running/brocade-zone/effective-configuration/checksum",
		"PATCH /rest/running/brocade-zone/effective-configuration/cfg-name/cfg1",
	}
	if got := strings.Join(calls, "\n"); got != strings.Join(want, "\n") {
		t.Fatalf("unexpected call order:\n%s", got)
	}
}

func TestUpdateZone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.UpdateZone("zone_A", []string{"10:00:00:00:c9:f8:04:35", "new-member"}, []string{})
	if err != nil {
		t.Fatalf("UpdateZone() error: %v", err)
	}
}

func TestRenameZone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		wantPath := "/rest/running/brocade-zone/defined-configuration/zone/zone-name/old%2Fzone"
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Errorf("expected escaped path %q, got %q", wantPath, got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var payload DefinedZoneAPI
		if err := xml.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if payload.Name != "new zone" {
			t.Fatalf("expected new zone name, got %q", payload.Name)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.RenameZone("old/zone", "new zone")
	if err != nil {
		t.Fatalf("RenameZone() error: %v", err)
	}
}

func TestReplaceZoneAndActivateWorkflow(t *testing.T) {
	var calls []string
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/checksum", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/yang-data+xml")
		if len(calls) == 1 {
			w.Write([]byte(checksumXML("old-checksum")))
			return
		}
		w.Write([]byte(checksumXML("new-checksum")))
	})
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-action-v2/save", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-name/cfg1", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.ReplaceZoneAndActivate("cfg1", "zone_A", []string{"10:00:00:00:00:00:00:01"}, nil)
	if err != nil {
		t.Fatalf("ReplaceZoneAndActivate() error: %v", err)
	}

	want := []string{
		"GET /rest/running/brocade-zone/effective-configuration/checksum",
		"PATCH /rest/running/brocade-zone/defined-configuration/zone",
		"PATCH /rest/running/brocade-zone/effective-configuration/cfg-action-v2/save",
		"GET /rest/running/brocade-zone/effective-configuration/checksum",
		"PATCH /rest/running/brocade-zone/effective-configuration/cfg-name/cfg1",
	}
	if got := strings.Join(calls, "\n"); got != strings.Join(want, "\n") {
		t.Fatalf("unexpected call order:\n%s", got)
	}
}

func TestDeleteZone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_A", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.DeleteZone("zone_A")
	if err != nil {
		t.Fatalf("DeleteZone() error: %v", err)
	}
}

func TestDeleteZoneAndActivateWorkflow(t *testing.T) {
	var calls []string
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.EscapedPath())
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/rest/running/brocade-zone/effective-configuration/checksum":
			w.Header().Set("Content-Type", "application/yang-data+xml")
			if len(calls) == 1 {
				w.Write([]byte(checksumXML("old-checksum")))
				return
			}
			w.Write([]byte(checksumXML("new-checksum")))
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fdelete":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPatch && r.URL.Path == "/rest/running/brocade-zone/effective-configuration/cfg-action-v2/save":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPatch && r.URL.Path == "/rest/running/brocade-zone/effective-configuration/cfg-name/cfg1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.String())
			w.WriteHeader(http.StatusNotFound)
		}
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.DeleteZoneAndActivate("cfg1", "zone/delete")
	if err != nil {
		t.Fatalf("DeleteZoneAndActivate() error: %v", err)
	}

	want := []string{
		"GET /rest/running/brocade-zone/effective-configuration/checksum",
		"DELETE /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fdelete",
		"PATCH /rest/running/brocade-zone/effective-configuration/cfg-action-v2/save",
		"GET /rest/running/brocade-zone/effective-configuration/checksum",
		"PATCH /rest/running/brocade-zone/effective-configuration/cfg-name/cfg1",
	}
	if got := strings.Join(calls, "\n"); got != strings.Join(want, "\n") {
		t.Fatalf("unexpected call order:\n%s", got)
	}
}

func TestZoneAndActivateValidation(t *testing.T) {
	c := NewClient("localhost", "admin", "password")

	tests := []struct {
		name string
		err  error
		run  func() error
	}{
		{name: "create requires cfg", run: func() error { return c.CreateZoneAndActivate("", "zone1", []string{"member"}, nil) }},
		{name: "create requires zone", run: func() error { return c.CreateZoneAndActivate("cfg1", "", []string{"member"}, nil) }},
		{name: "create requires member", run: func() error { return c.CreateZoneAndActivate("cfg1", "zone1", nil, nil) }},
		{name: "replace requires member", run: func() error { return c.ReplaceZoneAndActivate("cfg1", "zone1", nil, nil) }},
		{name: "delete requires zone", run: func() error { return c.DeleteZoneAndActivate("cfg1", "") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

func TestDeleteZoneEscapesZoneName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		wantPath := "/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fwith%20space%3F"
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Errorf("expected escaped path %q, got %q", wantPath, got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.DeleteZone("zone/with space?")
	if err != nil {
		t.Fatalf("DeleteZone() error: %v", err)
	}
}

func TestAbortZoneTransaction(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-action-v2/transaction-abort", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.AbortZoneTransaction()
	if err != nil {
		t.Fatalf("AbortZoneTransaction() error: %v", err)
	}
}

func TestGetZoneTransactionStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/transaction-token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(transactionTokenXML))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	status, err := c.GetZoneTransactionStatus()
	if err != nil {
		t.Fatalf("GetZoneTransactionStatus() error: %v", err)
	}
	if status.TransactionToken != 12345 {
		t.Errorf("expected transaction-token 12345, got %d", status.TransactionToken)
	}
}

func TestDeleteAliasEscapesAliasName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		wantPath := "/rest/running/brocade-zone/defined-configuration/alias/alias-name/alias%2Fwith%20space%3F"
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Errorf("expected escaped path %q, got %q", wantPath, got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.DeleteAlias("alias/with space?")
	if err != nil {
		t.Fatalf("DeleteAlias() error: %v", err)
	}
}

func TestRenameAlias(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		wantPath := "/rest/running/brocade-zone/defined-configuration/alias/alias-name/old%2Falias"
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Errorf("expected escaped path %q, got %q", wantPath, got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var payload DefinedAliasAPI
		if err := xml.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if payload.Name != "new alias" {
			t.Fatalf("expected new alias name, got %q", payload.Name)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.RenameAlias("old/alias", "new alias")
	if err != nil {
		t.Fatalf("RenameAlias() error: %v", err)
	}
}

func TestActivateZoneConfigEscapesConfigName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		wantPath := "/rest/running/brocade-zone/effective-configuration/cfg-name/cfg%2Fwith%20space%3F"
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Errorf("expected escaped path %q, got %q", wantPath, got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.ActivateZoneConfig("cfg/with space?", "abc")
	if err != nil {
		t.Fatalf("ActivateZoneConfig() error: %v", err)
	}
}
