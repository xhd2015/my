# Scenario

**Feature**: `my openclaw add data-dir` registers an existing `.openclaw` directory.

```
# resolve path, verify directory exists, append or update registry entry
my openclaw add data-dir <path> [--note] -> Config store
```

## Preconditions

- Registry may be empty or pre-seeded depending on leaf.

## Steps

1. Set `req.Subcommand` to `"add"`.
2. Descendants set `DataDirPath`, `Note`, and registry preconditions.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "add"
	return nil
}
```