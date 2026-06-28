package openclaw

import (
	"strings"
	"testing"
)

func TestGrokCLIInstallBootstrapToleratesClockSkew(t *testing.T) {
	for _, want := range []string{
		"Acquire::Check-Valid-Until=false",
		"Acquire::Check-Date=false",
		grokCLIInstallURL,
	} {
		if !strings.Contains(grokCLIInstallBootstrap, want) {
			t.Fatalf("bootstrap command missing %q: %s", want, grokCLIInstallBootstrap)
		}
	}
}