## Expected

- Exit code `1`.
- Stderr reports multiple running gateways and lists running paths.

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
		"multiple local gateways running",
		"pass --data-dir explicitly",
		"Running data dirs:",
	}, paths...) {
		if !strings.Contains(errOut, want) {
			t.Fatalf("stderr missing %q:\n%s", want, errOut)
		}
	}
}

```