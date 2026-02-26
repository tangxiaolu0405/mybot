package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"mybot/internal/config"
)

// RunGit 在 brain 基目录执行 git 命令
func RunGit(args ...string) (string, error) {
	baseDir := config.GetBrainBaseDir()
	
	// 确保目录存在
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return "", fmt.Errorf("brain base directory does not exist: %s", baseDir)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = baseDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("git command failed: %w", err)
	}

	return string(output), nil
}

// InitGitRepo 初始化 git 仓库（如果不存在）
func InitGitRepo() error {
	baseDir := config.GetBrainBaseDir()
	gitDir := filepath.Join(baseDir, ".git")

	// 如果已经是 git 仓库，直接返回
	if _, err := os.Stat(gitDir); err == nil {
		return nil
	}

	_, err := RunGit("init")
	return err
}

// AddAll 添加所有文件到 git
func AddAll() error {
	_, err := RunGit("add", "-A")
	return err
}

// Commit 提交更改
func Commit(message string) error {
	_, err := RunGit("commit", "-m", message)
	return err
}

// Status 获取 git 状态
func Status() (string, error) {
	return RunGit("status", "--short")
}

// IsGitRepo 检查是否是 git 仓库
func IsGitRepo() bool {
	baseDir := config.GetBrainBaseDir()
	gitDir := filepath.Join(baseDir, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

// GetBrainBaseDir 获取 brain 基目录（供外部使用）
func GetBrainBaseDir() string {
	return config.GetBrainBaseDir()
}
