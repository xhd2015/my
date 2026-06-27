# Scenario

**Feature**: gateway token falls back to `.env` when json has no token.

```
Token resolver <- .env (OPENCLAW_GATEWAY_TOKEN) -> podman run env
```

## Preconditions

- `openclaw.json` lacks `gateway.auth.token`.
- `.env` contains `OPENCLAW_GATEWAY_TOKEN=env-gateway-token`.

## Steps

1. Use `with-env-only` fixture.
2. Run command.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-env-only")
	req.PodmanImageExists = true
	return nil
}
```