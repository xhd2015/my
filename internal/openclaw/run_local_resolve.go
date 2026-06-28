package openclaw

import (
	"fmt"
	"strings"
)

func resolveRunDataDir(explicit string) (string, bool, error) {
	if explicit != "" {
		return explicit, false, nil
	}

	reg, err := loadRegistry()
	if err != nil {
		return "", false, err
	}

	if len(reg.DataDirs) == 0 {
		return "", false, fmt.Errorf("requires --data-dir <dir>, or add with 'my openclaw add data-dir <dir>'")
	}

	var running []string
	for _, entry := range reg.DataDirs {
		ok, err := localGatewayRunning(entry.Path)
		if err != nil {
			return "", false, err
		}
		if ok {
			running = append(running, entry.Path)
		}
	}

	switch len(running) {
	case 0:
		return "", false, formatRegisteredDataDirsError(reg)
	case 1:
		return running[0], true, nil
	default:
		return "", false, formatMultipleRunningError(running)
	}
}

func formatRegisteredDataDirsError(reg *Registry) error {
	var b strings.Builder
	fmt.Fprintf(&b, "requires --data-dir <dir> (running) or 'my openclaw run --data-dir <dir>'\n")
	fmt.Fprintf(&b, "Registered data dirs:\n")
	for _, entry := range reg.DataDirs {
		fmt.Fprintf(&b, "  %s\n", entry.Path)
	}
	return fmt.Errorf("%s", strings.TrimSuffix(b.String(), "\n"))
}

func formatMultipleRunningError(running []string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "multiple local gateways running; pass --data-dir explicitly\n")
	fmt.Fprintf(&b, "Running data dirs:\n")
	for _, path := range running {
		fmt.Fprintf(&b, "  %s\n", path)
	}
	return fmt.Errorf("%s", strings.TrimSuffix(b.String(), "\n"))
}