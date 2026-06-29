## Expected

- Non-zero exit.
- Error indicates no OIDC / usable grok entry.

## Exit Code

- Non-zero

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	assertErrorResponse(t, resp)
	combined := strings.ToLower(resp.Stdout + resp.Stderr)
	if !strings.Contains(combined, "oidc") && !strings.Contains(combined, "entry") {
		t.Fatalf("expected OIDC/entry error, got %q %q", resp.Stdout, resp.Stderr)
	}
}```
