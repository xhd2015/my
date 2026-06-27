## Expected

- Exit code `0`.
- Registry file created at `$MY_CONFIG_DIR/openclaw.json`.
- Stored `path` equals the absolute data directory path passed to `add`.

## Side Effects

- `data_dirs` contains exactly one entry.
- `added_at` is a non-empty RFC3339 timestamp.

## Errors

- None.

## Exit Code

- `0`

```go
import (
	"path/filepath"
	"strings"
)

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
	entry := reg.DataDirs[0]
	if entry.Path != req.DataDirPath {
		t.Fatalf("stored path = %q, want %q", entry.Path, req.DataDirPath)
	}
	if !filepath.IsAbs(entry.Path) {
		t.Fatalf("stored path is not absolute: %q", entry.Path)
	}
	if strings.TrimSpace(entry.AddedAt) == "" {
		t.Fatal("added_at is empty")
	}
}
```