## Expected

- Exit code `0`.
- Registry entry `note` field equals `"dev laptop"`.

## Side Effects

- `data_dirs` contains one entry with the given note.

## Errors

- None.

## Exit Code

- `0`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	reg, err := parseRegistry(resp.RegistryJSON)
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.DataDirs) != 1 {
		t.Fatalf("data_dirs count = %d, want 1", len(reg.DataDirs))
	}
	if reg.DataDirs[0].Note != "dev laptop" {
		t.Fatalf("note = %q, want %q", reg.DataDirs[0].Note, "dev laptop")
	}
}
```