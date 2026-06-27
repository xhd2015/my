# Scenario

**Feature**: re-adding an existing path warns on stderr and updates the note.

```
# path already in registry
my openclaw add data-dir <existing> --note "updated" -> stderr warning + note update, exit 0
```

## Preconditions

- Registry already contains the target path with note `"original"`.

## Steps

1. Seed registry with the data dir and note `"original"`.
2. Run `add` again with `--note "updated"`.

```go
import (
	"fmt"
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	seed := fmt.Sprintf(`{
  "data_dirs": [
    {
      "path": %q,
      "note": "original",
      "added_at": "2026-06-27T10:00:00+08:00"
    }
  ]
}`, req.DataDirPath)
	if err := os.WriteFile(filepath.Join(req.ConfigDir, "openclaw.json"), []byte(seed), 0o644); err != nil {
		return err
	}
	req.Note = "updated"
	return nil
}
```