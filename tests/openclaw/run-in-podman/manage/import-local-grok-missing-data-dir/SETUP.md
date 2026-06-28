# Scenario

**Feature**: `--import-local-grok` requires `--data-dir` when no gateway container is running.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.ImportLocalGrok = true
	req.GrokAuthPath = filepath.Join(DOCTEST_ROOT, "run-in-podman", "manage", "import-local-grok", "testdata", "grok-auth.json")
	return nil
}
```