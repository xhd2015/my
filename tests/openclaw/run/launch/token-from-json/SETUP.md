# Scenario

**Feature**: gateway token is taken from `openclaw.json` when present.

```
Token resolver <- openclaw.json (gateway.auth.token) -> OPENCLAW_GATEWAY_TOKEN in openclaw exec
```

## Preconditions

- `openclaw.json` contains `gateway.auth.token`.

## Steps

1. Use `with-token` fixture.
2. Run command.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	return nil
}
```
