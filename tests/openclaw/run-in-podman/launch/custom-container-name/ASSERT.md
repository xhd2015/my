## Expected

- Exit code `0`.
- Podman `stop`, `rm`, and `run` use `--name my-gateway`.
- Stdout logs hint references `my-gateway`.

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
	for _, want := range []string{
		"podman stop my-gateway",
		"podman rm my-gateway",
		"--name my-gateway",
	} {
		if !podmanCallsContain(resp.PodmanCalls, want) {
			t.Fatalf("missing %q in podman calls:\n%v", want, resp.PodmanCalls)
		}
	}
	if !strings.Contains(resp.Stdout, "podman logs -f my-gateway") {
		t.Fatalf("stdout missing logs hint for custom name:\n%s", resp.Stdout)
	}
}
```