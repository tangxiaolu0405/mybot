package llm

import (
	"encoding/json"
	"testing"
)

func TestParseToolArguments_looseJSON(t *testing.T) {
	raw := "{\"argv\": \n[\"mkdir\", \"-p\", \"/tmp/foo\"]\n\n}"
	var p struct {
		Argv []string `json:"argv"`
	}
	if err := ParseToolArguments(raw, &p); err != nil {
		t.Fatal(err)
	}
	if len(p.Argv) != 3 || p.Argv[0] != "mkdir" {
		t.Fatalf("argv=%v", p.Argv)
	}
}

func TestNormalizeToolArguments_truncatedAppendFile(t *testing.T) {
	raw := `{"path": "zhangtingban_analysis/lianban/连板票分析报告.md", "content": "# 标题\n部分正文`
	got := NormalizeToolArguments("append_file", raw)
	if got == "" {
		t.Fatal("expected repair")
	}
	var m struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := ParseToolArguments(got, &m); err != nil {
		t.Fatal(err)
	}
	if m.Path != "zhangtingban_analysis/lianban/连板票分析报告.md" {
		t.Fatalf("path=%q", m.Path)
	}
	if m.Content == "" {
		t.Fatal("expected partial content")
	}
}

func TestParseEmbeddedToolCalls_bracket(t *testing.T) {
	content := `说明文字\n\n[tool_call append_file] {"path": "a/b.csv", "content": "x,y\n"}`
	calls, stripped := ParseEmbeddedToolCalls(content)
	if len(calls) != 1 {
		t.Fatalf("calls=%d", len(calls))
	}
	if calls[0].Function.Name != "append_file" {
		t.Fatalf("name=%s", calls[0].Function.Name)
	}
	if stripped == content {
		t.Fatal("expected stripped content")
	}
	var m struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := ParseToolArguments(calls[0].Function.Arguments, &m); err != nil {
		t.Fatal(err)
	}
	if m.Path != "a/b.csv" || m.Content != "x,y\n" {
		t.Fatalf("path=%q content=%q", m.Path, m.Content)
	}
}

func TestSanitizeMessagesToolCalls_invalidHistory(t *testing.T) {
	msgs := []Message{{
		Role: "assistant",
		ToolCalls: []ToolCall{{
			ID: "c1",
			Function: ToolCallFunction{
				Name:      "append_file",
				Arguments: `{"path":"out.md","content":"hello`,
			},
		}},
	}}
	out := SanitizeMessagesToolCalls(msgs)
	if len(out[0].ToolCalls) != 1 {
		t.Fatalf("tool_calls=%d", len(out[0].ToolCalls))
	}
	if !json.Valid([]byte(out[0].ToolCalls[0].Function.Arguments)) {
		t.Fatalf("args=%q", out[0].ToolCalls[0].Function.Arguments)
	}
}
