# Scenario

**Feature**: explicit `--data-dir` wins over auto-resolve; no auto-select message.

## Steps

1. Seed registry with two absolute data-dir paths.
2. Write bookkeeping for only the first path (auto would pick it).
3. Run `my openclaw run --status --data-dir <second>`.

```go
func Setup(t *testing.T, req *Request) error {
	dirA := minimalDataDir(t, "explicit-a")
	dirB := minimalDataDir(t, "explicit-b")
	if err := seedRegistry(req, dirA, dirB); err != nil {
		return err
	}
	if err := markDirsRunning(t, req, dirA); err != nil {
		return err
	}
	req.RunDataDir = dirB
	return nil
}

```