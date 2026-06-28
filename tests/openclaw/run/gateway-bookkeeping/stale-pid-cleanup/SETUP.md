# Scenario

**Feature**: stale bookkeeping with dead pid is removed on read and gateway reports not running.

## Steps

1. Seed data dir with `openclaw.json` fixture.
2. Write `.my/gateway.json` with pid `999999` and port `18789`.
3. Set `OPENCLAW_STUB_BUSY_PORTS=18789` so the port check would pass if pid were alive.
4. Run `my openclaw run --status --data-dir <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Status = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.BusyPorts = "18789"
	return writeGatewayBookkeeping(t, req.RunDataDir, 999999, 18789)
}
```