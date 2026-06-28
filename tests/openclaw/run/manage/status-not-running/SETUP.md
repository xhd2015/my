# Scenario

**Feature**: `--status` errors when the local gateway is not running.

## Steps

1. Ensure gateway is not running (default stub state).
2. Run `my openclaw run --status`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Status = true
	return nil
}
```