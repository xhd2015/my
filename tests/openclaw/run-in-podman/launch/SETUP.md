# Scenario

**Feature**: successful launch runs gateway container with correct mounts and token.

```
# token resolved, image ready, container started with defaults or custom flags
my openclaw run-in-podman -> Podman stub (stop/rm/run) -> dashboard URL on stdout
```

## Preconditions

- Valid data directory fixture with required config files.
- Podman image exists unless leaf tests first-time build path separately.

## Steps

1. Copy leaf-specific fixture into temp data dir.
2. Set `req.RunDataDir`.
3. Leaf sets container name, port, and token source fixtures.

```go
func Setup(t *testing.T, req *Request) error {
	req.PodmanImageExists = true
	return nil
}
```