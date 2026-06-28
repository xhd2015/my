# Scenario

**Feature**: `--port` uses the requested port when available.

```
my openclaw run --port 19999 -> openclaw gateway --port 19999 and dashboard URL on 19999
```

## Preconditions

- Valid data dir with token.
- Port `19999` is available.

## Steps

1. Use `with-token` fixture.
2. Set `req.Port = "19999"`.
3. Run command.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.Port = "19999"
	return nil
}
```
