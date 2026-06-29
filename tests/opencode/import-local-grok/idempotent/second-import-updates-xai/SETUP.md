# Scenario

**Feature**: re-import overwrites xai access/refresh, preserves other providers if any.

```go
func Setup(t *testing.T, req *Request) error {
	req.GrokAuthPath = copyFixture(t, "testdata/grok-auth-valid.json")
	req.PreseedAuthJSON = []byte(`{"other":{"type":"api","key":"keep-me"}}`)
	req.SecondGrokPath = copyFixture(t, "testdata/grok-auth-valid-v2.json")
	return nil
}```
