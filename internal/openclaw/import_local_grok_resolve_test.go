package openclaw

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveContainerDataDirRequiresExplicitPath(t *testing.T) {
	const missingContainer = "openclaw-gateway-doctest-missing"

	_, err := resolveContainerDataDir("", missingContainer)
	if err == nil {
		t.Fatal("expected error when container is not running and --data-dir is omitted")
	}
}

func TestPrintImportLocalGrokPlan(t *testing.T) {
	dataDir := filepath.Join(t.TempDir(), ".openclaw")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Smoke test: plan printing should not panic.
	printImportLocalGrokPlan(dataDir, "/tmp/grok-auth.json", "openclaw-gateway-doctest-missing")
}