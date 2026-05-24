package sessionkit

import (
	"testing"

	"github.com/hollis-labs/go-agent-runtime/runtimekind"
	"github.com/hollis-labs/go-agent-runtime/turn"
	agentsessions "github.com/hollis-labs/go-agent-sessions/agentsessions"
)

func TestFirstTurnPolicyDoesNotDoubleSendOnResume(t *testing.T) {
	opts := agentsessions.StartOptions{AutoFireFirstTurn: true, FirstTurnPayload: []byte("old")}
	if err := ApplyFirstTurnPolicy(&opts, FirstTurnPolicy{Mode: ResumeWithoutFirstTurn}); err != nil {
		t.Fatal(err)
	}
	if opts.AutoFireFirstTurn || len(opts.FirstTurnPayload) != 0 {
		t.Fatalf("resume policy left auto-fire enabled: %#v", opts)
	}
}

func TestFirstTurnPolicyFramesAutoFire(t *testing.T) {
	var opts agentsessions.StartOptions
	if err := ApplyFirstTurnPolicy(&opts, FirstTurnPolicy{
		Mode:   AutoFireFirstTurn,
		Prompt: "hello",
		Turn:   turn.Options{Runtime: runtimekind.StreamingStdio},
	}); err != nil {
		t.Fatal(err)
	}
	if !opts.AutoFireFirstTurn || string(opts.FirstTurnPayload) == "hello" {
		t.Fatalf("auto-fire not framed: %#v", opts)
	}
}
