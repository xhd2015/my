## Expected

- Exit code `0`.
- Stub receives `gateway --bind lan --port 18791`.
- Stdout notes port bump from `18789` and uses dashboard URL on `18791`.

## Side Effects

- Auto-selected port skips busy defaults.

## Errors

- None.

## Exit Code

- `0`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "gateway --bind lan --port 18791") {
		t.Fatalf("openclaw should use port 18791:\n%v", resp.OpenClawCalls)
	}
	if !strings.Contains(resp.Stdout, "18789") || !strings.Contains(resp.Stdout, "18791") {
		t.Fatalf("stdout should note port bump:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "http://127.0.0.1:18791/") {
		t.Fatalf("stdout missing bumped dashboard URL:\n%s", resp.Stdout)
	}
}
```