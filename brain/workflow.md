# Workflow（演进与维护 SOP）

本文件与 `core.md` 的演进章节一致：**把一次「响应」变成可累积的「成长」**，并保证 **下一次 AI 读到的文档更准**。

## 1. 闭环（OODA）

**Observe** → **Orient / Decide** → **Act** → **Reflect**

| 阶段 | 做什么 | 产出 |
|------|--------|------|
| 分析 | 读 hot/short/long/archive 体量、evolution_log 近项 | `Snapshot`（`internal/evolve` Observe） |
| 决策 | 基于状态选单一主行动 | JSON：action、reason、learning、updates[] |
| 执行 | **文档补丁** 写入 `~/.cata/brain/` | 改 hot / short / long / core / workflow（白名单路径） |
| 反思 | 对照预期，**必须落盘** | 追加 `evolution_log.json`；learning 回答 workflow §4 三问 |

## 2. 决策启发（非硬编码 cron）

- archive 文件多且缺月度摘要 → **summarize**  
- 零散结论需固化身份/偏好 → **consolidate**（指向 hot 或 archive）  
- 同类失败 ≥2 → **reflect**，并考虑改 **hot.md** 禁忌或 **core.md** 一条规则  
- 复杂任务稳定成功 → 提炼 **Pattern** 入长期记忆  

（阈值由 Observe 快照 + LLM 判断；周期由 `evolution.cycle_interval` 驱动，**无单独 cron 二进制**。）

## 3. 任务类型（执行要点）

- **summarize**：压缩 archive，更新 summary 与索引。  
- **consolidate**：把会话/长期中的稳定结论写入 hot 或按日 archive。  
- **recall**：检索长期与索引。  
- **learn / optimize / reflect**：新技能、索引优化、规则反思；**reflect 必须连接「文档更新」**（见下）。  
- **idle**：明确记录「无需行动」及原因。

## 4. 学习反馈（必答三问）

1. 与预期的偏差？  
2. 是否有应写入 **长期** 的通用实践？  
3. 是否需要改 **core / workflow** 的权重或规则（一条即可，避免泛泛重写）？

## 5. 下一轮提升（**更新核心文档**）

当反思结论 **稳定、可复用** 时，**至少做一项**（由人或 Agent 落盘）：

1. **evolution_log.json**：追加一条（timestamp、action、reason、status、learning、doc_touched）。  
2. **hot.md**：更新目标、偏好、禁忌、技术栈一句——**下一轮注入最直接**。  
3. **core.md 或 workflow.md**：仅当规则级变更；**单行补丁式**增补，不删路径表、不破坏与 `paths.go` 的一致性。

这样下一次会话在 **boot → hot → 节选 core/workflow** 链路上会自动变强。

## 6. 演进日志示例

```json
{
  "timestamp": "2026-05-11T12:00:00Z",
  "action": "reflect",
  "reason": "Short-term memory grew past threshold.",
  "status": "completed",
  "learning": "Consolidated session notes into memory/long-term/patterns.md; hot unchanged.",
  "doc_touched": ["memory/short-term/current_session.md", "memory/long-term/patterns.md"]
}
```

## 7. 演进守则（三条）

1. **失败要进记忆**：风险写 hot 或长期，不只在聊天里消失。  
2. **剪枝**：过时规则与错误记忆优先 optimize / 归档。  
3. **不违背 core 根原则**；若要改根原则，必须在 evolution_log 写「核心变更」并最小 diff 改文档。
