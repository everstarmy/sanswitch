package san

import (
	"net/http"
	"testing"
)

const fabricSwitchXML = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <fabric-switch>
    <name>10:00:50:eb:1a:0b:00:00</name>
    <switch-user-friendly-name>san-sw01</switch-user-friendly-name>
    <chassis-wwn>10:00:50:eb:1a:0b:00:01</chassis-wwn>
    <chassis-user-friendly-name>chassis-01</chassis-user-friendly-name>
    <domain-id>1</domain-id>
    <fcid>0x010000</fcid>
    <fcid-hex>0x010000</fcid-hex>
    <ip-address>192.168.1.100</ip-address>
    <ipv6-address>fe80::1</ipv6-address>
    <firmware-version>v9.2.0a</firmware-version>
    <switch-model>G620</switch-model>
    <serial-number>FOS1234567</serial-number>
    <principal>1</principal>
    <is-principal>true</is-principal>
  </fabric-switch>
</Response>`

func TestGetFabricSwitches(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-fabric/fabric-switch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(fabricSwitchXML))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	switches, err := c.GetFabricSwitches()
	if err != nil {
		t.Fatalf("GetFabricSwitches() error: %v", err)
	}
	if len(switches) != 1 {
		t.Fatalf("expected 1 switch, got %d", len(switches))
	}

	sw := switches[0]
	if sw.Name != "10:00:50:eb:1a:0b:00:00" {
		t.Errorf("unexpected Name: %q", sw.Name)
	}
	if sw.SwitchUserFriendlyName != "san-sw01" {
		t.Errorf("unexpected SwitchUserFriendlyName: %q", sw.SwitchUserFriendlyName)
	}
	if sw.DomainID != 1 {
		t.Errorf("unexpected DomainID: %d", sw.DomainID)
	}
	if sw.IPAddress != "192.168.1.100" {
		t.Errorf("unexpected IPAddress: %q", sw.IPAddress)
	}
	if sw.FirmwareVersion != "v9.2.0a" {
		t.Errorf("unexpected FirmwareVersion: %q", sw.FirmwareVersion)
	}
	if sw.SwitchModel != "G620" {
		t.Errorf("unexpected SwitchModel: %q", sw.SwitchModel)
	}
	if sw.SerialNumber != "FOS1234567" {
		t.Errorf("unexpected SerialNumber: %q", sw.SerialNumber)
	}
	if !sw.IsPrincipal {
		t.Error("expected IsPrincipal = true")
	}
}

func TestGetSwitchInfo(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-fabric/fabric-switch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(fabricSwitchXML))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	info, err := c.GetSwitchInfo()
	if err != nil {
		t.Fatalf("GetSwitchInfo() error: %v", err)
	}
	if info.Name != "san-sw01" {
		t.Errorf("unexpected Name: %q", info.Name)
	}
	if info.WWN != "10:00:50:eb:1a:0b:00:00" {
		t.Errorf("unexpected WWN: %q", info.WWN)
	}
	if info.ChassisWWN != "10:00:50:eb:1a:0b:00:01" {
		t.Errorf("unexpected ChassisWWN: %q", info.ChassisWWN)
	}
	if info.FirmwareVersion != "v9.2.0a" {
		t.Errorf("unexpected FirmwareVersion: %q", info.FirmwareVersion)
	}
	if info.ModelName != "G620" {
		t.Errorf("unexpected ModelName: %q", info.ModelName)
	}
	if info.SerialNumber != "FOS1234567" {
		t.Errorf("unexpected SerialNumber: %q", info.SerialNumber)
	}
	if info.IPAddress != "192.168.1.100" {
		t.Errorf("unexpected IPAddress: %q", info.IPAddress)
	}
}

func TestGetSwitchInfoNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-fabric/fabric-switch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?><Response></Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	_, err := c.GetSwitchInfo()
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
