package openclaw

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	xaiOAuthClientID  = "b1a00492-073a-47ea-816f-4c329264a828"
	xaiOAuthIssuer    = "https://auth.x.ai"
	xaiTokenEndpoint  = "https://auth.x.ai/oauth2/token"
	xaiDefaultModel   = "xai/grok-4"
	xaiDefaultProfile = "xai:default"
)

type grokAuthEntry struct {
	Key          string `json:"key"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
	Email        string `json:"email"`
	UserID       string `json:"user_id"`
	AuthMode     string `json:"auth_mode"`
	OIDCClientID string `json:"oidc_client_id"`
	OIDCIssuer   string `json:"oidc_issuer"`
}

type openclawAuthProfileStore struct {
	Version  int                                      `json:"version"`
	Profiles map[string]openclawAuthProfileCredential `json:"profiles"`
	Order    map[string][]string                      `json:"order,omitempty"`
	LastGood map[string]string                        `json:"lastGood,omitempty"`
}

type openclawAuthProfileCredential struct {
	Type          string `json:"type"`
	Provider      string `json:"provider"`
	Access        string `json:"access"`
	Refresh       string `json:"refresh"`
	Expires       int64  `json:"expires"`
	Email         string `json:"email,omitempty"`
	AccountID     string `json:"accountId,omitempty"`
	Issuer        string `json:"issuer,omitempty"`
	TokenEndpoint string `json:"tokenEndpoint,omitempty"`
	ClientID      string `json:"clientId,omitempty"`
}

func resolveGrokAuthPath() (string, error) {
	if path := os.Getenv("MY_GROK_AUTH_PATH"); path != "" {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".grok", "auth.json"), nil
}

func loadGrokAuthEntry(path string) (grokAuthEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return grokAuthEntry{}, fmt.Errorf("grok auth file not found: %s", path)
		}
		return grokAuthEntry{}, err
	}

	var raw map[string]grokAuthEntry
	if err := json.Unmarshal(data, &raw); err != nil {
		return grokAuthEntry{}, fmt.Errorf("parse grok auth: %w", err)
	}
	if len(raw) == 0 {
		return grokAuthEntry{}, fmt.Errorf("grok auth file has no entries")
	}

	for _, entry := range raw {
		if entry.AuthMode == "oidc" || entry.Key != "" {
			return entry, nil
		}
	}
	return grokAuthEntry{}, fmt.Errorf("grok auth file has no OIDC entry")
}

func grokEntryToAuthStore(entry grokAuthEntry) (openclawAuthProfileStore, error) {
	if entry.Key == "" {
		return openclawAuthProfileStore{}, fmt.Errorf("grok auth entry is missing access token")
	}
	if entry.RefreshToken == "" {
		return openclawAuthProfileStore{}, fmt.Errorf("grok auth entry is missing refresh token")
	}

	expires, err := time.Parse(time.RFC3339Nano, entry.ExpiresAt)
	if err != nil {
		return openclawAuthProfileStore{}, fmt.Errorf("parse grok token expiry: %w", err)
	}

	issuer := entry.OIDCIssuer
	if issuer == "" {
		issuer = xaiOAuthIssuer
	}
	clientID := entry.OIDCClientID
	if clientID == "" {
		clientID = xaiOAuthClientID
	}

	return openclawAuthProfileStore{
		Version: 1,
		Profiles: map[string]openclawAuthProfileCredential{
			xaiDefaultProfile: {
				Type:          "oauth",
				Provider:      "xai",
				Access:        entry.Key,
				Refresh:       entry.RefreshToken,
				Expires:       expires.UnixMilli(),
				Email:         entry.Email,
				AccountID:     entry.UserID,
				Issuer:        issuer,
				TokenEndpoint: xaiTokenEndpoint,
				ClientID:      clientID,
			},
		},
		Order:    map[string][]string{"xai": {xaiDefaultProfile}},
		LastGood: map[string]string{"xai": xaiDefaultProfile},
	}, nil
}