package evolution

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mybot/internal/brain"
	"mybot/internal/memory"
)

var (
	EvolutionLogFilePath string
	CapabilitiesFilePath string
)

// 向后兼容
var (
	EvolutionLogFile = &EvolutionLogFilePath
	CapabilitiesFile = &CapabilitiesFilePath
)

// initEvolutionPaths 与 brain/core.md、workflow.md 一致
func initEvolutionPaths() {
	EvolutionLogFilePath = brain.EvolutionLogPath()
	CapabilitiesFilePath = brain.CapabilitiesPath()
}

func init() {
	initEvolutionPaths()
}

// SystemState 系统状态
type SystemState struct {
	MemoryState    MemoryState    `json:"memory_state"`
	TaskState      TaskState      `json:"task_state"`
	EvolutionState EvolutionState `json:"evolution_state"`
	Timestamp      string         `json:"timestamp"`
}

// MemoryState 记忆状态
type MemoryState struct {
	ArchiveFileCount int    `json:"archive_file_count"`
	ArchiveTotalSize int64  `json:"archive_total_size"`
	HotLastUpdated   string `json:"hot_last_updated"`
	IndexEntryCount  int    `json:"index_entry_count"`
	IndexComplete    bool   `json:"index_complete"`
	NeedsSummarize   bool   `json:"needs_summarize"`
	SummarizeReason  string `json:"summarize_reason"`
}

// TaskState 任务状态
type TaskState struct {
	RecentTasks      []TaskHistory `json:"recent_tasks"`
	SuccessRate      float64       `json:"success_rate"`
	PendingTasks     int           `json:"pending_tasks"`
	LastTaskTime     string        `json:"last_task_time"`
}

// TaskHistory 任务历史
type TaskHistory struct {
	TaskID    string `json:"task_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Completed string `json:"completed"`
	Result    string `json:"result"`
}

// EvolutionState 演进状态
type EvolutionState struct {
	Capabilities     []Capability `json:"capabilities"`
	LearningProgress map[string]float64 `json:"learning_progress"`
	LastEvolution    string       `json:"last_evolution"`
	ImprovementAreas []string     `json:"improvement_areas"`
}

// Capability 能力
type Capability struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      string   `json:"status"` // "learning" | "mastered" | "failed"
	Progress    float64  `json:"progress"` // 0.0 - 1.0
}

// StateAnalyzer 状态分析器
type StateAnalyzer struct {
	memMgr *memory.MemoryManager
}

// NewStateAnalyzer 创建状态分析器
func NewStateAnalyzer(memMgr *memory.MemoryManager) *StateAnalyzer {
	return &StateAnalyzer{
		memMgr: memMgr,
	}
}

// Analyze 分析当前状态
func (sa *StateAnalyzer) Analyze() (*SystemState, error) {
	memoryState, err := sa.analyzeMemoryState()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze memory state: %w", err)
	}

	taskState, err := sa.analyzeTaskState()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze task state: %w", err)
	}

	evolutionState, err := sa.analyzeEvolutionState()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze evolution state: %w", err)
	}

	return &SystemState{
		MemoryState:    *memoryState,
		TaskState:       *taskState,
		EvolutionState:  *evolutionState,
		Timestamp:       time.Now().Format(time.RFC3339),
	}, nil
}

// analyzeMemoryState 分析记忆状态
func (sa *StateAnalyzer) analyzeMemoryState() (*MemoryState, error) {
	// 统计 archive 文件
	archiveDir := memory.ArchiveDir
	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		return nil, err
	}

	fileCount := 0
	var totalSize int64

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		if filepath.Base(entry.Name()) == "summary" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileCount++
		totalSize += info.Size()
	}

	// 检查是否需要压缩
	needsSummarize, reason := sa.memMgr.CheckSummarizeTrigger()

	// 获取 hot.md 最后更新时间
	hotFile := memory.HotFile
	hotLastUpdated := ""
	if info, err := os.Stat(hotFile); err == nil {
		hotLastUpdated = info.ModTime().Format(time.RFC3339)
	}

	// 获取索引条目数
	index := sa.memMgr.GetIndex()
	indexEntryCount := 0
	if index != nil {
		indexEntryCount = len(index.Entries)
	}

	return &MemoryState{
		ArchiveFileCount: fileCount,
		ArchiveTotalSize: totalSize,
		HotLastUpdated:   hotLastUpdated,
		IndexEntryCount:  indexEntryCount,
		IndexComplete:    indexEntryCount > 0,
		NeedsSummarize:   needsSummarize,
		SummarizeReason:  reason,
	}, nil
}

// analyzeTaskState 分析任务状态
func (sa *StateAnalyzer) analyzeTaskState() (*TaskState, error) {
	// 读取任务历史（从 evolution_log.json）
	logFile := EvolutionLogFilePath
	recentTasks := []TaskHistory{}
	successCount := 0
	totalCount := 0
	lastTaskTime := ""

	if data, err := os.ReadFile(logFile); err == nil {
		var log EvolutionLog
		if err := json.Unmarshal(data, &log); err == nil {
			// 取最近 10 个任务
			start := len(log.Entries) - 10
			if start < 0 {
				start = 0
			}
			for i := start; i < len(log.Entries); i++ {
				entry := log.Entries[i]
				recentTasks = append(recentTasks, TaskHistory{
					TaskID:    entry.TaskID,
					Type:      entry.Action,
					Status:    entry.Status,
					Completed: entry.CompletedAt,
					Result:    entry.Result,
				})
				totalCount++
				if entry.Status == "completed" {
					successCount++
				}
				if entry.CompletedAt != "" {
					lastTaskTime = entry.CompletedAt
				}
			}
		}
	}

	successRate := 0.0
	if totalCount > 0 {
		successRate = float64(successCount) / float64(totalCount)
	}

	return &TaskState{
		RecentTasks:  recentTasks,
		SuccessRate:  successRate,
		PendingTasks: 0, // TODO: 从任务队列获取
		LastTaskTime: lastTaskTime,
	}, nil
}

// analyzeEvolutionState 分析演进状态
func (sa *StateAnalyzer) analyzeEvolutionState() (*EvolutionState, error) {
	capabilities := []Capability{}
	learningProgress := make(map[string]float64)
	lastEvolution := ""
	improvementAreas := []string{}

	// 读取 capabilities.json
	capFile := CapabilitiesFilePath
	if data, err := os.ReadFile(capFile); err == nil {
		var caps Capabilities
		if err := json.Unmarshal(data, &caps); err == nil {
			capabilities = caps.Capabilities
			learningProgress = caps.LearningProgress
		}
	}

	// 读取 evolution_log.json 获取最后演进时间
	logFile := EvolutionLogFilePath
	if data, err := os.ReadFile(logFile); err == nil {
		var log EvolutionLog
		if err := json.Unmarshal(data, &log); err == nil {
			if len(log.Entries) > 0 {
				lastEvolution = log.Entries[len(log.Entries)-1].Timestamp
			}
		}
	}

	return &EvolutionState{
		Capabilities:     capabilities,
		LearningProgress: learningProgress,
		LastEvolution:   lastEvolution,
		ImprovementAreas: improvementAreas,
	}, nil
}

// Capabilities 能力记录
type Capabilities struct {
	Capabilities    []Capability         `json:"capabilities"`
	LearningProgress map[string]float64 `json:"learning_progress"`
}

// EvolutionLog 演进日志
type EvolutionLog struct {
	Entries []EvolutionLogEntry `json:"entries"`
}

// EvolutionLogEntry 演进日志条目
type EvolutionLogEntry struct {
	Timestamp   string   `json:"timestamp"`
	TaskID      string   `json:"task_id"`
	Action      string   `json:"action"`
	Decision    string   `json:"decision"`
	Status      string   `json:"status"` // "pending" | "running" | "completed" | "failed"
	Result      string   `json:"result"`
	Learning    string   `json:"learning"`
	CompletedAt string   `json:"completed_at"`
	NextSteps   []string `json:"next_steps"`
}
