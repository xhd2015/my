package openclaw

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func isSQLiteSidecar(name string) bool {
	return strings.HasSuffix(name, ".sqlite-wal") || strings.HasSuffix(name, ".sqlite-shm")
}

func prepareSQLiteFilesForExport(dataDir string) ([]string, error) {
	var warnings []string
	err := filepath.WalkDir(dataDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sqlite") {
			return nil
		}
		warn, err := prepareSQLiteFileForExport(path)
		if err != nil {
			return err
		}
		if warn != "" {
			warnings = append(warnings, warn)
		}
		return nil
	})
	if err != nil {
		return warnings, err
	}
	for _, warn := range warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warn)
	}
	return warnings, nil
}

func prepareSQLiteFileForExport(path string) (string, error) {
	if err := checkpointSQLiteFile(path); err != nil {
		if isRegenerablePluginStateDB(path) {
			return quarantineCorruptSQLite(path)
		}
		return "", fmt.Errorf("checkpoint %s: %w", path, err)
	}

	ok, err := sqliteIntegrityOK(path)
	if err != nil {
		if isRegenerablePluginStateDB(path) {
			return quarantineCorruptSQLite(path)
		}
		return "", err
	}
	if ok {
		removeSQLiteSidecars(path)
		return "", nil
	}

	if isRegenerablePluginStateDB(path) {
		return quarantineCorruptSQLite(path)
	}

	return "", fmt.Errorf("sqlite integrity check failed for %s (export aborted)", path)
}

func quarantineCorruptSQLite(path string) (string, error) {
	backup := path + ".corrupt.bak"
	if err := os.Rename(path, backup); err != nil {
		return "", fmt.Errorf("quarantine corrupt plugin db %s: %w", path, err)
	}
	removeSQLiteSidecars(path)
	return fmt.Sprintf("removed corrupt plugin state database %s (backed up to %s; OpenClaw will recreate it)", path, backup), nil
}

func isRegenerablePluginStateDB(path string) bool {
	base := filepath.Base(path)
	dir := filepath.Base(filepath.Dir(path))
	return dir == "state" && base == "openclaw.sqlite"
}

func checkpointSQLiteFile(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(`PRAGMA wal_checkpoint(TRUNCATE)`); err != nil {
		return err
	}
	return nil
}

func sqliteIntegrityOK(path string) (bool, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var result string
	if err := db.QueryRow(`PRAGMA integrity_check`).Scan(&result); err != nil {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(result), "ok"), nil
}

func removeSQLiteSidecars(sqlitePath string) {
	_ = os.Remove(sqlitePath + "-wal")
	_ = os.Remove(sqlitePath + "-shm")
}