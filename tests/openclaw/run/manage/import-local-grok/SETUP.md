# Scenario

**Feature**: `--import-local-grok` copies local Grok Build OAuth into the OpenClaw data dir (host only).

## Steps

1. Copy a valid data-dir fixture with `openclaw.json` (reuse run-in-podman manage fixture).
2. Point `MY_GROK_AUTH_PATH` at leaf `testdata/grok-auth.json`.
3. Run `my openclaw run --import-local-grok --data-dir <fixture>`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.ImportLocalGrok = true
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run-in-podman", "manage", "import-local-grok", "testdata", "data-dir"))
	req.GrokAuthPath = filepath.Join(DOCTEST_ROOT, "run-in-podman", "manage", "import-local-grok", "testdata", "grok-auth.json")
	return nil
}
```