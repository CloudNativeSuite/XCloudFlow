package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Run or manage Xconfig plugins",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[Xconfig] Plugin system not implemented yet.")
	},
}
