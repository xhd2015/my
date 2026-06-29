# `my opencode --import-local-grok` — Doc-Style Test Tree

Imports local Grok OIDC credentials from `~/.grok/auth.json` (or `MY_GROK_AUTH_PATH`) into OpenCode's
`auth.json` under the OpenCode data directory, merging provider `xai` without disturbing other providers.

## Version

0.0.2

# DSN (Domain Specific Notion)

- **CLI binary** — `my` built from `cmd/my/` at module root (`github.com/xhd2015/my`).
- **Grok source** — JSON map of OIDC entries; default path `~/.grok/auth.json`, override `MY_GROK_AUTH_PATH`.
- **OpenCode destination** — `auth.json` in OpenCode global data dir; tests override with `MY_OPENCODE_DATA_DIR`
  pointing at an isolated directory (writer must create `auth.json` with mode `0600`).
- **xAI credential shape** — top-level key `xai` with `type: oauth`, `access`, `refresh`, `expires` (Unix ms).
- **Import operation** — read grok entry (`auth_mode: oidc`), validate tokens, merge into existing auth map,
  write file; stdout reports destination path and provider `xai` without secrets.

## How to Run

```sh
doctest vet ./tests/opencode/import-local-grok
doctest test -v ./tests/opencode/import-local-grok
```

## Decision Tree

```
outcome
├── happy-path/
│   └── empty-opencode-auth
├── merge/
│   └── preserves-deepseek
├── errors/
│   ├── missing-grok-file
│   ├── no-oidc-entry
│   ├── missing-access
│   └── missing-refresh
└── idempotent/
    └── second-import-updates-xai
```

## Test Leaf Index

| Leaf | Path |
|------|------|
| happy empty auth | `happy-path/empty-opencode-auth` |
| merge deepseek | `merge/preserves-deepseek` |
| missing grok file | `errors/missing-grok-file` |
| no OIDC entry | `errors/no-oidc-entry` |
| missing access | `errors/missing-access` |
| missing refresh | `errors/missing-refresh` |
| idempotent update | `idempotent/second-import-updates-xai` |

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type Request struct {
	CommandDir       string
	BinPath          string
	GrokAuthPath     string
	OpenCodeDataDir  string
	PreseedAuthJSON  []byte
	SecondGrokPath   string
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
	AuthJSON []byte
	AuthPath string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	args := []string{"opencode", "--import-local-grok"}
	if len(req.PreseedAuthJSON) > 0 {
		if err := os.MkdirAll(req.OpenCodeDataDir, 0o755); err != nil {
			return nil, err
		}
		authPath := filepath.Join(req.OpenCodeDataDir, "auth.json")
		if err := os.WriteFile(authPath, req.PreseedAuthJSON, 0o600); err != nil {
			return nil, err
		}
	}

	runOnce := func(grokPath string) (stdout, stderr string, exitCode int, runErr error) {
		envRun := append(os.Environ(), "MY_OPENCODE_DATA_DIR="+req.OpenCodeDataDir)
		if grokPath != "" {
			envRun = append(envRun, "MY_GROK_AUTH_PATH="+grokPath)
		}
		cmd := exec.Command(req.BinPath, args...)
		cmd.Dir = req.CommandDir
		cmd.Env = envRun
		var outB, errB strings.Builder
		cmd.Stdout = &outB
		cmd.Stderr = &errB
		runErr = cmd.Run()
		exitCode = 0
		if runErr != nil {
			if ee, ok := runErr.(*exec.ExitError); ok {
				exitCode = ee.ExitCode()
			} else {
				return "", "", 0, runErr
			}
		}
		return outB.String(), errB.String(), exitCode, nil
	}

	authPath := filepath.Join(req.OpenCodeDataDir, "auth.json")
	stdout, stderr, exitCode, err := runOnce(req.GrokAuthPath)
	if err != nil {
		return nil, err
	}
	if req.SecondGrokPath != "" {
		stdout, stderr, exitCode, err = runOnce(req.SecondGrokPath)
		if err != nil {
			return nil, err
		}
	}

	var authJSON []byte
	if data, readErr := os.ReadFile(authPath); readErr == nil {
		authJSON = data
	}

	return &Response{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
		AuthJSON: authJSON,
		AuthPath: authPath,
	}, nil
}

func parseAuth(t *testing.T, raw []byte) map[string]json.RawMessage {
	t.Helper()
	if len(raw) == 0 {
		return nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("parse auth.json: %v\n%s", err, raw)
	}
	return m
}

func xaiOAuth(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	var o map[string]any
	if err := json.Unmarshal(raw, &o); err != nil {
		t.Fatalf("parse xai oauth: %v", err)
	}
	return o
}

func authFileMode(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat auth: %v", err)
	}
	return info.Mode().Perm()
}

func assertNoSecretsInOutput(t *testing.T, out string) {
	t.Helper()
	for _, secret := range []string{
		"fixture-grok-access-token",
		"fixture-grok-refresh-token",
		"fixture-grok-access-token-2",
		"fixture-grok-refresh-token-2",
		"deepseek-secret-key",
	} {
		if strings.Contains(out, secret) {
			t.Fatalf("output must not contain secret %q", secret)
		}
	}
}

func copyFixture(t *testing.T, rel string) string {
	t.Helper()
	src := filepath.Join(DOCTEST_ROOT, rel)
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read fixture %s: %v", rel, err)
	}
	dst := filepath.Join(t.TempDir(), filepath.Base(rel))
	if err := os.WriteFile(dst, data, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return dst
}

func writeTempJSON(t *testing.T, name string, content []byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, content, 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

func mustRun(t *testing.T, req *Request) *Response {
	t.Helper()
	resp, err := Run(t, req)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	return resp
}

func assertStdoutMentionsImport(t *testing.T, stdout, authPath string) {
	t.Helper()
	lower := strings.ToLower(stdout)
	if !strings.Contains(lower, "xai") {
		t.Fatalf("stdout should mention provider xai: %q", stdout)
	}
	if !strings.Contains(stdout, authPath) && !strings.Contains(lower, "auth.json") {
		t.Fatalf("stdout should mention destination auth path: %q", stdout)
	}
}

func assertErrorResponse(t *testing.T, resp Response) {
	t.Helper()
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0 stdout=%q stderr=%q", resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + resp.Stderr
	if strings.TrimSpace(combined) == "" {
		t.Fatal("expected error message on stdout or stderr")
	}
}
```