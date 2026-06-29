## Expected

- Non-zero exit code.
- Error message references missing grok auth file.

## Errors

- Grok auth file not found.

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
	if !strings.Contains(combined, "grok") && !strings.Contains(combined, "not found") {
		t.Fatalf("expected grok/missing error, got stdout=%q stderr=%q", resp.Stdout, resp.Stderr)
	}
}```
