# Scenario

**Feature**: `--container-name` overrides default `openclaw-gateway`.

```
my openclaw run-in-podman --container-name my-gateway -> podman --name my-gateway
```

## Preconditions

- Valid data dir with token.

## Steps

1. Set `req.ContainerName = "my-gateway"`.
2. Run command.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-token")
	req.ContainerName = "my-gateway"
	req.PodmanImageExists = true
	return nil
}
```