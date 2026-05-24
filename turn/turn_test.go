package turn

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hollis-labs/go-agent-runtime/runtimekind"
)

func TestClaudeStreamingUserFrame(t *testing.T) {
	raw, err := ClaudeStreamingUserFrame("# Boot\nsay \"hi\"\n")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "\n") {
		t.Fatalf("serialized frame contains raw newline: %q", raw)
	}
	var got struct {
		Type    string `json:"type"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.Type != "user" || got.Message.Role != "user" || got.Message.Content != "# Boot\nsay \"hi\"\n" {
		t.Fatalf("decoded frame = %#v", got)
	}
}

func TestStreamingStdioDoesNotSendRawMarkdown(t *testing.T) {
	s := &captureSender{}
	if err := SendTurn(context.Background(), s, "# raw markdown", Options{Runtime: runtimekind.StreamingStdio}); err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(string(s.last), "# raw") {
		t.Fatalf("sent raw markdown to streaming stdio: %q", s.last)
	}
}

type captureSender struct{ last []byte }

func (c *captureSender) SendInput(_ context.Context, data []byte) error {
	c.last = append([]byte(nil), data...)
	return nil
}
