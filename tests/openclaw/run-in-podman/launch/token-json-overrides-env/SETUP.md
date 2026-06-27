# Scenario

**Feature**: json token takes precedence over `.env` when both are present.

```
Token resolver <- openclaw.json (preferred) over .env -> json token in podman run
```

## Preconditions

- Both `openclaw.json` token and `.env` token exist with different values.

## Steps

1. Use `with-both` fixture.
2. Run command.

```go
func Setup(t *testing.T, req *Request) error {
	req.RunDataDir = fixtureDataDir(t, "with-both")
	req.PodmanImageExists = true
	return nil
}
```