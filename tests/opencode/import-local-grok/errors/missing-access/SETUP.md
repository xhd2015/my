# Scenario

**Feature**: grok OIDC entry without access token fails.

```go
func Setup(t *testing.T, req *Request) error {
	req.GrokAuthPath = copyFixture(t, "testdata/grok-auth-missing-access.json")
	return nil
}```
