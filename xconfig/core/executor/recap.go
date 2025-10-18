package executor

import (
	"fmt"
	"sort"
)

// hostStats stores counts of task results per host.
type hostStats struct {
	OK          int
	Changed     int
	Failed      int
	Skipped     int
	Unreachable int
	Rescued     int
	Ignored     int
}

var (
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func colorize(val int, color string) string {
	return fmt.Sprintf("%s%d%s", color, val, colorReset)
}

// printRecap displays a PLAY RECAP summary similar to Ansible.
func printRecap(stats map[string]*hostStats) {
	if len(stats) == 0 {
		return
	}
	fmt.Println("\nPLAY RECAP ****************************************************************")
	hosts := make([]string, 0, len(stats))
	for h := range stats {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)
	for _, h := range hosts {
		s := stats[h]
		okStr := colorize(s.OK, colorGreen)
		changedStr := colorize(s.Changed, colorYellow)
		failedStr := colorize(s.Failed, colorRed)
		fmt.Printf("%-20s : ok=%s changed=%s unreachable=%d failed=%s skipped=%d rescued=%d ignored=%d\n",
			h, okStr, changedStr, s.Unreachable, failedStr, s.Skipped, s.Rescued, s.Ignored)
	}
}
