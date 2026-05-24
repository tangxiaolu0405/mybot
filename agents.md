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
| **`internal/config`** | `~/.cata/config.json`：LLM（`deepseek` / `qwen` 等 OpenAI 兼容）、exec、`evolution.enabled` / `cycle_interval` |

**已移除（方案 B）**：`internal/memory`、`internal/evolution`（旧任务引擎）、`internal/scheduler`、`internal/git`、`skills/` 服务端加载。

---

## 记忆流（摘要）

| 层 | 写入方 |
|----|--------|
| Socket 会话历史 | server（内存，chat_reset 清空） |
| `memory/short/current.md`（每格脑子） | **每轮 cata chat 成功后 server 规则追加**（`session_memory.go`） |
| mode `hot.md` 等 | **仅** `internal/evolve` 从 short-term 提炼 |
| `long-term/`、`archive/` | evolve |

详见 **`brain/constraints.md` §记忆分层**。

## 自主演进（摘要）

- **触发**：short-term 有新内容等门控（见 `internal/evolve`）；默认周期 600s。
- **无** 手动 `cata evolve` 命令。

---

## LLM（DeepSeek）

- **provider**：`deepseek`（[OpenAI 兼容](https://api-docs.deepseek.com/zh-cn/)，代码走 `OpenAIProvider`）
- **默认**：`https://api.deepseek.com/chat/completions`，模型 `deepseek-v4-flash`（更强用 `deepseek-v4-pro`）
- **密钥**：`llm.api_key` 或 `DEEPSEEK_API_KEY`
- **`llm.thinking`**：`auto`（默认，有 tools 时 `disabled`，避免 tool 轮次 400）、`enabled`、`disabled`
- 思考模式 + tool 调用时须回传 `reasoning_content`（已实现）；见 [Thinking Mode](https://api-docs.deepseek.com/guides/thinking_mode)
- 原千问配置备份在 `~/.cata/config.json` → `llm_previous_qwen`（不参与加载）

## MCP 与 Skill（已接入）

- **MCP browser**：`~/.cata/config.json` → `mcp.servers`（默认 `npx -y @playwright/mcp@latest`，name=`browser`）；`modes/*/capabilities.yaml` 用 `mcp: [browser]` 启用。
- **Skill**：`capabilities.yaml` → `skills: [id]`；`SKILL.md` 查找顺序：`brain/workspaces/<ws>/skills/` → `~/.cata/skills/` → `~/.cursor/skills-cursor/`。
- **run_skill**：执行脑子内 `skills/<id>/` 的 `manifest.yaml` + 脚本（cwd=产出区）；由演进 `crystallize_skill` 固化；**不删** `mcp: [browser]`。
- **crystallize_skill**：高 token / 重复 browser 任务后，evolve 写 `skills/<id>/` 并自动 append capabilities；下次 chat 生效。

## 产出区（Output Area / Workspace）

**原则**：程序不一定运行在产出区。用户可在任意目录启动 `cata`，通过 `--dir` 指定干活的目标目录。

```
cata chat                         # 产出区 = 当前目录（默认）
cata chat --dir ~/project         # 产出区 = ~/project
cata chat --dir ~/a --dir ~/b     # 多产出区，第一个是主产出区
```

**配置文件** (`~/.cata/config.json`)：

```json
{ "workspace": { "default_dir": "~/myproject" } }
```

**参考 Claude Code**：
- Claude 的 launch dir → cata 的 cwd 或第一个 `--dir`
- Claude 的 `--add-dir` → cata 的多个 `--dir`
- Claude 不允许 `--cwd` 改变工作目录；cata 走自己的 `--dir` 方案

**脑子绑定**：
- 产出区用于文件工具和 `run_command` 的操作范围
- 脑子格子（`workspaces/<id>/`）由 `focus_path`（git root / workspace.yaml）决定
- 多产出区共享同一个脑子（如果它们属于同一个 git 项目）
- 不同产出区可并行开多个 chat（各自独立 output lock）

**实现要点**：
- Client 将 `--dir` 解析后的路径作为 `cwd` 发给 Server
- Server 的 `ResolveWorkspace(cwd)` 逻辑不变
- 文件工具 `safePathUnder` 检查每个产出区都在允许范围内

---

## 交互层：事件输出规范

**核心原则**：`stdout` = AI 的回答正文（可被管道/重定向）。`stderr` = 元信息（工具、进度、错误）。两者绝不混淆。

### 事件类型与显示

| 事件 | 通道 | 显示条件 | 说明 |
|------|------|----------|------|
| `token` | **stdout** | 始终 | AI 生成的文本，逐字流式输出 |
| `thinking` | stderr | `--show-thinking` | DeepSeek reasoning_content |
| `tool:start` | stderr | 始终 | 工具名 + 参数摘要 |
| `tool:output` | stderr | 按 display 级别 | `silent`=隐藏，`normal`=截断摘要，`verbose`=完整 |
| `tool:done` | stderr | 仅出错或 `--verbose` | 退出码 / 耗时 |
| `progress` | stderr | 第 2+ 轮 | 第 1 轮不显示（"正在想"是默认预期） |
| `error` | stderr | 始终 | 错误信息 |
| `done` | 内部 | 不显示 | 流结束信号 |

### 工具输出的三级显示

| 级别 | 何时用 | 示例 |
|------|--------|------|
| `silent` | `read_file` 成功、纯上下文获取 | 用户不需要看文件内容 |
| `normal` | `search_replace` diff、`run_skill` 日志 | 显示变更摘要 |
| `verbose` | `run_command` 结果、任何工具出错 | 完整输出 |

### Client 覆盖

```
cata chat             # 默认：normal 级别
cata chat --quiet     # 所有工具输出静默，只显示 AI 文本
cata chat --verbose   # 所有工具输出完整显示
cata chat --show-thinking  # 输出 reasoning/thinking 内容到 stderr
```

### Server 事件格式（NDJSON）

```jsonld
{"type":"token","content":"我来帮你"}
{"type":"thinking","content":"需要先读取文件..."}
{"type":"tool:start","id":"c1","name":"read_file","args":{"path":"main.go"},"display":"silent"}
{"type":"tool:output","id":"c1","name":"read_file","content":"package main\n...","display":"silent"}
{"type":"tool:done","id":"c1","name":"read_file","ok":true}
{"type":"tool:start","id":"c2","name":"run_command","args":{"argv":["go","test","./..."]},"display":"verbose"}
{"type":"tool:output","id":"c2","name":"run_command","content":"ok  ...\n","display":"verbose"}
{"type":"tool:done","id":"c2","name":"run_command","ok":true,"exit_code":0}
{"type":"progress","message":"model round 2"}
{"type":"error","message":"something went wrong"}
{"type":"done","success":true}
```

**`display` 字段**：Server 给 Client 的提示。Client 可根据 `--quiet` / `--verbose` 覆盖。`silent` 的工具输出在 `--quiet` 下完全不出现；`verbose` 的输出在 `--verbose` 下完整显示。

### 推理/思考内容

- DeepSeek thinking (`reasoning_content`) 默认不向用户展示
- `--show-thinking` 时作为 `thinking` 事件流式输出到 stderr
- API 层始终回传 `reasoning_content`（协议要求），只是 UI 层过滤

---

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
2. `~/.cata/global/constraints.md`（或仓库模板 `brain/constraints.md`）  
3. `brain/behavior.md`  
4. 具体 `internal/server`、`internal/evolve` 源码  
