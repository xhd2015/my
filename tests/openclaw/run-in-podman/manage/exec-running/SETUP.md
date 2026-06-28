# Scenario

**Feature**: `--exec` runs a command verbatim inside the running gateway container.

## Steps

1. Seed running container.
2. Run `my openclaw run-in-podman --exec openclaw models auth login --provider xai --method oauth`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Exec = true
	req.ExecArgs = []string{
		"openclaw", "models", "auth", "login",
		"--provider", "xai", "--method", "oauth",
	}
	return markContainerRunning(req, "openclaw-gateway")
}
```