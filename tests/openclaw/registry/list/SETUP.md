# Scenario

**Feature**: `my openclaw list` prints registered data directories.

```
# read registry and format tab-separated table
my openclaw list <- Config store -> stdout table or empty message
```

## Preconditions

- Registry state is controlled by leaf setup (empty or pre-populated).

## Steps

1. Set `req.Subcommand` to `"list"`.
2. Leaf seeds registry when needed.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "list"
	return nil
}
```