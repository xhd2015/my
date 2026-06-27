## Expected

- Exit code `1`.
- Stderr indicates the data directory is missing or invalid.

## Side Effects

- No `podman run` invocation recorded.

## Errors

- Data directory not found.

## Exit Code

- `1`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if podmanCallsContain(resp.PodmanCalls, "podman run") {
		t.Fatalf("podman run should not be called:\n%v", resp.PodmanCalls)
	}
}
```