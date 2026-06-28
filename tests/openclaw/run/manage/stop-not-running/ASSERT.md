## Expected

- Exit code `0`.
- Stderr warns gateway is not running.
- No `gateway stop` call.

## Exit Code

- `0`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0", resp.ExitCode)
	}
	if !strings.Contains(resp.Stderr, "warning:") || !strings.Contains(resp.Stderr, "not running") {
		t.Fatalf("stderr missing warning:\n%s", resp.Stderr)
	}
	if openclawCallsContain(resp.OpenClawCalls, "gateway stop") {
		t.Fatalf("unexpected gateway stop:\n%v", resp.OpenClawCalls)
	}
}
```