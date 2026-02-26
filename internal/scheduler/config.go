package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultConfigFile = ".cata/skills.json"
)

// SkillConfig 技能配置
type SkillConfig struct {
	Enabled map[string]bool `json:"enabled"` // skill name -> enabled
}

// LoadSkillConfig 加载技能配置
func LoadSkillConfig() (*SkillConfig, error) {
	configPath := getConfigPath()
	
	// 如果配置文件不存在，返回默认配置（所有技能启用）
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &SkillConfig{
			Enabled: make(map[string]bool),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config SkillConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.Enabled == nil {
		config.Enabled = make(map[string]bool)
	}

	return &config, nil
}

// SaveSkillConfig 保存技能配置
func SaveSkillConfig(config *SkillConfig) error {
	configPath := getConfigPath()
	
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath 获取配置文件路径（使用当前工作目录）
func getConfigPath() string {
	// 使用当前工作目录，而不是用户目录
	wd, err := os.Getwd()
	if err != nil {
		// fallback 到相对路径
		return DefaultConfigFile
	}
	return filepath.Join(wd, DefaultConfigFile)
}

// IsSkillEnabled 检查技能是否启用（默认启用）
func (sc *SkillConfig) IsSkillEnabled(skillName string) bool {
	if enabled, ok := sc.Enabled[skillName]; ok {
		return enabled
	}
	// 默认启用
	return true
}

// SetSkillEnabled 设置技能启用状态
func (sc *SkillConfig) SetSkillEnabled(skillName string, enabled bool) {
	if sc.Enabled == nil {
		sc.Enabled = make(map[string]bool)
	}
	sc.Enabled[skillName] = enabled
}
