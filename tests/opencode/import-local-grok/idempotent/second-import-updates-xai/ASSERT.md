## Expected

- Exit code `0` after second import.
- `xai.access` and `xai.refresh` match v2 fixture tokens.
- `other` provider entry still present.

## Exit Code

- `0`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d after second import", resp.ExitCode)
	}
	m := parseAuth(t, resp.AuthJSON)
	if _, ok := m["other"]; !ok {
		t.Fatal("other provider removed on re-import")
	}
	o := xaiOAuth(t, m["xai"])
	if o["access"] != "fixture-grok-access-token-2" {
		t.Fatalf("xai.access = %v, want v2 token", o["access"])
	}
	if o["refresh"] != "fixture-grok-refresh-token-2" {
		t.Fatalf("xai.refresh = %v, want v2 token", o["refresh"])
	}
}```
