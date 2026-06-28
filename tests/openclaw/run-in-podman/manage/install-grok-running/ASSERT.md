## Expected

- Exit code `0`.
- Stdout mentions installing Grok CLI in the container.
- Stderr previews the concrete `podman exec` install command.
- Podman `exec` runs the official install script as root/node.

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
	if !strings.Contains(resp.Stdout, "Installing Grok CLI in container openclaw-gateway") {
		t.Fatalf("stdout missing install message:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "/home/node/.grok/bin/grok") {
		t.Fatalf("stdout missing grok run hint:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, "$ podman exec") {
		t.Fatalf("stderr missing command preview:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "x.ai/cli/install.sh") {
		t.Fatalf("stderr missing install script URL:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "Acquire::Check-Valid-Until=false") {
		t.Fatalf("stderr missing apt clock-skew workaround:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "--user root") {
		t.Fatalf("stderr missing root user for curl bootstrap:\n%s", resp.Stderr)
	}

	want := "podman exec --user root openclaw-gateway sh -c"
	if !podmanCallsContain(resp.PodmanCalls, want) {
		t.Fatalf("missing podman exec install call %q in:\n%v", want, resp.PodmanCalls)
	}
	if !podmanCallsContain(resp.PodmanCalls, "x.ai/cli/install.sh") {
		t.Fatalf("missing install script in podman calls:\n%v", resp.PodmanCalls)
	}
}
```