package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"mybot/internal/brain"
	"mybot/internal/config"
	"mybot/internal/llm"
)

type toolRoute struct {
	serverName string
	toolName   string
}

// Manager 管理已连接的 MCP server 与工具路由。
type Manager struct {
	mu        sync.RWMutex
	clients   map[string]*stdioClient
	routes    map[string]*toolRoute
	llmTools  []llm.Tool
	timeout   time.Duration
	maxOutput int
}

var (
	global     *Manager
	initMu     sync.Mutex
	lastMCPKey string
)

// 终端 chat 仅暴露高频 browser 工具，避免 20+ 工具撑爆上下文导致网关/进程异常。
var preferredBrowserTools = []string{
	"browser_navigate", "browser_snapshot", "browser_click", "browser_type",
	"browser_fill", "browser_tabs", "browser_wait_for", "browser_take_screenshot",
	"browser_scroll", "browser_select_option", "browser_press_key", "browser_navigate_back",
}

const maxExportedMCPTools = 14

// Init 按配置与 capabilities 启动 MCP；失败的服务器仅记日志。
func Init(cfg config.MCPConfig, caps brain.Capabilities) *Manager {
	mgr := &Manager{
		clients:   make(map[string]*stdioClient),
		routes:    make(map[string]*toolRoute),
		timeout:   time.Duration(cfg.ToolTimeoutSeconds) * time.Second,
		maxOutput: cfg.MaxOutputBytes,
	}
	if !cfg.Enabled {
		global = mgr
		return mgr
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	for _, s := range cfg.Servers {
		if !s.Enabled {
			continue
		}
		if !caps.AllowsMCPServer(s.Name) {
			continue
		}
		if err := mgr.connectServer(ctx, s); err != nil {
			log.Printf("MCP server %q: %v", s.Name, err)
		}
	}
	global = mgr
	if n := len(mgr.llmTools); n > 0 {
		log.Printf("MCP: %d tool(s) from %d server(s)", n, len(mgr.clients))
	}
	return mgr
}

func connectServer(mgr *Manager, ctx context.Context, s config.MCPServerEntry) error {
	c, err := startStdioClient(ctx, s.Name, s.Command, s.Args, s.Env)
	if err != nil {
		return err
	}
	listCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	tools, err := c.listTools(listCtx)
	if err != nil {
		_ = c.Close()
		return err
	}
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.clients[s.Name] = c
	exported := 0
	byName := make(map[string]listedTool)
	for _, t := range tools {
		if strings.TrimSpace(t.Name) != "" {
			byName[t.Name] = t
		}
	}
	for _, name := range preferredBrowserTools {
		if exported >= maxExportedMCPTools {
			break
		}
		t, ok := byName[name]
		if !ok {
			continue
		}
		mgr.routes[name] = &toolRoute{serverName: s.Name, toolName: name}
		mgr.llmTools = append(mgr.llmTools, toLLMTool(t))
		exported++
	}
	return nil
}

func mcpCapsKey(caps brain.Capabilities) string {
	parts := make([]string, len(caps.MCP))
	copy(parts, caps.MCP)
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.ToLower(strings.TrimSpace(parts[i]))
	}
	return strings.Join(parts, ",")
}

// EnsureInit 按 capabilities 延迟初始化 MCP；mcp 列表变化时重建。
func EnsureInit() {
	initMu.Lock()
	defer initMu.Unlock()
	if config.Config == nil || !config.Config.MCP.Enabled {
		global = &Manager{clients: make(map[string]*stdioClient), routes: make(map[string]*toolRoute)}
		lastMCPKey = ""
		return
	}
	caps := brain.LoadActiveCapabilities()
	key := mcpCapsKey(caps)
	if global != nil && key == lastMCPKey {
		return
	}
	shutdownLocked()
	caps = brain.LoadActiveCapabilities()
	Init(config.Config.MCP, caps)
	lastMCPKey = mcpCapsKey(caps)
}

// ReinitIfNeeded 在 capabilities.yaml 的 mcp 段变化后重建（新 chat 连接时调用）。
func ReinitIfNeeded() {
	EnsureInit()
}

func shutdownLocked() {
	if global == nil {
		return
	}
	for name, c := range global.clients {
		if err := c.Close(); err != nil {
			log.Printf("MCP close %q: %v", name, err)
		}
	}
	global.clients = nil
	global.routes = nil
	global.llmTools = nil
	global = nil
}

func (mgr *Manager) connectServer(ctx context.Context, s config.MCPServerEntry) error {
	return connectServer(mgr, ctx, s)
}

func toLLMTool(t listedTool) llm.Tool {
	params := t.InputSchema
	if len(params) == 0 {
		params = json.RawMessage(`{"type":"object","properties":{}}`)
	}
	desc := strings.TrimSpace(t.Description)
	if desc == "" {
		desc = "MCP tool " + t.Name
	}
	return llm.Tool{
		Type: "function",
		Function: llm.ToolFunction{
			Name:        t.Name,
			Description: desc + " (via MCP browser)",
			Parameters:  params,
		},
	}
}

// Global 返回已初始化的 MCP 管理器（可能为 nil 或无工具）。
func Global() *Manager {
	EnsureInit()
	return global
}

// Tools 供 LLM API 注册的 MCP 工具列表。
func (mgr *Manager) Tools() []llm.Tool {
	if mgr == nil {
		return nil
	}
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()
	out := make([]llm.Tool, len(mgr.llmTools))
	copy(out, mgr.llmTools)
	return out
}

// TryCall 若 name 为 MCP 工具则执行并返回 ok=true。
func (mgr *Manager) TryCall(ctx context.Context, name, argsJSON string) (out string, err error, ok bool) {
	if mgr == nil {
		return "", nil, false
	}
	mgr.mu.RLock()
	route, exists := mgr.routes[name]
	var client *stdioClient
	if exists && route != nil {
		client = mgr.clients[route.serverName]
	}
	mgr.mu.RUnlock()
	if !exists || route == nil || client == nil {
		return "", nil, false
	}
	var args map[string]interface{}
	if strings.TrimSpace(argsJSON) == "" || argsJSON == "null" {
		args = map[string]interface{}{}
	} else if err := llm.ParseToolArguments(argsJSON, &args); err != nil {
		return "", fmt.Errorf("mcp args: %w", err), true
	}
	callCtx := ctx
	if mgr.timeout > 0 {
		var cancel context.CancelFunc
		callCtx, cancel = context.WithTimeout(ctx, mgr.timeout)
		defer cancel()
	}
	text, err := client.callTool(callCtx, route.toolName, args)
	if mgr.maxOutput > 0 && len(text) > mgr.maxOutput {
		text = text[:mgr.maxOutput] + "\n…(truncated)"
	}
	return text, err, true
}

// Shutdown 关闭所有 MCP 子进程。
func Shutdown() {
	initMu.Lock()
	defer initMu.Unlock()
	shutdownLocked()
	lastMCPKey = ""
}
