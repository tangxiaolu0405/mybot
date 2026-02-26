package scheduler

import (
	"context"
	"fmt"
	"log"

	"mybot/internal/memory"
)

// DailyConsolidateSkill 今日固化技能
// 对应 brain/workflow.md 任务类型 consolidate：固化到 hot/archive。
// 定时触发：将当日输入/缓冲区经 LLM 整理后调用 memory.Consolidate 写入 brain/archive 或 hot.md（与 core.md 资源路径一致）。
type DailyConsolidateSkill struct {
	memMgr *memory.MemoryManager
}

func NewDailyConsolidateSkill(memMgr *memory.MemoryManager) *DailyConsolidateSkill {
	return &DailyConsolidateSkill{
		memMgr: memMgr,
	}
}

func (s *DailyConsolidateSkill) Name() string {
	return "daily-consolidate"
}

func (s *DailyConsolidateSkill) Run(ctx context.Context, args []string) error {
	log.Println("Running daily consolidate skill...")
	
	// TODO: 这里应该从缓冲区或外部输入获取当日内容
	// 目前简化实现：如果没有 LLM，可以手动调用 consolidate
	// 实际使用时，可以通过其他方式（如文件、API）获取当日输入
	
	// 示例：如果有缓冲区内容，可以这样处理
	// bufferContent := getBufferContent() // 需要实现
	// if bufferContent != "" {
	//     // 使用 LLM 整理内容（如果有）
	//     processedContent := processWithLLM(bufferContent) // 需要实现
	//     topic := fmt.Sprintf("每日总结 %s", time.Now().Format("2006-01-02"))
	//     return s.memMgr.Consolidate(topic, processedContent)
	// }
	
	log.Println("Daily consolidate: No buffer content to process (LLM integration pending)")
	return nil
}

func (s *DailyConsolidateSkill) CronSchedule() string {
	return "23:00" // 每日 23:00 执行
}

func (s *DailyConsolidateSkill) CLICommand() string {
	return "" // 不提供 CLI 命令
}

// PeriodicSummarizeSkill 周期摘要技能
// 对应 brain/workflow.md 任务类型 summarize：压缩 archive。
// 定时或阈值触发：多日 archive 合并摘要、写回 summary-YYYY-MM.md、更新 memory_index.json（与 core.md 档案与索引一致）。
type PeriodicSummarizeSkill struct {
	memMgr *memory.MemoryManager
}

func NewPeriodicSummarizeSkill(memMgr *memory.MemoryManager) *PeriodicSummarizeSkill {
	return &PeriodicSummarizeSkill{
		memMgr: memMgr,
	}
}

func (s *PeriodicSummarizeSkill) Name() string {
	return "periodic-summarize"
}

func (s *PeriodicSummarizeSkill) Run(ctx context.Context, args []string) error {
	log.Println("Running periodic summarize skill...")
	
	// 检查是否应该触发摘要
	shouldSummarize, reason := s.memMgr.CheckSummarizeTrigger()
	if !shouldSummarize {
		log.Printf("Summarize not needed: %s", reason)
		return nil
	}
	
	log.Printf("Summarize needed: %s", reason)
	
	// 调用 SummarizeAndRotate
	if err := s.memMgr.SummarizeAndRotate(); err != nil {
		return fmt.Errorf("failed to summarize and rotate: %w", err)
	}
	
	log.Println("Periodic summarize completed successfully")
	return nil
}

func (s *PeriodicSummarizeSkill) CronSchedule() string {
	return "02:00" // 每日 02:00 执行（在每日固化之后）
}

func (s *PeriodicSummarizeSkill) CLICommand() string {
	return "" // 不提供 CLI 命令
}
