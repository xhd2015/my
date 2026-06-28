## Expected

- Exit code `1`.
- Stderr says `--data-dir` is required when the container is not running.

## Exit Code

- `1`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1", resp.ExitCode)
	}
	if !strings.Contains(resp.Stderr, "--data-dir is required") {
		t.Fatalf("stderr missing data-dir requirement:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "not running") {
		t.Fatalf("stderr missing not-running context:\n%s", resp.Stderr)
	}
}
```