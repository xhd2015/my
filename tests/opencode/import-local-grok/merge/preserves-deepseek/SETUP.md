# Scenario

**Feature**: `deepseek` entry unchanged after xai import.

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
