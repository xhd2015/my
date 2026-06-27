## Expected

- Exit code `1`.
- Stderr mentions missing gateway token.

## Side Effects

- No `podman run` invocation recorded.

## Errors

- Token not found in json or `.env`.

## Exit Code

- `1`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	lower := strings.ToLower(resp.Stderr)
	if !strings.Contains(lower, "token") {
		t.Fatalf("stderr should mention token:\n%s", resp.Stderr)
	}
	if podmanCallsContain(resp.PodmanCalls, "podman run") {
		t.Fatalf("podman run should not be called:\n%v", resp.PodmanCalls)
	}
}
```