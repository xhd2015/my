# Scenario

**Feature**: import creates `xai` oauth entry and secure file mode.

## Steps

1. Use valid grok fixture from `testdata/grok-auth-valid.json`.

```go
func Setup(t *testing.T, req *Request) error {
	req.GrokAuthPath = copyFixture(t, "testdata/grok-auth-valid.json")
	return nil
}
```