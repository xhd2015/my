package opencode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xhd2015/my/internal/openclaw"
)

func ImportLocalGrok() int {
	grokAuthPath, err := openclaw.ResolveGrokAuthPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	entry, err := openclaw.LoadGrokAuthEntry(grokAuthPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if entry.Key == "" {
		fmt.Fprintf(os.Stderr, "error: grok auth entry is missing access token\n")
		return 1
	}
	if entry.RefreshToken == "" {
		fmt.Fprintf(os.Stderr, "error: grok auth entry is missing refresh token\n")
		return 1
	}

	expires, err := time.Parse(time.RFC3339Nano, entry.ExpiresAt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: parse grok token expiry: %v\n", err)
		return 1
	}

	dataDir := resolveOpenCodeDataDir()
	authPath := filepath.Join(dataDir, "auth.json")

	authMap := make(map[string]json.RawMessage)
	if existing, err := os.ReadFile(authPath); err == nil {
		if err := json.Unmarshal(existing, &authMap); err != nil {
			fmt.Fprintf(os.Stderr, "error: parse existing auth.json: %v\n", err)
			return 1
		}
	}

	xaiEntry := map[string]any{
		"type":    "oauth",
		"access":  entry.Key,
		"refresh": entry.RefreshToken,
		"expires": expires.UnixMilli(),
	}
	xaiJSON, err := json.Marshal(xaiEntry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: marshal xai entry: %v\n", err)
		return 1
	}
	authMap["xai"] = xaiJSON

	output, err := json.MarshalIndent(authMap, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: marshal auth.json: %v\n", err)
		return 1
	}
	output = append(output, '\n')

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error: create data dir: %v\n", err)
		return 1
	}
	if err := os.WriteFile(authPath, output, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "error: write auth.json: %v\n", err)
		return 1
	}

	fmt.Printf("Imported xai credentials to %s\n", authPath)
	return 0
}

func resolveOpenCodeDataDir() string {
	if dir := os.Getenv("MY_OPENCODE_DATA_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "opencode")
}
