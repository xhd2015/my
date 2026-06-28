# Scenario

**Feature**: `--show-tokens` requires `--data-dir` when no gateway container is running.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.ShowTokens = true
	return nil
}
```