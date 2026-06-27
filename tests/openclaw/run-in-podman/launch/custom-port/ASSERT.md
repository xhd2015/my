## Expected

- Exit code `0`.
- Podman `run` includes `-p 19001:18789`.
- Stdout dashboard URL uses port `19001`.

## Side Effects

- Host port mapping customized; container port remains `18789`.

## Errors

- None.

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
	if !podmanCallsContain(resp.PodmanCalls, "-p 19001:18789") {
		t.Fatalf("podman run missing custom port mapping:\n%v", resp.PodmanCalls)
	}
	if !strings.Contains(resp.Stdout, "http://127.0.0.1:19001/") {
		t.Fatalf("stdout missing custom dashboard URL:\n%s", resp.Stdout)
	}
}
```