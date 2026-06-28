# Scenario

**Feature**: `--restart` stops a running container and starts it again (auto-detects data dir from mount).

## Steps

1. Seed a running container with a known mounted data dir.
2. Run `my openclaw run-in-podman --restart` without `--data-dir`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Restart = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.PodmanContainerDataDir = req.RunDataDir
	req.PodmanImageExists = true
	hash, err := hashFileSHA256(filepath.Join(req.CommandDir, "internal", "openclaw", "container", "Containerfile"))
	if err != nil {
		return err
	}
	seed := `{
  "image": {
    "spec_hash": "` + hash + `",
    "built_at": "2026-06-27T10:05:00+08:00"
  }
}`
	if err := writeRegistrySeed(filepath.Join(req.ConfigDir, "openclaw.json"), seed); err != nil {
		return err
	}
	return markContainerRunning(req, "openclaw-gateway")
}
```