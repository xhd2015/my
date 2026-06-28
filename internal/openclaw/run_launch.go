package openclaw

import (
	"fmt"
	"os"
	"path/filepath"
)

func runLaunch(dataDir string, rebuild bool, containerName, port string) int {
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

	configPath := filepath.Join(absPath, "openclaw.json")
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: openclaw.json not found in %s\n", absPath)
			return 1
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	token, err := resolveToken(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	reg, err := loadRegistry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	shouldRebuild, err := needsRebuild(reg, rebuild)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if shouldRebuild {
		if err := buildImage(reg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	}

	workspaceDir := filepath.Join(absPath, "workspace")
	portMapping := fmt.Sprintf("%s:%s", port, containerPort)

	runArgs := []string{
		"run", "-d", "--replace",
		"--name", containerName,
		"-v", absPath + ":/home/node/.openclaw",
		"-v", workspaceDir + ":/home/node/.openclaw/workspace",
		"-e", "OPENCLAW_GATEWAY_TOKEN=" + token,
		"-p", portMapping,
		imageName,
		"gateway", "--bind", "lan",
	}
	if err := podmanRunPreviewed(runArgs...); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	printLaunchHelp(port, containerName, absPath, token)
	return 0
}