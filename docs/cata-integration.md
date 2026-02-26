# Cata 服务与 CLI

Cata 是实现了 brain/core.md 与 brain/workflow.md 规定流程的常驻服务：记忆管理、自主演进循环、任务队列。**catacli 仅用于发布任务与查看**，其余能力（记忆检索、固化、摘要、演进、技能等）由 **cataserver 内 LLM 自主决策**完成。

## 代码位置（已并入本仓库）

- **服务入口**：`cmd/cata/`（cata init / config / run / stop / upgrade / test）
- **客户端**：`cmd/catacli/`（仅 task create/list/status、ping）
- **实现**：`internal/config`、`internal/evolution`、`internal/memory`、`internal/llm`、`internal/server`、`internal/scheduler`、`internal/git`

构建与运行需在**项目根目录**（含 `go.mod`、`brain/`）下进行。

## 构建与运行

```bash
# 在项目根目录
cd /path/to/mybot

# 构建
go build -o cata ./cmd/cata
go build -o catacli ./cmd/catacli

# 初始化 brain 与配置（首次）
./cata init

# 启动服务（前台）
./cata run
```

另开终端，用 catacli **发布任务与查看**（会自动从当前目录向上查找项目根）：

```bash
./catacli ping
./catacli task create "帮我整理本周记忆摘要"
./catacli task create summarize --async
./catacli task list
./catacli task status <task-id>
```

## 与 core/workflow 的对应

- **核心**：brain/（core.md、memory、可选 hot.md、archive、memory_index.json、evolution_log.json、task_queue.json）。
- **自主演进**：分析状态 → LLM 决策 → 记录 → 生成任务 → 执行 → 学习反馈，见 brain/workflow.md；服务内由 `internal/evolution` 实现；**不需要通过 CLI 触发**，由服务端自主执行。
- **任务类型**：summarize / consolidate / recall / learn / optimize / reflect / idle，语义与 workflow 一致；用户通过 `task create "<需求>"` 或 `task create <type>` 下发，由 cata 解析并执行。
- **技能**：Agent 侧 skills 与 Cata 服务共用同一套 brain 数据与流程约定；技能的启用/调度由服务端内部决策。

## CLI 命令参考（仅发布任务与查看）

| 命令 | 语义 |
|------|------|
| task create \"\<需求描述>\" [--async] | 按需求创建任务，cata 解析并执行；--async 表示仅入队由服务自动执行 |
| task create \<type> [args...] [--async] | 按类型创建：summarize, consolidate, recall, learn, optimize, reflect, idle, integrate |
| task list | 任务列表 |
| task status \<task-id> | 任务状态 |
| ping | 连通检查 |
| help | 帮助 |

配置：`cata config show/get/set`，或编辑项目根目录下 `.cata/config.json`；环境变量见 `cata config show`。
