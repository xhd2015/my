## Expected

- Exit code `0`.
- Npm stub receives `install -g openclaw@latest`.
- OpenClaw stub receives `gateway --bind lan --port 18789`.
- Stdout contains dashboard URL `http://127.0.0.1:18789/`.

## Side Effects

- Npm install runs once, then gateway launches via newly installed openclaw stub.

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
		t.Fatalf("exit code = %d, want 0\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !npmCallsContain(resp.NpmCalls, "install -g openclaw@latest") {
		t.Fatalf("npm should run global openclaw install:\n%v", resp.NpmCalls)
	}
	if !openclawCallsContain(resp.OpenClawCalls, "gateway --bind lan --port 18789") {
		t.Fatalf("openclaw should launch gateway after install:\n%v", resp.OpenClawCalls)
	}
	if !strings.Contains(resp.Stdout, "http://127.0.0.1:18789/") {
		t.Fatalf("stdout missing dashboard URL:\n%s", resp.Stdout)
	}
}
```