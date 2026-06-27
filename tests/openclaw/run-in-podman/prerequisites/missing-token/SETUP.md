# Scenario

**Feature**: run-in-podman errors when gateway token is absent from json and .env.

```
my openclaw run-in-podman --data-dir <no-token-dir> -> token error, exit 1
```

## Preconditions

- `openclaw.json` exists without `gateway.auth.token`.
- No `.env` with `OPENCLAW_GATEWAY_TOKEN`.

## Steps

1. Copy `no-token` fixture to temp data dir.
2. Run command.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "no-token")
	return nil
}
```