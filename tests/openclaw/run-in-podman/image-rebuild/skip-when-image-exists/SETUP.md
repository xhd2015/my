# Scenario

**Feature**: no rebuild when image exists and stored spec hash matches current hash.

```
my openclaw run-in-podman -> Podman stub (images ok, no build, run only)
```

## Preconditions

- Podman reports image `my-openclaw:local` exists.
- Registry `image.spec_hash` matches embedded Containerfile SHA256.

## Steps

1. Compute current Containerfile hash via helper expectation placeholder.
2. Seed registry with matching hash.
3. Run command without `--rebuild`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	req.PodmanImageExists = true
	hash, err := currentContainerfileHash(req.CommandDir)
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

func currentContainerfileHash(moduleRoot string) (string, error) {
	path := filepath.Join(moduleRoot, "internal", "openclaw", "container", "Containerfile")
	return hashFileSHA256(path)
}
```