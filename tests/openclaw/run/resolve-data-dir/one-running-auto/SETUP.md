# Scenario

**Feature**: exactly one registered data dir is running — auto-select and print status.

## Steps

1. Seed registry with two absolute data-dir paths.
2. Write bookkeeping for only the first path (`os.Getpid()`, port `18789`, busy port hook).
3. Run `my openclaw run --status` without `--data-dir`.

```go
func Setup(t *testing.T, req *Request) error {
	dirA := minimalDataDir(t, "auto-a")
	dirB := minimalDataDir(t, "auto-b")
	if err := seedRegistry(req, dirA, dirB); err != nil {
		return err
	}
	return markDirsRunning(t, req, dirA)
}

```