# Scenario

**Feature**: `my openclaw run` launches the OpenClaw gateway on the host using a custom data directory.

```
# validate data dir, resolve token, select port, exec openclaw gateway
my openclaw run --data-dir <path> [--port PORT] -> openclaw stub -> foreground gateway
```

## Preconditions

- Stub `openclaw` on `PATH` records invocations to `OPENCLAW_STUB_LOG`.
- Port availability can be simulated via `OPENCLAW_STUB_BUSY_PORTS` (test-only hook).

## Steps

1. Set `req.Subcommand` to `"run"`.
2. Descendants configure data dir fixtures, port flags, and busy-port simulation.

```go
import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run"
	req.Port = ""
	req.BusyPorts = ""
	return nil
}

func fixtureDataDir(t *testing.T, name string) string {
	t.Helper()
	src := filepath.Join(DOCTEST_ROOT, "run-in-podman", "testdata", name)
	return copyFixtureDir(t, src)
}

type gatewayBookkeepingFile struct {
	PID       int    `json:"pid"`
	Port      int    `json:"port"`
	StartedAt string `json:"started_at"`
	Kind      string `json:"kind"`
}

func writeGatewayBookkeeping(t *testing.T, dataDir string, pid, port int) error {
	t.Helper()
	myDir := filepath.Join(dataDir, ".my")
	if err := os.MkdirAll(myDir, 0o755); err != nil {
		return err
	}
	body, err := json.Marshal(gatewayBookkeepingFile{
		PID:       pid,
		Port:      port,
		StartedAt: "2026-06-28T13:00:00+08:00",
		Kind:      "local",
	})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(myDir, "gateway.json"), body, 0o644)
}

func gatewayBookkeepingPath(dataDir string) string {
	return filepath.Join(dataDir, ".my", "gateway.json")
}

func readGatewayBookkeeping(dataDir string) (*gatewayBookkeepingFile, error) {
	data, err := os.ReadFile(gatewayBookkeepingPath(dataDir))
	if err != nil {
		return nil, err
	}
	var rec gatewayBookkeepingFile
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}
```
