---
name: skill-list-reader
description: Reads and lists all available skills from the skills directory. Invoke when user asks for available skills, needs to find a specific skill, or AI needs to know what skills exist.
version: 1.0.0
author: CLAW Team
tags: [skill-management, discovery, listing]
dependencies: []
---

# Skill List Reader

## 功能描述

读取skills目录下的所有技能，列出技能的名称、描述、版本、作者、标签和依赖关系，帮助AI和用户了解可用的技能资源。

## 使用场景

- 用户询问"有哪些技能可用"
- 用户需要查找特定功能的技能
- AI需要了解当前有哪些技能可以调用
- 用户想要浏览技能列表
- 需要检查技能的依赖关系

## 指令说明

### 基础指令

- **读取技能列表**: 扫描skills目录，读取所有SKILL.md文件
- **显示技能摘要**: 展示每个技能的基本信息（名称、描述、版本）
- **显示详细信息**: 展示技能的完整信息（包括标签、依赖等）
- **按标签过滤**: 根据标签筛选技能
- **按名称搜索**: 根据技能名称或描述关键词搜索

### 高级指令

- **检查依赖关系**: 分析技能之间的依赖关系
- **验证技能完整性**: 检查技能文件是否存在且格式正确
- **生成技能报告**: 生成技能列表的汇总报告

## 执行流程

1. 读取skills/skills-index.json索引文件
2. 如果索引文件不存在或过期，运行python scripts/generate_skills_index.py生成新索引
3. 从索引中提取所有技能信息
4. 格式化输出技能列表
5. 支持按标签过滤和按名称搜索
6. 使用tags_index进行快速标签查找

## 索引文件结构

skills-index.json包含以下结构：
```json
{
  "version": "1.0.0",
  "generated_at": "ISO时间戳",
  "skills": [技能列表],
  "tags_index": {标签到技能名称的映射}
}
```

## 输出格式

### 基础列表格式

```
Available Skills:
1. skill-name (v1.0.0)
   Description: 技能描述
   Author: 作者名
   Tags: [tag1, tag2]
```

### 详细列表格式

```
Skill: skill-name
Version: 1.0.0
Author: 作者名
Description: 技能描述
Tags: [tag1, tag2]
Dependencies: [dep1, dep2]
```

## 最佳实践

- 每次会话开始时可以读取技能列表以了解可用技能
- 当用户询问技能相关问题时，优先读取技能列表
- 使用标签过滤可以快速找到相关技能
- 定期检查技能的依赖关系确保技能可用

## 注意事项

- 只读取skills目录下的SKILL.md文件
- 如果SKILL.md文件格式错误，跳过该技能并记录错误
- 技能名称必须与目录名称一致
- 确保技能列表的输出清晰易读
- 处理读取错误时提供友好的错误信息
