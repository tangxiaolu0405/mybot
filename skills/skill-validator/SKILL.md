---
name: skill-validator
description: 验证并规范化技能文件的 frontmatter、结构与内容。发现问题时先向用户报告并获确认后再修复。
version: 1.0.0
author: CLAW Team
tags: [validation, skill-management, compliance]
dependencies: []
---

# skill-validator

## 触发

- 新建或修改技能后需校验
- 需批量校验 skills 目录下所有技能

## 约束

**禁止未获用户确认即修改文件**：先报告问题与修复方案，待确认后执行修复并复验。

## 验证项

### YAML frontmatter（必需）

- `name`：小写、连字符、唯一
- `description`：非空，建议 ≤200 字
- `version`：语义化版本 MAJOR.MINOR.PATCH
- `author`：非空
- `tags`：数组
- `dependencies`：技能名数组

缺失或格式错误：在 frontmatter 中补全/修正。

### 命名

- 仅小写、连字符；无空格与特殊字符。
- 纠错：大写→小写，下划线/空格→连字符，去掉特殊字符。

### 目录与文件

- 技能目录内必须有 `SKILL.md`；目录名与 `name` 一致。

## 流程

1. 单技能：读 SKILL.md → 解析 frontmatter → 校验上述项 → 输出结果。
2. 批量：遍历 skills 下子目录，对每个执行单技能流程，汇总报告。
3. 若需修复：向用户列出问题与修复方案 → 获确认 → 执行修改 → 再次验证。
