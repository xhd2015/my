## Expected

- Exit code `0`.
- Podman log contains `podman machine info`.
- Podman log contains `podman machine start` before `podman run`.

## Side Effects

- Machine started automatically on macOS.

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
	if !podmanCallsContain(resp.PodmanCalls, "podman machine info") {
		t.Fatalf("expected podman machine info, got:\n%v", resp.PodmanCalls)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman machine start") {
		t.Fatalf("expected podman machine start, got:\n%v", resp.PodmanCalls)
	}
	infoIdx, startIdx, runIdx := -1, -1, -1
	for i, call := range resp.PodmanCalls {
		switch {
		case call == "podman machine info":
			infoIdx = i
		case call == "podman machine start":
			startIdx = i
		case containsRun(call):
			runIdx = i
		}
	}
	if infoIdx < 0 || startIdx < 0 || runIdx < 0 {
		t.Fatalf("missing ordered machine lifecycle calls:\n%v", resp.PodmanCalls)
	}
	if !(infoIdx < startIdx && startIdx < runIdx) {
		t.Fatalf("want info < start < run order, got:\n%v", resp.PodmanCalls)
	}
}

func containsRun(call string) bool {
	return len(call) >= len("podman run") && call[:len("podman run")] == "podman run"
}
```