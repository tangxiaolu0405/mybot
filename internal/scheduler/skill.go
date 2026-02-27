package scheduler

import (
	"context"
	"fmt"
)

// Skill 定义可插拔的技能接口
type Skill interface {
	// Name 返回技能名称
	Name() string

	// Run 执行技能，args 为命令行参数或定时触发的参数
	Run(ctx context.Context, args []string) error

	// CronSchedule 返回定时任务的 Cron 表达式，返回空字符串表示不定时执行
	CronSchedule() string

	// CLICommand 返回对应的 CLI 子命令，返回空字符串表示不提供 CLI 命令
	CLICommand() string
}

// SkillRegistry 技能注册表
type SkillRegistry struct {
	skills map[string]Skill
	config *SkillConfig
}

// NewSkillRegistry 创建新的技能注册表
func NewSkillRegistry() (*SkillRegistry, error) {
	config, err := LoadSkillConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load skill config: %w", err)
	}

	return &SkillRegistry{
		skills: make(map[string]Skill),
		config: config,
	}, nil
}

// Register 注册一个技能（如果启用）
func (r *SkillRegistry) Register(skill Skill) error {
	name := skill.Name()
	if name == "" {
		return fmt.Errorf("skill name cannot be empty")
	}

	// 检查是否启用
	if !r.config.IsSkillEnabled(name) {
		return nil // 未启用，不注册
	}

	if _, exists := r.skills[name]; exists {
		return fmt.Errorf("skill %s already registered", name)
	}

	r.skills[name] = skill
	return nil
}

// Get 根据名称获取技能
func (r *SkillRegistry) Get(name string) (Skill, bool) {
	skill, ok := r.skills[name]
	return skill, ok
}

// GetByCLICommand 根据 CLI 命令获取技能
func (r *SkillRegistry) GetByCLICommand(cmd string) (Skill, bool) {
	for _, skill := range r.skills {
		if skill.CLICommand() == cmd {
			return skill, true
		}
	}
	return nil, false
}

// List 列出所有已注册的技能
func (r *SkillRegistry) List() []Skill {
	result := make([]Skill, 0, len(r.skills))
	for _, skill := range r.skills {
		result = append(result, skill)
	}
	return result
}

// GetAllWithCron 获取所有有定时任务的技能
func (r *SkillRegistry) GetAllWithCron() []Skill {
	var result []Skill
	for _, skill := range r.skills {
		if skill.CronSchedule() != "" {
			result = append(result, skill)
		}
	}
	return result
}

// EnableSkill 启用技能
func (r *SkillRegistry) EnableSkill(name string) error {
	r.config.SetSkillEnabled(name, true)
	return SaveSkillConfig(r.config)
}

// DisableSkill 禁用技能
func (r *SkillRegistry) DisableSkill(name string) error {
	r.config.SetSkillEnabled(name, false)
	// 从注册表中移除
	delete(r.skills, name)
	return SaveSkillConfig(r.config)
}

// GetConfig 获取配置
func (r *SkillRegistry) GetConfig() *SkillConfig {
	return r.config
}
