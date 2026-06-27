# Scenario

**Feature**: darwin auto-starts Podman machine when `podman machine info` reports stopped.

```
Podman stub (machine info stopped) -> my openclaw run-in-podman -> podman machine start -> launch
```

## Preconditions

- `runtime.GOOS == "darwin"`.
- `PodmanMachineRunning = false`.
- Valid data dir and existing image.

## Steps

1. Run `my openclaw run-in-podman` with stopped machine simulation.

```go
func Setup(t *testing.T, req *Request) error {
	if req.PodmanMachineRunning {
		t.Fatal("precondition: podman machine must be simulated as stopped")
	}
	return nil
}
```