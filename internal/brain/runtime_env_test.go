package brain

import (
	"strings"
	"testing"
)

func TestTerminalPathsIncludesRuntime(t *testing.T) {
	SetRuntimeEnv(&RuntimeEnv{
		OS: "windows", HostOS: "windows", Arch: "amd64", Shell: "cmd", Terminal: "Windows Terminal",
	})
	defer SetRuntimeEnv(nil)
	block := TerminalPathsSystemBlock()
	if !strings.Contains(block, "windows") || !strings.Contains(block, "cmd") {
		t.Fatalf("missing runtime in block: %s", block)
	}
}
