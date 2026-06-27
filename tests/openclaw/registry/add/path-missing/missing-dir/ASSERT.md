## Expected

- Exit code `1`.
- Stderr mentions the path was not found.

## Side Effects

- Registry file is not created or unchanged.

## Errors

- Path not found error on stderr.

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
	if !strings.Contains(strings.ToLower(resp.Stderr), "not found") && !strings.Contains(resp.Stderr, req.DataDirPath) {
		t.Fatalf("stderr should mention missing path:\n%s", resp.Stderr)
	}
}
```