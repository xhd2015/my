# Scenario

**Feature**: auto-port mode retries when openclaw fails at runtime with port in use.

```
port 18789 free at probe but fails at launch -> retry 18790 -> gateway succeeds
```

## Preconditions

- Valid data dir with token.
- Port `18789` appears available (no busy-port stub).
- `OPENCLAW_STUB_FAIL_PORTS=18789` makes the first openclaw launch fail with `EADDRINUSE`.

## Steps

1. Use `with-token` fixture.
2. Set `req.FailPorts = "18789"`.
3. Run without explicit `--port`.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.FailPorts = "18789"
	return nil
}
```