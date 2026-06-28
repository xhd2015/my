package openclaw

import (
	"fmt"
	"os"
	"os/exec"
)

func requireRunningLocalGateway(dataDir string) (string, error) {
	running, err := localGatewayRunning(dataDir)
	if err != nil {
		return "", err
	}
	if !running {
		return "", fmt.Errorf("local gateway is not running (start it with 'my openclaw run --data-dir <path>' first)")
	}
	if dataDir == "" {
		return "", fmt.Errorf("--data-dir is required")
	}
	return resolvePath(dataDir)
}

func resolveLocalGatewayPort(port string) string {
	if port != "" {
		return port
	}
	return defaultPort
}

func resolveLocalGatewayToken(dataDir string) string {
	if dataDir == "" {
		return ""
	}
	if sources, err := resolveTokenSources(dataDir); err == nil {
		return sources.Effective()
	}
	return ""
}

func runLocalStatus(dataDir, port string) int {
	running, err := localGatewayRunning(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if !running {
		fmt.Fprintf(os.Stderr, "error: local gateway is not running\n")
		return 1
	}

	if dataDir == "" {
		fmt.Fprintf(os.Stderr, "error: --data-dir is required\n")
		return 1
	}

	absPath, err := resolvePath(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resolvedPort, err := resolveLocalGatewayPortFromBookkeeping(absPath, port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	token := resolveLocalGatewayToken(absPath)
	printLocalGatewayStatus(absPath, resolvedPort, token)
	return 0
}

func runLocalRestart(dataDir, port string) int {
	if dataDir == "" {
		fmt.Fprintf(os.Stderr, "error: --data-dir is required\n")
		return 1
	}

	absPath, err := resolvePath(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if err := validateLocalDataDir(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	token, err := resolveToken(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if !ensureOpenclawAvailable() {
		return 1
	}

	selectedPort, bumped, err := selectRestartPort(port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	running, err := localGatewayRunning(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if running {
		if err := openclawExecPreviewed(absPath, token, "gateway", "stop"); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return exitErr.ExitCode()
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		if err := removeGatewayBookkeeping(absPath); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Println("Stopped local gateway")
	}

	portStr := fmt.Sprintf("%d", selectedPort)
	gatewayArgs := []string{"gateway", "--bind", "lan", "--port", portStr}
	cmd, err := openclawExecDetached(absPath, token, gatewayArgs...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if err := writeGatewayBookkeeping(absPath, cmd.Process.Pid, selectedPort); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	printLocalLaunchHelp(absPath, portStr, token, bumped)
	return 0
}