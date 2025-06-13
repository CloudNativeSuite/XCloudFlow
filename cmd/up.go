package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"PulumiGo/internal/modules"
	"PulumiGo/internal/pulumi"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "🚀 部署资源",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 正在部署资源...")

		pool := modules.NewPool(2)

		pool.Submit(func() {
			c := exec.Command("pulumi", "up", "--stack", env, "--non-interactive", "--yes")
			c.Env = append(os.Environ(), "PULUMI_CONFIG_PASSPHRASE_FILE="+os.Getenv("HOME")+"/.pulumi-passphrase")
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			_ = c.Run()
		})

		pool.Submit(func() {
			if err := modules.ExecuteTask(context.Background(), pulumi.DeployTask{}); err != nil {
				fmt.Println("❌ 部署失败:", err)
				os.Exit(1)
			}
		})

		pool.Wait()
		fmt.Println("✅ 部署完成")
	},
}
