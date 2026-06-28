## Expected

- Exit code `1`.
- Stderr mentions port `18789` is in use.

## Side Effects

- No `openclaw gateway` invocation recorded.

## Errors

- Requested port unavailable.

## Exit Code

- `1`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	lower := strings.ToLower(resp.Stderr)
	if !strings.Contains(lower, "18789") || !strings.Contains(lower, "in use") {
		t.Fatalf("stderr should mention port in use:\n%s", resp.Stderr)
	}
	if openclawCallsContain(resp.OpenClawCalls, "gateway") {
		t.Fatalf("openclaw gateway should not be called:\n%v", resp.OpenClawCalls)
	}
}
```
