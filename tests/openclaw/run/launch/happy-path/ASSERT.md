## Expected

- Exit code `0`.
- Stdout contains dashboard URL `http://127.0.0.1:18789/`.
- Stub receives `gateway --bind lan --port 18789`.
- Stub env `OPENCLAW_STATE_DIR` equals absolute data dir.
- Stderr previews `openclaw gateway` command with state dir and token.

## Side Effects

- Gateway launched via host `openclaw` binary.

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
	if !strings.Contains(resp.Stdout, "http://127.0.0.1:18789/") {
		t.Fatalf("stdout missing dashboard URL:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, "openclaw gateway") {
		t.Fatalf("stderr missing command preview:\n%s", resp.Stderr)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "gateway --bind lan --port 18789") {
		t.Fatalf("openclaw missing gateway invocation:\n%v", resp.OpenClawCalls)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "OPENCLAW_STATE_DIR="+req.RunDataDir) {
		t.Fatalf("openclaw missing state dir env:\n%v", resp.OpenClawCalls)
	}
}
```
