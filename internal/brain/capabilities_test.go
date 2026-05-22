package brain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendSkillToCapabilities_dedup(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CATA_HOME", home)
	ws := &Workspace{ID: "ws_testdedup", RootPath: home, ActiveMode: ModeDefaultID}
	modeDir := filepath.Join(home, DirBrain, DirWorkspaces, ws.ID, DirModes, ModeDefaultID)
	if err := os.MkdirAll(modeDir, 0755); err != nil {
		t.Fatal(err)
	}
	capPath := filepath.Join(modeDir, FileCapabilities)
	if err := os.WriteFile(capPath, []byte("skills: []\nmcp:\n  - browser\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := AppendSkillToCapabilities(ws, "foo-skill"); err != nil {
		t.Fatal(err)
	}
	if err := AppendSkillToCapabilities(ws, "foo-skill"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(capPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(data), "foo-skill") != 1 {
		t.Fatalf("expected one foo-skill entry: %s", data)
	}
}

func TestParseCapabilitiesYAML(t *testing.T) {
	c := ParseCapabilitiesYAML([]byte("skills:\n  - babysit\nmcp:\n  - browser\n"))
	if len(c.Skills) != 1 || c.Skills[0] != "babysit" {
		t.Fatalf("skills=%v", c.Skills)
	}
	if !c.AllowsMCPServer("browser") {
		t.Fatal("expected browser allowed")
	}
}
