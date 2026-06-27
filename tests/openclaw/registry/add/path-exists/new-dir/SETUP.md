# Scenario

**Feature**: first registration creates registry file with absolute path.

```
my openclaw add data-dir <new-dir> -> Config store (new entry, absolute path)
```

## Preconditions

- No prior registry file exists.
- Data directory exists.

## Steps

1. Run `my openclaw add data-dir <path>` without `--note`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Note = ""
	return nil
}
```