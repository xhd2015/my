# Scenario

**Feature**: `--show-tokens` prints gateway tokens and auth URLs.

## Steps

1. Copy `with-token` fixture.
2. Run `my openclaw run-in-podman --show-tokens --data-dir <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.ShowTokens = true
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run-in-podman", "testdata", "with-token"))
	return nil
}
```