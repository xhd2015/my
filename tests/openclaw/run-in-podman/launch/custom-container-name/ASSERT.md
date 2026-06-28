## Expected

- Exit code `0`.
- Podman `run` uses `--name my-gateway`.
- Stdout help references `my-gateway`.

## Side Effects

- Custom container name used throughout lifecycle.

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
	if !podmanCallsContain(resp.PodmanCalls, "--name my-gateway") {
		t.Fatalf("missing custom container name in podman calls:\n%v", resp.PodmanCalls)
	}
	if !strings.Contains(resp.Stdout, "--container-name my-gateway") {
		t.Fatalf("stdout missing my CLI help for custom name:\n%s", resp.Stdout)
	}
}
```