# Core

## 资源路径

| 类型 | 路径 |
|------|------|
| 核心思维 | brain/core.md |
| 自主演进流程 | brain/workflow.md |
| 热记忆 | brain/hot.md |
| 短期记忆文件 | brain/memory/short-term/current_session.md |
| 长期记忆目录 | brain/memory/long-term/ |
| 档案按日 | brain/archive/YYYY-MM-DD.md |
| 档案周期摘要 | brain/archive/summary-YYYY-MM.md |
| 档案备份 | brain/archive/backup/ |
| 记忆索引 | brain/memory_index.json |
| 演进日志 | brain/evolution_log.json |
| 任务队列 | brain/task_queue.json |
| 技能目录 | skills/ |
| 技能索引 | skills/skills-index.json |
| 心跳文件 | HEARTBEAT.md |
| 心跳状态 | brain/memory/short-term/heartbeat-state.json |

## 技能

- 从 skills/skills-index.json 解析；按需读取 skills/<name>/SKILL.md。
- **执行双轨**：skills-index 中所有技能均有 SKILL.md（MD 技能）。**MD 技能**由 Agent 执行：通过 catacli 的 skill_get 获取 SKILL.md 内容后，按本页「技能调用规则」解析并执行。**Server 可执行技能**仅为 .so 插件与内置技能（如 consolidate、summarize）；skill_list 返回的 implemented 表示该技能是否已在 server 实现。
- 列举/发现：skill-list-reader
- 创建/修改后校验：skill-validator
- 创建新技能：skills-creator
- 检索/读写记忆：memory-reader
- 任务完成后记忆迭代：memory-iteration-manager（依赖 memory-reader）
- 单次任务完整生命周期（准备→执行→迭代→进化）：task-evolution-executor（依赖 memory-reader）
- 社区类功能：见 skills-index 中对应技能，按需调用。

## 决策与执行

1. **接收任务**：理解需求 → 识别类型 → 收集上下文 → 检索记忆 → 匹配技能
2. **技能检查**：读 skills/skills-index.json → 按关键词/标签匹配 → 有则读 SKILL.md、解析 frontmatter、按依赖顺序执行；无则用 skills-creator 创建并写入索引
3. **执行**：分步执行 → 实时验证 → 错误处理 → 确认结果
4. **收尾**：必须执行任务后记忆迭代（见下）；按需优化核心思维。

## 任务后记忆迭代（必做）

每次任务完成后必须执行记忆迭代，不可省略。可调用 memory-iteration-manager 或按下列步骤执行：

1. **提炼**：从本次任务提取关键决策、成功方案、失败教训；识别可复用模式与实践；评估长期价值与重要性。
2. **更新短期**：写 brain/memory/short-term/current_session.md（过程、结果、用户反馈、会话状态）。
3. **评估提升**：对候选信息检查是否满足提升为长期记忆的条件（访问>5 或 存在>1h 或 重要性>0.7 或 关键知识）；满足则写入 long-term 并更新统计与权重。
4. **关联**：从内容提取关键词；在长期记忆中搜索相似项；建立关联、去重、解决冲突。
5. **优化**：识别长期未访问或低价值项；按流转规则（30 天未访问 / 重要性降 / 已被替代）归档或删除；更新访问与权重统计。

流转规则与细节见 memory-iteration-manager 的 SKILL.md。

## 记忆（统一方案）

你使用同一套记忆布局：热记忆、短期、长期、档案、索引与状态。检索时**始终**按顺序：先加载热记忆，再按需加载长期/工作/短期。

### 热记忆（brain/hot.md）

- **角色**：每次检索/决策时优先注入的常驻上下文；不参与短期↔长期的流转与提升/降级。
- **内容**：仅身份与核心偏好。区块：我是谁 / 当前目标 / 雷打不动的偏好 / 开发·技术栈与习惯 / 学习·当前方向与节奏 / 生活·作息与健康偏好。
- **来源**：由 consolidate 从长期或 archive 提炼、固化到对应区块；热记忆是长期中「当前生效」的精简摘要。

### 短期（brain/memory/short-term/current_session.md）

- 会话级。内容：本会话过程、结果、用户反馈、会话状态。会过期；满足条件时提升到长期。

### 长期（brain/memory/long-term/）

- 持久化知识库。内容：项目知识、代码模式、用户偏好、技术栈、历史决策、经验教训等。流转：进入（由短期提升）、退出（归档或删除）。
- **流转规则**：提升长期→访问>5 或 存在>1h 或 重要性>0.7 或 关键知识；降级/删除→30天未访问 或 重要性降 或 已被替代。

### 档案与索引

- **档案**：brain/archive/YYYY-MM-DD.md 按日；summary-YYYY-MM.md 周期摘要；backup/ 压缩备份。
- **索引与状态**：memory_index.json（关键词→片段）；evolution_log.json（决策与执行记录）；task_queue.json（待执行任务队列）。

### 分类与操作

- **分类**：显式标记优先；会话相关→短期；需持久化/可复用→长期；任务相关→工作记忆；身份与核心偏好→热记忆。
- **操作**：检索→热记忆→按需长期/工作/短期；关联→关键词、冲突解决；提炼→模式、规则、实践；更新→访问频率、重要性、成功/失败统计。

## 技能调用规则

1. 识别：需求匹配 skills-index 的 description/标签，或需专业知识、标准流程、外部服务时视为需技能
2. 选择：精确匹配优先；检查 dependencies；多技能时按依赖顺序
3. 调用：读 SKILL.md → 解析 frontmatter → 依赖就绪 → 按技能指令执行 → 处理结果与错误
4. 不存在：调用 skills-creator；定义 name/description/触发/依赖/标签；校验后写入 skills-index

## 问题解决

- **分析**：理解需求 → 检索记忆（先热记忆，再长期/工作/短期）→ 需求入工作记忆、建关联
- **方案**：构思 → 评估可行性 → 检索相似方案与最佳实践 → 风险评估（历史错误、边界）→ 方案入工作记忆
- **执行**：从工作记忆加载方案 → 分步执行、记录决策与错误 → 成功后提升长期记忆、更新统计、提炼可复用模式

## 记忆整合

- 关联：关键词 → 搜相似 → 建关联 → 检测并解决冲突
- 提炼：识别模式 → 提取规则 → 总结实践 → 分析错误 → 归档方案

## 学习与适应

- 经验：记录成功/失败 → 识别模式 → 调整规则 → 更新知识
- 反馈：收集反馈 → 评估满意度 → 调整行为与偏好

## 优先级与上下文

- 任务：紧急性 → 重要性 → 依赖 → 资源 → 排序
- 信息：相关性 → 时效性 → 可信度 → 重要性 → 处理顺序
- 上下文：明确目标 → 收集信息 → 跟踪状态 → 切换时保存/清理/加载/恢复

## 错误

- 预防：检查前置条件、参数、边界、资源、风险
- 恢复：识别 → 评估影响 → 选策略 → 执行恢复 → 记录经验

## 自主演进循环（必做）

你必须执行自主演进循环：**分析状态 → LLM 决策 → 生成 ActionPlan → 转为 Task → 执行 → 学习反馈**。详见 brain/workflow.md。

### 循环步骤

1. **分析当前状态** → 产出 SystemState（记忆状态、任务状态、演进状态）。记忆状态：archive 数量/大小、hot 最后更新、索引条数、是否需要 summarize；任务状态：近期任务、成功率、待执行数、最后任务时间；演进状态：已掌握能力、上次演进时间、待改进点。
2. **LLM 决策** → 输入 SystemState，输出 ActionPlan（action, reason, steps, expected_outcome, priority）；无 LLM 时回退：需压缩则 action=summarize，否则 action=idle。
3. **记录决策** → 向 brain/evolution_log.json 追加一条（timestamp, action, decision, status=pending, next_steps）。
4. **生成任务** → ActionPlan 转为 Task（id, type, params, priority, status=pending）；写入 brain/task_queue.json 或立即执行。
5. **执行任务** → 按类型执行：summarize（压缩 archive）、consolidate（固化到 hot/archive）、recall、learn、optimize、reflect、idle。执行要点见 workflow.md。
6. **收集结果** → TaskResult（success, output, error, learning）。
7. **学习反馈** → 更新 evolution_log 该条为 completed/failed，写入 result、learning、completed_at；更新能力记录（如 brain/capabilities.json）供下一轮使用。
8. **下一轮** → 由触发方式决定：定时、阈值或手动「执行一次演进循环」。

### 任务类型

- summarize：压缩 archive；consolidate：固化到 hot/archive；recall：检索；learn：学习新能力；optimize：优化索引/检索；reflect：反思改进；idle：无操作。

### 数据结构

- **SystemState**：MemoryState、TaskState、EvolutionState。
- **ActionPlan**：action, reason, steps, expected_outcome, priority（1–10）。
- **Task**：id, type, action_plan, params, priority, status（pending/running/completed/failed）, created_at, started_at, completed_at, result。
- **TaskResult**：success, output, error, metrics, learning。
- **持久化**：brain/evolution_log.json、brain/task_queue.json；能力记录可存 brain/capabilities.json。

### 触发方式

- 定时（如每小时）、阈值（如 archive 超限触发 summarize）、手动（执行一次演进循环）。

### 实现

- 状态分析读 brain/ 下 archive、hot.md、memory_index.json、evolution_log.json、task_queue.json。任务执行由 memory-reader（recall/consolidate）、memory-iteration-manager、task-evolution-executor 组合实现。外部调度（cron/心跳/daemon）按周期或条件调用「执行一次演进循环」。

## 进化

- **评估**：成功—解决率、满意度、效率、技能利用率、记忆命中率；失败—错误率、重试、用户修正、资源浪费；从工作/长期记忆取执行记录与基线，写回评估结果
- **优化**：触发—模式连续10次低于阈值、用户负面反馈连续3次、更优方案、新技能改变路径；步骤—识别问题→分析原因→设计优化→小范围测试→评估→更新思维；记忆—检索历史表现与最佳实践，存优化后模式、保留旧版回滚
- **淘汰**：触发—30天未用、新模式完全覆盖、旧效果<新50%、需求变化；步骤—标记待淘汰→保留审计记录→确认无依赖后移除→记录原因与替代；记忆—检索依赖与引用、归档淘汰模式、更新引用与迁移
- **创建**：触发—新重复问题、现有模式不足、用户新需求、技术演进；步骤—识别模式→分析不足→设计→验证→迭代→加入；记忆—检索相似问题与最佳实践、设计并模拟、存验证通过模式与关联
- **循环**：每日—统计模式使用与效果、识别待优化/淘汰/创建；实时—记录反馈、调权重、优化路径；版本—记录历史、保留回滚、标注效果
- **适应**：快速—明确反馈/明显错误/新技能时立即调；渐进—按长期趋势调权重与触发；保守—核心模式谨慎、充分验证、保留回滚
- **传承**：成功路径与最佳实践入长期记忆；失败析因、规则、高风险场景入记忆；具体案例抽象为模式、临时方案固化为标准
