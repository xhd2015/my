## Expected

- Exit code `0`.
- Only `podman run` is called (no stop/rm).
- Stdout shows gateway started.

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
	if !strings.Contains(resp.Stdout, "Gateway started: openclaw-gateway") {
		t.Fatalf("stdout missing launch confirmation:\n%s", resp.Stdout)
	}
	if podmanCallsContain(resp.PodmanCalls, "podman stop openclaw-gateway") {
		t.Fatalf("unexpected podman stop:\n%v", resp.PodmanCalls)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman run") {
		t.Fatalf("missing podman run:\n%v", resp.PodmanCalls)
	}
}
```