package openclaw

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestPrepareSQLiteFileForExportQuarantinesCorruptPluginDB(t *testing.T) {
	dir := t.TempDir()
	stateDir := filepath.Join(dir, "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatal(err)
	}
	dbPath := filepath.Join(stateDir, "openclaw.sqlite")
	if err := os.WriteFile(dbPath, []byte("not-a-sqlite-db"), 0o600); err != nil {
		t.Fatal(err)
	}

	warn, err := prepareSQLiteFileForExport(dbPath)
	if err != nil {
		t.Fatalf("prepareSQLiteFileForExport() error = %v", err)
	}
	if warn == "" {
		t.Fatal("expected warning about corrupt plugin db")
	}
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Fatalf("expected corrupt db to be quarantined, stat err = %v", err)
	}
	if _, err := os.Stat(dbPath + ".corrupt.bak"); err != nil {
		t.Fatalf("expected corrupt backup, err = %v", err)
	}
}

func TestPrepareSQLiteFileForExportCheckpointsWAL(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.sqlite")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE t (id INTEGER PRIMARY KEY, v TEXT)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO t (v) VALUES ('hello')`); err != nil {
		t.Fatal(err)
	}
	db.Close()

	walPath := dbPath + "-wal"
	if err := os.WriteFile(walPath, []byte("wal"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := prepareSQLiteFileForExport(dbPath); err != nil {
		t.Fatalf("prepareSQLiteFileForExport() error = %v", err)
	}
	if _, err := os.Stat(walPath); !os.IsNotExist(err) {
		t.Fatalf("expected wal sidecar removed, err = %v", err)
	}
}