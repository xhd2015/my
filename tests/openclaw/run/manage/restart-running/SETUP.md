# Scenario

**Feature**: `--restart` stops a running gateway and starts it again detached.

## Steps

1. Write `.my/gateway.json` with `os.Getpid()` and port `18789`; set busy port hook.
2. Run `my openclaw run --restart --data-dir <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Restart = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
	return markGatewayRunning(t, req)
}
```