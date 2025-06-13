package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var localPath string
var dbConfigFile string

type dbConfig struct {
	DSN string `yaml:"dsn"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "⚙️ 初始化依赖或环境",
	RunE: func(cmd *cobra.Command, args []string) error {
		if localPath != "" && dbConfigFile != "" {
			return fmt.Errorf("cannot use --local with --dbconfig")
		}

		if localPath == "" && dbConfigFile == "" {
			// default behavior: run go mod tidy
			fmt.Println("🔧 初始化依赖...")
			c := exec.Command("go", "mod", "tidy")
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		}

		if localPath != "" {
			return initLocal(localPath)
		}

		return initDB(dbConfigFile)
	},
}

func initLocal(dir string) error {
	dir = filepath.Clean(os.ExpandEnv(dir))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	files := []string{"Pulumi.yaml", "Pulumi.sit.yaml"}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		dst := filepath.Join(dir, f)
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
	}
	fmt.Printf("✅ 初始化文件写入 %s\n", dir)
	return nil
}

func initDB(cfgPath string) error {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	var cfg dbConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if cfg.DSN == "" {
		return fmt.Errorf("missing dsn in %s", cfgPath)
	}

	store := map[string]string{}
	if b, err := os.ReadFile(cfg.DSN); err == nil {
		_ = json.Unmarshal(b, &store)
	}

	files := []string{"Pulumi.yaml", "Pulumi.sit.yaml"}
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		store[f] = string(content)
	}

	out, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(cfg.DSN, out, 0o644); err != nil {
		return err
	}
	fmt.Printf("✅ 初始化文件写入 DB (%s)\n", cfg.DSN)
	return nil
}

func init() {
	initCmd.Flags().StringVar(&localPath, "local", "", "本地初始化路径")
	initCmd.Flags().StringVar(&dbConfigFile, "dbconfig", "", "数据库配置文件")
}
