# Scenario

**Feature**: `--status` prints gateway URLs and info for a running container.

## Steps

1. Seed a running container with a known mounted data dir.
2. Run `my openclaw run-in-podman --status`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Status = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.PodmanContainerDataDir = req.RunDataDir
	return markContainerRunning(req, "openclaw-gateway")
}
```