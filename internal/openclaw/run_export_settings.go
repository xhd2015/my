package openclaw

import (
	"fmt"
	"os"
	"path/filepath"
)

func runExportSettings(dataDir, containerName, outputZip string) int {
	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resolvedDataDir, err := requireRunningContainerDataDir(dataDir, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	absDataDir, err := resolvePath(resolvedDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	absOutput, err := resolvePath(outputZip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Exporting settings from %s\n", absDataDir)
	fmt.Printf("Writing archive: %s\n", absOutput)
	fmt.Fprintf(os.Stderr, "Preparing SQLite databases for export...\n")

	if _, err := prepareSQLiteFilesForExport(absDataDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fileCount, err := zipOpenClawDataDir(absDataDir, absOutput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Exported %d files under %s/\n", fileCount, exportSettingsZipRoot)
	fmt.Println()
	fmt.Println("Restore on another machine:")
	unzipParent := filepath.Dir(absOutput)
	restoredDataDir := filepath.Join(unzipParent, exportSettingsZipRoot)
	fmt.Printf("  unzip %s -d %s\n", shellQuote(absOutput), shellQuote(unzipParent))
	fmt.Printf("  my openclaw run --data-dir %s\n", shellQuote(restoredDataDir))
	fmt.Printf("  my openclaw run-in-podman --data-dir %s\n", shellQuote(restoredDataDir))
	fmt.Println()
	fmt.Println("The archive contains secrets (tokens, credentials). Store and share it carefully.")
	return 0
}