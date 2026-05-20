package brain

import (
	"fmt"
	"os"
	"strings"
)

// TerminalBundleSystemPrefix cata 终端脑子节选 system 消息前缀。
const TerminalBundleSystemPrefix = "【Cata 脑子节选"

// CompactExcessiveNewlines 将连续 3 个及以上的换行压成至多 2 个。
func CompactExcessiveNewlines(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	var b strings.Builder
	b.Grow(len(s))
	consecNL := 0
	for _, r := range s {
		if r == '\n' {
			consecNL++
			if consecNL <= 2 {
				b.WriteByte('\n')
			}
		} else {
			consecNL = 0
			b.WriteRune(r)
		}
	}
	return b.String()
}

// TerminalBrainSystemExtension 注入 global 约束/行为 + mode persona + persona.local。
func TerminalBrainSystemExtension(maxPerFile, maxTotal int) string {
	if maxPerFile <= 0 {
		maxPerFile = 6500
	}
	if maxTotal <= 0 {
		maxTotal = 20000
	}

	sections := []struct {
		title string
		path  string
	}{}

	if p := GlobalConstraintsPath(); fileExists(p) {
		sections = append(sections, struct{ title, path string }{"global/constraints", p})
	}
	if p := GlobalBehaviorPath(); fileExists(p) {
		sections = append(sections, struct{ title, path string }{"global/behavior", p})
	}
	if w := Active(); w != nil {
		if p := w.PersonaPath(); fileExists(p) {
			sections = append(sections, struct{ title, path string }{
				fmt.Sprintf("mode/%s/persona", w.modeID()), p,
			})
		}
		if p := w.PersonaLocalPath(); fileExists(p) {
			sections = append(sections, struct{ title, path string }{"brain/persona.local (focus)", p})
		}
	} else {
		// legacy fallback
		for _, rel := range []struct{ title, rel string }{
			{"brain/core.md", RelPathCore},
			{"brain/workflow.md", RelPathWorkflow},
			{"brain/hot.md", RelPathHot},
		} {
			p := filepathJoinBrain(rel.rel)
			if fileExists(p) {
				sections = append(sections, struct{ title, path string }{rel.title, p})
			}
		}
	}

	var blocks []string
	used := 0
	for _, sec := range sections {
		block := readSection(sec.title, sec.path, maxPerFile)
		if used+len(block) > maxTotal {
			blocks = append(blocks, "## (省略)\n后续 brain 摘录因 maxTotal 上限未载入。")
			break
		}
		blocks = append(blocks, block)
		used += len(block)
	}
	body := strings.Join(blocks, "\n\n")
	paths := TerminalPathsSystemBlock()
	if strings.TrimSpace(body) == "" {
		return paths
	}
	return paths + "\n\n" + TerminalBundleSystemPrefix + "（global + mode persona；均在 ~/.cata，非产出区）】\n\n" + body
}

func readSection(title, path string, maxPerFile int) string {
	b, err := os.ReadFile(path)
	var block string
	if err != nil {
		block = fmt.Sprintf("## %s\n(未能读取 %s: %v)", title, path, err)
	} else {
		body := CompactExcessiveNewlines(string(b))
		if len(body) > maxPerFile {
			body = body[:maxPerFile] + "\n…(truncated)"
		}
		block = "## " + title + "\n" + body
	}
	return block
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func filepathJoinBrain(rel string) string {
	return Path(rel)
}
