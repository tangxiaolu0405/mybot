// Package scheduler 中与 brain/core.md「技能」及 skills/skills-index.json 对齐的技能索引加载。

package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"mybot/internal/brain"
)

// SkillsIndex 与 skills/skills-index.json 结构对齐（brain/core.md：从 skills-index 解析，按需读 SKILL.md）。
type SkillsIndex struct {
	Version    string       `json:"version"`
	GeneratedAt string      `json:"generated_at"`
	Skills     []SkillMeta  `json:"skills"`
	TagsIndex  map[string][]string `json:"tags_index"`
}

// SkillMeta 技能元数据，与 skills-index 中单条一致。
type SkillMeta struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Description  string   `json:"description"`
	Version      string   `json:"version"`
	Author       string   `json:"author"`
	Tags         []string `json:"tags"`
	Dependencies []string `json:"dependencies"`
}

// SkillsIndexLoader 加载并缓存 skills-index.json，与 brain 路径一致。
type SkillsIndexLoader struct {
	path string
	mu   sync.RWMutex
	idx  *SkillsIndex
}

// NewSkillsIndexLoader 使用 brain/core.md 定义的技能索引路径创建加载器。
func NewSkillsIndexLoader() *SkillsIndexLoader {
	return &SkillsIndexLoader{
		path: brain.SkillsIndexPath(),
	}
}

// Load 读取并解析 skills-index.json，若文件不存在返回 nil 与 nil error。
func (l *SkillsIndexLoader) Load() (*SkillsIndex, error) {
	data, err := os.ReadFile(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read skills index: %w", err)
	}

	var idx SkillsIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse skills index: %w", err)
	}

	l.mu.Lock()
	l.idx = &idx
	l.mu.Unlock()

	return &idx, nil
}

// Get 返回已加载的索引（若未加载则先 Load）。
func (l *SkillsIndexLoader) Get() (*SkillsIndex, error) {
	l.mu.RLock()
	idx := l.idx
	l.mu.RUnlock()

	if idx != nil {
		return idx, nil
	}

	return l.Load()
}

// SkillNames 返回索引中所有技能名称（与 core.md 技能列举一致）。
func (l *SkillsIndexLoader) SkillNames() ([]string, error) {
	idx, err := l.Get()
	if err != nil || idx == nil {
		return nil, err
	}

	names := make([]string, 0, len(idx.Skills))
	for _, s := range idx.Skills {
		names = append(names, s.Name)
	}
	return names, nil
}

// SkillByName 按名称查找技能元数据。
func (l *SkillsIndexLoader) SkillByName(name string) (*SkillMeta, error) {
	idx, err := l.Get()
	if err != nil || idx == nil {
		return nil, err
	}

	for i := range idx.Skills {
		if idx.Skills[i].Name == name {
			return &idx.Skills[i], nil
		}
	}
	return nil, nil
}

// SkillsDir 返回技能目录路径（与 brain 一致，供 loader 扫描或按 path 读 SKILL.md）。
func SkillsDirFromBrain() string {
	return brain.SkillsDir()
}
