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

	"cata/internal/brain"
	"cata/internal/config"
	"cata/internal/execcmd"
)

type req struct {
	Command   string            `json:"command"`
	Text      string            `json:"text,omitempty"`
	Stream    bool              `json:"stream,omitempty"`
	ConfirmID string            `json:"confirm_id,omitempty"`
	Approved  bool              `json:"approved,omitempty"`
	Cwd       string            `json:"cwd,omitempty"`
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

	welcome()

	for {
		select {
		case <-sigCh:
			meta("\n")
			return
		default:
		}
		line, err := readLine()
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "/") {
			cmd := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(line), "/"))
			switch cmd {
			case "exit", "quit", "q":
				return
			case "clear", "reset":
				r, err := s.call(req{Command: "chat_reset"})
				if err != nil {
					errorMsg(err.Error())
					continue
				}
				if !r.Success {
					errorMsg(r.Message)
					continue
				}
				progressMsg(r.Message)
			case "config":
				meta("  config: %s%s%s\n", ansiYellow, config.GetConfigPath(), ansiReset)
			case "cls":
				meta("\033[H\033[2J")
			case "help":
				meta("  %scommands:%s\n", ansiBold, ansiReset)
				for _, c := range commands {
					meta("  %s/%s%s  %s%s%s\n", ansiBold, c.Name, ansiReset, ansiDim, c.Desc, ansiReset)
				}
			default:
				meta("  %sunknown:%s /%s (try /help)\n", ansiDim, ansiReset, cmd)
			}
			continue
		}

		outCwd, _ := os.Getwd()
		if err := s.write(req{Command: "chat", Text: line, Stream: true, Cwd: outCwd, Runtime: CollectRuntimeEnv()}); err != nil {
			errorMsg(err.Error())
			continue
		}
		if err := s.drainStream(); err != nil {
			errorMsg(err.Error())
			if connLost(err) {
				meta("  %s提示:%s 连接断开。直接发送下一条消息自动重连。\n", ansiDim, ansiReset)
				_ = s.conn.Close()
				_ = EnsureServer()
				if ns, derr := dial(); derr == nil {
					s.conn = ns.conn
					s.br = ns.br
					progressMsg("已重新连接 cata server")
				}
			}
		}
	}
}

func (s *session) drainStream() error {
	firstToken := true
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
			c, _ := ev["content"].(string)
			if c == "" {
				continue
			}
			if firstToken {
				firstToken = false
				outToken("\n")
			}
			outToken(c)

		case "thinking":
			c, _ := ev["content"].(string)
			if c != "" {
				meta("  %s…%s %s%s%s\n", ansiDim, ansiReset, ansiDim, truncate(c, 120), ansiReset)
			}

		case "progress":
			m, _ := ev["message"].(string)
			if m != "" {
				progressMsg(m)
			}

		case "tool_start":
			name, _ := ev["name"].(string)
			display, _ := ev["display"].(string)
			if name != "" {
				toolStart(name, display)
			}

		case "tool_result":
			out, _ := ev["output"].(string)
			name, _ := ev["name"].(string)
			display, _ := ev["display"].(string)
			if name == "run_command" {
				runCmdResult(s.lastExecCmd, s.lastExecCwd, out)
			} else if out != "" {
				toolOutput(name, out, display)
			}

		case "file_written":
			path, _ := ev["path"].(string)
			bytes, _ := ev["bytes"].(float64)
			fileWritten(path, int(bytes))

		case "diff":
			c, _ := ev["content"].(string)
			if c != "" {
				diffLine(c)
			}

		case "exec_confirm_required":
			id, _ := ev["confirm_id"].(string)
			cmd := execLine(ev)
			cwd, _ := ev["cwd"].(string)
			s.lastExecCmd = cmd
			s.lastExecCwd = cwd
			approved, err := confirmPrompt(id, cmd, cwd)
			if err != nil {
				return err
			}
			if err := s.write(req{Command: "exec_confirm", ConfirmID: id, Approved: approved}); err != nil {
				return err
			}
			if !approved {
				execDenied()
			}

		case "exec_denied":
			execDenied()

		case "exec_done":
			s.lastExecCmd = execLine(ev)
			if cwd, ok := ev["cwd"].(string); ok {
				s.lastExecCwd = cwd
			}
			exitCode := 0
			if ec, ok := ev["exit_code"].(float64); ok {
				exitCode = int(ec)
			}
			timedOut, _ := ev["timed_out"].(bool)
			execDone(s.lastExecCmd, exitCode, timedOut)

		case "error":
			m, _ := ev["message"].(string)
			errorMsg(m)

		case "user_choice":
			id, _ := ev["id"].(string)
			prompt, _ := ev["prompt"].(string)
			detail, _ := ev["detail"].(string)
			multi, _ := ev["multi"].(bool)
			rawOpts, _ := ev["options"].([]any)
			var opts []SelectOption
			for _, r := range rawOpts {
				if m, ok := r.(map[string]any); ok {
					opts = append(opts, SelectOption{
						ID:    str(m["id"]),
						Label: str(m["label"]),
						Desc:  str(m["desc"]),
					})
				}
			}
			if id == "" || len(opts) < 2 {
				errorMsg("invalid user_choice event")
				continue
			}
			var selected []string
			if multi {
				selected, _ = SelectMulti(prompt, detail, opts)
			} else {
				single, _ := Select(prompt, detail, opts)
				if single != "" {
					selected = []string{single}
				}
			}
			// Send response directly (not via s.write which wraps in req struct)
			type choiceResp struct {
				Command  string   `json:"command"`
				ChoiceID string   `json:"choice_id"`
				Selected []string `json:"selected"`
			}
			b, _ := json.Marshal(choiceResp{Command: "user_choice", ChoiceID: id, Selected: selected})
			s.conn.Write(append(b, '\n'))

		case "done":
			firstToken = true
			success, _ := ev["success"].(bool)
			if !success {
				return fmt.Errorf("chat failed")
			}
			outToken("\n")
			return nil
		}
	}
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
	return s[:n] + "…"
}

func str(v any) string {
	s, _ := v.(string)
	return s
}
