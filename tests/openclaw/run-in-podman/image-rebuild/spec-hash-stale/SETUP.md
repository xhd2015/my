# Scenario

**Feature**: stale `image.spec_hash` in registry triggers rebuild.

```
my openclaw run-in-podman -> Podman stub (build because hash mismatch)
```

## Preconditions

- Image may or may not exist; stored hash does not match current Containerfile hash.
- Valid data dir with token.

## Steps

1. Seed registry with outdated `spec_hash`.
2. Run without `--rebuild`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.PodmanImageExists = true
	seed := `{
  "image": {
    "spec_hash": "sha256:stale-hash-value",
    "built_at": "2026-01-01T00:00:00Z"
  }
}`
	return writeRegistrySeed(filepath.Join(req.ConfigDir, "openclaw.json"), seed)
}
```