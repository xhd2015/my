package openclaw

import (
	"encoding/json"
	"fmt"
	"os"
)

func patchOpenClawConfigForXAI(configPath, email string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse openclaw.json: %w", err)
	}

	auth, _ := cfg["auth"].(map[string]any)
	if auth == nil {
		auth = map[string]any{}
		cfg["auth"] = auth
	}

	profiles, _ := auth["profiles"].(map[string]any)
	if profiles == nil {
		profiles = map[string]any{}
		auth["profiles"] = profiles
	}

	profile := map[string]any{
		"provider": "xai",
		"mode":     "oauth",
	}
	if email != "" {
		profile["email"] = email
	}
	profiles[xaiDefaultProfile] = profile

	order, _ := auth["order"].(map[string]any)
	if order == nil {
		order = map[string]any{}
		auth["order"] = order
	}
	order["xai"] = []any{xaiDefaultProfile}

	agents, _ := cfg["agents"].(map[string]any)
	if agents == nil {
		agents = map[string]any{}
		cfg["agents"] = agents
	}
	defaults, _ := agents["defaults"].(map[string]any)
	if defaults == nil {
		defaults = map[string]any{}
		agents["defaults"] = defaults
	}
	defaults["model"] = xaiDefaultModel

	models, _ := cfg["models"].(map[string]any)
	if models == nil {
		models = map[string]any{}
		cfg["models"] = models
	}
	providers, _ := models["providers"].(map[string]any)
	if providers == nil {
		providers = map[string]any{}
		models["providers"] = providers
	}
	providers["xai"] = xaiProviderEntry()

	encoded, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')
	return os.WriteFile(configPath, encoded, 0o600)
}