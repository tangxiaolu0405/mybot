package brain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkillSearchPaths_workspaceBeforeGlobal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CATA_HOME", home)
	ws := &Workspace{ID: "ws_skillprompt", RootPath: home, ActiveMode: ModeDefaultID}
	SetActive(ws)

	skillID := "test-skill"
	wsDir := filepath.Join(ws.SkillDir(skillID))
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wsDir, FileSkillMD), []byte("from workspace brain"), 0644); err != nil {
		t.Fatal(err)
	}
	globalDir := filepath.Join(CataHome(), DirSkills, skillID)
	if err := os.MkdirAll(globalDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(globalDir, FileSkillMD), []byte("from global"), 0644); err != nil {
		t.Fatal(err)
	}

	body, from, err := loadSkillMarkdown(skillID)
	if err != nil {
		t.Fatal(err)
	}
	if body != "from workspace brain" {
		t.Fatalf("body=%q", body)
	}
	if from != ws.SkillMarkdownPath(skillID) {
		t.Fatalf("from=%q want %q", from, ws.SkillMarkdownPath(skillID))
	}
}
