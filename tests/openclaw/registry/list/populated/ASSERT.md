## Expected

- Exit code `0`.
- Stdout contains header with `path`, `note`, and `added_at`.
- Stdout contains the seeded path, note, and timestamp.

## Side Effects

- Registry unchanged.

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
	out := resp.Stdout
	for _, col := range []string{"path", "note", "added_at"} {
		if !strings.Contains(out, col) {
			t.Fatalf("stdout missing column %q:\n%s", col, out)
		}
	}
	if !strings.Contains(out, "work machine") {
		t.Fatalf("stdout missing note:\n%s", out)
	}
	if !strings.Contains(out, "2026-06-27T10:00:00+08:00") {
		t.Fatalf("stdout missing added_at:\n%s", out)
	}
}
```