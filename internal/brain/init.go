package brain

import (
	"fmt"
	"os"
)

// InitDirectory 创建 ~/.cata 布局、global 模板，并迁移旧版扁平 brain（若有）。
func InitDirectory() error {
	if err := EnsureCataLayout(); err != nil {
		return fmt.Errorf("cata layout: %w", err)
	}
	// 若从 cata init 调用，为当前目录注册 workspace
	if wd, err := os.Getwd(); err == nil {
		if _, err := ResolveWorkspace(wd); err != nil {
			return fmt.Errorf("brain: %w", err)
		}
	}
	return nil
}
