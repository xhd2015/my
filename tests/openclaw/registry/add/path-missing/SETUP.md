# Scenario

**Feature**: add rejects paths that do not exist.

```
# missing directory on disk
my openclaw add data-dir <missing-path> -> error on stderr, exit 1
```

## Preconditions

- Target path does not exist.

## Steps

1. Set `req.DataDirPath` to a non-existent absolute path.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.DataDirPath = filepath.Join(t.TempDir(), "does-not-exist", ".openclaw")
	return nil
}
```