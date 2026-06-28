## Expected

- Exit code `0`.
- Stdout shows container status, data dir, port, URLs, and auth URLs.
- No podman stop/run side effects.

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
	for _, want := range []string{
		"Container: openclaw-gateway (running)",
		"Data dir: " + req.PodmanContainerDataDir,
		"Port: 18789",
		"http://127.0.0.1:18789/",
		"http://127.0.0.1:18789/chat?session=main",
		"http://127.0.0.1:18789/#token=json-gateway-token",
		"Gateway token: configured",
		"my openclaw run-in-podman --logs --container-name openclaw-gateway",
	} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
	if podmanCallsContain(resp.PodmanCalls, "podman run") {
		t.Fatalf("unexpected podman run:\n%v", resp.PodmanCalls)
	}
}
```