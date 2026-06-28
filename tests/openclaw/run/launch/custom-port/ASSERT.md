## Expected

- Exit code `0`.
- Stub receives `gateway --bind lan --port 19999`.
- Stdout dashboard URL uses port `19999`.

## Side Effects

- Custom gateway port selected.

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
	if !openclawCallsContain(resp.OpenClawCalls, "gateway --bind lan --port 19999") {
		t.Fatalf("openclaw missing custom port:\n%v", resp.OpenClawCalls)
	}
	if !strings.Contains(resp.Stdout, "http://127.0.0.1:19999/") {
		t.Fatalf("stdout missing custom dashboard URL:\n%s", resp.Stdout)
	}
}
```
