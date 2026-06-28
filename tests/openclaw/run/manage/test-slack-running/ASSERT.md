## Expected

- Exit code `0`.
- Stdout confirms test message sent to configured `testTarget`.
- Slack stub used (`MY_OPENCLAW_SLACK_STUB=1`); no real HTTP required.

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
	if !strings.Contains(resp.Stdout, "Sending Slack test message to C01234567 (C01234567)") {
		t.Fatalf("stdout missing send line:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Test message sent.") {
		t.Fatalf("stdout missing confirmation:\n%s", resp.Stdout)
	}
}
```