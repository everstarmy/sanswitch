package san

import (
	"net/http"
	"testing"
)

const portXML = `<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <fibrechannel>
    <name>0/0</name>
    <wwn>20:00:50:eb:1a:0b:0c:0d</wwn>
    <operational-status>3</operational-status>
    <operational-status-string>Online</operational-status-string>
    <is-enabled-state>true</is-enabled-state>
    <user-friendly-name>Port 0</user-friendly-name>
    <protocol-speed>N8Gbps</protocol-speed>
    <max-protocol-speed>N16Gbps</max-protocol-speed>
    <auto-negotiate-v2>true</auto-negotiate-v2>
    <port-type-string>F-Port</port-type-string>
    <fcid>0x010100</fcid>
    <fcid-hex>0x010100</fcid-hex>
    <npiv-enabled-v2>true</npiv-enabled-v2>
    <blade-port-number>0</blade-port-number>
    <neighbor>
      <wwn>10:00:00:00:c9:f8:04:35</wwn>
    </neighbor>
    <neighbor-slot-port>1/0</neighbor-slot-port>
  </fibrechannel>
  <fibrechannel>
    <name>0/1</name>
    <wwn>20:00:50:eb:1a:0b:0c:0e</wwn>
    <operational-status>2</operational-status>
    <operational-status-string>No_Light</operational-status-string>
    <is-enabled-state>true</is-enabled-state>
    <user-friendly-name>Port 1</user-friendly-name>
    <protocol-speed>N4Gbps</protocol-speed>
    <max-protocol-speed>N16Gbps</max-protocol-speed>
    <auto-negotiate-v2>true</auto-negotiate-v2>
    <port-type-string>Unknown</port-type-string>
    <fcid>0x000000</fcid>
    <fcid-hex>0x000000</fcid-hex>
    <blade-port-number>1</blade-port-number>
  </fibrechannel>
</Response>`

func TestGetPorts(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-interface/fibrechannel", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(portXML))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	ports, err := c.GetPorts()
	if err != nil {
		t.Fatalf("GetPorts() error: %v", err)
	}
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}

	p0 := ports[0]
	if p0.Name != "0/0" {
		t.Errorf("unexpected port name: %q", p0.Name)
	}
	if p0.WWN != "20:00:50:eb:1a:0b:0c:0d" {
		t.Errorf("unexpected WWN: %q", p0.WWN)
	}
	if p0.OperationalStatusString != "Online" {
		t.Errorf("unexpected OperationalStatusString: %q", p0.OperationalStatusString)
	}
	if !p0.EnabledState {
		t.Error("expected EnabledState = true")
	}
	if p0.Speed != "N8Gbps" {
		t.Errorf("unexpected Speed: %q", p0.Speed)
	}
	if p0.MaxSpeed != "N16Gbps" {
		t.Errorf("unexpected MaxSpeed: %q", p0.MaxSpeed)
	}
	if p0.PortType != "F-Port" {
		t.Errorf("unexpected PortType: %q", p0.PortType)
	}
	if p0.FCID != "0x010100" {
		t.Errorf("unexpected FCID: %q", p0.FCID)
	}
	if !p0.AutoNegotiate {
		t.Error("expected AutoNegotiate = true")
	}
	if !p0.NPIVEnabled {
		t.Error("expected NPIVEnabled = true")
	}
	if len(p0.NeighborWWNs) != 1 || p0.NeighborWWNs[0] != "10:00:00:00:c9:f8:04:35" {
		t.Errorf("unexpected NeighborWWNs: %v", p0.NeighborWWNs)
	}
	if p0.NeighborSlotPort != "1/0" {
		t.Errorf("unexpected NeighborSlotPort: %q", p0.NeighborSlotPort)
	}

	p1 := ports[1]
	if p1.Name != "0/1" {
		t.Errorf("unexpected port name: %q", p1.Name)
	}
	if p1.OperationalStatusString != "No_Light" {
		t.Errorf("unexpected OperationalStatusString: %q", p1.OperationalStatusString)
	}
}
