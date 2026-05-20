package brain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArchiveLogFileIfExists(t *testing.T) {
	dir := t.TempDir()
	oldHome := os.Getenv("CATA_HOME")
	os.Setenv("CATA_HOME", dir)
	defer os.Setenv("CATA_HOME", oldHome)

	path := filepath.Join(dir, "llm.log")
	if err := os.WriteFile(path, []byte("old\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := archiveLogFileIfExists(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("original should be renamed")
	}
	entries, _ := os.ReadDir(dir)
	var archived bool
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "llm.") && strings.HasSuffix(e.Name(), ".log") && e.Name() != "llm.log" {
			archived = true
		}
	}
	if !archived {
		t.Fatalf("expected archived llm.*.log, got %v", entries)
	}
}
