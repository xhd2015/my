# Scenario

**Feature**: `--test-slack` errors when the local gateway is not running.

## Steps

1. Ensure gateway is not running.
2. Run `my openclaw run --test-slack`.

```go
func Setup(t *testing.T, req *Request) error {
	req.TestSlack = true
	return nil
}
```