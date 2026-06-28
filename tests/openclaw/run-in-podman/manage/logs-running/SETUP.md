# Scenario

**Feature**: `--logs` shows container logs when running.

## Steps

1. Seed running container.
2. Run `my openclaw run-in-podman --logs`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Logs = true
	return markContainerRunning(req, "openclaw-gateway")
}
```