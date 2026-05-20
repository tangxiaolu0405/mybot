package brain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendChatTurn_workspace(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CATA_HOME", home)
	unsetBrainDir(t)

	if err := EnsureCataLayout(); err != nil {
		t.Fatal(err)
	}
	proj := filepath.Join(home, "proj")
	if err := os.MkdirAll(proj, 0755); err != nil {
		t.Fatal(err)
	}
	ws, err := ResolveWorkspace(proj)
	if err != nil {
		t.Fatal(err)
	}
	if err := AppendChatTurn("hello", "world"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(ws.ShortTermPath())
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "**User:** hello") || !strings.Contains(s, "**Assistant:** world") {
		t.Fatalf("unexpected: %s", s)
	}
}

func unsetBrainDir(t *testing.T) {
	t.Helper()
	os.Unsetenv("CATA_BRAIN_DIR")
}
