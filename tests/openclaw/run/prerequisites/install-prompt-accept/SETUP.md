# Scenario

**Feature**: interactive install prompt accepts npm install when `openclaw` is missing on a TTY.

```
my openclaw run --data-dir <valid> -> prompt yes -> npm stub installs openclaw -> gateway launch
```

## Preconditions

- Valid data directory with gateway token.
- Stdin simulated as TTY (`MY_OPENCLAW_EXEC_INTERACTIVE=1`).
- Canned prompt answer `yes` (`MY_OPENCLAW_INSTALL_ANSWER=yes`).
- No `openclaw` binary initially; npm stub writes it on `install -g openclaw@latest`.

## Steps

1. Use `with-token` fixture.
2. Enable install prompt hooks and remove stub `openclaw` from test `binDir`.
3. Run command; npm stub should install openclaw, then gateway should launch.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.InstallPrompt = true
	req.InstallAnswer = "yes"
	req.BinDirOnlyPATH = true
	stub := filepath.Join(req.BinDir, "openclaw")
	if err := os.Remove(stub); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
```