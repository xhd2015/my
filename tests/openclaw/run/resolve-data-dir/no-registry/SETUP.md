# Scenario

**Feature**: `run --status` without `--data-dir` errors when the registry is empty.

## Steps

1. Do not seed registry (`RegistrySeed` unset).
2. Run `my openclaw run --status` with no `--data-dir`.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	req.RegistrySeed = nil
	err := os.Remove(filepath.Join(req.ConfigDir, "openclaw.json"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

```