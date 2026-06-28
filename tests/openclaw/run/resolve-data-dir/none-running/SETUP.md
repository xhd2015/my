# Scenario

**Feature**: registered data dirs exist but none are running — list dirs and exit 1.

## Steps

1. Seed registry with two absolute data-dir paths.
2. Do not write bookkeeping for any path.
3. Run `my openclaw run --status` without `--data-dir`.

```go
func Setup(t *testing.T, req *Request) error {
	dirA := minimalDataDir(t, "registered-a")
	dirB := minimalDataDir(t, "registered-b")
	return seedRegistry(req, dirA, dirB)
}

```