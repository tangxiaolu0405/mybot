package client

import (
	"fmt"
	"os"
	"strings"
)

// ANSI escape sequences.
const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiGray   = "\033[90m"
)

// --- Output to stdout (AI text, the "deliverable") ---

func out(s string)               { fmt.Print(s) }
func outToken(s string)          { out(s) }

// --- Output to stderr (meta: tools, progress, errors) ---

func meta(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// toolStart renders a tool invocation header.
func toolStart(name string, display string) {
	switch display {
	case "silent":
		meta("  %s●%s %s\n", ansiDim, ansiReset, name)
	default:
		meta("  %s✦%s %s%s%s\n", ansiCyan, ansiReset, ansiBold, name, ansiReset)
	}
}

// toolOutput renders tool result content.
func toolOutput(_ /*name*/, output, display string) {
	if strings.TrimSpace(output) == "" {
		return
	}
	switch display {
	case "silent":
		return
	case "verbose":
		meta("%s\n", strings.TrimRight(output, "\n"))
	default:
		trunc := truncate(output, 2000)
		meta("%s%s%s\n", ansiDim, strings.TrimRight(trunc, "\n"), ansiReset)
	}
}

// runCmdResult renders a run_command result block.
func runCmdResult(cmd, _ /*cwd*/, output string) {
	meta("\n")
	meta("  %s$ %s%s\n", ansiGreen, ansiReset, cmd)
	if strings.TrimSpace(output) != "" {
		for _, line := range strings.Split(strings.TrimRight(output, "\n"), "\n") {
			meta("  %s│%s %s\n", ansiDim, ansiReset, line)
		}
	}
	meta("\n")
}

// progressMsg renders a progress/status message.
func progressMsg(msg string) {
	meta("  %s⟳%s %s%s%s\n", ansiDim, ansiReset, ansiDim, msg, ansiReset)
}

// errorMsg renders an error.
func errorMsg(msg string) {
	meta("  %s✖%s %s%s%s\n", ansiRed, ansiReset, ansiRed, msg, ansiReset)
}

// confirmPrompt shows the exec confirmation UI and reads the user's choice.
func confirmPrompt(_ /*id*/, cmd, cwd string) (bool, error) {
	opts := []SelectOption{
		{ID: "run", Label: "Run"},
		{ID: "cancel", Label: "Cancel"},
	}
	id, err := Select("⚙ "+cmd, "cwd: "+cwd, opts)
	if err != nil || id == "" || id == "cancel" {
		return false, err
	}
	return true, nil
}

// execDenied renders a cancelled command.
func execDenied() {
	meta("  %s── cancelled%s\n", ansiDim, ansiReset)
}

// execDone renders an exec_done status line.
func execDone(cmd string, exitCode int, timedOut bool) {
	if timedOut {
		meta("  %s⏱ timeout%s  %s\n", ansiYellow, ansiReset, cmd)
	} else if exitCode != 0 {
		meta("  %s✖ exit %d%s  %s\n", ansiRed, exitCode, ansiReset, cmd)
	} else {
		meta("  %s✓%s  %s\n", ansiGreen, ansiReset, cmd)
	}
}

// fileWritten renders a file write confirmation.
func fileWritten(path string, bytes int) {
	meta("  %s✎%s wrote %s%s%s (%d bytes)\n", ansiGreen, ansiReset, ansiYellow, path, ansiReset, bytes)
}

// diffLine renders a single line of a diff.
func diffLine(content string) {
	content = strings.TrimRight(content, "\n\r")
	if strings.HasPrefix(content, "+") {
		meta("  %s%s%s\n", ansiGreen, content, ansiReset)
	} else if strings.HasPrefix(content, "-") {
		meta("  %s%s%s\n", ansiRed, content, ansiReset)
	} else {
		meta("  %s%s%s\n", ansiDim, content, ansiReset)
	}
}

// welcome prints the startup message.
func welcome() {
	meta("%s── cata ──────────────────────────────%s\n", ansiDim, ansiReset)
	meta("  /clear  reset  ·  /exit  quit  ·  %s\"\"\"%s multiline\n", ansiDim, ansiReset)
	meta("%s──────────────────────────────────────%s\n", ansiDim, ansiReset)
}
