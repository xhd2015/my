# Scenario

**Feature**: list with no registered data dirs prints empty message.

```
my openclaw list <- empty Config store -> "(no data dirs registered)"
```

## Preconditions

- No registry file exists.

## Steps

1. Run `my openclaw list` without seeding registry.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	_ = os.Remove(filepath.Join(req.ConfigDir, "openclaw.json"))
	return nil
}
```