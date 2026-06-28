# Scenario

**Feature**: detached `--restart` writes `.my/gateway.json` with the selected gateway port.

## Steps

1. Use token fixture data dir (gateway not running).
2. Run `my openclaw run --restart --data-dir <fixture>`.
3. Assert bookkeeping file created with gateway port.

```go
func Setup(t *testing.T, req *Request) error {
	req.Restart = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.BusyPorts = "18789"
	return nil
}
```