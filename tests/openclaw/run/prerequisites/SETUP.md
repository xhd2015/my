# Scenario

**Feature**: `my openclaw run` validates required inputs before launching the gateway.

```
# missing data dir, openclaw.json, token, or openclaw binary
my openclaw run --data-dir <path> -> validation error, exit 1
# TTY + missing openclaw -> install prompt (accept runs npm stub, decline exits 1)
```

## Preconditions

- Tests exercise invalid or incomplete data directory setups, a missing `openclaw` binary, or the
  interactive install prompt (`MY_OPENCLAW_EXEC_INTERACTIVE=1`, `MY_OPENCLAW_INSTALL_ANSWER`).

## Steps

1. Configure `req.RunDataDir`, fixtures, and PATH per leaf.
2. No successful `openclaw gateway` invocation should occur.

```go
func Setup(t *testing.T, req *Request) error {
	req.BusyPorts = ""
	req.Port = ""
	return nil
}
```
