package memory

import (
	"fmt"
	"strings"
)

// FormatMemoryPiecesForContext 将 MemoryPiece 列表格式化为文本，用于注入 LLM Context
func FormatMemoryPiecesForContext(pieces []MemoryPiece) string {
	if len(pieces) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("# 相关记忆片段\n\n")

	for i, piece := range pieces {
		builder.WriteString(fmt.Sprintf("## 片段 %d\n\n", i+1))
		
		// 元数据
		builder.WriteString(fmt.Sprintf("**来源**: %s\n", piece.Source))
		builder.WriteString(fmt.Sprintf("**类别**: %s\n", piece.Category))
		builder.WriteString(fmt.Sprintf("**优先级**: %d/10\n\n", piece.Priority))
		
		// 内容（限制长度，避免 Context 过长）
		content := piece.Content
		maxLength := 1000 // 每个片段最多 1000 字符
		if len(content) > maxLength {
			content = content[:maxLength] + "...\n\n(内容已截断)"
		}
		builder.WriteString(content)
		builder.WriteString("\n\n---\n\n")
	}

	return builder.String()
}

// FormatMemoryPiecesForSummary 将 MemoryPiece 列表格式化为摘要格式（更简洁）
func FormatMemoryPiecesForSummary(pieces []MemoryPiece) string {
	if len(pieces) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("相关记忆摘要：\n\n")

	for i, piece := range pieces {
		// 只包含来源和简短摘要
		content := piece.Content
		maxLength := 200
		if len(content) > maxLength {
			// 尝试找到句号截断
			if idx := strings.LastIndex(content[:maxLength], "。"); idx > 0 {
				content = content[:idx+3] + "..."
			} else {
				content = content[:maxLength] + "..."
			}
		}

		builder.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, piece.Source, content))
	}

	return builder.String()
}
