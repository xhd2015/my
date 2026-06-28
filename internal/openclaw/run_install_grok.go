package openclaw

import (
	"fmt"
	"os"
	"os/exec"
)

func runInstallGrok(containerName string) int {
	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	running, err := containerIsRunning(containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if !running {
		fmt.Fprintf(os.Stderr, "error: container %s is not running (start it with 'my openclaw run-in-podman --data-dir <path>' first)\n", containerName)
		return 1
	}

	plan, err := planGrokInstall(containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Installing Grok CLI in container %s\n", containerName)
	printCommand(plan.previewCommand())

	if err := podmanRun(plan.execArgs...); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Println()
	fmt.Println("Grok CLI installed. Run inside the container:")
	fmt.Println(grokRunHint(containerName))
	fmt.Println()
	fmt.Println("Tip: import local OAuth first so the installer can reuse ~/.grok/auth.json:")
	fmt.Println("  my openclaw run-in-podman --import-local-grok")
	return 0
}