## Expected

- Exit code `0`.
- No `gateway stop` call; detached gateway start logged.
- Stub receives detached `gateway --bind lan` invocation.

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
	if !strings.Contains(resp.Stdout, "Gateway started locally") {
		t.Fatalf("stdout missing launch confirmation:\n%s", resp.Stdout)
	}
	if openclawCallsContain(resp.OpenClawCalls, "gateway stop") {
		t.Fatalf("unexpected gateway stop:\n%v", resp.OpenClawCalls)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "gateway --bind lan --port 18789") {
		t.Fatalf("missing detached gateway start:\n%v", resp.OpenClawCalls)
	}
}
```