# Scenario

**Feature**: `--port` overrides default host port `18789`.

```
my openclaw run-in-podman --port 19001 -> -p 19001:18789 and dashboard URL on 19001
```

## Preconditions

- Valid data dir with token.

## Steps

1. Set `req.Port = "19001"`.
2. Run command.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.Port = "19001"
	req.PodmanImageExists = true
	return nil
}
```