本文件定义了系统的**闭环进化逻辑**。所有任务执行必须锚定此流程，确保从单纯的“响应”转变为“成长”。

## 1. 流程全景图

系统遵循 **感知(Observe) → 决策(Orient/Decide) → 执行(Act) → 反思(Reflect)** 的循环。

---

## 2. 核心阶段详解

### 第一阶段：分析状态 (Analyze System State)

在任何演进周期开始前，必须读取以下数据点生成 `SystemState`：

* **记忆状态**：扫描 `brain/archive/` 大小，检查 `memory_index.json` 的冗余度，确认 `hot.md` 的时效性。
* **任务状态**：从 `evolution_log.json` 提取最近 10 次任务的成功/失败分布。
* **能力状态**：识别用户最近频繁触达但当前 `skills/` 缺失的领域。

### 第二阶段：LLM 决策与 ActionPlan

基于状态，系统必须输出一个结构化的 `ActionPlan`：

> **示例决策逻辑**：如果 `archive` 文件超过 5 个且未汇总 $\rightarrow$ 执行 `summarize`；如果某项错误连续出现 3 次 $\rightarrow$ 执行 `reflect` 并修改 `hot.md`。

### 第三阶段：任务转换与执行 (Task Execution)

将 `ActionPlan` 分解为原子化任务写入 `task_queue.json`。

* **Recall**：从长期记忆中捞取相关模式。
* **Learn**：调用 `skills-creator` 构建新工具。
* **Consolidate**：将零散的知识固化为规则。

### 第四阶段：学习反馈 (Learning & Feedback)

**这是演进的关键点**。任务完成后必须回答：

1. 本次执行与预期结果的偏差是什么？
2. 是否有新的“最佳实践”需要写入长期记忆？
3. 是否需要调整 `brain/core.md` 中的决策权重？

---

## 3. 自发性触发阈值 (Triggers)

| 触发源 | 阈值条件 | 对应行动 |
| --- | --- | --- |
| **时间触发** | 每 24 小时 | 执行 `summarize` 与 `backup` |
| **记忆触发** | 短期记忆条目 > 20 条 | 启动 `memory-iteration-manager` 进行提升/归档 |
| **错误触发** | 同类 Error 出现次数 $\ge$ 2 | 强制执行 `reflect`（反思）并更新核心思维 |
| **成功触发** | 复杂任务高分完成 | 提取 `Pattern`（模式）并存入 `long-term` |

---

## 4. 演进日志格式 (Evolution Log)

所有演进记录必须严格遵守 `brain/evolution_log.json` 的格式，以便后续自检：

```json
{
  "timestamp": "2026-02-25T19:40:00Z",
  "action": "optimize",
  "decision": "Detected redundant coding patterns in Python tasks.",
  "status": "completed",
  "learning": "Consolidated 'FastAPI-auth-template' into long-term memory to reduce redundant generation by 40%.",
  "next_steps": ["Update skills-index to include the new template."]
}

```

---

## 5. 进化守则 (The Evolution Creed)

1. **绝不遗忘教训**：失败的尝试必须转化为“风险警告”存入热记忆。
2. **动态剪枝**：过时的技能和错误的记忆比没有记忆更危险，必须定期执行 `optimize` 淘汰。
3. **一致性优先**：新进化的规则不得与 `brain/core.md` 的根原则冲突，除非发起“核心重构”。

---
