# Scenario

**Feature**: default port scan bumps past busy ports starting at 18789.

```
ports 18789-18790 busy -> my openclaw run selects 18791 and notes bump on stdout
```

## Preconditions

- Valid data dir with token.
- `OPENCLAW_STUB_BUSY_PORTS=18789,18790`.

## Steps

1. Use `with-token` fixture.
2. Set `req.BusyPorts = "18789,18790"`.
3. Run without explicit `--port`.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.BusyPorts = "18789,18790"
	return nil
}
```
