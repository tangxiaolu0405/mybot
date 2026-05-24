package client

import "cata/internal/brain"

// CollectRuntimeEnv 采集当前终端/OS 信息，供 server 注入 LLM。
func CollectRuntimeEnv() *brain.RuntimeEnv {
	e := brain.DetectRuntimeEnvFromProcess()
	return &e
}
