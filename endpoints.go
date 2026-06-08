package san

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const legacyFOSVersion = "v8.2"

type endpoints struct {
	version fosMajorMinor
}

type fosMajorMinor struct {
	major int
	minor int
	known bool
}

func (c *Client) endpoints() endpoints {
	return endpoints{version: parseFOSMajorMinor(c.fosVersion)}
}

func (e endpoints) allowWrite() bool {
	if !e.version.known {
		return true
	}
	if e.version.major > 9 {
		return true
	}
	return e.version.major == 9 && e.version.minor >= 1
}

func (e endpoints) allowFRUHistoryLogSensor() bool {
	return e.allowFOS90Endpoint()
}

func (e endpoints) allowLogging() bool {
	return e.allowFOS90Endpoint()
}

func (e endpoints) allowFOS90Endpoint() bool {
	if !e.version.known {
		return true
	}
	return e.version.major >= 9
}

func (e endpoints) allowFirmwareHistory() bool {
	if !e.version.known {
		return true
	}
	if e.version.major > 9 {
		return true
	}
	return e.version.major == 9 && e.version.minor >= 1
}

func (e endpoints) ZoneSaveConfig() string {
	if e.version.known && e.version.major == 9 && e.version.minor == 1 {
		return "/brocade-zone/effective-configuration/cfg-action/1"
	}
	return "/brocade-zone/effective-configuration/cfg-action-v2/save"
}

func (endpoints) Login() string {
	return "/login"
}

func (endpoints) Logout() string {
	return "/logout"
}

func (endpoints) DefinedZones() string {
	return "/brocade-zone/defined-configuration/zone"
}

func (endpoints) DefinedZone(name string) string {
	return "/brocade-zone/defined-configuration/zone/zone-name/" + url.PathEscape(name)
}

func (endpoints) EffectiveZones() string {
	return "/brocade-zone/effective-configuration/enabled-zone"
}

func (endpoints) DefinedAliases() string {
	return "/brocade-zone/defined-configuration/alias"
}

func (endpoints) DefinedAlias(name string) string {
	return "/brocade-zone/defined-configuration/alias/alias-name/" + url.PathEscape(name)
}

func (endpoints) DefinedConfigs() string {
	return "/brocade-zone/defined-configuration/cfg"
}

func (endpoints) EffectiveConfig() string {
	return "/brocade-zone/effective-configuration"
}

func (endpoints) ZoneChecksum() string {
	return "/brocade-zone/effective-configuration/checksum"
}

func (endpoints) ZoneActivateConfig(name string) string {
	return "/brocade-zone/effective-configuration/cfg-name/" + url.PathEscape(name)
}

func (e endpoints) ZoneAbortTransaction() string {
	if e.version.known && (e.version.major < 9 || e.version.major == 9 && e.version.minor <= 1) {
		return "/brocade-zone/effective-configuration/cfg-action/4"
	}
	return "/brocade-zone/effective-configuration/cfg-action-v2/transaction-abort"
}

func (endpoints) ZoneTransactionStatus() string {
	return "/brocade-zone/effective-configuration/transaction-token"
}

func (endpoints) FabricSwitches() string {
	return "/brocade-fabric/fabric-switch"
}

func (endpoints) FibreChannelPorts() string {
	return "/brocade-interface/fibrechannel"
}

func (endpoints) FibreChannelStatistics() string {
	return "/brocade-interface/fibrechannel-statistics"
}

func (endpoints) LogicalSwitches() string {
	return "/brocade-fibrechannel-logical-switch/fibrechannel-logical-switch"
}

func (endpoints) Blades() string {
	return "/brocade-fru/blade"
}

func (endpoints) Fans() string {
	return "/brocade-fru/fan"
}

func (endpoints) PowerSupplies() string {
	return "/brocade-fru/power-supply"
}

func (endpoints) HistoryLogs() string {
	return "/brocade-fru/history-log"
}

func (endpoints) Sensors() string {
	return "/brocade-fru/sensor"
}

func (endpoints) Trunks() string {
	return "/brocade-fibrechannel-trunk/trunk"
}

func (endpoints) TrunkPerformances() string {
	return "/brocade-fibrechannel-trunk/performance"
}

func (endpoints) TrunkAreas() string {
	return "/brocade-fibrechannel-trunk/trunk-area"
}

func (endpoints) SNMPSystem() string {
	return "/brocade-snmp/system"
}

func (endpoints) SNMPv1Accounts() string {
	return "/brocade-snmp/v1-account"
}

func (endpoints) SNMPv1Traps() string {
	return "/brocade-snmp/v1-trap"
}

func (endpoints) SNMPv3Accounts() string {
	return "/brocade-snmp/v3-account"
}

func (endpoints) SNMPv3Traps() string {
	return "/brocade-snmp/v3-trap"
}

func (endpoints) Chassis() string {
	return "/brocade-chassis/chassis"
}

func (endpoints) TimeZone() string {
	return "/brocade-time/time-zone"
}

func (endpoints) ClockServer() string {
	return "/brocade-time/clock-server"
}

func (endpoints) ErrorLogs() string {
	return "/brocade-logging/error-log"
}

func (endpoints) AuditLogs() string {
	return "/brocade-logging/audit-log"
}

func (endpoints) SwitchStatusPolicyReport() string {
	return "/brocade-maps/switch-status-policy-report"
}

func (endpoints) SystemResources() string {
	return "/brocade-maps/system-resources"
}

func (endpoints) MediaRDPs() string {
	return "/brocade-media/media-rdp"
}

func (endpoints) NameServer() string {
	return "/brocade-name-server/fibrechannel-name-server"
}

func (endpoints) FDMIHBAs() string {
	return "/brocade-fdmi/hba"
}

func (endpoints) FDMIPorts() string {
	return "/brocade-fdmi/port"
}

func (endpoints) FirmwareHistory() string {
	return "/brocade-firmware/firmware-history"
}

func parseFOSMajorMinor(version string) fosMajorMinor {
	normalized := strings.TrimSpace(strings.ToLower(version))
	normalized = strings.TrimPrefix(normalized, "v")
	if normalized == "" {
		return fosMajorMinor{}
	}
	parts := strings.Split(normalized, ".")
	if len(parts) < 2 {
		return fosMajorMinor{}
	}
	major, err := strconv.Atoi(leadingDigits(parts[0]))
	if err != nil {
		return fosMajorMinor{}
	}
	minor, err := strconv.Atoi(leadingDigits(parts[1]))
	if err != nil {
		return fosMajorMinor{}
	}
	return fosMajorMinor{major: major, minor: minor, known: true}
}

func leadingDigits(value string) string {
	for i, r := range value {
		if r < '0' || r > '9' {
			return value[:i]
		}
	}
	return value
}

func (v fosMajorMinor) String() string {
	if !v.known {
		return "unknown"
	}
	return fmt.Sprintf("v%d.%d", v.major, v.minor)
}
