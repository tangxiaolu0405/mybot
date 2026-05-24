# Workflow（演进与维护 SOP）

## OODA 闭环

| 阶段 | 做什么 | 产出 |
|------|--------|------|
| Observe | 读 short/long/archive 体量、evolution_log 近项 | `Snapshot` |
| Decide | 基于状态选单一主行动 | JSON: action, reason, learning, updates[] |
| Act | 文档补丁写入白名单路径 | 改 persona / long；summarize 时移入 archive |
| Reflect | 对照预期，落盘 | evolution_log + learning 三问 |

## 决策启发

- long-term 文件过多（≥25）→ **summarize**：合并/压缩 long-term 旧条目，移入 archive（冷存储）
- 零散结论需固化身份/偏好 → **consolidate** → persona 或 long
- 同类失败 ≥2 → **reflect**，考虑改 persona 禁忌
- 复杂任务稳定成功 ≥2 → **crystallize** → `skills/<id>/`
- persona 中出现互斥的风格/偏好（如不同项目的习惯明显不同）→ **consolidate** 时考虑将当前 persona 分叉到 `modes/<new-id>/`，更新 `meta.json` active_mode
- 无事可做 → **idle**（明确记录原因）

## 任务类型

| 动作 | 效果 |
|------|------|
| `summarize` | 压缩 long-term，移入 archive（冷存储，不再参与 evolve 和 context） |
| `consolidate` | short-term → persona / long；少数情况下可附带 mode 分叉 |
| `crystallize` | 稳定流程 → `skills/<id>/`（SKILL.md + manifest + script） |
| `update` | 单条规则/偏好写入 persona |
| `reflect` | 评估近期决策质量，产出改进 |
| `idle` | 无需行动 |

## 学习反馈（必答三问）

1. 与预期的偏差？
2. 是否有应写入长期记忆的通用实践？
3. 是否需要改 persona 或约束（一条即可，避免泛泛重写）？

## 演进守则

1. **失败要进记忆**：风险写 persona 或 long，不只在聊天里消失
2. **剪枝**：过时规则与错误记忆优先 optimize / 归档
3. **最小变更**：改 global 约束须单行补丁 + evolution_log 标记"核心变更"
4. **不破坏能力**：crystallize 只追加 `skills:`，禁止模型改 `mcp:` 列表

## evolution_log 格式

```json
{
  "timestamp": "2026-05-22T12:00:00Z",
  "action": "consolidate",
  "reason": "short-term 积累超过阈值",
  "status": "completed",
  "learning": "用户偏好显式错误处理；写入 persona",
  "doc_touched": ["modes/_default/persona.md", "memory/long/patterns.md"]
}
```
