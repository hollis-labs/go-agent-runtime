# Changelog

All notable changes to go-agent-runtime are documented here. The format
follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Changed

- Raised the module's `go` directive to `1.26.6` (Go floor across the portfolio); CI now uses `go-version-file: go.mod`.

### Added

- Initial module scaffold from folio's `go-lib` preset — importable
  `agentruntime` package, CI workflow, MIT license.
