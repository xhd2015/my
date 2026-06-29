package openclaw

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGrokEntryToAuthStore(t *testing.T) {
	entry := GrokAuthEntry{
		Key:          "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour).Format(time.RFC3339Nano),
		Email:        "user@example.com",
		UserID:       "user-123",
		AuthMode:     "oidc",
		OIDCClientID: xaiOAuthClientID,
		OIDCIssuer:   xaiOAuthIssuer,
	}

	store, err := grokEntryToAuthStore(entry)
	if err != nil {
		t.Fatalf("grokEntryToAuthStore() error = %v", err)
	}

	profile := store.Profiles[xaiDefaultProfile]
	if profile.Access != "access-token" {
		t.Fatalf("access = %q, want access-token", profile.Access)
	}
	if profile.Refresh != "refresh-token" {
		t.Fatalf("refresh = %q, want refresh-token", profile.Refresh)
	}
	if profile.Email != "user@example.com" {
		t.Fatalf("email = %q", profile.Email)
	}
	if profile.AccountID != "user-123" {
		t.Fatalf("accountId = %q", profile.AccountID)
	}
}

func TestWriteAuthProfileStore(t *testing.T) {
	dir := t.TempDir()
	entry := GrokAuthEntry{
		Key:          "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour).Format(time.RFC3339Nano),
		AuthMode:     "oidc",
	}
	store, err := grokEntryToAuthStore(entry)
	if err != nil {
		t.Fatalf("grokEntryToAuthStore() error = %v", err)
	}
	if err := writeAuthProfileStore(dir, store); err != nil {
		t.Fatalf("writeAuthProfileStore() error = %v", err)
	}
	if _, err := os.Stat(agentAuthDBPath(dir)); err != nil {
		t.Fatalf("auth db missing: %v", err)
	}
}

func TestLoadGrokAuthEntryFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	content := `{
  "https://auth.x.ai::client": {
    "auth_mode": "oidc",
    "key": "jwt-access",
    "refresh_token": "jwt-refresh",
    "expires_at": "2026-12-31T00:00:00.000000Z",
    "email": "user@example.com",
    "user_id": "abc"
  }
}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	entry, err := LoadGrokAuthEntry(path)
	if err != nil {
		t.Fatalf("LoadGrokAuthEntry() error = %v", err)
	}
	if entry.Key != "jwt-access" {
		t.Fatalf("key = %q", entry.Key)
	}
}
