# Changelog

All notable changes to this module are recorded here. The project follows
Semantic Versioning.

## [Unreleased]

## [0.2.0] - 2026-08-31

### Breaking changes

- Removed the monolithic `SwitchAPI`; use the focused `Session`,
  `SwitchReader`, `ZoneReader`, `ZoneWriter`, `InventoryReader`, and
  `MonitoringReader` interfaces.
- Renamed `FDMIportInfo` to `FDMIPortInfo` and standardized the method name to
  `GetFDMIPorts`.

### Added

- Context-aware methods for every network-facing Client and SANSwitch
  operation.
- `Close` methods for releasing idle HTTP connections.
- A bounded response-body reader with `WithMaxResponseBodyBytes` and
  `ErrResponseBodyTooLarge`.
- External-package examples, package documentation, and CI module checks.

### Security and maintenance

- Credentials in documentation examples now come from environment variables.
- The module remains standard-library-only at runtime and has no `go.sum` or
  vendored third-party dependencies.
