## Expected

- Exit code `0`.
- Podman `stop` and `rm` called for default container.
- Stderr previews both commands.
- Stdout confirms stop.

## Exit Code

- `0`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0", resp.ExitCode)
	}
	if !strings.Contains(resp.Stderr, "$ podman stop") || !strings.Contains(resp.Stderr, "$ podman rm") {
		t.Fatalf("stderr missing command preview:\n%s", resp.Stderr)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman stop openclaw-gateway") {
		t.Fatalf("missing podman stop:\n%v", resp.PodmanCalls)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman rm openclaw-gateway") {
		t.Fatalf("missing podman rm:\n%v", resp.PodmanCalls)
	}
	if !strings.Contains(resp.Stdout, "Stopped and removed container openclaw-gateway") {
		t.Fatalf("stdout missing confirmation:\n%s", resp.Stdout)
	}
}
```