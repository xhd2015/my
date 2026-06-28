## Expected

- Exit code `0`.
- `<data-dir>/.my/gateway.json` exists with `port` `18789` and a live `pid`.
- Stdout confirms gateway start.

## Exit Code

- `0`

```go
import (
	"os"
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Gateway started locally") {
		t.Fatalf("stdout missing launch confirmation:\n%s", resp.Stdout)
	}
	if _, err := os.Stat(gatewayBookkeepingPath(req.RunDataDir)); err != nil {
		t.Fatalf("gateway.json missing: %v", err)
	}
	rec, err := readGatewayBookkeeping(req.RunDataDir)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Port != 18789 {
		t.Fatalf("bookkeeping port = %d, want 18789", rec.Port)
	}
	if rec.PID <= 0 {
		t.Fatalf("bookkeeping pid = %d, want positive", rec.PID)
	}
	if rec.Kind != "local" {
		t.Fatalf("bookkeeping kind = %q, want local", rec.Kind)
	}
}
```