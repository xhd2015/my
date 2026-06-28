# Scenario

**Feature**: `--install-grok` installs the Grok CLI inside the running gateway container.

## Steps

1. Seed running container.
2. Run `my openclaw run-in-podman --install-grok`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.InstallGrok = true
	return markContainerRunning(req, "openclaw-gateway")
}
```