package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "catacli 已废弃，请使用：")
	fmt.Fprintln(os.Stderr, "  cata        # 默认进入对话")
	fmt.Fprintln(os.Stderr, "  cata chat")
	fmt.Fprintln(os.Stderr, "  cata run    # 先启动后台 server（每台机器只需一个）")
	os.Exit(1)
}
