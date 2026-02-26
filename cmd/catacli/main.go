package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

const (
	DefaultSocketPath = ".cata/cata.sock"
)

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// help 命令不需要连接服务器
	if command == "help" || command == "--help" || command == "-h" {
		printUsage()
		os.Exit(0)
	}

	// 连接到服务器
	conn, err := connectToServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to connect to server: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure the server is running with 'cata run'\n")
		os.Exit(1)
	}
	defer conn.Close()

	// 处理命令（仅保留：发布任务、查看、ping；其余由 cataserver 内 LLM 自主决策）
	var req Request
	switch command {
	case "task":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: task requires a subcommand\n")
			fmt.Fprintf(os.Stderr, "Usage: catacli task <create|list|status> [args]\n")
			os.Exit(1)
		}
		req = Request{
			Command: "task",
			Args:    os.Args[2:],
		}

	case "ping":
		req = Request{
			Command: "ping",
			Args:    []string{},
		}

	case "skill":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: skill requires a subcommand\n")
			fmt.Fprintf(os.Stderr, "Usage: catacli skill <list|get <name>>\n")
			os.Exit(1)
		}
		skillSub := os.Args[2]
		switch skillSub {
		case "list":
			req = Request{Command: "skill_list", Args: []string{}}
			command = "skill_list"
		case "get":
			if len(os.Args) < 4 {
				fmt.Fprintf(os.Stderr, "Error: skill get requires skill name\n")
				fmt.Fprintf(os.Stderr, "Usage: catacli skill get <skill-name>\n")
				os.Exit(1)
			}
			req = Request{Command: "skill_get", Args: os.Args[3:4]}
			command = "skill_get"
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown skill subcommand: %s\n", skillSub)
			fmt.Fprintf(os.Stderr, "Usage: catacli skill <list|get <name>>\n")
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}

	// 发送请求
	reqData, err := json.Marshal(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to marshal request: %v\n", err)
		os.Exit(1)
	}

	if _, err := conn.Write(append(reqData, '\n')); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to send request: %v\n", err)
		os.Exit(1)
	}

	// 读取响应
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		fmt.Fprintf(os.Stderr, "Error: Failed to read response\n")
		os.Exit(1)
	}

	var resp Response
	if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse response: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	if resp.Success {
		if resp.Data != nil {
			// 格式化输出数据
			outputData(resp.Data, command)
		} else {
			fmt.Println(resp.Message)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", resp.Message)
		os.Exit(1)
	}
}

func connectToServer() (net.Conn, error) {
	socketPath := getSocketPath()
	return net.Dial("unix", socketPath)
}

func getSocketPath() string {
	root := findProjectRoot()
	if root != "" {
		return filepath.Join(root, ".cata", "cata.sock")
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, DefaultSocketPath)
}

// findProjectRoot 向上查找包含 go.mod 或 .git 的项目根目录（与 cata 服务端一致）
func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func outputData(data interface{}, command string) {
	switch command {
	case "skill_list":
		if entries, ok := data.([]interface{}); ok {
			fmt.Println("\n=== Skills (registry + skills-index) ===")
			for _, e := range entries {
				m, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				name, _ := m["name"].(string)
				desc, _ := m["description"].(string)
				impl, _ := m["implemented"].(bool)
				implStr := "MD (agent)"
				if impl {
					implStr = "server"
				}
				fmt.Printf("- %s [%s]\n", name, implStr)
				if desc != "" {
					fmt.Printf("  %s\n", desc)
				}
			}
			fmt.Println()
		} else {
			dataJSON, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(dataJSON))
		}
	case "skill_get":
		if m, ok := data.(map[string]interface{}); ok {
			if content, ok := m["content"].(string); ok {
				fmt.Println(content)
			} else {
				dataJSON, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(dataJSON))
			}
		} else {
			dataJSON, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(dataJSON))
		}
	case "task":
		// task 命令的输出处理
		if m, ok := data.(map[string]interface{}); ok {
			// 任务创建结果
			if taskID, ok := m["task_id"].(string); ok {
				fmt.Println("\n=== Task Execution Result ===")
				fmt.Printf("Task ID: %s\n", taskID)
				if output, ok := m["output"].(string); ok {
					fmt.Printf("Output: %s\n", output)
				}
				if learning, ok := m["learning"].(string); ok && learning != "" {
					fmt.Printf("Learning: %s\n", learning)
				}
				if success, ok := m["success"].(bool); ok {
					if success {
						fmt.Println("Status: ✓ Success")
					} else {
						fmt.Println("Status: ✗ Failed")
					}
				}
				fmt.Println()
			} else {
				// 任务状态详情（QueuedTask）
				fmt.Println("\n=== Task Status ===")
				
				if taskID, ok := m["id"].(string); ok {
					fmt.Printf("Task ID: %s\n", taskID)
				}
				
				if taskType, ok := m["type"].(string); ok {
					fmt.Printf("Type: %s\n", taskType)
				}
				
				if createdAt, ok := m["created_at"].(string); ok {
					fmt.Printf("Created At: %s\n", createdAt)
				}
				
				if createdBy, ok := m["created_by"].(string); ok {
					fmt.Printf("Created By: %s\n", createdBy)
				}
				
				if status, ok := m["status"].(string); ok {
					statusIcon := "⏳"
					switch status {
					case "completed":
						statusIcon = "✓"
					case "failed":
						statusIcon = "✗"
					case "running":
						statusIcon = "▶"
					}
					fmt.Printf("Status: %s %s\n", statusIcon, status)
				}
				
				if priority, ok := m["priority"].(float64); ok {
					fmt.Printf("Priority: %d\n", int(priority))
				}
				
				if actionPlan, ok := m["action_plan"].(map[string]interface{}); ok {
					if reason, ok := actionPlan["reason"].(string); ok && reason != "" {
						fmt.Printf("Reason: %s\n", reason)
					}
					if steps, ok := actionPlan["steps"].([]interface{}); ok && len(steps) > 0 {
						fmt.Printf("Steps:\n")
						for i, step := range steps {
							if s, ok := step.(string); ok {
								fmt.Printf("  %d. %s\n", i+1, s)
							}
						}
					}
				}
				
				if startedAt, ok := m["started_at"].(string); ok && startedAt != "" {
					fmt.Printf("Started At: %s\n", startedAt)
				}
				
				if completedAt, ok := m["completed_at"].(string); ok && completedAt != "" {
					fmt.Printf("Completed At: %s\n", completedAt)
				}
				
				if result, ok := m["result"].(map[string]interface{}); ok {
					fmt.Printf("\nResult:\n")
					if success, ok := result["success"].(bool); ok {
						fmt.Printf("  Success: %v\n", success)
					}
					if output, ok := result["output"].(string); ok && output != "" {
						fmt.Printf("  Output: %s\n", output)
					}
					if err, ok := result["error"].(string); ok && err != "" {
						fmt.Printf("  Error: %s\n", err)
					}
					if learning, ok := result["learning"].(string); ok && learning != "" {
						fmt.Printf("  Learning: %s\n", learning)
					}
				}
				
				fmt.Println()
			}
		} else if entries, ok := data.([]interface{}); ok {
			// 任务列表
			fmt.Printf("\n=== Task List (%d tasks) ===\n\n", len(entries))
			for i, entry := range entries {
				if m, ok := entry.(map[string]interface{}); ok {
					// 获取任务信息（QueuedTask 结构）
					taskID := ""
					if id, ok := m["id"].(string); ok {
						taskID = id
					}
					
					taskType := ""
					if t, ok := m["type"].(string); ok {
						taskType = t
					}
					
					createdAt := ""
					if ca, ok := m["created_at"].(string); ok {
						createdAt = ca
					}
					
					status := ""
					if s, ok := m["status"].(string); ok {
						status = s
					}
					
					createdBy := ""
					if cb, ok := m["created_by"].(string); ok {
						createdBy = cb
					}
					
					priority := 0
					if p, ok := m["priority"].(float64); ok {
						priority = int(p)
					}
					
					// 格式化输出
					fmt.Printf("[%d] %s - %s\n", i+1, createdAt, taskType)
					if taskID != "" {
						fmt.Printf("     Task ID: %s\n", taskID)
					}
					if status != "" {
						statusIcon := "⏳"
						switch status {
						case "completed":
							statusIcon = "✓"
						case "failed":
							statusIcon = "✗"
						case "running":
							statusIcon = "▶"
						}
						fmt.Printf("     Status: %s %s\n", statusIcon, status)
					}
					if priority > 0 {
						fmt.Printf("     Priority: %d\n", priority)
					}
					if createdBy != "" {
						fmt.Printf("     Created by: %s\n", createdBy)
					}
					
					// 如果有 ActionPlan，显示 reason
					if actionPlan, ok := m["action_plan"].(map[string]interface{}); ok {
						if reason, ok := actionPlan["reason"].(string); ok && reason != "" {
							fmt.Printf("     Reason: %s\n", reason)
						}
					}
					
					// 如果有结果，显示简要信息
					if result, ok := m["result"].(map[string]interface{}); ok {
						if success, ok := result["success"].(bool); ok {
							if success {
								if output, ok := result["output"].(string); ok && output != "" {
									// 只显示前 60 个字符
									if len(output) > 60 {
										output = output[:60] + "..."
									}
									fmt.Printf("     Result: %s\n", output)
								}
							} else {
								if err, ok := result["error"].(string); ok && err != "" {
									fmt.Printf("     Error: %s\n", err)
								}
							}
						}
					}
					
					fmt.Println()
				}
			}
		} else {
			// 默认 JSON 输出
			dataJSON, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(dataJSON))
		}
	default:
		// 默认 JSON 输出
		dataJSON, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(dataJSON))
	}
}

func printUsage() {
	fmt.Println("Cata CLI - 仅用于向 Cata 服务发布任务与查看（其余能力由 cataserver 内 LLM 自主决策）")
	fmt.Println()
	fmt.Println("Usage: catacli <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  task create \"<需求描述>\" [--async]  Create task by requirement; cata parses and runs (--async = queue for server)")
	fmt.Println("  task create <type> [args...] [--async]  Or by type: summarize, consolidate, recall, learn, optimize, reflect, idle, integrate")
	fmt.Println("  task list                       List recent tasks")
	fmt.Println("  task status <task-id>          Show task status")
	fmt.Println("  skill list                      List skills (registry + skills-index); implemented=true = server, false = MD only (agent)")
	fmt.Println("  skill get <name>                Get SKILL.md content for agent to execute (MD skills)")
	fmt.Println("  ping                            Check server connection")
	fmt.Println("  help                            Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  catacli task create \"帮我整理本周记忆摘要\"")
	fmt.Println("  catacli task create summarize --async")
	fmt.Println("  catacli task list")
	fmt.Println("  catacli task status <task-id>")
	fmt.Println("  catacli skill list")
	fmt.Println("  catacli skill get memory-reader")
	fmt.Println("  catacli ping")
	fmt.Println()
	fmt.Println("Note: Run 'cata run' first to start the server.")
}
