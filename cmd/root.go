package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var env string
var configPath string

var rootCmd = &cobra.Command{
	Use:   "PulumiGo",
	Short: "🧰 PulumiGo - 多环境自动化管理器 (Go + Pulumi Native)",
	Long: `📖 用法:    
  PulumiGo [命令]
支持命令:
  init      ⚙️ 初始化依赖
  up        🚀 部署资源
  down      🔥 销毁资源
  export    📤 导出 stack 状态
  import    📥 导入 stack 状态
  ansible   🧪 执行 ansible-playbook
  help      📖 显示帮助`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Ensure other modules relying on environment variables can
		// access the CLI provided values.
		os.Setenv("STACK_ENV", env)
		os.Setenv("CONFIG_PATH", configPath)

		fmt.Println("🌍 当前环境:", env)
		fmt.Println("📁 当前配置路径:", configPath)
		fmt.Println("🔐 Pulumi 密码文件已加载:", os.Getenv("HOME")+"/.pulumi-passphrase")
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&env, "env", "sit", "指定环境")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./config/", "指定配置路径")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(ansibleCmd)

	// Customize root help output to avoid repeating Cobra's default sections
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
		fmt.Println("\nFlags:")
		cmd.Flags().PrintDefaults()
		fmt.Println("ENV:")
		fmt.Println("  STACK_ENV=prod")
		fmt.Println("  CONFIG_PATH=<path>config")
		fmt.Println("\nexample:")
		fmt.Println("    STACK_ENV=prod CONFIG_PATH=config/  PulumiGo up")
		fmt.Println("    Or")
		fmt.Println("    PulumiGo up --config <path>/config/sit")
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
}
