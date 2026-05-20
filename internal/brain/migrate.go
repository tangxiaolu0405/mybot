package brain

import (
	"os"
	"path/filepath"
	"time"

	"mybot/internal/config"
)

// MigrateLegacyBrain 将旧版扁平 ~/.cata/brain/{hot,core,memory,...} 迁入首个 workspace。
func MigrateLegacyBrain() error {
	entries, err := ListRegistryEntries()
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return nil
	}

	legacyHot := filepath.Join(brainRoot(), RelPathHot)
	legacyShort := filepath.Join(brainRoot(), RelPathShortTermCurrent)
	if _, err := os.Stat(legacyHot); os.IsNotExist(err) {
		if _, err2 := os.Stat(legacyShort); os.IsNotExist(err2) {
			return nil
		}
	}

	root := fallbackWorkspaceRoot()
	root, _ = filepath.Abs(root)
	kind := KindEphemeral
	if findGitRoot(root) == root {
		kind = KindGit
	} else if findMarkedRoot(root) == root {
		kind = KindMarked
	}
	ws := &Workspace{
		ID:         workspaceID(root),
		RootPath:   root,
		Kind:       kind,
		ActiveMode: ModeDefaultID,
	}
	if err := ws.EnsureScaffold(); err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if err := upsertRegistryEntry(RegistryEntry{
		ID: ws.ID, RootPath: ws.RootPath, Kind: ws.Kind,
		CreatedAt: now, LastSeenAt: now, ActiveMode: ws.ActiveMode,
	}); err != nil {
		return err
	}

	_ = copyIfExists(legacyHot, ws.PersonaPath())
	_ = copyIfExists(legacyShort, ws.ShortTermPath())
	_ = copyTreeIfExists(filepath.Join(brainRoot(), RelPathLongTerm), ws.LongTermDir())
	_ = copyIfExists(filepath.Join(brainRoot(), RelPathEvolutionLog), ws.EvolutionLogPath())
	_ = copyTreeIfExists(filepath.Join(brainRoot(), RelPathArchive), ws.ArchiveDir())
	return nil
}

func fallbackWorkspaceRoot() string {
	_ = config.InitBrainPath()
	if b := config.GetBrainBaseDir(); b != "" {
		if gitRoot := findGitRoot(b); gitRoot != "" {
			return gitRoot
		}
		return b
	}
	if r := findGitRoot(mustGetwd()); r != "" {
		return r
	}
	return mustGetwd()
}

func mustGetwd() string {
	wd, _ := os.Getwd()
	if wd == "" {
		return CataHome()
	}
	return wd
}

func copyIfExists(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return nil
	}
	_ = os.MkdirAll(filepath.Dir(dst), 0755)
	return os.WriteFile(dst, data, 0644)
}

func copyTreeIfExists(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		return copyIfExists(path, filepath.Join(dst, rel))
	})
}
