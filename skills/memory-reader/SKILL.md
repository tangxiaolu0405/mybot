---
name: memory-reader
description: 读取与管理 CLAW 记忆（长期/短期/工作）。支持按关键词检索、按类型过滤、摘要。任务前后加载或更新记忆时使用。
version: 1.0.0
author: CLAW Team
tags: [memory, retrieval, knowledge-management]
dependencies: []
---

# memory-reader

## 触发

- 任务前需加载相关记忆
- 需按关键词或类型检索记忆
- 需读/写短期或长期记忆
- 任务后需更新记忆

## 路径

- 长期记忆：`brain/memory/long-term/`
- 短期记忆：`brain/memory/short-term/current_session.md`
- 核心思维：`brain/core.md`

## 指令

- **读长期记忆**：读 long-term 目录下文件；可按类型/关键词过滤。
- **读短期记忆**：读 current_session.md。
- **按关键词检索**：在指定范围（长期/短期/全部）内搜索关键词，返回匹配片段与路径。
- **按类型过滤**：类型包括项目知识、代码模式、用户偏好、技术栈、历史决策、经验教训（长期）；当前会话、任务队列、上下文、临时变量、中间结果、会话历史（短期）。
- **写短期记忆**：更新 current_session.md，保持既有结构。
- **写长期记忆**：在 long-term 下新增或更新文件，遵循命名与分类约定。

## 输出约定

- 检索结果：条数 + 每条 [类型] 标题、关键信息、可选关联。
- 详细格式：类型、路径、更新时间、摘要、关键信息、关联记忆。

## 约束

- 文件不存在时按空处理并可选创建
- 格式异常时兼容解析并记录
- 不泄露敏感信息；更新遵循 core.md 中的流转规则
