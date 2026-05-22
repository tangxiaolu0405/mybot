package client

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mybot/internal/brain"
	"mybot/internal/clock"
)

type lockMeta struct {
	PID       int    `json:"pid"`
	OutputCwd string `json:"output_cwd"`
	StartedAt string `json:"started_at"`
}

// AcquireOutputLock 同一产出区（绝对 cwd）仅允许一个 cata chat；不同目录可并行。
func AcquireOutputLock(outputCwd string) (release func(), err error) {
	abs, err := filepath.Abs(strings.TrimSpace(outputCwd))
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(brain.CataHome(), "locks"), 0755); err != nil {
		return nil, err
	}
	path := lockFilePath(abs)

	for attempt := 0; attempt < 2; attempt++ {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
		if err == nil {
			meta := lockMeta{
				PID:       os.Getpid(),
				OutputCwd: abs,
				StartedAt: clock.RFC3339(),
			}
			b, _ := json.Marshal(meta)
			_, _ = f.Write(b)
			return func() {
				_ = f.Close()
				_ = os.Remove(path)
			}, nil
		}
		if !os.IsExist(err) {
			return nil, fmt.Errorf("lock %s: %w", path, err)
		}
		meta, readErr := readLockMeta(path)
		if readErr == nil && meta.PID > 0 && processAlive(meta.PID) {
			return nil, fmt.Errorf("cata: 产出区已有会话（%s，PID %d）。请先退出后再开，或换目录", abs, meta.PID)
		}
		_ = os.Remove(path)
	}
	return nil, fmt.Errorf("cata: 无法获取产出区锁：%s", abs)
}

func lockFilePath(absCwd string) string {
	h := sha256.Sum256([]byte(filepath.Clean(absCwd)))
	name := "out_" + hex.EncodeToString(h[:8]) + ".lock"
	return filepath.Join(brain.CataHome(), "locks", name)
}

func readLockMeta(path string) (lockMeta, error) {
	var m lockMeta
	b, err := os.ReadFile(path)
	if err != nil {
		return m, err
	}
	return m, json.Unmarshal(b, &m)
}
