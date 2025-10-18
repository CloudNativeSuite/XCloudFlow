package ssh

import (
	"fmt"
	"sort"
	"strings"
)

func AggregatedPrint(results []CommandResult) {
	type key struct {
		ReturnMsg  string
		ReturnCode int
		Output     string
	}

	grouped := make(map[key][]string)

	for _, r := range results {
		k := key{
			ReturnMsg:  r.ReturnMsg,
			ReturnCode: r.ReturnCode,
			Output:     r.Output,
		}
		grouped[k] = append(grouped[k], r.Host)
	}

	for k, hosts := range grouped {
		sort.Strings(hosts)
		fmt.Printf("%s | %s | rc=%d >>\n", strings.Join(hosts, ","), k.ReturnMsg, k.ReturnCode)
		fmt.Print(k.Output)
	}
}

