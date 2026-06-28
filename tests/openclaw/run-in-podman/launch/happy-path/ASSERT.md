## Expected

- Exit code `0`.
- Stdout contains dashboard URL `http://127.0.0.1:18789/`.
- Stdout contains `podman logs -f openclaw-gateway`.
- Podman `run` includes data-dir mount, workspace mount, token env, `--bind lan`, and image `my-openclaw:local`.
- Stderr previews `podman run` command.

## Side Effects

- No `podman stop` or `podman rm` before start when container is not already running.

## Errors

- None.

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
	if !strings.Contains(resp.Stdout, "http://127.0.0.1:18789/") {
		t.Fatalf("stdout missing dashboard URL:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "podman logs -f openclaw-gateway") {
		t.Fatalf("stdout missing logs hint:\n%s", resp.Stdout)
	}

	if !strings.Contains(resp.Stderr, "$ podman run") {
		t.Fatalf("stderr missing command preview:\n%s", resp.Stderr)
	}

	var runLine string
	for _, call := range resp.PodmanCalls {
		if strings.Contains(call, "podman run") {
			runLine = call
			break
		}
	}
	if runLine == "" {
		t.Fatalf("missing podman run call:\n%v", resp.PodmanCalls)
	}
	for _, want := range []string{
		req.RunDataDir + ":/home/node/.openclaw",
		req.RunDataDir + "/workspace:/home/node/.openclaw/workspace",
		"OPENCLAW_GATEWAY_TOKEN=json-gateway-token",
		"--name openclaw-gateway",
		"-p 18789:18789",
		"gateway --bind lan",
		"my-openclaw:local",
	} {
		if !strings.Contains(runLine, want) {
			t.Fatalf("podman run missing %q in:\n%s", want, runLine)
		}
	}
	if podmanCallsContain(resp.PodmanCalls, "podman stop openclaw-gateway") {
		t.Fatalf("unexpected podman stop on fresh start:\n%v", resp.PodmanCalls)
	}
}
```