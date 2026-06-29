## Expected

- Exit code `0`.
- Both `deepseek` and `xai` keys present.
- `deepseek` api key unchanged.

## Side Effects

- `xai` oauth populated from grok fixture.

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
	m := parseAuth(t, resp.AuthJSON)
	if _, ok := m["deepseek"]; !ok {
		t.Fatal("deepseek missing after merge")
	}
	ds := xaiOAuth(t, m["deepseek"])
	if ds["type"] != "api" || ds["key"] != "deepseek-secret-key" {
		t.Fatalf("deepseek entry altered: %v", ds)
	}
	if _, ok := m["xai"]; !ok {
		t.Fatal("xai missing after merge")
	}
	assertNoSecretsInOutput(t, resp.Stdout+resp.Stderr)
}```
