# Cata — 终端原生 AI Agent

Go 编写的终端个人 AI 助手。单二进制，Unix socket 架构，后台自主演进记忆。

## 快速开始

```bash
# 构建
go build -o cata ./cmd/cata

# 初始化脑子
./cata init

# 开始对话
./cata
```

## 架构

```
cata chat ──Unix Socket──▶ cata run (server) ──HTTP──▶ LLM (DeepSeek / OpenAI 兼容)
                                │
                                ├── 脑子 (~/.cata/brain/workspaces/<id>/)
                                └── 后台演进 (evolve engine, 默认 600s)
```

## 配置

`~/.cata/config.json`：

```json
{
  "llm": {
    "provider": "deepseek",
    "model": "deepseek-v4-flash",
    "api_key": "sk-..."
  },
  "exec": { "enabled": true },
  "evolution": { "enabled": true, "cycle_interval": 600 }
}
```

## 目录

| 位置 | 用途 |
|------|------|
| `cmd/cata/` | CLI 入口 (`chat`, `init`, `run`, `config`) |
| `internal/server/` | Unix socket 服务端，聊天循环，工具执行 |
| `internal/client/` | 终端客户端，REPL，事件渲染 |
| `internal/llm/` | OpenAI 兼容 LLM 客户端 |
| `internal/brain/` | 脑子路径、工作区解析、上下文组装 |
| `internal/evolve/` | 后台自主演进引擎 |
| `internal/config/` | 配置加载与校验 |
| `skills/` | 结晶化技能示例 |
| `brain/` | 脑子模板种子（仓库模板，非运行时） |

## 设计文档

- `agents.md` — 项目边界与 AI 约束
- `design.md` — 完整系统设计（架构、交互层、产出区、记忆分层）
- `brain/constraints.md` — 行为宪法（种子 → `global/constraints.md`）
- `brain/behavior.md` — 演进 SOP（种子 → `global/behavior.md`）
- `brain/boot-assembler.md` — 运行时引导（种子 → `global/boot-assembler.md`）

## 依赖

仅需 Go 1.21+。无 Python、Node.js 依赖。

## License

MIT
