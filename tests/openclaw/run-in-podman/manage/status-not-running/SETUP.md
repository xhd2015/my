# Scenario

**Feature**: `--status` errors when the gateway container is not running.

## Steps

1. Ensure container is not running.
2. Run `my openclaw run-in-podman --status`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Status = true
	return nil
}
```