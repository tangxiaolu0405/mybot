package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// DefaultMaxArchiveFiles 默认最大 archive 文件数
	DefaultMaxArchiveFiles = 30
	// DefaultMaxArchiveSize 默认最大 archive 总大小（字节）
	DefaultMaxArchiveSize = 10 * 1024 * 1024 // 10MB
)

// SummarizeTrigger 摘要触发条件
type SummarizeTrigger struct {
	MaxArchiveFiles int   // 最大 archive 文件数
	MaxArchiveSize  int64 // 最大 archive 总大小（字节）
	LastCheck       time.Time
}

// NewSummarizeTrigger 创建摘要触发器
func NewSummarizeTrigger(maxFiles int, maxSize int64) *SummarizeTrigger {
	return &SummarizeTrigger{
		MaxArchiveFiles: maxFiles,
		MaxArchiveSize:  maxSize,
		LastCheck:       time.Now(),
	}
}

// ShouldSummarize 检查是否应该触发摘要
func (st *SummarizeTrigger) ShouldSummarize() (bool, string) {
	// 检查 archive 文件数和大小
	fileCount, totalSize, err := st.getArchiveStats()
	if err != nil {
		return false, fmt.Sprintf("Failed to get archive stats: %v", err)
	}

	reasons := []string{}
	if fileCount > st.MaxArchiveFiles {
		reasons = append(reasons, fmt.Sprintf("file count (%d > %d)", fileCount, st.MaxArchiveFiles))
	}
	if totalSize > st.MaxArchiveSize {
		reasons = append(reasons, fmt.Sprintf("total size (%d > %d)", totalSize, st.MaxArchiveSize))
	}

	if len(reasons) > 0 {
		return true, fmt.Sprintf("Archive limits exceeded: %v", reasons)
	}

	return false, ""
}

// getArchiveStats 获取 archive 统计信息
func (st *SummarizeTrigger) getArchiveStats() (fileCount int, totalSize int64, err error) {
	entries, err := os.ReadDir(ArchiveDir)
	if err != nil {
		return 0, 0, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 只统计 .md 文件
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		// 跳过 summary 文件
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

	return fileCount, totalSize, nil
}

// CheckAndSummarize 检查并触发摘要（如果满足条件）
func (m *MemoryManager) CheckAndSummarize() (bool, error) {
	trigger := NewSummarizeTrigger(DefaultMaxArchiveFiles, DefaultMaxArchiveSize)
	shouldSummarize, _ := trigger.ShouldSummarize()
	
	if !shouldSummarize {
		return false, nil
	}

	// 调用 SummarizeAndRotate
	if err := m.SummarizeAndRotate(); err != nil {
		return true, fmt.Errorf("failed to summarize: %w", err)
	}

	return true, nil
}
