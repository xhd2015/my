## Expected

- Exit code `0`.
- Podman `run` uses `OPENCLAW_GATEWAY_TOKEN=json-wins-token`, not the `.env` value.

## Side Effects

- Json token precedence enforced.

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
		t.Fatalf("exit code = %d, want 0", resp.ExitCode)
	}
	if !podmanCallsContain(resp.PodmanCalls, "OPENCLAW_GATEWAY_TOKEN=json-wins-token") {
		t.Fatalf("podman run missing json-preferred token:\n%v", resp.PodmanCalls)
	}
	if podmanCallsContain(resp.PodmanCalls, "OPENCLAW_GATEWAY_TOKEN=env-loses-token") {
		t.Fatalf("podman run should not use env token when json token exists:\n%v", resp.PodmanCalls)
	}
}
```