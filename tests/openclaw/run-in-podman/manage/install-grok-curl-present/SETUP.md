# Scenario

**Feature**: `--install-grok` skips apt bootstrap when curl is already in the container.

## Steps

1. Seed running container with curl present.
2. Run `my openclaw run-in-podman --install-grok`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.InstallGrok = true
	t.Setenv("PODMAN_STUB_HAS_CURL", "1")
	return markContainerRunning(req, "openclaw-gateway")
}
```