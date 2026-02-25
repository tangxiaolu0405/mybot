---
name: skill-validator
description: 验证和规范化技能文件的格式、结构和内容。发现问题时需用户确认后再修复。
version: 1.0.0
author: CLAW Team
tags: [validation, skill-management, compliance]
dependencies: []
---

# Skill Validator

**重要说明**: 本技能用于指导AI agent验证技能文件。发现问题时必须先向用户报告，获得用户确认后再进行修复。

## 功能描述

验证技能文件是否符合CLAW项目规范，包括YAML frontmatter、必需字段、Markdown格式、命名约定和目录结构。

## 使用场景

- 创建新技能后验证格式
- 修改现有技能后检查
- 批量验证所有技能
- 技能规范化

## 验证规则

### YAML Frontmatter（必需）

```yaml
---
name: skill-name
description: 技能的详细描述
version: 1.0.0
author: 作者名
tags: [tag1, tag2]
dependencies: []
---
```

**必需字段检查**：
- `name`: 小写，连字符分隔，唯一标识符
- `description`: 清晰描述技能功能和调用时机，200字符以内
- `version`: 遵循语义化版本规范（MAJOR.MINOR.PATCH）
- `author`: 作者名称
- `tags`: 标签列表，用于分类
- `dependencies`: 依赖的其他技能列表

**修复方法**：
- 缺失字段：在YAML frontmatter中添加缺失的字段
- 格式错误：修正字段值格式
- 版本号错误：修改为正确的语义化版本格式

### 命名约定

**规则**：
- 只使用小写字母
- 单词间用连字符（-）分隔
- 避免特殊字符和空格

**正确示例**: `code-generator`, `document-parser`
**错误示例**: `CodeGenerator`, `code_generator`, `code generator`

**修复方法**：
- 大写转小写：`CodeGenerator` → `code-generator`
- 下划线转连字符：`code_generator` → `code-generator`
- 空格转连字符：`code generator` → `code-generator`
- 移除特殊字符：`code@generator` → `code-generator`

### Markdown 内容结构

**必需部分**：
- 功能描述
- 使用场景
- 指令说明
- 最佳实践
- 注意事项

**修复方法**：
- 缺失部分：添加缺失的章节
- 格式错误：修正Markdown语法
- 层次混乱：调整标题层级

### 目录结构

**标准结构**：
```
skill-name/
├── SKILL.md              # 必需
└── resources/            # 可选
```

**修复方法**：
- 缺少SKILL.md：创建SKILL.md文件
- 目录名称错误：重命名为与技能名称一致

## 验证流程

### 单个技能验证

1. 读取技能的SKILL.md文件
2. 解析YAML frontmatter
3. 检查必需字段
4. 验证命名约定
5. 检查Markdown内容结构
6. 验证目录结构
7. 输出验证结果

### 批量验证

1. 扫描skills目录下所有子目录
2. 对每个技能执行验证流程
3. 收集所有验证结果
4. 生成汇总报告

## 问题处理流程

**重要**: 发现问题时必须遵循以下流程：

1. **报告问题**: 向用户详细说明发现的问题
2. **提供修复方案**: 明确说明如何修复每个问题
3. **等待确认**: 等待用户确认是否进行修复
4. **执行修复**: 获得确认后执行修复操作
5. **重新验证**: 修复后重新验证确保问题已解决

## 常见问题及修复

### 问题1: 缺少必需字段

**症状**: YAML frontmatter中缺少name、description、version等字段

**修复方法**:
```yaml
# 在YAML frontmatter中添加缺失字段
---
name: your-skill-name
description: 技能描述
version: 1.0.0
author: 作者名
tags: [tag1, tag2]
dependencies: []
---
```

### 问题2: 命名不符合约定

**症状**: 技能名称包含大写字母、下划线或空格

**修复方法**:
- 将所有字母转为小写
- 将下划线和空格替换为连字符
- 移除特殊字符

### 问题3: 版本号格式错误

**症状**: 版本号不符合语义化版本规范

**修复方法**:
- 使用格式: MAJOR.MINOR.PATCH
- 示例: 1.0.0, 1.2.3, 2.0.1

### 问题4: Markdown内容结构不完整

**症状**: 缺少功能描述、使用场景、指令说明等必需章节

**修复方法**:
- 添加缺失的章节
- 确保包含所有必需部分
- 保持正确的Markdown格式

### 问题5: 目录结构错误

**症状**: 缺少SKILL.md文件或目录名称不匹配

**修复方法**:
- 确保技能目录包含SKILL.md文件
- 确保目录名称与技能名称一致

## 最佳实践

- 创建新技能后立即验证
- 修改技能后重新验证
- 定期批量验证所有技能
- 发现问题及时修复
- 保持命名约定一致性

## 注意事项

- 验证过程不会修改文件，只进行检查和报告
- 修复操作必须获得用户确认
- 必需字段缺失必须修复
- 命名约定不符合规范可能导致技能无法被正确识别
- 版本号应遵循语义化版本规范
