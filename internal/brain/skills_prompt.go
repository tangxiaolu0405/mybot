package brain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mybot/internal/config"
)

const (
	skillFileName       = "SKILL.md"
	maxBytesPerSkill    = 8000
	maxBytesSkillsTotal = 16000
	SkillsSystemPrefix  = "【Cata Skills】"
)

// SkillsPromptBlock 将 capabilities 中的 skill 名对应的 SKILL.md 拼成 system 段。
func SkillsPromptBlock(skillNames []string) string {
	if len(skillNames) == 0 {
		return ""
	}
	var blocks []string
	used := 0
	for _, name := range skillNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		body, path, err := loadSkillMarkdown(name)
		if err != nil {
			blocks = append(blocks, fmt.Sprintf("## skill:%s\n(load failed: %v)", name, err))
			continue
		}
		if len(body) > maxBytesPerSkill {
			body = body[:maxBytesPerSkill] + "\n…(truncated)"
		}
		block := fmt.Sprintf("## skill:%s\n(from %s)\n\n%s", name, path, body)
		if used+len(block) > maxBytesSkillsTotal {
			blocks = append(blocks, "## (skills omitted)\n后续 skill 因体积上限未载入。")
			break
		}
		blocks = append(blocks, block)
		used += len(block)
	}
	if len(blocks) == 0 {
		return ""
	}
	return SkillsSystemPrefix + "\n\n" + strings.Join(blocks, "\n\n")
}

func loadSkillMarkdown(name string) (body, from string, err error) {
	for _, p := range skillSearchPaths(name) {
		data, e := os.ReadFile(p)
		if e == nil {
			return CompactExcessiveNewlines(strings.TrimSpace(string(data))), p, nil
		}
		if !os.IsNotExist(e) {
			err = e
		}
	}
	if err != nil {
		return "", "", err
	}
	return "", "", fmt.Errorf("SKILL.md not found for %q", name)
}

func skillSearchPaths(name string) []string {
	var paths []string
	home := CataHome()
	paths = append(paths, filepath.Join(home, "skills", name, skillFileName))
	if root := config.FindProjectRoot(); root != "" {
		paths = append(paths, filepath.Join(root, "skills", name, skillFileName))
	}
	if base := config.GetBrainBaseDir(); base != "" {
		paths = append(paths, filepath.Join(base, "skills", name, skillFileName))
	}
	if h, e := os.UserHomeDir(); e == nil && h != "" {
		paths = append(paths, filepath.Join(h, ".cursor", "skills-cursor", name, skillFileName))
	}
	return paths
}
