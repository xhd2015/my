# Scenario

**Feature**: interactive install prompt declines when `openclaw` is missing on a TTY.

```
my openclaw run --data-dir <valid> -> openclaw not found + prompt -> user declines -> exit 1
```

## Preconditions

- Valid data directory with gateway token.
- Stdin simulated as TTY (`MY_OPENCLAW_EXEC_INTERACTIVE=1`).
- Canned prompt answer `no` (`MY_OPENCLAW_INSTALL_ANSWER=no`).
- No `openclaw` binary on `PATH` (stub removed after root setup).

## Steps

1. Use `with-token` fixture.
2. Enable install prompt hooks and remove stub `openclaw` from test `binDir`.
3. Run command.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.InstallPrompt = true
	req.InstallAnswer = "no"
	req.BinDirOnlyPATH = true
	stub := filepath.Join(req.BinDir, "openclaw")
	if err := os.Remove(stub); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
```