## Expected

- Exit code `0`.
- Podman log does **not** contain `podman build`.
- Podman log contains `podman run`.

## Side Effects

- Existing image reused.

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
	if podmanCallsCount(resp.PodmanCalls, "podman build") != 0 {
		t.Fatalf("podman build should be skipped, got:\n%v", resp.PodmanCalls)
	}
	if podmanCallsCount(resp.PodmanCalls, "podman run") < 1 {
		t.Fatalf("expected podman run, got:\n%v", resp.PodmanCalls)
	}
}
```