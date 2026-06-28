## Expected

- Exit code `1`.
- Stderr mentions container is not running.
- No `podman exec` call.

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
	if !strings.Contains(resp.Stderr, "not running") {
		t.Fatalf("stderr missing not-running message:\n%s", resp.Stderr)
	}
	if podmanCallsContain(resp.PodmanCalls, "podman exec") {
		t.Fatalf("unexpected podman exec when container not running:\n%v", resp.PodmanCalls)
	}
}
```