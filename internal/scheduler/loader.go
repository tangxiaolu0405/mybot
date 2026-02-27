package scheduler

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

const (
	DefaultSkillsDir = "skills"
)

// SkillLoader 技能加载器
type SkillLoader struct {
	skillsDir string
	registry  *SkillRegistry
}

// NewSkillLoader 创建技能加载器
func NewSkillLoader(skillsDir string, registry *SkillRegistry) *SkillLoader {
	if skillsDir == "" {
		skillsDir = DefaultSkillsDir
	}
	return &SkillLoader{
		skillsDir: skillsDir,
		registry:  registry,
	}
}

// LoadSkills 加载所有技能
func (sl *SkillLoader) LoadSkills() error {
	// 检查目录是否存在
	if _, err := os.Stat(sl.skillsDir); os.IsNotExist(err) {
		// 目录不存在，创建它
		if err := os.MkdirAll(sl.skillsDir, 0755); err != nil {
			return fmt.Errorf("failed to create skills directory: %w", err)
		}
		log.Printf("Created skills directory: %s", sl.skillsDir)
		return nil // 没有技能可加载
	}

	// 扫描目录
	return filepath.WalkDir(sl.skillsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和非 .so 文件（Go plugin）
		if d.IsDir() || !strings.HasSuffix(path, ".so") {
			return nil
		}

		// 尝试加载 plugin
		if err := sl.loadPlugin(path); err != nil {
			log.Printf("Warning: failed to load skill from %s: %v", path, err)
			// 继续加载其他技能，不中断
			return nil
		}

		return nil
	})
}

// loadPlugin 加载 Go plugin
func (sl *SkillLoader) loadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// 查找 NewSkill 函数
	symbol, err := p.Lookup("NewSkill")
	if err != nil {
		return fmt.Errorf("plugin does not export NewSkill: %w", err)
	}

	// 调用 NewSkill 函数
	newSkillFunc, ok := symbol.(func() (Skill, error))
	if !ok {
		return fmt.Errorf("NewSkill has wrong signature")
	}

	skill, err := newSkillFunc()
	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}

	// 注册技能
	if err := sl.registry.Register(skill); err != nil {
		return fmt.Errorf("failed to register skill: %w", err)
	}

	log.Printf("Loaded skill: %s from %s", skill.Name(), path)
	return nil
}

// LoadBuiltinSkills 加载内置技能（不通过 plugin）
func (sl *SkillLoader) LoadBuiltinSkills() error {
	// 注意：这里需要 MemoryManager，但目前 loader 没有访问权限
	// 实际应该在 Server 中注册内置技能
	// 这里先返回空，内置技能在 Server.Start() 中注册
	return nil
}
