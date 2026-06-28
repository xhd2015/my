package openclaw

import (
	"fmt"
	"os"
	"runtime"
)

func runRestart(explicitDataDir string, rebuild bool, containerName, port string) int {
	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resolvedDataDir, err := resolveContainerDataDir(explicitDataDir, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	running, err := containerIsRunning(containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if running {
		if err := stopAndRemoveContainer(containerName); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Printf("Stopped and removed container %s\n", containerName)
	}

	return runLaunch(resolvedDataDir, rebuild, containerName, port)
}

func runStop(containerName string) int {
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
		fmt.Fprintf(os.Stderr, "warning: container %s is not running\n", containerName)
		return 0
	}

	if err := stopAndRemoveContainer(containerName); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	fmt.Printf("Stopped and removed container %s\n", containerName)
	return 0
}

func runStatus(explicitDataDir, port, containerName string) int {
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
		fmt.Fprintf(os.Stderr, "error: container %s is not running\n", containerName)
		return 1
	}

	dataDir, err := resolveContainerDataDir(explicitDataDir, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resolvedPort := resolveGatewayPort(port, containerName)
	token := resolveGatewayToken(dataDir, containerName)
	printGatewayStatus(containerName, dataDir, resolvedPort, token)
	return 0
}

func runLogs(containerName string) int {
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
		fmt.Fprintf(os.Stderr, "error: container %s is not running\n", containerName)
		return 1
	}

	if err := showContainerLogs(containerName); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}

func runShowTokens(dataDir, port, containerName string) int {
	absPath, err := resolvePath(dataDir)
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

	sources, err := resolveTokenSources(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	printCommand(fmt.Sprintf("read gateway tokens from %s", shellQuote(absPath)))
	printTokenSources(absPath, port, containerName, sources)
	return 0
}

func resolveDashboardToken(dataDir, containerName string) (string, error) {
	if dataDir != "" {
		absPath, err := resolvePath(dataDir)
		if err != nil {
			return "", err
		}
		sources, err := resolveTokenSources(absPath)
		if err != nil {
			return "", err
		}
		if token := sources.Effective(); token != "" {
			return token, nil
		}
	}

	running, err := containerIsRunning(containerName)
	if err != nil {
		return "", err
	}
	if running {
		token, err := containerEnvToken(containerName)
		if err != nil {
			return "", err
		}
		if token != "" {
			return token, nil
		}
	}

	return "", fmt.Errorf("gateway token not found; pass --data-dir or start the container first")
}

func runDashboard(dataDir, port, containerName string) int {
	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	token, err := resolveDashboardToken(dataDir, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	url := fmt.Sprintf("http://127.0.0.1:%s/#token=%s", port, token)
	openCmd := "open " + shellQuote(url)
	if runtime.GOOS == "linux" {
		openCmd = "xdg-open " + shellQuote(url)
	}
	printCommand(openCmd)

	if err := openURL(url); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}

	fmt.Println(url)
	return 0
}