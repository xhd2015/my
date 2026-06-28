# Scenario

**Feature**: `my openclaw run` auto-resolves `--data-dir` from the registry when omitted.

```
registry entries + openclaw gateway status per path -> auto-select or error
my openclaw run --status [--data-dir <path>]
```

## Preconditions

- Registry seeded via `RegistrySeed` or `seedRegistry` helper.
- Per-dir running simulation via `writeGatewayBookkeeping(t, dataDir, os.Getpid(), 18789)` and
  `OPENCLAW_STUB_BUSY_PORTS=18789` (pid alive + port in use).
- Registered data dirs are absolute temp paths with minimal `openclaw.json` fixtures.

## Steps

1. Set `req.Subcommand` to `"run"` and `req.Status` to `true` (leaf may override).
2. Seed registry and write bookkeeping for dirs that should appear running.
3. Omit `req.RunDataDir` to exercise auto-resolve (except `explicit-override`).

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const defaultGatewayPort = 18789

func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run"
	req.Status = true
	return nil
}

func minimalDataDir(t *testing.T, name string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `{
  "gateway": {
    "auth": {
      "token": "json-gateway-token"
    }
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "openclaw.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

func seedRegistry(req *Request, paths ...string) error {
	entries := make([]map[string]string, len(paths))
	for i, p := range paths {
		entries[i] = map[string]string{
			"path":     p,
			"note":     fmt.Sprintf("dir-%d", i+1),
			"added_at": fmt.Sprintf("2026-06-28T10:0%d:00+08:00", i),
		}
	}
	body, err := json.Marshal(map[string]any{"data_dirs": entries})
	if err != nil {
		return err
	}
	req.RegistrySeed = body
	return os.WriteFile(filepath.Join(req.ConfigDir, "openclaw.json"), body, 0o644)
}

func markDirsRunning(t *testing.T, req *Request, paths ...string) error {
	t.Helper()
	req.BusyPorts = "18789"
	for _, p := range paths {
		if err := writeGatewayBookkeeping(t, p, os.Getpid(), defaultGatewayPort); err != nil {
			return err
		}
	}
	return nil
}

func registeredPaths(req *Request) ([]string, error) {
	reg, err := parseRegistry(req.RegistrySeed)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(reg.DataDirs))
	for _, entry := range reg.DataDirs {
		paths = append(paths, entry.Path)
	}
	return paths, nil
}

```