package openclaw

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiagnoseSlackConfigSocketModeOK(t *testing.T) {
	dir := t.TempDir()
	config := `{
  "channels": {
    "slack": {
      "enabled": true,
      "mode": "socket",
      "botToken": "xoxb-test",
      "appToken": "xapp-test",
      "testTarget": "C123TEST"
    }
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "openclaw.json"), []byte(config), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	view, checks, err := diagnoseSlackConfig(dir)
	if err != nil {
		t.Fatalf("diagnoseSlackConfig() error = %v", err)
	}
	if view.mode != "socket" {
		t.Fatalf("mode = %q, want socket", view.mode)
	}
	if view.testTarget != "C123TEST" {
		t.Fatalf("testTarget = %q", view.testTarget)
	}
	for _, check := range checks {
		if check.missing {
			t.Fatalf("unexpected missing check: %s (%s)", check.name, check.detail)
		}
	}
}

func TestDiagnoseSlackConfigEnvSecretRef(t *testing.T) {
	dir := t.TempDir()
	config := `{
  "channels": {
    "slack": {
      "enabled": true,
      "botToken": { "source": "env", "id": "SLACK_BOT_TOKEN" },
      "appToken": { "source": "env", "id": "SLACK_APP_TOKEN" }
    }
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "openclaw.json"), []byte(config), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("SLACK_BOT_TOKEN=xoxb-env\nSLACK_APP_TOKEN=xapp-env\n"), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	_, checks, err := diagnoseSlackConfig(dir)
	if err != nil {
		t.Fatalf("diagnoseSlackConfig() error = %v", err)
	}
	for _, check := range checks {
		if check.missing {
			t.Fatalf("unexpected missing check: %s (%s)", check.name, check.detail)
		}
	}
}

func TestDiagnoseSlackConfigMissingTokens(t *testing.T) {
	dir := t.TempDir()
	config := `{
  "channels": {
    "slack": {
      "enabled": true
    }
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "openclaw.json"), []byte(config), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, checks, err := diagnoseSlackConfig(dir)
	if err != nil {
		t.Fatalf("diagnoseSlackConfig() error = %v", err)
	}
	missing := 0
	for _, check := range checks {
		if check.missing {
			missing++
		}
	}
	if missing < 2 {
		t.Fatalf("expected missing token checks, got %d missing checks", missing)
	}
}