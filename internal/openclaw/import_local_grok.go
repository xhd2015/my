package openclaw

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func runImportLocalGrok(dataDir, containerName string) int {
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

	grokAuthPath, err := resolveGrokAuthPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	printImportLocalGrokPlan(absPath, grokAuthPath, containerName)

	entry, err := loadGrokAuthEntry(grokAuthPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	store, err := grokEntryToAuthStore(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if time.Now().After(time.UnixMilli(store.Profiles[xaiDefaultProfile].Expires)) {
		fmt.Fprintf(os.Stderr, "warning: grok access token is expired; OpenClaw will try to refresh it on first use\n")
	}

	if err := writeAuthProfileStore(absPath, store); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
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

	if err := patchOpenClawConfigForXAI(configPath, entry.Email); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if err := syncAgentModelsJSONForXAI(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if err := resetMainSessionModelForXAI(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if containerName != "" {
		if err := copyGrokAuthToContainer(containerName, grokAuthPath); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	}

	fmt.Println("Import complete.")
	if containerName != "" {
		if running, err := containerIsRunning(containerName); err == nil && running {
			fmt.Println("Restart the gateway so it reloads openclaw.json and session defaults:")
			fmt.Printf("  my openclaw run-in-podman --restart --data-dir %s\n", absPath)
		}
	} else {
		fmt.Println("Restart the gateway so it reloads openclaw.json and session defaults:")
		fmt.Printf("  my openclaw run --restart --data-dir %s\n", absPath)
	}
	return 0
}