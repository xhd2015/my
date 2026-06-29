## Expected

- Exit code `0`.
- `auth.json` exists with top-level `xai` oauth credential.
- `access`, `refresh`, and `expires` match grok fixture (expires as Unix ms).
- Stdout mentions `xai` and destination path; no token values printed.

## Side Effects

- File mode `0600` on `auth.json`.

## Errors

- None.

## Exit Code

- `0`

```go
import "math"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertNoSecretsInOutput(t, resp.Stdout+resp.Stderr)
	assertStdoutMentionsImport(t, resp.Stdout, resp.AuthPath)

	m := parseAuth(t, resp.AuthJSON)
	if len(m) != 1 {
		t.Fatalf("auth keys = %d, want 1 (xai only)", len(m))
	}
	raw, ok := m["xai"]
	if !ok {
		t.Fatal("missing xai key in auth.json")
	}
	o := xaiOAuth(t, raw)
	if o["type"] != "oauth" {
		t.Fatalf("xai.type = %v, want oauth", o["type"])
	}
	if o["access"] != "fixture-grok-access-token" {
		t.Fatalf("xai.access mismatch")
	}
	if o["refresh"] != "fixture-grok-refresh-token" {
		t.Fatalf("xai.refresh mismatch")
	}
	exp, ok := o["expires"].(float64)
	if !ok {
		t.Fatalf("xai.expires type %T", o["expires"])
	}
	wantMs := float64(4102444800000) // 2099-12-31T00:00:00Z
	if math.Abs(exp-wantMs) > 1 {
		t.Fatalf("xai.expires = %v, want ~%v", exp, wantMs)
	}
	if mode := authFileMode(t, resp.AuthPath); mode != 0o600 {
		t.Fatalf("auth mode = %o, want 0600", mode)
	}
}```
