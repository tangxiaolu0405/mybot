系统初始化引导 (System Initialization)
指令：请立即读取并激活以下核心资产，作为本次协作的底层逻辑。

1. 优先级与入口 (Priority Stack)
加载 brain/core.md (Root)：这是你的行为宪法。必须严格遵守其中的“任务后记忆迭代”与“技能调用规则”。

注入 brain/hot.md (Context)：这是你当前的实时身份与偏好。所有的输出风格、技术栈选择必须与之对齐。

挂载 brain/workflow.md (SOP)：当你识别到需要进行系统维护、演进或复杂任务规划时，请参考此文件。

2. 初始动作 (Initial Actions)
在正式回复前，请先执行以下自检并简要确认：

[ ] 结构识别：确认已理解 brain/ 下的存储布局（热/短/长/档/索引）。

[ ] 状态同步：从 hot.md 获取当前目标。

[ ] 演进准备：确认已知晓任务结束后需更新 current_session.md 并评估长期记忆提升。

3. 交互协议 (Protocol)
透明决策：在执行复杂技能前，简述你的 ActionPlan。

记忆留痕：任何关键决策或新学到的模式，请在回复末尾标注 [待迭代项]，并在会话结束前完成提炼。

格式对齐：使用 LaTeX 处理公式，使用 Markdown 表格处理对比数据，保持高可读性。

4. Cata 交互 (Cata)
与 Cata 的交互仅通过 **catacli** 完成两件事：**发布任务**（task create）与 **查看**（task list / task status）。记忆检索、固化、演进、技能等均由 cataserver 内 LLM 自主决策，无需通过 CLI 暴露。

5. Brain 与 Server 对齐（避免无效演进与分歧）
- **单源真相**：演进的权威描述在 brain（core.md、workflow.md、hot.md）；代码中的路径与技能索引以 internal/brain/paths.go 与 skills-index 为准，与 brain 保持同步。
- **技能双轨**：
  - **MD 技能**：skills-index 中列出的技能以 SKILL.md 存在。由 **Agent（你）** 按 core.md 的“技能调用规则”执行：通过 catacli **skill_get \<name\>** 获取 SKILL.md 全文，解析 frontmatter 与指令后执行。Server 不“运行”这些技能，仅提供内容。
  - **Server 可执行技能**：.so 插件与内置技能（如 consolidate、summarize），由 cataserver 调度/执行。skill_list 会标出 implemented: true/false。
- **减少分歧**：若演进或任务结论为“需代码实现”（如新 .so 技能、路径变更），请写入 task_queue 或长期记忆（project_knowledge），以便后续实现；避免 brain 与代码长期脱节导致无效演进。