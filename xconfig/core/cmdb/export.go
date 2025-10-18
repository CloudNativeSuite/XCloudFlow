package cmdb

import (
	"encoding/json"
	"os"
)

type Node struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Labels     []string          `json:"labels,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func ExportGraph(filename string, graph Graph) error {
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
