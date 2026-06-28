package openclaw

import (
	"fmt"
	"os"
)

func runList(args []string) int {
	if wantsHelp(args) {
		printHelp(listHelpText)
		return 0
	}
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: my openclaw list\n")
		return 1
	}

	reg, err := loadRegistry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if len(reg.DataDirs) == 0 {
		fmt.Println("(no data dirs registered)")
		return 0
	}

	fmt.Println("path\tnote\tadded_at")
	for _, entry := range reg.DataDirs {
		fmt.Printf("%s\t%s\t%s\n", entry.Path, entry.Note, entry.AddedAt)
	}
	return 0
}