# Scenario

**Feature**: default local launch execs gateway on port 18789 with state dir and dashboard URL.

```
my openclaw run --data-dir <valid> -> openclaw gateway --bind lan --port 18789 -> dashboard URL
```

## Preconditions

- Valid data dir with token in `openclaw.json`.
- Port `18789` is available.

## Steps

1. Use `with-token` fixture.
2. Run with default port selection.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	return nil
}
```
