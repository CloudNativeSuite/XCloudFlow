package cmd

import (
	"fmt"
	"strings"
	"sync"

	"craftweave/core/executor"
	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
	"github.com/spf13/cobra"
)

var module, args string

var remoteCmd = &cobra.Command{
	Use:   "remote [target]",
	Short: "Run ad-hoc tasks on target hosts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, targets []string) {
		hosts, err := inventory.Parse(InventoryPath, targets[0])
		if err != nil {
			fmt.Println("Failed to parse inventory:", err)
			return
		}

		task := parser.Task{}
		switch module {
		case "shell":
			task.Shell = args
		case "script":
			task.Script = args
		case "template":
			parts := strings.SplitN(args, ":", 2)
			if len(parts) == 2 {
				task.Template = &parser.Template{Src: parts[0], Dest: parts[1]}
			}
		}

		collector := &executor.MemoryCollector{}
		var wg sync.WaitGroup
		sem := make(chan struct{}, MaxWorkers)

		for _, h := range hosts {
			wg.Add(1)
			go func(h inventory.Host) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				if CheckMode {
					fmt.Printf("%s | SKIPPED\n", h.Name)
					return
				}

				res := executor.ExecuteTask(task, h, nil, DiffMode)
				collector.Collect(res)
			}(h)
		}
		wg.Wait()

		results := collector.Results

		if AggregateOutput {
			ssh.AggregatedPrint(results)
		} else {
			for _, r := range results {
				fmt.Printf("%s | %s | rc=%d >>\n%s\n", r.Host, r.ReturnMsg, r.ReturnCode, r.Output)
			}
		}
	},
}

func init() {
	remoteCmd.Flags().StringVarP(&InventoryPath, "inventory", "i", "hosts.yaml", "Inventory file")
	remoteCmd.Flags().StringVarP(&module, "module", "m", "shell", "Module to execute")
	remoteCmd.Flags().StringVarP(&args, "args", "a", "", "Arguments for the module")
	remoteCmd.Flags().IntVarP(&MaxWorkers, "forks", "f", 5, "Max parallel tasks")
	remoteCmd.Flags().BoolVarP(&CheckMode, "check", "C", false, "Check mode (dry-run)")
	remoteCmd.Flags().BoolVarP(&AggregateOutput, "aggregate", "A", false, "Aggregate identical output")
	rootCmd.AddCommand(remoteCmd)
}
