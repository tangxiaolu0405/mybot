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

var global *Manager

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
	for _, t := range tools {
		name := strings.TrimSpace(t.Name)
		if name == "" {
			continue
		}
		mgr.routes[name] = &toolRoute{serverName: s.Name, toolName: name}
		mgr.llmTools = append(mgr.llmTools, toLLMTool(t))
	}
	return nil
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
	client := mgr.clients[route.serverName]
	mgr.mu.RUnlock()
	if !exists || client == nil {
		return "", nil, false
	}
	var args map[string]interface{}
	if strings.TrimSpace(argsJSON) == "" || argsJSON == "null" {
		args = map[string]interface{}{}
	} else if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
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
	if global == nil {
		return
	}
	global.mu.Lock()
	defer global.mu.Unlock()
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
