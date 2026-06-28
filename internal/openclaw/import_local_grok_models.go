package openclaw

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	xaiBaseURL                = "https://api.x.ai/v1"
	xaiModelID                = "grok-4"
	xaiModelName              = "Grok 4"
	xaiComposerModelID        = "grok-composer-2.5-fast"
	xaiComposerModelName      = "Composer 2.5 Fast"
	xaiBuildModelID           = "grok-build"
	xaiBuildModelName         = "Grok Build"
	xaiContextWindow          = 131072
	xaiBuildContextWindow     = 256000
	xaiMaxTokens              = 8192
	xaiBuildMaxTokens         = 16384
	mainAgentSessionKey       = "agent:main:main"
	dashboardSessionKeyPrefix = "agent:main:dashboard:"
)

type xaiModelSpec struct {
	id            string
	name          string
	contextWindow int
	maxTokens     int
	reasoning     bool
}

func xaiCatalogModels() []xaiModelSpec {
	return []xaiModelSpec{
		{id: xaiModelID, name: xaiModelName, contextWindow: xaiContextWindow, maxTokens: xaiMaxTokens, reasoning: false},
		{id: xaiComposerModelID, name: xaiComposerModelName, contextWindow: xaiBuildContextWindow, maxTokens: xaiBuildMaxTokens, reasoning: false},
		{id: xaiBuildModelID, name: xaiBuildModelName, contextWindow: xaiBuildContextWindow, maxTokens: xaiBuildMaxTokens, reasoning: false},
	}
}

func agentModelsJSONPath(dataDir string) string {
	return filepath.Join(dataDir, "agents", "main", "agent", "models.json")
}

func mainSessionsJSONPath(dataDir string) string {
	return filepath.Join(dataDir, "agents", "main", "sessions", "sessions.json")
}

func xaiModelDefinition(spec xaiModelSpec) map[string]any {
	return map[string]any{
		"id":        spec.id,
		"name":      spec.name,
		"reasoning": spec.reasoning,
		"input":     []any{"text"},
		"cost": map[string]any{
			"input":      0,
			"output":     0,
			"cacheRead":  0,
			"cacheWrite": 0,
		},
		"contextWindow": spec.contextWindow,
		"maxTokens":     spec.maxTokens,
	}
}

func xaiProviderEntry() map[string]any {
	models := make([]any, 0, len(xaiCatalogModels()))
	for _, spec := range xaiCatalogModels() {
		models = append(models, xaiModelDefinition(spec))
	}
	return map[string]any{
		"baseUrl": xaiBaseURL,
		"api":     "openai-completions",
		"models":  models,
	}
}

func syncAgentModelsJSONForXAI(dataDir string) error {
	path := agentModelsJSONPath(dataDir)
	providers := map[string]any{}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		var existing map[string]any
		if err := json.Unmarshal(data, &existing); err != nil {
			return fmt.Errorf("parse models.json: %w", err)
		}
		if raw, ok := existing["providers"].(map[string]any); ok {
			providers = raw
		}
	}

	providers["xai"] = xaiProviderEntry()
	payload := map[string]any{"providers": providers}

	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, encoded, 0o600)
}

func shouldResetSessionKeyForXAI(key string) bool {
	return key == mainAgentSessionKey || strings.HasPrefix(key, dashboardSessionKeyPrefix)
}

func applyXAISessionModel(entry map[string]any) {
	entry["modelProvider"] = "xai"
	entry["model"] = xaiModelID
	delete(entry, "providerOverride")
	delete(entry, "modelOverride")
	delete(entry, "modelOverrideSource")
	delete(entry, "liveModelSwitchPending")
}

func resetSessionsForXAI(dataDir string) error {
	path := mainSessionsJSONPath(dataDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var sessions map[string]map[string]any
	if err := json.Unmarshal(data, &sessions); err != nil {
		return fmt.Errorf("parse sessions.json: %w", err)
	}

	sessionFiles := make(map[string]struct{})
	for key, session := range sessions {
		if !shouldResetSessionKeyForXAI(key) {
			continue
		}
		applyXAISessionModel(session)
		sessions[key] = session
		if sessionFile, ok := session["sessionFile"].(string); ok && sessionFile != "" {
			sessionFiles[sessionFile] = struct{}{}
		}
	}

	encoded, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')
	if err := os.WriteFile(path, encoded, 0o600); err != nil {
		return err
	}

	for sessionFile := range sessionFiles {
		resolved := sessionFile
		if !filepath.IsAbs(resolved) {
			resolved = filepath.Join(dataDir, resolved)
		}
		if err := patchSessionTranscriptModelForXAI(resolved); err != nil {
			return err
		}
	}
	return nil
}

func resetMainSessionModelForXAI(dataDir string) error {
	return resetSessionsForXAI(dataDir)
}

func patchSessionTranscriptModelForXAI(sessionFile string) error {
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return nil
	}

	lastID := ""
	provider := ""
	modelID := ""
	for _, line := range lines {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if id, ok := entry["id"].(string); ok && id != "" {
			lastID = id
		}
		switch entry["type"] {
		case "model_change":
			if p, ok := entry["provider"].(string); ok {
				provider = p
			}
			if m, ok := entry["modelId"].(string); ok {
				modelID = m
			}
		case "custom":
			if entry["customType"] != "model-snapshot" {
				continue
			}
			snapshot, _ := entry["data"].(map[string]any)
			if p, ok := snapshot["provider"].(string); ok {
				provider = p
			}
			if m, ok := snapshot["modelId"].(string); ok {
				modelID = m
			}
		}
	}

	if provider == "xai" && modelID == xaiModelID {
		return nil
	}

	changeID, err := randomSessionEntryID()
	if err != nil {
		return err
	}
	var parentID any = lastID
	if lastID == "" {
		parentID = nil
	}
	changeEntry := map[string]any{
		"type":      "model_change",
		"id":        changeID,
		"parentId":  parentID,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"provider":  "xai",
		"modelId":   xaiModelID,
	}
	encoded, err := json.Marshal(changeEntry)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(sessionFile, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		return err
	}
	return nil
}

func randomSessionEntryID() (string, error) {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf[:]), nil
}

func sessionTranscriptActiveModel(path string) (provider, modelID string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}
		switch entry["type"] {
		case "model_change":
			if p, ok := entry["provider"].(string); ok {
				provider = p
			}
			if m, ok := entry["modelId"].(string); ok {
				modelID = m
			}
		case "custom":
			if entry["customType"] != "model-snapshot" {
				continue
			}
			snapshot, _ := entry["data"].(map[string]any)
			if p, ok := snapshot["provider"].(string); ok {
				provider = p
			}
			if m, ok := snapshot["modelId"].(string); ok {
				modelID = m
			}
		}
	}
	return provider, modelID, scanner.Err()
}