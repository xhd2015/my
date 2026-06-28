package openclaw

import (
	"fmt"
	"os"
	"os/exec"
)

func runExec(execArgs []string) int {
	if len(execArgs) == 0 {
		fmt.Fprintf(os.Stderr, "error: --exec requires a command\n")
		return 1
	}

	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	containerName := defaultContainerName
	running, err := containerIsRunning(containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if !running {
		fmt.Fprintf(os.Stderr, "error: container %s is not running (start it with 'my openclaw run-in-podman --data-dir <path>' first)\n", containerName)
		return 1
	}

	if err := podmanExecInteractivePreviewed(containerName, execArgs...); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}