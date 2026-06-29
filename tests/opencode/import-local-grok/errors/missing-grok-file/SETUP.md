# Scenario

**Feature**: missing grok auth file returns error.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.GrokAuthPath = filepath.Join(t.TempDir(), "does-not-exist", "auth.json")
	return nil
}```
