# 长期记忆存储

## 项目知识

### 项目架构
- 项目名称: CLAW (Cognitive Learning and Adaptive Workspace)
- 核心目标: 构建具有学习能力和自适应机制的AI代理系统
- 设计哲学: 好品味、实用主义、简洁执念

### 代码模式
- 链表删除优化: 10行带if判断 → 4行无条件分支
- 消除边界情况: 使特殊情况成为正常情况
- 函数设计: 短小精悍，只做一件事

### 最佳实践
- 向后兼容性不可破坏
- 解决实际问题而非假想威胁
- 避免超过3层缩进
- 优先组合而非继承

## 用户偏好

### 代码风格
- 语言: Python/JavaScript
- 风格: 简洁、实用、高效
- 注释: 最少化，代码自解释

### 工作习惯
- 优先使用现有工具
- 避免过度工程化
- 快速迭代，持续改进

### AI交互偏好
- AI需要完全自主，无需用户二次确认
- AI需要自主检查技能，不存在则创建
- AI需要自主迭代记忆系统

## 技术栈

### 核心技术
- AI框架: Claude API
- 存储: Markdown文件系统
- 版本控制: Git

### 开发工具
- IDE: Trae
- 包管理: npm/pip
- 测试: pytest/jest

## 资源路径

### 核心文件
- 核心思维模式: brain/core.md
- 项目Agent配置: agent.md

### 记忆系统
- 长期记忆: brain/memory/long-term/
- 短期记忆: brain/memory/short-term/current_session.md

### 技能系统
- 技能目录: skills/
- 技能索引: skills/skills-index.json
- 技能规范: skills/skills-creator/SKILL.md

### 基本技能
- 技能列表读取: skills/skill-list-reader/SKILL.md
- 技能验证: skills/skill-validator/SKILL.md
- 技能创建: skills/skills-creator/SKILL.md
- 记忆读取: skills/memory-reader/SKILL.md

## 历史决策

### 架构决策
- 采用模块化设计
- 使用插件式技能系统
- 实现分层记忆架构
- 定义AI自主决策流程

### 技术选型
- 选择Claude API而非自研模型
- 使用Markdown而非数据库存储
- 采用异步处理机制

### 功能决策
- 创建memory-reader技能用于记忆管理
- 定义AI自主决策流程：接收任务→技能检查→执行任务→任务完成
- 在core.md中添加资源路径定义

## 经验教训

### 成功案例
- 简化复杂逻辑提升可维护性
- 消除边界情况减少bug
- 模块化设计提高扩展性
- AI自主决策流程提高效率

### 失败教训
- 过度设计导致复杂度增加
- 忽视向后兼容性破坏用户体验
- 过早优化影响开发进度

## AI自主决策流程

### 1. 接收任务
理解用户需求 → 识别任务类型 → 收集上下文 → 检索记忆 → 匹配技能

### 2. 技能检查
检查技能是否存在 → 读取技能文件 → 检查依赖关系 → 执行技能 或 创建新技能

### 3. 执行任务
分步执行 → 实时验证 → 错误处理 → 确认结果 → 更新记忆

### 4. 任务完成
迭代记忆系统 → 更新短期记忆 → 更新长期记忆 → 更新核心思维模式
