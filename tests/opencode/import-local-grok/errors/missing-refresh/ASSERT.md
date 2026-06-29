## Expected

- Non-zero exit.
- Error mentions missing refresh token.

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
	if !strings.Contains(combined, "refresh") {
		t.Fatalf("expected refresh token error, got %q %q", resp.Stdout, resp.Stderr)
	}
}```
