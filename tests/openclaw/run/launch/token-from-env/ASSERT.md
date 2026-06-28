## Expected

- Exit code `0`.
- Stub env `OPENCLAW_GATEWAY_TOKEN=env-gateway-token`.

## Side Effects

- Token resolved from `.env` file.

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
	if !openclawCallsContain(resp.OpenClawCalls, "OPENCLAW_GATEWAY_TOKEN=env-gateway-token") {
		t.Fatalf("openclaw missing env token:\n%v", resp.OpenClawCalls)
	}
}
```
