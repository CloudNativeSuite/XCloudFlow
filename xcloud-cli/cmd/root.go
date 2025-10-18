package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var env string
var configPath string
var targetCloud string
var targetRegion string

var rootCmd = &cobra.Command{
	Use:   "xcloud",
	Short: "🧰 xcloud-cli - 多环境自动化管理器 (Go + Pulumi Native)",
	Long: `📖 用法:
  xcloud [命令]
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
		if targetCloud != "" {
			os.Setenv("STACK_CLOUD", targetCloud)
		}
		if targetRegion != "" {
			os.Setenv("STACK_REGION", targetRegion)
		}

		fmt.Println("🌍 当前环境:", env)
		fmt.Println("📁 当前配置路径:", configPath)
		fmt.Println("🔐 Pulumi 密码文件已加载:", os.Getenv("HOME")+"/.pulumi-passphrase")
		if targetCloud != "" {
			fmt.Println("☁️  目标云:", targetCloud)
		}
		if targetRegion != "" {
			fmt.Println("📍 目标区域:", targetRegion)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&env, "env", "sit", "指定环境")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./config/", "指定配置路径")
	rootCmd.PersistentFlags().StringVar(&targetCloud, "cloud", "", "仅部署指定云 (matrix 覆盖)")
	rootCmd.PersistentFlags().StringVar(&targetRegion, "region", "", "仅部署指定区域")
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
		fmt.Println("  STACK_CLOUD=aws")
		fmt.Println("  STACK_REGION=ap-northeast-1")
		fmt.Println("\nexample:")
		fmt.Println("    STACK_ENV=prod CONFIG_PATH=config/ xcloud up")
		fmt.Println("    Or")
		fmt.Println("    xcloud up --config <path>/config/sit --cloud aws --region ap-northeast-1")
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
}
