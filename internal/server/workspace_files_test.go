package server

import (
	"os"
	"path/filepath"
	"testing"

	"mybot/internal/config"
)

func TestWorkspaceFileTools(t *testing.T) {
	dir := t.TempDir()
	config.Config = &config.AppConfig{
		Brain: config.BrainConfig{BaseDir: dir},
		WorkspaceFiles: config.WorkspaceFilesConfig{
			MaxReadBytes:  4096,
			MaxWriteBytes: 4096,
		},
	}
	config.BrainBaseDir = dir

	p := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(p, []byte("hello\nworld\n"), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := toolReadFile(`{"path":"a.txt","offset":2,"limit":1}`)
	if err != nil || out == "" {
		t.Fatalf("read: %v %q", err, out)
	}
	if !contains(out, "world") {
		t.Fatalf("read offset: %q", out)
	}

	out, err = toolSearchReplace(`{"path":"a.txt","old_string":"hello","new_string":"hi"}`)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(p)
	if string(data) != "hi\nworld\n" {
		t.Fatalf("replace: %q", data)
	}

	out, err = toolAppendFile(`{"path":"a.txt","content":"!"}`)
	if err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(p)
	if string(data) != "hi\nworld\n!" {
		t.Fatalf("append: %q", data)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
