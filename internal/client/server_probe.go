package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"mybot/internal/config"
)

// PingServer 检测本机 cata.sock 是否已有存活 server。
func PingServer() error {
	if err := config.InitBrainPath(); err != nil {
		return err
	}
	path := config.ResolvedSocketPath()
	conn, err := net.DialTimeout("unix", path, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	req, _ := json.Marshal(map[string]string{"command": "ping"})
	if _, err := conn.Write(append(req, '\n')); err != nil {
		return err
	}
	br := bufio.NewReader(conn)
	line, err := br.ReadBytes('\n')
	if err != nil {
		return err
	}
	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(line), &resp); err != nil {
		return err
	}
	if !resp.Success || resp.Message != "pong" {
		return fmt.Errorf("unexpected ping response: %s", string(line))
	}
	return nil
}
