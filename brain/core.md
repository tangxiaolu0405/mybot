# Core（行为宪法）

本页是 **单源真相**：角色、路径、记忆与自主演进。细节流程见 `brain/workflow.md`。

## 脑子 vs 产出区（必读）

完整规划见 **`brain/directory-plan.md`**。

| | **脑子** | **产出区** |
|--|----------|------------|
| 是什么 | Agent 的记忆、人格、演进（`~/.cata`） | 当前工作目录 **cwd**：代码、文件、`run_command` 结果 |
| 路径 | `CATA_HOME`（默认 `~/.cata/`） | `cata` 启动时 cwd / `brain.base_dir` |
| 不是什么 | 项目仓库、交付物目录 | 脑子正文、persona、short-term |

- **`~/.cata/brain/<id>/`**：一格脑子（绑定某个 **focus_path**，用于选哪套记忆；不等于产出放这里）。  
- **项目内 `.cata/`**：可选门牌（`workspace.yaml` / `link`），**不是**脑子搬项目里。  
- **仓库 `brain/`**：模板种子，不是运行时脑子。

| 概念 | 路径 | 环境变量 |
|------|------|----------|
| 脑子根 | `~/.cata/` | `CATA_HOME` |
| 脑子分区 | `brain/workspaces/<id>/` | 由 focus_path（git/yaml/cwd）解析 |
| 产出区 | 当前 cwd | 请求 `cwd`；`brain.base_dir` |

## ~/.cata 存储布局（定稿）

```text
~/.cata/
├── config.json
├── cata.sock
├── registry/
│   └── workspaces.json          # 所有 workspace（含临时 cwd）
├── global/
│   ├── constraints.md           # 全局约束（种子来自 core.md）
│   ├── behavior.md              # 全局 SOP（种子来自 workflow.md）
│   └── boot-assembler.md        # 注入顺序（种子来自 boot-leader.md）
└── brain/
    └── workspaces/
        └── <workspace-id>/
            ├── meta.json
            ├── persona.local.md
            ├── evolution_log.json   # 本 workspace 演进审计（多份，非全局）
            ├── modes/
            │   ├── _default/
            │   │   ├── persona.md      # ≈ 原 hot，由 evolve 维护
            │   │   ├── behavior.md
            │   │   ├── constraints.md
            │   │   └── capabilities.yaml
            │   └── <mode-id>/          # 项目内迭代长出，无内置固定 mode
            └── memory/
                ├── short/current.md    # server 每轮对话追加
                ├── long/*.md
                ├── archive/
                └── index.json          # 演进维护；对话注入摘要，按需 read_file 展开
```

**项目目录（可选，不在 home 内）**

| 路径 | 作用 |
|------|------|
| `.git/` | 有则 workspace 根 = git root |
| `.cata/workspace.yaml` | 有则 workspace 根 = 该目录；可写 `name`、`active_mode` |
| `.cata/workspace.link` | 可选；`id: ws_xxx` |

**无 git、无 yaml**：cwd 作为 workspace（`kind: ephemeral`），**自动**在 `~/.cata` 创建存储。

## 工作区解析顺序

1. 从 `cata chat` 传入的 **cwd** 向上找 **git root**  
2. 否则找 **`.cata/workspace.yaml`**  
3. 否则 **cwd** 本身（临时项目，仍自动建库）

## 记忆分层

| 层 | 位置 | 写入 | 读取 |
|----|------|------|------|
| 会话历史 | 内存（socket） | server | 同连接 LLM |
| short-term | `memory/short/current.md` | server 每轮 | evolve |
| persona | `modes/<mode>/persona.md` | **仅 evolve** | 每轮对话注入 |
| persona.local | `persona.local.md` | evolve / 人 | 每轮注入 |
| long / archive | `memory/long`、`memory/archive` | evolve | 按需（index） |
| global | `~/.cata/global/` | 人 / 种子 | 约束+行为节选 |
| **skills** | `skills/<id>/`（脑子内） | **仅 evolve** `crystallize_skill` | `capabilities.yaml` + `run_skill` / SKILL 注入 |

**Skill**：脚本与 `manifest.yaml` 只在脑子；**产出**写在 cwd。查找：workspace `skills/` → `~/.cata/skills/`。`mcp: [browser]` 保留给未固化站点。

对话注入：`boot-assembler` → global 约束/行为 → mode persona → persona.local（见 `internal/brain/terminal_context.go`）。

## 自主演进

- 周期：`evolution.cycle_interval`（默认 600s）  
- **每个 workspace 单独** Observe / 门控 / LLM / 补丁 / **`evolution_log.json`**  
- **会话压缩**：估算本连接 **socket history**（含将注入的 boot/global/persona）达到 `llm.context_window × evolution.context_compress_ratio`（默认 **85%**）时，触发 consolidate（short-term → persona），并把 history 裁到约 **40%** 窗口留空给回复  
- 路径白名单：`persona.local.md`、`modes/*/persona|behavior|constraints`、`memory/**`、`skills/**`（SKILL.md、manifest、脚本）  
- **`crystallize_skill`**：高 token 压缩后 / 重复 browser 等 → 写脑子 `skills/<id>/`，**代码** append `capabilities.yaml` 的 `skills:`（禁止模型改 `mcp:`）  
- 周期演进门控：short-term 活动等（见 `internal/evolve`）

## 终端对话上下文

| 存储 | 位置 | 作用 |
|------|------|------|
| **socket history** | 服务端内存（每连接） | 多轮 user/assistant/tool，发给 LLM；`/clear` 清空 |
| **short-term** | `memory/short/current.md` | 磁盘流水，供演进提炼；与 history 不同步截断 |

- **Tool 调用**：无固定轮次上限；exec 需用户确认。  
- **不再**按「56 条消息」机械截断；仅在接近上下文上限时做演进压缩 + 按 token 预算裁 history。

## Mode

- 不内置小说/游戏等；新建 workspace 仅有 **`_default`**。  
- 新 mode 由 **evolve 分叉** 或用户后续扩展；均在 `modes/<mode-id>/` 下。

## 进化与质量

- 成功实践 → persona / long；规则级变更极少改 global。  
- 多 CLI / 多 Agent 见 `todo.md`。
