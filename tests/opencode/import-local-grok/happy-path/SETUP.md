# Scenario

**Feature**: first import into empty OpenCode auth store.

## Steps

1. Leaf tests set `MY_GROK_AUTH_PATH` to grok fixtures.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreseedAuthJSON = nil
	return nil
}
```