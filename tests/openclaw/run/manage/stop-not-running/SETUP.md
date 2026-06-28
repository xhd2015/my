# Scenario

**Feature**: `--stop` warns when the local gateway is not running.

## Steps

1. Seed data dir with `openclaw.json` fixture (no bookkeeping).
2. Run `my openclaw run --stop --data-dir <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Stop = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
	return nil
}
```