# Scenario

**Feature**: `--import-local-grok` errors when the grok auth file is missing.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.ImportLocalGrok = true
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run-in-podman", "testdata", "with-token"))
	req.GrokAuthPath = filepath.Join(t.TempDir(), "missing-grok-auth.json")
	return nil
}
```