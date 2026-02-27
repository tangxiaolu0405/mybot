---
name: data-create
description: 创建新的记录、条目或记忆。用于持久化存储新产生的信息。
version: 1.0.0
author: CLAW Team
tags: [basic, crud, write]
dependencies: []
---

# data-create

## 触发条件
- 任务完成产生新结论
- 用户要求记住某项信息
- 产生新的标准化流程

## 执行步骤
- [ ] 1. 验证待创建数据的完整性（检查必需字段）
- [ ] 2. 分配全局唯一 ID (UUID)
- [ ] 3. 写入目标存储介质
- [ ] 4. 确认写入成功并返回条目摘要

## 输入/输出约定
- **输入**: `payload` (JSON 对象)
- **输出**: `status: success/fail`, `created_id`