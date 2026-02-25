---
name: skills-creator
description: 定义CLAW项目自定义技能的标准格式和规范。AI agent在创建技能时必须遵循此规范。
version: 1.0.0
author: CLAW Team
tags: [agent, skill-creation, specification]
dependencies: []
---

# Skills Creator

**重要说明**: 本技能用于指导AI agent创建和管理CLAW项目中的自定义技能。所有技能描述应保持简洁，避免冗长的示例、模板和详细说明，以减少context消耗。

## 技能格式规范

### YAML Frontmatter（必需）

```yaml
---
name: skill-name
description: 技能的详细描述
version: 1.0.0
author: CLAW Team
tags: [tag1, tag2]
dependencies: []
---
```

**必需字段**：
- `name`: 技能唯一标识符，小写，连字符分隔
- `description`: 技能的详细描述
- `version`: 版本号（语义化版本）
- `author`: 作者名称
- `tags`: 技能标签列表
- `dependencies`: 依赖的其他技能名称列表

### Markdown 内容结构

必须包含以下部分：

```markdown
# 技能名称

## 功能描述
[描述技能的功能、解决的问题、主要特性和适用范围]

## 使用场景
- 场景1: 描述
- 场景2: 描述

## 指令说明
### 基础指令
- 指令1: 说明
- 指令2: 说明

### 高级指令
- 指令3: 说明

## 最佳实践
- 实践1
- 实践2

## 注意事项
- 注意1
- 注意2
```

## 命名约定

- 技能名称：小写字母，连字符分隔，描述性强且简洁
- 正确示例：`code-generator`, `document-parser`, `data-analyzer`
- 错误示例：`CodeGenerator`, `code_generator`, `code generator`

## 目录结构

```
skill-name/
├── SKILL.md              # 技能主文件（必需）
└── resources/            # 资源文件（可选）
    ├── config.json
    ├── templates/
    └── data/
```

## 技能调用约定

### 调用触发条件
1. 用户明确请求技能功能
2. 用户需求与某个技能的功能匹配
3. 当前技能依赖于其他技能

### 调用流程
```
用户请求 → 分析请求 → 识别技能 → 检查依赖 → 执行调用 → 返回结果
```

### 调用约定
1. 加载技能的 SKILL.md 文件
2. 解析 YAML frontmatter
3. 检查依赖关系
4. 根据用户请求匹配指令
5. 按照指令说明执行操作
6. 返回结果

## 依赖管理

在 YAML frontmatter 中声明依赖：

```yaml
dependencies:
  - skill-a
  - skill-b
```

依赖规则：
- 只声明必需的依赖
- 避免循环依赖
- 明确依赖版本

## 版本控制

遵循语义化版本：`MAJOR.MINOR.PATCH`
- `MAJOR`: 不兼容的 API 修改
- `MINOR`: 向下兼容的功能性新增
- `PATCH`: 向下兼容的问题修正

## 质量标准

### 格式检查清单
- YAML frontmatter 完整
- 必需字段齐全
- Markdown 格式正确
- 无语法错误

### 内容检查清单
- 描述清晰准确
- 指令明确
- 文档完整
