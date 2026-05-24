package bootdir

import (
	"errors"
	"testing"

	agentlaunch "github.com/hollis-labs/go-agent-launch/agentlaunch"
)

func TestBuildInjectionRejectsUnsafePaths(t *testing.T) {
	_, _, err := BuildInjection(Request{Overlays: map[string]string{"../escape": "x"}})
	if !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("err = %v, want ErrUnsafePath", err)
	}
}

func TestTaskBundleRejectsUnsafePaths(t *testing.T) {
	_, err := TaskBundle("tasks", map[string]string{"../bad.md": "x"})
	if err == nil {
		t.Fatal("expected unsafe path error")
	}
}

func TestBuildInjectionNativeFiles(t *testing.T) {
	_, _, err := BuildInjection(Request{NativeFiles: []agentlaunch.NativeFile{{
		Kind: agentlaunch.NativeFileRaw, RelPath: "tasks/README.md", Content: "read me",
	}}})
	if err != nil {
		t.Fatal(err)
	}
}
