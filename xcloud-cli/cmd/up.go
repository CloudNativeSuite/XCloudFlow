package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"xcloud-cli/internal/modules"
	"xcloud-cli/internal/pulumi"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "🚀 部署资源",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 正在部署资源...")

		if err := modules.ExecuteTask(context.Background(), pulumi.DeployTask{}); err != nil {
			fmt.Println("❌ 部署失败:", err)
			os.Exit(1)
		}

		fmt.Println("✅ 部署完成")
	},
}
