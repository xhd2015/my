## Expected

- Exit code `0`.
- Stdout shows json token and auth URLs.
- Stderr previews token read command.

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
	if !strings.Contains(resp.Stderr, "$ read gateway tokens from") {
		t.Fatalf("stderr missing command preview:\n%s", resp.Stderr)
	}
	for _, want := range []string{
		"openclaw.json gateway.auth.token: json-gateway-token",
		"Effective token: json-gateway-token",
		"#token=json-gateway-token",
	} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
}
```