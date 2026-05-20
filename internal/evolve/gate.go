package evolve

import (
	"fmt"
	"strings"
	"time"
)

// shouldInvokeLLM 在调用模型前做确定性门控，避免空转与重复扫描。
// hot.md 由本模块写入，不作为「hot 已改」触发条件。
func shouldInvokeLLM(snap *Snapshot, cooldownUntil time.Time, lastFingerprint string) (bool, string) {
	if time.Now().Before(cooldownUntil) {
		return false, "cooldown after last patch"
	}
	fp := snap.Fingerprint()
	if fp != "" && fp == lastFingerprint {
		return false, "inputs unchanged since last cycle (short/archive/long)"
	}
	if len(snap.Triggers) == 0 {
		return false, "no triggers (need short-term activity or large archive)"
	}
	return true, strings.Join(snap.Triggers, ", ")
}

// computeTriggers 根据「待提炼进 hot 的输入」判断，不观测 hot 是否被改过。
func computeTriggers(s *Snapshot) {
	s.Triggers = nil

	if s.ShortTermBytes >= shortTermTriggerBytes {
		s.Triggers = append(s.Triggers, fmt.Sprintf("short_term>=%dB", shortTermTriggerBytes))
		return // 大短期记忆优先 consolidate，一条 trigger 即可
	}

	if s.ShortTermBytes >= shortTermActivityBytes {
		if s.LastEvolutionAt == "" {
			s.Triggers = append(s.Triggers, fmt.Sprintf("short_term>=%dB_first", shortTermActivityBytes))
		} else if s.ShortTermModTime != "" && s.ShortTermModTime > s.LastEvolutionAt {
			s.Triggers = append(s.Triggers, "short_term_updated_since_evolution")
		}
	}

	if s.ArchiveFileCount >= archiveSummarizeMinFiles {
		s.Triggers = append(s.Triggers, fmt.Sprintf("archive>=%d", archiveSummarizeMinFiles))
	}
}
