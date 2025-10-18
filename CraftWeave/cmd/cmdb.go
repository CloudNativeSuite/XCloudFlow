package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var cmdbCmd = &cobra.Command{
	Use:   "cmdb",
	Short: "Export architecture as a graph model",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[Xconfig] CMDB export not implemented yet.")
	},
}
