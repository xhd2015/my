## Expected

- Exit code `0`.
- Stderr contains `data dir already registered: <path>`.
- Registry note updated to `"updated"`.
- Still exactly one registry entry for the path.

## Side Effects

- No duplicate rows added for the same path.

## Errors

- Warning on stderr only; command still succeeds.

## Exit Code

- `0`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	wantWarn := "data dir already registered: " + req.DataDirPath
	if !strings.Contains(resp.Stderr, wantWarn) {
		t.Fatalf("stderr missing warning %q:\n%s", wantWarn, resp.Stderr)
	}

	reg, err := parseRegistry(resp.RegistryJSON)
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.DataDirs) != 1 {
		t.Fatalf("data_dirs count = %d, want 1", len(reg.DataDirs))
	}
	if reg.DataDirs[0].Note != "updated" {
		t.Fatalf("note = %q, want %q", reg.DataDirs[0].Note, "updated")
	}
}
```