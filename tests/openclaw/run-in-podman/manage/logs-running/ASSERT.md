## Expected

- Exit code `0`.
- Podman `logs` called.
- Stderr previews `podman logs` command.
- Stdout contains stub log line.

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
	if !strings.Contains(resp.Stderr, "$ podman logs") {
		t.Fatalf("stderr missing command preview:\n%s", resp.Stderr)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman logs openclaw-gateway") {
		t.Fatalf("missing podman logs:\n%v", resp.PodmanCalls)
	}
	if !strings.Contains(resp.Stdout, "stub log line") {
		t.Fatalf("stdout missing log output:\n%s", resp.Stdout)
	}
}
```