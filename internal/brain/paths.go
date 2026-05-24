// Package brain 提供与 brain/core.md 一致的核心路径与资源定义，确保代码与文档单一数据源。
package brain

import (
	"os"
	"path/filepath"

	"cata/internal/config"
)

// 仓库模板 brain/ 内相对路径（种子来源，与 ~/.cata/global/ 名称一致）。
const (
	RelPathBootAssembler = "boot-assembler.md"
	RelPathConstraints   = "constraints.md"
	RelPathBehavior      = "behavior.md"

	RelPathMemoryIndex   = "memory_index.json"
	RelPathEvolutionLog  = "evolution_log.json"
)

// legacy paths kept for migration / backward-compatible fallbacks
const (
	RelPathHot        = "hot.md"                              // migrated → modes/_default/persona.md
	RelPathShortTermCurrent = "memory/short-term/current_session.md" // migrated → memory/short/current.md
	RelPathLongTerm   = "memory/long-term"                    // migrated → memory/long/
	RelPathArchive    = "archive"                             // migrated → memory/archive/
	RelPathLearnSystemPrompt = "context/learn_system_prompt.md"
)

// Dir 返回 brain 根（~/.cata/brain），兼容旧代码。
func Dir() string {
	return brainRoot()
}

// BaseDir 返回项目根（用于 run_command 等）。
func BaseDir() string {
	return config.GetBrainBaseDir()
}

// Path 在活跃 workspace 下解析相对路径；无 workspace 时回退到 legacy brain 根。
func Path(rel string) string {
	if w := Active(); w != nil {
		return w.Path(rel)
	}
	return filepath.Join(brainRoot(), rel)
}

// BootLeaderPath 全局启动组装说明。
func BootLeaderPath() string {
	p := filepath.Join(globalDir(), FileGlobalBoot)
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return filepath.Join(brainRoot(), RelPathBootAssembler)
}

// HotPath 当前 workspace 活跃 mode 的 persona。
func HotPath() string {
	if w := Active(); w != nil {
		return w.PersonaPath()
	}
	return filepath.Join(brainRoot(), RelPathHot)
}

// ArchiveDir 当前 workspace 的 archive。
func ArchiveDir() string {
	if w := Active(); w != nil {
		return w.ArchiveDir()
	}
	return filepath.Join(brainRoot(), RelPathArchive)
}

// EvolutionLogPath 当前 workspace 的演进日志。
func EvolutionLogPath() string {
	if w := Active(); w != nil {
		return w.EvolutionLogPath()
	}
	return filepath.Join(brainRoot(), RelPathEvolutionLog)
}

// ShortTermCurrentPath 当前 workspace 短期记忆。
func ShortTermCurrentPath() string {
	if w := Active(); w != nil {
		return w.ShortTermPath()
	}
	return filepath.Join(brainRoot(), RelPathShortTermCurrent)
}

// LearnSystemPromptPath legacy。
func LearnSystemPromptPath() string {
	return filepath.Join(brainRoot(), RelPathLearnSystemPrompt)
}

// LongTermDir 当前 workspace 长期记忆目录。
func LongTermDir() string {
	if w := Active(); w != nil {
		return w.LongTermDir()
	}
	return filepath.Join(brainRoot(), RelPathLongTerm)
}

// GlobalConstraintsPath 全局约束。
func GlobalConstraintsPath() string {
	return filepath.Join(globalDir(), FileGlobalConstraints)
}

// GlobalBehaviorPath 全局行为 SOP。
func GlobalBehaviorPath() string {
	return filepath.Join(globalDir(), FileGlobalBehavior)
}

// PersonaLocalPath 当前 workspace 项目说明。
func PersonaLocalPath() string {
	if w := Active(); w != nil {
		return w.PersonaLocalPath()
	}
	return filepath.Join(brainRoot(), RelPersonaLocal)
}

// ArchiveSummaryFilename 返回周期摘要文件名。
func ArchiveSummaryFilename(yearMonth string) string {
	return "summary-" + yearMonth + ".md"
}

const ArchiveBackupDir = "backup"
