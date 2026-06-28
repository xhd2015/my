## Expected

- Exit code `1`.
- Stderr says `--data-dir` is required and mentions `add data-dir`.

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
	if !strings.Contains(resp.Stderr, "requires --data-dir") {
		t.Fatalf("stderr missing requires --data-dir:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "add data-dir") {
		t.Fatalf("stderr missing add data-dir hint:\n%s", resp.Stderr)
	}
	if strings.Contains(resp.Stdout, "Using data dir:") {
		t.Fatalf("stdout should not auto-select data dir:\n%s", resp.Stdout)
	}
}

```