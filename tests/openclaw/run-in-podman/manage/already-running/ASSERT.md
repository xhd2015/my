## Expected

- Exit code `1`.
- Stderr mentions already running and `--stop`.

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
	if !strings.Contains(resp.Stderr, "already running") {
		t.Fatalf("stderr missing already running:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "--stop") {
		t.Fatalf("stderr missing --stop hint:\n%s", resp.Stderr)
	}
	if podmanCallsContain(resp.PodmanCalls, "podman run") {
		t.Fatalf("unexpected podman run when already running:\n%v", resp.PodmanCalls)
	}
}
```