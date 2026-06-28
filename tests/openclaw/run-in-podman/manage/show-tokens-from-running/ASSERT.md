## Expected

- Exit code `0`.
- Stderr previews reading tokens from the auto-detected mounted data dir.
- Stdout shows the gateway token from that data dir.

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
	if !strings.Contains(resp.Stderr, "read gateway tokens from") || !strings.Contains(resp.Stderr, req.PodmanContainerDataDir) {
		t.Fatalf("stderr missing auto-detected data dir %q:\n%s", req.PodmanContainerDataDir, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Effective token: json-gateway-token") {
		t.Fatalf("stdout missing effective token:\n%s", resp.Stdout)
	}
}
```