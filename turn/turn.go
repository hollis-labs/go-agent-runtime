// Package turn frames user turns before they are sent to go-agent-sessions.
package turn

import (
	"context"
	"encoding/json"
	"errors"

	agentlaunch "github.com/hollis-labs/go-agent-launch/agentlaunch"
	"github.com/hollis-labs/go-agent-runtime/runtimekind"
)

var ErrUnsupportedRuntime = errors.New("turn: unsupported runtime")

type Options struct {
	Provider string
	Runtime  agentlaunch.RuntimeKind
	// JSONRPCMethod is a provider-specific wire detail. Public callers should
	// prefer SendTurn and let the provider binding/adapter choose the real
	// JSON-RPC method. Tests and custom adapters can override it here.
	JSONRPCMethod string
}

type Sender interface {
	SendInput(ctx context.Context, data []byte) error
}

type JSONRPCSender interface {
	Call(ctx context.Context, method string, params any) (json.RawMessage, error)
}

// Frame returns the payload for raw SendInput. JSON-RPC stdio callers should
// prefer SendTurn so typed calls do not go through the raw byte escape hatch.
func Frame(text string, opts Options) ([]byte, error) {
	switch runtimekind.Parse(string(opts.Runtime)) {
	case runtimekind.StreamingStdio:
		return ClaudeStreamingUserFrame(text)
	case runtimekind.Subprocess, runtimekind.PTY, runtimekind.PTYDebug:
		return []byte(text), nil
	case runtimekind.API:
		return []byte(text), nil
	case runtimekind.JSONRPCStdio:
		params := map[string]any{"message": text}
		return json.Marshal(map[string]any{"method": method(opts), "params": params})
	default:
		return nil, ErrUnsupportedRuntime
	}
}

// SendTurn applies runtime-specific framing and delivery. JSON-RPC stdio uses a
// typed call when the sender exposes JSONRPCSender; otherwise it falls back to a
// serialized request-shaped frame for adapters that own the final method.
func SendTurn(ctx context.Context, sender Sender, text string, opts Options) error {
	if runtimekind.Parse(string(opts.Runtime)) == runtimekind.JSONRPCStdio {
		if rpc, ok := sender.(JSONRPCSender); ok {
			_, err := rpc.Call(ctx, method(opts), map[string]any{"message": text})
			return err
		}
	}
	frame, err := Frame(text, opts)
	if err != nil {
		return err
	}
	return sender.SendInput(ctx, frame)
}

// ClaudeStreamingUserFrame emits the exact NDJSON object Claude Code streaming
// stdio consumes. It deliberately returns no trailing newline; sessions appends
// the line break at write time.
func ClaudeStreamingUserFrame(text string) ([]byte, error) {
	type userMsg struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type frame struct {
		Type    string  `json:"type"`
		Message userMsg `json:"message"`
	}
	return json.Marshal(frame{Type: "user", Message: userMsg{Role: "user", Content: text}})
}

func method(opts Options) string {
	if opts.JSONRPCMethod != "" {
		return opts.JSONRPCMethod
	}
	return "turn.send"
}
