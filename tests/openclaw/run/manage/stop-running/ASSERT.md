## Expected

- Exit code `0`.
- `gateway stop` invoked; bookkeeping removed.
- Stdout confirms stop.

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
	if !strings.Contains(resp.Stdout, "Stopped local gateway") {
		t.Fatalf("stdout missing stop confirmation:\n%s", resp.Stdout)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "gateway stop") {
		t.Fatalf("missing gateway stop:\n%v", resp.OpenClawCalls)
	}
	if _, err := os.Stat(gatewayBookkeepingPath(req.RunDataDir)); !os.IsNotExist(err) {
		t.Fatalf("gateway.json should be removed after stop, stat err: %v", err)
	}
}
```