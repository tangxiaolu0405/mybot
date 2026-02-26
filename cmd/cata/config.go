package main

import (
	"encoding/json"
	"fmt"
	"os"

	"mybot/internal/config"
)

func handleConfigCommand(args []string) {
	if len(args) < 1 {
		printConfigUsage()
		os.Exit(1)
	}

	subcommand := args[0]

	switch subcommand {
	case "show":
		handleConfigShow()
	case "set":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: config set requires key and value\n")
			fmt.Fprintf(os.Stderr, "Usage: cata config set <key> <value>\n")
			os.Exit(1)
		}
		handleConfigSet(args[1], args[2])
	case "get":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: config get requires key\n")
			fmt.Fprintf(os.Stderr, "Usage: cata config get <key>\n")
			os.Exit(1)
		}
		handleConfigGet(args[1])
	case "edit":
		handleConfigEdit()
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown config subcommand: %s\n", subcommand)
		printConfigUsage()
		os.Exit(1)
	}
}

func printConfigUsage() {
	fmt.Println("Configuration Management")
	fmt.Println()
	fmt.Println("Usage: cata config <subcommand> [args]")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  show              Show current configuration")
	fmt.Println("  get <key>         Get a configuration value")
	fmt.Println("  set <key> <value> Set a configuration value")
	fmt.Println("  edit              Open config file in editor")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cata config show")
	fmt.Println("  cata config get brain.dir")
	fmt.Println("  cata config set llm.provider qwen")
	fmt.Println("  cata config set llm.api_key sk-xxx")
	fmt.Println("  cata config set llm.api_url https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions")
	fmt.Println("  cata config set llm.model qwen-turbo")
	fmt.Println("  cata config set evolution.enabled false")
}

func handleConfigShow() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 隐藏敏感信息
	displayConfig := *cfg
	if displayConfig.LLM.APIKey != "" {
		displayConfig.LLM.APIKey = "***hidden***"
	}

	data, err := json.MarshalIndent(displayConfig, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}

func handleConfigGet(key string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	value := getConfigValue(cfg, key)
	if value == nil {
		fmt.Fprintf(os.Stderr, "Error: key not found: %s\n", key)
		os.Exit(1)
	}

	fmt.Println(value)
}

func handleConfigSet(key, value string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := setConfigValue(cfg, key, value); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting config: %v\n", err)
		os.Exit(1)
	}

	if err := config.SaveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration updated: %s = %s\n", key, value)
}

func handleConfigEdit() {
	configPath := config.GetConfigPath()
	fmt.Printf("Config file: %s\n", configPath)
	fmt.Println("Please edit the file manually or use 'cata config set' command")
}

// getConfigValue 获取配置值（支持嵌套键，如 "brain.dir"）
func getConfigValue(cfg *config.AppConfig, key string) interface{} {
	switch key {
	case "brain.dir":
		return cfg.Brain.Dir
	case "brain.base_dir":
		return cfg.Brain.BaseDir
	case "llm.provider":
		return cfg.LLM.Provider
	case "llm.api_key":
		if cfg.LLM.APIKey != "" {
			return "***hidden***"
		}
		return ""
	case "llm.api_url":
		return cfg.LLM.APIURL
	case "llm.model":
		return cfg.LLM.Model
	case "llm.max_tokens":
		return cfg.LLM.MaxTokens
	case "llm.timeout":
		return cfg.LLM.Timeout
	case "llm.enabled":
		return cfg.LLM.Enabled
	case "server.socket_path":
		return cfg.Server.SocketPath
	case "server.log_level":
		return cfg.Server.LogLevel
	case "evolution.enabled":
		return cfg.Evolution.Enabled
	case "evolution.cycle_interval":
		return cfg.Evolution.CycleInterval
	case "evolution.task_queue_interval":
		return cfg.Evolution.TaskQueueInterval
	default:
		return nil
	}
}

// setConfigValue 设置配置值（支持嵌套键）
func setConfigValue(cfg *config.AppConfig, key, value string) error {
	switch key {
	case "brain.dir":
		cfg.Brain.Dir = value
	case "brain.base_dir":
		cfg.Brain.BaseDir = value
	case "llm.provider":
		cfg.LLM.Provider = value
	case "llm.api_key":
		cfg.LLM.APIKey = value
		cfg.LLM.Enabled = value != ""
	case "llm.api_url":
		cfg.LLM.APIURL = value
	case "llm.model":
		cfg.LLM.Model = value
	case "llm.max_tokens":
		var v int
		if _, err := fmt.Sscanf(value, "%d", &v); err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		cfg.LLM.MaxTokens = v
	case "llm.timeout":
		var v int
		if _, err := fmt.Sscanf(value, "%d", &v); err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		cfg.LLM.Timeout = v
	case "llm.enabled":
		cfg.LLM.Enabled = value == "true" || value == "1"
	case "server.socket_path":
		cfg.Server.SocketPath = value
	case "server.log_level":
		cfg.Server.LogLevel = value
	case "evolution.enabled":
		cfg.Evolution.Enabled = value == "true" || value == "1"
	case "evolution.cycle_interval":
		var v int
		if _, err := fmt.Sscanf(value, "%d", &v); err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		cfg.Evolution.CycleInterval = v
	case "evolution.task_queue_interval":
		var v int
		if _, err := fmt.Sscanf(value, "%d", &v); err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		cfg.Evolution.TaskQueueInterval = v
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}
