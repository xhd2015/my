# Scenario

**Feature**: port busy without bookkeeping must not count as running — auto-select the truly running dir.

Only the dir with valid `.my/gateway.json` (alive pid + port in use) counts as running;
a busy default port alone is insufficient.

## Steps

1. Seed registry with two absolute data-dir paths.
2. Write bookkeeping for only the first path (`os.Getpid()`, port `18789`).
3. Set `OPENCLAW_STUB_BUSY_PORTS=18789` so the bookkeeping port appears in use.
4. Run `my openclaw run --status` without `--data-dir`.

```go
func Setup(t *testing.T, req *Request) error {
	dirA := minimalDataDir(t, "port-busy-a")
	dirB := minimalDataDir(t, "port-busy-b")
	if err := seedRegistry(req, dirA, dirB); err != nil {
		return err
	}
	return markDirsRunning(t, req, dirA)
}

```