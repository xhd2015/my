package openclaw

import (
	"fmt"
	"path/filepath"
	"strings"
)

const containerOpenClawMount = "/home/node/.openclaw"

func containerOpenClawDataDir(name string) (string, error) {
	format := `{{range .Mounts}}{{if eq .Destination "` + containerOpenClawMount + `"}}{{.Source}}{{end}}{{end}}`
	out, err := podmanOutput("inspect", name, "--format", format)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func resolveContainerDataDir(explicitDataDir, containerName string) (string, error) {
	if explicitDataDir != "" {
		return explicitDataDir, nil
	}
	running, err := containerIsRunning(containerName)
	if err != nil {
		return "", err
	}
	if running {
		dataDir, err := containerOpenClawDataDir(containerName)
		if err != nil {
			return "", err
		}
		if dataDir == "" {
			return "", fmt.Errorf("could not determine data dir from running container %s", containerName)
		}
		return dataDir, nil
	}
	return "", fmt.Errorf("--data-dir is required when container %s is not running", containerName)
}

func requireRunningContainerDataDir(explicitDataDir, containerName string) (string, error) {
	running, err := containerIsRunning(containerName)
	if err != nil {
		return "", err
	}
	if !running {
		return "", fmt.Errorf("container %s is not running (start it with 'my openclaw run-in-podman --data-dir <path>' first)", containerName)
	}
	return resolveContainerDataDir(explicitDataDir, containerName)
}

func resolveImportLocalGrokDataDir(explicitDataDir, containerName string) (string, error) {
	return resolveContainerDataDir(explicitDataDir, containerName)
}

func printImportLocalGrokPlan(dataDir, grokAuthPath, containerName string) {
	fmt.Printf("Reading Grok OAuth from: %s\n", grokAuthPath)
	fmt.Printf("Writing to data dir: %s\n", dataDir)
	fmt.Println("Will update:")
	fmt.Printf("  - %s — add or refresh the xai:default OAuth profile\n", agentAuthDBPath(dataDir))
	fmt.Printf("  - %s — set auth.profiles/order for xai and agents.defaults.model to %s\n",
		filepath.Join(dataDir, "openclaw.json"), xaiDefaultModel)
	fmt.Printf("  - %s — add or refresh the xai provider with %s for the dashboard model picker\n",
		agentModelsJSONPath(dataDir), xaiDefaultModel)
	fmt.Printf("  - %s — reset %s and dashboard session models to %s\n",
		mainSessionsJSONPath(dataDir), mainAgentSessionKey, xaiDefaultModel)
	if containerName != "" {
		if running, err := containerIsRunning(containerName); err == nil && running {
			fmt.Printf("  - %s — copy Grok CLI OAuth credentials for `grok` inside the container\n", containerGrokAuthPath)
		}
	}
}