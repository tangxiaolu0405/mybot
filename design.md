## Cata 系统设计

### 架构概览

```
┌─ 用户 ─┐
    │  cata chat [--dir <产出区>]
    ▼
┌─ Client (internal/client) ─────────────────────────────┐
│  REPL 循环 → Unix Socket → NDJSON 事件流 → 终端渲染     │
└────────────────────────────────────────────────────────┘
    │  Unix Socket (~/.cata/cata.sock)
    ▼
┌─ Server (internal/server) ─────────────────────────────┐
│  连接管理 → 脑子解析 → 聊天循环 → 工具执行               │
│  (managed 模式: 最后一个 chat 断开后自动退出)            │
└────────────────────────────────────────────────────────┘
    │
    ├── LLM (internal/llm) ─── OpenAI 兼容 API
    │
    ├── 脑子 (~/.cata/brain/workspaces/<id>/)
    │   ├── memory/short/current.md    ← 每轮写入
    │   ├── memory/long/               ← evolve 归档
    │   ├── memory/index.json          ← 记忆索引
    │   ├── modes/<mode>/persona.md    ← evolve 维护
    │   └── skills/                    ← evolve 固化
    │
    └── Evolve (internal/evolve) ─── 后台异步
        观察 → LLM 决策 → 文档补丁 → 索引同步
```

### 核心模型

**两层循环**：

```
外层 while true              对话级，等待触发（用户输入 / cron / 外部事件）
    └── 内层 while true      任务级，LLM + 工具链执行
            └── break: LLM 返回最终文本（无 tool_calls）或超限
```

**两种触发**：
| 触发源 | 说明 |
|--------|------|
| 用户消息 | `cata chat` 交互式对话 |
| 定时演进 | `evolve.cycle_interval`（默认 600s）后台自主运行 |

---

### 产出区设计（参考 Claude Code）

核心问题：用户不一定在项目目录下运行 `cata`，需要能指定"在哪个目录干活"。

**方案**：

```
cata chat                         # 产出区 = 当前目录 (默认，向后兼容)
cata chat --dir ~/project         # 产出区 = ~/project
cata chat --dir ~/a --dir ~/b     # 多产出区，第一个是主产出区
```

**配置文件** (`~/.cata/config.json`)：

```json
{
  "workspace": {
    "default_dir": "~/myproject"
  }
}
```

**产出区 vs 脑子的关系**：

```
产出区 (output dirs)          脑子 (~/.cata)
─────────────────────         ─────────────────────
文件工具操作范围              persona / 记忆 / 技能
run_command 执行目录          演进日志 / 注册表
项目 .git 检测起点            不存用户代码
```

- **产出区** = 用户的项目文件所在位置（代码、文档、构建产物）
- **脑子** = Agent 的记忆和 persona（永远在 `~/.cata/`）
- **focus_path** = 从产出区向上查找 `.git` 或 `.cata/workspace.yaml`，决定绑定哪个脑子格子
- 借鉴 Claude Code：`--add-dir` → cata 的多个 `--dir`；Claude 的 launch dir → cata 的第一个 `--dir` 或 cwd

**规则**：
1. `--dir` 指定的目录成为文件工具和 `run_command` 的操作根目录
2. 文件工具只能访问产出区内的路径（`safePathUnder` 检查）
3. 脑子格子选择基于主产出区解析（`focus_path` 逻辑不变）
4. 同一个产出区目录只能开一个 chat session（output lock 不变）
5. 不同产出区可以并行开多个 chat

---

### 交互层设计：对话交付

核心原则：**stdout 是 AI 的回答正文，stderr 是元信息**。

#### 事件类型与显示策略

```
事件流 (Server → Client NDJSON)
═══════════════════════════════════════════════════════

token           → stdout    AI 文本流（唯一向 stdout 输出的东西）
thinking        → stderr    推理/思考过程（默认隐藏，--show-thinking 开启）
tool:start      → stderr    工具名 + 参数摘要
tool:output     → stderr    工具输出（根据 display 级别决定）
tool:done       → stderr    退出码/状态（仅出错或 verbose 时显示）
progress        → stderr    第 2+ 轮提示（第 1 轮不显示）
error           → stderr    错误信息（始终显示）
done            → 内部      流结束信号（不显示）
```

#### 工具输出的三级显示

| 级别 | 含义 | 适用工具 |
|------|------|----------|
| `silent` | 不显示输出内容 | `read_file` 成功时（AI 在阅读，用户不需要看原文） |
| `normal` | 显示摘要/截断输出 | `search_replace` diff、`run_skill` 日志 |
| `verbose` | 显示完整输出 | `run_command` 结果、任何工具出错时 |

Server 在事件中携带 `display` 提示，Client 可以根据 `--quiet` / `--verbose` 覆盖。

```
# 正常模式（默认）
› add tests for auth
⟳ round 2                                    ← stderr
📄 reading internal/auth/auth.go              ← stderr, tool:start
✏ editing internal/auth/auth_test.go          ← stderr, tool:start
  + func TestAuthenticate(t *testing.T) {     ← stderr, tool:output (normal)
  +   ...
⚙ go test ./internal/auth/...                 ← stderr, tool:start
ok  internal/auth  0.234s                     ← stderr, tool:output (verbose)
我已经添加了测试...                            ← stdout, token stream

# 安静模式 (--quiet)
› add tests for auth
我已经添加了测试...                            ← stdout only，工具静默

# 详细模式 (--verbose)
› add tests for auth
⟳ round 1                                    ← stderr
📄 reading internal/auth/auth.go              ← stderr
  (145 lines)                                 ← stderr, 输出摘要
✏ editing internal/auth/auth_test.go          ← stderr
  + func TestAuthenticate...                  ← stderr, 完整 diff
⚙ go test ./internal/auth/...                 ← stderr
ok  internal/auth  0.234s                     ← stderr, 完整命令输出
我已经添加了测试...                            ← stdout
```

#### 推理/思考（DeepSeek thinking）

- **默认**：不向用户展示 `reasoning_content`
- **`--show-thinking`**：将思考内容作为 `thinking` 事件输出到 stderr
- 思考内容始终包含在 LLM 请求的 `reasoning_content` 回传中（API 要求），只是不展示

#### 文件操作确认

- `search_replace`：默认不确认（可逆操作）
- `append_file`：默认不确认
- `run_command`：黑名单命令或 `require_confirm` 时弹出确认

---

### 存储层结构

```
~/.cata/
├── registry/workspaces.json      # 工作区注册表
├── global/
│   ├── constraints.md            # 全局约束
│   ├── behavior.md               # 全局行为 SOP
│   └── boot-assembler.md         # Boot leader 指令
├── brain/workspaces/<ws_id>/
│   ├── meta.json
│   ├── persona.local.md          # 聚焦上下文
│   ├── evolution_log.json        # 演进日志
│   ├── memory/
│   │   ├── index.json            # 记忆索引（常驻 context）
│   │   ├── short/current.md      # 短期记忆（每轮写入）
│   │   ├── long/                 # 长期记忆（evolve 归档）
│   │   └── archive/              # 冷记忆
│   ├── modes/<mode>/
│   │   ├── persona.md            # 模式 persona（evolve 维护）
│   │   ├── behavior.md
│   │   ├── constraints.md
│   │   └── capabilities.yaml     # MCP + skills 声明
│   └── skills/<id>/
│       ├── SKILL.md
│       ├── manifest.yaml
│       └── script.py
├── skills/                       # 全局共享技能
├── locks/                        # 产出区锁文件
└── cata.sock                     # Unix socket
```

### 记忆分层（与 design.md 对齐）

| 层 | 位置 | 写入方 | 作用 |
|----|------|--------|------|
| Socket 会话历史 | server 内存 | 每轮对话 | 当前 session 上下文，chat_reset 清空 |
| short/current.md | 每格脑子 | 每轮 chat 成功后追加 | 对话原文，evolve 的输入 |
| memory/index.json | 每格脑子 | evolve 同步 | 摘要索引，常驻 context（< 2800 bytes） |
| persona.md（hot） | modes/<mode>/ | evolve 提炼 | 偏好 + 流程，全量注入 context |
| long/ + archive/ | 每格脑子 | evolve 归档 | 低频事实，按需召回 |

### Context 组装（每次 LLM 调用重新组装）

```
固定层（每次必有）:
    boot-assembler.md          引导指令
    路径块                      脑子 vs 产出区 + 运行时环境
    global/constraints.md      全局约束
    global/behavior.md         全局行为
    mode persona.md            当前模式 persona
    persona.local.md           聚焦上下文
    memory/index.json          记忆索引
    技能 SKILL.md 块           当前模式启用的技能

动态层（按需注入）:
    对话历史（in-memory）      最近轮次 user/assistant/tool
    压缩摘要                   历史超出阈值时触发 session compress

硬限制:
    persona 块      < 6500 bytes/文件，总计 < 20000 bytes
    记忆索引         < 2800 bytes
    skills 块       < 8000 bytes/skill，总计 < 16000 bytes
    历史压缩后       < context_window × 40%
```

### 自主演进

```
触发条件:
    short-term > shortTermTriggerBytes（有足够新内容）
    或 short-term 自上次演进后有变化 + >= shortTermActivityBytes
    或 archive 文件数 >= archiveSummarizeMinFiles

周期:
    默认 600s，由 evolve.cycle_interval 控制

动作:
    observe → LLM 决策 (idle|update|consolidate|crystallize) → 文档补丁 → 索引同步

无手动演进命令。
```

### 保护机制

```
内层循环:
    上下文超 context_window × 85% → 触发 session compress → 历史截断到 40%
    单轮最大 tool 轮次限制（隐式，由 token 预算控制）

工具执行:
    run_command 输出上限 256KB
    命令黑名单 + 用户确认
    路径遍历防护 (safePathUnder)

并发:
    同一产出区只能一个 chat session（output lock）
    同一机只能一个 server（socket 文件锁）
    记忆读写由 evolve 引擎串行化（单 goroutine）

记忆膨胀:
    short/current.md 上限 96KB → 触发 trim
    persona.md 超 6500 bytes → 触发 consolidate
    index.json 超 2800 bytes → 触发 summary 压缩
```

### 刻意排除

- 无手动演进命令（纯后台自主运行）
- 无任务队列、无 scheduler
- 无内置 git 操作
- 无 Web UI
- 无多机分布式
- CLI `catacli` 已废弃
