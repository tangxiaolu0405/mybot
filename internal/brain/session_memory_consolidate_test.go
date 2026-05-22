package brain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFinalizeShortTermAfterConsolidate(t *testing.T) {
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
	body := shortTermFileHeader + strings.Repeat("x", 3000) + "\n\n## 2026-01-01T00:00:00Z\n\n**User:** hi\n\n"
	if err := os.WriteFile(ws.ShortTermPath(), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	rel, err := FinalizeShortTermAfterConsolidate(512)
	if err != nil {
		t.Fatal(err)
	}
	if rel == "" {
		t.Fatal("expected archive path")
	}
	arch, err := os.ReadFile(ws.Path(rel))
	if err != nil || !strings.Contains(string(arch), "Short-term archive") {
		t.Fatalf("archive: %v %q", err, arch)
	}
	cur, err := os.ReadFile(ws.ShortTermPath())
	if err != nil {
		t.Fatal(err)
	}
	if len(cur) > 1500 {
		t.Fatalf("short-term should shrink, got %d bytes", len(cur))
	}
	if !strings.Contains(string(cur), "Last consolidated") {
		t.Fatalf("missing marker: %s", cur)
	}
}
