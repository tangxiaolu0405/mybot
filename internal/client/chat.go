package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"mybot/internal/brain"
	"mybot/internal/config"
	"mybot/internal/execcmd"
)

type req struct {
	Command   string `json:"command"`
	Text      string `json:"text,omitempty"`
	Stream    bool   `json:"stream,omitempty"`
	ConfirmID string `json:"confirm_id,omitempty"`
	Approved  bool   `json:"approved,omitempty"`
	Cwd       string `json:"cwd,omitempty"`
	Runtime   *brain.RuntimeEnv `json:"runtime,omitempty"`
}

type resp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type session struct {
	conn        net.Conn
	br          *bufio.Reader
	lastExecCmd string
	lastExecCwd string
}

func dial() (*session, error) {
	if err := config.InitBrainPath(); err != nil {
		return nil, err
	}
	conn, err := net.Dial("unix", config.ResolvedSocketPath())
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return &session{conn: conn, br: bufio.NewReader(conn)}, nil
}

func (s *session) write(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(append(b, '\n'))
	return err
}

func (s *session) readLine() ([]byte, error) {
	line, err := s.br.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(line), nil
}

func (s *session) call(r req) (resp, error) {
	if err := s.write(r); err != nil {
		return resp{}, err
	}
	line, err := s.readLine()
	if err != nil {
		return resp{}, err
	}
	var out resp
	return out, json.Unmarshal(line, &out)
}

// RunChat 启动终端交互（默认 cata / cata chat）。
func RunChat() {
	if err := config.InitBrainPath(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	release, err := AcquireOutputLock(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer release()

	if err := EnsureServer(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s, err := dial()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer s.conn.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	fmt.Println("cata — /clear reset history, /exit quit, \"\"\" … \"\"\" for multiline")

	sc := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-sigCh:
			return
		default:
		}
		fmt.Print("› ")
		if !sc.Scan() {
			break
		}
		line := readInput(sc, sc.Text())
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "/") {
			cmd := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(line), "/"))
			switch cmd {
			case "exit", "quit", "q":
				return // 断开连接；managed server 在最后一个客户端退出后自动停止
			case "clear", "reset":
				r, err := s.call(req{Command: "chat_reset"})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				if !r.Success {
					fmt.Fprintln(os.Stderr, r.Message)
					continue
				}
				fmt.Println(r.Message)
			default:
				fmt.Fprintf(os.Stderr, "unknown: /%s (try /clear, /exit)\n", cmd)
			}
			continue
		}

		outCwd, _ := os.Getwd()
		if err := s.write(req{Command: "chat", Text: line, Stream: true, Cwd: outCwd, Runtime: CollectRuntimeEnv()}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if err := s.drainStream(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			if connLost(err) {
				fmt.Fprintln(os.Stderr, "提示: 与 cata server 的连接已断开。请直接再发一条消息（将自动重连 server），勿在 › 下输入 1 确认命令。")
				_ = s.conn.Close()
				_ = EnsureServer()
				if ns, derr := dial(); derr == nil {
					s.conn = ns.conn
					s.br = ns.br
					fmt.Fprintln(os.Stderr, "# 已重新连接 cata server")
				}
			}
		}
	}
}

func readInput(sc *bufio.Scanner, first string) string {
	if strings.TrimSpace(first) != `"""` {
		return first
	}
	var b strings.Builder
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimSpace(line) == `"""` {
			break
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(line)
	}
	return b.String()
}

func (s *session) drainStream() error {
	for {
		line, err := s.readLine()
		if err != nil {
			return err
		}
		if len(line) == 0 {
			continue
		}
		var ev map[string]any
		if err := json.Unmarshal(line, &ev); err != nil {
			return err
		}
		switch ev["type"] {
		case "token":
			if c, _ := ev["content"].(string); c != "" {
				fmt.Print(c)
			}
		case "progress":
			if m, _ := ev["message"].(string); m != "" {
				fmt.Fprintf(os.Stderr, "# %s\n", m)
			}
		case "tool_start":
			if n, _ := ev["name"].(string); n != "" && n != "run_command" {
				fmt.Fprintf(os.Stderr, "▶ %s\n", n)
			}
		case "tool_result":
			out, _ := ev["output"].(string)
			name, _ := ev["name"].(string)
			if name == "run_command" {
				printExecBlock(s.lastExecCmd, s.lastExecCwd, out)
			} else if out != "" {
				fmt.Fprintf(os.Stderr, "── %s ──\n%s\n", name, truncate(out, 4000))
			}
		case "file":
			fmt.Fprintf(os.Stderr, "── %v ──\n%v\n", ev["path"], ev["content"])
		case "file_written":
			fmt.Fprintf(os.Stderr, "── wrote %v (%v bytes) ──\n", ev["path"], ev["bytes"])
		case "exec_confirm_required":
			s.lastExecCmd = execLine(ev)
			s.lastExecCwd, _ = ev["cwd"].(string)
			if err := s.confirmExec(ev); err != nil {
				return err
			}
		case "exec_denied":
			fmt.Fprintln(os.Stderr, "── command cancelled ──")
		case "exec_done":
			s.lastExecCmd = execLine(ev)
			if cwd, ok := ev["cwd"].(string); ok {
				s.lastExecCwd = cwd
			}
		case "error":
			fmt.Fprintf(os.Stderr, "error: %v\n", ev["message"])
		case "done":
			fmt.Println()
			if ok, _ := ev["success"].(bool); !ok {
				return fmt.Errorf("chat failed")
			}
			return nil
		}
	}
}

func (s *session) confirmExec(ev map[string]any) error {
	id, _ := ev["confirm_id"].(string)
	if id == "" {
		return fmt.Errorf("missing confirm_id")
	}
	cmd, cwd := execLine(ev), ""
	if c, ok := ev["cwd"].(string); ok {
		cwd = c
	}
	fmt.Fprintf(os.Stderr, "\nrun?  $ %s\n", cmd)
	if cwd != "" {
		fmt.Fprintf(os.Stderr, "cwd: %s\n", cwd)
	}
	fmt.Fprint(os.Stderr, "[1] run  [2] cancel › ")
	var choice string
	fmt.Scanln(&choice)
	approved := strings.TrimSpace(choice) == "1"
	return s.write(req{Command: "exec_confirm", ConfirmID: id, Approved: approved})
}

func execLine(ev map[string]any) string {
	if s, ok := ev["command_line"].(string); ok && s != "" {
		return s
	}
	argv, ok := ev["argv"].([]any)
	if !ok {
		return ""
	}
	parts := make([]string, 0, len(argv))
	for _, a := range argv {
		if s, ok := a.(string); ok {
			parts = append(parts, s)
		}
	}
	return execcmd.FormatLine(parts)
}

func printExecBlock(cmd, cwd, output string) {
	fmt.Fprintln(os.Stderr, "\n── command output ──")
	if cmd != "" {
		fmt.Fprintf(os.Stderr, "  $ %s\n", cmd)
	}
	if cwd != "" {
		fmt.Fprintf(os.Stderr, "  cwd: %s\n", cwd)
	}
	if strings.TrimSpace(output) != "" {
		fmt.Fprintf(os.Stderr, "\n%s\n", strings.TrimRight(output, "\n"))
	}
}

func connLost(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "forcibly closed") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "use of closed network")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "\n…"
}
