package openclaw

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestSyncAgentModelsJSONForXAI(t *testing.T) {
	dir := t.TempDir()

	if err := syncAgentModelsJSONForXAI(dir); err != nil {
		t.Fatalf("syncAgentModelsJSONForXAI() error = %v", err)
	}

	data, err := os.ReadFile(agentModelsJSONPath(dir))
	if err != nil {
		t.Fatalf("read models.json: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("parse models.json: %v", err)
	}
	providers, _ := parsed["providers"].(map[string]any)
	xai, _ := providers["xai"].(map[string]any)
	if xai["baseUrl"] != xaiBaseURL {
		t.Fatalf("baseUrl = %v, want %s", xai["baseUrl"], xaiBaseURL)
	}
	models, _ := xai["models"].([]any)
	if len(models) != 3 {
		t.Fatalf("models len = %d, want 3", len(models))
	}
	ids := make([]string, 0, len(models))
	for _, raw := range models {
		model, _ := raw.(map[string]any)
		ids = append(ids, model["id"].(string))
	}
	want := []string{xaiModelID, xaiComposerModelID, xaiBuildModelID}
	for i, id := range want {
		if ids[i] != id {
			t.Fatalf("model ids = %v, want %v", ids, want)
		}
	}
}

func TestSyncAgentModelsJSONForXAIPreservesExistingProviders(t *testing.T) {
	dir := t.TempDir()
	modelsDir := filepath.Join(dir, "agents", "main", "agent")
	if err := os.MkdirAll(modelsDir, 0o700); err != nil {
		t.Fatal(err)
	}
	existing := `{
  "providers": {
    "compass": {
      "baseUrl": "http://example.com/v1",
      "apiKey": "keep-me",
      "api": "openai-completions",
      "models": [{"id": "gpt-5-mini-2025-08-07", "name": "GPT-5 Mini"}]
    }
  }
}`
	if err := os.WriteFile(agentModelsJSONPath(dir), []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := syncAgentModelsJSONForXAI(dir); err != nil {
		t.Fatalf("syncAgentModelsJSONForXAI() error = %v", err)
	}

	data, err := os.ReadFile(agentModelsJSONPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}
	providers, _ := parsed["providers"].(map[string]any)
	compass, _ := providers["compass"].(map[string]any)
	if compass["apiKey"] != "keep-me" {
		t.Fatalf("compass apiKey = %v, want keep-me", compass["apiKey"])
	}
	if providers["xai"] == nil {
		t.Fatal("expected xai provider to be added")
	}
}

func TestResetSessionsForXAIDashboard(t *testing.T) {
	dir := t.TempDir()
	sessionsDir := filepath.Join(dir, "agents", "main", "sessions")
	if err := os.MkdirAll(sessionsDir, 0o700); err != nil {
		t.Fatal(err)
	}
	transcript := filepath.Join(sessionsDir, "dash.jsonl")
	if err := os.WriteFile(transcript, []byte(`{"type":"session","version":3,"id":"dash","timestamp":"2026-06-27T00:00:00.000Z"}
{"type":"model_change","id":"dash-model","parentId":null,"timestamp":"2026-06-27T00:00:01.000Z","provider":"compass","modelId":"gpt-5-mini-2025-08-07"}
`), 0o600); err != nil {
		t.Fatal(err)
	}
	existing := fmt.Sprintf(`{
  "agent:main:dashboard:abc": {
    "sessionId": "dash",
    "sessionFile": %q,
    "modelProvider": "compass",
    "model": "gpt-5-mini-2025-08-07"
  }
}`, transcript)
	if err := os.WriteFile(mainSessionsJSONPath(dir), []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := resetSessionsForXAI(dir); err != nil {
		t.Fatalf("resetSessionsForXAI() error = %v", err)
	}

	data, err := os.ReadFile(mainSessionsJSONPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	var sessions map[string]map[string]any
	if err := json.Unmarshal(data, &sessions); err != nil {
		t.Fatal(err)
	}
	dash, _ := sessions["agent:main:dashboard:abc"]
	if dash["modelProvider"] != "xai" || dash["model"] != xaiModelID {
		t.Fatalf("dashboard session model = %v/%v, want xai/%s", dash["modelProvider"], dash["model"], xaiModelID)
	}

	provider, model, err := sessionTranscriptActiveModel(transcript)
	if err != nil {
		t.Fatal(err)
	}
	if provider != "xai" || model != xaiModelID {
		t.Fatalf("transcript active model = %s/%s, want xai/%s", provider, model, xaiModelID)
	}
}

func TestResetMainSessionModelForXAI(t *testing.T) {
	dir := t.TempDir()
	sessionsDir := filepath.Join(dir, "agents", "main", "sessions")
	if err := os.MkdirAll(sessionsDir, 0o700); err != nil {
		t.Fatal(err)
	}
	existing := `{
  "agent:main:main": {
    "sessionId": "sess-1",
    "modelProvider": "compass",
    "model": "gpt-5-mini-2025-08-07",
    "providerOverride": "compass",
    "modelOverride": "gpt-5-mini-2025-08-07",
    "modelOverrideSource": "user",
    "liveModelSwitchPending": true
  },
  "agent:main:slack:channel:abc": {
    "modelProvider": "compass",
    "model": "gpt-5-mini-2025-08-07"
  }
}`
	if err := os.WriteFile(mainSessionsJSONPath(dir), []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := resetMainSessionModelForXAI(dir); err != nil {
		t.Fatalf("resetMainSessionModelForXAI() error = %v", err)
	}

	data, err := os.ReadFile(mainSessionsJSONPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	var sessions map[string]map[string]any
	if err := json.Unmarshal(data, &sessions); err != nil {
		t.Fatal(err)
	}
	main, _ := sessions[mainAgentSessionKey]
	if main["modelProvider"] != "xai" {
		t.Fatalf("modelProvider = %v, want xai", main["modelProvider"])
	}
	if main["model"] != xaiModelID {
		t.Fatalf("model = %v, want %s", main["model"], xaiModelID)
	}
	for _, key := range []string{"providerOverride", "modelOverride", "modelOverrideSource", "liveModelSwitchPending"} {
		if _, ok := main[key]; ok {
			t.Fatalf("expected %s to be cleared", key)
		}
	}

	slack, _ := sessions["agent:main:slack:channel:abc"]
	if slack["modelProvider"] != "compass" {
		t.Fatalf("slack modelProvider = %v, want compass", slack["modelProvider"])
	}
}

func TestResetMainSessionModelForXAISkipsMissingFile(t *testing.T) {
	if err := resetMainSessionModelForXAI(t.TempDir()); err != nil {
		t.Fatalf("resetMainSessionModelForXAI() error = %v", err)
	}
}