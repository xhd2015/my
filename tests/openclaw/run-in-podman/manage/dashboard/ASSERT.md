## Expected

- Exit code `0`.
- Stderr previews `open` command with auth URL.
- Stdout prints authenticated dashboard URL.

## Exit Code

- `0`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr: %s", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "$ open ") {
		t.Fatalf("stderr missing open command preview:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "http://127.0.0.1:18789/#token=json-gateway-token") {
		t.Fatalf("stdout missing auth dashboard URL:\n%s", resp.Stdout)
	}
}
```