# Scenario

**Feature**: successful local launch execs `openclaw gateway` with correct env, port, and output.

```
# token resolved, port selected, gateway started on host
my openclaw run --data-dir <valid> -> openclaw stub -> dashboard URL on stdout
```

## Preconditions

- Valid data directory fixture with required config files.
- Stub `openclaw` installed on `PATH`.
- Runtime port retry uses `OPENCLAW_STUB_FAIL_PORTS` so the first launch fails with `EADDRINUSE`.

## Steps

1. Copy leaf-specific fixture into temp data dir.
2. Set `req.RunDataDir`, port, and busy-port simulation per leaf.
3. Invoke shared `Run` from `DOCTEST.md`.

```go
func Setup(t *testing.T, req *Request) error {
	req.BusyPorts = ""
	return nil
}
```
