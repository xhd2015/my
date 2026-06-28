## Expected

- Exit code `1`.
- Stderr mentions `openclaw.json`.

## Side Effects

- No `openclaw gateway` invocation recorded.

## Errors

- Missing required config file.

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
	if !strings.Contains(resp.Stderr, "openclaw.json") {
		t.Fatalf("stderr missing openclaw.json mention:\n%s", resp.Stderr)
	}
	if openclawCallsContain(resp.OpenClawCalls, "gateway") {
		t.Fatalf("openclaw gateway should not be called:\n%v", resp.OpenClawCalls)
	}
}
```
