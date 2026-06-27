# Scenario

**Feature**: `my openclaw run-in-podman` launches OpenClaw gateway in a Podman container.

```
# validate data dir, resolve token, rebuild image, run container
my openclaw run-in-podman --data-dir <path> -> Podman stub -> gateway container
```

## Preconditions

- Stub `podman` on PATH records all invocations.
- Podman machine is running by default (`PodmanMachineRunning=true`).

## Steps

1. Set `req.Subcommand` to `"run-in-podman"`.
2. Descendants configure data dir fixtures, flags, and registry image state.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.ContainerName = ""
	req.Port = ""
	req.Rebuild = false
	return nil
}

func fixtureDataDir(t *testing.T, name string) string {
	t.Helper()
	src := filepath.Join(DOCTEST_ROOT, "run-in-podman", "testdata", name)
	return copyFixtureDir(t, src)
}
```