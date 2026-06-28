package openclaw

import (
	"errors"
	"fmt"
	"os"
	"time"

	lessflags "github.com/xhd2015/less-flags"
)

func runAdd(args []string) int {
	if len(args) == 0 || wantsHelp(args) {
		printHelp(addHelpText)
		return 0
	}
	if args[0] != "data-dir" {
		fmt.Fprintf(os.Stderr, "usage: my openclaw add data-dir <path> [--note \"...\"]\n")
		return 1
	}
	if len(args) == 1 || wantsHelp(args[1:]) {
		printHelp(addHelpText)
		return 0
	}

	var note string
	remain, err := lessFlags.String("--note", &note).
		Help("-h,--help", addHelpText).
		HelpNoExit().
		Parse(args[2:])
	if errors.Is(err, lessflags.ErrHelp) {
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if len(remain) != 0 {
		fmt.Fprintf(os.Stderr, "usage: my openclaw add data-dir <path> [--note \"...\"]\n")
		return 1
	}

	absPath, err := resolvePath(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: data directory not found: %s\n", absPath)
			return 1
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "error: path is not a directory: %s\n", absPath)
		return 1
	}

	reg, err := loadRegistry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if _, existing := findDataDir(reg, absPath); existing != nil {
		fmt.Fprintf(os.Stderr, "data dir already registered: %s\n", absPath)
		existing.Note = note
	} else {
		reg.DataDirs = append(reg.DataDirs, DataDirEntry{
			Path:    absPath,
			Note:    note,
			AddedAt: time.Now().Format(time.RFC3339),
		})
	}

	if err := saveRegistry(reg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}