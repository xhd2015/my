# Scenario

**Feature**: add fails when the directory does not exist.

```
my openclaw add data-dir <missing> -> stderr path not found, exit 1
```

## Preconditions

- Target path does not exist on disk.

## Steps

1. Run `my openclaw add data-dir <missing-path>`.

```go
import "os"

func Setup(t *testing.T, req *Request) error {
	if _, err := os.Stat(req.DataDirPath); !os.IsNotExist(err) {
		t.Fatalf("precondition: data dir must not exist before add: %v", err)
	}
	return nil
}
```