# go-agent-runtime

Shared runtime substrate that composes the lower-level agent libraries into
reusable boot, turn, session and resume primitives. It is not a daemon, does
not own app persistence, and owns neither session rows nor a shared database
schema — `go-agent-sessions` still owns the process, `go-agent-launch` still
owns the launch plan, and the app still owns its data.

## Start Here

- `README.md` carries the package table and the relationship to the libraries
  underneath.
- `turn/` frames a user turn per runtime kind; `SendTurn` is the public
  abstraction over the raw `SendInput` escape hatch.
- `runtimekind/` is the runtime vocabulary; `runtimebind/` owns binding
  defaults and app override hooks.
- `sessionkit/` projects first-turn policy onto `StartOptions`.
- `bootdir/` holds the native file/overlay helpers and the path-safety
  boundary.
- `loopback/` builds MCP loopback descriptors and renders `.mcp.json`.
- `checkpoint/` holds provider session ID and resume hint types.

## Commands

```bash
go vet ./...
go test -race -count=1 ./...
```

CI runs both.

## Boundaries

This module was absorbed into `agentkit` as `agentkit/agentruntime` at agentkit
v0.1.0, and this repo is maintenance-only. New work belongs in `agentkit`.

`CHANGELOG.md` and the git tags are the authority for what has shipped here.

`turn` exists so streaming-stdio and JSON-RPC sessions never receive raw
markdown. Claude streaming emits pinned NDJSON and Codex app-server speaks
JSON-RPC; only serve-http, subprocess, API and PTY pass through raw. Reaching
past `SendTurn` to `SendInput` for a framed runtime is the mistake this package
prevents.

`bootdir` validates every path it writes and reuses `go-agent-launch`'s
validation rather than duplicating it. `TestBuildInjectionRejectsUnsafePaths`
and `TestEmptyBootDirRejected` guard the boundary; `TestPlanInjectionSpecMatchesActualWriteOrder`
guards that the planned order is the written order.

Resume must not replay the first turn — `TestFirstTurnPolicyDoesNotDoubleSendOnResume`
is the guard. `.mcp.json` rendering is deterministic across multiple
descriptors (`TestRenderMCPJSONMultipleDescriptorsDeterministic`).
