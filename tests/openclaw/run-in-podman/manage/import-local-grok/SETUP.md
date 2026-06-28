# Scenario

**Feature**: `--import-local-grok` copies local Grok Build OAuth into the OpenClaw data dir.

## Steps

1. Copy a valid data-dir fixture with `openclaw.json`.
2. Point `MY_GROK_AUTH_PATH` at leaf `testdata/grok-auth.json`.
3. Run `my openclaw run-in-podman --import-local-grok --data-dir <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.ImportLocalGrok = true
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run-in-podman", "manage", "import-local-grok", "testdata", "data-dir"))
	req.GrokAuthPath = filepath.Join(DOCTEST_ROOT, "run-in-podman", "manage", "import-local-grok", "testdata", "grok-auth.json")
	return nil
}
```