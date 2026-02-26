package evolution

import (
	"math/rand"
	"time"
)

// Task 任务
type Task struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`        // summarize, consolidate, recall, learn, optimize, reflect
	ActionPlan *ActionPlan            `json:"action_plan"` // 原始行动方案
	Params     map[string]interface{} `json:"params"`     // 执行参数
	Priority   int                    `json:"priority"`   // 优先级 1-10
	Status     string                 `json:"status"`     // "pending" | "running" | "completed" | "failed"
	CreatedAt  string                `json:"created_at"`
	StartedAt  string                `json:"started_at"`
	CompletedAt string               `json:"completed_at"`
	Result     *TaskResult            `json:"result"`
}

// ActionPlan 行动方案（由 LLM 生成）
type ActionPlan struct {
	Action         string   `json:"action"`          // 行动类型
	Reason         string   `json:"reason"`          // 决策理由
	Steps          []string `json:"steps"`           // 执行步骤
	ExpectedOutcome string  `json:"expected_outcome"` // 预期结果
	Priority       int      `json:"priority"`         // 优先级 1-10
}

// TaskResult 任务结果
type TaskResult struct {
	Success  bool                   `json:"success"`
	Output   string                 `json:"output"`
	Error    string                 `json:"error"`
	Metrics  map[string]interface{} `json:"metrics"`
	Learning string                 `json:"learning"` // 学到的经验
}

// NewTask 创建新任务
func NewTask(actionPlan *ActionPlan) *Task {
	return &Task{
		ID:         generateTaskID(),
		Type:       actionPlan.Action,
		ActionPlan: actionPlan,
		Params:     make(map[string]interface{}),
		Priority:   actionPlan.Priority,
		Status:     "pending",
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
}

// generateTaskID 生成任务 ID
func generateTaskID() string {
	return "task-" + time.Now().Format("20060102-150405") + "-" + randomString(6)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
