# Scenario

**Feature**: `my openclaw run` rejects a non-existent data directory.

```
my openclaw run --data-dir <missing> -> validation error, exit 1
```

## Preconditions

- `--data-dir` points to a path that does not exist.

## Steps

1. Set `req.RunDataDir` to a non-existent absolute path.
2. Run command.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = filepath.Join(t.TempDir(), "missing-data-dir")
	return nil
}
```
