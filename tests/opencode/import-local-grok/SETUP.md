# Scenario

**Feature**: `my opencode --import-local-grok` writes merged OpenCode `auth.json` from local Grok OIDC file.

```
# build my + isolated MY_OPENCODE_DATA_DIR + MY_GROK_AUTH_PATH fixture
my opencode --import-local-grok -> auth.json with xai oauth
```

## Preconditions

- Module root is `filepath.Join(DOCTEST_ROOT, "..", "..")`.
- OpenCode data dir is isolated per test via `MY_OPENCODE_DATA_DIR` (never the real user data dir).
- Grok auth path is set via `MY_GROK_AUTH_PATH` to leaf fixtures under `testdata/`.
- The `my` CLI is built to `binPath` before each test run.

## Steps

1. Create temp dirs for OpenCode data and grok auth fixtures.
2. Build `my` from `./cmd/my/`.
3. Leaf `Setup` sets `Request` fields (grok path, preseed auth, second import).
4. Invoke shared `Run` from `DOCTEST.md`.

## Context

- Destination file: `<MY_OPENCODE_DATA_DIR>/auth.json`.
- Tests do not call live xAI APIs.

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.CommandDir = filepath.Clean(filepath.Join(DOCTEST_ROOT, "..", ".."))
	tmp := t.TempDir()
	req.OpenCodeDataDir = filepath.Join(tmp, "opencode-data")
	req.BinPath = filepath.Join(tmp, "my")

	build := exec.Command("go", "build", "-o", req.BinPath, "./cmd/my/")
	build.Dir = req.CommandDir
	if output, err := build.CombinedOutput(); err != nil {
		return fmt.Errorf("go build my: %w: %s", err, strings.TrimSpace(string(output)))
	}
	if err := os.MkdirAll(req.OpenCodeDataDir, 0o755); err != nil {
		return err
	}
	return nil
}```
