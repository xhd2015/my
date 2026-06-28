# Scenario

**Feature**: `--dashboard` opens authenticated dashboard URL.

## Steps

1. Copy `with-token` fixture.
2. Run `my openclaw run-in-podman --dashboard --data-dir <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Dashboard = true
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run-in-podman", "testdata", "with-token"))
	return nil
}
```