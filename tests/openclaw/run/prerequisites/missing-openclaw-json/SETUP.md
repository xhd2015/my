# Scenario

**Feature**: `my openclaw run` requires `<data-dir>/openclaw.json`.

```
my openclaw run --data-dir <dir-without-json> -> stderr mentions openclaw.json, exit 1
```

## Preconditions

- Data directory exists but lacks `openclaw.json`.

## Steps

1. Create empty data directory without config file.
2. Run command.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := filepath.Join(t.TempDir(), "no-json")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	req.RunDataDir = abs
	return nil
}
```
