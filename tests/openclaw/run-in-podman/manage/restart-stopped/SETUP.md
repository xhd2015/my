# Scenario

**Feature**: `--restart` starts gateway when container is not running (requires `--data-dir`).

## Steps

1. Ensure container is not running.
2. Run `my openclaw run-in-podman --restart --data-dir <fixture>`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "run-in-podman"
	req.Restart = true
	req.RunDataDir = fixtureDataDir(t, "with-token")
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
	return writeRegistrySeed(filepath.Join(req.ConfigDir, "openclaw.json"), seed)
}
```