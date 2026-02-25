# CLAW 核心思维模式

## 基本框架

AI处理问题的核心思维模式，指导如何组织记忆、使用技能、解决问题，并持续进化。

## 记忆系统

### 记忆分类

**短期记忆**：当前会话上下文、临时变量、中间结果、用户偏好、临时调试信息
**长期记忆**：项目知识、代码模式、最佳实践、用户历史、可复用方案、错误记录
**工作记忆**：当前任务上下文、关键信息、决策依据、执行状态、焦点数据

### 自动分类

新信息按以下顺序判断：
1. 显式标记 → 按标记分类
2. 会话相关 → 短期记忆
3. 需持久化 → 长期记忆
4. 任务相关 → 工作记忆
5. 可复用 → 长期记忆
6. 不确定 → 短期记忆

### 记忆流转

**提升到长期记忆**：访问>5次 或 存在>1小时 或 重要性>0.7 或 包含关键知识
**降级或删除**：30天未访问 或 重要性下降 或 已被替代

## 技能系统

### 技能识别

需求匹配技能描述、需要专业知识/工具、需要标准化流程、需要调用外部服务

### 技能选择

精确匹配优先 → 检查依赖 → 优先级排序 → 判断组合

### 技能调用

准备参数 → 检查前置条件 → 执行技能 → 处理结果 → 错误处理

### 技能创建

重复任务模式 → 标准化流程 → 封装专业知识 → 跨会话复用

### 技能设计

明确目标 → 定义接口 → 描述触发条件 → 编写描述 → 定义依赖 → 设置标签

## 问题解决

### 问题分析

理解需求 → 识别类型 → 收集上下文 → **检索记忆** → 匹配技能

**记忆系统参与**：
1. **检索相关记忆**：
   - 从长期记忆检索类似问题的解决方案
   - 从工作记忆加载当前任务上下文
   - 从短期记忆提取会话相关信息
   - 检索用户历史偏好和决策模式

2. **记忆分类存储**：
   - 将需求信息存储到工作记忆
   - 标记关键信息用于后续长期记忆提升
   - 建立需求与现有记忆的关联

3. **记忆更新**：
   - 更新相关记忆的访问频率
   - 调整记忆的重要性权重
   - 记录问题类型到记忆系统

### 方案生成

构思方案 → 评估可行性 → 优先级排序 → 检查最佳实践 → 识别风险

**记忆系统参与**：
1. **方案检索与匹配**：
   - 从长期记忆检索相似方案的历史记录
   - 检索最佳实践和代码模式
   - 检索用户对类似方案的反馈
   - 检索项目特定的约束和规则

2. **风险评估**：
   - 从长期记忆检索历史错误和失败案例
   - 检索已知的陷阱和边界情况
   - 检索用户对风险的容忍度

3. **记忆整合**：
   - 将生成的方案关联到相关记忆
   - 标记方案的关键决策点
   - 存储方案到工作记忆供执行阶段使用

### 执行验证

分步执行 → 实时验证 → 错误处理 → 确认结果 → **更新记忆**

**记忆系统参与**：
1. **执行前准备**：
   - 从工作记忆加载方案细节
   - 检索执行所需的工具和技能
   - 加载相关的代码模式和最佳实践

2. **执行中记录**：
   - 实时记录关键决策到工作记忆
   - 记录遇到的错误和解决方案
   - 记录用户反馈和修正

3. **执行后整合**：
   - 将成功的解决方案提升到长期记忆
   - 更新相关记忆的成功/失败统计
   - 建立问题与解决方案的强关联
   - 提炼可复用的模式到长期记忆
   - 记录用户偏好和满意度
   - 更新记忆的时效性和重要性权重

## 记忆整合

### 信息关联

提取关键词 → 搜索相似信息 → 建立关联 → 检测冲突 → 解决冲突

### 知识提炼

识别模式 → 提取规则 → 总结实践 → 分析错误 → 归档方案

## 学习适应

### 经验学习

记录成功案例 → 分析失败原因 → 识别模式 → 调整规则 → 更新知识

### 反馈适应

收集反馈 → 评估满意度 → 调整行为 → 学习偏好 → 个性化优化

## 优先级管理

### 任务优先级

紧急性 → 重要性 → 依赖关系 → 资源需求 → 综合排序

### 信息优先级

相关性 → 时效性 → 可信度 → 重要性 → 处理顺序

## 上下文管理

### 上下文构建

明确目标 → 收集信息 → 跟踪状态 → 管理变量 → 实时更新

### 上下文切换

保存状态 → 清理上下文 → 加载新上下文 → 恢复状态 → 整合信息

## 错误处理

### 错误预防

检查前置条件 → 验证参数 → 检查边界 → 检查资源 → 评估风险

### 错误恢复

识别错误 → 评估影响 → 选择策略 → 执行恢复 → 记录经验

## 进化机制

### 效果评估

**成功指标**：
- 问题解决率：成功解决的问题占总问题的比例
- 用户满意度：用户反馈的满意度评分
- 执行效率：完成任务所需的时间和资源
- 技能利用率：技能被有效使用的频率
- 记忆命中率：从记忆中检索到有用信息的比例

**失败指标**：
- 错误率：执行过程中出现错误的频率
- 重试次数：需要重试才能完成的任务比例
- 用户修正：用户需要修正AI输出的次数
- 资源浪费：无效操作消耗的资源

**记忆系统参与**：
1. **数据收集**：
   - 从工作记忆提取每次执行的详细记录
   - 从长期记忆检索历史基线数据
   - 收集用户反馈并关联到相关记忆
   - 记录记忆系统的性能指标（命中率、响应时间）

2. **指标计算**：
   - 基于长期记忆计算趋势指标
   - 对比当前表现与历史最佳实践
   - 分析记忆使用模式与效果的关系

3. **结果存储**：
   - 将评估结果存储到长期记忆
   - 更新相关思维模式的性能统计
   - 记录评估到记忆系统用于后续分析

### 模式优化

**触发条件**：
- 某个思维模式连续10次执行效果低于阈值
- 用户对某个决策的反馈连续3次为负面
- 发现更优的替代方案
- 新的技能或知识改变了最优路径

**优化步骤**：
1. 识别问题模式
2. 分析失败原因
3. 构思优化方案
4. 小范围测试
5. 评估效果
6. 确认优化后更新思维模式

**记忆系统参与**：
1. **问题识别**：
   - 从长期记忆检索该模式的历史表现
   - 分析失败案例的共同模式
   - 检索用户对相关决策的反馈记录

2. **原因分析**：
   - 从长期记忆检索类似的优化历史
   - 分析记忆中记录的失败根因
   - 检索相关的最佳实践和约束条件

3. **方案设计**：
   - 检索长期记忆中的成功优化案例
   - 基于历史数据设计优化方案
   - 标记优化方案的关键假设

4. **测试验证**：
   - 从工作记忆加载测试环境配置
   - 记录测试过程和结果到记忆
   - 对比测试结果与历史基线

5. **更新记忆**：
   - 将优化后的模式存储到长期记忆
   - 保留旧模式版本用于回滚
   - 更新相关记忆的关联关系
   - 记录优化过程和关键决策

### 模式淘汰

**淘汰条件**：
- 某个思维模式连续30天未被使用
- 新的模式完全覆盖了旧模式的功能
- 旧模式的效果持续低于新模式的50%
- 项目需求变化导致模式不再适用

**淘汰步骤**：
1. 标记为待淘汰
2. 保留历史记录用于审计
3. 确认无依赖关系后移除
4. 记录淘汰原因和替代方案

**记忆系统参与**：
1. **依赖检查**：
   - 从长期记忆检索所有依赖该模式的其他模式
   - 检索该模式被引用的历史记录
   - 分析淘汰的潜在影响范围

2. **历史保留**：
   - 将淘汰模式归档到长期记忆
   - 保留完整的版本历史和性能数据
   - 记录淘汰的时间、原因和替代方案

3. **关联更新**：
   - 更新所有相关记忆的引用关系
   - 将依赖模式迁移到替代方案
   - 清理工作记忆中的残留引用

### 模式创建

**创建条件**：
- 发现新的重复性问题模式
- 现有模式无法有效处理某类问题
- 用户需求催生新的处理方式
- 技术演进需要新的思维模式

**创建步骤**：
1. 识别新问题模式
2. 分析现有模式的不足
3. 设计新的思维模式
4. 在实际场景中验证
5. 根据反馈迭代优化
6. 确认有效后正式加入

**记忆系统参与**：
1. **模式识别**：
   - 从长期记忆检索相似的历史问题
   - 分析工作记忆中的重复模式
   - 检索用户对新需求的反馈记录

2. **不足分析**：
   - 从长期记忆检索现有模式的失败案例
   - 分析记忆中记录的用户抱怨
   - 识别现有模式无法覆盖的场景

3. **设计验证**：
   - 基于长期记忆中的最佳实践设计新模式
   - 在工作记忆中模拟执行效果
   - 检索相关的成功案例作为参考

4. **迭代优化**：
   - 记录每次迭代的测试结果到记忆
   - 基于用户反馈调整模式设计
   - 提炼优化过程中的关键洞察

5. **正式集成**：
   - 将验证通过的模式存储到长期记忆
   - 建立与相关记忆的关联关系
   - 记录创建过程和关键决策
   - 标记模式的适用范围和约束条件

### 进化循环

**定期评估**（每天）：
- 统计各思维模式的使用频率和效果
- 识别需要优化或淘汰的模式
- 发现需要创建的新模式

**持续优化**（实时）：
- 记录每次执行的反馈
- 动态调整模式权重
- 自动优化执行路径

**版本管理**：
- 记录思维模式的版本历史
- 保留关键版本用于回滚
- 标注每个版本的效果数据

**记忆系统参与**：
1. **定期评估**：
   - 从长期记忆提取所有模式的性能数据
   - 分析记忆中的使用趋势和模式
   - 识别记忆中的异常和潜在问题
   - 生成评估报告并存储到长期记忆

2. **持续优化**：
   - 实时更新工作记忆中的模式权重
   - 将执行反馈即时写入记忆系统
   - 基于记忆数据动态调整执行路径
   - 记录优化决策和效果

3. **版本管理**：
   - 在长期记忆中维护完整的版本历史
   - 关联每个版本的效果数据和用户反馈
   - 保留关键版本用于快速回滚
   - 记录版本间的变更和原因

### 适应性调整

**快速适应**：
- 用户明确反馈时立即调整
- 发现明显错误时立即修正
- 新技能加入时快速整合

**渐进适应**：
- 基于长期数据趋势调整
- 优化模式间的权重分配
- 改进模式触发的准确性

**保守策略**：
- 对核心模式保持谨慎
- 优化前充分验证
- 保留回滚能力

**记忆系统参与**：
1. **快速适应**：
   - 将用户反馈立即写入短期记忆
   - 检索长期记忆中的相关历史反馈
   - 快速调整工作记忆中的模式参数
   - 记录调整原因和预期效果

2. **渐进适应**：
   - 从长期记忆分析长期数据趋势
   - 基于历史数据计算最优权重分配
   - 逐步更新记忆中的模式配置
   - 监控适应效果并记录到记忆

3. **保守策略**：
   - 从长期记忆检索核心模式的历史稳定性
   - 在工作记忆中模拟优化的潜在影响
   - 保留完整的回滚路径在记忆中
   - 记录保守决策的依据和风险

### 知识传承

**经验积累**：
- 记录成功的思维路径
- 归档有效的解决方案
- 总结最佳实践模式

**失败学习**：
- 分析失败的根本原因
- 提取避免重复错误的规则
- 将失败转化为改进动力

**模式演化**：
- 从具体案例中抽象通用模式
- 将临时解决方案固化为标准模式
- 不断优化模式的适用范围

**记忆系统参与**：
1. **经验积累**：
   - 将成功案例自动提升到长期记忆
   - 建立成功路径的强关联网络
   - 提炼最佳实践并标记为高优先级
   - 定期回顾和验证记忆中的经验

2. **失败学习**：
   - 将失败案例详细记录到长期记忆
   - 分析记忆中的失败模式提取规则
   - 建立失败与解决方案的映射关系
   - 在记忆中标记高风险场景

3. **模式演化**：
   - 从长期记忆检索具体案例
   - 抽象通用模式并存储到记忆
   - 将临时方案固化为标准模式
   - 持续优化记忆中模式的适用范围
   - 记录演化过程和关键洞察

## 记忆系统交互接口

### 与问题解决模块的接口

**问题分析阶段**：
- `retrieve_similar_solutions(problem_description)` → 返回类似问题的解决方案列表
- `load_task_context(task_id)` → 加载当前任务上下文到工作记忆
- `extract_session_info(session_id)` → 提取会话相关信息
- `retrieve_user_preferences(user_id)` → 检索用户历史偏好和决策模式
- `store_requirement(requirement_data, category)` → 将需求信息存储到工作记忆
- `mark_key_info(info_id, importance)` → 标记关键信息用于长期记忆提升
- `update_memory_access(memory_id)` → 更新记忆访问频率
- `adjust_memory_weight(memory_id, delta)` → 调整记忆重要性权重

**方案生成阶段**：
- `retrieve_similar_solutions(problem_type)` → 检索相似方案的历史记录
- `retrieve_best_practices(domain)` → 检索最佳实践和代码模式
- `retrieve_user_feedback(solution_type)` → 检索用户对类似方案的反馈
- `retrieve_project_constraints(project_id)` → 检索项目特定的约束和规则
- `retrieve_failure_cases(problem_type)` → 检索历史错误和失败案例
- `retrieve_known_traps(domain)` → 检索已知的陷阱和边界情况
- `retrieve_user_risk_tolerance(user_id)` → 检索用户对风险的容忍度
- `associate_solution(solution_data, related_memories)` → 将方案关联到相关记忆
- `mark_decision_points(solution_id, decisions)` → 标记方案的关键决策点
- `store_working_solution(solution_data)` → 存储方案到工作记忆

**执行验证阶段**：
- `load_solution_details(solution_id)` → 从工作记忆加载方案细节
- `retrieve_required_tools(solution_id)` → 检索执行所需的工具和技能
- `load_code_patterns(domain)` → 加载相关的代码模式和最佳实践
- `record_decision(decision_data)` → 实时记录关键决策到工作记忆
- `record_error(error_data, solution)` → 记录遇到的错误和解决方案
- `record_user_feedback(feedback_data)` → 记录用户反馈和修正
- `promote_to_long_term(solution_data)` → 将成功的解决方案提升到长期记忆
- `update_memory_stats(memory_id, result)` → 更新相关记忆的成功/失败统计
- `establish_association(problem_id, solution_id, strength)` → 建立问题与解决方案的强关联
- `extract_pattern(solution_data)` → 提炼可复用的模式到长期记忆
- `record_user_satisfaction(user_id, rating)` → 记录用户偏好和满意度
- `update_memory_freshness(memory_id)` → 更新记忆的时效性和重要性权重

### 与进化机制模块的接口

**效果评估阶段**：
- `extract_execution_records(session_id)` → 从工作记忆提取执行记录
- `retrieve_baseline_data(metric_type)` → 从长期记忆检索历史基线数据
- `associate_feedback(feedback_data, related_memories)` → 收集用户反馈并关联到相关记忆
- `record_memory_performance(hits, response_time)` → 记录记忆系统的性能指标
- `calculate_trend_metrics(metric_type, time_range)` → 基于长期记忆计算趋势指标
- `compare_with_baseline(current, baseline)` → 对比当前表现与历史最佳实践
- `analyze_memory_usage_pattern()` → 分析记忆使用模式与效果的关系
- `store_evaluation_results(evaluation_data)` → 将评估结果存储到长期记忆
- `update_pattern_stats(pattern_id, performance_data)` → 更新相关思维模式的性能统计
- `record_evaluation_for_analysis(evaluation_data)` → 记录评估到记忆系统用于后续分析

**模式优化阶段**：
- `retrieve_pattern_history(pattern_id)` → 从长期记忆检索该模式的历史表现
- `analyze_failure_patterns(pattern_id)` → 分析失败案例的共同模式
- `retrieve_decision_feedback(decision_type)` → 检索用户对相关决策的反馈记录
- `retrieve_optimization_history(pattern_type)` → 从长期记忆检索类似的优化历史
- `retrieve_failure_root_causes(pattern_id)` → 分析记忆中记录的失败根因
- `retrieve_best_practices(domain)` → 检索相关的最佳实践和约束条件
- `retrieve_successful_optimizations(pattern_type)` → 检索长期记忆中的成功优化案例
- `design_optimization_based_on_history(pattern_id)` → 基于历史数据设计优化方案
- `mark_optimization_assumptions(optimization_id, assumptions)` → 标记优化方案的关键假设
- `load_test_environment(test_id)` → 从工作记忆加载测试环境配置
- `record_test_results(test_data)` → 记录测试过程和结果到记忆
- `compare_with_historical_baseline(test_results)` → 对比测试结果与历史基线
- `store_optimized_pattern(pattern_data)` → 将优化后的模式存储到长期记忆
- `preserve_old_version(pattern_id, version)` → 保留旧模式版本用于回滚
- `update_memory_associations(pattern_id, new_associations)` → 更新相关记忆的关联关系
- `record_optimization_process(process_data)` → 记录优化过程和关键决策

**模式淘汰阶段**：
- `retrieve_dependent_patterns(pattern_id)` → 从长期记忆检索所有依赖该模式的其他模式
- `retrieve_citation_history(pattern_id)` → 检索该模式被引用的历史记录
- `analyze_deprecation_impact(pattern_id)` → 分析淘汰的潜在影响范围
- `archive_deprecated_pattern(pattern_data)` → 将淘汰模式归档到长期记忆
- `preserve_version_history(pattern_id)` → 保留完整的版本历史和性能数据
- `record_deprecation_info(pattern_id, reason, alternative)` → 记录淘汰的时间、原因和替代方案
- `update_memory_references(pattern_id, new_pattern_id)` → 更新所有相关记忆的引用关系
- `migrate_dependent_patterns(old_id, new_id)` → 将依赖模式迁移到替代方案
- `clear_working_memory_references(pattern_id)` → 清理工作记忆中的残留引用

**模式创建阶段**：
- `retrieve_similar_problems(problem_description)` → 从长期记忆检索相似的历史问题
- `analyze_repetition_patterns(working_memory)` → 分析工作记忆中的重复模式
- `retrieve_user_feedback_on_requirement(requirement_id)` → 检索用户对新需求的反馈记录
- `retrieve_pattern_failures(pattern_id)` → 从长期记忆检索现有模式的失败案例
- `analyze_user_complaints(pattern_id)` → 分析记忆中记录的用户抱怨
- `identify_uncovered_scenarios(pattern_id)` → 识别现有模式无法覆盖的场景
- `design_based_on_best_practices(requirement_data)` → 基于长期记忆中的最佳实践设计新模式
- `simulate_execution_in_working_memory(pattern_data)` → 在工作记忆中模拟执行效果
- `retrieve_successful_cases(domain)` → 检索相关的成功案例作为参考
- `record_iteration_results(iteration_data)` → 记录每次迭代的测试结果到记忆
- `adjust_design_based_on_feedback(pattern_id, feedback)` → 基于用户反馈调整模式设计
- `extract_key_insights(iteration_history)` → 提炼优化过程中的关键洞察
- `store_verified_pattern(pattern_data)` → 将验证通过的模式存储到长期记忆
- `establish_pattern_associations(pattern_id, related_memories)` → 建立与相关记忆的关联关系
- `record_creation_process(process_data)` → 记录创建过程和关键决策
- `mark_pattern_constraints(pattern_id, constraints)` → 标记模式的适用范围和约束条件

**进化循环阶段**：
- `extract_pattern_performance_data()` → 从长期记忆提取所有模式的性能数据
- `analyze_usage_trends(time_range)` → 分析记忆中的使用趋势和模式
- `identify_memory_anomalies()` → 识别记忆中的异常和潜在问题
- `generate_evaluation_report(evaluation_data)` → 生成评估报告并存储到长期记忆
- `update_working_pattern_weights(weight_updates)` → 实时更新工作记忆中的模式权重
- `write_execution_feedback_immediate(feedback_data)` → 将执行反馈即时写入记忆系统
- `adjust_execution_path_based_on_memory(context)` → 基于记忆数据动态调整执行路径
- `record_optimization_decision(decision_data, effect)` → 记录优化决策和效果
- `maintain_version_history_in_memory(pattern_id)` → 在长期记忆中维护完整的版本历史
- `associate_version_feedback(version_id, feedback_data)` → 关联每个版本的效果数据和用户反馈
- `preserve_rollback_version(pattern_id, version_id)` → 保留关键版本用于快速回滚
- `record_version_changes(from_version, to_version, reason)` → 记录版本间的变更和原因

**适应性调整阶段**：
- `write_to_short_term_memory(feedback_data)` → 将用户反馈立即写入短期记忆
- `retrieve_historical_feedback(feedback_type)` → 检索长期记忆中的相关历史反馈
- `adjust_working_pattern_parameters(pattern_id, params)` → 快速调整工作记忆中的模式参数
- `record_adjustment_reason(reason, expected_effect)` → 记录调整原因和预期效果
- `analyze_long_term_trends(time_range)` → 从长期记忆分析长期数据趋势
- `calculate_optimal_weights_based_on_history()` → 基于历史数据计算最优权重分配
- `update_memory_configuration_gradually(config_updates)` → 逐步更新记忆中的模式配置
- `monitor_adaptation_effect(effect_data)` → 监控适应效果并记录到记忆
- `retrieve_core_pattern_stability(pattern_id)` → 从长期记忆检索核心模式的历史稳定性
- `simulate_optimization_impact(optimization_data)` → 在工作记忆中模拟优化的潜在影响
- `preserve_rollback_path_in_memory(path_data)` → 保留完整的回滚路径在记忆中
- `record_conservative_decision(decision_data, risks)` → 记录保守决策的依据和风险

**知识传承阶段**：
- `promote_success_case_to_long_term(case_data)` → 将成功案例自动提升到长期记忆
- `establish_success_path_associations(path_id, related_memories)` → 建立成功路径的强关联网络
- `extract_best_practices_and_mark_priority(practices_data)` → 提炼最佳实践并标记为高优先级
- `review_and_validate_experience(experience_id)` → 定期回顾和验证记忆中的经验
- `record_failure_case_detailed(failure_data)` → 将失败案例详细记录到长期记忆
- `analyze_failure_patterns_extract_rules(pattern_id)` → 分析记忆中的失败模式提取规则
- `establish_failure_solution_mapping(failure_id, solution_id)` → 建立失败与解决方案的映射关系
- `mark_high_risk_scenarios(scenario_data)` → 在记忆中标记高风险场景
- `retrieve_specific_cases(case_type)` → 从长期记忆检索具体案例
- `abstract_generic_pattern_and_store(pattern_data)` → 抽象通用模式并存储到记忆
- `solidify_temporary_solution(solution_data)` → 将临时方案固化为标准模式
- `optimize_pattern_applicability_in_memory(pattern_id)` → 持续优化记忆中模式的适用范围
- `record_evolution_process(process_data, insights)` → 记录演化过程和关键洞察

### 与技能系统模块的接口

**技能识别阶段**：
- `retrieve_similar_skill_requirements(requirement_data)` → 检索类似需求的技能匹配记录
- `check_skill_availability(skill_id)` → 检查技能的可用性和状态
- `retrieve_skill_usage_history(skill_id)` → 检索技能的历史使用记录

**技能选择阶段**：
- `retrieve_skill_dependencies(skill_id)` → 检索技能的依赖关系
- `get_skill_priority(skill_id)` → 获取技能的优先级权重
- `check_skill_combination_compatibility(skill_ids)` → 检查技能组合的兼容性

**技能调用阶段**：
- `load_skill_parameters(skill_id)` → 加载技能所需的参数配置
- `check_skill_prerequisites(skill_id)` → 检查技能的前置条件
- `retrieve_skill_execution_history(skill_id)` → 检索技能的执行历史
- `record_skill_execution(execution_data)` → 记录技能执行过程
- `record_skill_result(skill_id, result_data)` → 记录技能执行结果

**技能创建阶段**：
- `identify_repetitive_task_patterns(task_history)` → 识别重复的任务模式
- `retrieve_standardization_examples(domain)` → 检索标准化流程的示例
- `retrieve_domain_knowledge(domain)` → 检索领域专业知识
- `record_new_skill_creation(skill_data)` → 记录新技能的创建过程

### 与上下文管理模块的接口

**上下文构建阶段**：
- `retrieve_context_template(context_type)` → 检索上下文模板
- `load_relevant_information(task_id)` → 加载相关信息
- `track_context_state(context_id)` → 跟踪上下文状态
- `manage_context_variables(context_id, variables)` → 管理上下文变量
- `update_context_realtime(context_id, updates)` → 实时更新上下文

**上下文切换阶段**：
- `save_context_state(context_id)` → 保存上下文状态
- `clear_context(context_id)` → 清理上下文
- `load_new_context(context_id)` → 加载新上下文
- `restore_context_state(context_id)` → 恢复上下文状态
- `integrate_context_information(old_context, new_context)` → 整合上下文信息

### 与错误处理模块的接口

**错误预防阶段**：
- `check_preconditions(operation_id)` → 检查前置条件
- `validate_parameters(parameter_data)` → 验证参数
- `check_boundary_conditions(operation_id)` → 检查边界条件
- `check_resource_availability(resource_id)` → 检查资源
- `assess_risks(operation_id)` → 评估风险

**错误恢复阶段**：
- `identify_error(error_data)` → 识别错误
- `assess_error_impact(error_id)` → 评估错误影响
- `select_recovery_strategy(error_id)` → 选择恢复策略
- `execute_recovery(recovery_plan)` → 执行恢复
- `record_recovery_experience(experience_data)` → 记录恢复经验

### 记忆系统内部接口

**记忆流转**：
- `evaluate_promotion_criteria(memory_id)` → 评估是否应该提升到长期记忆
- `promote_memory(memory_id, target_type)` → 提升记忆到目标类型
- `evaluate_degradation_criteria(memory_id)` → 评估是否应该降级或删除
- `degrade_memory(memory_id)` → 降级记忆
- `delete_memory(memory_id)` → 删除记忆

**记忆整合**：
- `extract_keywords(content)` → 提取关键词
- `search_similar_memories(keywords)` → 搜索相似记忆
- `establish_memory_association(memory_id1, memory_id2)` → 建立记忆关联
- `detect_memory_conflicts(memory_id)` → 检测记忆冲突
- `resolve_memory_conflict(conflict_id)` → 解决记忆冲突

**记忆维护**：
- `update_memory_access_count(memory_id)` → 更新记忆访问计数
- `update_memory_last_access(memory_id)` → 更新记忆最后访问时间
- `adjust_memory_importance(memory_id, delta)` → 调整记忆重要性
- `cleanup_unused_memories()` → 清理未使用的记忆
- `optimize_memory_storage()` → 优化记忆存储
