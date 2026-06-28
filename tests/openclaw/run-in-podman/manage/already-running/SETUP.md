# Scenario

**Feature**: start errors when container already running.

## Steps

1. Seed running container and valid data dir fixture.
2. Run start command.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.RunDataDir = copyFixtureDir(t, filepath.Join(DOCTEST_ROOT, "run-in-podman", "testdata", "with-token"))
	return markContainerRunning(req, "openclaw-gateway")
}
```