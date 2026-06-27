# Scenario

**Feature**: `--rebuild` forces `podman build` even when image exists with matching hash.

```
my openclaw run-in-podman --rebuild -> Podman stub (build + run)
```

## Preconditions

- Image already exists (`PodmanImageExists=true`).
- Registry has matching spec hash.

## Steps

1. Seed registry with current spec hash placeholder.
2. Set `req.Rebuild = true`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.Rebuild = true
	req.PodmanImageExists = true
	req.RegistrySeed = []byte(`{
  "image": {
    "spec_hash": "sha256:current-hash",
    "built_at": "2026-06-27T10:05:00+08:00"
  }
}`)
	if err := writeRegistrySeed(filepath.Join(req.ConfigDir, "openclaw.json"), string(req.RegistrySeed)); err != nil {
		return err
	}
	return nil
}
```