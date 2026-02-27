// Package brain 提供与 brain/core.md 一致的核心路径与资源定义，确保代码与文档单一数据源。
package brain

import (
	"path/filepath"

	"mybot/internal/config"
)

// 以下常量与 brain/core.md「资源路径」表一一对应，修改时请同步更新 core.md。

const (
	// RelPathBootLeader 系统初始化引导
	RelPathBootLeader = "boot-leader.md"
	// RelPathCore 核心思维
	RelPathCore = "core.md"
	// RelPathWorkflow 自主演进流程
	RelPathWorkflow = "workflow.md"
	// RelPathHot 热记忆
	RelPathHot = "hot.md"
	// RelPathMemoryIndex 记忆索引
	RelPathMemoryIndex = "memory_index.json"
	// RelPathEvolutionLog 演进日志
	RelPathEvolutionLog = "evolution_log.json"
	// RelPathTaskQueue 任务队列
	RelPathTaskQueue = "task_queue.json"
	// RelPathCapabilities 能力记录（workflow 可选）
	RelPathCapabilities = "capabilities.json"
	// RelPathArchive archive 目录
	RelPathArchive = "archive"
	// RelPathShortTermCurrent 短期记忆当前会话
	RelPathShortTermCurrent = "memory/short-term/current_session.md"
	// RelPathLongTerm 长期记忆目录
	RelPathLongTerm = "memory/long-term"
	// RelPathHeartbeatState 心跳状态
	RelPathHeartbeatState = "memory/short-term/heartbeat-state.json"
)

// 技能相关路径在项目根（BaseDir），与 core.md「技能目录 / 技能索引」一致。
const (
	// RelPathSkillsDir 技能目录（相对于 BaseDir）
	RelPathSkillsDir = "skills"
	// RelPathSkillsIndex 技能索引（相对于 BaseDir）
	RelPathSkillsIndex = "skills/skills-index.json"
)

// Dir 返回 brain 目录绝对路径。
func Dir() string {
	return config.GetBrainDir()
}

// BaseDir 返回项目根目录（用于 skills、git 等）。
func BaseDir() string {
	return config.GetBrainBaseDir()
}

// Path 返回 brain 内相对路径的绝对路径。
func Path(rel string) string {
	return filepath.Join(Dir(), rel)
}

// BootLeaderPath 返回 boot-leader.md 绝对路径。
func BootLeaderPath() string {
	return Path(RelPathBootLeader)
}

// HotPath 热记忆文件路径。
func HotPath() string {
	return Path(RelPathHot)
}

// ArchiveDir 返回 archive 目录绝对路径。
func ArchiveDir() string {
	return Path(RelPathArchive)
}

// MemoryIndexPath 记忆索引文件路径。
func MemoryIndexPath() string {
	return Path(RelPathMemoryIndex)
}

// EvolutionLogPath 演进日志文件路径。
func EvolutionLogPath() string {
	return Path(RelPathEvolutionLog)
}

// TaskQueuePath 任务队列文件路径。
func TaskQueuePath() string {
	return Path(RelPathTaskQueue)
}

// CapabilitiesPath 能力记录文件路径。
func CapabilitiesPath() string {
	return Path(RelPathCapabilities)
}

// ShortTermCurrentPath 短期记忆当前会话文件路径。
func ShortTermCurrentPath() string {
	return Path(RelPathShortTermCurrent)
}

// LongTermDir 长期记忆目录路径。
func LongTermDir() string {
	return Path(RelPathLongTerm)
}

// SkillsDir 返回技能目录绝对路径（位于项目根下）。
func SkillsDir() string {
	return filepath.Join(BaseDir(), RelPathSkillsDir)
}

// SkillsIndexPath 返回 skills-index.json 绝对路径。
func SkillsIndexPath() string {
	return filepath.Join(BaseDir(), RelPathSkillsIndex)
}

// ArchiveSummaryFilename 返回周期摘要文件名格式 summary-YYYY-MM.md。
func ArchiveSummaryFilename(yearMonth string) string {
	return "summary-" + yearMonth + ".md"
}

// ArchiveBackupDir archive 备份子目录名。
const ArchiveBackupDir = "backup"
