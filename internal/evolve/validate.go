package evolve

import (
	"path/filepath"
	"strings"
	"unicode/utf8"

	"mybot/internal/brain"
)

func filterUpdates(updates []DocUpdate) []DocUpdate {
	var out []DocUpdate
	for _, u := range updates {
		content := strings.TrimSpace(u.Content)
		if u.Mode != "write" && u.Mode != "overwrite" && utf8.RuneCountInString(content) < minPatchContentRunes {
			continue
		}
		rel := strings.TrimPrefix(strings.TrimSpace(u.Path), "brain/")
		rel = filepath.ToSlash(filepath.Clean(rel))
		if strings.Contains(rel, brain.DirModes+"/") {
			if strings.HasSuffix(rel, "/"+brain.FilePersona) || strings.HasSuffix(rel, "/"+brain.FileBehavior) {
				// allow append to persona/behavior
			}
		}
		if strings.HasSuffix(rel, "/"+brain.FileConstraints) &&
			(strings.EqualFold(u.Mode, "write") || strings.EqualFold(u.Mode, "overwrite")) &&
			utf8.RuneCountInString(content) > 600 {
			continue
		}
		if (rel == brain.RelPersonaLocal || strings.Contains(rel, "/"+brain.FilePersona)) &&
			(strings.EqualFold(u.Mode, "write") || strings.EqualFold(u.Mode, "overwrite")) &&
			utf8.RuneCountInString(content) > 2000 {
			continue
		}
		out = append(out, u)
	}
	if len(out) > maxUpdatesPerCycle {
		out = out[:maxUpdatesPerCycle]
	}
	return out
}

func isMeaningfulDecision(dec *Decision, touched []string) bool {
	action := strings.ToLower(strings.TrimSpace(dec.Action))
	if action == "" || action == "idle" {
		return len(touched) > 0
	}
	return len(touched) > 0 || utf8.RuneCountInString(strings.TrimSpace(dec.Learning)) >= minPatchContentRunes
}
