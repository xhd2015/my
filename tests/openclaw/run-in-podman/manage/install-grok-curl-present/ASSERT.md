## Expected

- Exit code `0`.
- Stderr previews direct curl install without apt or `--user root`.

## Exit Code

- `0`

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0", resp.ExitCode)
	}
	if strings.Contains(resp.Stderr, "apt-get") {
		t.Fatalf("stderr should not run apt when curl exists:\n%s", resp.Stderr)
	}
	if strings.Contains(resp.Stderr, "--user root") {
		t.Fatalf("stderr should not use root when curl exists:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "x.ai/cli/install.sh") {
		t.Fatalf("stderr missing install script URL:\n%s", resp.Stderr)
	}
	want := "podman exec openclaw-gateway sh -c"
	if !podmanCallsContain(resp.PodmanCalls, want) {
		t.Fatalf("missing direct install call %q in:\n%v", want, resp.PodmanCalls)
	}
}
```