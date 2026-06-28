## Expected

- Exit code `0`.
- Stdout writes to the mounted data dir from the running container.
- Import updates land in that auto-detected data dir.
- `agents/main/agent/models.json` includes an `xai` provider with `grok-4`.
- `agents/main/sessions/sessions.json` resets `agent:main:main` to `xai/grok-4`.

## Exit Code

- `0`

```go
import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Writing to data dir: "+req.PodmanContainerDataDir) {
		t.Fatalf("stdout missing auto-detected data dir:\n%s", resp.Stdout)
	}

	dbPath := filepath.Join(req.RunDataDir, "agents", "main", "agent", "openclaw-agent.sqlite")
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("auth db missing: %v", err)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open auth db: %v", err)
	}
	defer db.Close()

	var storeJSON string
	if err := db.QueryRow(`SELECT store_json FROM auth_profile_store WHERE store_key = 'primary'`).Scan(&storeJSON); err != nil {
		t.Fatalf("read auth profile store: %v", err)
	}
	if !strings.Contains(storeJSON, "fixture-grok-access-token") {
		t.Fatalf("auth store missing imported access token:\n%s", storeJSON)
	}
	if !strings.Contains(resp.Stdout, "/home/node/.grok/auth.json") {
		t.Fatalf("stdout missing container auth path:\n%s", resp.Stdout)
	}
	if !podmanCallsContain(resp.PodmanCalls, "podman cp") {
		t.Fatalf("missing podman cp to container auth.json:\n%v", resp.PodmanCalls)
	}
	if !podmanCallsContain(resp.PodmanCalls, ":/home/node/.grok/auth.json") {
		t.Fatalf("missing container auth destination:\n%v", resp.PodmanCalls)
	}

	modelsData, err := os.ReadFile(filepath.Join(req.RunDataDir, "agents", "main", "agent", "models.json"))
	if err != nil {
		t.Fatalf("read models.json: %v", err)
	}
	var models map[string]any
	if err := json.Unmarshal(modelsData, &models); err != nil {
		t.Fatalf("parse models.json: %v", err)
	}
	providers, _ := models["providers"].(map[string]any)
	if providers["xai"] == nil {
		t.Fatalf("models.json missing xai provider:\n%s", modelsData)
	}

	sessionsData, err := os.ReadFile(filepath.Join(req.RunDataDir, "agents", "main", "sessions", "sessions.json"))
	if err != nil {
		t.Fatalf("read sessions.json: %v", err)
	}
	var sessions map[string]map[string]any
	if err := json.Unmarshal(sessionsData, &sessions); err != nil {
		t.Fatalf("parse sessions.json: %v", err)
	}
	main, _ := sessions["agent:main:main"]
	if main["modelProvider"] != "xai" || main["model"] != "grok-4" {
		t.Fatalf("main session model = %v/%v, want xai/grok-4", main["modelProvider"], main["model"])
	}
}
```