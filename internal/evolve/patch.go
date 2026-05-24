package evolve

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cata/internal/brain"
	"cata/internal/clock"
)

// DocUpdate LLM 决策输出的文档补丁（路径相对当前 workspace 根）。
type DocUpdate struct {
	Path    string `json:"path"`
	Mode    string `json:"mode"`
	Section string `json:"section,omitempty"`
	Content string `json:"content"`
}

// ApplyUpdates 将补丁写入当前活跃 workspace。
func ApplyUpdates(updates []DocUpdate) ([]string, error) {
	w, err := brain.MustActive()
	if err != nil {
		return nil, err
	}
	var touched []string
	for _, u := range updates {
		if strings.TrimSpace(u.Content) == "" && u.Mode != "write" && u.Mode != "overwrite" {
			continue
		}
		rel, err := normalizeWorkspaceRel(u.Path)
		if err != nil {
			return touched, err
		}
		abs := w.Path(rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
			return touched, err
		}
		switch strings.ToLower(strings.TrimSpace(u.Mode)) {
		case "write", "overwrite":
			if err := os.WriteFile(abs, []byte(u.Content), 0644); err != nil {
				return touched, err
			}
		case "append", "":
			f, err := os.OpenFile(abs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return touched, err
			}
			_, err = f.WriteString("\n\n" + strings.TrimSpace(u.Content) + "\n")
			_ = f.Close()
			if err != nil {
				return touched, err
			}
		case "append_section":
			body, _ := os.ReadFile(abs)
			section := strings.TrimSpace(u.Section)
			if section == "" {
				section = "Notes"
			}
			block := fmt.Sprintf("\n\n## %s\n\n%s\n", section, strings.TrimSpace(u.Content))
			if err := os.WriteFile(abs, append(body, []byte(block)...), 0644); err != nil {
				return touched, err
			}
		default:
			return touched, fmt.Errorf("unknown patch mode: %s", u.Mode)
		}
		touched = append(touched, rel)
	}
	return touched, nil
}

func normalizeWorkspaceRel(p string) (string, error) {
	p = strings.TrimSpace(p)
	p = strings.TrimPrefix(p, "brain/")
	p = filepath.ToSlash(filepath.Clean(p))
	if p == ".." || strings.HasPrefix(p, "../") {
		return "", fmt.Errorf("path not allowed: %s", p)
	}

	if p == brain.RelPersonaLocal {
		return p, nil
	}

	modePrefix := brain.DirModes + "/"
	if strings.HasPrefix(p, modePrefix) {
		rest := strings.TrimPrefix(p, modePrefix)
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) == 2 {
			file := parts[1]
			if file == brain.FilePersona || file == brain.FileBehavior || file == brain.FileConstraints || file == brain.FileCapabilities {
				return p, nil
			}
		}
	}

	if p == brain.RelMetaJSON {
		return p, nil
	}

	if p == brain.RelShortCurrent || strings.HasPrefix(p, "memory/short/") {
		return brain.RelShortCurrent, nil
	}
	if strings.HasPrefix(p, brain.RelMemoryLong+"/") && strings.HasSuffix(p, ".md") {
		return filepath.ToSlash(p), nil
	}
	if strings.HasPrefix(p, brain.RelMemoryArchive+"/") && strings.HasSuffix(p, ".md") {
		return p, nil
	}

	if strings.HasPrefix(p, brain.DirSkills+"/") {
		parts := strings.Split(p, "/")
		if len(parts) < 3 {
			return "", fmt.Errorf("invalid skill path: %s", p)
		}
		file := parts[2]
		if file == brain.FileSkillMD || file == brain.FileSkillManifest || isSkillScriptFile(file) {
			return p, nil
		}
		return "", fmt.Errorf("skill file not allowed: %s", file)
	}

	// legacy aliases
	legacy := map[string]string{
		brain.RelPathHot:              brain.DirModes + "/" + brain.ModeDefaultID + "/" + brain.FilePersona,
		brain.RelPathShortTermCurrent: brain.RelShortCurrent,
	}
	if mapped, ok := legacy[p]; ok {
		return mapped, nil
	}

	return "", fmt.Errorf("path not in evolution whitelist: %s", p)
}

func isSkillScriptFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".py", ".js", ".sh", ".ps1":
		return true
	default:
		return false
	}
}

// TouchArchiveDay 创建当日 archive 占位。
func TouchArchiveDay() (string, error) {
	w, err := brain.MustActive()
	if err != nil {
		return "", err
	}
	rel := filepath.Join(brain.RelMemoryArchive, clock.Format("2006-01-02")+".md")
	abs := w.Path(rel)
	if _, err := os.Stat(abs); err == nil {
		return rel, nil
	}
	content := fmt.Sprintf("# %s\n\n", clock.Format("2006-01-02"))
	if err := os.WriteFile(abs, []byte(content), 0644); err != nil {
		return "", err
	}
	return rel, nil
}
