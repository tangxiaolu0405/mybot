# Skills（提示词型）

在脑子分区的 `modes/<mode>/capabilities.yaml` 中启用：

```yaml
skills:
  - create-rule
mcp:
  - browser
```

查找顺序（`SKILL.md`）：

1. `~/.cata/skills/<name>/SKILL.md`
2. `<项目根>/skills/<name>/SKILL.md`
3. `~/.cursor/skills-cursor/<name>/SKILL.md`（与 Cursor 共用）

执行型 skill（脚本）尚未接入。
