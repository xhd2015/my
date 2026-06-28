## Expected

- Exit code `0`.
- Podman `exec` called with container name and verbatim command args.
- Stderr previews `podman exec` command.

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
	if !strings.Contains(resp.Stderr, "$ podman exec") {
		t.Fatalf("stderr missing command preview:\n%s", resp.Stderr)
	}

	want := "podman exec openclaw-gateway openclaw models auth login --provider xai --method oauth"
	if !podmanCallsContain(resp.PodmanCalls, want) {
		t.Fatalf("missing podman exec call %q in:\n%v", want, resp.PodmanCalls)
	}
}
```