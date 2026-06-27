## Expected

- Exit code `0`.
- Stdout is exactly `(no data dirs registered)` (trimmed).

## Side Effects

- None.

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
	if strings.TrimSpace(resp.Stdout) != "(no data dirs registered)" {
		t.Fatalf("stdout = %q, want %q", strings.TrimSpace(resp.Stdout), "(no data dirs registered)")
	}
}
```