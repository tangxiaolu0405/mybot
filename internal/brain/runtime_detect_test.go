package brain

import (
	"strings"
	"testing"
)

func TestWSLPathForOutput(t *testing.T) {
	got := WSLPathForOutput(`D:\software\mybot`)
	if got != "/mnt/d/software/mybot" {
		t.Fatalf("got %q", got)
	}
}

func TestShellLineToArgv_WSL(t *testing.T) {
	SetRuntimeEnv(&RuntimeEnv{OS: "linux", HostOS: "windows", Shell: "bash", Terminal: "wsl:Ubuntu"})
	defer SetRuntimeEnv(nil)
	argv := ShellLineToArgv("mkdir -p foo")
	if len(argv) < 4 || argv[0] != "wsl.exe" || argv[3] != "-lc" {
		t.Fatalf("argv=%v", argv)
	}
}

func TestRunCommandHints_WSL_NoPowerShell(t *testing.T) {
	e := RuntimeEnv{OS: "linux", HostOS: "windows", Shell: "bash", Terminal: "wsl:Ubuntu"}
	h := e.runCommandHints()
	if strings.Contains(strings.ToLower(h), "powershell") && !strings.Contains(h, "禁止") {
		t.Fatalf("hints should not promote powershell: %s", h)
	}
	if !strings.Contains(h, "mkdir -p") {
		t.Fatalf("missing bash hints: %s", h)
	}
}
