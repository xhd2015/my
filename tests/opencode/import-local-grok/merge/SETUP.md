# Scenario

**Feature**: import merges `xai` into existing OpenCode auth without removing other providers.

## Steps

1. Preseed `auth.json` with a `deepseek` api entry.
2. Import from valid grok fixture.

```go
func Setup(t *testing.T, req *Request) error {
	req.GrokAuthPath = copyFixture(t, "testdata/grok-auth-valid.json")
	req.PreseedAuthJSON = []byte(`{
  "deepseek": {
    "type": "api",
    "key": "deepseek-secret-key"
  }
}`)
	return nil
}```
