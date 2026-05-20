# MyBot Project Documentation

## Overview
MyBot is an advanced chatbot designed to assist users with various tasks by utilizing machine learning and natural language processing techniques. The bot aims to provide seamless interactions and enhance user productivity.

## Brain & 自主演进

- **核心**：认知与记忆以 `brain/core.md` 为核心，资源路径、技能、记忆规则、进化机制见其中。
- **自主演进**：分析状态→LLM 决策→任务执行→学习反馈，见 `brain/workflow.md`。在 Agent 内通过 skills（memory-reader、memory-iteration-manager、task-evolution-executor）即可完成同等能力。
- **Cata**：单二进制。直接 `cata` 会自动后台拉起 server，退出 chat（`/exit` 或最后一个会话结束）后关闭 managed server；`cata run` 仅用于手动常驻（**catacli 已废弃**）。构建：`go build -o cata ./cmd/cata`；运行：`./cata init`、`./cata`。说明见 `agents.md` 与 `brain/core.md`。

## Features
- Conversational AI: Engage in human-like conversations.
- Task Automation: Perform automated tasks based on user commands.
- Multi-Platform Support: Available on various platforms like Slack, Discord, and web applications.

## Installation
### Prerequisites
- Python 3.6+
- Node.js and npm

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/tangxiaolu0405/mybot.git
   ```
2. Navigate into the directory:
   ```bash
   cd mybot
   ```
3. Install the dependencies:
   ```bash
   npm install
   pip install -r requirements.txt
   ```
4. Start the bot:
   ```bash
   node index.js
   ```

## Usage
- To interact with the bot, simply type your queries or commands in the chat interface.
- Use specific commands to access various features.

## Contributing
1. Fork the repository.
2. Create a new branch (e.g., `feature-xyz`).
3. Make your changes and commit them.
4. Push to your fork and submit a pull request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact
For any inquiries, please reach out to [your-email@example.com](mailto:your-email@example.com).