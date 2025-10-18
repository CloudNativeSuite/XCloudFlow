package main

import (
	"fmt"

	"craftweave/cmd"
)

func main() {
	fmt.Println("🧶 欢迎使用：Xconfig - 任务与架构编织工具")
	cmd.Execute() // ✅ 正确方式：不接收返回值
}
