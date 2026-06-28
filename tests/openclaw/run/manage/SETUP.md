# Scenario

**Feature**: `my openclaw run` management actions (`--status`, `--restart`, `--import-local-grok`, `--test-slack`).

```
my openclaw run --<action> -> openclaw stub gateway state -> stdout/stderr
```

## Preconditions

- Running simulation via `writeGatewayBookkeeping(t, dataDir, os.Getpid(), 18789)` and
  `OPENCLAW_STUB_BUSY_PORTS=18789`.
- OpenClaw stub still handles `gateway stop` / detached `gateway` start for `--restart`.
- Slack tests may set `MY_OPENCLAW_SLACK_STUB=1` to skip real Slack HTTP.

```go
import "os"

func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run"
	return nil
}

func markGatewayRunning(t *testing.T, req *Request) error {
	req.BusyPorts = "18789"
	return writeGatewayBookkeeping(t, req.RunDataDir, os.Getpid(), 18789)
}
```