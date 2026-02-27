---
name: data-update
description: 修改已存在的记录。用于更正错误、更新状态或迭代记忆。
version: 1.0.0
author: CLAW Team
tags: [basic, crud, update]
dependencies: [data-query]
---

# data-update

## 触发条件
- 发现既有数据过时或错误
- 任务状态发生变更（如：从“进行中”变为“已完成”）

## 执行步骤
- [ ] 1. 调用 `data-query` 确认目标记录存在
- [ ] 2. 对比新旧数据，仅更新变化的部分（Patch 模式）
- [ ] 3. 执行更新操作并记录修改日志
- [ ] 4. 返回更新后的完整对象预览

## 注意事项
- 严禁在未确认 ID 的情况下进行大范围覆盖。