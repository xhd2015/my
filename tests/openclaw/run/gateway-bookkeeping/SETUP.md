# Scenario

**Feature**: per-data-dir `.my/gateway.json` bookkeeping for local gateway running detection.

```
<data-dir>/.my/gateway.json  ->  pid alive + port in use  ->  running
```

## Preconditions

- Running simulation: `writeGatewayBookkeeping(t, dataDir, os.Getpid(), port)` plus
  `OPENCLAW_STUB_BUSY_PORTS` including the bookkeeping port (default `18789`).
- Stale bookkeeping (dead pid or free port) is removed lazily on read.

## Steps

1. Set `req.Subcommand` to `"run"`.
2. Leaf configures bookkeeping fixtures and management flags (`--status`, `--restart`).

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run"
	return nil
}
```