# Scenario

**Feature**: `--stop` stops a running local gateway.

## Steps

1. Write `.my/gateway.json` with `os.Getpid()` and port `18789`; set busy port hook.
2. Run `my openclaw run --stop --data-dir <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Stop = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
	return markGatewayRunning(t, req)
}
```