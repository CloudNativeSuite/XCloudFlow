// cmd/playbook.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"craftweave/core/executor"
	"craftweave/core/parser"
)

var inventoryPath string

var playbookCmd = &cobra.Command{
	Use:   "playbook [file]",
	Short: "Run a Xconfig playbook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		fmt.Printf("📜 Executing playbook: %s\n", file)

		plays, err := parser.LoadPlaybook(file)
		if err != nil {
			fmt.Printf("❌ Failed to load playbook: %v\n", err)
			os.Exit(1)
		}

		exec := executor.New(AggregateOutput, CheckMode, DiffMode)
		exec.MaxWorkers = MaxWorkers
		exec.Execute(plays, inventoryPath)
	},
}

func init() {
	playbookCmd.Flags().StringVarP(&inventoryPath, "inventory", "i", "hosts.yaml", "Inventory file")
	playbookCmd.Flags().IntVarP(&MaxWorkers, "forks", "f", 5, "Max parallel tasks")
	playbookCmd.Flags().BoolVarP(&AggregateOutput, "aggregate", "A", false, "Aggregate output from identical results")
	playbookCmd.Flags().BoolVarP(&CheckMode, "check", "C", false, "Dry-run mode")
	rootCmd.AddCommand(playbookCmd)
}
