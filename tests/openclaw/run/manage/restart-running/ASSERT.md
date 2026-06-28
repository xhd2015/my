## Expected

- Exit code `0`.
- `gateway stop` then detached `gateway --bind lan` start.
- Stdout confirms stop and gateway start.

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
	if !strings.Contains(resp.Stdout, "Stopped local gateway") {
		t.Fatalf("stdout missing stop confirmation:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Gateway started locally") {
		t.Fatalf("stdout missing launch confirmation:\n%s", resp.Stdout)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "gateway stop") {
		t.Fatalf("missing gateway stop:\n%v", resp.OpenClawCalls)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "gateway --bind lan --port 18789") {
		t.Fatalf("missing detached gateway start:\n%v", resp.OpenClawCalls)
	}
}
```