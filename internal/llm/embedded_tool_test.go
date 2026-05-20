package llm

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseEmbeddedToolCalls_MiniMax(t *testing.T) {
	content := `说明文字
<tool_call>
<tool name="run_command">
<param name="command">cd "D:\software\mybot" &amp;&amp; mkdir xiaohongshu</param>
</tool>
</tool_call>`
	calls, stripped := ParseEmbeddedToolCalls(content)
	if len(calls) != 1 {
		t.Fatalf("calls=%d", len(calls))
	}
	if calls[0].Function.Name != "run_command" {
		t.Fatalf("name=%q", calls[0].Function.Name)
	}
	var p struct {
		Argv []string `json:"argv"`
	}
	if err := json.Unmarshal([]byte(calls[0].Function.Arguments), &p); err != nil {
		t.Fatal(err)
	}
	if len(p.Argv) < 2 || p.Argv[0] != "cmd.exe" {
		t.Fatalf("argv=%v", p.Argv)
	}
	if !strings.Contains(p.Argv[len(p.Argv)-1], "mkdir xiaohongshu") {
		t.Fatalf("argv=%v", p.Argv)
	}
	if strings.Contains(stripped, "<tool_call>") {
		t.Fatalf("stripped still has tool_call: %q", stripped)
	}
}

func TestParseEmbeddedToolCalls_LooseTool(t *testing.T) {
	content := `<tool name="run_command"><param name="command">cd D:\a &amp;&amp; dir</param></tool>`
	calls, _ := ParseEmbeddedToolCalls(content)
	if len(calls) != 1 {
		t.Fatalf("calls=%d", len(calls))
	}
}
