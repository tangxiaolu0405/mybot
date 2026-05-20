package brain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxTurnUserRunes       = 600
	maxTurnAssistantRunes  = 2000
	maxShortTermFileBytes  = 96 * 1024
	shortTermFileHeader    = "# Current session (short-term)\n\n> Appended by cata after each chat turn. Consolidated into persona by autonomous evolution.\n\n"
)

// AppendChatTurn 在对话成功结束后写入当前 workspace 的 short-term。
func AppendChatTurn(userText, assistantText string) error {
	w, err := MustActive()
	if err != nil {
		return err
	}
	return AppendChatTurnFor(w, userText, assistantText)
}

// AppendChatTurnFor 向指定 workspace 追加回合。
func AppendChatTurnFor(w *Workspace, userText, assistantText string) error {
	userText = truncateRunes(strings.TrimSpace(userText), maxTurnUserRunes)
	assistantText = truncateRunes(strings.TrimSpace(assistantText), maxTurnAssistantRunes)
	if userText == "" && assistantText == "" {
		return nil
	}
	block := formatTurnBlock(userText, assistantText)
	return appendToShortTerm(w.ShortTermPath(), block)
}

// AppendSessionBoundary 在 chat_reset 时写入会话边界。
func AppendSessionBoundary() error {
	w, err := MustActive()
	if err != nil {
		return err
	}
	path := w.ShortTermPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	block := fmt.Sprintf("\n\n---\n\n**Session ended** `%s` (socket history cleared)\n\n", ts)
	return appendToShortTerm(path, block)
}

func formatTurnBlock(user, assistant string) string {
	ts := time.Now().UTC().Format(time.RFC3339)
	var b strings.Builder
	b.WriteString("\n\n## ")
	b.WriteString(ts)
	b.WriteString("\n\n")
	if user != "" {
		b.WriteString("**User:** ")
		b.WriteString(user)
		b.WriteString("\n\n")
	}
	if assistant != "" {
		b.WriteString("**Assistant:** ")
		b.WriteString(assistant)
		b.WriteString("\n\n")
	}
	return b.String()
}

func appendToShortTerm(path, block string) error {
	var existing []byte
	if data, err := os.ReadFile(path); err == nil {
		existing = data
	} else if !os.IsNotExist(err) {
		return err
	}

	body := existing
	if len(body) == 0 {
		body = []byte(shortTermFileHeader)
	} else if !strings.HasPrefix(string(body), "# ") {
		body = append([]byte(shortTermFileHeader), body...)
	}

	combined := append(body, []byte(block)...)
	if len(combined) > maxShortTermFileBytes {
		combined = trimShortTermToMax(combined, maxShortTermFileBytes)
	}
	return os.WriteFile(path, combined, 0644)
}

func trimShortTermToMax(data []byte, max int) []byte {
	if len(data) <= max {
		return data
	}
	marker := []byte("\n\n…(older turns trimmed)\n\n")
	keep := max - len(marker)
	if keep < len(shortTermFileHeader)+256 {
		keep = len(shortTermFileHeader) + 256
	}
	if keep > len(data) {
		keep = len(data)
	}
	out := make([]byte, 0, max)
	out = append(out, data[:len(shortTermFileHeader)]...)
	out = append(out, marker...)
	out = append(out, data[len(data)-keep+len(shortTermFileHeader):]...)
	if len(out) > max {
		out = out[len(out)-max:]
	}
	return out
}

func truncateRunes(s string, max int) string {
	if max <= 0 || s == "" {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "…"
}

// EnsureShortTermFileFor 确保 workspace short-term 存在。
func EnsureShortTermFileFor(w *Workspace) error {
	path := w.ShortTermPath()
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return os.WriteFile(path, []byte(shortTermFileHeader), 0644)
	}
	return err
}

// EnsureShortTermFile 活跃 workspace 或 legacy 路径。
func EnsureShortTermFile() error {
	if w := Active(); w != nil {
		return EnsureShortTermFileFor(w)
	}
	path := ShortTermCurrentPath()
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return os.WriteFile(path, []byte(shortTermFileHeader), 0644)
	}
	return err
}
