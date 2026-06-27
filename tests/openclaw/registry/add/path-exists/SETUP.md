# Scenario

**Feature**: add succeeds when the target path exists and is a directory.

```
# path resolves to absolute existing directory
my openclaw add data-dir <existing-dir> -> registry entry created or updated
```

## Preconditions

- Target data directory exists on disk before `add` runs.

## Steps

1. Create or copy a fixture data directory for `req.DataDirPath`.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := filepath.Join(t.TempDir(), "data-dir")
	if err := os.MkdirAll(filepath.Join(dir, "workspace"), 0o755); err != nil {
		return err
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	req.DataDirPath = abs
	return nil
}
```