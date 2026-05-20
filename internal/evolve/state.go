package evolve

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"mybot/internal/brain"
)

// Snapshot 自主演进 Observe 阶段的只读状态（仅元数据，不把整库塞进 LLM）。
type Snapshot struct {
	ObservedAt            string   `json:"observed_at"`
	HotModTime            string   `json:"hot_mod_time,omitempty"`
	ShortTermModTime      string   `json:"short_term_mod_time,omitempty"`
	ShortTermBytes        int64    `json:"short_term_bytes"`
	LongTermFileCount     int      `json:"long_term_file_count"`
	ArchiveFileCount      int      `json:"archive_file_count"`
	LastEvolutionAt       string   `json:"last_evolution_at,omitempty"`
	LastEvolutionAction   string   `json:"last_evolution_action,omitempty"`
	RecentLogSummary      string   `json:"recent_log_summary,omitempty"`
	Triggers              []string `json:"triggers,omitempty"`
}

// Fingerprint 仅跟踪演进「输入」信号（不含 hot：hot 由演进写出，不应作为触发依据）。
func (s *Snapshot) Fingerprint() string {
	return fmt.Sprintf("st:%s|sb:%d|lt:%d|ar:%d",
		s.ShortTermModTime, s.ShortTermBytes, s.LongTermFileCount, s.ArchiveFileCount)
}

// Observe 读取 ~/.cata/brain 元数据（不读 workflow/core 全文）。
func Observe() (*Snapshot, error) {
	s := &Snapshot{ObservedAt: time.Now().UTC().Format(time.RFC3339)}

	if info, err := os.Stat(brain.HotPath()); err == nil {
		s.HotModTime = info.ModTime().UTC().Format(time.RFC3339)
	}

	shortPath := brain.ShortTermCurrentPath()
	if info, err := os.Stat(shortPath); err == nil {
		s.ShortTermModTime = info.ModTime().UTC().Format(time.RFC3339)
	}
	if data, err := os.ReadFile(shortPath); err == nil {
		s.ShortTermBytes = int64(len(data))
	}

	longDir := brain.LongTermDir()
	if entries, err := os.ReadDir(longDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				s.LongTermFileCount++
			}
		}
	}

	archiveDir := brain.ArchiveDir()
	if entries, err := os.ReadDir(archiveDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				s.ArchiveFileCount++
			}
		}
	}

	loadLastEvolutionMeta(s)
	s.RecentLogSummary = summarizeRecentLog(2, 80)
	computeTriggers(s)
	return s, nil
}

func loadLastEvolutionMeta(s *Snapshot) {
	data, err := os.ReadFile(brain.EvolutionLogPath())
	if err != nil {
		return
	}
	var log EvolutionLog
	if err := json.Unmarshal(data, &log); err != nil || len(log.Entries) == 0 {
		return
	}
	last := log.Entries[len(log.Entries)-1]
	s.LastEvolutionAt = last.Timestamp
	s.LastEvolutionAction = last.Action
}

func summarizeRecentLog(n int, maxLearning int) string {
	data, err := os.ReadFile(brain.EvolutionLogPath())
	if err != nil {
		return ""
	}
	var log EvolutionLog
	if err := json.Unmarshal(data, &log); err != nil || len(log.Entries) == 0 {
		return ""
	}
	start := len(log.Entries) - n
	if start < 0 {
		start = 0
	}
	var b strings.Builder
	for i := start; i < len(log.Entries); i++ {
		e := log.Entries[i]
		b.WriteString(e.Action)
		if e.Learning != "" {
			b.WriteString(": ")
			b.WriteString(truncate(e.Learning, maxLearning))
		}
		if i < len(log.Entries)-1 {
			b.WriteString(" | ")
		}
	}
	return b.String()
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
