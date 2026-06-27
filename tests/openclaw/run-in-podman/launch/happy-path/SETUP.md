# Scenario

**Feature**: default launch runs gateway with standard mounts, token, and output URLs.

```
my openclaw run-in-podman --data-dir <valid> -> podman stop/rm/run -> dashboard URL + logs hint
```

## Preconditions

- Valid data dir with token in `openclaw.json`.
- Image exists with matching hash (no rebuild).

## Steps

1. Use `with-token` fixture.
2. Run with default container name and port.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
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