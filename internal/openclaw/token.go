package openclaw

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type openclawConfig struct {
	Gateway struct {
		Auth struct {
			Token string `json:"token"`
		} `json:"auth"`
	} `json:"gateway"`
}

func resolveToken(dataDir string) (string, error) {
	configPath := filepath.Join(dataDir, "openclaw.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("openclaw.json not found in %s", dataDir)
		}
		return "", err
	}

	var cfg openclawConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("parse openclaw.json: %w", err)
	}
	if token := strings.TrimSpace(cfg.Gateway.Auth.Token); token != "" {
		return token, nil
	}

	envPath := filepath.Join(dataDir, ".env")
	envData, err := os.ReadFile(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("gateway token not found in openclaw.json or .env")
		}
		return "", err
	}

	for _, line := range strings.Split(string(envData), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if ok && key == "OPENCLAW_GATEWAY_TOKEN" {
			value = strings.TrimSpace(value)
			if value != "" {
				return value, nil
			}
		}
	}

	return "", fmt.Errorf("gateway token not found in openclaw.json or .env")
}