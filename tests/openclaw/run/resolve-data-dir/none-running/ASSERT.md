## Expected

- Exit code `1`.
- Stderr requires `--data-dir` (running) and lists both registered paths.

## Exit Code

- `1`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	paths, err := registeredPaths(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 {
		t.Fatalf("registered paths = %v, want 2", paths)
	}
	errOut := resp.Stderr
	for _, want := range append([]string{
		"requires --data-dir",
		"(running)",
		"Registered data dirs:",
	}, paths...) {
		if !strings.Contains(errOut, want) {
			t.Fatalf("stderr missing %q:\n%s", want, errOut)
		}
	}
}

```