package evolution

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"mybot/internal/brain"
)

var (
	TaskQueueFilePath string
)

// initQueuePaths 与 brain/core.md 任务队列路径一致
func initQueuePaths() {
	TaskQueueFilePath = brain.TaskQueuePath()
}

func init() {
	initQueuePaths()
}

// TaskQueue 任务队列
type TaskQueue struct {
	mu    sync.RWMutex
	tasks []*QueuedTask
}

// QueuedTask 队列中的任务
type QueuedTask struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	ActionPlan  *ActionPlan            `json:"action_plan"`
	Params      map[string]interface{} `json:"params"`
	Priority    int                    `json:"priority"`
	Status      string                 `json:"status"` // "pending" | "running" | "completed" | "failed"
	CreatedAt   string                 `json:"created_at"`
	CreatedBy   string                 `json:"created_by"` // "user" | "system"
	StartedAt   string                 `json:"started_at"`
	CompletedAt string                 `json:"completed_at"`
	Result     *TaskResult             `json:"result"`
}

// NewTaskQueue 创建任务队列
func NewTaskQueue() *TaskQueue {
	queue := &TaskQueue{
		tasks: []*QueuedTask{},
	}
	queue.load()
	return queue
}

// Enqueue 添加任务到队列
func (q *TaskQueue) Enqueue(actionPlan *ActionPlan, createdBy string) (*QueuedTask, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	task := &QueuedTask{
		ID:         generateTaskID(),
		Type:       actionPlan.Action,
		ActionPlan: actionPlan,
		Params:     make(map[string]interface{}),
		Priority:   actionPlan.Priority,
		Status:     "pending",
		CreatedAt:  time.Now().Format(time.RFC3339),
		CreatedBy:  createdBy,
	}

	q.tasks = append(q.tasks, task)
	q.save()

	return task, nil
}

// Dequeue 从队列中取出优先级最高的待执行任务
func (q *TaskQueue) Dequeue() *QueuedTask {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 找到优先级最高且状态为 pending 的任务
	var bestTask *QueuedTask
	maxPriority := -1

	for _, task := range q.tasks {
		if task.Status == "pending" && task.Priority > maxPriority {
			maxPriority = task.Priority
			bestTask = task
		}
	}

	if bestTask != nil {
		bestTask.Status = "running"
		bestTask.StartedAt = time.Now().Format(time.RFC3339)
		q.save()
	}

	return bestTask
}

// GetPendingTasks 获取所有待执行的任务
func (q *TaskQueue) GetPendingTasks() []*QueuedTask {
	q.mu.RLock()
	defer q.mu.RUnlock()

	pending := []*QueuedTask{}
	for _, task := range q.tasks {
		if task.Status == "pending" {
			pending = append(pending, task)
		}
	}
	return pending
}

// GetTask 根据 ID 获取任务
func (q *TaskQueue) GetTask(taskID string) *QueuedTask {
	q.mu.RLock()
	defer q.mu.RUnlock()

	for _, task := range q.tasks {
		if task.ID == taskID {
			return task
		}
	}
	return nil
}

// UpdateTask 更新任务状态
func (q *TaskQueue) UpdateTask(taskID string, status string, result *TaskResult) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, task := range q.tasks {
		if task.ID == taskID {
			task.Status = status
			task.Result = result
			if status == "completed" || status == "failed" {
				task.CompletedAt = time.Now().Format(time.RFC3339)
			}
			q.save()
			return nil
		}
	}
	return fmt.Errorf("task not found: %s", taskID)
}

// ListTasks 列出所有任务（可选过滤状态）
func (q *TaskQueue) ListTasks(statusFilter string, limit int) []*QueuedTask {
	q.mu.RLock()
	defer q.mu.RUnlock()

	result := []*QueuedTask{}
	count := 0

	// 从后往前遍历（最新的在前）
	for i := len(q.tasks) - 1; i >= 0 && count < limit; i-- {
		task := q.tasks[i]
		if statusFilter == "" || task.Status == statusFilter {
			result = append(result, task)
			count++
		}
	}

	return result
}

// load 从文件加载任务队列
func (q *TaskQueue) load() {
	data, err := os.ReadFile(TaskQueueFilePath)
	if err != nil {
		// 文件不存在，使用空队列
		return
	}

	var tasks []*QueuedTask
	if err := json.Unmarshal(data, &tasks); err != nil {
		// 解析失败，使用空队列
		return
	}

	q.tasks = tasks
}

// save 保存任务队列到文件
func (q *TaskQueue) save() error {
	data, err := json.MarshalIndent(q.tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(TaskQueueFilePath, data, 0644)
}
