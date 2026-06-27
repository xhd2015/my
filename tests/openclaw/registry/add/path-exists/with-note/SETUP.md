# Scenario

**Feature**: optional `--note` is persisted with the registry entry.

```
my openclaw add data-dir <dir> --note "..." -> Config store (note field)
```

## Preconditions

- No prior registry file.
- Data directory exists.

## Steps

1. Run `my openclaw add data-dir <path> --note "dev laptop"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Note = "dev laptop"
	return nil
}
```