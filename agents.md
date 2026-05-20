# CLAW / MyBot — 项目级 AI 上下文（`agents.md`）

本文档定义**终端优先个人 Agent** 的目标与实现边界，与当前 Go 代码对齐。

# 最重要的事情：严格遵守第一性原理
---

## 愿景

**终端原生 AI 助手**：编排与记忆可审计、可 fork；推理外置为可配置 HTTP LLM；主入口为 `cata`（默认 `chat`）流式对话。

---

## 目录边界（重要）

见 **`brain/directory-plan.md`**。

| 位置 | 角色 |
|------|------|
| **`~/.cata/`** | **脑子**：记忆、persona、演进、config、socket |
| **当前工作目录 cwd** | **产出区**：代码与命令结果；`run_command` 在此 |
| **项目内 `.cata/`** | 脑子门牌（yaml/link），不是脑子正文 |
| **仓库 `brain/`** | 脑子模板种子 |

`focus_path`（git 根 / yaml 目录 / cwd）只用于**选中哪一格脑子**，不把产出存进 `~/.cata`。

---

## 当前实现

| 区域 | 作用 |
|------|------|
| **`cmd/cata`** | `init`（初始化 ~/.cata/brain）、`config`、`run`（socket + 后台演进） |
| **`cmd/cata`（`chat`）** | 默认流式 LLM 客户端；协议：`chat`（`stream:true`）、`chat_reset`、`ping`；`cmd/catacli` 已废弃 |
| **`internal/server`** | Unix socket、终端 chat 工具循环 |
| **`internal/llm`** | OpenAI 兼容 Chat；注入 boot-leader + **路径块**（脑子 vs 产出区）+ 脑子节选 |
| **`internal/brain`** | 路径常量、`InitDirectory`、终端节选 |
| **`internal/evolve`** | **仅后台**自主演进：Observe → LLM → 文档补丁 → `evolution_log.json`（无手动 CLI） |
| **`internal/config`** | `~/.cata/config.json`：LLM、exec、`evolution.enabled` / `cycle_interval` |

**已移除（方案 B）**：`internal/memory`、`internal/evolution`（旧任务引擎）、`internal/scheduler`、`internal/git`、`skills/` 服务端加载。

---

## 记忆流（摘要）

| 层 | 写入方 |
|----|--------|
| Socket 会话历史 | server（内存，chat_reset 清空） |
| `memory/short/current.md`（每格脑子） | **每轮 cata chat 成功后 server 规则追加**（`session_memory.go`） |
| mode `hot.md` 等 | **仅** `internal/evolve` 从 short-term 提炼 |
| `long-term/`、`archive/` | evolve |

详见 **`brain/core.md` §记忆分层**。

## 自主演进（摘要）

- **触发**：short-term 有新内容等门控（见 `internal/evolve`）；默认周期 600s。
- **无** 手动 `cata evolve` 命令。

---

## MCP 与 Skill（已接入）

- **MCP browser**：`~/.cata/config.json` → `mcp.servers`（默认 `npx -y @playwright/mcp@latest`，name=`browser`）；`modes/*/capabilities.yaml` 用 `mcp: [browser]` 启用。
- **Skill（提示词型）**：`capabilities.yaml` → `skills: [name]`；从 `~/.cata/skills/`、项目 `skills/`、或 `~/.cursor/skills-cursor/` 加载 `SKILL.md` 注入 system（非 API tool）。
- **执行型 skill**（脚本）：暂未接，等同后续 `run_command` / 子进程扩展。

## 刻意排除

- **旧 `skills/` 服务端调度、`scripts/` 主线**：已废弃；仅保留 MD 提示词加载。
- **手动演进命令、任务队列、MemoryManager 索引**：已废弃。

---

## 给 AI 的约束

1. 改功能先看 **`internal/server`**、**`cmd/cata`**、**`internal/client`**、**`internal/evolve`**。
2. 路径以 **`internal/brain/paths.go`**、`context_paths.go` 与 **`~/.cata/global/`** 为准；产出区 = `cata` 启动时的 cwd。
3. **同机一个 server**（`cata` 自动 `run --managed` 或手动 `cata run`）；**同一产出区目录只能开一个 chat**；**最后一个 chat 断开**后 managed server 自动退出。
4. 勿虚构路径；勿把仓库 `brain/`（模板）与 `~/.cata`（脑子）混为一谈；勿把 focus_path 当成产出区。

---

## 建议阅读顺序

1. 本文件  
2. `~/.cata/brain/core.md`（或仓库模板 `brain/core.md`）  
3. `brain/workflow.md`  
4. 具体 `internal/server`、`internal/evolve` 源码  
