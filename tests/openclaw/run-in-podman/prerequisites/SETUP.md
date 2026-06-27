# Scenario

**Feature**: run-in-podman validates required inputs before launching Podman.

```
# missing data dir, openclaw.json, or gateway token
my openclaw run-in-podman --data-dir <path> -> validation error, exit 1
```

## Preconditions

- Tests exercise invalid or incomplete data directory setups.

## Steps

1. Configure `req.RunDataDir` and fixtures per leaf.
2. No successful `podman run` should occur.

```go
func Setup(t *testing.T, req *Request) error {
	req.Rebuild = false
	req.PodmanImageExists = false
	return nil
}
```