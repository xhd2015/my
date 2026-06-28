# Scenario

**Feature**: explicit `--port` errors when the requested port is busy.

```
--port 18789 with 18789 busy -> exit 1, stderr mentions port in use
```

## Preconditions

- Valid data dir with token.
- Port `18789` marked busy via `OPENCLAW_STUB_BUSY_PORTS`.

## Steps

1. Use `with-token` fixture.
2. Set `req.Port = "18789"` and `req.BusyPorts = "18789"`.
3. Run command.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.Port = "18789"
	req.BusyPorts = "18789"
	return nil
}
```
