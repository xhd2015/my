# Scenario

**Feature**: container image rebuild follows missing image, stale hash, or `--rebuild`.

```
# compare registry image.spec_hash with Containerfile SHA256
my openclaw run-in-podman -> Podman stub (images/build/run)
```

## Preconditions

- Valid data directory with `openclaw.json` containing gateway token.
- Podman image presence controlled via `PodmanImageExists` and registry seed.

## Steps

1. Copy `testdata/with-token` fixture to temp data dir.
2. Set `req.RunDataDir` to fixture path.
3. Leaf configures rebuild flags and registry image metadata.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	return nil
}
```