# Scenario

**Feature**: `--restart` starts gateway when it is not running (requires `--data-dir`).

## Steps

1. Ensure gateway is not running.
2. Run `my openclaw run --restart --data-dir <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Restart = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
	return nil
}
```