# Scenario

**Feature**: `--stop` warns when container is not running.

## Steps

1. Set `req.Subcommand = "run-in-podman"` and `req.Stop = true`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Stop = true
	return nil
}
```