package memory

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// HotSection 热记忆区块映射
var HotSectionMap = map[string]string{
	"我是谁":              "## 我是谁",
	"身份":               "## 我是谁",
	"自我认知":             "## 我是谁",
	"当前目标":             "## 当前目标",
	"目标":               "## 当前目标",
	"雷打不动的偏好":         "## 雷打不动的偏好",
	"偏好":               "## 雷打不动的偏好",
	"开发":               "## 开发 · 技术栈与习惯",
	"技术栈":              "## 开发 · 技术栈与习惯",
	"开发习惯":             "## 开发 · 技术栈与习惯",
	"dev":               "## 开发 · 技术栈与习惯",
	"学习":               "## 学习 · 当前方向与节奏",
	"学习方向":             "## 学习 · 当前方向与节奏",
	"学习节奏":             "## 学习 · 当前方向与节奏",
	"learning":           "## 学习 · 当前方向与节奏",
	"生活":               "## 生活 · 作息与健康偏好",
	"作息":               "## 生活 · 作息与健康偏好",
	"健康":               "## 生活 · 作息与健康偏好",
	"life":              "## 生活 · 作息与健康偏好",
}

// determineHotSection 根据 topic 确定应该写入的 hot 区块
func determineHotSection(topic string) string {
	topicLower := strings.ToLower(topic)
	
	// 精确匹配
	for key, section := range HotSectionMap {
		if strings.Contains(topicLower, strings.ToLower(key)) {
			return section
		}
	}

	// 默认写入到"雷打不动的偏好"
	return "## 雷打不动的偏好"
}

// appendToHotFileWithSection 追加内容到 hot.md 的指定区块
func appendToHotFileWithSection(topic, content, section string) error {
	data, err := os.ReadFile(HotFile)
	if err != nil {
		return err
	}

	fileContent := string(data)
	
	// 查找区块位置
	sectionPattern := regexp.MustCompile(fmt.Sprintf(`(%s)\s*\n([^\n]*\n)*`, regexp.QuoteMeta(section)))
	matches := sectionPattern.FindStringSubmatch(fileContent)

	if len(matches) > 0 {
		// 找到区块，在区块末尾追加
		sectionEnd := matches[0]
		insertPos := strings.Index(fileContent, sectionEnd) + len(sectionEnd)
		
		// 检查是否已有内容，如果有则添加分隔
		beforeInsert := fileContent[:insertPos]
		afterInsert := fileContent[insertPos:]
		
		// 如果区块后没有内容或只有空行，直接追加
		trimmedAfter := strings.TrimSpace(afterInsert)
		if trimmedAfter == "" || strings.HasPrefix(trimmedAfter, "##") {
			// 在区块内追加
			newContent := fmt.Sprintf("%s\n### %s\n\n%s\n\n", beforeInsert, topic, content)
			return os.WriteFile(HotFile, []byte(newContent+afterInsert), 0644)
		} else {
			// 在区块末尾追加
			newContent := fmt.Sprintf("%s\n### %s\n\n%s\n\n", beforeInsert, topic, content)
			return os.WriteFile(HotFile, []byte(newContent+afterInsert), 0644)
		}
	}

	// 如果找不到区块，在文件末尾追加新区块
	newSection := fmt.Sprintf("\n%s\n\n### %s\n\n%s\n", section, topic, content)
	return os.WriteFile(HotFile, append(data, []byte(newSection)...), 0644)
}
