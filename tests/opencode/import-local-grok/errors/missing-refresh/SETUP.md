# Scenario

**Feature**: grok OIDC entry without refresh token fails.

```go
func Setup(t *testing.T, req *Request) error {
	req.GrokAuthPath = copyFixture(t, "testdata/grok-auth-missing-refresh.json")
	return nil
}```
