# Scenario

**Feature**: macOS auto-starts Podman machine when stopped.

```
# machine info shows not running
my openclaw run-in-podman -> Podman stub (machine start) -> container launch
```

## Preconditions

- `runtime.GOOS == "darwin"` for machine check behavior.
- Valid data directory fixture.

## Steps

1. Set `req.PodmanMachineRunning = false`.
2. Copy valid data dir fixture and run command.

```go
import "runtime"

func Setup(t *testing.T, req *Request) error {
	if runtime.GOOS != "darwin" {
		t.Skip("podman machine auto-start is darwin-only")
	}
	req.PodmanMachineRunning = false
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.PodmanImageExists = true
	return nil
}
```