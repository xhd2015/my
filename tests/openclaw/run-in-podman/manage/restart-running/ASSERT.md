## Expected

- Exit code `0`.
- Podman `stop`, `rm`, and `run` called for default container.
- Stdout confirms stop and gateway start.

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
	if !strings.Contains(resp.Stdout, "Stopped and removed container openclaw-gateway") {
		t.Fatalf("stdout missing stop confirmation:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Gateway started: openclaw-gateway") {
		t.Fatalf("stdout missing launch confirmation:\n%s", resp.Stdout)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman stop openclaw-gateway") {
		t.Fatalf("missing podman stop:\n%v", resp.PodmanCalls)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman rm openclaw-gateway") {
		t.Fatalf("missing podman rm:\n%v", resp.PodmanCalls)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman run") {
		t.Fatalf("missing podman run:\n%v", resp.PodmanCalls)
	}
	var runLine string
	for _, call := range resp.PodmanCalls {
		if strings.Contains(call, "podman run") {
			runLine = call
			break
		}
	}
	if !strings.Contains(runLine, req.PodmanContainerDataDir+":/home/node/.openclaw") {
		t.Fatalf("podman run missing data dir mount in:\n%s", runLine)
	}
}
```