package san

import (
	"net/http"
	"testing"
)

// ==================== Trunk Tests ====================

func TestGetTrunks(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-fibrechannel-trunk/trunk", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <trunk>
        <group>1</group>
        <source-port>132</source-port>
        <destination-port>47</destination-port>
        <neighbor-wwn>10:00:d1:1f:cc:1a:f1:1e</neighbor-wwn>
        <neighbor-switch-name>G720</neighbor-switch-name>
        <neighbor-domain-id>238</neighbor-domain-id>
        <deskew>5</deskew>
        <master>true</master>
        <trunk-type>inter-switch-link</trunk-type>
    </trunk>
    <trunk>
        <group>6</group>
        <source-port>251</source-port>
        <destination-port>9</destination-port>
        <neighbor-wwn>10:00:28:24:21:22:2b:d2</neighbor-wwn>
        <neighbor-switch-name>B7810</neighbor-switch-name>
        <neighbor-domain-id>237</neighbor-domain-id>
        <deskew>0</deskew>
        <master>true</master>
        <trunk-type>inter-switch-link</trunk-type>
    </trunk>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	trunks, err := c.GetTrunks()
	if err != nil {
		t.Fatalf("GetTrunks failed: %v", err)
	}
	if len(trunks) != 2 {
		t.Fatalf("expected 2 trunks, got %d", len(trunks))
	}
	if trunks[0].Group != 1 {
		t.Errorf("expected group=1, got %d", trunks[0].Group)
	}
	if trunks[0].SourcePort != 132 {
		t.Errorf("expected source-port=132, got %d", trunks[0].SourcePort)
	}
	if trunks[0].NeighborWWN != "10:00:d1:1f:cc:1a:f1:1e" {
		t.Errorf("unexpected neighbor-wwn: %s", trunks[0].NeighborWWN)
	}
	if !trunks[0].Master {
		t.Error("expected master=true")
	}
	if trunks[0].TrunkType != "inter-switch-link" {
		t.Errorf("unexpected trunk-type: %s", trunks[0].TrunkType)
	}
	if trunks[1].Group != 6 {
		t.Errorf("expected group=6, got %d", trunks[1].Group)
	}
	if trunks[1].NeighborSwitchName != "B7810" {
		t.Errorf("unexpected neighbor-switch-name: %s", trunks[1].NeighborSwitchName)
	}
}

func TestGetTrunkPerformances(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-fibrechannel-trunk/performance", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <performance>
        <group>1</group>
        <tx-bandwidth>32</tx-bandwidth>
        <tx-throughput>8000</tx-throughput>
        <tx-percentage>0.00</tx-percentage>
        <rx-bandwidth>32</rx-bandwidth>
        <rx-throughput>0</rx-throughput>
        <rx-percentage>0.00</rx-percentage>
        <txrx-bandwidth>64</txrx-bandwidth>
        <txrx-throughput>8000</txrx-throughput>
        <txrx-percentage>0.00</txrx-percentage>
    </performance>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	perfs, err := c.GetTrunkPerformances()
	if err != nil {
		t.Fatalf("GetTrunkPerformances failed: %v", err)
	}
	if len(perfs) != 1 {
		t.Fatalf("expected 1 performance entry, got %d", len(perfs))
	}
	if perfs[0].Group != 1 {
		t.Errorf("expected group=1, got %d", perfs[0].Group)
	}
	if perfs[0].TxBandwidth != 32 {
		t.Errorf("expected tx-bandwidth=32, got %d", perfs[0].TxBandwidth)
	}
	if perfs[0].TxThroughput != 8000 {
		t.Errorf("expected tx-throughput=8000, got %d", perfs[0].TxThroughput)
	}
	if perfs[0].TxRxBandwidth != 64 {
		t.Errorf("expected txrx-bandwidth=64, got %d", perfs[0].TxRxBandwidth)
	}
}

func TestGetTrunkAreas(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-fibrechannel-trunk/trunk-area", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <trunk-area>
        <trunk-index>0</trunk-index>
        <trunk-active>true</trunk-active>
        <master-port>0/0</master-port>
        <trunk-members>
            <trunk-member>0/0</trunk-member>
            <trunk-member>0/6</trunk-member>
            <trunk-member>0/7</trunk-member>
        </trunk-members>
    </trunk-area>
    <trunk-area>
        <trunk-index>40</trunk-index>
        <trunk-active>false</trunk-active>
        <master-port/>
        <trunk-members>
            <trunk-member>0/40</trunk-member>
            <trunk-member>0/41</trunk-member>
        </trunk-members>
    </trunk-area>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	areas, err := c.GetTrunkAreas()
	if err != nil {
		t.Fatalf("GetTrunkAreas failed: %v", err)
	}
	if len(areas) != 2 {
		t.Fatalf("expected 2 trunk areas, got %d", len(areas))
	}
	if areas[0].TrunkIndex != 0 {
		t.Errorf("expected trunk-index=0, got %d", areas[0].TrunkIndex)
	}
	if !areas[0].TrunkActive {
		t.Error("expected trunk-active=true")
	}
	if areas[0].MasterPort != "0/0" {
		t.Errorf("expected master-port=0/0, got %s", areas[0].MasterPort)
	}
	if len(areas[0].TrunkMembers) != 3 {
		t.Errorf("expected 3 trunk members, got %d", len(areas[0].TrunkMembers))
	}
	if areas[1].TrunkIndex != 40 {
		t.Errorf("expected trunk-index=40, got %d", areas[1].TrunkIndex)
	}
	if areas[1].TrunkActive {
		t.Error("expected trunk-active=false for second trunk area")
	}
	if len(areas[1].TrunkMembers) != 2 {
		t.Errorf("expected 2 trunk members, got %d", len(areas[1].TrunkMembers))
	}
}

// ==================== Firmware History Test ====================

func TestGetFirmwareHistory(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-firmware/firmware-history", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <firmware-history>
        <sequence-number>1</sequence-number>
        <time-stamp>Mon Mar 1 21:11:55 2021</time-stamp>
        <switch-name>unknown(0)</switch-name>
        <slot-number>2</slot-number>
        <process-id>2063</process-id>
        <firmware-version>Fabos Version 9.1.0</firmware-version>
    </firmware-history>
    <firmware-history>
        <sequence-number>2</sequence-number>
        <time-stamp>Tue Mar 2 12:02:55 2021</time-stamp>
        <switch-name>V21SW</switch-name>
        <slot-number>2</slot-number>
        <process-id>2197</process-id>
        <firmware-version>Fabos Version 9.0.0</firmware-version>
    </firmware-history>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	history, err := c.GetFirmwareHistory()
	if err != nil {
		t.Fatalf("GetFirmwareHistory failed: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("expected 2 history entries, got %d", len(history))
	}
	if history[0].SequenceNumber != 1 {
		t.Errorf("expected sequence-number=1, got %d", history[0].SequenceNumber)
	}
	if history[0].SlotNumber != 2 {
		t.Errorf("expected slot-number=2, got %d", history[0].SlotNumber)
	}
	if history[0].ProcessID != 2063 {
		t.Errorf("expected process-id=2063, got %d", history[0].ProcessID)
	}
	if history[1].SwitchName != "V21SW" {
		t.Errorf("expected switch-name=V21SW, got %s", history[1].SwitchName)
	}
	if history[1].FirmwareVersion != "Fabos Version 9.0.0" {
		t.Errorf("unexpected firmware-version: %s", history[1].FirmwareVersion)
	}
}

// ==================== SNMP Tests ====================

func TestGetSNMPSystem(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-snmp/system", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <system>
        <description>Brocade G620</description>
        <location>DC1-Rack5</location>
        <contact>admin@example.com</contact>
        <informs-enabled>false</informs-enabled>
        <encryption-enabled>false</encryption-enabled>
        <audit-interval>0</audit-interval>
        <security-get-level-string>authNoPriv</security-get-level-string>
        <security-set-level-string>authNoPriv</security-set-level-string>
    </system>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	sys, err := c.GetSNMPSystem()
	if err != nil {
		t.Fatalf("GetSNMPSystem failed: %v", err)
	}
	if sys.Description != "Brocade G620" {
		t.Errorf("unexpected description: %s", sys.Description)
	}
	if sys.Location != "DC1-Rack5" {
		t.Errorf("unexpected location: %s", sys.Location)
	}
	if sys.Contact != "admin@example.com" {
		t.Errorf("unexpected contact: %s", sys.Contact)
	}
	if sys.InformsEnabled {
		t.Error("expected informs-enabled=false")
	}
}

func TestGetSNMPv1Accounts(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-snmp/v1-account", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <v1-account>
        <index>1</index>
        <community-group>read-write</community-group>
        <community-name>private</community-name>
    </v1-account>
    <v1-account>
        <index>4</index>
        <community-group>read-only</community-group>
        <community-name>public</community-name>
    </v1-account>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	accounts, err := c.GetSNMPv1Accounts()
	if err != nil {
		t.Fatalf("GetSNMPv1Accounts failed: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}
	if accounts[0].Index != 1 {
		t.Errorf("expected index=1, got %d", accounts[0].Index)
	}
	if accounts[0].CommunityName != "private" {
		t.Errorf("unexpected community-name: %s", accounts[0].CommunityName)
	}
	if accounts[0].CommunityGroup != "read-write" {
		t.Errorf("unexpected community-group: %s", accounts[0].CommunityGroup)
	}
	if accounts[1].CommunityGroup != "read-only" {
		t.Errorf("expected read-only, got %s", accounts[1].CommunityGroup)
	}
}

func TestGetSNMPv1Traps(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-snmp/v1-trap", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <v1-trap>
        <index>1</index>
        <host>0.0.0.0</host>
        <trap-severity-level>none</trap-severity-level>
        <port-number>162</port-number>
    </v1-trap>
    <v1-trap>
        <index>2</index>
        <host>10.10.10.10</host>
        <trap-severity-level>warning</trap-severity-level>
        <port-number>162</port-number>
    </v1-trap>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	traps, err := c.GetSNMPv1Traps()
	if err != nil {
		t.Fatalf("GetSNMPv1Traps failed: %v", err)
	}
	if len(traps) != 2 {
		t.Fatalf("expected 2 traps, got %d", len(traps))
	}
	if traps[0].TrapSeverityLevel != "none" {
		t.Errorf("expected severity=none, got %s", traps[0].TrapSeverityLevel)
	}
	if traps[1].Host != "10.10.10.10" {
		t.Errorf("unexpected host: %s", traps[1].Host)
	}
	if traps[1].TrapSeverityLevel != "warning" {
		t.Errorf("expected severity=warning, got %s", traps[1].TrapSeverityLevel)
	}
}

func TestGetSNMPv3Accounts(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-snmp/v3-account", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <v3-account>
        <index>1</index>
        <user-name>admin</user-name>
        <user-group>read-write</user-group>
        <authentication-protocol>sha</authentication-protocol>
        <privacy-protocol>des</privacy-protocol>
        <authentication-password/>
        <privacy-password/>
        <manager-engine-id>00:00:00:00:00:00:00:00:00</manager-engine-id>
    </v3-account>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	accounts, err := c.GetSNMPv3Accounts()
	if err != nil {
		t.Fatalf("GetSNMPv3Accounts failed: %v", err)
	}
	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}
	if accounts[0].UserName != "admin" {
		t.Errorf("unexpected user-name: %s", accounts[0].UserName)
	}
	if accounts[0].AuthenticationProtocol != "sha" {
		t.Errorf("unexpected auth protocol: %s", accounts[0].AuthenticationProtocol)
	}
	if accounts[0].PrivacyProtocol != "des" {
		t.Errorf("unexpected privacy protocol: %s", accounts[0].PrivacyProtocol)
	}
	if accounts[0].ManagerEngineID != "00:00:00:00:00:00:00:00:00" {
		t.Errorf("unexpected engine-id: %s", accounts[0].ManagerEngineID)
	}
}

func TestGetSNMPv3Traps(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-snmp/v3-trap", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <v3-trap>
        <trap-index>1</trap-index>
        <usm-index>1</usm-index>
        <host>0.0.0.0</host>
        <port-number>162</port-number>
        <trap-severity-level>none</trap-severity-level>
        <informs-enabled>false</informs-enabled>
    </v3-trap>
    <v3-trap>
        <trap-index>2</trap-index>
        <usm-index>2</usm-index>
        <host>192.168.1.100</host>
        <port-number>162</port-number>
        <trap-severity-level>error</trap-severity-level>
        <informs-enabled>true</informs-enabled>
    </v3-trap>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	traps, err := c.GetSNMPv3Traps()
	if err != nil {
		t.Fatalf("GetSNMPv3Traps failed: %v", err)
	}
	if len(traps) != 2 {
		t.Fatalf("expected 2 traps, got %d", len(traps))
	}
	if traps[0].TrapIndex != 1 {
		t.Errorf("expected trap-index=1, got %d", traps[0].TrapIndex)
	}
	if traps[0].USMIndex != 1 {
		t.Errorf("expected usm-index=1, got %d", traps[0].USMIndex)
	}
	if traps[0].InformsEnabled {
		t.Error("expected informs-enabled=false")
	}
	if traps[1].Host != "192.168.1.100" {
		t.Errorf("unexpected host: %s", traps[1].Host)
	}
	if traps[1].TrapSeverityLevel != "error" {
		t.Errorf("expected severity=error, got %s", traps[1].TrapSeverityLevel)
	}
	if !traps[1].InformsEnabled {
		t.Error("expected informs-enabled=true for second trap")
	}
}

// ==================== Time / NTP Tests ====================

func TestGetTimeZone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-time/time-zone", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <time-zone>
        <name>America/Toronto</name>
        <gmt-offset-hours>-5</gmt-offset-hours>
        <gmt-offset-minutes>0</gmt-offset-minutes>
    </time-zone>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	tz, err := c.GetTimeZone()
	if err != nil {
		t.Fatalf("GetTimeZone failed: %v", err)
	}
	if tz.Name != "America/Toronto" {
		t.Errorf("unexpected timezone name: %s", tz.Name)
	}
	if tz.GMTOffsetHours != -5 {
		t.Errorf("expected gmt-offset-hours=-5, got %d", tz.GMTOffsetHours)
	}
	if tz.GMTOffsetMinutes != 0 {
		t.Errorf("expected gmt-offset-minutes=0, got %d", tz.GMTOffsetMinutes)
	}
}

func TestGetClockServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/running/brocade-time/clock-server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+xml")
		w.Write([]byte(`<?xml version="1.0"?>
<Response>
    <clock-server>
        <ntp-server-address>
            <server-address>10.0.0.1</server-address>
            <server-address>10.0.0.2</server-address>
        </ntp-server-address>
        <active-server>10.0.0.1</active-server>
        <ts-auth-spec>noauth</ts-auth-spec>
        <ts-legacy-mode>true</ts-legacy-mode>
    </clock-server>
</Response>`))
	})

	ts := newMockFOS(t, mux)
	c := newTestClient(t, ts)

	cs, err := c.GetClockServer()
	if err != nil {
		t.Fatalf("GetClockServer failed: %v", err)
	}
	if len(cs.NTPServerAddresses) != 2 {
		t.Fatalf("expected 2 NTP servers, got %d", len(cs.NTPServerAddresses))
	}
	if cs.NTPServerAddresses[0] != "10.0.0.1" {
		t.Errorf("unexpected first NTP server: %s", cs.NTPServerAddresses[0])
	}
	if cs.ActiveServer != "10.0.0.1" {
		t.Errorf("unexpected active server: %s", cs.ActiveServer)
	}
	if cs.TSAuthSpec != "noauth" {
		t.Errorf("unexpected ts-auth-spec: %s", cs.TSAuthSpec)
	}
	if !cs.TSLegacyMode {
		t.Error("expected ts-legacy-mode=true")
	}
}
