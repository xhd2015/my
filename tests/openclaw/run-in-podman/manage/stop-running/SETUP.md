# Scenario

**Feature**: `--stop` stops a running container.

## Steps

1. Seed running container in stub state.
2. Run `my openclaw run-in-podman --stop`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Stop = true
	return markContainerRunning(req, "openclaw-gateway")
}
```