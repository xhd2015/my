# Scenario

**Feature**: grok file without OIDC entry fails.

```go
func Setup(t *testing.T, req *Request) error {
	req.GrokAuthPath = copyFixture(t, "testdata/grok-auth-no-oidc.json")
	return nil
}```
