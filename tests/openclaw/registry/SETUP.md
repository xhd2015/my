# Scenario

**Feature**: registry subcommands manage explicit OpenClaw data-dir bookkeeping.

```
# add registers paths; list prints registered data dirs
my openclaw add/list -> Config store (registry JSON)
```

## Preconditions

- Isolated `MY_CONFIG_DIR` from root setup.
- `my` binary built and on PATH via `BinDir`.

## Steps

1. Set `req.Subcommand` to either `add` or `list` in descendant nodes.

```go
func Setup(t *testing.T, req *Request) error {
	req.Rebuild = false
	req.PodmanImageExists = false
	req.ContainerName = ""
	req.Port = ""
	return nil
}
```