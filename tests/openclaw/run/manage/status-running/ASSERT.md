## Expected

- Exit code `0`.
- Stdout shows local gateway status, data dir, port, URLs, and auth URLs.
- No gateway stop or relaunch side effects.

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
		"Gateway running locally",
		"Data dir: " + req.RunDataDir,
		"Port: 18789",
		"http://127.0.0.1:18789/",
		"http://127.0.0.1:18789/chat?session=main",
		"http://127.0.0.1:18789/#token=json-gateway-token",
		"Gateway token: configured",
	} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
	if openclawCallsContain(resp.OpenClawCalls, "gateway stop") {
		t.Fatalf("unexpected gateway stop:\n%v", resp.OpenClawCalls)
	}
}
```