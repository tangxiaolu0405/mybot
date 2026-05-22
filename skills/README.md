# Skills（开发模板）

运行时 **可执行 Skill** 只存在于脑子，不在产出区：

```text
~/.cata/brain/workspaces/<ws-id>/skills/<skill-id>/
├── SKILL.md
├── manifest.yaml
└── script.py   # 或其它 entry
```

全局兜底（可选）：`~/.cata/skills/<skill-id>/`

## 查找顺序（chat 注入 SKILL.md）

1. workspace 脑子 `skills/<id>/SKILL.md`
2. `~/.cata/skills/<id>/SKILL.md`
3. `~/.cursor/skills-cursor/<id>/SKILL.md`

## 执行

- 对话工具 **`run_skill`**：在产出区 cwd（`brain.base_dir` 或当前 cwd）执行脑子内脚本。
- **`capabilities.yaml`** 的 `skills:` 列表启用 id；由演进 `crystallize_skill` 后 **代码自动 append**，模型不得 patch 该文件。
- **`mcp:`**（如 browser）保留，用于未固化站点；skill 是已知任务的捷径。

## 固化

高 token 会话压缩后 + 重复 browser/任务关键词等门控 → `internal/evolve` 的 `crystallize_skill` 写入上述目录。

## 本目录

`skills/zhangtingban-lianban/` 等为**仓库模板**，供演进参考；`cata init` 不会自动拷贝到用户脑子。复制逻辑由 evolve 写入 workspace `skills/`。

参考脚本（开发用）：`scripts/fetch_zhangtingban_lianban.py`
