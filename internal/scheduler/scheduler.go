package scheduler

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// Scheduler 任务调度器
type Scheduler struct {
	registry *SkillRegistry
	cron     *SimpleCron
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewScheduler 创建新的调度器
func NewScheduler(registry *SkillRegistry) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	cron := NewSimpleCron(ctx)
	return &Scheduler{
		registry: registry,
		cron:     cron,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start 启动调度器（注册定时任务并开始运行）
func (s *Scheduler) Start() error {
	// 注册所有有定时任务的技能
	skills := s.registry.GetAllWithCron()
	for _, skill := range skills {
		schedule := skill.CronSchedule()
		if schedule == "" {
			continue
		}

		// 解析 Cron 表达式（简化版：仅支持 "HH:MM" 格式，表示每日执行）
		hour, minute, err := parseDailySchedule(schedule)
		if err != nil {
			log.Printf("Warning: invalid schedule '%s' for skill %s: %v", schedule, skill.Name(), err)
			continue
		}

		// 创建技能副本以避免闭包问题
		skillCopy := skill
		err = s.cron.AddDailyTask(hour, minute, func() {
			log.Printf("Running scheduled skill: %s", skillCopy.Name())
			if err := skillCopy.Run(s.ctx, nil); err != nil {
				log.Printf("Error running skill %s: %v", skillCopy.Name(), err)
			}
		})
		if err != nil {
			return fmt.Errorf("failed to register cron for skill %s: %w", skill.Name(), err)
		}
		log.Printf("Registered daily schedule '%s' for skill: %s", schedule, skill.Name())
	}

	log.Println("Scheduler started")
	return nil
}

// parseDailySchedule 解析每日时间格式 "HH:MM"
func parseDailySchedule(schedule string) (hour, minute int, err error) {
	parts := strings.Split(schedule, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid schedule format, expected HH:MM")
	}

	var h, m int
	if _, err := fmt.Sscanf(parts[0], "%d", &h); err != nil {
		return 0, 0, fmt.Errorf("invalid hour: %w", err)
	}
	if _, err := fmt.Sscanf(parts[1], "%d", &m); err != nil {
		return 0, 0, fmt.Errorf("invalid minute: %w", err)
	}

	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, fmt.Errorf("hour must be 0-23, minute must be 0-59")
	}

	return h, m, nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cron.Stop()
	s.cancel()
	log.Println("Scheduler stopped")
}

// RunSkill 执行指定的技能（用于 CLI 调用）
func (s *Scheduler) RunSkill(skillName string, args []string) error {
	skill, ok := s.registry.Get(skillName)
	if !ok {
		return fmt.Errorf("skill not found: %s", skillName)
	}

	return skill.Run(s.ctx, args)
}

// RunSkillByCLICommand 根据 CLI 命令执行技能
func (s *Scheduler) RunSkillByCLICommand(cmd string, args []string) error {
	skill, ok := s.registry.GetByCLICommand(cmd)
	if !ok {
		return fmt.Errorf("CLI command not found: %s", cmd)
	}

	return skill.Run(s.ctx, args)
}

// Wait 等待调度器停止（阻塞）
func (s *Scheduler) Wait() {
	<-s.ctx.Done()
}

// IsRunning 检查调度器是否正在运行
func (s *Scheduler) IsRunning() bool {
	select {
	case <-s.ctx.Done():
		return false
	default:
		return true
	}
}
