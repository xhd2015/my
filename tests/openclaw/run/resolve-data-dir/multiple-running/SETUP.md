# Scenario

**Feature**: more than one registered data dir is running — must pass `--data-dir` explicitly.

## Steps

1. Seed registry with two absolute data-dir paths.
2. Write bookkeeping for both paths (`os.Getpid()`, port `18789`, busy port hook).
3. Run `my openclaw run --status` without `--data-dir`.

```go
func Setup(t *testing.T, req *Request) error {
	dirA := minimalDataDir(t, "running-a")
	dirB := minimalDataDir(t, "running-b")
	if err := seedRegistry(req, dirA, dirB); err != nil {
		return err
	}
	return markDirsRunning(t, req, dirA, dirB)
}

```