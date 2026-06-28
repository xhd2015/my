## Expected

- Exit code `1`.
- Stderr contains `openclaw not found` and `npm install -g openclaw@latest`.

## Side Effects

- No `openclaw gateway` invocation recorded.

## Errors

- OpenClaw CLI missing from `PATH`.

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
	if !strings.Contains(lower, "openclaw not found") {
		t.Fatalf("stderr missing openclaw not found:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "npm install -g openclaw@latest") {
		t.Fatalf("stderr missing install instructions:\n%s", resp.Stderr)
	}
	if len(resp.OpenClawCalls) > 0 {
		t.Fatalf("openclaw should not be invoked:\n%v", resp.OpenClawCalls)
	}
}
```
