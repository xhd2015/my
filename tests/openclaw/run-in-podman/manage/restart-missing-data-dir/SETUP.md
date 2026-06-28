# Scenario

**Feature**: `--restart` requires `--data-dir` when no gateway container is running.

## Steps

1. Ensure container is not running.
2. Run `my openclaw run-in-podman --restart` without `--data-dir`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Restart = true
	return nil
}
```