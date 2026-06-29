## Expected

- Non-zero exit.
- Error mentions missing access token.

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
	if !strings.Contains(combined, "access") {
		t.Fatalf("expected access token error, got %q %q", resp.Stdout, resp.Stderr)
	}
}
```