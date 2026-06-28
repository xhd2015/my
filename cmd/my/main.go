package main

import (
	"fmt"
	"os"

	"github.com/xhd2015/my/internal/openclaw"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(openclaw.TopHelp())
		return
	}

	switch args[0] {
	case "openclaw":
		os.Exit(openclaw.Run(args[1:]))
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		os.Exit(1)
	}
}