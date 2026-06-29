package opencode

import (
	"fmt"
	"os"
)

const opencodeHelpText = `Usage: my opencode [options]

Options:
  --import-local-grok   Import local Grok OIDC credentials into auth.json
  -h, --help            Show this help
`

func Run(args []string) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(opencodeHelpText)
		return 0
	}

	for _, arg := range args {
		switch arg {
		case "--import-local-grok":
			return ImportLocalGrok()
		case "-h", "--help":
			fmt.Print(opencodeHelpText)
			return 0
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n", arg)
			return 1
		}
	}

	fmt.Print(opencodeHelpText)
	return 0
}
