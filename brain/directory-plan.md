# 目录规划：脑子（.cata）vs 产出区（工作目录）

> **核心区分（用户定义）**  
> - **`.cata` / `~/.cata`**：Agent 的**脑子**——记忆、人格、约束、演进、本机服务。  
> - **当前工作目录（cwd）**：**产出存放处**——代码、文件、命令执行结果等交付物。  
>  
> 二者分离：脑子不替代项目文件夹；项目文件夹也不承载脑子正文。

---

## 1. 两个世界

```text
┌──────────────────────────────────────┐
│  脑子 · CATA_HOME（默认 ~/.cata/）      │
│  · 想、记、演进                        │
│  · config / socket / global / state   │
│  · 不当作「项目仓库」用                  │
└──────────────────────────────────────┘
              ▲
              │ 绑定（用哪个脑子处理当前活）
              │ 键：git 根 / 声明文件 / cwd
              │
┌──────────────────────────────────────┐
│  产出区 · 当前工作目录（cwd）            │
│  · 写代码、改文稿、go build、git commit  │
│  · run_command 的 cwd 默认在这里         │
│  · Agent 工具产生的结果落在这里          │
└──────────────────────────────────────┘
```

| | **脑子（~/.cata）** | **产出区（cwd）** |
|--|---------------------|-------------------|
| 隐喻 | 记忆、习惯、人格 | 手的工作成果 |
| 典型内容 | persona、short-term 流水、evolution_log | 源码、配置、构建产物、用户文档 |
| 是否进用户 git | 否（本机私有） | 是（用户的项目） |
| 谁写 | server 追加 short-term；evolve 改 persona | 用户 + `run_command` + 外部工具 |
| 对话注入读哪 | `~/.cata/global` + `state/<id>/modes/...` | 不默认把产出区全文注入 |

---

## 2. ~/.cata：脑子在哪（结构）

```text
~/.cata/                          # CATA_HOME = 脑子根
├── config.json                   # 本机服务（LLM、exec、evolution）
├── cata.sock
│
├── global/                       # 全机级思维底层（约束、行为、boot 顺序）
│   ├── constraints.md
│   ├── behavior.md
│   └── boot-assembler.md
│
├── registry/
│   └── workspaces.json           # 索引：「关注路径」→ 哪一格脑子
│       # 关注路径 ≠ 产出区定义权，见 §3
│
└── brain/                        # 脑子分区（实现名；可改称 state/）
    └── <brain-id>/               # 一格脑子，对应一个「关注对象」
        ├── meta.json             # 绑定的 focus_path、kind、active_mode
        ├── persona.local.md      # 对当前关注对象的说明
        ├── evolution_log.json
        ├── modes/<mode-id>/      # 人格/模式（项目内迭代，非内置）
        └── memory/
            ├── short/current.md  # 对话流水（原料）
            ├── long/
            └── archive/
```

**脑子里的 `memory/`**：是 Agent **内部记忆**，不是把你的小说/代码拷贝进脑子目录。  
提炼结果进 `persona`；交付物仍在 **cwd**。

---

## 3. 产出区：当前工作目录

- **cata chat** 每次请求带 **cwd**（客户端 `os.Getwd()`）。
- **`run_command`**、`brain.base_dir`（配置）对齐 **产出区**：命令在这里执行，生成物写在这里。
- 用户用编辑器打开的文件、git 操作的对象，都在产出区。

**脑子如何知道「用哪一格」：**  
解析时用一个 **关注路径（focus_path）** 选 `~/.cata/brain/<id>/`，规则仍可沿用：

1. 从 cwd 向上 → **git 根**（有则 `focus_path = git root`）
2. 否则 → 含 `.cata/workspace.yaml` 的目录
3. 否则 → **`focus_path = cwd`**（临时目录也有一格脑子，但产出仍只在 cwd）

> **重要**：`focus_path` 只决定 **加载哪一格脑子**；**不**把产出区迁移到 `~/.cata`。  
> 在 monorepo 子目录干活时：cwd 可以是 `repo/services/api`，脑子仍可绑 `repo` 根（若按 git 根）。

---

## 4. 项目里的 `.cata/` 是什么？

**不是脑子本身**，只是脑子上的**标签**（可选）：

```text
<产出区某目录>/          # 通常是 git 根，也可以是 cwd
└── .cata/
    ├── workspace.yaml   # 可选：name、active_mode（给人看/给团队）
    └── workspace.link   # 可选：brain_id: ws_xxx → 指向 ~/.cata/brain/<id>
```

- 可提交 git 的只有**声明**，不是 persona/short-term 正文。  
- **脑子正文始终在 `~/.cata`**。

---

## 5. 仓库 `mybot/brain/` 是什么？

| 位置 | 角色 |
|------|------|
| `mybot/brain/` | **脑子模板的出厂复印件**（init 时拷到 `~/.cata/global/` 等） |
| `~/.cata/` | **运行时脑子** |
| **cwd** | **产出区** |

三者不要互称「工作区 brain」。

---

## 6. 运行时数据流（对齐实现）

```text
用户在 产出区/cwd 下执行 cata
    │
    ├─► 解析 focus_path → 选中 ~/.cata/brain/<id>/  （脑子）
    │
    ├─► LLM 注入：global + modes/persona + persona.local （来自脑子）
    │
    ├─► run_command 在 产出区（base_dir / cwd）执行          （产出）
    │
    ├─► 每轮成功：AppendChatTurn → 脑子/memory/short/     （记一笔）
    │
    └─► 上下文 ≥85% window：evolve 压缩脑子 → persona      （脑子整理）
         socket history 裁短（会话缓存，不在磁盘产出区）
```

---

## 7. 命名建议（对话/文档）

| 避免 | 改用 |
|------|------|
| 「工作区 = ~/.cata」 | **脑子 = ~/.cata**；**产出区 = cwd**（或配置的 base_dir） |
| 「workspace 数据在 home」 | **脑子分区 `<brain-id>` 在 home** |
| 「brain 目录在用户项目里」 | 项目里只有 **`.cata` 声明**；脑子在 home |

代码里 `workspace.RootPath` 宜理解为 **`focus_path`（绑定脑子用）**，不是「产出根」。

---

## 8. 配置与环境变量

| 项 | 作用 |
|----|------|
| `CATA_HOME` | 脑子根（~/.cata） |
| `brain.base_dir` | **产出区默认根**（exec、找 go.mod 等）；应跟 cwd 或 git 根对齐，不是脑子路径 |
| `CATA_BRAIN_DIR` | 遗留名；若存在，应等于或指向 **脑子根**，勿与 cwd 混淆 |

---

## 9. 实施备注（代码现状）

- 已实现：`ResolveWorkspace(cwd)` → 选脑子分区；`AppendChatTurn` 写脑子；exec 用 `BrainBaseDir`（产出）。  
- 文档/日志建议显式打印：  
  - `brain: ~/.cata/brain/ws_xxx`  
  - `output cwd: /path/to/project`  
- 可选重命名：`brain/workspaces/` → `brain/instances/` 或 `brain/cells/`，避免「workspace = 用户文件夹」歧义。

---

## 10. 一句话

- **~/.cata = 脑子**（想、记、演进）。  
- **cwd = 产出区**（做、写、生成）。  
- 项目内 `.cata/` = 给脑子贴的**门牌号**，不是脑子搬进门牌里。
