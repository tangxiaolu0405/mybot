package brain

import (
	"fmt"
	"os"
	"strings"
)

// ReadPromptFile 读取 brain 目录下的提示词文本（去首尾空白）；空文件视为错误。
func ReadPromptFile(relPath string) (string, error) {
	p := Path(relPath)
	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("read brain prompt %s: %w", relPath, err)
	}
	s := strings.TrimSpace(string(data))
	if s == "" {
		return "", fmt.Errorf("brain prompt %s is empty", relPath)
	}
	return s, nil
}
