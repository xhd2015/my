# Scenario

**Feature**: `--test-slack` sends a test message when the local gateway is running and Slack is configured.

## Steps

1. Write `.my/gateway.json` with `os.Getpid()` and port `18789`; set busy port hook and `MY_OPENCLAW_SLACK_STUB=1`.
2. Copy slack-enabled data-dir fixture from leaf `testdata/with-slack`.
3. Run `my openclaw run --test-slack --data-dir <fixture>`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.TestSlack = true
	req.SlackStub = true
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run", "manage", "test-slack-running", "testdata", "with-slack"))
	return markGatewayRunning(t, req)
}
```