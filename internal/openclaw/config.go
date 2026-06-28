package openclaw

import (
	"os"
	"path/filepath"
)

func configDir() (string, error) {
	if dir := os.Getenv("MY_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "my"), nil
}

func registryPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "openclaw.json"), nil
}