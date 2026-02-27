---
name: data-delete
description: 移除指定的记录。用于清理冗余信息或执行用户撤回请求。
version: 1.0.0
author: CLAW Team
tags: [basic, crud, delete]
dependencies: [data-query]
---

# data-delete

## 触发条件
- 用户明确要求删除
- 缓存过期或达到存储上限

## 执行步骤
- [ ] 1. 确认待删除记录的 ID
- [ ] 2. 执行“软删除”（标记 deleted 字段）或“硬删除”
- [ ] 3. 更新 `skills-index.json`（如果是技能删除）
- [ ] 4. 反馈清理结果

## 警示
- 删除操作不可逆，需在 Prompt 中引导 Agent 进行二次确认。