package main

import (
	"fmt"
	"os"

	"xcloud-cli/cmd"
)

func main() {

	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "💥 panic: %v\n", err)
		os.Exit(2)
	}
	cmd.Execute()
}
