package openclaw

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type gatewayBookkeeping struct {
	PID       int    `json:"pid"`
	Port      int    `json:"port"`
	StartedAt string `json:"started_at"`
	Kind      string `json:"kind"`
}

func gatewayBookkeepingPath(dataDir string) string {
	return filepath.Join(dataDir, ".my", "gateway.json")
}

func readGatewayBookkeeping(dataDir string) (*gatewayBookkeeping, error) {
	data, err := os.ReadFile(gatewayBookkeepingPath(dataDir))
	if err != nil {
		return nil, err
	}
	var rec gatewayBookkeeping
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func writeGatewayBookkeeping(dataDir string, pid, port int) error {
	myDir := filepath.Join(dataDir, ".my")
	if err := os.MkdirAll(myDir, 0o755); err != nil {
		return err
	}
	body, err := json.Marshal(gatewayBookkeeping{
		PID:       pid,
		Port:      port,
		StartedAt: time.Now().Format(time.RFC3339),
		Kind:      "local",
	})
	if err != nil {
		return err
	}
	return os.WriteFile(gatewayBookkeepingPath(dataDir), body, 0o644)
}

func removeGatewayBookkeeping(dataDir string) error {
	err := os.Remove(gatewayBookkeepingPath(dataDir))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func localGatewayRunning(dataDir string) (bool, error) {
	if dataDir == "" {
		return false, nil
	}

	absPath, err := resolvePath(dataDir)
	if err != nil {
		return false, err
	}

	rec, err := readGatewayBookkeeping(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if !isProcessAlive(rec.PID) {
		_ = removeGatewayBookkeeping(absPath)
		return false, nil
	}

	if isPortAvailable(rec.Port) {
		_ = removeGatewayBookkeeping(absPath)
		return false, nil
	}

	return true, nil
}

func resolveLocalGatewayPortFromBookkeeping(dataDir, port string) (string, error) {
	if port != "" {
		return port, nil
	}

	absPath, err := resolvePath(dataDir)
	if err != nil {
		return "", err
	}

	rec, err := readGatewayBookkeeping(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return resolveLocalGatewayPort(port), nil
		}
		return "", err
	}

	return fmt.Sprintf("%d", rec.Port), nil
}