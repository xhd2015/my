## Expected

- Exit code `0`.
- Stdout prints `Using data dir:` for the auto-selected path, then normal status output.

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
	paths, err := registeredPaths(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) < 1 {
		t.Fatal("registered paths empty")
	}
	autoDir := paths[0]
	out := resp.Stdout
	if !strings.Contains(out, "Using data dir: "+autoDir) {
		t.Fatalf("stdout missing auto-select message for %q:\n%s", autoDir, out)
	}
	for _, want := range []string{
		"Gateway running locally",
		"Data dir: " + autoDir,
		"Port: 18789",
		"http://127.0.0.1:18789/",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
}

```