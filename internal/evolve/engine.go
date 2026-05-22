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
		if err := e.runCycle(ctx, ws, false, false); err != nil {
			log.Printf("Autonomous evolution [%s]: %v", ws.ID, err)
		}
	}
}

func (e *Engine) runCycle(ctx context.Context, ws *brain.Workspace, sessionCompress, crystallize bool) error {
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

	if crystallize {
		if snap.ShortTermBytes < crystallizeMinShortBytes {
			log.Printf("Autonomous evolution [%s]: crystallize skipped (short-term too small)", ws.ID)
			return nil
		}
		if excerpt, err := readFileCap(brain.ShortTermCurrentPath(), maxShortExcerptBytes); err == nil {
			appendCrystallizeTriggers(snap, excerpt)
		}
		snap.Triggers = append(snap.Triggers, "high_token_session")
		if !shouldInvokeCrystallize(snap) {
			log.Printf("Autonomous evolution [%s]: crystallize skipped (no triggers)", ws.ID)
			return nil
		}
		log.Printf("Autonomous evolution [%s]: crystallize_skill (%s)", ws.ID, strings.Join(snap.Triggers, ","))
	} else if sessionCompress {
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

	prompt := buildDecisionPrompt(snap, sessionCompress, crystallize)
	sys := evolutionSystemPrompt()
	if sessionCompress {
		sys = evolutionSessionCompressPrompt()
	} else if crystallize {
		sys = evolutionCrystallizePrompt()
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
	if crystallize {
		dec.Updates = filterUpdatesCrystallize(dec.Updates)
	} else {
		dec.Updates = filterUpdates(dec.Updates)
	}

	var touched []string
	action := strings.ToLower(strings.TrimSpace(dec.Action))
	if action != "idle" && len(dec.Updates) > 0 {
		touched, err = ApplyUpdates(dec.Updates)
		if err != nil {
			return fmt.Errorf("apply: %w", err)
		}
	}
	if crystallize && (action == "crystallize_skill" || len(touched) > 0) {
		ingestCrystallizedSkills(ws, touched)
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
	entry := LogEntry{
		WorkspaceID: ws.ID,
		ModeID:      ws.ActiveMode,
		Action:      dec.Action,
		Reason:      dec.Reason,
		Learning:    learning,
		DocTouched:  touched,
	}
	if shouldFinalizeShortTerm(dec, touched, snap, sessionCompress) {
		if arch, err := brain.FinalizeShortTermAfterConsolidate(brain.DefaultKeepRecentAfterConsolidate); err != nil {
			log.Printf("Autonomous evolution [%s]: short-term finalize: %v", ws.ID, err)
		} else if arch != "" {
			entry.DocTouched = append(entry.DocTouched, arch)
			log.Printf("Autonomous evolution [%s]: short-term archived to %s", ws.ID, arch)
			if fresh, err := Observe(); err == nil {
				snap = fresh
			}
		}
	}
	if err := brain.SyncMemoryIndexAfterEvolution(entry.DocTouched, learning, archRel(entry.DocTouched)); err != nil {
		log.Printf("Autonomous evolution [%s]: memory index: %v", ws.ID, err)
	}
	if err := AppendLog(entry); err != nil {
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
- path 相对 workspace 根，例如：modes/_default/persona.md、persona.local.md、memory/long/note.md
- consolidate：把 short_term excerpt 中的**新事实**写入 modes/<mode>/persona.md（append），细节写 memory/long/*.md；**不要** patch memory/short/current.md（服务端会归档并更新 memory/index.json）
- 禁止整篇重写 constraints；persona 只 append 不重复已有段落

默认 idle。`
}

func evolutionSessionCompressPrompt() string {
	return evolutionSystemPrompt() + `

本轮为「对话轮次阈值」触发的强制压缩：action 应为 consolidate；将 short-term 中的新事实写入 modes/<mode>/persona.md，细节摘要可 append 到 memory/long/*.md；不要 idle；不要 patch short/current.md。`
}

func evolutionCrystallizePrompt() string {
	return `你是 Cata 自主演进模块（固化 skill）。

将 short-term 中**已验证**的探索流程固化为脑子内可执行 skill，供后续 run_skill 复用。

输出单个 JSON：action, reason, learning, updates[]
- action 应为 crystallize_skill（无合适固化则 idle）
- path 相对 workspace 根，仅允许：
  - skills/<skill-id>/SKILL.md（流程：何时用 run_skill、不适用时仍用 browser）
  - skills/<skill-id>/manifest.yaml（runner: python, entry: script.py）
  - skills/<skill-id>/script.py（从 excerpt 成功命令提炼，标准库优先）
- skill-id 用小写英文与连字符，如 zhangtingban-lianban
- **禁止** patch modes/*/capabilities.yaml（服务端会自动 append skills 列表）
- **禁止** 写入 mcp: [] 或删除 browser；未覆盖站点仍依赖 browser 基础能力
- SKILL 中写明：适用场景（如东财 A 站）、输出路径（相对产出区 cwd）、禁止 browser_snapshot 整页抓取`
}

func buildDecisionPrompt(snap *Snapshot, sessionCompress, crystallize bool) string {
	var b strings.Builder
	b.WriteString("triggers: ")
	b.WriteString(strings.Join(snap.Triggers, ", "))
	if sessionCompress {
		b.WriteString(" (session-driven compress)")
	}
	if crystallize {
		b.WriteString(" (crystallize_skill)")
	}
	if len(snap.SkillIDs) > 0 {
		b.WriteString("\nexisting_skills: ")
		b.WriteString(strings.Join(snap.SkillIDs, ", "))
	}
	b.WriteString("\nstate: ")
	compact, _ := json.Marshal(snap)
	b.Write(compact)

	if snap.ShortTermBytes >= shortTermActivityBytes {
		includeExcerpt := true
		if snap.LastEvolutionAt != "" && snap.ShortTermModTime != "" &&
			snap.ShortTermModTime <= snap.LastEvolutionAt && snap.ShortTermBytes < shortTermTriggerBytes {
			includeExcerpt = false
		}
		if includeExcerpt {
			if excerpt, err := readFileCap(brain.ShortTermCurrentPath(), maxShortExcerptBytes); err == nil && excerpt != "" {
				b.WriteString("\n\nshort_term excerpt:\n")
				b.WriteString(excerpt)
			}
		} else {
			b.WriteString("\n\n(short_term unchanged since last evolution; excerpt omitted)\n")
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
	return NewEngine(interval).runCycle(ctx, ws, false, false)
}
