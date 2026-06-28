# Scenario

**Feature**: `--show-tokens` auto-detects the data dir from a running gateway container.

## Steps

1. Seed a running container with a known mounted data dir.
2. Run `my openclaw run-in-podman --show-tokens` without `--data-dir`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.ShowTokens = true
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run-in-podman", "testdata", "with-token"))
	req.PodmanContainerDataDir = req.RunDataDir
	return markContainerRunning(req, "openclaw-gateway")
}
```