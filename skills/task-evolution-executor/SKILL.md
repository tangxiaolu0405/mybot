---
name: task-evolution-executor
description: 单次任务全生命周期：准备→执行→记忆迭代→进化评估。需记忆支持且要迭代与效果评估时使用。依赖 memory-reader。
version: 1.0.0
author: CLAW Team
tags: [evolution, task-execution, memory-iteration, self-improvement]
dependencies: [memory-reader]
---

# task-evolution-executor

## 触发

- 任务需记忆支持且完成后需迭代与进化评估
- 需记录执行中的决策与反馈并反哺记忆与模式

## 依赖

使用 memory-reader 做检索与写入；记忆迭代规则见 memory-iteration-manager。

## 四阶段流程

### 1. 准备

- 用 memory-reader 检索长期/短期记忆中与任务相关项
- 从 skills-index 识别所需技能及依赖
- 将任务目标与上下文写入工作记忆（短期或会话内结构）

### 2. 执行

- 按任务要求与技能指令执行；记录关键决策与理由、错误与解决、用户反馈
- 执行后验证结果并收集满意度信号

### 3. 记忆迭代

- 提炼关键信息与模式；更新短期记忆（过程、结果、反馈）
- 评估是否满足提升为长期记忆的条件（见 memory-iteration-manager）；满足则写入长期并建立关联

### 4. 进化评估

- 计算成功指标：解决率、满意度、效率、技能利用率、记忆命中率
- 计算失败指标：错误率、重试次数、用户修正、资源浪费
- 与历史基线对比；识别模式优化机会（如连续 N 次低于阈值或连续负面反馈）；将评估结果与优化建议写入记忆，供后续模式优化使用

## 约束

- 记忆操作遵循 core.md 与 memory-iteration-manager 的流转规则
- 进化结论基于客观数据；模式调整需谨慎并保留回滚能力
