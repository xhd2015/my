package openclaw

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TokenSources struct {
	JSON string
	Env  string
}

func (s TokenSources) Effective() string {
	if token := strings.TrimSpace(s.JSON); token != "" {
		return token
	}
	return strings.TrimSpace(s.Env)
}

func resolveTokenSources(dataDir string) (TokenSources, error) {
	var sources TokenSources

	configPath := filepath.Join(dataDir, "openclaw.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return sources, fmt.Errorf("openclaw.json not found in %s", dataDir)
		}
		return sources, err
	}

	var cfg openclawConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return sources, fmt.Errorf("parse openclaw.json: %w", err)
	}
	sources.JSON = strings.TrimSpace(cfg.Gateway.Auth.Token)

	envPath := filepath.Join(dataDir, ".env")
	envData, err := os.ReadFile(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return sources, nil
		}
		return sources, err
	}

	for _, line := range strings.Split(string(envData), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if ok && key == "OPENCLAW_GATEWAY_TOKEN" {
			sources.Env = strings.TrimSpace(value)
			break
		}
	}

	return sources, nil
}

func printTokenSources(dataDir, port, containerName string, sources TokenSources) {
	fmt.Printf("Data dir: %s\n", dataDir)
	if sources.JSON != "" {
		fmt.Printf("openclaw.json gateway.auth.token: %s\n", sources.JSON)
	} else {
		fmt.Println("openclaw.json gateway.auth.token: (not set)")
	}
	if sources.Env != "" {
		fmt.Printf(".env OPENCLAW_GATEWAY_TOKEN: %s\n", sources.Env)
	} else {
		fmt.Println(".env OPENCLAW_GATEWAY_TOKEN: (not set)")
	}

	effective := sources.Effective()
	if effective == "" {
		fmt.Println("Effective token: (none)")
		return
	}

	fmt.Printf("Effective token: %s\n", effective)
	base := fmt.Sprintf("http://127.0.0.1:%s", port)
	fmt.Printf("Auth dashboard: %s/#token=%s\n", base, effective)
	fmt.Printf("Auth chat:      %s/chat?session=main#token=%s\n", base, effective)

	running, err := containerIsRunning(containerName)
	if err == nil && running {
		if containerToken, err := containerEnvToken(containerName); err == nil && containerToken != "" {
			fmt.Printf("Container OPENCLAW_GATEWAY_TOKEN: %s\n", containerToken)
			if containerToken != effective {
				fmt.Fprintf(os.Stderr, "warning: container token does not match effective data-dir token\n")
			}
		}
	}
}