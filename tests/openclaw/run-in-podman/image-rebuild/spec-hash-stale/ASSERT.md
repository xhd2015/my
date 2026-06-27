## Expected

- Exit code `0`.
- Podman log contains `podman build` due to hash mismatch.

## Side Effects

- Registry `image.spec_hash` updated after build (implementation detail).

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
	if podmanCallsCount(resp.PodmanCalls, "podman build") < 1 {
		t.Fatalf("expected podman build for stale hash, got:\n%v", resp.PodmanCalls)
	}
}
```