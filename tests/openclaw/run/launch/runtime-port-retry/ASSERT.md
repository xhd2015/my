## Expected

- Exit code `0`.
- First openclaw attempt uses port `18789`; relaunch uses `18790`.
- Stdout notes port bump from `18789` and uses dashboard URL on `18790`.

## Side Effects

- Runtime port retry after launch-time `EADDRINUSE` on default port.

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
	if !openclawCallsContain(resp.OpenClawCalls, "gateway --bind lan --port 18789") {
		t.Fatalf("openclaw should attempt default port 18789 first:\n%v", resp.OpenClawCalls)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "gateway --bind lan --port 18790") {
		t.Fatalf("openclaw should retry on port 18790:\n%v", resp.OpenClawCalls)
	}
	if !strings.Contains(resp.Stdout, "18789") || !strings.Contains(resp.Stdout, "18790") {
		t.Fatalf("stdout should note runtime port bump:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "http://127.0.0.1:18790/") {
		t.Fatalf("stdout missing retried dashboard URL:\n%s", resp.Stdout)
	}
}
```