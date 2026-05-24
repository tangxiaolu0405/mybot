# Windows 用户配置指南

## 1. 前置条件

- **Go 1.21+**：[下载安装](https://go.dev/dl/)，装完确认：

```powershell
go version
```

## 2. 构建

```powershell
cd D:\project\mybot
go build -o cata.exe .\cmd\cata
```

## 3. 配置

配置文件位置：`%USERPROFILE%\.cata\config.json`（即 `C:\Users\<你的用户名>\.cata\config.json`）。

### 方式一：自动初始化（推荐）

```powershell
.\cata.exe init
```

这会自动创建 `~\.cata\` 目录结构和默认配置文件。

### 方式二：手动创建

在 `%USERPROFILE%\.cata\config.json` 写入：

```json
{
  "llm": {
    "provider": "deepseek",
    "model": "deepseek-v4-flash",
    "api_key": "sk-你的key"
  },
  "server": {
    "timezone": "Asia/Shanghai"
  },
  "exec": {
    "enabled": true
  },
  "evolution": {
    "enabled": true,
    "cycle_interval": 600
  }
}
```

## 4. 通过环境变量配置

不想写配置文件到磁盘时，用环境变量（PowerShell）：

```powershell
# 方式 A：DeepSeek（国内推荐，内置默认）
$env:DEEPSEEK_API_KEY = "sk-你的key"
.\cata.exe init
.\cata.exe

# 方式 B：OpenAI 兼容接口
$env:OPENAI_API_KEY = "sk-你的key"
$env:LLM_API_URL = "https://你的代理地址/v1/chat/completions"
$env:LLM_MODEL = "gpt-4o"
.\cata.exe init
.\cata.exe

# 方式 C：阿里通义
$env:DASHSCOPE_API_KEY = "sk-你的key"
.\cata.exe init
```

支持的环境变量一览：

| 环境变量 | 作用 |
|---|---|
| `CATA_HOME` | 覆盖默认的 `~\.cata` 目录 |
| `CATA_CONFIG_FILE` | 指定配置文件路径 |
| `DEEPSEEK_API_KEY` | DeepSeek API Key |
| `OPENAI_API_KEY` | OpenAI API Key |
| `ANTHROPIC_API_KEY` | Anthropic API Key |
| `DASHSCOPE_API_KEY` | 阿里通义 API Key |
| `LLM_API_URL` | 自定义 API 地址 |
| `LLM_MODEL` | 模型名 |
| `LLM_PROVIDER` | 提供商：`deepseek`/`openai`/`claude`/`qwen` |
| `CATA_EXEC_ENABLED` | `1`/`true` 启用命令执行 |

## 5. Windows 特别说明

### 时区

Windows 下一定要设 `server.timezone`，否则时间戳会乱：

```json
{
  "server": {
    "timezone": "Asia/Shanghai"
  }
}
```

### 路径格式

配置文件里的路径用正斜杠 `/` 或双反斜杠 `\\` 都可以：

```json
{
  "brain": {
    "dir": "D:/mybot-data/brain",
    "base_dir": "D:/projects/mybot"
  }
}
```

### 终端 ANSI 颜色

代码里已自动启用 `ENABLE_VIRTUAL_TERMINAL_PROCESSING`，Windows Terminal / PowerShell 7 直接用。旧版 `cmd.exe` 或未开启 VT 的 conhost 无颜色，建议用 Windows Terminal。

### 命令白名单

默认已包含 Windows 命令（`dir`、`cmd`、`powershell`、`wsl`、`code` 等），不需要额外配置。

## 6. 验证

```powershell
.\cata.exe
```

输入一句话测试 LLM 连通性。没问题就说明配置成功。

> **注意**：当前版本服务端走 Unix socket（`net.Listen("unix", ...)`），Windows 原生不支持。如果启动报 socket 相关错误，说明这一块还在移植中 —— 可以用 WSL2 跑，或者等后续 Windows named pipe 适配。
