# Scenario

**Feature**: `--import-local-grok` auto-detects the data dir from a running gateway container.

## Steps

1. Seed a running container with a known mounted data dir.
2. Run `my openclaw run-in-podman --import-local-grok` without `--data-dir`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.ImportLocalGrok = true
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run-in-podman", "manage", "import-local-grok", "testdata", "data-dir"))
	req.PodmanContainerDataDir = req.RunDataDir
	req.GrokAuthPath = filepath.Join(DOCTEST_ROOT, "run-in-podman", "manage", "import-local-grok", "testdata", "grok-auth.json")
	return markContainerRunning(req, "openclaw-gateway")
}
```