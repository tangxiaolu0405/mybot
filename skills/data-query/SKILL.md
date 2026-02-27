---
name: data-query
description: 根据关键词或条件检索结构化数据或历史记录。当需要获取特定信息以支持决策时调用。
version: 1.0.0
author: CLAW Team
tags: [basic, crud, read]
dependencies: []
---

# data-query

## 触发条件
- 接收到需要背景资料的任务
- 用户询问特定历史数据或状态

## 执行步骤
- [ ] 1. 提取用户查询中的核心实体（Keywords）
- [ ] 2. 调用底层存储接口（如向量数据库或 JSON 文件）
- [ ] 3. 过滤并排序结果（按相关性或时间）
- [ ] 4. 格式化输出给 LLM 进行下一步推理

## 输入/输出约定
- **输入**: `query_string` (字符串), `limit` (整数)
- **输出**: 匹配到的对象数组或“未找到”消息