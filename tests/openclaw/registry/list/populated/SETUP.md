# Scenario

**Feature**: list prints header and tab-separated rows for registered dirs.

```
my openclaw list <- populated Config store -> header + path/note/added_at rows
```

## Preconditions

- Registry contains at least one data dir entry.

## Steps

1. Seed registry with one entry before running list.

```go
import (
	"fmt"
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := filepath.Join(t.TempDir(), "registered")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	seed := fmt.Sprintf(`{
  "data_dirs": [
    {
      "path": %q,
      "note": "work machine",
      "added_at": "2026-06-27T10:00:00+08:00"
    }
  ]
}`, abs)
	return os.WriteFile(filepath.Join(req.ConfigDir, "openclaw.json"), []byte(seed), 0o644)
}
```