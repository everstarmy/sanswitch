package sanswitch

import (
	"fmt"
	"strconv"
	"strings"
)

// Version identifies a Fabric OS release. A zero Version means that the
// switch did not report a usable version.
type Version struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses Fabric OS forms such as v9.2.0a and 9.1.1.
func ParseVersion(value string) (Version, error) {
	normalized := strings.TrimSpace(strings.ToLower(value))
	normalized = strings.TrimPrefix(normalized, "v")
	if normalized == "" {
		return Version{}, fmt.Errorf("%w: empty version", ErrInvalidVersion)
	}

	parts := strings.Split(normalized, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return Version{}, fmt.Errorf("%w: %q", ErrInvalidVersion, value)
	}

	values := [3]int{}
	for i, part := range parts {
		digits := leadingDigits(part)
		if digits == "" {
			return Version{}, fmt.Errorf("%w: %q", ErrInvalidVersion, value)
		}
		parsed, err := strconv.Atoi(digits)
		if err != nil || parsed < 0 {
			return Version{}, fmt.Errorf("%w: %q", ErrInvalidVersion, value)
		}
		values[i] = parsed
	}

	version := Version{Major: values[0], Minor: values[1], Patch: values[2]}
	if !version.Valid() {
		return Version{}, fmt.Errorf("%w: %q", ErrInvalidVersion, value)
	}
	return version, nil
}

// Valid reports whether v represents a known Fabric OS version.
func (v Version) Valid() bool {
	return v.Major > 0
}

// AtLeast reports whether v is greater than or equal to major.minor.
func (v Version) AtLeast(major, minor int) bool {
	return v.Valid() && (v.Major > major || v.Major == major && v.Minor >= minor)
}

// String returns v in normalized vX.Y.Z form or "unknown".
func (v Version) String() string {
	if !v.Valid() {
		return "unknown"
	}
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Capabilities describes version-gated Fabric OS operations.
type Capabilities struct {
	VersionKnown       bool
	ZoneWrite          bool
	FRUHistory         bool
	Logging            bool
	FirmwareHistory    bool
	ModernZoneEndpoint bool
}

func capabilitiesFor(version Version, allowUnknownWrites bool) Capabilities {
	if !version.Valid() {
		return Capabilities{ZoneWrite: allowUnknownWrites}
	}
	return Capabilities{
		VersionKnown:       true,
		ZoneWrite:          version.AtLeast(9, 1),
		FRUHistory:         version.AtLeast(9, 0),
		Logging:            version.AtLeast(9, 0),
		FirmwareHistory:    version.AtLeast(9, 1),
		ModernZoneEndpoint: version.AtLeast(9, 2),
	}
}

func leadingDigits(value string) string {
	for i, r := range value {
		if r < '0' || r > '9' {
			return value[:i]
		}
	}
	return value
}
