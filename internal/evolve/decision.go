package evolve

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Decision LLM 一轮自主演进决策。
type Decision struct {
	Action   string      `json:"action"`
	Reason   string      `json:"reason"`
	Learning string      `json:"learning"`
	Updates  []DocUpdate `json:"updates"`
}

func parseDecision(raw string) (*Decision, error) {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end <= start {
		return nil, fmt.Errorf("no JSON object in LLM response")
	}
	raw = raw[start : end+1]

	var d Decision
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		return nil, fmt.Errorf("parse decision JSON: %w", err)
	}
	if d.Action == "" {
		d.Action = "idle"
	}
	return &d, nil
}
