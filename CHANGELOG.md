# Changelog

All notable changes to this module are recorded here. The project follows
Semantic Versioning.

## [0.3.0] - 2026-08-31

### Breaking changes

- Raised the minimum Go version to 1.24 and renamed the package to
  `sanswitch`.
- Removed `SANSwitch`, producer-owned capability interfaces, context-free
  network methods, and all `WithContext` method suffixes.
- Replaced `NewClient` with validating `New` and authenticated `Open`.
- `Login` now accepts short-lived `Credentials`; clients no longer retain
  usernames or passwords.
- Renamed public models and read methods to concise Go names such as `Port`,
  `Zone`, `Ports`, and `DefinedZones`.
- Replaced composite Zone-and-activate helpers with explicit
  `ZoneTransaction` operations and `Commit`/`Abort`.

### Added

- Typed `Version` and `Capabilities` APIs.
- `WithTransport` for proxy, tracing, and deterministic testing integration.
- Retry jitter and transient failure classification.
- Contract, fuzz, concurrent session, transaction, and constructor validation
  tests.

### Maintenance

- XML/YANG wire types are private implementation details.
- CI now tests Go 1.24 and stable with randomized tests, race detection, and a
  coverage floor.

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
