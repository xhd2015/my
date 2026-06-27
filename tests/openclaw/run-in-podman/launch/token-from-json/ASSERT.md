## Expected

- Exit code `0`.
- Podman `run` passes `-e OPENCLAW_GATEWAY_TOKEN=json-gateway-token`.

## Side Effects

- Token from json used; `.env` not required.

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
	if !podmanCallsContain(resp.PodmanCalls, "OPENCLAW_GATEWAY_TOKEN=json-gateway-token") {
		t.Fatalf("podman run missing json token env:\n%v", resp.PodmanCalls)
	}
}
```