package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// BuildIndexFromMarkdown 扫描 brain/**/*.md 并构建索引
func BuildIndexFromMarkdown() (*MemoryIndex, error) {
	index := &MemoryIndex{
		Version:   IndexVersion,
		UpdatedAt: time.Now().Format(time.RFC3339),
		Entries:   []IndexEntry{},
	}

	// 扫描 brain 目录下的所有 .md 文件
	err := filepath.Walk(BrainDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和非 .md 文件
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// 跳过索引文件本身（如果存在）
		if strings.Contains(path, "memory_index.json") {
			return nil
		}

		// 读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// 提取关键词和摘要
		entry := extractEntry(path, string(content))
		index.Entries = append(index.Entries, entry)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan markdown files: %w", err)
	}

	return index, nil
}

// LoadOrBuildIndex 加载索引，如果不存在或无效则重建
func LoadOrBuildIndex() (*MemoryIndex, error) {
	// 尝试加载现有索引
	index, err := LoadIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to load index: %w", err)
	}

	// 如果索引为空或版本不匹配，重建索引
	if len(index.Entries) == 0 || index.Version != IndexVersion {
		index, err = BuildIndexFromMarkdown()
		if err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}

		// 保存新构建的索引
		if err := SaveIndex(index); err != nil {
			return nil, fmt.Errorf("failed to save index: %w", err)
		}
	}

	return index, nil
}

// extractEntry 从 Markdown 内容中提取索引条目
func extractEntry(source, content string) IndexEntry {
	entry := IndexEntry{
		Source:   source,
		Keywords: []string{},
		Summary:  "",
		Category: determineCategory(source, content),
		Priority: determinePriority(source),
	}

	// 提取标题作为关键词
	titleRegex := regexp.MustCompile(`^#+\s+(.+)$`)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if matches := titleRegex.FindStringSubmatch(line); len(matches) > 1 {
			title := strings.TrimSpace(matches[1])
			// 移除 Markdown 格式标记
			title = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`).ReplaceAllString(title, "$1")
			if title != "" {
				entry.Keywords = append(entry.Keywords, title)
			}
		}
	}

	// 提取摘要：取前 200 个字符（去除标题和空行）
	summary := extractSummaryFromContent(content)
	if len(summary) > 200 {
		summary = summary[:200] + "..."
	}
	entry.Summary = summary

	// 从内容中提取更多关键词（简单实现：提取常见的中文词汇和英文单词）
	keywords := extractKeywordsFromContent(content)
	entry.Keywords = append(entry.Keywords, keywords...)

	return entry
}

// extractSummaryFromContent 提取摘要（去除标题和空行后的前几行）
func extractSummaryFromContent(content string) string {
	lines := strings.Split(content, "\n")
	var summaryLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和标题行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// 跳过引用块
		if strings.HasPrefix(line, ">") {
			continue
		}
		summaryLines = append(summaryLines, line)
		if len(summaryLines) >= 3 {
			break
		}
	}

	return strings.Join(summaryLines, " ")
}

// extractKeywordsFromContent 从内容中提取关键词
func extractKeywordsFromContent(content string) []string {
	keywords := []string{}
	
	// 提取中文词汇（简单实现：2-4字词组）
	chineseWordRegex := regexp.MustCompile(`[\x{4e00}-\x{9fa5}]{2,4}`)
	chineseWords := chineseWordRegex.FindAllString(content, -1)
	
	// 去重并限制数量
	seen := make(map[string]bool)
	for _, word := range chineseWords {
		if !seen[word] && len(keywords) < 10 {
			keywords = append(keywords, word)
			seen[word] = true
		}
	}

	// 提取英文单词（3个字符以上）
	englishWordRegex := regexp.MustCompile(`\b[a-zA-Z]{3,}\b`)
	englishWords := englishWordRegex.FindAllString(content, -1)
	
	for _, word := range englishWords {
		word = strings.ToLower(word)
		// 跳过常见停用词
		stopWords := map[string]bool{
			"the": true, "and": true, "for": true, "are": true,
			"but": true, "not": true, "you": true, "all": true,
		}
		if stopWords[word] {
			continue
		}
		if !seen[word] && len(keywords) < 15 {
			keywords = append(keywords, word)
			seen[word] = true
		}
	}

	return keywords
}

// determineCategory 根据文件路径和内容确定类别
func determineCategory(source, content string) string {
	// hot.md 通常是 preference
	if strings.Contains(source, "hot.md") {
		return "preference"
	}

	// archive 文件根据内容判断
	if strings.Contains(content, "目标") || strings.Contains(content, "偏好") {
		return "preference"
	}
	if strings.Contains(content, "项目") || strings.Contains(content, "设计") || strings.Contains(content, "架构") {
		return "logic"
	}

	// 默认为 fact
	return "fact"
}

// determinePriority 根据文件路径确定优先级
func determinePriority(source string) int {
	// hot.md 优先级最高
	if strings.Contains(source, "hot.md") {
		return 10
	}

	// archive 文件按新鲜度（这里简化处理，实际可以根据文件日期）
	if strings.Contains(source, "archive") {
		return 5
	}

	return 3
}
