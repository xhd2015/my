## Expected

- Exit code `1`.
- Stderr indicates the data directory is missing or invalid.

## Side Effects

- No `openclaw gateway` invocation recorded.

## Errors

- Data directory not found.

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
	if !strings.Contains(lower, "data") || (!strings.Contains(lower, "dir") && !strings.Contains(lower, "directory")) {
		t.Fatalf("stderr should mention data directory:\n%s", resp.Stderr)
	}
	if openclawCallsContain(resp.OpenClawCalls, "gateway") {
		t.Fatalf("openclaw gateway should not be called:\n%v", resp.OpenClawCalls)
	}
}
```
