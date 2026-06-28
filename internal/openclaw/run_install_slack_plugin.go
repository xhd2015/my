package openclaw

import (
	"fmt"
	"os"
	"os/exec"
)

const slackPluginPackage = "@openclaw/slack"

func runInstallSlackPlugin(dataDir, containerName string) int {
	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resolvedDataDir, err := requireRunningContainerDataDir(dataDir, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	absPath, err := resolvePath(resolvedDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Installing Slack plugin into %s\n", absPath)
	installCmd := "OPENCLAW_STATE_DIR=" + shellQuote(absPath) + " npx --yes openclaw@latest plugins install " + slackPluginPackage
	printCommand(installCmd)

	cmd := exec.Command("npx", "--yes", "openclaw@latest", "plugins", "install", slackPluginPackage)
	cmd.Env = append(os.Environ(), "OPENCLAW_STATE_DIR="+absPath)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Println()
	fmt.Println("Slack plugin installed into the running session data dir.")
	fmt.Println("Restart the gateway so the container loads it:")
	fmt.Printf("  my openclaw run-in-podman --restart --container-name %s\n", containerName)
	return 0
}