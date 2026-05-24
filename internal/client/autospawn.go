package client

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"cata/internal/brain"
)

const spawnWaitTimeout = 20 * time.Second

// EnsureServer 若本机无 cata server 则后台拉起 `cata run --managed`。
func EnsureServer() error {
	if err := PingServer(); err == nil {
		return nil
	}
	return withSpawnLock(spawnManagedIfNeeded)
}

func spawnManagedIfNeeded() error {
	if err := PingServer(); err == nil {
		return nil
	}
	if err := startManagedServer(); err != nil {
		return err
	}
	deadline := time.Now().Add(spawnWaitTimeout)
	for time.Now().Before(deadline) {
		if err := PingServer(); err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("cata server did not become ready within %s", spawnWaitTimeout)
}

func startManagedServer() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, "run", "--managed")
	// 日志由子进程 runServer 归档后写入新 cata-server.log，避免与旧文件混写
	cmd.Stdout = nil
	cmd.Stderr = nil
	detachCmd(cmd)
	return cmd.Start()
}

func withSpawnLock(fn func() error) error {
	if err := os.MkdirAll(filepath.Join(brain.CataHome(), "locks"), 0755); err != nil {
		return err
	}
	path := filepath.Join(brain.CataHome(), "locks", "spawn.lock")
	for attempt := 0; attempt < 30; attempt++ {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			_, _ = fmt.Fprintf(f, "%d\n", os.Getpid())
			err := fn()
			_ = f.Close()
			_ = os.Remove(path)
			return err
		}
		if !os.IsExist(err) {
			return err
		}
		if err := PingServer(); err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("cata: timed out waiting for spawn lock")
}
