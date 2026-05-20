package evolve

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"mybot/internal/brain"
	"mybot/internal/config"
	"mybot/internal/llm"
)

// Engine 后台自主演进。
type Engine struct {
	interval time.Duration

	mu              sync.Mutex
	lastFingerprint map[string]string
	cooldownUntil   map[string]time.Time
}

// NewEngine 创建演进引擎。
func NewEngine(interval time.Duration) *Engine {
	if interval <= 0 {
		interval = DefaultCycleSeconds * time.Second
	}
	return &Engine{
		interval:        interval,
		lastFingerprint: make(map[string]string),
		cooldownUntil:   make(map[string]time.Time),
	}
}

// Start 周期执行；对每个已注册 workspace 分别门控与演进。
func (e *Engine) Start(ctx context.Context) {
	if config.Config == nil || !config.Config.LLM.Enabled {
		log.Println("Autonomous evolution: skipped (LLM not enabled)")
		return
	}
	if config.Config != nil && !config.Config.Evolution.Enabled {
		log.Println("Autonomous evolution: disabled in config")
		return
	}

	log.Printf("Autonomous evolution: started (interval %s, per-workspace)", e.interval)

	go func() {
		ticker := time.NewTicker(e.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Autonomous evolution: stopped")
				return
			case <-ticker.C:
				e.runAll(ctx)
			}
		}
	}()
}

func (e *Engine) runAll(ctx context.Context) {
	_ = brain.EnsureCataLayout()
	list, err := brain.ListWorkspaces()
	if err != nil {
		log.Printf("Autonomous evolution: list workspaces: %v", err)
		return
	}
	if len(list) == 0 {
		return
	}
	for _, ws := range list {
		brain.SetActive(ws)
		if err := e.runCycle(ctx, ws, false); err != nil {
			log.Printf("Autonomous evolution [%s]: %v", ws.ID, err)
		}
	}
}

func (e *Engine) runCycle(ctx context.Context, ws *brain.Workspace, sessionCompress bool) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	snap, err := Observe()
	if err != nil {
		return fmt.Errorf("observe: %w", err)
	}

	e.mu.Lock()
	cooldown := e.cooldownUntil[ws.ID]
	lastFP := e.lastFingerprint[ws.ID]
	e.mu.Unlock()

	if sessionCompress {
		if snap.ShortTermBytes < sessionCompressMinShortBytes {
			log.Printf("Autonomous evolution [%s]: session compress skipped (short-term too small)", ws.ID)
			return nil
		}
		snap.Triggers = append([]string{"session_turn_threshold"}, snap.Triggers...)
		log.Printf("Autonomous evolution [%s]: session compress (turn threshold)", ws.ID)
	} else {
		ok, reason := shouldInvokeLLM(snap, cooldown, lastFP)
		if !ok {
			log.Printf("Autonomous evolution [%s]: skip LLM (%s)", ws.ID, reason)
			return nil
		}
	}

	client, err := llm.NewClientForRole(llm.RoleEvolution)
	if err != nil {
		return fmt.Errorf("LLM: %w", err)
	}

	prompt := buildDecisionPrompt(snap, sessionCompress)
	sys := evolutionSystemPrompt()
	if sessionCompress {
		sys = evolutionSessionCompressPrompt()
	}
	messages := []llm.Message{
		{Role: "system", Content: sys},
		{Role: "user", Content: prompt},
	}

	reply, err := client.ChatEvolution(messages)
	if err != nil {
		return fmt.Errorf("decide: %w", err)
	}

	dec, err := parseDecision(reply)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	dec.Updates = filterUpdates(dec.Updates)

	var touched []string
	action := strings.ToLower(strings.TrimSpace(dec.Action))
	if action != "idle" && len(dec.Updates) > 0 {
		touched, err = ApplyUpdates(dec.Updates)
		if err != nil {
			return fmt.Errorf("apply: %w", err)
		}
	}

	if !isMeaningfulDecision(dec, touched) {
		log.Printf("Autonomous evolution [%s]: no-op (action=%s)", ws.ID, dec.Action)
		e.mu.Lock()
		e.lastFingerprint[ws.ID] = snap.Fingerprint()
		e.mu.Unlock()
		return nil
	}

	learning := strings.TrimSpace(dec.Learning)
	if learning == "" {
		learning = dec.Reason
	}
	if err := AppendLog(LogEntry{
		WorkspaceID: ws.ID,
		ModeID:      ws.ActiveMode,
		Action:      dec.Action,
		Reason:      dec.Reason,
		Learning:    learning,
		DocTouched:  touched,
	}); err != nil {
		return err
	}

	e.mu.Lock()
	e.lastFingerprint[ws.ID] = snap.Fingerprint()
	if len(touched) > 0 {
		e.cooldownUntil[ws.ID] = time.Now().Add(e.interval)
	}
	e.mu.Unlock()

	log.Printf("Autonomous evolution [%s]: action=%s files=%v", ws.ID, dec.Action, touched)
	return nil
}

func evolutionSystemPrompt() string {
	return `你是 Cata 自主演进模块。只修改 ~/.cata 脑子目录内 Markdown；产出区（用户 cwd）不在此写入。

persona 由你从脑子内 short-term 提炼；终端对话不直接改 persona。

只在 triggers 成立时修改当前脑子分区内文件；否则 action=idle 且 updates=[]。

输出：单个 JSON 对象。
字段：action, reason, learning, updates[]
- path 相对 workspace 根，例如：modes/_default/persona.md、persona.local.md、memory/short/current.md、memory/long/note.md
- consolidate：short → modes/<mode>/persona.md，必要时写 memory/long/
- 禁止整篇重写 constraints；persona 用 append

默认 idle。`
}

func evolutionSessionCompressPrompt() string {
	return evolutionSystemPrompt() + `

本轮为「对话轮次阈值」触发的强制压缩：action 应为 consolidate；将 short-term 稳定内容写入 modes/<mode>/persona.md，重复内容可写入 memory/long/ 或从 short-term 删减；不要 idle。`
}

func buildDecisionPrompt(snap *Snapshot, sessionCompress bool) string {
	var b strings.Builder
	b.WriteString("triggers: ")
	b.WriteString(strings.Join(snap.Triggers, ", "))
	if sessionCompress {
		b.WriteString(" (session-driven compress)")
	}
	b.WriteString("\nstate: ")
	compact, _ := json.Marshal(snap)
	b.Write(compact)

	if snap.ShortTermBytes >= shortTermActivityBytes {
		if excerpt, err := readFileCap(brain.ShortTermCurrentPath(), maxShortExcerptBytes); err == nil && excerpt != "" {
			b.WriteString("\n\nshort_term excerpt:\n")
			b.WriteString(excerpt)
		}
		if hot, err := readFileCap(brain.HotPath(), 1200); err == nil && hot != "" {
			b.WriteString("\n\ncurrent mode persona (merge here, append only):\n")
			b.WriteString(hot)
		}
	}
	if snap.RecentLogSummary != "" {
		b.WriteString("\n\nrecent evolution: ")
		b.WriteString(snap.RecentLogSummary)
	}
	b.WriteString("\n\nRespond with decision JSON only.")
	return b.String()
}

func readFileCap(path string, max int) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	s := string(data)
	if len(s) > max {
		return s[:max] + "\n…(truncated)", nil
	}
	return s, nil
}

// RunCycle 对当前活跃 workspace 执行一轮（测试用）。
func RunCycle(ctx context.Context) error {
	ws, err := brain.MustActive()
	if err != nil {
		return err
	}
	interval := DefaultCycleSeconds * time.Second
	if config.Config != nil && config.Config.Evolution.CycleInterval > 0 {
		interval = time.Duration(config.Config.Evolution.CycleInterval) * time.Second
	}
	return NewEngine(interval).runCycle(ctx, ws, false)
}
