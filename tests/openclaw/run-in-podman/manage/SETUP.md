# Scenario

**Feature**: `run-in-podman` management actions (`--stop`, `--restart`, `--logs`, `--status`, `--show-tokens`, `--exec`, `--import-local-grok`, `--install-grok`) and already-running guard.

```
my openclaw run-in-podman --<action> -> Podman stub state -> stdout/stderr
```

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	return nil
}
```