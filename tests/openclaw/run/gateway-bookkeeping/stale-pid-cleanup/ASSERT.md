## Expected

- Exit code `1` (gateway not running).
- Stale `<data-dir>/.my/gateway.json` removed during running check.

## Exit Code

- `1`

```go
import "os"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if _, err := os.Stat(gatewayBookkeepingPath(req.RunDataDir)); !os.IsNotExist(err) {
		t.Fatalf("stale gateway.json should be removed, stat err: %v", err)
	}
}
```