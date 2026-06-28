package openclaw

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestZipOpenClawDataDir(t *testing.T) {
	source := t.TempDir()
	if err := os.WriteFile(filepath.Join(source, "openclaw.json"), []byte(`{"gateway":{}}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(source, "agents", "main"), 0o755); err != nil {
		t.Fatalf("mkdir agents: %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, ".env"), []byte("OPENCLAW_GATEWAY_TOKEN=test\n"), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "openclaw-settings.zip")
	count, err := zipOpenClawDataDir(source, zipPath)
	if err != nil {
		t.Fatalf("zipOpenClawDataDir() error = %v", err)
	}
	if count != 2 {
		t.Fatalf("file count = %d, want 2", count)
	}

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("OpenReader() error = %v", err)
	}
	defer reader.Close()

	got := make(map[string]string)
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "/") {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			t.Fatalf("open %s: %v", file.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("read %s: %v", file.Name, err)
		}
		got[file.Name] = string(data)
	}

	want := map[string]string{
		"openclaw/openclaw.json": `{"gateway":{}}`,
		"openclaw/.env":          "OPENCLAW_GATEWAY_TOKEN=test\n",
	}
	for path, content := range want {
		if got[path] != content {
			t.Fatalf("%s = %q, want %q", path, got[path], content)
		}
	}
}