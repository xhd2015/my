## Expected

- Exit code `0`.
- Stdout mentions writing to data dir, planned file updates, and import completion.
- Data dir contains `agents/main/agent/openclaw-agent.sqlite` with `xai:default` profile.
- `openclaw.json` default model is `xai/grok-4` and auth profile metadata is present.
- `agents/main/agent/models.json` includes an `xai` provider with `grok-4`, `grok-composer-2.5-fast`, and `grok-build` while preserving other providers.
- `agents/main/sessions/sessions.json` resets `agent:main:main` and dashboard sessions to `xai/grok-4` and clears model overrides.
- Session transcript files append an `xai/grok-4` model change when still pinned to compass.
- No container `~/.grok` copy; restart hint uses `my openclaw run --restart`.

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
	if !strings.Contains(resp.Stdout, "Writing to data dir:") {
		t.Fatalf("stdout missing data dir line:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "/home/node/.grok/auth.json") {
		t.Fatalf("stdout should not mention container auth for local import:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Will update:") {
		t.Fatalf("stdout missing change plan:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "openclaw-agent.sqlite") {
		t.Fatalf("stdout missing auth store plan:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "xai/grok-4") {
		t.Fatalf("stdout missing default model plan:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "models.json") {
		t.Fatalf("stdout missing models.json plan:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "sessions.json") {
		t.Fatalf("stdout missing sessions.json plan:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Import complete.") {
		t.Fatalf("stdout missing completion line:\n%s", resp.Stdout)
	}

	dbPath := filepath.Join(req.RunDataDir, "agents", "main", "agent", "openclaw-agent.sqlite")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open auth db: %v", err)
	}
	defer db.Close()

	var storeJSON string
	if err := db.QueryRow(`SELECT store_json FROM auth_profile_store WHERE store_key = 'primary'`).Scan(&storeJSON); err != nil {
		t.Fatalf("read auth profile store: %v", err)
	}
	if !strings.Contains(storeJSON, "xai:default") {
		t.Fatalf("auth store missing xai:default profile:\n%s", storeJSON)
	}
	if !strings.Contains(storeJSON, "fixture-grok-access-token") {
		t.Fatalf("auth store missing imported access token:\n%s", storeJSON)
	}

	cfgData, err := os.ReadFile(filepath.Join(req.RunDataDir, "openclaw.json"))
	if err != nil {
		t.Fatalf("read openclaw.json: %v", err)
	}
	var cfg map[string]any
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		t.Fatalf("parse openclaw.json: %v", err)
	}
	agents, _ := cfg["agents"].(map[string]any)
	defaults, _ := agents["defaults"].(map[string]any)
	if defaults["model"] != "xai/grok-4" {
		t.Fatalf("default model = %v, want xai/grok-4", defaults["model"])
	}
	auth, _ := cfg["auth"].(map[string]any)
	profiles, _ := auth["profiles"].(map[string]any)
	profile, _ := profiles["xai:default"].(map[string]any)
	if profile["mode"] != "oauth" {
		t.Fatalf("xai profile mode = %v, want oauth", profile["mode"])
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
	compass, _ := providers["compass"].(map[string]any)
	if compass["apiKey"] != "compass-test-key" {
		t.Fatalf("compass provider apiKey = %v, want compass-test-key", compass["apiKey"])
	}
	xai, _ := providers["xai"].(map[string]any)
	xaiModels, _ := xai["models"].([]any)
	if len(xaiModels) == 0 {
		t.Fatalf("xai provider has no models:\n%s", modelsData)
	}
	xaiIDs := make([]string, 0, len(xaiModels))
	for _, raw := range xaiModels {
		model, _ := raw.(map[string]any)
		xaiIDs = append(xaiIDs, model["id"].(string))
	}
	for _, want := range []string{"grok-4", "grok-composer-2.5-fast", "grok-build"} {
		if !strings.Contains(strings.Join(xaiIDs, ","), want) {
			t.Fatalf("xai models = %v, missing %s", xaiIDs, want)
		}
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
	if main["modelProvider"] != "xai" {
		t.Fatalf("main session modelProvider = %v, want xai", main["modelProvider"])
	}
	if main["model"] != "grok-4" {
		t.Fatalf("main session model = %v, want grok-4", main["model"])
	}
	for _, key := range []string{"providerOverride", "modelOverride", "modelOverrideSource", "liveModelSwitchPending"} {
		if _, ok := main[key]; ok {
			t.Fatalf("main session still has %s after import", key)
		}
	}

	dashboard, _ := sessions["agent:main:dashboard:fixture-dashboard"]
	if dashboard["modelProvider"] != "xai" || dashboard["model"] != "grok-4" {
		t.Fatalf("dashboard session model = %v/%v, want xai/grok-4", dashboard["modelProvider"], dashboard["model"])
	}

	mainTranscript, err := os.ReadFile(filepath.Join(req.RunDataDir, "agents", "main", "sessions", "fixture-main-session.jsonl"))
	if err != nil {
		t.Fatalf("read main transcript: %v", err)
	}
	if !strings.Contains(string(mainTranscript), `"provider":"xai"`) || !strings.Contains(string(mainTranscript), `"modelId":"grok-4"`) {
		t.Fatalf("main transcript missing xai/grok-4 model change:\n%s", mainTranscript)
	}
}
```