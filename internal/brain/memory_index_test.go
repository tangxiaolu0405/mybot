package brain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncMemoryIndexAfterEvolution(t *testing.T) {
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
	rel := filepath.ToSlash(filepath.Join(DirModes, ws.modeID(), FilePersona))
	persona := ws.Path(rel)
	if err := os.WriteFile(persona, []byte("# Persona\n\nUser prefers Go and terminal agents.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := SyncMemoryIndexAfterEvolution([]string{rel}, "Learned user likes minimal diffs.", ""); err != nil {
		t.Fatal(err)
	}
	idx, err := LoadMemoryIndex()
	if err != nil {
		t.Fatal(err)
	}
	if len(idx.Entries) < 2 {
		t.Fatalf("expected persona + learning entries, got %d", len(idx.Entries))
	}
	data, _ := os.ReadFile(ws.MemoryIndexPath())
	if !strings.Contains(string(data), `"version"`) {
		t.Fatalf("index should be object: %s", data)
	}
	var parsed MemoryIndex
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}
	block := MemoryIndexPromptBlock(4000)
	if !strings.Contains(block, "【Cata 记忆索引】") {
		t.Fatalf("prompt block: %q", block)
	}
}
