## Expected

- Exit code `1`.
- Stderr mentions grok auth file not found.

## Exit Code

- `1`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1", resp.ExitCode)
	}
	if !strings.Contains(resp.Stderr, "grok auth file not found") {
		t.Fatalf("stderr missing grok auth error:\n%s", resp.Stderr)
	}
}
```