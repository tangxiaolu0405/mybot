---
name: skill-list-reader
description: 从 skills 目录读取并列出所有可用技能。当需列举技能、按名称/标签查找技能或确认可调用技能时使用。
version: 1.0.0
author: CLAW Team
tags: [skill-management, discovery, listing]
dependencies: []
---

# skill-list-reader

## 触发

- 需列举当前可用技能
- 按名称或关键词查找技能
- 按标签筛选技能
- 检查技能依赖或完整性

## 执行

1. 读取 `skills/skills-index.json`；若不存在或需重建，执行 `python scripts/generate_skills_index.py` 后重读。
2. 从 `skills` 数组与 `tags_index` 取数据；按需按 tag 或 name/description 过滤。
3. 输出：基础列表（name, version, description, author, tags）或详细列表（含 dependencies）。

## 索引结构

```json
{
  "version": "1.0.0",
  "generated_at": "ISO8601",
  "skills": [{ "name", "path", "description", "version", "author", "tags", "dependencies" }],
  "tags_index": { "tag": ["skill-name", ...] }
}
```

## 输出约定

- 基础：`Available Skills: 1. <name> (v<version>) Description: ... Author: ... Tags: [...]`
- 详细：每技能单独块，含 Dependencies。
- 仅解析 skills 目录下 SKILL.md；frontmatter 解析失败则跳过该条并记录。
