# Scenario

**Feature**: import surfaces errors for invalid or missing grok auth inputs.

## Context

- Each leaf sets `MY_GROK_AUTH_PATH` to a fixture or missing path.

```go
func Setup(t *testing.T, req *Request) error {
	req.SecondGrokPath = ""
	return nil
}
```