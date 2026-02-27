package memory

// MemoryPiece 表示一个记忆片段
type MemoryPiece struct {
	Content  string `json:"content"`  // MD 原文或片段
	Category string `json:"category"`  // 类别：preference、fact、logic
	Source   string `json:"source"`    // 来源文件路径，如 /brain/hot.md、/brain/archive/2025-02-13.md
	Priority int    `json:"priority"`  // 优先级 0–10，越高越优先参与 Recall
}

// IndexEntry 表示 memory_index.json 中的一个条目
type IndexEntry struct {
	Source   string   `json:"source"`   // 文件路径
	Keywords []string `json:"keywords"` // 关键词列表
	Summary  string   `json:"summary"`  // 摘要
	Category string   `json:"category"` // 类别：preference、fact、logic
	Priority int      `json:"priority"` // 优先级 0–10
}

// MemoryIndex 表示整个索引结构
type MemoryIndex struct {
	Version   int          `json:"version"`   // 版本号
	UpdatedAt string       `json:"updated_at"` // ISO 8601 时间戳
	Entries   []IndexEntry `json:"entries"`   // 索引条目列表
}
