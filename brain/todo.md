# 路线图（规划 · 当前未实现）

本文档记录**后续**方向；当前仓库仅实现 **单 cata 守护进程 + catacli 终端对话 + 后台自主演进（文档补丁）**。

## 多 CLI

- [ ] 除 `catacli` 外的其它终端前端（脚本化、管道、非交互）
- [ ] 统一的 socket 协议版本与能力协商

## 多 Agent

- [ ] 多个 Agent 角色/模型配置并行（不同 `llm.models` 角色已预留，编排未做）
- [ ] Agent 间分工：对话 Agent vs 演进 Agent vs 专用工具 Agent（演进已独立为 `internal/evolve`）
- [ ] 外部 IDE Agent（如 Cursor）与 cata 共享 `~/.cata/brain` 的协作约定

## 多任务

- [ ] 任务队列与异步作业（历史 `task_queue.json` 已废弃）
- [ ] 跨会话任务状态、取消、优先级
- [ ] 与 CI / 定时器集成的触发器（非 cata 内置 cron）

## 记忆增强（可选）

- [ ] `memory_index.json` 与检索 API
- [ ] archive 自动 summarize 阈值任务
- [x] 对话结束自动写 short-term（`internal/brain/session_memory.go`）
- [ ] short-term → hot 的 LLM 摘要质量（可选每 N 轮用轻量模型压缩再追加）

---

实现任一項前请先更新本文件与 `agents.md`，避免与「终端优先、文档补丁演进」主线冲突。
