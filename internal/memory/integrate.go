package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"mybot/internal/config"
)

// BuildIntegratedSystemPrompt 将整个 brain 与分散记忆整合为一份 system prompt：
// core、workflow、hot、archive、long-term、short-term、graph_memory 合并写入 brain/context/integrated_system_prompt.md
func BuildIntegratedSystemPrompt() (outPath string, size int, err error) {
	brainDir := config.GetBrainDir()
	ctxDir := config.GetBrainPath("context")
	if err := os.MkdirAll(ctxDir, 0755); err != nil {
		return "", 0, fmt.Errorf("create context dir: %w", err)
	}
	outPath = config.GetBrainPath("context/integrated_system_prompt.md")

	var b strings.Builder

	// 1. 核心思维与流程（作为 system prompt 骨架）
	corePath := config.GetBrainPath("core.md")
	if data, e := os.ReadFile(corePath); e == nil {
		b.WriteString("---\n# Core（核心思维）\n---\n\n")
		b.Write(data)
		b.WriteString("\n\n")
	}
	workflowPath := config.GetBrainPath("workflow.md")
	if data, e := os.ReadFile(workflowPath); e == nil {
		b.WriteString("---\n# Workflow（自主演进流程）\n---\n\n")
		b.Write(data)
		b.WriteString("\n\n")
	}

	// 2. 热记忆
	if data, e := os.ReadFile(HotFile); e == nil {
		b.WriteString("---\n# Hot（热记忆）\n---\n\n")
		b.Write(data)
		b.WriteString("\n\n")
	}

	// 3. 长期记忆 long-term
	longDir := filepath.Join(brainDir, "memory", "long-term")
	if entries, e := os.ReadDir(longDir); e == nil {
		var names []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				names = append(names, e.Name())
			}
		}
		sort.Strings(names)
		for _, n := range names {
			fpath := filepath.Join(longDir, n)
			if data, e := os.ReadFile(fpath); e == nil {
				b.WriteString("---\n# Long-term / " + n + "\n---\n\n")
				b.Write(data)
				b.WriteString("\n\n")
			}
		}
	}

	// 4. 短期记忆 short-term
	shortPath := config.GetBrainPath("memory/short-term/current_session.md")
	if data, e := os.ReadFile(shortPath); e == nil {
		b.WriteString("---\n# Short-term（当前会话）\n---\n\n")
		b.Write(data)
		b.WriteString("\n\n")
	}

	// 5. 图记忆 / 知识图谱
	graphPath := config.GetBrainPath("knowledge_graph_memory.md")
	if data, e := os.ReadFile(graphPath); e == nil {
		b.WriteString("---\n# Graph Memory（知识图谱记忆）\n---\n\n")
		b.Write(data)
		b.WriteString("\n\n")
	}

	// 6. 档案 archive（按文件名排序，跳过 backup）
	if entries, e := os.ReadDir(ArchiveDir); e == nil {
		var names []string
		for _, e := range entries {
			if e.IsDir() {
				if e.Name() == "backup" {
					continue
				}
				continue
			}
			if strings.HasSuffix(e.Name(), ".md") {
				names = append(names, e.Name())
			}
		}
		sort.Strings(names)
		b.WriteString("---\n# Archive（档案）\n---\n\n")
		for _, n := range names {
			fpath := filepath.Join(ArchiveDir, n)
			if data, e := os.ReadFile(fpath); e == nil {
				b.WriteString("## " + n + "\n\n")
				b.Write(data)
				b.WriteString("\n\n")
			}
		}
	}

	content := []byte(strings.TrimSpace(b.String()))
	if err := os.WriteFile(outPath, content, 0644); err != nil {
		return "", 0, fmt.Errorf("write integrated prompt: %w", err)
	}
	return outPath, len(content), nil
}
