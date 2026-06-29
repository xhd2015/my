# Scenario

**Feature**: second import updates only `xai` credentials.

## Steps

1. First import from `grok-auth-valid.json`.
2. Second import from `grok-auth-valid-v2.json` via `SecondGrokPath`.

```go
func Setup(t *testing.T, req *Request) error {
	req.GrokAuthPath = copyFixture(t, "testdata/grok-auth-valid.json")
	req.SecondGrokPath = copyFixture(t, "testdata/grok-auth-valid-v2.json")
	return nil
}```
