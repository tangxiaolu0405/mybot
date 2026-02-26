package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mybot/internal/brain"
)

const (
	IndexVersion = 1
)

var (
	// BrainDir brain 目录路径（与 brain/core.md 一致，来源 internal/brain）
	BrainDir   string
	HotFile    string
	ArchiveDir string
	IndexFile  string
)

// initPaths 初始化路径，与 brain/core.md 资源路径表对齐
func initPaths() {
	BrainDir = brain.Dir()
	HotFile = brain.HotPath()
	ArchiveDir = brain.ArchiveDir()
	IndexFile = brain.MemoryIndexPath()
}

func init() {
	initPaths()
}

// InitBrainDirectory 初始化 brain 目录结构，与 brain/core.md 资源路径一致。
func InitBrainDirectory() error {
	// 创建 brain 目录
	if err := os.MkdirAll(BrainDir, 0755); err != nil {
		return fmt.Errorf("failed to create brain directory: %w", err)
	}

	// 创建 archive 目录
	if err := os.MkdirAll(ArchiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// 与 core.md 一致：短期记忆、长期记忆目录
	if err := os.MkdirAll(filepath.Join(BrainDir, "memory", "short-term"), 0755); err != nil {
		return fmt.Errorf("failed to create short-term memory dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(BrainDir, "memory", "long-term"), 0755); err != nil {
		return fmt.Errorf("failed to create long-term memory dir: %w", err)
	}

	// 初始化 hot.md（如果不存在）
	if _, err := os.Stat(HotFile); os.IsNotExist(err) {
		if err := initHotFile(); err != nil {
			return fmt.Errorf("failed to initialize hot.md: %w", err)
		}
	}

	return nil
}

// initHotFile 创建初始的 hot.md 文件，包含占位区块
func initHotFile() error {
	content := `# Cata · 热记忆

> 身份认同、当前目标、不可妥协的偏好（与 brain/core.md 可选存储一致）。

## 我是谁

（在这里记录你的身份认同、核心特质）

## 当前目标

（在这里记录当前的主要目标和方向）

## 雷打不动的偏好

（在这里记录不可妥协的偏好和原则）

---

## 开发 · 技术栈与习惯

（在这里记录开发相关的技术栈、习惯和偏好）

## 学习 · 当前方向与节奏

（在这里记录学习相关的方向和节奏）

## 生活 · 作息与健康偏好

（在这里记录生活相关的作息和健康偏好）
`

	return os.WriteFile(HotFile, []byte(content), 0644)
}

// GetArchivePath 根据日期返回 archive 文件路径
func GetArchivePath(date time.Time) string {
	return filepath.Join(ArchiveDir, date.Format("2006-01-02")+".md")
}

// EnsureArchiveFile 确保指定日期的 archive 文件存在，不存在则创建
func EnsureArchiveFile(date time.Time) (string, error) {
	path := GetArchivePath(date)
	
	// 如果文件已存在，直接返回
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// 创建新文件，包含日期标题
	content := fmt.Sprintf("# %s\n\n", date.Format("2006-01-02"))
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create archive file: %w", err)
	}

	return path, nil
}

// LoadIndex 从文件加载 memory_index.json
func LoadIndex() (*MemoryIndex, error) {
	data, err := os.ReadFile(IndexFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 索引不存在，返回空索引
			return &MemoryIndex{
				Version:   IndexVersion,
				UpdatedAt: time.Now().Format(time.RFC3339),
				Entries:   []IndexEntry{},
			}, nil
		}
		return nil, fmt.Errorf("failed to read index file: %w", err)
	}

	var index MemoryIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index file: %w", err)
	}

	return &index, nil
}

// SaveIndex 将索引保存到文件
func SaveIndex(index *MemoryIndex) error {
	index.UpdatedAt = time.Now().Format(time.RFC3339)
	index.Version = IndexVersion

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if err := os.WriteFile(IndexFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}
