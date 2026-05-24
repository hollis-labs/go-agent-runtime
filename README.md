# go-agent-runtime

Shared Go agent runtime substrate for Hollis Labs apps.

This library is not a daemon and does not own app persistence. It composes the
lower-level Hollis Labs agent libraries into reusable boot, launch, turn,
session, and resume primitives that Agridd, Torque, Tether, and future apps can
import without sharing an application process.

## Packages

| Package | Responsibility |
|---|---|
| `runtimekind` | Runtime-kind vocabulary and compatibility aliases: `api`, `subprocess`, `streaming-stdio`, `jsonrpc-stdio`, `serve-http`, `pty`, `pty-debug`, `unknown`. |
| `runtimebind` | Provider/runtime binding defaults, explicit generic-subprocess opt-in, app override hooks, and shared Codex headless policy shape. |
| `turn` | Runtime-aware user turn framing. Claude streaming-stdio emits pinned NDJSON, jsonrpc-stdio can use typed calls, subprocess/API/PTY stay raw. `SendTurn` is the public abstraction; provider-specific JSON-RPC methods remain adapter/binding details. |
| `sessionkit` | First-turn policy projection onto `go-agent-sessions.StartOptions` without owning session rows. |
| `bootdir` | App-owned native file/task/overlay helpers with the shared path-safety boundary. |
| `loopback` | Provider-neutral MCP loopback descriptors for subprocess, HTTP/SSE, and mux proxy entries. |
| `checkpoint` | Provider session ID and checkpoint/resume hint types without a shared database schema. |
| `smoke` | Anchor for cross-package smoke fixtures; executable tests live beside the packages they protect. |

## Relationship To Existing Libraries

`go-agent-launch` remains the launch-plan, compile, prepare, and providerplant
library. `go-agent-runtime/bootdir` reuses its `InjectionSpec`, `NativeFile`,
and path validation instead of duplicating the planter.

`go-agent-sessions` remains the process/session owner: `Start`, `SendInput`,
`Attach`, `Wait`, `Stop`, `JsonRpcCall`, event fanout, and provider session ID
callbacks. `go-agent-runtime/turn` sits above its raw `SendInput` escape hatch
so streaming/jsonrpc sessions do not receive unsafe raw markdown.

`go-providers` remains the adapter and provider boot-file library. Runtime
binding here names the intended provider/runtime behavior but does not replace
provider adapter constructors.

`go-sandbox` remains the OS sandbox substrate. This library only captures
runtime policy shape such as Codex jsonrpc/app-server bypass handling; apps
still decide writable roots and sandbox posture.

Agridd, Torque, and Tether continue to own product content: prompts, task
bundles, MCP tool surfaces, launch catalogs, scheduler state, chat messages,
checkpoint tables, and attach UX. This library gives them shared primitives for
the recurring mechanics and bug classes.

Unknown provider IDs fail closed with `runtimebind.ErrUnsupportedBinding`.
Callers that intentionally host a generic subprocess adapter must set
`AllowGenericSubprocess`; this keeps accidental provider typos from silently
launching in the wrong runtime.

`bootdir.Writer` owns the mechanical boot directory write lifecycle for apps
that render their own content: full `agentlaunch.InjectionSpec`
populate/re-populate, slot-only `WriteFiles`, deterministic dry-run planning,
default atomic writes, an atomic writer hook, sorted overlay-wins-last ordering,
and sentinel errors for unsafe paths, unsupported native file kinds, empty boot
dirs, and empty relative paths.

## Install

```sh
go get github.com/hollis-labs/go-agent-runtime
```
