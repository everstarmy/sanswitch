package sanswitch_test

import (
	"reflect"
	"slices"
	"testing"

	sanswitch "github.com/everstarmy/sanswitch"
)

func TestPublicMethodSurface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		typeOf reflect.Type
		want   []string
	}{
		{
			name:   "Client",
			typeOf: reflect.TypeFor[*sanswitch.Client](),
			want: []string{
				"AbortZoneTransaction", "ActivateZoneConfig", "AuditLogs",
				"BeginZoneTransaction", "Blades", "Capabilities", "ClockServer", "Close",
				"CreateAlias", "CreateZone", "DefinedAliases", "DefinedConfigs", "DefinedZone",
				"DefinedZones", "DeleteAlias", "DeleteZone", "EffectiveConfig", "EffectiveZones",
				"ErrorLogs", "FDMIHBAs", "FDMIPorts", "FabricSwitches", "Fans",
				"FibreChannelNameServers", "FibreChannelStatistics", "FirmwareHistory", "Hardware",
				"HistoryLogs", "IsLoggedIn", "LogicalSwitches", "Login", "Logout", "MediaRDPs",
				"Ports", "PowerSupplies", "RenameAlias", "RenameZone", "RetryCount", "SaveZoneConfig",
				"Sensors", "SetLogOutput", "SetVFID", "SetVerbose", "SNMPSystem", "SNMPv1Accounts",
				"SNMPv1Traps", "SNMPv3Accounts", "SNMPv3Traps", "Switch", "SwitchStatus",
				"SystemResources", "Timeout", "TimeZone", "TrunkAreas", "TrunkPerformances", "Trunks",
				"UpdateAlias", "UpdateDefinedConfig", "UpdateZone", "Version", "ZoneChecksum",
				"ZoneDatabase", "ZoneTransactionStatus",
			},
		},
		{
			name:   "ZoneTransaction",
			typeOf: reflect.TypeFor[*sanswitch.ZoneTransaction](),
			want: []string{
				"Abort", "AddZoneToConfig", "Commit", "CreateZone", "DeleteZone",
				"RemoveZoneFromConfig", "ReplaceZone",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := make([]string, 0, test.typeOf.NumMethod())
			for index := range test.typeOf.NumMethod() {
				got = append(got, test.typeOf.Method(index).Name)
			}
			slices.Sort(test.want)
			if !slices.Equal(got, test.want) {
				t.Fatalf("public methods:\n got: %v\nwant: %v", got, test.want)
			}
		})
	}
}
