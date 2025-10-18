package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Encrypt/decrypt secrets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[Xconfig] Vault feature is not implemented yet.")
	},
}
