# Scenario

**Feature**: `my openclaw run` errors when `openclaw` is not on `PATH`.

```
my openclaw run --data-dir <valid> -> openclaw not found + install instructions, exit 1
```

## Preconditions

- Valid data directory with gateway token.
- No `openclaw` binary on `PATH` (stub removed after root setup).

## Steps

1. Use `with-token` fixture.
2. Remove stub `openclaw` from test `binDir`.
3. Run command.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.BinDirOnlyPATH = true
	stub := filepath.Join(req.BinDir, "openclaw")
	if err := os.Remove(stub); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
```
