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
  <effective-configuration>
	<enabled-zone>
		<zone-name>zone_A</zone-name>
		<zone-type>0</zone-type>
		<zone-type-string>zone</zone-type-string>
		<member-entry>
		<entry-name>10:00:00:00:c9:f8:04:35</entry-name>
		<entry-name>20:00:00:00:c9:f8:04:35</entry-name>
		</member-entry>
	</enabled-zone>
  </effective-configuration>
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

func writeNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><errors><error><error-message>not found</error-message></error></errors>`))
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
	if len(z.Members.MemberEntries) != 2 {
		t.Errorf("expected 2 members, got %d", len(z.Members.MemberEntries))
	}

	z2 := zones[1]
	if z2.Name != "zone_B" {
		t.Errorf("unexpected zone name: %q", z2.Name)
	}
	// zone_B has 2 entry-name + 1 principal-entry-name = 3 members
	if len(z2.Members.MemberEntries)+len(z2.Members.PrincipalEntries) != 3 {
		t.Errorf("expected 3 members for zone_B, got %d", len(z2.Members.MemberEntries)+len(z2.Members.PrincipalEntries))
	}
}

func TestGetDefinedZone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/qos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
  <zone>
    <zone-name>qos</zone-name>
    <zone-type>1</zone-type>
    <zone-type-string>user-created-peer-zone</zone-type-string>
    <member-entry>
      <principal-entry-name>10:10:10:27:f8:8f:44:cd</principal-entry-name>
      <entry-name>10:10:10:27:f8:f0:2a:e8</entry-name>
      <entry-name>10:10:10:27:f8:f0:3a:70</entry-name>
      <entry-name>10:10:10:27:f8:f0:38:70</entry-name>
    </member-entry>
  </zone>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	zone, err := c.GetDefinedZone("qos")
	if err != nil {
		t.Fatalf("GetDefinedZone() error: %v", err)
	}
	if zone.Name != "qos" || zone.Type != "1" || zone.TypeString != ZoneTypeUserCreatedPeerZone {
		t.Fatalf("unexpected zone: %+v", zone)
	}
	if got := strings.Join(zone.Members.PrincipalEntries, ","); got != "10:10:10:27:f8:8f:44:cd" {
		t.Fatalf("unexpected principal entries: %s", got)
	}
	if got := strings.Join(zone.Members.MemberEntries, ","); got != "10:10:10:27:f8:f0:2a:e8,10:10:10:27:f8:f0:3a:70,10:10:10:27:f8:f0:38:70" {
		t.Fatalf("unexpected member entries: %s", got)
	}
}

func TestGetDefinedZoneEscapesZoneName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wantPath := "/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fwith%20space%3F"
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Fatalf("expected escaped path %q, got %q", wantPath, got)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><zone><zone-name>zone/with space?</zone-name><zone-type-string>zone</zone-type-string></zone></Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	if _, err := c.GetDefinedZone("zone/with space?"); err != nil {
		t.Fatalf("GetDefinedZone() error: %v", err)
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
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/new_zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeNotFound(w)
	})
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

func TestCreateZoneRejectsExistingZone(t *testing.T) {
	var postCalled bool
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/existing_zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response><zone><zone-name>existing_zone</zone-name><zone-type-string>zone</zone-type-string></zone></Response>`))
	})
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		postCalled = true
		w.WriteHeader(http.StatusCreated)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	if err := c.CreateZone("existing_zone", []string{"member"}, nil); err == nil {
		t.Fatal("expected CreateZone to reject existing zone")
	}
	if postCalled {
		t.Fatal("expected CreateZone to return before POST")
	}
}

func TestCreateZoneWithPrincipalMembersCreatesPeerZone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/peer_zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeNotFound(w)
	})
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var payload DefinedZoneAPI
		if err := xml.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if strings.Contains(string(body), "<zone-type>") {
			t.Fatalf("create payload should not include zone-type: %s", string(body))
		}
		if payload.ZoneTypeString != ZoneTypeUserCreatedPeerZone {
			t.Fatalf("expected zone type %q, got %q", ZoneTypeUserCreatedPeerZone, payload.ZoneTypeString)
		}
		if got := strings.Join(payload.PrincipalEntryNames, ","); got != "10:10:10:27:f8:8f:44:cd" {
			t.Fatalf("unexpected principal entries: %s", got)
		}
		if got := strings.Join(payload.MemberEntryNames, ","); got != "10:10:10:27:f8:f0:2a:e8,10:10:10:27:f8:f0:3a:70,10:10:10:27:f8:f0:38:65" {
			t.Fatalf("unexpected member entries: %s", got)
		}
		w.WriteHeader(http.StatusCreated)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.CreateZone(
		"peer_zone",
		[]string{"10:10:10:27:f8:f0:2a:e8", "10:10:10:27:f8:f0:3a:70", "10:10:10:27:f8:f0:38:65"},
		[]string{"10:10:10:27:f8:8f:44:cd"},
	)
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
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_new", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeNotFound(w)
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
		if got := string(body); got != "<checksum>old-checksum</checksum>" {
			t.Fatalf("unexpected save body: %s", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-name/cfg1", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read activate body: %v", err)
		}
		if got := string(body); got != "<checksum>new-checksum</checksum>" {
			t.Fatalf("unexpected activate body: %s", got)
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
		"GET /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_new",
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
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_A", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <zone>
    <zone-name>zone_A</zone-name>
    <zone-type-string>zone</zone-type-string>
  </zone>
</Response>`))
	})
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var payload DefinedZoneAPI
		if err := xml.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if strings.Contains(string(body), "<zone-type>") {
			t.Fatalf("update payload should not include zone-type: %s", string(body))
		}
		if payload.ZoneTypeString != ZoneTypeZone {
			t.Fatalf("expected zone type %q, got %q", ZoneTypeZone, payload.ZoneTypeString)
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

func TestUpdateZoneKeepsExistingPeerZoneType(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/peer_zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <zone>
    <zone-name>peer_zone</zone-name>
    <zone-type-string>user-created-peer-zone</zone-type-string>
  </zone>
</Response>`))
	})
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var payload DefinedZoneAPI
		if err := xml.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if payload.ZoneTypeString != ZoneTypeUserCreatedPeerZone {
			t.Fatalf("expected zone type %q, got %q", ZoneTypeUserCreatedPeerZone, payload.ZoneTypeString)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.UpdateZone(
		"peer_zone",
		[]string{"10:10:10:27:f8:f0:2a:e8"},
		[]string{"10:10:10:27:f8:8f:44:cd"},
	)
	if err != nil {
		t.Fatalf("UpdateZone() error: %v", err)
	}
}

func TestUpdateZoneRejectsPrincipalMembersForExistingNormalZone(t *testing.T) {
	var patchCalled bool
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_A", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <zone>
    <zone-name>zone_A</zone-name>
    <zone-type-string>zone</zone-type-string>
  </zone>
</Response>`))
	})
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone", func(w http.ResponseWriter, r *http.Request) {
		patchCalled = true
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.UpdateZone("zone_A", []string{"member"}, []string{"principal"})
	if err == nil {
		t.Fatal("expected UpdateZone to reject principal members for existing normal zone")
	}
	if patchCalled {
		t.Fatal("expected UpdateZone to return before PATCH")
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
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_A", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <zone>
    <zone-name>zone_A</zone-name>
    <zone-type-string>zone</zone-type-string>
  </zone>
</Response>`))
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
		"GET /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone_A",
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
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/yang-data+xml")
			w.Write([]byte(`<?xml version="1.0"?><Response><zone><zone-name>zone_A</zone-name><zone-type-string>zone</zone-type-string></zone></Response>`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.DeleteZone("zone_A")
	if err != nil {
		t.Fatalf("DeleteZone() error: %v", err)
	}
}

func TestDeleteZoneRejectsMissingZone(t *testing.T) {
	var deleteCalled bool
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/defined-configuration/zone/zone-name/missing_zone", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeNotFound(w)
		case http.MethodDelete:
			deleteCalled = true
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	if err := c.DeleteZone("missing_zone"); err == nil {
		t.Fatal("expected DeleteZone to reject missing zone")
	}
	if deleteCalled {
		t.Fatal("expected DeleteZone to return before DELETE")
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
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fdelete":
			w.Header().Set("Content-Type", "application/yang-data+xml")
			w.Write([]byte(`<?xml version="1.0"?><Response><zone><zone-name>zone/delete</zone-name><zone-type-string>zone</zone-type-string></zone></Response>`))
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
		"GET /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fdelete",
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
	var calls []string
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.EscapedPath())
		wantPath := "/rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fwith%20space%3F"
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Errorf("expected escaped path %q, got %q", wantPath, got)
		}
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/yang-data+xml")
			w.Write([]byte(`<?xml version="1.0"?><Response><zone><zone-name>zone/with space?</zone-name><zone-type-string>zone</zone-type-string></zone></Response>`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	err := c.DeleteZone("zone/with space?")
	if err != nil {
		t.Fatalf("DeleteZone() error: %v", err)
	}
	want := []string{
		"GET /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fwith%20space%3F",
		"DELETE /rest/running/brocade-zone/defined-configuration/zone/zone-name/zone%2Fwith%20space%3F",
	}
	if got := strings.Join(calls, "\n"); got != strings.Join(want, "\n") {
		t.Fatalf("unexpected calls:\n%s", got)
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

func TestSaveZoneConfigSendsChecksumOnly(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-action-v2/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if got := string(body); got != "<checksum>abc</checksum>" {
			t.Fatalf("unexpected body: %s", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	if err := c.SaveZoneConfig("abc"); err != nil {
		t.Fatalf("SaveZoneConfig() error: %v", err)
	}
}

func TestSaveZoneConfigUsesFOS91Endpoint(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-action/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if got := string(body); got != "<checksum>abc</checksum>" {
			t.Fatalf("unexpected body: %s", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "password", WithFOSVersion("v9.1.1"))
	c.baseURL = ts.URL + "/rest/running"

	if err := c.SaveZoneConfig("abc"); err != nil {
		t.Fatalf("SaveZoneConfig() error: %v", err)
	}
}

func TestSaveZoneConfigDoesNotTreatFOS910AsFOS91(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-action-v2/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := NewClient("localhost", "admin", "password", WithFOSVersion("v9.10.0"))
	c.baseURL = ts.URL + "/rest/running"

	if err := c.SaveZoneConfig("abc"); err != nil {
		t.Fatalf("SaveZoneConfig() error: %v", err)
	}
}

func TestActivateZoneConfigSendsChecksumOnly(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-zone/effective-configuration/cfg-name/cfg1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if got := string(body); got != "<checksum>abc</checksum>" {
			t.Fatalf("unexpected body: %s", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	if err := c.ActivateZoneConfig("cfg1", "abc"); err != nil {
		t.Fatalf("ActivateZoneConfig() error: %v", err)
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
