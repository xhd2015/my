## Expected

- Exit code `1` (explicit second dir is not running).
- Stdout does not print `Using data dir:`.
- Stderr reports gateway not running.

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
	if strings.Contains(resp.Stdout, "Using data dir:") {
		t.Fatalf("stdout should not auto-select when --data-dir is explicit:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, "not running") {
		t.Fatalf("stderr missing not-running error:\n%s", resp.Stderr)
	}
}

```